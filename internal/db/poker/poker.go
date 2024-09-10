package poker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/microcosm-cc/bluemonday"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// Service represents a PostgreSQL implementation of thunderdome.PokerDataSvc.
type Service struct {
	DB                  *sql.DB
	Logger              *otelzap.Logger
	AESHashKey          string
	HTMLSanitizerPolicy *bluemonday.Policy
}

// CreateGame creates a new story pointing session
func (d *Service) CreateGame(ctx context.Context, FacilitatorID string, Name string, EstimationScaleID string, PointValuesAllowed []string, Stories []*thunderdome.Story, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, FacilitatorCode string, HideVoterIdentity bool) (*thunderdome.Poker, error) {
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create poker encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("create poker encrypt leader_code error: %v", codeErr)
		}
		encryptedLeaderCode = EncryptedCode
	}

	var b = &thunderdome.Poker{
		Name:                 Name,
		Users:                make([]*thunderdome.PokerUser, 0),
		Stories:              make([]*thunderdome.Story, 0),
		VotingLocked:         true,
		PointValuesAllowed:   PointValuesAllowed,
		AutoFinishVoting:     AutoFinishVoting,
		PointAverageRounding: PointAverageRounding,
		HideVoterIdentity:    HideVoterIdentity,
		Facilitators:         make([]string, 0),
		JoinCode:             JoinCode,
		FacilitatorCode:      FacilitatorCode,
		EstimationScaleID:    EstimationScaleID,
	}
	b.Facilitators = append(b.Facilitators, FacilitatorID)

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		d.Logger.Error("create poker error", zap.Error(err))
	}

	// Insert into poker table
	err = tx.QueryRowContext(ctx, `
            INSERT INTO thunderdome.poker 
            (owner_id, name, estimation_scale_id, point_values_allowed, auto_finish_voting, point_average_rounding,
             hide_voter_identity, join_code, leader_code)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            RETURNING id
        `, FacilitatorID, Name, EstimationScaleID, PointValuesAllowed, AutoFinishVoting,
		PointAverageRounding, HideVoterIdentity, encryptedJoinCode, encryptedLeaderCode,
	).Scan(&b.Id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			d.Logger.Error("update drivers: unable to rollback", zap.Error(rollbackErr))
		}
		d.Logger.Error("create poker error", zap.Error(err))
		return nil, fmt.Errorf("failed to insert into poker table: %v", err)
	}

	// Insert into poker_facilitator table
	_, err = tx.Exec(`
            INSERT INTO thunderdome.poker_facilitator (poker_id, user_id)
            VALUES ($1, $2)
        `, &b.Id, FacilitatorID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			d.Logger.Error("update drivers: unable to rollback", zap.Error(rollbackErr))
		}
		d.Logger.Error("create poker error", zap.Error(err))
		return nil, fmt.Errorf("failed to insert into poker_facilitator table: %v", err)
	}

	// Insert into poker_user table
	_, err = tx.Exec(`
            INSERT INTO thunderdome.poker_user (poker_id, user_id)
            VALUES ($1, $2)
        `, &b.Id, FacilitatorID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			d.Logger.Error("update drivers: unable to rollback", zap.Error(rollbackErr))
		}
		d.Logger.Error("create poker error", zap.Error(err))
		return nil, fmt.Errorf("failed to insert into poker_user table: %v", err)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		d.Logger.Error("update drivers: unable to commit", zap.Error(commitErr))
		return nil, fmt.Errorf("failed to create poker game: %v", commitErr)
	}

	for _, plan := range Stories {
		plan.Votes = make([]*thunderdome.Vote, 0)

		e := d.DB.QueryRowContext(ctx,
			`INSERT INTO thunderdome.poker_story (poker_id, name, type, reference_id, link, description, acceptance_criteria, position) 
					VALUES ($1, $2, $3, $4, $5, $6, $7, (
					  coalesce(
						(select max(position) from thunderdome.poker_story where poker_id = $1),
						-1
					  ) + 1
					)) RETURNING id`,
			b.Id,
			plan.Name,
			plan.Type,
			plan.ReferenceId,
			plan.Link,
			plan.Description,
			plan.AcceptanceCriteria,
		).Scan(&plan.Id)
		if e != nil {
			d.Logger.Error("insert stories error", zap.Error(e))
		}
	}

	b.Stories = Stories

	return b, nil
}

