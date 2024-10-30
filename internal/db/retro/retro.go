package retro

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// Service represents the database service for retros
type Service struct {
	DB         *sql.DB
	Logger     *otelzap.Logger
	AESHashKey string
}

func (d *Service) CreateRetro(ctx context.Context, ownerID, teamID string, retroName, joinCode, facilitatorCode string, maxVotes int, brainstormVisibility string, phaseTimeLimitMin int, phaseAutoAdvance bool, allowCumulativeVoting bool, templateID string) (*thunderdome.Retro, error) {
	var encryptedFacilitatorCode string
	var encryptedJoinCode string
	var retro = &thunderdome.Retro{
		OwnerID:               ownerID,
		TeamID:                teamID,
		Name:                  retroName,
		Phase:                 "intro",
		PhaseTimeLimitMin:     phaseTimeLimitMin,
		PhaseAutoAdvance:      phaseAutoAdvance,
		Users:                 make([]*thunderdome.RetroUser, 0),
		Items:                 make([]*thunderdome.RetroItem, 0),
		ActionItems:           make([]*thunderdome.RetroAction, 0),
		BrainstormVisibility:  brainstormVisibility,
		MaxVotes:              maxVotes,
		TemplateID:            templateID,
		AllowCumulativeVoting: allowCumulativeVoting,
	}

	if joinCode != "" {
		encryptedCode, codeErr := db.Encrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create retro encrypt join code error: %v", codeErr)
		}
		encryptedJoinCode = encryptedCode
	}

	if facilitatorCode != "" {
		encryptedCode, codeErr := db.Encrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create retro encrypt facilitator code error: %v", codeErr)
		}
		encryptedFacilitatorCode = encryptedCode
	}

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		d.Logger.Error("create retro error", zap.Error(err))
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
		INSERT INTO thunderdome.retro (
			owner_id, team_id, name, join_code, facilitator_code,
			max_votes, brainstorm_visibility, phase_time_limit_min, phase_auto_advance,
			allow_cumulative_voting, template_id
		)
		VALUES ($1, NULLIF($2::text, '')::uuid, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_date, updated_date;
	`, ownerID, teamID, retroName, encryptedJoinCode, encryptedFacilitatorCode, maxVotes, brainstormVisibility,
		phaseTimeLimitMin, phaseAutoAdvance, allowCumulativeVoting, templateID).Scan(
		&retro.ID, &retro.CreatedDate, &retro.UpdatedDate,
	)

	if err != nil {
		d.Logger.Error("create retro error", zap.Error(err),
			zap.String("owner_id", ownerID), zap.String("name", retroName))
		return nil, fmt.Errorf("failed to insert into retro table: %v", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO thunderdome.retro_facilitator (retro_id, user_id)
		VALUES ($1, $2)
	`, retro.ID, ownerID)

	if err != nil {
		d.Logger.Error("create retro error", zap.Error(err))
		return nil, fmt.Errorf("failed to insert into retro_facilitator table: %v", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO thunderdome.retro_user (retro_id, user_id)
		VALUES ($1, $2)
	`, retro.ID, ownerID)

	if err != nil {
		d.Logger.Error("create retro error", zap.Error(err))
		return nil, fmt.Errorf("failed to insert into retro_user table: %v", err)
	}

	if err = tx.Commit(); err != nil {
		d.Logger.Error("create retro error", zap.Error(err))
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return retro, nil
}

// EditRetro updates the retro by ID
func (d *Service) EditRetro(retroID string, retroName string, joinCode string, facilitatorCode string, maxVotes int, brainstormVisibility string, phaseAutoAdvance bool) error {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if joinCode != "" {
		encryptedCode, codeErr := db.Encrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return fmt.Errorf("edit retro encrypt join code error: %v", codeErr)
		}
		encryptedJoinCode = encryptedCode
	}

	if facilitatorCode != "" {
		encryptedCode, codeErr := db.Encrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return fmt.Errorf("edit retro encrypt join facilitator error: %v", codeErr)
		}
		encryptedFacilitatorCode = encryptedCode
	}

	if _, err := d.DB.Exec(`UPDATE thunderdome.retro
    SET name = $2, join_code = $3, facilitator_code = $4, max_votes = $5,
        brainstorm_visibility = $6, phase_auto_advance = $7, updated_date = NOW()
    WHERE id = $1;`,
		retroID, retroName, encryptedJoinCode, encryptedFacilitatorCode,
		maxVotes, brainstormVisibility, phaseAutoAdvance,
	); err != nil {
		return fmt.Errorf("edit retro query error: %v", err)
	}

	return nil
}

// RetroGetByID gets a retro by ID
func (d *Service) RetroGetByID(retroID string, userID string) (*thunderdome.Retro, error) {
	var b = &thunderdome.Retro{
		ID:           retroID,
		Users:        make([]*thunderdome.RetroUser, 0),
		Items:        make([]*thunderdome.RetroItem, 0),
		Groups:       make([]*thunderdome.RetroGroup, 0),
		ActionItems:  make([]*thunderdome.RetroAction, 0),
		Votes:        make([]*thunderdome.RetroVote, 0),
		Facilitators: make([]string, 0),
		ReadyUsers:   make([]string, 0),
	}

	// get retro
	var joinCode string
	var facilitatorCode string
	var facilitators string
	var readyUsers string
	var template string
	err := d.DB.QueryRow(
		`SELECT
			r.id, r.name, r.owner_id, COALESCE(r.team_id::TEXT, ''), r.phase, r.phase_time_limit_min, r.phase_time_start, r.phase_auto_advance,
			 COALESCE(r.join_code, ''), COALESCE(r.facilitator_code, ''), r.allow_cumulative_voting,
			r.max_votes, r.brainstorm_visibility, r.ready_users, r.created_date, r.updated_date, r.template_id,
			CASE WHEN COUNT(rf) = 0 THEN '[]'::json ELSE array_to_json(array_agg(rf.user_id)) END AS facilitators,
			(SELECT row_to_json(t.*) as template FROM thunderdome.retro_template t WHERE t.id = r.template_id) AS template
		FROM thunderdome.retro r
		LEFT JOIN thunderdome.retro_facilitator rf ON r.id = rf.retro_id
		WHERE r.id = $1
		GROUP BY r.id`,
		retroID,
	).Scan(
		&b.ID,
		&b.Name,
		&b.OwnerID,
		&b.TeamID,
		&b.Phase,
		&b.PhaseTimeLimitMin,
		&b.PhaseTimeStart,
		&b.PhaseAutoAdvance,
		&joinCode,
		&facilitatorCode,
		&b.AllowCumulativeVoting,
		&b.MaxVotes,
		&b.BrainstormVisibility,
		&readyUsers,
		&b.CreatedDate,
		&b.UpdatedDate,
		&b.TemplateID,
		&facilitators,
		&template,
	)
	if err != nil {
		d.Logger.Error("get retro error", zap.Error(err))
		return nil, fmt.Errorf("get retro query error: %v", err)
	}

	facilError := json.Unmarshal([]byte(facilitators), &b.Facilitators)
	if facilError != nil {
		d.Logger.Error("facilitators json error", zap.Error(facilError))
	}
	isFacilitator := db.Contains(b.Facilitators, userID)

	if joinCode != "" {
		decryptedCode, codeErr := db.Decrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("get retro decrypt join code error: %v", codeErr)
		}
		b.JoinCode = decryptedCode
	}

	if facilitatorCode != "" && isFacilitator {
		decryptedCode, codeErr := db.Decrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("get retro decrypt join facilitator error: %v", codeErr)
		}
		b.FacilitatorCode = decryptedCode
	}

	readyUsersError := json.Unmarshal([]byte(readyUsers), &b.ReadyUsers)
	if readyUsersError != nil {
		d.Logger.Error("ready users json error", zap.Error(readyUsersError))
	}

	templateError := json.Unmarshal([]byte(template), &b.Template)
	if templateError != nil {
		d.Logger.Error("retro template json error", zap.Error(templateError))
		return nil, fmt.Errorf("get retro template error: %v", templateError)
	}

	b.Items = d.GetRetroItems(retroID)
	b.Groups = d.GetRetroGroups(retroID)
	b.Users = d.RetroGetUsers(retroID)
	b.ActionItems = d.GetRetroActions(retroID)
	b.Votes = d.GetRetroVotes(retroID)

	return b, nil
}

// RetroGetByUser gets a list of retros by UserID
func (d *Service) RetroGetByUser(userID string, limit int, offset int) ([]*thunderdome.Retro, int, error) {
	var retros = make([]*thunderdome.Retro, 0)
	var count int

	e := d.DB.QueryRow(`
		WITH user_teams AS (
			SELECT t.id FROM thunderdome.team_user tu
			LEFT JOIN thunderdome.team t ON t.id = tu.team_id
			WHERE tu.user_id = $1
		),
		team_retros AS (
			SELECT id FROM thunderdome.retro WHERE team_id IN (SELECT id FROM user_teams)
		),
		user_retros AS (
			SELECT u.retro_id AS id FROM thunderdome.retro_user u
			WHERE u.user_id = $1 AND u.abandoned = false
		),
		retros AS (
			SELECT id from user_retros UNION SELECT id FROM team_retros
		)
		SELECT COUNT(*) FROM retros;
	`, userID).Scan(
		&count,
	)
	if e != nil {
		return nil, count, fmt.Errorf("get retros by user count query error: %v", e)
	}

	retroRows, retrosErr := d.DB.Query(`
		WITH user_teams AS (
			SELECT t.id, t.name FROM thunderdome.team_user tu
			LEFT JOIN thunderdome.team t ON t.id = tu.team_id
			WHERE tu.user_id = $1
		),
		team_retros AS (
			SELECT id FROM thunderdome.retro WHERE team_id IN (SELECT id FROM user_teams)
		),
		user_retros AS (
			SELECT u.retro_id AS id FROM thunderdome.retro_user u
			WHERE u.user_id = $1 AND u.abandoned = false
		),
		retros AS (
			SELECT id from user_retros UNION SELECT id FROM team_retros
		)
		SELECT r.id, r.name, r.owner_id, COALESCE(r.team_id::TEXT, ''), r.phase, r.phase_time_limit_min, r.phase_auto_advance, r.template_id,
		 r.allow_cumulative_voting, r.created_date, r.updated_date,
		  MIN(COALESCE(t.name, '')) as teamName,
		  (SELECT row_to_json(t.*) as template FROM thunderdome.retro_template t WHERE t.id = r.template_id) AS template
		FROM thunderdome.retro r
		LEFT JOIN user_teams t ON t.id = r.team_id
		WHERE r.id IN (SELECT id FROM retros)
		GROUP BY r.id, r.created_date ORDER BY r.created_date DESC LIMIT $2 OFFSET $3;
	`, userID, limit, offset)
	if retrosErr != nil {
		d.Logger.Error("get retros by user error", zap.Error(retrosErr))
		return nil, count, fmt.Errorf("get retro by user query error: %v", retrosErr)
	}

	defer retroRows.Close()
	for retroRows.Next() {
		var b = &thunderdome.Retro{
			Users: make([]*thunderdome.RetroUser, 0),
		}
		var Template string
		if err := retroRows.Scan(
			&b.ID,
			&b.Name,
			&b.OwnerID,
			&b.TeamID,
			&b.Phase,
			&b.PhaseTimeLimitMin,
			&b.PhaseAutoAdvance,
			&b.TemplateID,
			&b.AllowCumulativeVoting,
			&b.CreatedDate,
			&b.UpdatedDate,
			&b.TeamName,
			&Template,
		); err != nil {
			d.Logger.Error("get retro by user error", zap.Error(err))
		} else {
			templateError := json.Unmarshal([]byte(Template), &b.Template)
			if templateError != nil {
				d.Logger.Error("retro template json error", zap.Error(templateError))
				return nil, count, fmt.Errorf("get retro by user template error: %v", templateError)
			}

			retros = append(retros, b)
		}
	}

	return retros, count, nil
}

// RetroAdvancePhase sets the phase for the retro
func (d *Service) RetroAdvancePhase(retroID string, phase string) (*thunderdome.Retro, error) {
	var b thunderdome.Retro
	err := d.DB.QueryRow(
		`UPDATE thunderdome.retro
			SET updated_date = NOW(), phase = $2, phase_time_start = NOW(), ready_users = '[]'::jsonb
			WHERE id = $1 RETURNING name, phase_time_start, template_id;`,
		retroID, phase,
	).Scan(&b.Name, &b.PhaseTimeStart, &b.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("retro advance phase query error: %v", err)
	}

	b.ID = retroID
	b.Items = d.GetRetroItems(retroID)
	b.Groups = d.GetRetroGroups(retroID)
	b.ActionItems = d.GetRetroActions(retroID)
	b.Votes = d.GetRetroVotes(retroID)
	b.Phase = phase

	return &b, nil
}

// RetroDelete removes all retro associations and the retro itself from DB by Id
func (d *Service) RetroDelete(retroID string) error {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro WHERE id = $1;`, retroID); err != nil {
		return fmt.Errorf("delete retro query error: %v", err)
	}

	return nil
}

