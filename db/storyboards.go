package db

import (
	"encoding/json"
	"errors"
	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

//CreateStoryboard adds a new storyboard to the db
func (d *Database) CreateStoryboard(OwnerID string, StoryboardName string, JoinCode string) (*model.Storyboard, error) {
	var encryptedJoinCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, codeErr
		}
		encryptedJoinCode = EncryptedCode
	}

	var b = &model.Storyboard{
		StoryboardID:   "",
		OwnerID:        OwnerID,
		StoryboardName: StoryboardName,
		Users:          make([]*model.StoryboardUser, 0),
	}

	e := d.db.QueryRow(
		`SELECT * FROM create_storyboard($1, $2, $3);`,
		OwnerID,
		StoryboardName,
		encryptedJoinCode,
	).Scan(&b.StoryboardID)
	if e != nil {
		d.logger.Error("create_storyboard query error", zap.Error(e))
		return nil, errors.New("error creating storyboard")
	}

	// if a join code is set than add owner to retro_user
	// this prevents them from having to enter join code on initial create to enter
	if JoinCode != "" {
		if _, err := d.db.Exec(
			`INSERT INTO storyboard_user (storyboard_id, user_id)
		VALUES ($1, $2)`,
			b.StoryboardID,
			OwnerID,
		); err != nil {
			d.logger.Error("insert storyboard user error", zap.Error(err))
		}
	}

	return b, nil
}

// EditStoryboard updates the storyboard by ID
func (d *Database) EditStoryboard(StoryboardID string, StoryboardName string, JoinCode string) error {
	var encryptedJoinCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return errors.New("unable to revise storyboard join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if _, err := d.db.Exec(`call edit_storyboard($1, $2, $3);`,
		StoryboardID, StoryboardName, encryptedJoinCode,
	); err != nil {
		d.logger.Error("update storyboard error", zap.Error(err))
		return errors.New("unable to edit storyboard")
	}

	return nil
}

// GetStoryboard gets a storyboard by ID
func (d *Database) GetStoryboard(StoryboardID string) (*model.Storyboard, error) {
	var cl string
	var JoinCode string
	var b = &model.Storyboard{
		StoryboardID:   StoryboardID,
		OwnerID:        "",
		StoryboardName: "",
		Users:          make([]*model.StoryboardUser, 0),
		Goals:          make([]*model.StoryboardGoal, 0),
		ColorLegend:    make([]*model.Color, 0),
		Personas:       make([]*model.StoryboardPersona, 0),
	}

	// get storyboard
	e := d.db.QueryRow(
		`SELECT id, name, owner_id, color_legend, COALESCE(join_code, ''), created_date, updated_date FROM storyboard WHERE id = $1`,
		StoryboardID,
	).Scan(
		&b.StoryboardID,
		&b.StoryboardName,
		&b.OwnerID,
		&cl,
		&JoinCode,
		&b.CreatedDate,
		&b.UpdatedDate,
	)
	if e != nil {
		d.logger.Error("get storyboard query error", zap.Error(e))
		return nil, errors.New("Not found")
	}

	clErr := json.Unmarshal([]byte(cl), &b.ColorLegend)
	if clErr != nil {
		d.logger.Error("color legend json error", zap.Error(clErr))
	}

	b.Users = d.GetStoryboardUsers(StoryboardID)
	b.Goals = d.GetStoryboardGoals(StoryboardID)
	b.Personas = d.GetStoryboardPersonas(StoryboardID)

	if JoinCode != "" {
		DecryptedCode, codeErr := decrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, errors.New("unable to decode join_code")
		}
		b.JoinCode = DecryptedCode
	}

	return b, nil
}

// GetStoryboardsByUser gets a list of storyboards by UserID
func (d *Database) GetStoryboardsByUser(UserID string) ([]*model.Storyboard, int, error) {
	var storyboards = make([]*model.Storyboard, 0)
	storyboardRows, storyboardsErr := d.db.Query(`
		SELECT * FROM get_storyboards_by_user($1);
	`, UserID)
	if storyboardsErr != nil {
		return nil, 0, errors.New("Not found")
	}

	defer storyboardRows.Close()
	for storyboardRows.Next() {
		var b = &model.Storyboard{
			StoryboardID:   "",
			OwnerID:        "",
			StoryboardName: "",
			Users:          make([]*model.StoryboardUser, 0),
		}
		if err := storyboardRows.Scan(
			&b.StoryboardID,
			&b.StoryboardName,
			&b.OwnerID,
			&b.CreatedDate,
			&b.UpdatedDate,
		); err != nil {
			d.logger.Error("get_storyboards_by_user query scan error", zap.Error(err))
		} else {
			storyboards = append(storyboards, b)
		}
	}

	return storyboards, 0, nil
}

// ConfirmStoryboardOwner confirms the user is infact owner of the storyboard
func (d *Database) ConfirmStoryboardOwner(StoryboardID string, userID string) error {
	var ownerID string
	e := d.db.QueryRow("SELECT owner_id FROM storyboard WHERE id = $1", StoryboardID).Scan(&ownerID)
	if e != nil {
		d.logger.Error("get owner_id from storyboard query error", zap.Error(e))
		return errors.New("Storyboard Not found")
	}

	if ownerID != userID {
		return errors.New("Not Owner")
	}

	return nil
}

