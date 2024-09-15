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

// Service represents a PostgreSQL implementation of thunderdome.RetroDataSvc.
type Service struct {
	DB         *sql.DB
	Logger     *otelzap.Logger
	AESHashKey string
}

func (d *Service) CreateRetro(ctx context.Context, OwnerID, TeamID string, RetroName, JoinCode, FacilitatorCode string, MaxVotes int, BrainstormVisibility string, PhaseTimeLimitMin int, PhaseAutoAdvance bool, AllowCumulativeVoting bool, TemplateID string) (*thunderdome.Retro, error) {
	var encryptedFacilitatorCode string
	var encryptedJoinCode string
	var retro = &thunderdome.Retro{
		OwnerID:               OwnerID,
		TeamID:                TeamID,
		Name:                  RetroName,
		Phase:                 "intro",
		PhaseTimeLimitMin:     PhaseTimeLimitMin,
		PhaseAutoAdvance:      PhaseAutoAdvance,
		Users:                 make([]*thunderdome.RetroUser, 0),
		Items:                 make([]*thunderdome.RetroItem, 0),
		ActionItems:           make([]*thunderdome.RetroAction, 0),
		BrainstormVisibility:  BrainstormVisibility,
		MaxVotes:              MaxVotes,
		TemplateID:            TemplateID,
		AllowCumulativeVoting: AllowCumulativeVoting,
	}

	if JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create retro encrypt join code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create retro encrypt facilitator code error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
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
		RETURNING id
	`, OwnerID, TeamID, RetroName, encryptedJoinCode, encryptedFacilitatorCode, MaxVotes, BrainstormVisibility,
		PhaseTimeLimitMin, PhaseAutoAdvance, AllowCumulativeVoting, TemplateID).Scan(&retro.Id)

	if err != nil {
		d.Logger.Error("create retro error", zap.Error(err),
			zap.String("owner_id", OwnerID), zap.String("name", RetroName))
		return nil, fmt.Errorf("failed to insert into retro table: %v", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO thunderdome.retro_facilitator (retro_id, user_id)
		VALUES ($1, $2)
	`, retro.Id, OwnerID)

	if err != nil {
		d.Logger.Error("create retro error", zap.Error(err))
		return nil, fmt.Errorf("failed to insert into retro_facilitator table: %v", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO thunderdome.retro_user (retro_id, user_id)
		VALUES ($1, $2)
	`, retro.Id, OwnerID)

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
func (d *Service) EditRetro(RetroID string, RetroName string, JoinCode string, FacilitatorCode string, maxVotes int, brainstormVisibility string, phaseAutoAdvance bool) error {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return fmt.Errorf("edit retro encrypt join code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return fmt.Errorf("edit retro encrypt join facilitator error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	if _, err := d.DB.Exec(`UPDATE thunderdome.retro
    SET name = $2, join_code = $3, facilitator_code = $4, max_votes = $5,
        brainstorm_visibility = $6, phase_auto_advance = $7, updated_date = NOW()
    WHERE id = $1;`,
		RetroID, RetroName, encryptedJoinCode, encryptedFacilitatorCode,
		maxVotes, brainstormVisibility, phaseAutoAdvance,
	); err != nil {
		return fmt.Errorf("edit retro query error: %v", err)
	}

	return nil
}

// RetroGet gets a retro by ID
func (d *Service) RetroGet(RetroID string, UserID string) (*thunderdome.Retro, error) {
	var b = &thunderdome.Retro{
		Id:           RetroID,
		Users:        make([]*thunderdome.RetroUser, 0),
		Items:        make([]*thunderdome.RetroItem, 0),
		Groups:       make([]*thunderdome.RetroGroup, 0),
		ActionItems:  make([]*thunderdome.RetroAction, 0),
		Votes:        make([]*thunderdome.RetroVote, 0),
		Facilitators: make([]string, 0),
		ReadyUsers:   make([]string, 0),
	}

	// get retro
	var JoinCode string
	var FacilitatorCode string
	var Facilitators string
	var ReadyUsers string
	var Template string
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
		RetroID,
	).Scan(
		&b.Id,
		&b.Name,
		&b.OwnerID,
		&b.TeamID,
		&b.Phase,
		&b.PhaseTimeLimitMin,
		&b.PhaseTimeStart,
		&b.PhaseAutoAdvance,
		&JoinCode,
		&FacilitatorCode,
		&b.AllowCumulativeVoting,
		&b.MaxVotes,
		&b.BrainstormVisibility,
		&ReadyUsers,
		&b.CreatedDate,
		&b.UpdatedDate,
		&b.TemplateID,
		&Facilitators,
		&Template,
	)
	if err != nil {
		d.Logger.Error("get retro error", zap.Error(err))
		return nil, fmt.Errorf("get retro query error: %v", err)
	}

	facilError := json.Unmarshal([]byte(Facilitators), &b.Facilitators)
	if facilError != nil {
		d.Logger.Error("facilitators json error", zap.Error(facilError))
	}
	isFacilitator := db.Contains(b.Facilitators, UserID)

	if JoinCode != "" {
		DecryptedCode, codeErr := db.Decrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("get retro decrypt join code error: %v", codeErr)
		}
		b.JoinCode = DecryptedCode
	}

	if FacilitatorCode != "" && isFacilitator {
		DecryptedCode, codeErr := db.Decrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("get retro decrypt join facilitator error: %v", codeErr)
		}
		b.FacilitatorCode = DecryptedCode
	}

	readyUsersError := json.Unmarshal([]byte(ReadyUsers), &b.ReadyUsers)
	if readyUsersError != nil {
		d.Logger.Error("ready users json error", zap.Error(readyUsersError))
	}

	templateError := json.Unmarshal([]byte(Template), &b.Template)
	if templateError != nil {
		d.Logger.Error("retro template json error", zap.Error(templateError))
		return nil, fmt.Errorf("get retro template error: %v", templateError)
	}

	b.Items = d.GetRetroItems(RetroID)
	b.Groups = d.GetRetroGroups(RetroID)
	b.Users = d.RetroGetUsers(RetroID)
	b.ActionItems = d.GetRetroActions(RetroID)
	b.Votes = d.GetRetroVotes(RetroID)

	return b, nil
}

// RetroGetByUser gets a list of retros by UserID
func (d *Service) RetroGetByUser(UserID string, Limit int, Offset int) ([]*thunderdome.Retro, int, error) {
	var retros = make([]*thunderdome.Retro, 0)
	var Count int

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
	`, UserID).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, fmt.Errorf("get retros by user count query error: %v", e)
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
	`, UserID, Limit, Offset)
	if retrosErr != nil {
		d.Logger.Error("get retros by user error", zap.Error(retrosErr))
		return nil, Count, fmt.Errorf("get retro by user query error: %v", retrosErr)
	}

	defer retroRows.Close()
	for retroRows.Next() {
		var b = &thunderdome.Retro{
			Users: make([]*thunderdome.RetroUser, 0),
		}
		var Template string
		if err := retroRows.Scan(
			&b.Id,
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
				return nil, Count, fmt.Errorf("get retro by user template error: %v", templateError)
			}

			retros = append(retros, b)
		}
	}

	return retros, Count, nil
}

// RetroAdvancePhase sets the phase for the retro
func (d *Service) RetroAdvancePhase(RetroID string, Phase string) (*thunderdome.Retro, error) {
	var b thunderdome.Retro
	err := d.DB.QueryRow(
		`UPDATE thunderdome.retro 
			SET updated_date = NOW(), phase = $2, phase_time_start = NOW(), ready_users = '[]'::jsonb 
			WHERE id = $1 RETURNING name, phase_time_start, template_id;`,
		RetroID, Phase,
	).Scan(&b.Name, &b.PhaseTimeStart, &b.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("retro advance phase query error: %v", err)
	}

	b.Id = RetroID
	b.Items = d.GetRetroItems(RetroID)
	b.Groups = d.GetRetroGroups(RetroID)
	b.ActionItems = d.GetRetroActions(RetroID)
	b.Votes = d.GetRetroVotes(RetroID)
	b.Phase = Phase

	return &b, nil
}

// RetroDelete removes all retro associations and the retro itself from DB by Id
func (d *Service) RetroDelete(RetroID string) error {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro WHERE id = $1;`, RetroID); err != nil {
		return fmt.Errorf("delete retro query error: %v", err)
	}

	return nil
}

