package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// StoryboardService represents a PostgreSQL implementation of thunderdome.StoryboardService.
type StoryboardService struct {
	DB         *sql.DB
	Logger     *otelzap.Logger
	AESHashKey string
}

// CreateStoryboard adds a new storyboard
func (d *StoryboardService) CreateStoryboard(ctx context.Context, OwnerID string, StoryboardName string, JoinCode string, FacilitatorCode string) (*thunderdome.Storyboard, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, codeErr
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, codeErr
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	var b = &thunderdome.Storyboard{
		Id:      "",
		OwnerID: OwnerID,
		Name:    StoryboardName,
		Users:   make([]*thunderdome.StoryboardUser, 0),
	}

	e := d.DB.QueryRowContext(ctx,
		`SELECT * FROM create_storyboard($1, $2, $3, $4);`,
		OwnerID,
		StoryboardName,
		encryptedJoinCode,
		encryptedFacilitatorCode,
	).Scan(&b.Id)
	if e != nil {
		d.Logger.Error("create_storyboard query error", zap.Error(e))
		return nil, errors.New("error creating storyboard")
	}

	return b, nil
}

// TeamCreateStoryboard adds a new storyboard associated to a team
func (d *StoryboardService) TeamCreateStoryboard(ctx context.Context, TeamID string, OwnerID string, StoryboardName string, JoinCode string, FacilitatorCode string) (*thunderdome.Storyboard, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, codeErr
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, codeErr
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	var b = &thunderdome.Storyboard{
		Id:      "",
		OwnerID: OwnerID,
		Name:    StoryboardName,
		Users:   make([]*thunderdome.StoryboardUser, 0),
	}

	e := d.DB.QueryRowContext(ctx,
		`SELECT * FROM team_create_storyboard($1, $2, $3, $4, $5);`,
		TeamID,
		OwnerID,
		StoryboardName,
		encryptedJoinCode,
		encryptedFacilitatorCode,
	).Scan(&b.Id)
	if e != nil {
		d.Logger.Error("team_create_storyboard query error", zap.Error(e))
		return nil, errors.New("error creating storyboard")
	}

	return b, nil
}

// EditStoryboard updates the storyboard by ID
func (d *StoryboardService) EditStoryboard(StoryboardID string, StoryboardName string, JoinCode string, FacilitatorCode string) error {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return errors.New("unable to revise storyboard join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return errors.New("unable to revise storyboard facilitator_code")
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	if _, err := d.DB.Exec(`call edit_storyboard($1, $2, $3, $4);`,
		StoryboardID, StoryboardName, encryptedJoinCode, encryptedFacilitatorCode,
	); err != nil {
		d.Logger.Error("update storyboard error", zap.Error(err))
		return errors.New("unable to edit storyboard")
	}

	return nil
}

// GetStoryboard gets a storyboard by ID
func (d *StoryboardService) GetStoryboard(StoryboardID string, UserID string) (*thunderdome.Storyboard, error) {
	var cl string
	var JoinCode string
	var facilitators string
	var FacilitatorCode string
	var b = &thunderdome.Storyboard{
		Id:          StoryboardID,
		OwnerID:     "",
		Name:        "",
		Users:       make([]*thunderdome.StoryboardUser, 0),
		Goals:       make([]*thunderdome.StoryboardGoal, 0),
		ColorLegend: make([]*thunderdome.Color, 0),
		Personas:    make([]*thunderdome.StoryboardPersona, 0),
	}

	// get storyboard
	e := d.DB.QueryRow(
		`SELECT
				s.id, s.name, s.owner_id, s.color_legend, COALESCE(s.join_code, ''), COALESCE(s.facilitator_code, ''),
				 s.created_date, s.updated_date,
				COALESCE(json_agg(sf.user_id) FILTER (WHERE sf.storyboard_id IS NOT NULL), '[]') AS facilitators
				FROM storyboard s
				LEFT JOIN storyboard_facilitator sf ON sf.storyboard_id = s.id
				WHERE s.id = $1
				GROUP BY s.id`,
		StoryboardID,
	).Scan(
		&b.Id,
		&b.Name,
		&b.OwnerID,
		&cl,
		&JoinCode,
		&FacilitatorCode,
		&b.CreatedDate,
		&b.UpdatedDate,
		&facilitators,
	)
	if e != nil {
		d.Logger.Error("get storyboard query error", zap.Error(e))
		return nil, errors.New("Not found")
	}

	clErr := json.Unmarshal([]byte(cl), &b.ColorLegend)
	if clErr != nil {
		d.Logger.Error("color legend json error", zap.Error(clErr))
	}

	facilError := json.Unmarshal([]byte(facilitators), &b.Facilitators)
	if facilError != nil {
		d.Logger.Error("facilitators json error", zap.Error(facilError))
	}
	isFacilitator := contains(b.Facilitators, UserID)

	b.Users = d.GetStoryboardUsers(StoryboardID)
	b.Goals = d.GetStoryboardGoals(StoryboardID)
	b.Personas = d.GetStoryboardPersonas(StoryboardID)

	if JoinCode != "" {
		DecryptedCode, codeErr := decrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to decode join_code")
		}
		b.JoinCode = DecryptedCode
	}

	if FacilitatorCode != "" && isFacilitator {
		DecryptedCode, codeErr := decrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to decode facilitator_code")
		}
		b.FacilitatorCode = DecryptedCode
	}

	return b, nil
}

