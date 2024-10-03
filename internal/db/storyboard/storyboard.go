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

// Service represents the storyboard database service
type Service struct {
	DB         *sql.DB
	Logger     *otelzap.Logger
	AESHashKey string
}

// CreateStoryboard adds a new storyboard
func (d *Service) CreateStoryboard(ctx context.Context, ownerID string, storyboardName string, joinCode string, facilitatorCode string) (*thunderdome.Storyboard, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if joinCode != "" {
		encryptedCode, codeErr := db.Encrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create storyboard encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = encryptedCode
	}

	if facilitatorCode != "" {
		encryptedCode, codeErr := db.Encrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create storyboard encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = encryptedCode
	}

	var b = &thunderdome.Storyboard{
		ID:      "",
		OwnerID: ownerID,
		Name:    storyboardName,
		Users:   make([]*thunderdome.StoryboardUser, 0),
	}

	err := d.DB.QueryRowContext(ctx,
		`SELECT * FROM thunderdome.sb_create($1, $2, $3, $4, null);`,
		ownerID,
		storyboardName,
		encryptedJoinCode,
		encryptedFacilitatorCode,
	).Scan(&b.ID)
	if err != nil {
		return nil, fmt.Errorf("create storyboard query error: %v", err)
	}

	return b, nil
}

