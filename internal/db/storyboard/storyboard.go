package storyboard

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

// Service represents a PostgreSQL implementation of thunderdome.StoryboardDataSvc.
type Service struct {
	DB         *sql.DB
	Logger     *otelzap.Logger
	AESHashKey string
}

// CreateStoryboard adds a new storyboard
func (d *Service) CreateStoryboard(ctx context.Context, OwnerID string, StoryboardName string, JoinCode string, FacilitatorCode string) (*thunderdome.Storyboard, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create storyboard encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create storyboard encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	var b = &thunderdome.Storyboard{
		Id:      "",
		OwnerID: OwnerID,
		Name:    StoryboardName,
		Users:   make([]*thunderdome.StoryboardUser, 0),
	}

	err := d.DB.QueryRowContext(ctx,
		`SELECT * FROM thunderdome.sb_create($1, $2, $3, $4, null);`,
		OwnerID,
		StoryboardName,
		encryptedJoinCode,
		encryptedFacilitatorCode,
	).Scan(&b.Id)
	if err != nil {
		return nil, fmt.Errorf("create storyboard query error: %v", err)
	}

	return b, nil
}

// TeamCreateStoryboard adds a new storyboard associated to a team
func (d *Service) TeamCreateStoryboard(ctx context.Context, TeamID string, OwnerID string, StoryboardName string, JoinCode string, FacilitatorCode string) (*thunderdome.Storyboard, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("team create storyboard encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("team create storyboard encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	var b = &thunderdome.Storyboard{
		Id:      "",
		OwnerID: OwnerID,
		Name:    StoryboardName,
		Users:   make([]*thunderdome.StoryboardUser, 0),
	}

	err := d.DB.QueryRowContext(ctx,
		`SELECT * FROM thunderdome.sb_create($1, $2, $3, $4, $5);`,
		OwnerID,
		StoryboardName,
		encryptedJoinCode,
		encryptedFacilitatorCode,
		TeamID,
	).Scan(&b.Id)
	if err != nil {
		return nil, fmt.Errorf("create storyboard query error: %v", err)
	}

	return b, nil
}

// EditStoryboard updates the storyboard by ID
func (d *Service) EditStoryboard(StoryboardID string, StoryboardName string, JoinCode string, FacilitatorCode string) error {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return fmt.Errorf("edit storyboard encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return fmt.Errorf("edit storyboard encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	if _, err := d.DB.Exec(`UPDATE thunderdome.storyboard
        SET name = $2, join_code = $3, facilitator_code = $4, updated_date = NOW()
        WHERE id = $1;`,
		StoryboardID, StoryboardName, encryptedJoinCode, encryptedFacilitatorCode,
	); err != nil {
		return fmt.Errorf("edit storyboard query error: %v", err)
	}

	return nil
}

// GetStoryboard gets a storyboard by ID
func (d *Service) GetStoryboard(StoryboardID string, UserID string) (*thunderdome.Storyboard, error) {
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
				FROM thunderdome.storyboard s
				LEFT JOIN thunderdome.storyboard_facilitator sf ON sf.storyboard_id = s.id
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
		return nil, fmt.Errorf("get storyboard query error: %v", e)
	}

	clErr := json.Unmarshal([]byte(cl), &b.ColorLegend)
	if clErr != nil {
		d.Logger.Error("color legend json error", zap.Error(clErr))
	}

	facilError := json.Unmarshal([]byte(facilitators), &b.Facilitators)
	if facilError != nil {
		d.Logger.Error("facilitators json error", zap.Error(facilError))
	}
	isFacilitator := db.Contains(b.Facilitators, UserID)

	b.Users = d.GetStoryboardUsers(StoryboardID)
	b.Goals = d.GetStoryboardGoals(StoryboardID)
	b.Personas = d.GetStoryboardPersonas(StoryboardID)

	if JoinCode != "" {
		DecryptedCode, codeErr := db.Decrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("get storyboard encrypt join_code error: %v", codeErr)
		}
		b.JoinCode = DecryptedCode
	}

	if FacilitatorCode != "" && isFacilitator {
		DecryptedCode, codeErr := db.Decrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("get storyboard decrypt facilitator_code error: %v", codeErr)
		}
		b.FacilitatorCode = DecryptedCode
	}

	return b, nil
}