// GetStoryboardsByUser gets a list of storyboards by UserID
func (d *StoryboardService) GetStoryboardsByUser(UserID string) ([]*thunderdome.Storyboard, int, error) {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	storyboardRows, storyboardsErr := d.DB.Query(`
		SELECT * FROM get_storyboards_by_user($1);
	`, UserID)
	if storyboardsErr != nil {
		return nil, 0, errors.New("Not found")
	}

	defer storyboardRows.Close()
	for storyboardRows.Next() {
		var b = &thunderdome.Storyboard{
			Id:      "",
			OwnerID: "",
			Name:    "",
			Users:   make([]*thunderdome.StoryboardUser, 0),
		}
		if err := storyboardRows.Scan(
			&b.Id,
			&b.Name,
			&b.OwnerID,
			&b.CreatedDate,
			&b.UpdatedDate,
		); err != nil {
			d.Logger.Error("get_storyboards_by_user query scan error", zap.Error(err))
		} else {
			storyboards = append(storyboards, b)
		}
	}

	return storyboards, 0, nil
}

// ConfirmStoryboardFacilitator confirms the user is a facilitator of the storyboard
func (d *StoryboardService) ConfirmStoryboardFacilitator(StoryboardID string, UserID string) error {
	var facilitatorId string
	var role string
	err := d.DB.QueryRow("SELECT type FROM users WHERE id = $1", UserID).Scan(&role)
	if err != nil {
		d.Logger.Error("error getting user role", zap.Error(err))
		return errors.New("unable to get user role")
	}

	err = d.DB.QueryRow(
		`SELECT user_id FROM storyboard_facilitator WHERE storyboard_id = $1 AND user_id = $2;`,
		StoryboardID, UserID).Scan(&facilitatorId)
	if err != nil && role != "ADMIN" {
		d.Logger.Error("get ConfirmStoryboardFacilitator error", zap.Error(err))
		return errors.New("storyboard facilitator not found")
	}

	return nil
}

// GetStoryboardUsers retrieves the users for a given storyboard from db
func (d *StoryboardService) GetStoryboardUsers(StoryboardID string) []*thunderdome.StoryboardUser {
	var users = make([]*thunderdome.StoryboardUser, 0)
	rows, err := d.DB.Query(
		`SELECT * FROM get_storyboard_users($1);`,
		StoryboardID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w thunderdome.StoryboardUser
			if err := rows.Scan(&w.Id, &w.Name, &w.Active, &w.Avatar, &w.GravatarHash); err != nil {
				d.Logger.Error("get_storyboard_users query scan error", zap.Error(err))
			} else {
				if w.GravatarHash != "" {
					w.GravatarHash = createGravatarHash(w.GravatarHash)
				} else {
					w.GravatarHash = createGravatarHash(w.Id)
				}
				users = append(users, &w)
			}
		}
	}

	return users
}

// GetStoryboardPersonas retrieves the personas for a given storyboard from db
func (d *StoryboardService) GetStoryboardPersonas(StoryboardID string) []*thunderdome.StoryboardPersona {
	var personas = make([]*thunderdome.StoryboardPersona, 0)
	rows, err := d.DB.Query(
		`SELECT * FROM get_storyboard_personas($1);`,
		StoryboardID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p thunderdome.StoryboardPersona
			if err := rows.Scan(&p.Id, &p.Name, &p.Role, &p.Description); err != nil {
				d.Logger.Error("get_storyboard_personas query scan error", zap.Error(err))
			} else {
				personas = append(personas, &p)
			}
		}
	}

	return personas
}