// TeamCreateGame creates a new story pointing session associated to a team
func (d *Service) TeamCreateGame(ctx context.Context, TeamID string, FacilitatorID string, Name string, EstimationScaleID string, PointValuesAllowed []string, Stories []*thunderdome.Story, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, FacilitatorCode string, HideVoterIdentity bool) (*thunderdome.Poker, error) {
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("team create poker encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("team create poker encrypt leader_code error: %v", codeErr)
		}
		encryptedLeaderCode = EncryptedCode
	}

	var b = &thunderdome.Poker{
		Name:                 Name,
		Users:                make([]*thunderdome.PokerUser, 0),
		Stories:              make([]*thunderdome.Story, 0),
		VotingLocked:         true,
		PointValuesAllowed:   PointValuesAllowed,
		AutoFinishVoting:     AutoFinishVoting,
		PointAverageRounding: PointAverageRounding,
		HideVoterIdentity:    HideVoterIdentity,
		Facilitators:         make([]string, 0),
		JoinCode:             JoinCode,
		FacilitatorCode:      FacilitatorCode,
		EstimationScaleID:    EstimationScaleID,
		TeamID:               TeamID,
	}
	b.Facilitators = append(b.Facilitators, FacilitatorID)

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		d.Logger.Error("create poker error", zap.Error(err))
	}

	// Insert into poker table
	err = tx.QueryRowContext(ctx, `
            INSERT INTO thunderdome.poker 
            (owner_id, name, estimation_scale_id, point_values_allowed, auto_finish_voting, point_average_rounding,
             hide_voter_identity, join_code, leader_code, team_id)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            RETURNING id
        `, FacilitatorID, Name, EstimationScaleID, PointValuesAllowed, AutoFinishVoting,
		PointAverageRounding, HideVoterIdentity, encryptedJoinCode, encryptedLeaderCode, TeamID,
	).Scan(&b.Id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			d.Logger.Error("update drivers: unable to rollback", zap.Error(rollbackErr))
		}
		return nil, fmt.Errorf("failed to insert into poker table: %v", err)
	}

	// Insert into poker_facilitator table
	_, err = tx.Exec(`
            INSERT INTO thunderdome.poker_facilitator (poker_id, user_id)
            VALUES ($1, $2)
        `, &b.Id, FacilitatorID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			d.Logger.Error("update drivers: unable to rollback", zap.Error(rollbackErr))
		}
		return nil, fmt.Errorf("failed to insert into poker_facilitator table: %v", err)
	}

	// Insert into poker_user table
	_, err = tx.Exec(`
            INSERT INTO thunderdome.poker_user (poker_id, user_id)
            VALUES ($1, $2)
        `, &b.Id, FacilitatorID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			d.Logger.Error("update drivers: unable to rollback", zap.Error(rollbackErr))
		}
		return nil, fmt.Errorf("failed to insert into poker_user table: %v", err)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		d.Logger.Error("update drivers: unable to commit", zap.Error(commitErr))
		return nil, fmt.Errorf("failed to create poker game: %v", commitErr)
	}

	for _, plan := range Stories {
		plan.Votes = make([]*thunderdome.Vote, 0)

		e := d.DB.QueryRowContext(ctx,
			`INSERT INTO thunderdome.poker_story (poker_id, name, type, reference_id, link, description, acceptance_criteria, position) 
					VALUES ($1, $2, $3, $4, $5, $6, $7, (
					  coalesce(
						(select max(position) from thunderdome.poker_story where poker_id = $1),
						-1
					  ) + 1
					)) RETURNING id`,
			b.Id,
			plan.Name,
			plan.Type,
			plan.ReferenceId,
			plan.Link,
			plan.Description,
			plan.AcceptanceCriteria,
		).Scan(&plan.Id)
		if e != nil {
			d.Logger.Error("insert stories error", zap.Error(e))
		}
	}

	b.Stories = Stories

	return b, nil
}

// UpdateGame updates the game by ID
func (d *Service) UpdateGame(PokerID string, Name string, PointValuesAllowed []string, AutoFinishVoting bool, PointAverageRounding string, HideVoterIdentity bool, JoinCode string, FacilitatorCode string, TeamID string) error {
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return fmt.Errorf("update poker encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return fmt.Errorf("update poker encrypt leader_code error: %v", codeErr)
		}
		encryptedLeaderCode = EncryptedCode
	}

	if _, err := d.DB.Exec(`
		UPDATE thunderdome.poker
		SET name = $2, point_values_allowed = $3, auto_finish_voting = $4, point_average_rounding = $5,
		 hide_voter_identity = $6, join_code = $7, leader_code = $8, updated_date = NOW(), team_id = NULLIF($9, '')::uuid
		WHERE id = $1`,
		PokerID, Name, PointValuesAllowed, AutoFinishVoting, PointAverageRounding,
		HideVoterIdentity, encryptedJoinCode, encryptedLeaderCode, TeamID,
	); err != nil {
		return fmt.Errorf("update poker query error: %v", err)
	}

	return nil
}

