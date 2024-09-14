package storyboard

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
				s.id, s.name, s.owner_id, COALESCE(s.team_id::TEXT, ''), s.color_legend, COALESCE(s.join_code, ''), COALESCE(s.facilitator_code, ''),
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
		&b.TeamID,
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
func (d *Service) GetStoryboardsByUser(UserID string, Limit int, Offset int) ([]*thunderdome.Storyboard, int, error) {
	var Count int
	var storyboards = make([]*thunderdome.Storyboard, 0)

	e := d.DB.QueryRow(`
		WITH user_teams AS (
			SELECT t.id FROM thunderdome.team_user tu
			LEFT JOIN thunderdome.team t ON t.id = tu.team_id
			WHERE tu.user_id = $1
		),
		team_storyboards AS (
			SELECT id FROM thunderdome.storyboard WHERE team_id IN (SELECT id FROM user_teams)
		),
		user_storyboards AS (
			SELECT u.storyboard_id AS id FROM thunderdome.storyboard_user u
			WHERE u.user_id = $1 AND u.abandoned = false
		),
		storyboards AS (
			SELECT id from user_storyboards UNION SELECT id FROM team_storyboards
		)
		SELECT COUNT(*) FROM storyboards;
	`, UserID).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, fmt.Errorf("get storyboards by user count query error: %v", e)
	}

	storyboardRows, storyboardsErr := d.DB.Query(`
		WITH user_teams AS (
			SELECT t.id, t.name FROM thunderdome.team_user tu
			LEFT JOIN thunderdome.team t ON t.id = tu.team_id
			WHERE tu.user_id = $1
		),
		team_storyboards AS (
			SELECT id FROM thunderdome.storyboard WHERE team_id IN (SELECT id FROM user_teams)
		),
		user_storyboards AS (
			SELECT u.storyboard_id AS id FROM thunderdome.storyboard_user u
			WHERE u.user_id = $1 AND u.abandoned = false
		),
		storyboards AS (
			SELECT id from user_storyboards UNION SELECT id FROM team_storyboards
		)
		SELECT s.id, s.name, s.owner_id, COALESCE(s.team_id::TEXT, ''), s.created_date, s.updated_date,
		  min(COALESCE(t.name, '')) as team_name
		FROM thunderdome.storyboard s
		LEFT JOIN user_teams t ON t.id = s.team_id
		WHERE s.id IN (SELECT id FROM storyboards)
		GROUP BY s.id ORDER BY s.created_date DESC LIMIT $2 OFFSET $3;
	`, UserID, Limit, Offset)
	if storyboardsErr != nil {
		return nil, Count, fmt.Errorf("get storyboards by user query error: %v", storyboardsErr)
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
			&b.TeamID,
			&b.CreatedDate,
			&b.UpdatedDate,
			&b.TeamName,
		); err != nil {
			d.Logger.Error("get_storyboards_by_user query scan error", zap.Error(err))
		} else {
			storyboards = append(storyboards, b)
		}
	}

	return storyboards, Count, nil
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
		SELECT s.id, s.name, COALESCE(s.team_id::TEXT, ''), s.created_date, s.updated_date
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
			&b.TeamID,
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
		SELECT s.id, s.name, COALESCE(s.team_id::TEXT, ''), s.created_date, s.updated_date
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
			&b.TeamID,
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