// AddUserToStoryboard adds a user by ID to the storyboard by ID
func (d *StoryboardService) AddUserToStoryboard(StoryboardID string, UserID string) ([]*thunderdome.StoryboardUser, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO storyboard_user (storyboard_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (storyboard_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		StoryboardID,
		UserID,
	); err != nil {
		d.Logger.Error("insert storybaord user error", zap.Error(err))
	}

	users := d.GetStoryboardUsers(StoryboardID)

	return users, nil
}

// RetreatStoryboardUser removes a user from the current storyboard by ID
func (d *StoryboardService) RetreatStoryboardUser(StoryboardID string, UserID string) []*thunderdome.StoryboardUser {
	if _, err := d.DB.Exec(
		`UPDATE storyboard_user SET active = false WHERE storyboard_id = $1 AND user_id = $2`, StoryboardID, UserID); err != nil {
		d.Logger.Error("set storyboard user active false error", zap.Error(err))
	}

	if _, err := d.DB.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("set user last active error", zap.Error(err))
	}

	users := d.GetStoryboardUsers(StoryboardID)

	return users
}

// GetStoryboardUserActiveStatus checks storyboard active status of User for given storyboard
func (d *StoryboardService) GetStoryboardUserActiveStatus(StoryboardID string, UserID string) error {
	var active bool

	err := d.DB.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM storyboard_user
		WHERE user_id = $2 AND storyboard_id = $1;`,
		StoryboardID,
		UserID,
	).Scan(
		&active,
	)
	if err != nil {
		return err
	}

	if active {
		return errors.New("DUPLICATE_STORYBOARD_USER")
	}

	return nil
}

// AbandonStoryboard removes a user from the current storyboard by ID and sets abandoned true
func (d *StoryboardService) AbandonStoryboard(StoryboardID string, UserID string) ([]*thunderdome.StoryboardUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE storyboard_user SET active = false, abandoned = true WHERE storyboard_id = $1 AND user_id = $2`, StoryboardID, UserID); err != nil {
		d.Logger.Error("set storyboard user active false error", zap.Error(err))
		return nil, err
	}

	if _, err := d.DB.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("set user last active error", zap.Error(err))
		return nil, err
	}

	users := d.GetStoryboardUsers(StoryboardID)

	return users, nil
}

// StoryboardReviseColorLegend revises the storyboard color legend by StoryboardID
func (d *StoryboardService) StoryboardReviseColorLegend(StoryboardID string, UserID string, ColorLegend string) (*thunderdome.Storyboard, error) {
	if _, err := d.DB.Exec(
		`call revise_color_legend($1, $2);`,
		StoryboardID,
		ColorLegend,
	); err != nil {
		d.Logger.Error("call revise_color_legend error", zap.Error(err))
		return nil, err
	}

	storyboard, err := d.GetStoryboard(StoryboardID, "")
	if err != nil {
		return nil, errors.New("Unable to promote owner")
	}

	return storyboard, nil
}

// DeleteStoryboard removes all storyboard associations and the storyboard itself from DB by StoryboardID
func (d *StoryboardService) DeleteStoryboard(StoryboardID string, userID string) error {
	if _, err := d.DB.Exec(
		`call delete_storyboard($1);`, StoryboardID); err != nil {
		d.Logger.Error("call delete_storyboard error", zap.Error(err))
		return err
	}

	return nil
}