// GetStoryboardsByUser gets a list of storyboards by UserID
func (d *Service) GetStoryboardsByUser(UserID string) ([]*thunderdome.Storyboard, int, error) {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	storyboardRows, storyboardsErr := d.DB.Query(`
		SELECT b.id, b.name, b.owner_id, b.created_date, b.updated_date
		FROM thunderdome.storyboard b
		LEFT JOIN thunderdome.storyboard_user su ON b.id = su.storyboard_id WHERE su.user_id = $1 AND su.abandoned = false
		GROUP BY b.id ORDER BY b.created_date DESC;
	`, UserID)
	if storyboardsErr != nil {
		return nil, 0, fmt.Errorf("get storyboards by user query error: %v", storyboardsErr)
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
func (d *Service) ConfirmStoryboardFacilitator(StoryboardID string, UserID string) error {
	var facilitatorId string
	var role string
	err := d.DB.QueryRow("SELECT type FROM thunderdome.users WHERE id = $1", UserID).Scan(&role)
	if err != nil {
		return fmt.Errorf("confirm storyboard facilitator user role query error:%v", err)
	}

	err = d.DB.QueryRow(
		`SELECT user_id FROM thunderdome.storyboard_facilitator WHERE storyboard_id = $1 AND user_id = $2;`,
		StoryboardID, UserID).Scan(&facilitatorId)
	if err != nil && role != "ADMIN" {
		return fmt.Errorf("confirm storyboard facilitator query error:%v", err)
	}

	return nil
}

// GetStoryboardUsers retrieves the users for a given storyboard from db
func (d *Service) GetStoryboardUsers(StoryboardID string) []*thunderdome.StoryboardUser {
	var users = make([]*thunderdome.StoryboardUser, 0)
	rows, err := d.DB.Query(
		`SELECT
			w.id, w.name, su.active, w.avatar, COALESCE(w.email, '')
		FROM thunderdome.storyboard_user su
		LEFT JOIN thunderdome.users w ON su.user_id = w.id
		WHERE su.storyboard_id = $1
		ORDER BY w.name;`,
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
					w.GravatarHash = db.CreateGravatarHash(w.GravatarHash)
				} else {
					w.GravatarHash = db.CreateGravatarHash(w.Id)
				}
				users = append(users, &w)
			}
		}
	}

	return users
}

