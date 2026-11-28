// Package healthcheck provides team health check data operations
package healthcheck

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// Service represents a PostgreSQL implementation of the health check data service
type Service struct {
	DB *sql.DB
}

const healthCheckColumns = `hc.id, hc.team_id, hc.name, hc.areas, hc.deadline, hc.deadline <= CURRENT_TIMESTAMP,
	COALESCE(hc.created_by::text, ''),
	(SELECT count(*) FROM thunderdome.health_check_response r WHERE r.health_check_id = hc.id),
	hc.created_at, hc.updated_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanHealthCheck(row rowScanner) (*thunderdome.HealthCheck, error) {
	hc := &thunderdome.HealthCheck{}
	var areas string

	if err := row.Scan(
		&hc.ID, &hc.TeamID, &hc.Name, &areas, &hc.Deadline, &hc.Closed, &hc.CreatedBy,
		&hc.ResponseCount, &hc.CreatedAt, &hc.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(areas), &hc.Areas); err != nil {
		return nil, fmt.Errorf("health check areas unmarshal error: %v", err)
	}

	return hc, nil
}

// HealthCheckList gets a team's health checks, newest first
func (d *Service) HealthCheckList(ctx context.Context, teamID string, limit int, offset int) ([]*thunderdome.HealthCheck, int, error) {
	checks := make([]*thunderdome.HealthCheck, 0)
	var count int

	if err := d.DB.QueryRowContext(ctx,
		`SELECT count(*) FROM thunderdome.health_check WHERE team_id = $1;`,
		teamID,
	).Scan(&count); err != nil {
		return nil, 0, fmt.Errorf("health check list count query error: %v", err)
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT `+healthCheckColumns+`
		FROM thunderdome.health_check hc
		WHERE hc.team_id = $1
		ORDER BY hc.created_at DESC
		LIMIT $2 OFFSET $3;`,
		teamID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("health check list query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		hc, err := scanHealthCheck(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("health check list scan error: %v", err)
		}
		checks = append(checks, hc)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("health check list rows error: %v", err)
	}

	return checks, count, nil
}

// HealthCheckGet gets one of a team's health checks
func (d *Service) HealthCheckGet(ctx context.Context, teamID string, healthCheckID string) (*thunderdome.HealthCheck, error) {
	hc, err := scanHealthCheck(d.DB.QueryRowContext(ctx,
		`SELECT `+healthCheckColumns+`
		FROM thunderdome.health_check hc
		WHERE hc.team_id = $1 AND hc.id = $2;`,
		teamID, healthCheckID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, thunderdome.ErrHealthCheckNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("health check get query error: %v", err)
	}

	return hc, nil
}

// HealthCheckCreate creates a health check with the areas of a template available to the team,
// that is a global template or one belonging to the team or its organization
func (d *Service) HealthCheckCreate(ctx context.Context, teamID string, templateID string, name string, deadline time.Time, createdBy string) (string, error) {
	var id string

	err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.health_check (team_id, name, areas, deadline, created_by)
		SELECT $1, $3, ht.areas, $4, $5
		FROM thunderdome.health_check_template ht
		WHERE ht.id = $2 AND `+availableTemplateCondition+`
		RETURNING id;`,
		teamID, templateID, name, deadline, createdBy,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", thunderdome.ErrHealthCheckTemplateNotFound
	}
	if err != nil {
		return "", fmt.Errorf("health check create query error: %v", err)
	}

	return id, nil
}

// HealthCheckUpdate renames an open health check and changes its deadline
func (d *Service) HealthCheckUpdate(ctx context.Context, teamID string, healthCheckID string, name string, deadline time.Time) error {
	return d.updateOpenHealthCheck(ctx, "health check update",
		`UPDATE thunderdome.health_check
		SET name = $3, deadline = $4, updated_at = CURRENT_TIMESTAMP
		WHERE team_id = $1 AND id = $2 AND deadline > CURRENT_TIMESTAMP;`,
		teamID, healthCheckID, name, deadline,
	)
}

// HealthCheckClose ends an open health check now
func (d *Service) HealthCheckClose(ctx context.Context, teamID string, healthCheckID string) error {
	return d.updateOpenHealthCheck(ctx, "health check close",
		`UPDATE thunderdome.health_check
		SET deadline = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE team_id = $1 AND id = $2 AND deadline > CURRENT_TIMESTAMP;`,
		teamID, healthCheckID,
	)
}

// updateOpenHealthCheck runs an update that only matches open health checks. Closed ones are never changed,
// reopening them would let results be compared before and after a response to reveal who gave it.
func (d *Service) updateOpenHealthCheck(ctx context.Context, action string, query string, args ...any) error {
	result, err := d.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s query error: %v", action, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s rows affected error: %v", action, err)
	}
	if count == 0 {
		return thunderdome.ErrHealthCheckClosed
	}

	return nil
}

// HealthCheckDelete removes a health check and its responses
func (d *Service) HealthCheckDelete(ctx context.Context, teamID string, healthCheckID string) error {
	result, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.health_check WHERE team_id = $1 AND id = $2;`,
		teamID, healthCheckID,
	)
	if err != nil {
		return fmt.Errorf("health check delete query error: %v", err)
	}

	if count, err := result.RowsAffected(); err == nil && count == 0 {
		return thunderdome.ErrHealthCheckNotFound
	}

	return nil
}

