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
		`SELECT id, user_id, customer_id, active, expires, created_date, updated_date
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
		`SELECT id, user_id, customer_id, active, expires, created_date, updated_date
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
			&sub.Active,
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
func (s *Service) GetSubscriptionByCustomerID(ctx context.Context, customerId string) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`SELECT id, user_id, customer_id, active, expires, created_date, updated_date
 				FROM thunderdome.subscription WHERE customer_id = $1;`,
		customerId,
	).Scan(
		&sub.ID, &sub.UserID, &sub.CustomerID, &sub.Active, &sub.Expires,
		&sub.CreatedDate, &sub.UpdatedDate,
	)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return sub, fmt.Errorf("no subscription found for customer id %s", customerId)
	case err != nil:
		return sub, fmt.Errorf("error encountered finding customer id subscription: %v", err)
	}

	return sub, nil
}
func (s *Service) CreateSubscription(ctx context.Context, userId string, customerId string, expires time.Time) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.subscription 
				(user_id, customer_id, expires)
				VALUES ($1, $2, $3)
				RETURNING id, user_id, customer_id, active, expires, created_date, updated_date;`,
		userId, customerId, expires,
	).Scan(
		&sub.ID, &sub.UserID, &sub.CustomerID, &sub.Active, &sub.Expires,
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
func (s *Service) UpdateSubscription(ctx context.Context, id string, active bool, expires time.Time) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.subscription SET active = $2, expires = $3, updated_date = NOW() WHERE id = $1
				RETURNING id, user_id, customer_id, active, expires, created_date, updated_date;`,
		id, active, expires,
	).Scan(
		&sub.ID, &sub.UserID, &sub.CustomerID, &sub.Active, &sub.Expires,
		&sub.CreatedDate, &sub.UpdatedDate,
	)
	if err != nil {
		return sub, fmt.Errorf("error encountered updating subscription: %v", err)
	}

	result, err := s.DB.ExecContext(ctx,
		`UPDATE thunderdome.users SET subscribed = $2, updated_date = NOW() WHERE id = $1;`,
		sub.UserID, active,
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