// GetGame gets a game by ID
func (d *Service) GetGame(PokerID string, UserID string) (*thunderdome.Poker, error) {
	var b = &thunderdome.Poker{
		Id:           PokerID,
		Users:        make([]*thunderdome.PokerUser, 0),
		Stories:      make([]*thunderdome.Story, 0),
		VotingLocked: true,
		Facilitators: make([]string, 0),
	}

	// get game
	var facilitators string
	var JoinCode string
	var FacilitatorCode string
	var estimationScaleJSON []byte
	var vArray pgtype.Array[string]
	m := pgtype.NewMap()
	e := d.DB.QueryRow(
		`
		SELECT b.id, b.name, b.voting_locked, COALESCE(b.active_story_id::text, ''), b.auto_finish_voting, 
		b.point_average_rounding, b.hide_voter_identity, COALESCE(b.join_code, ''), COALESCE(b.leader_code, ''),
		b.estimation_scale_id, b.point_values_allowed, COALESCE(b.team_id::text, ''), b.created_date, b.updated_date,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS leaders,
		COALESCE(
			json_build_object(
				'id', es.id,
				'name', es.name,
				'description', es.description,
				'scale_type', es.scale_type,
				'values', es.values,
				'created_by', es.created_by,
				'created_at', es.created_at,
				'updated_at', es.updated_at,
				'is_public', es.is_public,
				'organization_id', es.organization_id,
				'team_id', es.team_id,
				'default_scale', es.default_scale
			)::jsonb,
			'{}'::jsonb
		) AS estimation_scale
		FROM thunderdome.poker b
		LEFT JOIN thunderdome.poker_facilitator bl ON b.id = bl.poker_id
		LEFT JOIN thunderdome.estimation_scale es ON b.estimation_scale_id = es.id
		WHERE b.id = $1
		GROUP BY b.id, es.id`,
		PokerID,
	).Scan(
		&b.Id,
		&b.Name,
		&b.VotingLocked,
		&b.ActiveStoryID,
		&b.AutoFinishVoting,
		&b.PointAverageRounding,
		&b.HideVoterIdentity,
		&JoinCode,
		&FacilitatorCode,
		&b.EstimationScaleID,
		m.SQLScanner(&vArray),
		&b.TeamID,
		&b.CreatedDate,
		&b.UpdatedDate,
		&facilitators,
		&estimationScaleJSON,
	)
	if e != nil {
		return nil, fmt.Errorf("get poker query error: %v", e)
	}

	b.PointValuesAllowed = vArray.Elements

	_ = json.Unmarshal([]byte(facilitators), &b.Facilitators)

	// Unmarshal the estimation scale JSON into the EstimationScale field
	if len(estimationScaleJSON) > 0 {
		var estimationScale thunderdome.EstimationScale
		if err := json.Unmarshal(estimationScaleJSON, &estimationScale); err != nil {
			return nil, fmt.Errorf("error unmarshaling estimation scale: %v", err)
		}
		b.EstimationScale = &estimationScale
	}

	isFacilitator := db.Contains(b.Facilitators, UserID)

	if JoinCode != "" {
		DecryptedCode, codeErr := db.Decrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("get poker decode join_code error: %v", codeErr)
		}
		b.JoinCode = DecryptedCode
	}

	if FacilitatorCode != "" && isFacilitator {
		DecryptedCode, codeErr := db.Decrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("get poker decode leader_code error: %v", codeErr)
		}
		b.FacilitatorCode = DecryptedCode
	}

	b.Users = d.GetUsers(PokerID)
	b.Stories = d.GetStories(PokerID, UserID)

	return b, nil
}