// GetStoryboardUser gets a user from db by ID and checks storyboard active status
func (d *Database) GetStoryboardUser(StoryboardID string, UserID string) (*model.StoryboardUser, error) {
	var active bool
	var w model.StoryboardUser

	e := d.db.QueryRow(
		`SELECT * FROM get_storyboard_user($1, $2);`,
		StoryboardID,
		UserID,
	).Scan(
		&w.UserID,
		&w.UserName,
		&active,
	)
	if e != nil {
		d.logger.Error("get_storyboard_user query error", zap.Error(e))
		return nil, errors.New("User Not found")
	}

	if active {
		return nil, errors.New("User Already Active in Storyboard")
	}

	return &w, nil
}

// GetStoryboardUsers retrieves the users for a given storyboard from db
func (d *Database) GetStoryboardUsers(StoryboardID string) []*model.StoryboardUser {
	var users = make([]*model.StoryboardUser, 0)
	rows, err := d.db.Query(
		`SELECT * FROM get_storyboard_users($1);`,
		StoryboardID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w model.StoryboardUser
			if err := rows.Scan(&w.UserID, &w.UserName, &w.Active, &w.Avatar, &w.GravatarHash); err != nil {
				d.logger.Error("get_storyboard_users query scan error", zap.Error(err))
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

// GetStoryboardPersonas retrieves the personas for a given storyboard from db
func (d *Database) GetStoryboardPersonas(StoryboardID string) []*model.StoryboardPersona {
	var personas = make([]*model.StoryboardPersona, 0)
	rows, err := d.db.Query(
		`SELECT * FROM get_storyboard_personas($1);`,
		StoryboardID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p model.StoryboardPersona
			if err := rows.Scan(&p.PersonaID, &p.Name, &p.Role, &p.Description); err != nil {
				d.logger.Error("get_storyboard_personas query scan error", zap.Error(err))
			} else {
				personas = append(personas, &p)
			}
		}
	}

	return personas
}

// AddUserToStoryboard adds a user by ID to the storyboard by ID
func (d *Database) AddUserToStoryboard(StoryboardID string, UserID string) ([]*model.StoryboardUser, error) {
	if _, err := d.db.Exec(
		`INSERT INTO storyboard_user (storyboard_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (storyboard_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		StoryboardID,
		UserID,
	); err != nil {
		d.logger.Error("insert storybaord user error", zap.Error(err))
	}

	users := d.GetStoryboardUsers(StoryboardID)

	return users, nil
}

// RetreatStoryboardUser removes a user from the current storyboard by ID
func (d *Database) RetreatStoryboardUser(StoryboardID string, UserID string) []*model.StoryboardUser {
	if _, err := d.db.Exec(
		`UPDATE storyboard_user SET active = false WHERE storyboard_id = $1 AND user_id = $2`, StoryboardID, UserID); err != nil {
		d.logger.Error("set storyboard user active false error", zap.Error(err))
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.logger.Error("set user last active error", zap.Error(err))
	}

	users := d.GetStoryboardUsers(StoryboardID)

	return users
}

// GetStoryboardUserActiveStatus checks storyboard active status of User for given storyboard
func (d *Database) GetStoryboardUserActiveStatus(StoryboardID string, UserID string) error {
	var active bool

	err := d.db.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM storyboard_user
		WHERE user_id = $2 AND storyboard_id = $1;`,
		StoryboardID,
		UserID,
	).Scan(
		&active,
	)
	if err != nil {
		d.logger.Error("get storyboard user active status error", zap.Error(err))
		return err
	}

	if active {
		return errors.New("DUPLICATE_STORYBOARD_USER")
	}

	return nil
}

// AbandonStoryboard removes a user from the current storyboard by ID and sets abandoned true
func (d *Database) AbandonStoryboard(StoryboardID string, UserID string) ([]*model.StoryboardUser, error) {
	if _, err := d.db.Exec(
		`UPDATE storyboard_user SET active = false, abandoned = true WHERE storyboard_id = $1 AND user_id = $2`, StoryboardID, UserID); err != nil {
		d.logger.Error("set storyboard user active false error", zap.Error(err))
		return nil, err
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.logger.Error("set user last active error", zap.Error(err))
		return nil, err
	}

	users := d.GetStoryboardUsers(StoryboardID)

	return users, nil
}

// SetStoryboardOwner sets the ownerId for the storyboard
func (d *Database) SetStoryboardOwner(StoryboardID string, userID string, OwnerID string) (*model.Storyboard, error) {
	err := d.ConfirmStoryboardOwner(StoryboardID, userID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call set_storyboard_owner($1, $2);`, StoryboardID, OwnerID); err != nil {
		d.logger.Error("call set_storyboard_owner error", zap.Error(err))
	}

	storyboard, err := d.GetStoryboard(StoryboardID)
	if err != nil {
		return nil, errors.New("Unable to promote owner")
	}

	return storyboard, nil
}

// StoryboardReviseColorLegend revises the storyboard color legend by StoryboardID
func (d *Database) StoryboardReviseColorLegend(StoryboardID string, UserID string, ColorLegend string) (*model.Storyboard, error) {
	err := d.ConfirmStoryboardOwner(StoryboardID, UserID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call revise_color_legend($1, $2);`,
		StoryboardID,
		ColorLegend,
	); err != nil {
		d.logger.Error("call revise_color_legend error", zap.Error(err))
		return nil, err
	}

	storyboard, err := d.GetStoryboard(StoryboardID)
	if err != nil {
		return nil, errors.New("Unable to promote owner")
	}

	return storyboard, nil
}

// DeleteStoryboard removes all storyboard associations and the storyboard itself from DB by StoryboardID
func (d *Database) DeleteStoryboard(StoryboardID string, userID string) error {
	err := d.ConfirmStoryboardOwner(StoryboardID, userID)
	if err != nil {
		return errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call delete_storyboard($1);`, StoryboardID); err != nil {
		d.logger.Error("call delete_storyboard error", zap.Error(err))
		return err
	}

	return nil
}

// AddStoryboardPersona adds a persona to a storyboard
func (d *Database) AddStoryboardPersona(StoryboardID string, UserID string, Name string, Role string, Description string) ([]*model.StoryboardPersona, error) {
	err := d.ConfirmStoryboardOwner(StoryboardID, UserID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call persona_add($1, $2, $3, $4);`,
		StoryboardID,
		Name,
		Role,
		Description,
	); err != nil {
		d.logger.Error("call persona_add error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(StoryboardID)

	return personas, nil
}

// UpdateStoryboardPersona updates a storyboard persona
func (d *Database) UpdateStoryboardPersona(StoryboardID string, UserID string, PersonaID string, Name string, Role string, Description string) ([]*model.StoryboardPersona, error) {
	err := d.ConfirmStoryboardOwner(StoryboardID, UserID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call persona_edit($1, $2, $3, $4, $5);`,
		StoryboardID,
		PersonaID,
		Name,
		Role,
		Description,
	); err != nil {
		d.logger.Error("call persona_edit error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(StoryboardID)

	return personas, nil
}

// DeleteStoryboardPersona deletes a storyboard persona
func (d *Database) DeleteStoryboardPersona(StoryboardID string, UserID string, PersonaID string) ([]*model.StoryboardPersona, error) {
	err := d.ConfirmStoryboardOwner(StoryboardID, UserID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call persona_delete($1, $2);`,
		StoryboardID,
		PersonaID,
	); err != nil {
		d.logger.Error("call persona_delete error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(StoryboardID)

	return personas, nil
}

// GetStoryboards gets a list of storyboards
func (d *Database) GetStoryboards(Limit int, Offset int) ([]*model.Storyboard, int, error) {
	var storyboards = make([]*model.Storyboard, 0)
	var Count int

	e := d.db.QueryRow(
		"SELECT COUNT(*) FROM storyboard;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	rows, storyboardErr := d.db.Query(`
		SELECT s.id, s.name, s.created_date, s.updated_date
		FROM storyboard s
		GROUP BY s.id ORDER BY s.created_date DESC
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if storyboardErr != nil {
		return nil, Count, storyboardErr
	}

	defer rows.Close()
	for rows.Next() {
		var b = &model.Storyboard{
			Users: make([]*model.StoryboardUser, 0),
		}
		if err := rows.Scan(
			&b.StoryboardID,
			&b.StoryboardName,
			&b.CreatedDate,
			&b.UpdatedDate,
		); err != nil {
			d.logger.Error("get storyboards error", zap.Error(err))
		} else {
			storyboards = append(storyboards, b)
		}
	}

	return storyboards, Count, nil
}

// GetActiveStoryboards gets a list of active storyboards
func (d *Database) GetActiveStoryboards(Limit int, Offset int) ([]*model.Storyboard, int, error) {
	var storyboards = make([]*model.Storyboard, 0)
	var Count int

	e := d.db.QueryRow(
		"SELECT COUNT(DISTINCT su.storyboard_id) FROM storyboard_user su WHERE su.active IS TRUE;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	rows, err := d.db.Query(`
		SELECT s.id, s.name, s.created_date, s.updated_date
		FROM storyboard_user su
		LEFT JOIN storyboard s ON s.id = su.storyboard_id
		WHERE su.active IS TRUE GROUP BY s.id
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if err != nil {
		return nil, Count, err
	}

	defer rows.Close()
	for rows.Next() {
		var b = &model.Storyboard{
			Users: make([]*model.StoryboardUser, 0),
		}
		if err := rows.Scan(
			&b.StoryboardID,
			&b.StoryboardName,
			&b.CreatedDate,
			&b.UpdatedDate,
		); err != nil {
			d.logger.Error("get active storyboards error", zap.Error(err))
		} else {
			storyboards = append(storyboards, b)
		}
	}

	return storyboards, Count, nil
}