// HealthCheckResponseSubmit saves or replaces a team member's answers while the health check is open
func (d *Service) HealthCheckResponseSubmit(ctx context.Context, teamID string, healthCheckID string, userID string, answers []*thunderdome.HealthCheckAnswer) error {
	data, err := json.Marshal(answers)
	if err != nil {
		return fmt.Errorf("health check response marshal error: %v", err)
	}

	return d.updateOpenHealthCheck(ctx, "health check response submit",
		`INSERT INTO thunderdome.health_check_response (health_check_id, user_id, answers)
		SELECT hc.id, tu.user_id, $4::jsonb
		FROM thunderdome.health_check hc
		JOIN thunderdome.team_user tu ON tu.team_id = hc.team_id AND tu.user_id = $3
		WHERE hc.team_id = $1 AND hc.id = $2 AND hc.deadline > CURRENT_TIMESTAMP
		ON CONFLICT (health_check_id, user_id)
		DO UPDATE SET answers = EXCLUDED.answers, updated_at = CURRENT_TIMESTAMP;`,
		teamID, healthCheckID, userID, string(data),
	)
}

// HealthCheckUserResponse gets a user's own answers, empty if they haven't responded
func (d *Service) HealthCheckUserResponse(ctx context.Context, healthCheckID string, userID string) ([]*thunderdome.HealthCheckAnswer, error) {
	answers := make([]*thunderdome.HealthCheckAnswer, 0)
	var data string

	err := d.DB.QueryRowContext(ctx,
		`SELECT answers FROM thunderdome.health_check_response WHERE health_check_id = $1 AND user_id = $2;`,
		healthCheckID, userID,
	).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return answers, nil
	}
	if err != nil {
		return nil, fmt.Errorf("health check user response query error: %v", err)
	}

	if err := json.Unmarshal([]byte(data), &answers); err != nil {
		return nil, fmt.Errorf("health check user response unmarshal error: %v", err)
	}

	return answers, nil
}

// HealthCheckResponses gets every response without saying who gave it
func (d *Service) HealthCheckResponses(ctx context.Context, healthCheckID string) ([][]*thunderdome.HealthCheckAnswer, error) {
	responses := make([][]*thunderdome.HealthCheckAnswer, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT answers FROM thunderdome.health_check_response WHERE health_check_id = $1;`,
		healthCheckID,
	)
	if err != nil {
		return nil, fmt.Errorf("health check responses query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("health check responses scan error: %v", err)
		}

		var answers []*thunderdome.HealthCheckAnswer
		if err := json.Unmarshal([]byte(data), &answers); err != nil {
			return nil, fmt.Errorf("health check responses unmarshal error: %v", err)
		}
		responses = append(responses, answers)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("health check responses rows error: %v", err)
	}

	return responses, nil
}

// HealthCheckPendingUsers gets the team members who haven't responded
func (d *Service) HealthCheckPendingUsers(ctx context.Context, teamID string, healthCheckID string) ([]*thunderdome.HealthCheckParticipant, error) {
	users := make([]*thunderdome.HealthCheckParticipant, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT u.id, u.name
		FROM thunderdome.team_user tu
		JOIN thunderdome.users u ON u.id = tu.user_id
		WHERE tu.team_id = $1 AND NOT EXISTS (
			SELECT 1 FROM thunderdome.health_check_response r
			WHERE r.health_check_id = $2 AND r.user_id = tu.user_id
		)
		ORDER BY lower(u.name);`,
		teamID, healthCheckID,
	)
	if err != nil {
		return nil, fmt.Errorf("health check pending users query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &thunderdome.HealthCheckParticipant{}
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, fmt.Errorf("health check pending users scan error: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("health check pending users rows error: %v", err)
	}

	return users, nil
}
