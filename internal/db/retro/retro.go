package retro

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

// RetroCreate adds a new retro
func (d *Service) RetroCreate(OwnerID string, RetroName string, JoinCode string, FacilitatorCode string, MaxVotes int, BrainstormVisibility string, PhaseTimeLimitMin int, PhaseAutoAdvance bool, TemplateID string) (*thunderdome.Retro, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

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

	var retro = &thunderdome.Retro{
		OwnerID:              OwnerID,
		Name:                 RetroName,
		Phase:                "intro",
		PhaseTimeLimitMin:    PhaseTimeLimitMin,
		PhaseAutoAdvance:     PhaseAutoAdvance,
		Users:                make([]*thunderdome.RetroUser, 0),
		Items:                make([]*thunderdome.RetroItem, 0),
		ActionItems:          make([]*thunderdome.RetroAction, 0),
		BrainstormVisibility: BrainstormVisibility,
		MaxVotes:             MaxVotes,
		TemplateID:           TemplateID,
	}

	err := d.DB.QueryRow(
		`SELECT * FROM thunderdome.retro_create($1, $2, $3, $4, $5, $6, $7, $8, null, $9);`,
		OwnerID,
		RetroName,
		encryptedJoinCode,
		encryptedFacilitatorCode,
		MaxVotes,
		BrainstormVisibility,
		PhaseTimeLimitMin,
		PhaseAutoAdvance,
		TemplateID,
	).Scan(&retro.Id)
	if err != nil {
		return nil, fmt.Errorf("create retro query error: %v", err)
	}

	return retro, nil
}

