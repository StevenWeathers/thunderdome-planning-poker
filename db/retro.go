package db

import (
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// RetroCreate adds a new retro to the db
func (d *Database) RetroCreate(OwnerID string, RetroName string, Format string, JoinCode string, MaxVotes int, BrainstormVisibility string) (*model.Retro, error) {
	var encryptedJoinCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, codeErr
		}
		encryptedJoinCode = EncryptedCode
	}

	var b = &model.Retro{
		OwnerID:     OwnerID,
		Name:        RetroName,
		Format:      "worked_improve_question",
		Phase:       "intro",
		Users:       make([]*model.RetroUser, 0),
		Items:       make([]*model.RetroItem, 0),
		ActionItems: make([]*model.RetroAction, 0),
	}

	e := d.db.QueryRow(
		`SELECT * FROM create_retro($1, $2, $3, $4, $5, $6);`,
		OwnerID,
		RetroName,
		Format,
		encryptedJoinCode,
		MaxVotes,
		BrainstormVisibility,
	).Scan(&b.Id)
	if e != nil {
		d.logger.Error("create retro error", zap.Error(e))
		return nil, e
	}

	// if a join code is set than add owner to retro_user
	// this prevents them from having to enter join code on initial create to enter
	if JoinCode != "" {
		if _, err := d.db.Exec(
			`INSERT INTO retro_user (retro_id, user_id)
		VALUES ($1, $2)`,
			b.Id,
			OwnerID,
		); err != nil {
			d.logger.Error("insert retro user error", zap.Error(err))
		}
	}

	return b, nil
}

// EditRetro updates the retro by ID
func (d *Database) EditRetro(RetroID string, RetroName string, JoinCode string, maxVotes int, brainstormVisibility string) error {
	var encryptedJoinCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return errors.New("unable to revise retro join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if _, err := d.db.Exec(`call edit_retro($1, $2, $3, $4, $5);`,
		RetroID, RetroName, encryptedJoinCode, maxVotes, brainstormVisibility,
	); err != nil {
		d.logger.Error("update retro error", zap.Error(err))
		return errors.New("unable to edit retro")
	}

	return nil
}

// RetroGet gets a retro by ID
func (d *Database) RetroGet(RetroID string) (*model.Retro, error) {
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
	e := d.db.QueryRow(
		`SELECT
			id, name, owner_id, format, phase, COALESCE(join_code, ''), max_votes, brainstorm_visibility, created_date, updated_date
		FROM retro WHERE id = $1`,
		RetroID,
	).Scan(
		&b.Id,
		&b.Name,
		&b.OwnerID,
		&b.Format,
		&b.Phase,
		&b.JoinCode,
		&b.MaxVotes,
		&b.BrainstormVisibility,
		&b.CreatedDate,
		&b.UpdatedDate,
	)
	if e != nil {
		return nil, e
	}

	if b.JoinCode != "" {
		DecryptedCode, codeErr := decrypt(b.JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, errors.New("unable to decode join_code")
		}
		b.JoinCode = DecryptedCode
	}

	b.Items = d.GetRetroItems(RetroID)
	b.Groups = d.GetRetroGroups(RetroID)
	b.Users = d.RetroGetUsers(RetroID)
	b.ActionItems = d.GetRetroActions(RetroID)
	b.Votes = d.GetRetroVotes(RetroID)
	b.Facilitators = d.GetRetroFacilitators(RetroID)

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
			&b.JoinCode,
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
		"SELECT COALESCE(user_id, '') FROM retro_facilitator WHERE retro_id = $1 AND user_id = $2",
		RetroID, userID).Scan(&facilitatorId)
	if err != nil {
		d.logger.Error("get RetroConfirmFacilitator error", zap.Error(err))
		return errors.New("retro Not found")
	}

	if facilitatorId == "" {
		return errors.New("not a facilitator")
	}

	return nil
}

// RetroGetUser gets a user from db by ID and checks retro active status
func (d *Database) RetroGetUser(RetroID string, UserID string) (*model.RetroUser, error) {
	var active bool
	var w model.RetroUser

	err := d.db.QueryRow(
		`SELECT * FROM get_retro_user($1, $2);`,
		RetroID,
		UserID,
	).Scan(
		&w.UserID,
		&w.UserName,
		&active,
	)
	if err != nil {
		d.logger.Error("get retro user error", zap.Error(err))
		return nil, errors.New("User Not found")
	}

	if active {
		return nil, errors.New("User Already Active in Retro")
	}

	return &w, nil
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
			if err := rows.Scan(&w.UserID, &w.UserName, &w.Active, &w.Avatar, &w.GravatarHash); err != nil {
				d.logger.Error("get retro users error", zap.Error(err))
			} else {
				if w.GravatarHash != "" {
					w.GravatarHash = createGravatarHash(w.GravatarHash)
				} else {
					w.GravatarHash = createGravatarHash(w.UserID)
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
func (d *Database) RetroFacilitatorAdd(RetroID string, UserID string) (*model.Retro, error) {
	if _, err := d.db.Exec(
		`call retro_add_facilitator($1, $2);`, RetroID, UserID); err != nil {
		d.logger.Error("call retro_add_facilitator error", zap.Error(err))
		return nil, errors.New("unable to add facilitator")
	}

	retro, err := d.RetroGet(RetroID)
	if err != nil {
		return nil, err
	}

	return retro, nil
}

// RetroFacilitatorRemove removes a retro facilitator
func (d *Database) RetroFacilitatorRemove(RetroID string, UserID string) (*model.Retro, error) {
	if _, err := d.db.Exec(
		`call retro_remove_facilitator($1, $2);`, RetroID, UserID); err != nil {
		d.logger.Error("call retro_remove_facilitator error", zap.Error(err))
		return nil, errors.New("unable to remove facilitator")
	}

	retro, err := d.RetroGet(RetroID)
	if err != nil {
		return nil, err
	}

	return retro, nil
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
	if _, err := d.db.Exec(
		`call set_retro_phase($1, $2);`, RetroID, Phase); err != nil {
		d.logger.Error("call set_retro_phase error", zap.Error(err))
		return nil, errors.New("Unable to advance phase")
	}

	retro, err := d.RetroGet(RetroID)
	if err != nil {
		return nil, err
	}

	return retro, nil
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
		d.logger.Error("get retro user active status error", zap.Error(err))
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
