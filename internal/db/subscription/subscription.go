package subscription

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// CheckActiveSubscriber looks for an active subscription for the user
func (s *Service) CheckActiveSubscriber(ctx context.Context, userId string) error {
	activeSub := false
	currentTime := time.Now()

	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, user_id, customer_id, subscription_id, active, expires, created_date, updated_date
 				FROM thunderdome.subscription WHERE user_id = $1 AND active = true;`,
		userId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no active subscription found for user id  %s", userId)
		}

		return fmt.Errorf("error encountered finding user %s active subscriptions:  %v", userId, err)
	}

	defer rows.Close()
	for rows.Next() {
		var sub thunderdome.Subscription
		if err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.CustomerID,
			&sub.SubscriptionID,
			&sub.Active,
			&sub.Expires,
			&sub.CreatedDate,
			&sub.UpdatedDate,
		); err != nil {
			return fmt.Errorf("error encountered finding user %s active subscriptions:  %v", userId, err)
		}
		if currentTime.After(sub.Expires) {
			_, updateErr := s.DB.ExecContext(ctx,
				`UPDATE thunderdome.users SET subscribed = false, updated_date = NOW() WHERE id = $1;`,
				sub.UserID,
			)
			if updateErr != nil {
				s.Logger.Ctx(ctx).Error(fmt.Sprintf("error updating user %s subscribed to false", userId),
					zap.Error(updateErr), zap.String("customer_id", sub.CustomerID))
			}
		} else {
			activeSub = true
		}
	}

	if !activeSub {
		return fmt.Errorf("no active subscription found for user id  %s", userId)
	}

	return nil
}

func (s *Service) GetSubscriptionsByUserID(ctx context.Context, userId string) ([]thunderdome.Subscription, error) {
	subs := make([]thunderdome.Subscription, 0)

	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, user_id, customer_id, subscription_id, active, type, expires, created_date, updated_date
 				FROM thunderdome.subscription WHERE user_id = $1;`,
		userId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return subs, nil
		}

		return subs, fmt.Errorf("error getting user %s subscriptions:  %v", userId, err)
	}

	defer rows.Close()
	for rows.Next() {
		var sub thunderdome.Subscription
		if err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.CustomerID,
			&sub.SubscriptionID,
			&sub.Active,
			&sub.Type,
			&sub.Expires,
			&sub.CreatedDate,
			&sub.UpdatedDate,
		); err != nil {
			return subs, fmt.Errorf("error getting user %s subscriptions:  %v", userId, err)
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func (s *Service) GetSubscriptionByID(ctx context.Context, id string) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`SELECT s.id, s.user_id, s.customer_id, s.subscription_id, s.active, s.type, s.expires,
 				s.created_date, s.updated_date, u.id, u.name, u.email
 				FROM thunderdome.subscription s
 				LEFT JOIN thunderdome.users u ON s.user_id = u.id
 				WHERE s.id = $1;`,
		id,
	).Scan(
		&sub.ID, &sub.UserID, &sub.CustomerID, &sub.SubscriptionID, &sub.Active, &sub.Type, &sub.Expires,
		&sub.CreatedDate, &sub.UpdatedDate, &sub.User.Id, &sub.User.Name, &sub.User.Email,
	)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return sub, fmt.Errorf("subscription not found %s", id)
	case err != nil:
		return sub, fmt.Errorf("error encountered finding subscription: %v", err)
	}

	return sub, nil
}

func (s *Service) GetSubscriptionBySubscriptionID(ctx context.Context, subscriptionId string) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`SELECT id, user_id, customer_id, subscription_id, active, type, expires, created_date, updated_date
 				FROM thunderdome.subscription WHERE subscription_id = $1;`,
		subscriptionId,
	).Scan(
		&sub.ID, &sub.UserID, &sub.CustomerID, &sub.SubscriptionID, &sub.Active, &sub.Type, &sub.Expires,
		&sub.CreatedDate, &sub.UpdatedDate,
	)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return sub, fmt.Errorf("no subscription %s", subscriptionId)
	case err != nil:
		return sub, fmt.Errorf("error encountered finding subscription: %v", err)
	}

	return sub, nil
}