// TeamRetroCreate adds a new retro associated to a team
func (d *Service) TeamRetroCreate(ctx context.Context, TeamID string, OwnerID string, RetroName string, JoinCode string, FacilitatorCode string, MaxVotes int, BrainstormVisibility string, PhaseTimeLimitMin int, PhaseAutoAdvance bool, TemplateID string) (*thunderdome.Retro, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create team retro encrypt joincode error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create team retro encrypt facilitator code error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	var b = &thunderdome.Retro{
		OwnerID:              OwnerID,
		Name:                 RetroName,
		Phase:                "intro",
		PhaseTimeLimitMin:    0,
		PhaseAutoAdvance:     PhaseAutoAdvance,
		Users:                make([]*thunderdome.RetroUser, 0),
		Items:                make([]*thunderdome.RetroItem, 0),
		ActionItems:          make([]*thunderdome.RetroAction, 0),
		BrainstormVisibility: BrainstormVisibility,
		MaxVotes:             MaxVotes,
		TemplateID:           TemplateID,
	}

	err := d.DB.QueryRowContext(ctx,
		`SELECT * FROM thunderdome.retro_create($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		OwnerID,
		RetroName,
		encryptedJoinCode,
		encryptedFacilitatorCode,
		MaxVotes,
		BrainstormVisibility,
		PhaseTimeLimitMin,
		PhaseAutoAdvance,
		TeamID,
		TemplateID,
	).Scan(&b.Id)
	if err != nil {
		return nil, fmt.Errorf("create team retro query error: %v", err)
	}

	return b, nil
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
			r.id, r.name, r.owner_id, r.phase, r.phase_time_limit_min, r.phase_time_start, r.phase_auto_advance,
			 COALESCE(r.join_code, ''), COALESCE(r.facilitator_code, ''),
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
		&b.Phase,
		&b.PhaseTimeLimitMin,
		&b.PhaseTimeStart,
		&b.PhaseAutoAdvance,
		&JoinCode,
		&FacilitatorCode,
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
		SELECT r.id, r.name, r.owner_id, r.phase, r.phase_time_limit_min, r.phase_auto_advance, r.template_id, r.created_date, r.updated_date,
		  MIN(COALESCE(t.name, '')) as teamName,
		  (SELECT row_to_json(t.*) as template FROM thunderdome.retro_template t WHERE t.id = r.template_id) AS template
		FROM thunderdome.retro r
		LEFT JOIN user_teams t ON t.id = r.team_id
		WHERE r.id IN (SELECT id FROM retros)
		GROUP BY r.id ORDER BY r.created_date DESC LIMIT $2 OFFSET $3;
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
			&b.Phase,
			&b.PhaseTimeLimitMin,
			&b.PhaseAutoAdvance,
			&b.TemplateID,
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

// RetroConfirmFacilitator confirms the user is a facilitator of the retro
func (d *Service) RetroConfirmFacilitator(RetroID string, userID string) error {
	var facilitatorId string
	var role string
	err := d.DB.QueryRow("SELECT type FROM thunderdome.users WHERE id = $1", userID).Scan(&role)
	if err != nil {
		return fmt.Errorf("retro confirm facilitator get user role error: %v", err)
	}

	err = d.DB.QueryRow(
		"SELECT user_id FROM thunderdome.retro_facilitator WHERE retro_id = $1 AND user_id = $2",
		RetroID, userID).Scan(&facilitatorId)
	if err != nil && role != thunderdome.AdminUserType {
		return fmt.Errorf("get retro facilitator error: %v", err)
	}

	return nil
}

// RetroGetUsers retrieves the users for a given retro from db
func (d *Service) RetroGetUsers(RetroID string) []*thunderdome.RetroUser {
	var users = make([]*thunderdome.RetroUser, 0)
	rows, err := d.DB.Query(
		`SELECT
			u.id, u.name, su.active, u.avatar, COALESCE(u.email, ''), COALESCE(u.picture, '')
		FROM thunderdome.retro_user su
		LEFT JOIN thunderdome.users u ON su.user_id = u.id
		WHERE su.retro_id = $1
		ORDER BY u.name;`,
		RetroID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w thunderdome.RetroUser
			if err := rows.Scan(&w.ID, &w.Name, &w.Active, &w.Avatar, &w.Email, &w.PictureURL); err != nil {
				d.Logger.Error("get retro users error", zap.Error(err))
			} else {
				if w.Email != "" {
					w.GravatarHash = db.CreateGravatarHash(w.Email)
				} else {
					w.GravatarHash = db.CreateGravatarHash(w.ID)
				}
				users = append(users, &w)
			}
		}
	}

	return users
}

// GetRetroFacilitators gets a list of retro facilitator ids
func (d *Service) GetRetroFacilitators(RetroID string) []string {
	var facilitators = make([]string, 0)
	rows, err := d.DB.Query(
		`SELECT user_id FROM thunderdome.retro_facilitator WHERE retro_id = $1;`,
		RetroID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var facilitator string
			if err := rows.Scan(&facilitator); err != nil {
				d.Logger.Error("get retro facilitators error", zap.Error(err))
			} else {
				facilitators = append(facilitators, facilitator)
			}
		}
	}

	return facilitators
}

// RetroAddUser adds a user by ID to the retro by ID
func (d *Service) RetroAddUser(RetroID string, UserID string) ([]*thunderdome.RetroUser, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_user (retro_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (retro_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		RetroID,
		UserID,
	); err != nil {
		d.Logger.Error("insert retro user error", zap.Error(err))
	}

	users := d.RetroGetUsers(RetroID)

	return users, nil
}

// RetroFacilitatorAdd adds a retro facilitator
func (d *Service) RetroFacilitatorAdd(RetroID string, UserID string) ([]string, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES ($1, $2);`,
		RetroID, UserID); err != nil {
		return nil, fmt.Errorf("retro add facilitator query error: %v", err)
	}

	facilitators := d.GetRetroFacilitators(RetroID)

	return facilitators, nil
}

// RetroFacilitatorRemove removes a retro facilitator
func (d *Service) RetroFacilitatorRemove(RetroID string, UserID string) ([]string, error) {
	facilitatorCount := 0
	err := d.DB.QueryRow(
		`SELECT count(user_id) FROM thunderdome.retro_facilitator WHERE retro_id = $1;`,
		RetroID,
	).Scan(&facilitatorCount)
	if err != nil {
		return nil, fmt.Errorf("retro remove facilitator query error: %v", err)
	}

	if facilitatorCount == 1 {
		return nil, fmt.Errorf("ONLY_FACILITATOR")
	}

	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_facilitator WHERE retro_id = $1 AND user_id = $2;`,
		RetroID, UserID); err != nil {
		return nil, fmt.Errorf("retro remove facilitator query error: %v", err)
	}

	facilitators := d.GetRetroFacilitators(RetroID)

	return facilitators, nil
}

// RetroRetreatUser removes a user from the current retro by ID
func (d *Service) RetroRetreatUser(RetroID string, UserID string) []*thunderdome.RetroUser {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.retro_user SET active = false WHERE retro_id = $1 AND user_id = $2`, RetroID, UserID); err != nil {
		d.Logger.Error("update retro user active false error", zap.Error(err))
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("update user last active timestamp error", zap.Error(err))
	}

	users := d.RetroGetUsers(RetroID)

	return users
}