// GetRetros gets a list of retros
func (d *Service) GetRetros(Limit int, Offset int) ([]*thunderdome.Retro, int, error) {
	var retros = make([]*thunderdome.Retro, 0)
	var Count int

	err := d.DB.QueryRow(
		"SELECT COUNT(*) FROM thunderdome.retro;",
	).Scan(
		&Count,
	)
	if err != nil {
		return nil, Count, fmt.Errorf("get retros count query error: %v", err)
	}

	rows, retrosErr := d.DB.Query(`
		SELECT r.id, COALESCE(r.team_id::TEXT, ''), r.name, r.phase, r.phase_time_limit_min, r.phase_auto_advance, r.allow_cumulative_voting,
		 r.created_date, r.updated_date, r.template_id,
		 (SELECT row_to_json(t.*) as template FROM thunderdome.retro_template t WHERE t.id = r.template_id) AS template
		FROM thunderdome.retro r
		GROUP BY r.id ORDER BY r.created_date DESC
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if retrosErr != nil {
		return nil, Count, fmt.Errorf("get retros query error: %v", retrosErr)
	}

	defer rows.Close()
	for rows.Next() {
		var b = &thunderdome.Retro{
			Users: make([]*thunderdome.RetroUser, 0),
		}
		var Template string
		if err := rows.Scan(
			&b.Id,
			&b.TeamID,
			&b.Name,
			&b.Phase,
			&b.PhaseTimeLimitMin,
			&b.PhaseAutoAdvance,
			&b.AllowCumulativeVoting,
			&b.CreatedDate,
			&b.UpdatedDate,
			&b.TemplateID,
			&Template,
		); err != nil {
			d.Logger.Error("get retros error", zap.Error(err))
		} else {
			templateError := json.Unmarshal([]byte(Template), &b.Template)
			if templateError != nil {
				d.Logger.Error("retro template json error", zap.Error(templateError))
				return nil, Count, fmt.Errorf("get retro by user template error: %v", templateError)
			}

			retros = append(retros, b)
		}
	}

	return retros, Count, nil
}

// GetActiveRetros gets a list of active retros
func (d *Service) GetActiveRetros(Limit int, Offset int) ([]*thunderdome.Retro, int, error) {
	var retros = make([]*thunderdome.Retro, 0)
	var Count int

	err := d.DB.QueryRow(
		"SELECT COUNT(DISTINCT ru.retro_id) FROM thunderdome.retro_user ru WHERE ru.active IS TRUE;",
	).Scan(
		&Count,
	)
	if err != nil {
		return nil, Count, fmt.Errorf("get active retros count query error: %v", err)
	}

	rows, retrosErr := d.DB.Query(`
		SELECT r.id, COALESCE(r.team_id::TEXT, ''), r.name, r.phase, r.phase_time_limit_min, r.phase_auto_advance, r.allow_cumulative_voting, 
		r.created_date, r.updated_date,
		r.template_id, (SELECT row_to_json(t.*) as template FROM thunderdome.retro_template t WHERE t.id = r.template_id) AS template
		FROM thunderdome.retro_user ru
		LEFT JOIN thunderdome.retro r ON r.id = ru.retro_id
		WHERE ru.active IS TRUE GROUP BY r.id
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if retrosErr != nil {
		return nil, Count, fmt.Errorf("get active retros query error: %v", retrosErr)
	}

	defer rows.Close()
	for rows.Next() {
		var b = &thunderdome.Retro{
			Users: make([]*thunderdome.RetroUser, 0),
		}
		var Template string
		if err := rows.Scan(
			&b.Id,
			&b.TeamID,
			&b.Name,
			&b.Phase,
			&b.PhaseTimeLimitMin,
			&b.PhaseAutoAdvance,
			&b.AllowCumulativeVoting,
			&b.CreatedDate,
			&b.UpdatedDate,
			&b.TemplateID,
			Template,
		); err != nil {
			d.Logger.Error("get active retros error", zap.Error(err))
		} else {
			templateError := json.Unmarshal([]byte(Template), &b.Template)
			if templateError != nil {
				d.Logger.Error("retro template json error", zap.Error(templateError))
				return nil, Count, fmt.Errorf("get retro by user template error: %v", templateError)
			}

			retros = append(retros, b)
		}
	}

	return retros, Count, nil
}
