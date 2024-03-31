package subscription

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// CheckActiveSubscriber looks for an active subscription for the user
// also checks if the user belongs to a team or organization with an active subscription
func (s *Service) CheckActiveSubscriber(ctx context.Context, userId string) error {
	count := 0

	err := s.DB.QueryRowContext(ctx,
		`
				WITH user_teams AS (
					SELECT team_id FROM thunderdome.team_user WHERE user_id = $1
				),
				user_organizations AS (
					SELECT organization_id FROM thunderdome.organization_user WHERE user_id = $1
				)
				SELECT count(id)
 				FROM thunderdome.subscription 
 				WHERE (user_id = $1 AND active = true AND expires > NOW()) 
 				OR (team_id IN (SELECT team_id FROM user_teams) AND active = true AND expires > NOW()) 
 				OR (organization_id IN (SELECT organization_id FROM user_organizations) AND active = true AND expires > NOW());
				`,
		userId,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("error encountered finding users %s active subscriptions:  %v", userId, err)
	}

	if count == 0 {
		return fmt.Errorf("no active subscription(s) found for user id  %s", userId)
	}

	return nil
}

func (s *Service) GetActiveSubscriptionsByUserID(ctx context.Context, userId string) ([]thunderdome.Subscription, error) {
	subs := make([]thunderdome.Subscription, 0)

	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, user_id, COALESCE(team_id::text, ''), COALESCE(organization_id::text, ''),
				customer_id, subscription_id, active, type, expires, created_date, updated_date
 				FROM thunderdome.subscription WHERE user_id = $1 AND active = true AND expires > NOW();`,
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
			&sub.TeamID,
			&sub.OrganizationID,
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
		`SELECT s.id, s.user_id, COALESCE(team_id::text, ''), COALESCE(organization_id::text, ''), s.customer_id, s.subscription_id, s.active, s.type, s.expires,
 				s.created_date, s.updated_date, u.id, u.name, u.email
 				FROM thunderdome.subscription s
 				LEFT JOIN thunderdome.users u ON s.user_id = u.id
 				WHERE s.id = $1;`,
		id,
	).Scan(
		&sub.ID, &sub.UserID, &sub.TeamID, &sub.OrganizationID, &sub.CustomerID, &sub.SubscriptionID,
		&sub.Active, &sub.Type, &sub.Expires, &sub.CreatedDate, &sub.UpdatedDate,
		&sub.User.Id, &sub.User.Name, &sub.User.Email,
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
		`SELECT id, user_id, COALESCE(team_id::text, ''), COALESCE(organization_id::text, ''), customer_id, subscription_id, active, type, expires, created_date, updated_date
 				FROM thunderdome.subscription WHERE subscription_id = $1;`,
		subscriptionId,
	).Scan(
		&sub.ID, &sub.UserID, &sub.TeamID, &sub.OrganizationID,
		&sub.CustomerID, &sub.SubscriptionID, &sub.Active, &sub.Type, &sub.Expires,
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

func (s *Service) CreateSubscription(ctx context.Context, subscription thunderdome.Subscription) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.subscription 
				(user_id, team_id, organization_id, customer_id, subscription_id, type, expires)
				VALUES ($1, NULLIF($2::text,'')::uuid, NULLIF($3::text,'')::uuid, $4, $5, $6, $7)
				RETURNING id, user_id, COALESCE(team_id::text, ''), COALESCE(organization_id::text, ''), customer_id, subscription_id, active, type, expires, created_date, updated_date;`,
		subscription.UserID, subscription.TeamID, subscription.OrganizationID,
		subscription.CustomerID, subscription.SubscriptionID, subscription.Type, subscription.Expires,
	).Scan(
		&sub.ID, &sub.UserID, &sub.TeamID, &sub.OrganizationID,
		&sub.CustomerID, &sub.SubscriptionID, &sub.Active, &sub.Type, &sub.Expires,
		&sub.CreatedDate, &sub.UpdatedDate,
	)
	if err != nil {
		return sub, fmt.Errorf("error encountered creating subscription: %v", err)
	}

	return sub, nil
}

func (s *Service) UpdateSubscription(ctx context.Context, id string, subscription thunderdome.Subscription) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.subscription SET customer_id = $2, subscription_id = $3, active = $4,
 				type = $5, expires = $6, team_id = NULLIF($7::text,'')::uuid, organization_id = NULLIF($8::text,'')::uuid,
 				updated_date = NOW()
 				WHERE id = $1
				RETURNING id, user_id, COALESCE(team_id::text, ''), COALESCE(organization_id::text, ''), customer_id, subscription_id, active, type, expires, created_date, updated_date;`,
		id, subscription.CustomerID, subscription.SubscriptionID, subscription.Active,
		subscription.Type, subscription.Expires, subscription.TeamID, subscription.OrganizationID,
	).Scan(
		&sub.ID, &sub.UserID, &sub.TeamID, &sub.OrganizationID,
		&sub.CustomerID, &sub.SubscriptionID, &sub.Active, &sub.Type, &sub.Expires,
		&sub.CreatedDate, &sub.UpdatedDate,
	)
	if err != nil {
		return sub, fmt.Errorf("error encountered updating subscription: %v", err)
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
		`SELECT s.id, s.user_id, COALESCE(s.team_id::text, ''), COALESCE(s.organization_id::text, ''),
 				s.customer_id, s.subscription_id, s.active, s.type, s.expires,
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
			&sub.TeamID,
			&sub.OrganizationID,
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
