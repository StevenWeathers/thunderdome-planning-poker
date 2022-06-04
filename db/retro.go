package db

import (
	"errors"
	"log"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
)

// RetroCreate adds a new retro to the db
func (d *Database) RetroCreate(OwnerID string, RetroName string, Format string, JoinCode string) (*model.Retro, error) {
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
		`SELECT * FROM create_retro($1, $2, $3, $4);`,
		OwnerID,
		RetroName,
		Format,
		encryptedJoinCode,
	).Scan(&b.Id)
	if e != nil {
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
			log.Println(err)
		}
	}

	return b, nil
}

// RetroGet gets a retro by ID
func (d *Database) RetroGet(RetroID string) (*model.Retro, error) {
	var b = &model.Retro{
		Id:          RetroID,
		Users:       make([]*model.RetroUser, 0),
		Items:       make([]*model.RetroItem, 0),
		Groups:      make([]*model.RetroGroup, 0),
		ActionItems: make([]*model.RetroAction, 0),
		Votes:       make([]*model.RetroVote, 0),
	}

	// get retro
	e := d.db.QueryRow(
		`SELECT
			id, name, owner_id, format, phase, COALESCE(join_code, '')
		FROM retro WHERE id = $1`,
		RetroID,
	).Scan(
		&b.Id,
		&b.Name,
		&b.OwnerID,
		&b.Format,
		&b.Phase,
		&b.JoinCode,
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
		); err != nil {
			log.Println(err)
		} else {
			retros = append(retros, b)
		}
	}

	return retros, nil
}

// RetroConfirmOwner confirms the user is infact owner of the retro
func (d *Database) RetroConfirmOwner(RetroID string, userID string) error {
	var ownerID string
	e := d.db.QueryRow("SELECT owner_id FROM retro WHERE id = $1", RetroID).Scan(&ownerID)
	if e != nil {
		log.Println(e)
		return errors.New("Retro Not found")
	}

	if ownerID != userID {
		return errors.New("Not Owner")
	}

	return nil
}

// RetroGetUser gets a user from db by ID and checks retro active status
func (d *Database) RetroGetUser(RetroID string, UserID string) (*model.RetroUser, error) {
	var active bool
	var w model.RetroUser

	e := d.db.QueryRow(
		`SELECT * FROM get_retro_user($1, $2);`,
		RetroID,
		UserID,
	).Scan(
		&w.UserID,
		&w.UserName,
		&active,
	)
	if e != nil {
		log.Println(e)
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
				log.Println(err)
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

// RetroAddUser adds a user by ID to the retro by ID
func (d *Database) RetroAddUser(RetroID string, UserID string) ([]*model.RetroUser, error) {
	if _, err := d.db.Exec(
		`INSERT INTO retro_user (retro_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (retro_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		RetroID,
		UserID,
	); err != nil {
		log.Println(err)
	}

	users := d.RetroGetUsers(RetroID)

	return users, nil
}

// RetroRetreatUser removes a user from the current retro by ID
func (d *Database) RetroRetreatUser(RetroID string, UserID string) []*model.RetroUser {
	if _, err := d.db.Exec(
		`UPDATE retro_user SET active = false WHERE retro_id = $1 AND user_id = $2`, RetroID, UserID); err != nil {
		log.Println(err)
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		log.Println(err)
	}

	users := d.RetroGetUsers(RetroID)

	return users
}

// RetroAbandon removes a user from the current retro by ID and sets abandoned true
func (d *Database) RetroAbandon(RetroID string, UserID string) ([]*model.RetroUser, error) {
	if _, err := d.db.Exec(
		`UPDATE retro_user SET active = false, abandoned = true WHERE retro_id = $1 AND user_id = $2`, RetroID, UserID); err != nil {
		log.Println(err)
		return nil, err
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		log.Println(err)
		return nil, err
	}

	users := d.RetroGetUsers(RetroID)

	return users, nil
}

// RetroSetOwner sets the ownerId for the retro
func (d *Database) RetroSetOwner(RetroID string, userID string, OwnerID string) (*model.Retro, error) {
	if _, err := d.db.Exec(
		`call set_retro_owner($1, $2);`, RetroID, OwnerID); err != nil {
		log.Println(err)
	}

	retro, err := d.RetroGet(RetroID)
	if err != nil {
		return nil, errors.New("Unable to promote owner")
	}

	return retro, nil
}

// RetroAdvancePhase sets the phase for the retro
func (d *Database) RetroAdvancePhase(RetroID string, Phase string) (*model.Retro, error) {
	if _, err := d.db.Exec(
		`call set_retro_phase($1, $2);`, RetroID, Phase); err != nil {
		log.Println(err)
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
		log.Println(err)
		return err
	}

	return nil
}

// GetRetroUserActiveStatus checks retro active status of User for given retro
func (d *Database) GetRetroUserActiveStatus(RetroID string, UserID string) error {
	var active bool

	e := d.db.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM retro_user
		WHERE user_id = $2 AND retro_id = $1;`,
		RetroID,
		UserID,
	).Scan(
		&active,
	)
	if e != nil {
		return e
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
			log.Println(err)
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
			log.Println(err)
		} else {
			retros = append(retros, b)
		}
	}

	return retros, Count, nil
}