// GetStoryboardPersonas retrieves the personas for a given storyboard from db
func (d *Service) GetStoryboardPersonas(StoryboardID string) []*thunderdome.StoryboardPersona {
	var personas = make([]*thunderdome.StoryboardPersona, 0)
	rows, err := d.DB.Query(
		`SELECT
			p.id, p.name, p.role, p.description
		FROM thunderdome.storyboard_persona p
		WHERE p.storyboard_id = $1;`,
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
func (d *Service) AddUserToStoryboard(StoryboardID string, UserID string) ([]*thunderdome.StoryboardUser, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_user (storyboard_id, user_id, active)
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
func (d *Service) RetreatStoryboardUser(StoryboardID string, UserID string) []*thunderdome.StoryboardUser {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_user SET active = false WHERE storyboard_id = $1 AND user_id = $2`, StoryboardID, UserID); err != nil {
		d.Logger.Error("set storyboard user active false error", zap.Error(err))
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("set user last active error", zap.Error(err))
	}

	users := d.GetStoryboardUsers(StoryboardID)

	return users
}

// GetStoryboardUserActiveStatus checks storyboard active status of User for given storyboard
func (d *Service) GetStoryboardUserActiveStatus(StoryboardID string, UserID string) error {
	var active bool

	err := d.DB.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM thunderdome.storyboard_user
		WHERE user_id = $2 AND storyboard_id = $1;`,
		StoryboardID,
		UserID,
	).Scan(
		&active,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get storyboard user active status query error: %v", err)
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if active {
		return errors.New("DUPLICATE_STORYBOARD_USER")
	}

	return nil
}

// AbandonStoryboard removes a user from the current storyboard by ID and sets abandoned true
func (d *Service) AbandonStoryboard(StoryboardID string, UserID string) ([]*thunderdome.StoryboardUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_user SET active = false, abandoned = true WHERE storyboard_id = $1 AND user_id = $2`, StoryboardID, UserID); err != nil {
		return nil, fmt.Errorf("set storyboard user active false error: %v", err)
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		return nil, fmt.Errorf("set user last active error: %v", err)
	}

	users := d.GetStoryboardUsers(StoryboardID)

	return users, nil
}

// StoryboardReviseColorLegend revises the storyboard color legend by StoryboardID
func (d *Service) StoryboardReviseColorLegend(StoryboardID string, UserID string, ColorLegend string) (*thunderdome.Storyboard, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard SET updated_date = NOW(), color_legend = $2 WHERE id = $1;`,
		StoryboardID,
		ColorLegend,
	); err != nil {
		return nil, fmt.Errorf("storyboard revise color legend query error: %v", err)
	}

	storyboard, err := d.GetStoryboard(StoryboardID, "")
	if err != nil {
		return nil, fmt.Errorf("storyboard revise color legend get storyboards query error: %v", err)
	}

	return storyboard, nil
}

// DeleteStoryboard removes all storyboard associations and the storyboard itself from DB by StoryboardID
func (d *Service) DeleteStoryboard(StoryboardID string, userID string) error {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard WHERE id = $1;`, StoryboardID); err != nil {
		return fmt.Errorf("storyboard delete query error: %v", err)
	}

	return nil
}

// AddStoryboardPersona adds a persona to a storyboard
func (d *Service) AddStoryboardPersona(StoryboardID string, UserID string, Name string, Role string, Description string) ([]*thunderdome.StoryboardPersona, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_persona (storyboard_id, name, role, description) VALUES ($1, $2, $3, $4);`,
		StoryboardID,
		Name,
		Role,
		Description,
	); err != nil {
		d.Logger.Error("CALL thunderdome.persona_add error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(StoryboardID)

	return personas, nil
}

// UpdateStoryboardPersona updates a storyboard persona
func (d *Service) UpdateStoryboardPersona(StoryboardID string, UserID string, PersonaID string, Name string, Role string, Description string) ([]*thunderdome.StoryboardPersona, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_persona SET name = $2, role = $3, description = $4, updated_date = NOW() WHERE id = $1;`,
		PersonaID,
		Name,
		Role,
		Description,
	); err != nil {
		d.Logger.Error("CALL thunderdome.persona_edit error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(StoryboardID)

	return personas, nil
}

// DeleteStoryboardPersona deletes a storyboard persona
func (d *Service) DeleteStoryboardPersona(StoryboardID string, UserID string, PersonaID string) ([]*thunderdome.StoryboardPersona, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard_persona WHERE id = $1;`,
		PersonaID,
	); err != nil {
		d.Logger.Error("CALL thunderdome.persona_delete error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(StoryboardID)

	return personas, nil
}

// GetStoryboards gets a list of storyboards
func (d *Service) GetStoryboards(Limit int, Offset int) ([]*thunderdome.Storyboard, int, error) {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	var Count int

	e := d.DB.QueryRow(
		"SELECT COUNT(*) FROM thunderdome.storyboard;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, fmt.Errorf("get storyboards count query error: %v", e)
	}

	rows, storyboardErr := d.DB.Query(`
		SELECT s.id, s.name, s.created_date, s.updated_date
		FROM thunderdome.storyboard s
		GROUP BY s.id ORDER BY s.created_date DESC
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if storyboardErr != nil {
		return nil, Count, fmt.Errorf("get storyboards query error: %v", storyboardErr)
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
func (d *Service) GetActiveStoryboards(Limit int, Offset int) ([]*thunderdome.Storyboard, int, error) {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	var Count int

	e := d.DB.QueryRow(
		"SELECT COUNT(DISTINCT su.storyboard_id) FROM thunderdome.storyboard_user su WHERE su.active IS TRUE;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, fmt.Errorf("get active storyboards count query error: %v", e)
	}

	rows, err := d.DB.Query(`
		SELECT s.id, s.name, s.created_date, s.updated_date
		FROM thunderdome.storyboard_user su
		LEFT JOIN thunderdome.storyboard s ON s.id = su.storyboard_id
		WHERE su.active IS TRUE GROUP BY s.id
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if err != nil {
		return nil, Count, fmt.Errorf("get active storyboards query error: %v", err)
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
func (d *Service) StoryboardFacilitatorAdd(StoryboardId string, UserID string) (*thunderdome.Storyboard, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_facilitator (storyboard_id, user_id) VALUES ($1, $2);`,
		StoryboardId, UserID); err != nil {
		return nil, fmt.Errorf("storyboard add faciliator query error: %v", err)
	}

	storyboard, err := d.GetStoryboard(StoryboardId, "")
	if err != nil {
		return nil, fmt.Errorf("storyboard add facilitator get storyboard error: %v", err)
	}

	return storyboard, nil
}

// StoryboardFacilitatorRemove removes a storyboard facilitator
func (d *Service) StoryboardFacilitatorRemove(StoryboardId string, UserID string) (*thunderdome.Storyboard, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard_facilitator WHERE storyboard_id = $1 AND user_id = $2;`,
		StoryboardId, UserID); err != nil {
		return nil, fmt.Errorf("storyboard remove facilitator query error: %v", err)
	}

	storyboard, err := d.GetStoryboard(StoryboardId, "")
	if err != nil {
		return nil, fmt.Errorf("storyboard remove facilitator get storyboard error: %v", err)
	}

	return storyboard, nil
}

// GetStoryboardFacilitatorCode retrieve the storyboard facilitator code
func (d *Service) GetStoryboardFacilitatorCode(StoryboardID string) (string, error) {
	var EncryptedCode string

	if err := d.DB.QueryRow(`
		SELECT COALESCE(facilitator_code, '') FROM thunderdome.storyboard
		WHERE id = $1`,
		StoryboardID,
	).Scan(&EncryptedCode); err != nil {
		return "", fmt.Errorf("get storyboard facilitator_code query error: %v", err)
	}

	if EncryptedCode == "" {
		return "", fmt.Errorf("storyboard facilitator_code not set")
	}
	DecryptedCode, codeErr := db.Decrypt(EncryptedCode, d.AESHashKey)
	if codeErr != nil {
		return "", fmt.Errorf("get storyboard facilitator_code decrypt error: %v", codeErr)
	}

	return DecryptedCode, nil
}

// CleanStoryboards deletes storyboards older than {DaysOld} days
func (d *Service) CleanStoryboards(ctx context.Context, DaysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.storyboard WHERE updated_date < (NOW() - $1 * interval '1 day');`,
		DaysOld,
	); err != nil {
		return fmt.Errorf("clean storyboards query error: %v", err)
	}

	return nil
}
