package jira

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

func (s *Service) FindInstancesByUserId(ctx context.Context, userId string) ([]thunderdome.JiraInstance, error) {
	instances := make([]thunderdome.JiraInstance, 0)

	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, user_id, host, client_mail, access_token, created_date, updated_date
 				FROM thunderdome.jira_instance WHERE user_id = $1;`,
		userId,
	)
	if err != nil {
		return instances, err
	}
	defer rows.Close()

	for rows.Next() {
		instance := thunderdome.JiraInstance{}
		if err := rows.Scan(
			&instance.ID, &instance.UserID, &instance.Host, &instance.ClientMail, &instance.AccessToken,
			&instance.CreatedDate, &instance.UpdatedDate,
		); err != nil {
			return instances, err
		}
		instance.AccessToken, err = db.Decrypt(instance.AccessToken, s.AESHashKey)
		if err != nil {
			return instances, fmt.Errorf("error decrypting jira_instance %s access_token:  %v", instance.ID, err)
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

func (s *Service) GetInstanceById(ctx context.Context, instanceId string) (thunderdome.JiraInstance, error) {
	instance := thunderdome.JiraInstance{}

	err := s.DB.QueryRowContext(ctx,
		`SELECT id, user_id, host, client_mail, access_token, created_date, updated_date
 				FROM thunderdome.jira_instance WHERE id = $1;`,
		instanceId,
	).Scan(
		&instance.ID, &instance.UserID, &instance.Host, &instance.ClientMail, &instance.AccessToken,
		&instance.CreatedDate, &instance.UpdatedDate,
	)
	if err != nil {
		return instance, fmt.Errorf("error encountered getting jira_instance %s:  %v", instanceId, err)
	}
	instance.AccessToken, err = db.Decrypt(instance.AccessToken, s.AESHashKey)
	if err != nil {
		return instance, fmt.Errorf("error decrypting jira_instance %s access_token:  %v", instanceId, err)
	}

	return instance, nil
}

func (s *Service) CreateInstance(ctx context.Context, userId string, host string, clientMail string, accessToken string) (thunderdome.JiraInstance, error) {
	instance := thunderdome.JiraInstance{}
	at, err := db.Encrypt(accessToken, s.AESHashKey)
	if err != nil {
		return instance, fmt.Errorf("error encountered creating jira_instance:  %v", err)
	}

	err = s.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.jira_instance 
				(user_id, host, client_mail, access_token)
				VALUES ($1, $2, $3, $4)
				RETURNING id, user_id, host, client_mail, access_token, created_date, updated_date;`,
		userId, host, clientMail, at,
	).Scan(
		&instance.ID, &instance.UserID, &instance.Host, &instance.ClientMail, &instance.AccessToken,
		&instance.CreatedDate, &instance.UpdatedDate,
	)
	if err != nil {
		return instance, fmt.Errorf("error encountered creating jira_instance:  %v", err)
	}

	return instance, nil
}

func (s *Service) UpdateInstance(ctx context.Context, instanceId string, host string, clientMail string, accessToken string) (thunderdome.JiraInstance, error) {
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
		instanceId, host, clientMail, at,
	).Scan(
		&instance.ID, &instance.UserID, &instance.Host, &instance.ClientMail, &instance.AccessToken,
		&instance.CreatedDate, &instance.UpdatedDate,
	)
	if err != nil {
		return instance, fmt.Errorf("error encountered updating jira_instance:  %v", err)
	}

	return instance, nil
}

func (s *Service) DeleteInstance(ctx context.Context, instanceId string) error {
	result, err := s.DB.ExecContext(ctx, `DELETE FROM thunderdome.jira_instance WHERE id = $1;`, instanceId)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}
