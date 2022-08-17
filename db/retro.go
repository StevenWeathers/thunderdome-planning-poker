package db

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// RetroCreate adds a new retro
func (d *Database) RetroCreate(OwnerID string, RetroName string, Format string, JoinCode string, FacilitatorCode string, MaxVotes int, BrainstormVisibility string) (*model.Retro, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, codeErr
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := encrypt(FacilitatorCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, codeErr
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	var b = &model.Retro{
		OwnerID:              OwnerID,
		Name:                 RetroName,
		Format:               Format,
		Phase:                "intro",
		Users:                make([]*model.RetroUser, 0),
		Items:                make([]*model.RetroItem, 0),
		ActionItems:          make([]*model.RetroAction, 0),
		BrainstormVisibility: BrainstormVisibility,
		MaxVotes:             MaxVotes,
	}

	e := d.db.QueryRow(
		`SELECT * FROM create_retro($1, $2, $3, $4, $5, $6, $7);`,
		OwnerID,
		RetroName,
		Format,
		encryptedJoinCode,
		encryptedFacilitatorCode,
		MaxVotes,
		BrainstormVisibility,
	).Scan(&b.Id)
	if e != nil {
		d.logger.Error("create_retro query error", zap.Error(e))
		return nil, errors.New("error creating retro")
	}

	return b, nil
}

// TeamRetroCreate adds a new retro associated to a team
func (d *Database) TeamRetroCreate(ctx context.Context, TeamID string, OwnerID string, RetroName string, Format string, JoinCode string, FacilitatorCode string, MaxVotes int, BrainstormVisibility string) (*model.Retro, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, codeErr
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := encrypt(FacilitatorCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, codeErr
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	var b = &model.Retro{
		OwnerID:              OwnerID,
		Name:                 RetroName,
		Format:               Format,
		Phase:                "intro",
		Users:                make([]*model.RetroUser, 0),
		Items:                make([]*model.RetroItem, 0),
		ActionItems:          make([]*model.RetroAction, 0),
		BrainstormVisibility: BrainstormVisibility,
		MaxVotes:             MaxVotes,
	}

	e := d.db.QueryRowContext(ctx,
		`SELECT * FROM team_create_retro($1, $2, $3, $4, $5, $6, $7, $8);`,
		TeamID,
		OwnerID,
		RetroName,
		Format,
		encryptedJoinCode,
		encryptedFacilitatorCode,
		MaxVotes,
		BrainstormVisibility,
	).Scan(&b.Id)
	if e != nil {
		d.logger.Error("team_create_retro query error", zap.Error(e))
		return nil, errors.New("error creating retro")
	}

	return b, nil
}

// EditRetro updates the retro by ID
func (d *Database) EditRetro(RetroID string, RetroName string, JoinCode string, FacilitatorCode string, maxVotes int, brainstormVisibility string) error {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return errors.New("unable to revise retro join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := encrypt(FacilitatorCode, d.config.AESHashkey)
		if codeErr != nil {
			return errors.New("unable to revise retro facilitator_code")
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	if _, err := d.db.Exec(`call edit_retro($1, $2, $3, $4, $5, $6);`,
		RetroID, RetroName, encryptedJoinCode, encryptedFacilitatorCode, maxVotes, brainstormVisibility,
	); err != nil {
		d.logger.Error("update retro error", zap.Error(err))
		return errors.New("unable to edit retro")
	}

	return nil
}

// RetroGet gets a retro by ID
func (d *Database) RetroGet(RetroID string, UserID string) (*model.Retro, error) {
	var b = &model.Retro{
		Id:           RetroID,
		Users:        make([]*model.RetroUser, 0),
		Items:        make([]*model.RetroItem, 0),
		Groups:       make([]*model.RetroGroup, 0),
		ActionItems:  make([]*model.RetroAction, 0),
		Votes:        make([]*model.RetroVote, 0),
		Facilitators: make([]string, 0),
	}

	// get retro
	var JoinCode string
	var FacilitatorCode string
	var Facilitators string
	e := d.db.QueryRow(
		`SELECT
			r.id, r.name, r.owner_id, r.format, r.phase, COALESCE(r.join_code, ''), COALESCE(r.facilitator_code, ''),
			r.max_votes, r.brainstorm_visibility, r.created_date, r.updated_date,
			CASE WHEN COUNT(rf) = 0 THEN '[]'::json ELSE array_to_json(array_agg(rf.user_id)) END AS facilitators
		FROM retro r 
		LEFT JOIN retro_facilitator rf ON r.id = rf.retro_id
		WHERE id = $1
		GROUP BY r.id`,
		RetroID,
	).Scan(
		&b.Id,
		&b.Name,
		&b.OwnerID,
		&b.Format,
		&b.Phase,
		&JoinCode,
		&FacilitatorCode,
		&b.MaxVotes,
		&b.BrainstormVisibility,
		&b.CreatedDate,
		&b.UpdatedDate,
		&Facilitators,
	)
	if e != nil {
		d.logger.Error("", zap.Error(e))
		return nil, e
	}

	facilError := json.Unmarshal([]byte(Facilitators), &b.Facilitators)
	if facilError != nil {
		d.logger.Error("facilitators json error", zap.Error(facilError))
	}
	isFacilitator := contains(b.Facilitators, UserID)

	if JoinCode != "" {
		DecryptedCode, codeErr := decrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, errors.New("unable to decode join_code")
		}
		b.JoinCode = DecryptedCode
	}

	if FacilitatorCode != "" && isFacilitator {
		DecryptedCode, codeErr := decrypt(FacilitatorCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, errors.New("unable to decode facilitator_code")
		}
		b.FacilitatorCode = DecryptedCode
	}

	b.Items = d.GetRetroItems(RetroID)
	b.Groups = d.GetRetroGroups(RetroID)
	b.Users = d.RetroGetUsers(RetroID)
	b.ActionItems = d.GetRetroActions(RetroID)
	b.Votes = d.GetRetroVotes(RetroID)

	return b, nil
}

// RetroGetByUser gets a list of retros by UserID
func (d *Database) RetroGetByUser(UserID string) ([]*model.Retro, error) {
	var retros = make([]*model.Retro, 0)
	retroRows, retrosErr := d.db.Query(`
		SELECT * FROM get_retros_by_user($1);
	`, UserID)
	if retrosErr != nil {
		return nil, retrosErr
	}

	defer retroRows.Close()
	for retroRows.Next() {
		var b = &model.Retro{
			Users: make([]*model.RetroUser, 0),
		}
		if err := retroRows.Scan(
			&b.Id,
			&b.Name,
			&b.OwnerID,
			&b.Format,
			&b.Phase,
			&b.CreatedDate,
			&b.UpdatedDate,
		); err != nil {
			d.logger.Error("get retro by user error", zap.Error(err))
		} else {
			retros = append(retros, b)
		}
	}

	return retros, nil
}

// RetroConfirmFacilitator confirms the user is a facilitator of the retro
func (d *Database) RetroConfirmFacilitator(RetroID string, userID string) error {
	var facilitatorId string
	err := d.db.QueryRow(
		"SELECT user_id FROM retro_facilitator WHERE retro_id = $1 AND user_id = $2",
		RetroID, userID).Scan(&facilitatorId)
	if err != nil {
		d.logger.Error("get RetroConfirmFacilitator error", zap.Error(err))
		return errors.New("retro facilitator not found")
	}

	return nil
}

// RetroGetUsers retrieves the users for a given retro from db
func (d *Database) RetroGetUsers(RetroID string) []*model.RetroUser {
	var users = make([]*model.RetroUser, 0)
	rows, err := d.db.Query(
		`SELECT * FROM get_retro_users($1);`,
		RetroID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w model.RetroUser
			if err := rows.Scan(&w.ID, &w.Name, &w.Active, &w.Avatar, &w.GravatarHash); err != nil {
				d.logger.Error("get retro users error", zap.Error(err))
			} else {
				if w.GravatarHash != "" {
					w.GravatarHash = createGravatarHash(w.GravatarHash)
				} else {
					w.GravatarHash = createGravatarHash(w.ID)
				}
				users = append(users, &w)
			}
		}
	}

	return users
}

// GetRetroFacilitators gets a list of retro facilitator ids
func (d *Database) GetRetroFacilitators(RetroID string) []string {
	var facilitators = make([]string, 0)
	rows, err := d.db.Query(
		`SELECT user_id FROM retro_facilitator WHERE retro_id = $1;`,
		RetroID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var facilitator string
			if err := rows.Scan(&facilitator); err != nil {
				d.logger.Error("get retro facilitators error", zap.Error(err))
			} else {
				facilitators = append(facilitators, facilitator)
			}
		}
	}

	return facilitators
}

// RetroAddUser adds a user by ID to the retro by ID
func (d *Database) RetroAddUser(RetroID string, UserID string) ([]*model.RetroUser, error) {
	if _, err := d.db.Exec(
		`INSERT INTO retro_user (retro_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (retro_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		RetroID,
		UserID,
	); err != nil {
		d.logger.Error("insert retro user error", zap.Error(err))
	}

	users := d.RetroGetUsers(RetroID)

	return users, nil
}

// RetroFacilitatorAdd adds a retro facilitator
func (d *Database) RetroFacilitatorAdd(RetroID string, UserID string) ([]string, error) {
	if _, err := d.db.Exec(
		`call retro_add_facilitator($1, $2);`, RetroID, UserID); err != nil {
		d.logger.Error("call retro_add_facilitator error", zap.Error(err))
		return nil, errors.New("unable to add facilitator")
	}

	facilitators := d.GetRetroFacilitators(RetroID)

	return facilitators, nil
}

// RetroFacilitatorRemove removes a retro facilitator
func (d *Database) RetroFacilitatorRemove(RetroID string, UserID string) ([]string, error) {
	if _, err := d.db.Exec(
		`call retro_remove_facilitator($1, $2);`, RetroID, UserID); err != nil {
		d.logger.Error("call retro_remove_facilitator error", zap.Error(err))
		return nil, errors.New("unable to remove facilitator")
	}

	facilitators := d.GetRetroFacilitators(RetroID)

	return facilitators, nil
}

// RetroRetreatUser removes a user from the current retro by ID
func (d *Database) RetroRetreatUser(RetroID string, UserID string) []*model.RetroUser {
	if _, err := d.db.Exec(
		`UPDATE retro_user SET active = false WHERE retro_id = $1 AND user_id = $2`, RetroID, UserID); err != nil {
		d.logger.Error("update retro user active false error", zap.Error(err))
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.logger.Error("update user last active timestamp error", zap.Error(err))
	}

	users := d.RetroGetUsers(RetroID)

	return users
}

// RetroAbandon removes a user from the current retro by ID and sets abandoned true
func (d *Database) RetroAbandon(RetroID string, UserID string) ([]*model.RetroUser, error) {
	if _, err := d.db.Exec(
		`UPDATE retro_user SET active = false, abandoned = true WHERE retro_id = $1 AND user_id = $2`, RetroID, UserID); err != nil {
		d.logger.Error("update retro user abandoned true error", zap.Error(err))
		return nil, err
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.logger.Error("update user last active timestamp error", zap.Error(err))
		return nil, err
	}

	users := d.RetroGetUsers(RetroID)

	return users, nil
}

// RetroAdvancePhase sets the phase for the retro
func (d *Database) RetroAdvancePhase(RetroID string, Phase string) (*model.Retro, error) {
	var b model.Retro
	if _, err := d.db.Exec(
		`call set_retro_phase($1, $2);`, RetroID, Phase); err != nil {
		d.logger.Error("call set_retro_phase error", zap.Error(err))
		return nil, errors.New("Unable to advance phase")
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
func (d *Database) RetroDelete(RetroID string) error {
	if _, err := d.db.Exec(
		`call delete_retro($1);`, RetroID); err != nil {
		d.logger.Error("call delete_retro error", zap.Error(err))
		return err
	}

	return nil
}

// GetRetroUserActiveStatus checks retro active status of User for given retro
func (d *Database) GetRetroUserActiveStatus(RetroID string, UserID string) error {
	var active bool

	err := d.db.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM retro_user
		WHERE user_id = $2 AND retro_id = $1;`,
		RetroID,
		UserID,
	).Scan(
		&active,
	)
	if err != nil {
		return err
	}

	if active {
		return errors.New("DUPLICATE_RETRO_USER")
	}

	return nil
}

// GetRetros gets a list of retros
func (d *Database) GetRetros(Limit int, Offset int) ([]*model.Retro, int, error) {
	var retros = make([]*model.Retro, 0)
	var Count int

	e := d.db.QueryRow(
		"SELECT COUNT(*) FROM retro;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	rows, retrosErr := d.db.Query(`
		SELECT r.id, r.name, r.format, r.phase, r.created_date, r.updated_date
		FROM retro r
		GROUP BY r.id ORDER BY r.created_date DESC
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if retrosErr != nil {
		return nil, Count, retrosErr
	}

	defer rows.Close()
	for rows.Next() {
		var b = &model.Retro{
			Users: make([]*model.RetroUser, 0),
		}
		if err := rows.Scan(
			&b.Id,
			&b.Name,
			&b.Format,
			&b.Phase,
			&b.CreatedDate,
			&b.UpdatedDate,
		); err != nil {
			d.logger.Error("get retros error", zap.Error(err))
		} else {
			retros = append(retros, b)
		}
	}

	return retros, Count, nil
}

// GetActiveRetros gets a list of active retros
func (d *Database) GetActiveRetros(Limit int, Offset int) ([]*model.Retro, int, error) {
	var retros = make([]*model.Retro, 0)
	var Count int

	e := d.db.QueryRow(
		"SELECT COUNT(DISTINCT ru.retro_id) FROM retro_user ru WHERE ru.active IS TRUE;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	rows, retrosErr := d.db.Query(`
		SELECT r.id, r.name, r.format, r.phase, r.created_date, r.updated_date
		FROM retro_user ru
		LEFT JOIN retro r ON r.id = ru.retro_id
		WHERE ru.active IS TRUE GROUP BY r.id
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if retrosErr != nil {
		return nil, Count, retrosErr
	}

	defer rows.Close()
	for rows.Next() {
		var b = &model.Retro{
			Users: make([]*model.RetroUser, 0),
		}
		if err := rows.Scan(
			&b.Id,
			&b.Name,
			&b.Format,
			&b.Phase,
			&b.CreatedDate,
			&b.UpdatedDate,
		); err != nil {
			d.logger.Error("get active retros error", zap.Error(err))
		} else {
			retros = append(retros, b)
		}
	}

	return retros, Count, nil
}

// GetRetroFacilitatorCode retrieve the retro facilitator code
func (d *Database) GetRetroFacilitatorCode(RetroID string) (string, error) {
	var EncryptedCode string

	if err := d.db.QueryRow(`
		SELECT COALESCE(facilitator_code, '') FROM retro
		WHERE id = $1`,
		RetroID,
	).Scan(&EncryptedCode); err != nil {
		d.logger.Error("get retro facilitator_code error", zap.Error(err))
		return "", errors.New("unable to retrieve retro facilitator_code")
	}

	if EncryptedCode == "" {
		return "", errors.New("unable to retrieve retro facilitator_code")
	}
	DecryptedCode, codeErr := decrypt(EncryptedCode, d.config.AESHashkey)
	if codeErr != nil {
		return "", errors.New("unable to retrieve retro facilitator_code")
	}

	return DecryptedCode, nil
}