// GetGamesByUser gets a list of games by UserID
func (d *Service) GetGamesByUser(UserID string, Limit int, Offset int) ([]*thunderdome.Poker, int, error) {
	var Count int
	var games = make([]*thunderdome.Poker, 0)

	e := d.DB.QueryRow(`
		WITH user_teams AS (
			SELECT t.id FROM thunderdome.team_user tu
			LEFT JOIN thunderdome.team t ON t.id = tu.team_id
			WHERE tu.user_id = $1
		),
		team_games AS (
			SELECT id FROM thunderdome.poker WHERE team_id IN (SELECT id FROM user_teams)
		),
		user_games AS (
			SELECT u.poker_id AS id FROM thunderdome.poker_user u
			WHERE u.user_id = $1 AND u.abandoned = false
		),
		games AS (
			SELECT id from user_games UNION SELECT id FROM team_games
		)
		SELECT COUNT(*) FROM games;
	`, UserID).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, fmt.Errorf("get poker by user count query error: %v", e)
	}

	gameRows, gamesErr := d.DB.Query(`
		WITH user_teams AS (
			SELECT t.id, t.name FROM thunderdome.team_user tu
			LEFT JOIN thunderdome.team t ON t.id = tu.team_id
			WHERE tu.user_id = $1
		),
		team_games AS (
			SELECT id FROM thunderdome.poker WHERE team_id IN (SELECT id FROM user_teams)
		),
		user_games AS (
			SELECT u.poker_id AS id FROM thunderdome.poker_user u
			WHERE u.user_id = $1 AND u.abandoned = false
		),
		games AS (
			SELECT id from user_games UNION SELECT id FROM team_games
		),
		stories AS (
			SELECT poker_id, points FROM thunderdome.poker_story WHERE poker_id IN (SELECT poker_id FROM games)
		),
		facilitators AS (
			SELECT poker_id, user_id FROM thunderdome.poker_facilitator WHERE poker_id IN (SELECT poker_id FROM games)
		)
		SELECT p.id, p.name, p.voting_locked, COALESCE(p.active_story_id::text, ''), p.point_values_allowed, p.auto_finish_voting,
		  p.point_average_rounding, p.created_date, p.updated_date,
		  CASE WHEN COUNT(s) = 0 THEN '[]'::json ELSE array_to_json(array_agg(row_to_json(s))) END AS stories,
		  CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS facilitators,
		  min(COALESCE(t.name, '')) as team_name, p.estimation_scale_id,
		  COALESCE(
			json_build_object(
				'id', es.id,
				'name', es.name,
				'description', es.description,
				'scale_type', es.scale_type,
				'values', es.values,
				'created_by', es.created_by,
				'created_at', es.created_at,
				'updated_at', es.updated_at,
				'is_public', es.is_public,
				'organization_id', es.organization_id,
				'team_id', es.team_id,
				'default_scale', es.default_scale
			)::jsonb,
			'{}'::jsonb
		) AS estimation_scale
		FROM thunderdome.poker p
		LEFT JOIN stories AS s ON s.poker_id = p.id
		LEFT JOIN facilitators AS bl ON bl.poker_id = p.id
		LEFT JOIN user_teams t ON t.id = p.team_id
		LEFT JOIN thunderdome.estimation_scale es ON p.estimation_scale_id = es.id
		WHERE p.id IN (SELECT id FROM games)
		GROUP BY p.id, p.created_date, es.id
		ORDER BY p.created_date DESC
		LIMIT $2 OFFSET $3
	`, UserID, Limit, Offset)
	if gamesErr != nil {
		d.Logger.Error("get poker by user query error", zap.Error(gamesErr))
		return nil, Count, fmt.Errorf("get poker by user query error: %v", gamesErr)
	}

	defer gameRows.Close()
	for gameRows.Next() {
		var stories string
		var estimationScale string
		var vArray pgtype.Array[string]
		m := pgtype.NewMap()
		var facilitators string
		var b = &thunderdome.Poker{
			Users:              make([]*thunderdome.PokerUser, 0),
			Stories:            make([]*thunderdome.Story, 0),
			VotingLocked:       true,
			PointValuesAllowed: make([]string, 0),
			AutoFinishVoting:   true,
			Facilitators:       make([]string, 0),
		}
		if err := gameRows.Scan(
			&b.Id,
			&b.Name,
			&b.VotingLocked,
			&b.ActiveStoryID,
			m.SQLScanner(&vArray),
			&b.AutoFinishVoting,
			&b.PointAverageRounding,
			&b.CreatedDate,
			&b.UpdatedDate,
			&stories,
			&facilitators,
			&b.TeamName,
			&b.EstimationScaleID,
			&estimationScale,
		); err != nil {
			d.Logger.Error("error getting poker by user", zap.Error(err))
		} else {
			_ = json.Unmarshal([]byte(stories), &b.Stories)
			_ = json.Unmarshal([]byte(facilitators), &b.Facilitators)
			_ = json.Unmarshal([]byte(estimationScale), &b.EstimationScale)
			b.PointValuesAllowed = vArray.Elements

			games = append(games, b)
		}
	}

	return games, Count, nil
}