// AddStoryboardPersona adds a persona to a storyboard
func (d *StoryboardService) AddStoryboardPersona(StoryboardID string, UserID string, Name string, Role string, Description string) ([]*thunderdome.StoryboardPersona, error) {
	if _, err := d.DB.Exec(
		`call persona_add($1, $2, $3, $4);`,
		StoryboardID,
		Name,
		Role,
		Description,
	); err != nil {
		d.Logger.Error("call persona_add error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(StoryboardID)

	return personas, nil
}

// UpdateStoryboardPersona updates a storyboard persona
func (d *StoryboardService) UpdateStoryboardPersona(StoryboardID string, UserID string, PersonaID string, Name string, Role string, Description string) ([]*thunderdome.StoryboardPersona, error) {
	if _, err := d.DB.Exec(
		`call persona_edit($1, $2, $3, $4, $5);`,
		StoryboardID,
		PersonaID,
		Name,
		Role,
		Description,
	); err != nil {
		d.Logger.Error("call persona_edit error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(StoryboardID)

	return personas, nil
}

// DeleteStoryboardPersona deletes a storyboard persona
func (d *StoryboardService) DeleteStoryboardPersona(StoryboardID string, UserID string, PersonaID string) ([]*thunderdome.StoryboardPersona, error) {
	if _, err := d.DB.Exec(
		`call persona_delete($1, $2);`,
		StoryboardID,
		PersonaID,
	); err != nil {
		d.Logger.Error("call persona_delete error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(StoryboardID)

	return personas, nil
}

// GetStoryboards gets a list of storyboards
func (d *StoryboardService) GetStoryboards(Limit int, Offset int) ([]*thunderdome.Storyboard, int, error) {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	var Count int

	e := d.DB.QueryRow(
		"SELECT COUNT(*) FROM storyboard;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	rows, storyboardErr := d.DB.Query(`
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
		var b = &thunderdome.Storyboard{
			Users: make([]*thunderdome.StoryboardUser, 0),
		}
		if err := rows.Scan(
			&b.Id,
			&b.Name,
			&b.CreatedDate,
			&b.UpdatedDate,
		); err != nil {
			d.Logger.Error("get storyboards error", zap.Error(err))
		} else {
			storyboards = append(storyboards, b)
		}
	}

	return storyboards, Count, nil
}

// GetActiveStoryboards gets a list of active storyboards
func (d *StoryboardService) GetActiveStoryboards(Limit int, Offset int) ([]*thunderdome.Storyboard, int, error) {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	var Count int

	e := d.DB.QueryRow(
		"SELECT COUNT(DISTINCT su.storyboard_id) FROM storyboard_user su WHERE su.active IS TRUE;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	rows, err := d.DB.Query(`
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
		var b = &thunderdome.Storyboard{
			Users: make([]*thunderdome.StoryboardUser, 0),
		}
		if err := rows.Scan(
			&b.Id,
			&b.Name,
			&b.CreatedDate,
			&b.UpdatedDate,
		); err != nil {
			d.Logger.Error("get active storyboards error", zap.Error(err))
		} else {
			storyboards = append(storyboards, b)
		}
	}

	return storyboards, Count, nil
}

// StoryboardFacilitatorAdd adds a storyboard facilitator
func (d *StoryboardService) StoryboardFacilitatorAdd(StoryboardId string, UserID string) (*thunderdome.Storyboard, error) {
	if _, err := d.DB.Exec(
		`call sb_facilitator_add($1, $2);`, StoryboardId, UserID); err != nil {
		d.Logger.Error("call sb_facilitator_add error", zap.Error(err))
		return nil, errors.New("unable to add facilitator")
	}

	storyboard, err := d.GetStoryboard(StoryboardId, "")
	if err != nil {
		return nil, err
	}

	return storyboard, nil
}

// StoryboardFacilitatorRemove removes a storyboard facilitator
func (d *StoryboardService) StoryboardFacilitatorRemove(StoryboardId string, UserID string) (*thunderdome.Storyboard, error) {
	if _, err := d.DB.Exec(
		`call sb_facilitator_remove($1, $2);`, StoryboardId, UserID); err != nil {
		d.Logger.Error("call sb_facilitator_remove error", zap.Error(err))
		return nil, errors.New("unable to remove facilitator")
	}

	storyboard, err := d.GetStoryboard(StoryboardId, "")
	if err != nil {
		return nil, err
	}

	return storyboard, nil
}

// GetStoryboardFacilitatorCode retrieve the storyboard facilitator code
func (d *StoryboardService) GetStoryboardFacilitatorCode(StoryboardID string) (string, error) {
	var EncryptedCode string

	if err := d.DB.QueryRow(`
		SELECT COALESCE(facilitator_code, '') FROM storyboard
		WHERE id = $1`,
		StoryboardID,
	).Scan(&EncryptedCode); err != nil {
		d.Logger.Error("get retro facilitator_code error", zap.Error(err))
		return "", errors.New("unable to retrieve storyboard facilitator_code")
	}

	if EncryptedCode == "" {
		return "", errors.New("unable to retrieve storyboard facilitator_code")
	}
	DecryptedCode, codeErr := decrypt(EncryptedCode, d.AESHashKey)
	if codeErr != nil {
		return "", errors.New("unable to retrieve storyboard facilitator_code")
	}

	return DecryptedCode, nil
}

// CleanStoryboards deletes storyboards older than {DaysOld} days
func (d *StoryboardService) CleanStoryboards(ctx context.Context, DaysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`call clean_storyboards($1);`,
		DaysOld,
	); err != nil {
		d.Logger.Ctx(ctx).Error("call clean_storyboards", zap.Error(err))
		return errors.New("error attempting to clean storyboards")
	}

	return nil
}