// GetRetros gets a list of retros
func (d *Service) GetRetros(limit int, offset int) ([]*thunderdome.Retro, int, error) {
	var retros = make([]*thunderdome.Retro, 0)
	var count int

	err := d.DB.QueryRow(
		"SELECT COUNT(*) FROM thunderdome.retro;",
	).Scan(
		&count,
	)
	if err != nil {
		return nil, count, fmt.Errorf("get retros count query error: %v", err)
	}

	rows, retrosErr := d.DB.Query(`
		SELECT r.id, COALESCE(r.team_id::TEXT, ''), r.name, r.phase, r.phase_time_limit_min, r.phase_auto_advance, r.allow_cumulative_voting,
		 r.created_date, r.updated_date, r.template_id,
		 (SELECT row_to_json(t.*) as template FROM thunderdome.retro_template t WHERE t.id = r.template_id) AS template
		FROM thunderdome.retro r
		GROUP BY r.id ORDER BY r.created_date DESC
		LIMIT $1 OFFSET $2;
	`, limit, offset)
	if retrosErr != nil {
		return nil, count, fmt.Errorf("get retros query error: %v", retrosErr)
	}

	defer rows.Close()
	for rows.Next() {
		var b = &thunderdome.Retro{
			Users: make([]*thunderdome.RetroUser, 0),
		}
		var template string
		if err := rows.Scan(
			&b.ID,
			&b.TeamID,
			&b.Name,
			&b.Phase,
			&b.PhaseTimeLimitMin,
			&b.PhaseAutoAdvance,
			&b.AllowCumulativeVoting,
			&b.CreatedDate,
			&b.UpdatedDate,
			&b.TemplateID,
			&template,
		); err != nil {
			d.Logger.Error("get retros error", zap.Error(err))
		} else {
			templateError := json.Unmarshal([]byte(template), &b.Template)
			if templateError != nil {
				d.Logger.Error("retro template json error", zap.Error(templateError))
				return nil, count, fmt.Errorf("get retro by user template error: %v", templateError)
			}

			retros = append(retros, b)
		}
	}

	return retros, count, nil
}