// DeleteGame removes all game associations and the game itself by PokerID
func (d *Service) DeleteGame(PokerID string) error {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.poker WHERE id = $1;`, PokerID); err != nil {
		return fmt.Errorf("poker delete query error: %v", err)
	}

	return nil
}

// GetGames gets a list of games
func (d *Service) GetGames(Limit int, Offset int) ([]*thunderdome.Poker, int, error) {
	var games = make([]*thunderdome.Poker, 0)
	var Count int

	e := d.DB.QueryRow(
		"SELECT COUNT(*) FROM thunderdome.poker;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, fmt.Errorf("get poker games count query error: %v", e)
	}

	rows, gamesErr := d.DB.Query(`
		SELECT b.id, b.name, b.voting_locked, b.active_story_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, b.created_date, b.updated_date,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS leaders
		FROM thunderdome.poker b
		LEFT JOIN thunderdome.poker_facilitator bl ON b.id = bl.poker_id
		GROUP BY b.id ORDER BY b.created_date DESC
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if gamesErr != nil {
		return nil, Count, fmt.Errorf("get poker games query error: %v", gamesErr)
	}

	defer rows.Close()
	for rows.Next() {
		var vArray pgtype.Array[string]
		m := pgtype.NewMap()
		var facilitators string
		var ActiveStoryID sql.NullString
		var b = &thunderdome.Poker{
			Users:              make([]*thunderdome.PokerUser, 0),
			Stories:            make([]*thunderdome.Story, 0),
			VotingLocked:       true,
			PointValuesAllowed: make([]string, 0),
			AutoFinishVoting:   true,
			Facilitators:       make([]string, 0),
		}
		if err := rows.Scan(
			&b.Id,
			&b.Name,
			&b.VotingLocked,
			&ActiveStoryID,
			m.SQLScanner(&vArray),
			&b.AutoFinishVoting,
			&b.PointAverageRounding,
			&b.CreatedDate,
			&b.UpdatedDate,
			&facilitators,
		); err != nil {
			d.Logger.Error("get poker games query error", zap.Error(err))
		} else {
			b.PointValuesAllowed = vArray.Elements
			_ = json.Unmarshal([]byte(facilitators), &b.Facilitators)
			b.ActiveStoryID = ActiveStoryID.String
			games = append(games, b)
		}
	}

	return games, Count, nil
}

// GetActiveGames gets a list of active games
func (d *Service) GetActiveGames(Limit int, Offset int) ([]*thunderdome.Poker, int, error) {
	var games = make([]*thunderdome.Poker, 0)
	var Count int

	e := d.DB.QueryRow(
		"SELECT COUNT(DISTINCT pu.poker_id) FROM thunderdome.poker_user pu WHERE pu.active IS TRUE;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, fmt.Errorf("get active poker games count query error: %v", e)
	}

	rows, gamesErr := d.DB.Query(`
		SELECT b.id, b.name, b.voting_locked, b.active_story_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, b.created_date, b.updated_date,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS leaders
		FROM thunderdome.poker_user bu
		LEFT JOIN thunderdome.poker b ON b.id = bu.poker_id
		LEFT JOIN thunderdome.poker_facilitator bl ON b.id = bl.poker_id
		WHERE bu.active IS TRUE GROUP BY b.id
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if gamesErr != nil {
		return nil, Count, fmt.Errorf("get active poker games query error: %v", gamesErr)
	}

	defer rows.Close()
	for rows.Next() {
		var vArray pgtype.Array[string]
		m := pgtype.NewMap()
		var facilitators string
		var ActiveStoryID sql.NullString
		var b = &thunderdome.Poker{
			Users:              make([]*thunderdome.PokerUser, 0),
			Stories:            make([]*thunderdome.Story, 0),
			VotingLocked:       true,
			PointValuesAllowed: make([]string, 0),
			AutoFinishVoting:   true,
			Facilitators:       make([]string, 0),
		}
		if err := rows.Scan(
			&b.Id,
			&b.Name,
			&b.VotingLocked,
			&ActiveStoryID,
			m.SQLScanner(&vArray),
			&b.AutoFinishVoting,
			&b.PointAverageRounding,
			&b.CreatedDate,
			&b.UpdatedDate,
			&facilitators,
		); err != nil {
			d.Logger.Error("get active poker games query error", zap.Error(err))
		} else {
			b.PointValuesAllowed = vArray.Elements
			_ = json.Unmarshal([]byte(facilitators), &b.Facilitators)
			b.ActiveStoryID = ActiveStoryID.String
			games = append(games, b)
		}
	}

	return games, Count, nil
}

// PurgeOldGames deletes games older than {DaysOld} days
func (d *Service) PurgeOldGames(ctx context.Context, DaysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.poker WHERE last_active < (NOW() - $1 * interval '1 day');`,
		DaysOld,
	); err != nil {
		return fmt.Errorf("clean poker games query error: %v", err)
	}

	return nil
}
