package subscription

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

func (s *Service) GetSubscriptionByUserID(ctx context.Context, userId string) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`SELECT id, user_id, customer_id, active, expires, created_date, updated_date
 				FROM thunderdome.subscription WHERE user_id = $1;`,
		userId,
	).Scan(
		&sub.ID, &sub.UserID, &sub.CustomerID, &sub.Active, &sub.Expires,
		&sub.CreatedDate, &sub.UpdatedDate,
	)
	switch {
	case err == sql.ErrNoRows:
		return sub, fmt.Errorf("no subscription found for user id  %d", userId)
	case err != nil:
		return sub, fmt.Errorf("error encountered finding user id subscription:  %v", err)
	}

	return sub, nil
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
	case err == sql.ErrNoRows:
		return sub, fmt.Errorf("no subscription found for customer id  %s", customerId)
	case err != nil:
		return sub, fmt.Errorf("error encountered finding customer id subscription:  %v", err)
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
		return sub, fmt.Errorf("error encountered creating subscription:  %v", err)
	}

	return sub, nil
}
func (s *Service) UpdateSubscription(ctx context.Context, id string, active bool, expires time.Time) (thunderdome.Subscription, error) {
	sub := thunderdome.Subscription{}

	err := s.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.subscription SET active = $2, expires = $3 WHERE id = $1
				RETURNING id, user_id, customer_id, active, expires, created_date, updated_date;`,
		id, active, expires,
	).Scan(
		&sub.ID, &sub.UserID, &sub.CustomerID, &sub.Active, &sub.Expires,
		&sub.CreatedDate, &sub.UpdatedDate,
	)
	if err != nil {
		return sub, fmt.Errorf("error encountered updating subscription:  %v", err)
	}

	return sub, nil
}