// GetActiveRetros gets a list of active retros
func (d *Service) GetActiveRetros(limit int, offset int) ([]*thunderdome.Retro, int, error) {
	var retros = make([]*thunderdome.Retro, 0)
	var count int

	err := d.DB.QueryRow(
		"SELECT COUNT(DISTINCT ru.retro_id) FROM thunderdome.retro_user ru WHERE ru.active IS TRUE;",
	).Scan(
		&count,
	)
	if err != nil {
		return nil, count, fmt.Errorf("get active retros count query error: %v", err)
	}

	rows, retrosErr := d.DB.Query(`
		SELECT r.id, COALESCE(r.team_id::TEXT, ''), r.name, r.phase, r.phase_time_limit_min, r.phase_auto_advance, r.allow_cumulative_voting,
		r.created_date, r.updated_date,
		r.template_id, (SELECT row_to_json(t.*) as template FROM thunderdome.retro_template t WHERE t.id = r.template_id) AS template
		FROM thunderdome.retro_user ru
		LEFT JOIN thunderdome.retro r ON r.id = ru.retro_id
		WHERE ru.active IS TRUE GROUP BY r.id
		LIMIT $1 OFFSET $2;
	`, limit, offset)
	if retrosErr != nil {
		return nil, count, fmt.Errorf("get active retros query error: %v", retrosErr)
	}

	defer rows.Close()
	for rows.Next() {
		var b = &thunderdome.Retro{
			Users: make([]*thunderdome.RetroUser, 0),
		}
		var template string
		if err := rows.Scan(
			&b.ID,
			&b.TeamID,
			&b.Name,
			&b.Phase,
			&b.PhaseTimeLimitMin,
			&b.PhaseAutoAdvance,
			&b.AllowCumulativeVoting,
			&b.CreatedDate,
			&b.UpdatedDate,
			&b.TemplateID,
			template,
		); err != nil {
			d.Logger.Error("get active retros error", zap.Error(err))
		} else {
			templateError := json.Unmarshal([]byte(template), &b.Template)
			if templateError != nil {
				d.Logger.Error("retro template json error", zap.Error(templateError))
				return nil, count, fmt.Errorf("get retro by user template error: %v", templateError)
			}

			retros = append(retros, b)
		}
	}

	return retros, count, nil
}