func (s *Service) CreateSubscription(ctx context.Context, userId string, customerId string, subscriptionId string, subType string, expires time.Time) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.subscription 
				(user_id, customer_id, subscription_id, type, expires)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, user_id, customer_id, subscription_id, active, type, expires, created_date, updated_date;`,
		userId, customerId, subscriptionId, subType, expires,
	).Scan(
		&sub.ID, &sub.UserID, &sub.CustomerID, &sub.SubscriptionID, &sub.Active, &sub.Type, &sub.Expires,
		&sub.CreatedDate, &sub.UpdatedDate,
	)
	if err != nil {
		return sub, fmt.Errorf("error encountered creating subscription: %v", err)
	}

	result, err := s.DB.ExecContext(ctx,
		`UPDATE thunderdome.users SET subscribed = true, updated_date = NOW() WHERE id = $1;`,
		userId,
	)
	if err != nil {
		return sub, fmt.Errorf("error encountered updating user subscription status: %v", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return sub, fmt.Errorf("error encountered updating user subscription status: %v", err)
	}
	if rows != 1 {
		return sub, fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return sub, nil
}

func (s *Service) UpdateSubscription(ctx context.Context, id string, subscription thunderdome.Subscription) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.subscription SET customer_id = $2, subscription_id = $3, active = $4,
 				type = $5, expires = $6, updated_date = NOW() WHERE id = $1
				RETURNING id, user_id, customer_id, subscription_id, active, type, expires, created_date, updated_date;`,
		id, subscription.CustomerID, subscription.SubscriptionID, subscription.Active,
		subscription.Type, subscription.Expires,
	).Scan(
		&sub.ID, &sub.UserID, &sub.CustomerID, &sub.SubscriptionID, &sub.Active, &sub.Type, &sub.Expires,
		&sub.CreatedDate, &sub.UpdatedDate,
	)
	if err != nil {
		return sub, fmt.Errorf("error encountered updating subscription: %v", err)
	}

	stillActive := subscription.Active
	stillActiveErr := s.CheckActiveSubscriber(ctx, subscription.UserID)
	if stillActiveErr != nil {
		stillActive = false
	}
	result, err := s.DB.ExecContext(ctx,
		`UPDATE thunderdome.users SET subscribed = $2, updated_date = NOW() WHERE id = $1;`,
		sub.UserID, stillActive,
	)
	if err != nil {
		return sub, fmt.Errorf("error encountered updating user subscription status: %v", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return sub, fmt.Errorf("error encountered updating user subscription status: %v", err)
	}
	if rows != 1 {
		return sub, fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return sub, nil
}

func (s *Service) GetSubscriptions(ctx context.Context, Limit int, Offset int) ([]thunderdome.Subscription, int, error) {
	var count int

	e := s.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM thunderdome.subscription;",
	).Scan(
		&count,
	)
	if e != nil {
		s.Logger.Ctx(ctx).Error("GetSubscriptions query scan error", zap.Error(e))
	}

	subs := make([]thunderdome.Subscription, 0)

	rows, err := s.DB.QueryContext(ctx,
		`SELECT s.id, s.user_id, s.customer_id, s.subscription_id, s.active, s.type, s.expires,
 				s.created_date, s.updated_date, u.id, u.name, u.email
 				FROM thunderdome.subscription s
 				LEFT JOIN thunderdome.users u ON s.user_id = u.id 
 				ORDER BY s.created_date ASC LIMIT $1 OFFSET $2 ;`,
		Limit, Offset,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return subs, count, nil
		}

		return subs, count, fmt.Errorf("error getting subscriptions:  %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var sub thunderdome.Subscription
		if err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.CustomerID,
			&sub.SubscriptionID,
			&sub.Active,
			&sub.Type,
			&sub.Expires,
			&sub.CreatedDate,
			&sub.UpdatedDate,
			&sub.User.Id,
			&sub.User.Name,
			&sub.User.Email,
		); err != nil {
			return subs, count, fmt.Errorf("error getting subscriptions:  %v", err)
		}

		subs = append(subs, sub)
	}

	return subs, count, nil
}

func (s *Service) DeleteSubscription(ctx context.Context, id string) error {
	if _, err := s.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.subscription WHERE id = $1;`,
		id); err != nil {
		return fmt.Errorf("error deleting subscription: %v", err)
	}

	return nil
}