// RetroAbandon removes a user from the current retro by ID and sets abandoned true
func (d *Service) RetroAbandon(RetroID string, UserID string) ([]*thunderdome.RetroUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.retro_user SET active = false, abandoned = true WHERE retro_id = $1 AND user_id = $2`, RetroID, UserID); err != nil {
		return nil, fmt.Errorf("abandon retro query error: %v", err)
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		return nil, fmt.Errorf("abandon retro update user query error: %v", err)
	}

	users := d.RetroGetUsers(RetroID)

	return users, nil
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

// GetRetroUserActiveStatus checks retro active status of User for given retro
func (d *Service) GetRetroUserActiveStatus(RetroID string, UserID string) error {
	var active bool

	err := d.DB.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM thunderdome.retro_user
		WHERE user_id = $2 AND retro_id = $1;`,
		RetroID,
		UserID,
	).Scan(
		&active,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get retro user active status query error: %v", err)
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if active {
		return errors.New("DUPLICATE_RETRO_USER")
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
		SELECT r.id, r.name, r.phase, r.phase_time_limit_min, r.phase_auto_advance, r.created_date,
		 r.updated_date, r.template_id,
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
			&b.Name,
			&b.Phase,
			&b.PhaseTimeLimitMin,
			&b.PhaseAutoAdvance,
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
		SELECT r.id, r.name, r.phase, r.phase_time_limit_min, r.phase_auto_advance, r.created_date, r.updated_date,
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
			&b.Name,
			&b.Phase,
			&b.PhaseTimeLimitMin,
			&b.PhaseAutoAdvance,
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

// GetRetroFacilitatorCode retrieve the retro facilitator code
func (d *Service) GetRetroFacilitatorCode(RetroID string) (string, error) {
	var EncryptedCode string

	if err := d.DB.QueryRow(`
		SELECT COALESCE(facilitator_code, '') FROM thunderdome.retro
		WHERE id = $1`,
		RetroID,
	).Scan(&EncryptedCode); err != nil {
		return "", fmt.Errorf("get retro facilitator_code query error: %v", err)
	}

	if EncryptedCode == "" {
		return "", fmt.Errorf("retro facilitator_code not set")
	}
	DecryptedCode, codeErr := db.Decrypt(EncryptedCode, d.AESHashKey)
	if codeErr != nil {
		return "", fmt.Errorf("retrieve retro facilitator_code decrypt error: %v", codeErr)
	}

	return DecryptedCode, nil
}

// CleanRetros deletes retros older than {DaysOld} days
func (d *Service) CleanRetros(ctx context.Context, DaysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.retro WHERE updated_date < (NOW() - $1 * interval '1 day');`,
		DaysOld,
	); err != nil {
		return fmt.Errorf("clean retros query error: %v", err)
	}

	return nil
}

// MarkUserReady marks a user as ready for next phase
func (d *Service) MarkUserReady(RetroID string, userID string) ([]string, error) {
	var rawReadyUsers string
	readyUsers := make([]string, 0)

	err := d.DB.QueryRow(
		`UPDATE thunderdome.retro 
		SET updated_date = NOW(), ready_users = ready_users::jsonb || to_jsonb(array[$2])
		WHERE id = $1 
		RETURNING ready_users;`,
		RetroID, userID,
	).Scan(&rawReadyUsers)
	if err != nil {
		return readyUsers, fmt.Errorf("retro MarkUserReady query error: %v", err)
	}

	err = json.Unmarshal([]byte(rawReadyUsers), &readyUsers)
	if err != nil {
		d.Logger.Error("ready_users json error", zap.Error(err))
	}

	return readyUsers, nil
}

// UnmarkUserReady un-marks a user as ready for next phase
func (d *Service) UnmarkUserReady(RetroID string, userID string) ([]string, error) {
	var rawReadyUsers string
	readyUsers := make([]string, 0)

	err := d.DB.QueryRow(
		`UPDATE thunderdome.retro 
		SET updated_date = NOW(), ready_users = ready_users::jsonb - $2
		WHERE id = $1 
		RETURNING ready_users;`,
		RetroID, userID,
	).Scan(&rawReadyUsers)
	if err != nil {
		return readyUsers, fmt.Errorf("retro UnmarkUserReady query error: %v", err)
	}

	err = json.Unmarshal([]byte(rawReadyUsers), &readyUsers)
	if err != nil {
		d.Logger.Error("ready_users json error", zap.Error(err))
	}

	return readyUsers, nil
}
