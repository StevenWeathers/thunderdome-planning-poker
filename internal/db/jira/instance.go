package jira

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// FindInstancesByUserID returns all JiraInstances for a given user ID.
func (s *Service) FindInstancesByUserID(ctx context.Context, userID string) ([]thunderdome.JiraInstance, error) {
	instances := make([]thunderdome.JiraInstance, 0)

	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, user_id, host, client_mail, access_token, jira_data_center, created_date, updated_date
 				FROM thunderdome.jira_instance WHERE user_id = $1;`,
		userID,
	)
	if err != nil {
		return instances, fmt.Errorf("find jira instance by user id query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		instance := thunderdome.JiraInstance{}
		if err := rows.Scan(
			&instance.ID, &instance.UserID, &instance.Host, &instance.ClientMail, &instance.AccessToken, &instance.JiraDataCenter,
			&instance.CreatedDate, &instance.UpdatedDate,
		); err != nil {
			return instances, fmt.Errorf("find jira instance by user id row scan error: %v", err)
		}
		instance.AccessToken, err = db.Decrypt(instance.AccessToken, s.AESHashKey)
		if err != nil {
			return instances, fmt.Errorf("error decrypting jira_instance %s access_token:  %v", instance.ID, err)
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// GetInstanceByID returns a JiraInstance for a given instance ID.
func (s *Service) GetInstanceByID(ctx context.Context, instanceID string) (thunderdome.JiraInstance, error) {
	instance := thunderdome.JiraInstance{}

	err := s.DB.QueryRowContext(ctx,
		`SELECT id, user_id, host, client_mail, access_token, jira_data_center, created_date, updated_date
 				FROM thunderdome.jira_instance WHERE id = $1;`,
		instanceID,
	).Scan(
		&instance.ID, &instance.UserID, &instance.Host, &instance.ClientMail, &instance.AccessToken, &instance.JiraDataCenter,
		&instance.CreatedDate, &instance.UpdatedDate,
	)
	if err != nil {
		return instance, fmt.Errorf("error encountered getting jira_instance %s:  %v", instanceID, err)
	}
	instance.AccessToken, err = db.Decrypt(instance.AccessToken, s.AESHashKey)
	if err != nil {
		return instance, fmt.Errorf("error decrypting jira_instance %s access_token:  %v", instanceID, err)
	}

	return instance, nil
}

// CreateInstance creates a new JiraInstance.
func (s *Service) CreateInstance(ctx context.Context, userID string, host string, clientMail string, accessToken string, jiraDataCenter bool) (thunderdome.JiraInstance, error) {
	instance := thunderdome.JiraInstance{}
	secureToken, err := db.Encrypt(accessToken, s.AESHashKey)
	if err != nil {
		return instance, fmt.Errorf("error encountered creating jira_instance:  %v", err)
	}

	err = s.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.jira_instance
				(user_id, host, client_mail, access_token, jira_data_center)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, user_id, host, client_mail, access_token, jira_data_center, created_date, updated_date;`,
		userID, host, clientMail, secureToken, jiraDataCenter,
	).Scan(
		&instance.ID, &instance.UserID, &instance.Host, &instance.ClientMail, &instance.AccessToken, &instance.JiraDataCenter,
		&instance.CreatedDate, &instance.UpdatedDate,
	)
	if err != nil {
		return instance, fmt.Errorf("error encountered creating jira_instance:  %v", err)
	}

	return instance, nil
}

// UpdateInstance updates an existing JiraInstance.
func (s *Service) UpdateInstance(ctx context.Context, instanceID string, host string, clientMail string, accessToken string) (thunderdome.JiraInstance, error) {
	instance := thunderdome.JiraInstance{}
	at, err := db.Encrypt(accessToken, s.AESHashKey)
	if err != nil {
		return instance, fmt.Errorf("error encountered updating jira_instance:  %v", err)
	}

	err = s.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.jira_instance
				SET host = $2, client_mail = $3, access_token = $4
				WHERE id = $1
				RETURNING id, user_id, host, client_mail, access_token, created_date, updated_date;`,
		instanceID, host, clientMail, at,
	).Scan(
		&instance.ID, &instance.UserID, &instance.Host, &instance.ClientMail, &instance.AccessToken,
		&instance.CreatedDate, &instance.UpdatedDate,
	)
	if err != nil {
		return instance, fmt.Errorf("error encountered updating jira_instance:  %v", err)
	}

	return instance, nil
}

// DeleteInstance deletes an existing JiraInstance.
func (s *Service) DeleteInstance(ctx context.Context, instanceID string) error {
	result, err := s.DB.ExecContext(ctx, `DELETE FROM thunderdome.jira_instance WHERE id = $1;`, instanceID)
	if err != nil {
		return fmt.Errorf("delete jira instance query error: %v", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete jira instance rows error: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("delete jira instance expected to affect 1 row, affected %d", rows)
	}

	return nil
}