// TeamCreateStoryboard adds a new storyboard associated to a team
func (d *Service) TeamCreateStoryboard(ctx context.Context, teamID string, ownerID string, storyboardName string, joinCode string, facilitatorCode string) (*thunderdome.Storyboard, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if joinCode != "" {
		encryptedCode, codeErr := db.Encrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("team create storyboard encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = encryptedCode
	}

	if facilitatorCode != "" {
		encryptedCode, codeErr := db.Encrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("team create storyboard encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = encryptedCode
	}

	var b = &thunderdome.Storyboard{
		ID:      "",
		OwnerID: ownerID,
		Name:    storyboardName,
		Users:   make([]*thunderdome.StoryboardUser, 0),
	}

	err := d.DB.QueryRowContext(ctx,
		`SELECT * FROM thunderdome.sb_create($1, $2, $3, $4, $5);`,
		ownerID,
		storyboardName,
		encryptedJoinCode,
		encryptedFacilitatorCode,
		teamID,
	).Scan(&b.ID)
	if err != nil {
		return nil, fmt.Errorf("create storyboard query error: %v", err)
	}

	return b, nil
}

// EditStoryboard updates the storyboard by ID
func (d *Service) EditStoryboard(storyboardID string, storyboardName string, joinCode string, facilitatorCode string) error {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if joinCode != "" {
		encryptedCode, codeErr := db.Encrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return fmt.Errorf("edit storyboard encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = encryptedCode
	}

	if facilitatorCode != "" {
		encryptedCode, codeErr := db.Encrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return fmt.Errorf("edit storyboard encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = encryptedCode
	}

	if _, err := d.DB.Exec(`UPDATE thunderdome.storyboard
        SET name = $2, join_code = $3, facilitator_code = $4, updated_date = NOW()
        WHERE id = $1;`,
		storyboardID, storyboardName, encryptedJoinCode, encryptedFacilitatorCode,
	); err != nil {
		return fmt.Errorf("edit storyboard query error: %v", err)
	}

	return nil
}

// GetStoryboardByID gets a storyboard by ID
func (d *Service) GetStoryboardByID(storyboardID string, userID string) (*thunderdome.Storyboard, error) {
	var cl string
	var joinCode string
	var facilitators string
	var facilitatorCode string
	var b = &thunderdome.Storyboard{
		ID:          storyboardID,
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
		storyboardID,
	).Scan(
		&b.ID,
		&b.Name,
		&b.OwnerID,
		&b.TeamID,
		&cl,
		&joinCode,
		&facilitatorCode,
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
	isFacilitator := db.Contains(b.Facilitators, userID)

	b.Users = d.GetStoryboardUsers(storyboardID)
	b.Goals = d.GetStoryboardGoals(storyboardID)
	b.Personas = d.GetStoryboardPersonas(storyboardID)

	if joinCode != "" {
		decryptedCode, codeErr := db.Decrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("get storyboard encrypt join_code error: %v", codeErr)
		}
		b.JoinCode = decryptedCode
	}

	if facilitatorCode != "" && isFacilitator {
		decryptedCode, codeErr := db.Decrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("get storyboard decrypt facilitator_code error: %v", codeErr)
		}
		b.FacilitatorCode = decryptedCode
	}

	return b, nil
}

// GetStoryboardsByUser gets a list of storyboards by user ID
func (d *Service) GetStoryboardsByUser(userID string, limit int, offset int) ([]*thunderdome.Storyboard, int, error) {
	var count int
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
	`, userID).Scan(
		&count,
	)
	if e != nil {
		return nil, count, fmt.Errorf("get storyboards by user count query error: %v", e)
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
	`, userID, limit, offset)
	if storyboardsErr != nil {
		return nil, count, fmt.Errorf("get storyboards by user query error: %v", storyboardsErr)
	}

	defer storyboardRows.Close()
	for storyboardRows.Next() {
		var b = &thunderdome.Storyboard{
			ID:      "",
			OwnerID: "",
			Name:    "",
			Users:   make([]*thunderdome.StoryboardUser, 0),
		}
		if err := storyboardRows.Scan(
			&b.ID,
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

	return storyboards, count, nil
}

// StoryboardReviseColorLegend revises the storyboard color legend by StoryboardID
func (d *Service) StoryboardReviseColorLegend(storyboardID string, userID string, colorLegend string) (*thunderdome.Storyboard, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard SET updated_date = NOW(), color_legend = $2 WHERE id = $1;`,
		storyboardID,
		colorLegend,
	); err != nil {
		return nil, fmt.Errorf("storyboard revise color legend query error: %v", err)
	}

	storyboard, err := d.GetStoryboardByID(storyboardID, "")
	if err != nil {
		return nil, fmt.Errorf("storyboard revise color legend get storyboards query error: %v", err)
	}

	return storyboard, nil
}

// DeleteStoryboard removes all storyboard associations and the storyboard itself
func (d *Service) DeleteStoryboard(storyboardID string, userID string) error {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard WHERE id = $1;`, storyboardID); err != nil {
		return fmt.Errorf("storyboard delete query error: %v", err)
	}

	return nil
}

// GetStoryboards gets a list of storyboards
func (d *Service) GetStoryboards(limit int, offset int) ([]*thunderdome.Storyboard, int, error) {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	var count int

	e := d.DB.QueryRow(
		"SELECT COUNT(*) FROM thunderdome.storyboard;",
	).Scan(
		&count,
	)
	if e != nil {
		return nil, count, fmt.Errorf("get storyboards count query error: %v", e)
	}

	rows, storyboardErr := d.DB.Query(`
		SELECT s.id, s.name, COALESCE(s.team_id::TEXT, ''), s.created_date, s.updated_date
		FROM thunderdome.storyboard s
		GROUP BY s.id ORDER BY s.created_date DESC
		LIMIT $1 OFFSET $2;
	`, limit, offset)
	if storyboardErr != nil {
		return nil, count, fmt.Errorf("get storyboards query error: %v", storyboardErr)
	}

	defer rows.Close()
	for rows.Next() {
		var sb = &thunderdome.Storyboard{
			Users: make([]*thunderdome.StoryboardUser, 0),
		}
		if err := rows.Scan(
			&sb.ID,
			&sb.Name,
			&sb.TeamID,
			&sb.CreatedDate,
			&sb.UpdatedDate,
		); err != nil {
			d.Logger.Error("get storyboards error", zap.Error(err))
		} else {
			storyboards = append(storyboards, sb)
		}
	}

	return storyboards, count, nil
}

// GetActiveStoryboards gets a list of active storyboards
func (d *Service) GetActiveStoryboards(limit int, offset int) ([]*thunderdome.Storyboard, int, error) {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	var count int

	e := d.DB.QueryRow(
		"SELECT COUNT(DISTINCT su.storyboard_id) FROM thunderdome.storyboard_user su WHERE su.active IS TRUE;",
	).Scan(
		&count,
	)
	if e != nil {
		return nil, count, fmt.Errorf("get active storyboards count query error: %v", e)
	}

	rows, err := d.DB.Query(`
		SELECT s.id, s.name, COALESCE(s.team_id::TEXT, ''), s.created_date, s.updated_date
		FROM thunderdome.storyboard_user su
		LEFT JOIN thunderdome.storyboard s ON s.id = su.storyboard_id
		WHERE su.active IS TRUE GROUP BY s.id
		LIMIT $1 OFFSET $2;
	`, limit, offset)
	if err != nil {
		return nil, count, fmt.Errorf("get active storyboards query error: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var sb = &thunderdome.Storyboard{
			Users: make([]*thunderdome.StoryboardUser, 0),
		}
		if err := rows.Scan(
			&sb.ID,
			&sb.Name,
			&sb.TeamID,
			&sb.CreatedDate,
			&sb.UpdatedDate,
		); err != nil {
			d.Logger.Error("get active storyboards error", zap.Error(err))
		} else {
			storyboards = append(storyboards, sb)
		}
	}

	return storyboards, count, nil
}
