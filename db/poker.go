package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/microcosm-cc/bluemonday"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// PokerService represents a PostgreSQL implementation of thunderdome.PokerDataSvc.
type PokerService struct {
	DB                  *sql.DB
	Logger              *otelzap.Logger
	AESHashKey          string
	HTMLSanitizerPolicy *bluemonday.Policy
}

// CreateGame creates a new story pointing session
func (d *PokerService) CreateGame(ctx context.Context, FacilitatorID string, Name string, PointValuesAllowed []string, Stories []*thunderdome.Story, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, FacilitatorCode string, HideVoterIdentity bool) (*thunderdome.Poker, error) {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to create poker join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to create poker leader_code")
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
	}
	b.Facilitators = append(b.Facilitators, FacilitatorID)

	e := d.DB.QueryRowContext(ctx,
		`SELECT pokerid FROM thunderdome.poker_create($1, $2, $3, $4, $5, $6, $7, $8);`,
		FacilitatorID,
		Name,
		string(pointValuesJSON),
		AutoFinishVoting,
		PointAverageRounding,
		HideVoterIdentity,
		encryptedJoinCode,
		encryptedLeaderCode,
	).Scan(&b.Id)
	if e != nil {
		d.Logger.Error("poker_create query error", zap.Error(e))
		return nil, errors.New("error creating poker")
	}

	for _, plan := range Stories {
		plan.Votes = make([]*thunderdome.Vote, 0)

		e := d.DB.QueryRowContext(ctx,
			`INSERT INTO thunderdome.poker_story (poker_id, name, type, reference_id, link, description, acceptance_criteria) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
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
func (d *PokerService) TeamCreateGame(ctx context.Context, TeamID string, FacilitatorID string, Name string, PointValuesAllowed []string, Stories []*thunderdome.Story, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, FacilitatorCode string, HideVoterIdentity bool) (*thunderdome.Poker, error) {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to create poker join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to create poker leader_code")
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
	}
	b.Facilitators = append(b.Facilitators, FacilitatorID)

	e := d.DB.QueryRowContext(ctx,
		`SELECT pokerid FROM thunderdome.team_create_poker($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		TeamID,
		FacilitatorID,
		Name,
		string(pointValuesJSON),
		AutoFinishVoting,
		PointAverageRounding,
		HideVoterIdentity,
		encryptedJoinCode,
		encryptedLeaderCode,
	).Scan(&b.Id)
	if e != nil {
		d.Logger.Error("team_create_poker query error", zap.Error(e))
		return nil, errors.New("error creating poker")
	}

	for _, plan := range Stories {
		plan.Votes = make([]*thunderdome.Vote, 0)

		e := d.DB.QueryRowContext(ctx,
			`INSERT INTO thunderdome.poker_story (poker_id, name, type, reference_id, link, description, acceptance_criteria) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
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
func (d *PokerService) UpdateGame(PokerID string, Name string, PointValuesAllowed []string, AutoFinishVoting bool, PointAverageRounding string, HideVoterIdentity bool, JoinCode string, FacilitatorCode string) error {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return errors.New("unable to revise poker join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if FacilitatorCode != "" {
		EncryptedCode, codeErr := encrypt(FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return errors.New("unable to revise poker leadercode")
		}
		encryptedLeaderCode = EncryptedCode
	}

	if _, err := d.DB.Exec(`
		UPDATE thunderdome.poker
		SET name = $2, point_values_allowed = $3, auto_finish_voting = $4, point_average_rounding = $5, hide_voter_identity = $6, join_code = $7, leader_code = $8, updated_date = NOW()
		WHERE id = $1`,
		PokerID, Name, string(pointValuesJSON), AutoFinishVoting, PointAverageRounding, HideVoterIdentity, encryptedJoinCode, encryptedLeaderCode,
	); err != nil {
		d.Logger.Error("update poker error", zap.Error(err))
		return errors.New("unable to revise poker")
	}

	return nil
}

// GetFacilitatorCode retrieve the game leader_code
func (d *PokerService) GetFacilitatorCode(PokerID string) (string, error) {
	var EncryptedLeaderCode string

	if err := d.DB.QueryRow(`
		SELECT COALESCE(leader_code, '') FROM thunderdome.poker
		WHERE id = $1`,
		PokerID,
	).Scan(&EncryptedLeaderCode); err != nil {
		d.Logger.Error("get poker leadercode error", zap.Error(err))
		return "", errors.New("unable to retrieve poker leader_code")
	}

	if EncryptedLeaderCode == "" {
		return "", errors.New("unable to retrieve poker leader_code")
	}
	DecryptedCode, codeErr := decrypt(EncryptedLeaderCode, d.AESHashKey)
	if codeErr != nil {
		return "", errors.New("unable to retrieve poker leader_code")
	}

	return DecryptedCode, nil
}

// GetGame gets a game by ID
func (d *PokerService) GetGame(PokerID string, UserID string) (*thunderdome.Poker, error) {
	var b = &thunderdome.Poker{
		Id:                 PokerID,
		Users:              make([]*thunderdome.PokerUser, 0),
		Stories:            make([]*thunderdome.Story, 0),
		VotingLocked:       true,
		PointValuesAllowed: make([]string, 0),
		AutoFinishVoting:   true,
		Facilitators:       make([]string, 0),
	}

	// get battle
	var ActivePlanID sql.NullString
	var pv string
	var leaders string
	var JoinCode string
	var LeaderCode string
	e := d.DB.QueryRow(
		`
		SELECT b.id, b.name, b.voting_locked, b.active_story_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, b.hide_voter_identity, COALESCE(b.join_code, ''), COALESCE(b.leader_code, ''), b.created_date, b.updated_date,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS leaders
		FROM thunderdome.poker b
		LEFT JOIN thunderdome.poker_facilitator bl ON b.id = bl.poker_id
		WHERE b.id = $1
		GROUP BY b.id`,
		PokerID,
	).Scan(
		&b.Id,
		&b.Name,
		&b.VotingLocked,
		&ActivePlanID,
		&pv,
		&b.AutoFinishVoting,
		&b.PointAverageRounding,
		&b.HideVoterIdentity,
		&JoinCode,
		&LeaderCode,
		&b.CreatedDate,
		&b.UpdatedDate,
		&leaders,
	)
	if e != nil {
		d.Logger.Error("error getting poker", zap.Error(e))
		return nil, errors.New("not found")
	}

	_ = json.Unmarshal([]byte(leaders), &b.Facilitators)
	_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
	b.ActiveStoryID = ActivePlanID.String

	isBattleLeader := contains(b.Facilitators, UserID)

	if JoinCode != "" {
		DecryptedCode, codeErr := decrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to decode join_code")
		}
		b.JoinCode = DecryptedCode
	}

	if LeaderCode != "" && isBattleLeader {
		DecryptedCode, codeErr := decrypt(LeaderCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to decode leader_code")
		}
		b.FacilitatorCode = DecryptedCode
	}

	b.Users = d.GetUsers(PokerID)
	b.Stories = d.GetStories(PokerID, UserID)

	return b, nil
}

// GetGamesByUser gets a list of games by UserID
func (d *PokerService) GetGamesByUser(UserID string, Limit int, Offset int) ([]*thunderdome.Poker, int, error) {
	var Count int
	var battles = make([]*thunderdome.Poker, 0)

	e := d.DB.QueryRow(`
		SELECT COUNT(*) FROM thunderdome.poker b
		LEFT JOIN thunderdome.poker_user bw ON b.id = bw.poker_id
		WHERE bw.user_id = $1 AND bw.abandoned = false;
	`, UserID).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	battleRows, battlesErr := d.DB.Query(`
		SELECT b.id, b.name, b.voting_locked, b.active_story_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, b.created_date, b.updated_date,
		CASE WHEN COUNT(p) = 0 THEN '[]'::json ELSE array_to_json(array_agg(row_to_json(p))) END AS stories,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS facilitators
		FROM thunderdome.poker b
		LEFT JOIN thunderdome.poker_story p ON b.id = p.poker_id
		LEFT JOIN thunderdome.poker_facilitator bl ON b.id = bl.poker_id
		LEFT JOIN thunderdome.poker_user bw ON b.id = bw.poker_id
		WHERE bw.user_id = $1 AND bw.abandoned = false
		GROUP BY b.id ORDER BY b.created_date DESC
		LIMIT $2 OFFSET $3
	`, UserID, Limit, Offset)
	if battlesErr != nil {
		return nil, Count, errors.New("not found")
	}

	defer battleRows.Close()
	for battleRows.Next() {
		var plans string
		var pv string
		var leaders string
		var ActivePlanID sql.NullString
		var b = &thunderdome.Poker{
			Users:              make([]*thunderdome.PokerUser, 0),
			Stories:            make([]*thunderdome.Story, 0),
			VotingLocked:       true,
			PointValuesAllowed: make([]string, 0),
			AutoFinishVoting:   true,
			Facilitators:       make([]string, 0),
		}
		if err := battleRows.Scan(
			&b.Id,
			&b.Name,
			&b.VotingLocked,
			&ActivePlanID,
			&pv,
			&b.AutoFinishVoting,
			&b.PointAverageRounding,
			&b.CreatedDate,
			&b.UpdatedDate,
			&plans,
			&leaders,
		); err != nil {
			d.Logger.Error("error getting poker by user", zap.Error(e))
		} else {
			_ = json.Unmarshal([]byte(plans), &b.Stories)
			_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
			_ = json.Unmarshal([]byte(leaders), &b.Facilitators)
			b.ActiveStoryID = ActivePlanID.String
			battles = append(battles, b)
		}
	}

	return battles, Count, nil
}

// ConfirmFacilitator confirms the user is a facilitator of the game
func (d *PokerService) ConfirmFacilitator(BattleID string, UserID string) error {
	var leaderID string
	var role string
	err := d.DB.QueryRow("SELECT type FROM thunderdome.users WHERE id = $1", UserID).Scan(&role)
	if err != nil {
		d.Logger.Error("error getting user role", zap.Error(err))
		return errors.New("unable to get user role")
	}

	e := d.DB.QueryRow("SELECT user_id FROM thunderdome.poker_facilitator WHERE poker_id = $1 AND user_id = $2", BattleID, UserID).Scan(&leaderID)
	if e != nil && role != "ADMIN" {
		d.Logger.Error("error confirming poker facilitator", zap.Error(e))
		return errors.New("not a poker facilitator")
	}

	return nil
}

// GetUserActiveStatus checks game active status of User
func (d *PokerService) GetUserActiveStatus(PokerID string, UserID string) error {
	var active bool

	e := d.DB.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM thunderdome.poker_user
		WHERE user_id = $2 AND poker_id = $1;`,
		PokerID,
		UserID,
	).Scan(
		&active,
	)
	if e != nil {
		return e
	}

	if active {
		return errors.New("DUPLICATE_BATTLE_USER")
	}

	return nil
}

// GetUsers retrieves the users for a given game
func (d *PokerService) GetUsers(PokerID string) []*thunderdome.PokerUser {
	var users = make([]*thunderdome.PokerUser, 0)
	rows, err := d.DB.Query(
		`SELECT
			u.id, u.name, u.type, u.avatar, pu.active, pu.spectator, COALESCE(u.email, '')
		FROM thunderdome.poker_user pu
		LEFT JOIN thunderdome.users u ON pu.user_id = u.id
		WHERE pu.poker_id = $1
		ORDER BY u.name`,
		PokerID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w thunderdome.PokerUser
			if err := rows.Scan(&w.Id, &w.Name, &w.Type, &w.Avatar, &w.Active, &w.Spectator, &w.GravatarHash); err != nil {
				d.Logger.Error("error getting poker users", zap.Error(err))
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

// GetActiveUsers retrieves the active users for a given game
func (d *PokerService) GetActiveUsers(PokerID string) []*thunderdome.PokerUser {
	var users = make([]*thunderdome.PokerUser, 0)
	rows, err := d.DB.Query(
		`SELECT
			w.id, w.name, w.type, w.avatar, bw.active, bw.spectator, COALESCE(w.email, '')
		FROM thunderdome.poker_user bw
		LEFT JOIN thunderdome.users w ON bw.user_id = w.id
		WHERE bw.poker_id = $1 AND bw.active = true
		ORDER BY w.name`,
		PokerID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w thunderdome.PokerUser
			if err := rows.Scan(&w.Id, &w.Name, &w.Type, &w.Avatar, &w.Active, &w.Spectator, &w.GravatarHash); err != nil {
				d.Logger.Error("error getting active poker users", zap.Error(err))
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

// AddUser adds a user by ID to the game by ID
func (d *PokerService) AddUser(PokerID string, UserID string) ([]*thunderdome.PokerUser, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.poker_user (poker_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (poker_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		PokerID,
		UserID,
	); err != nil {
		d.Logger.Error("error adding user to poker", zap.Error(err))
	}

	users := d.GetUsers(PokerID)

	return users, nil
}

// RetreatUser removes a user from the current game by ID
func (d *PokerService) RetreatUser(PokerID string, UserID string) []*thunderdome.PokerUser {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.poker_user SET active = false WHERE poker_id = $1 AND user_id = $2`, PokerID, UserID); err != nil {
		d.Logger.Error("error updating poker user to active false", zap.Error(err))
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("error updating user last active timestamp", zap.Error(err))
	}

	users := d.GetUsers(PokerID)

	return users
}

// AbandonGame removes a user from the current game by ID and sets abandoned true
func (d *PokerService) AbandonGame(PokerID string, UserID string) ([]*thunderdome.PokerUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.poker_user SET active = false, abandoned = true WHERE poker_id = $1 AND user_id = $2`, PokerID, UserID); err != nil {
		d.Logger.Error("error updating battle user to abandoned", zap.Error(err))
		return nil, err
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("error updating user last active timestamp", zap.Error(err))
		return nil, err
	}

	users := d.GetUsers(PokerID)

	return users, nil
}

// AddFacilitator makes a user a facilitator of the game
func (d *PokerService) AddFacilitator(PokerID string, UserID string) ([]string, error) {
	leaders := make([]string, 0)

	// set battle leader
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.poker_facilitator (poker_id, user_id) VALUES ($1, $2);`,
		PokerID, UserID); err != nil {
		d.Logger.Error("set poker facilitator query error", zap.Error(err))
		return nil, errors.New("unable to make facilitator")
	}

	leaderRows, leadersErr := d.DB.Query(`
		SELECT user_id FROM thunderdome.poker_facilitator WHERE poker_id = $1;
	`, PokerID)
	if leadersErr != nil {
		return leaders, nil
	}

	defer leaderRows.Close()
	for leaderRows.Next() {
		var leader string
		if err := leaderRows.Scan(
			&leader,
		); err != nil {
			d.Logger.Error("poker_facilitator query scan error", zap.Error(err))
		} else {
			leaders = append(leaders, leader)
		}
	}

	return leaders, nil
}

// RemoveFacilitator removes a user from game facilitators
func (d *PokerService) RemoveFacilitator(PokerID string, UserID string) ([]string, error) {
	leaders := make([]string, 0)

	// set battle leader
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.poker_facilitator WHERE poker_id = $1 AND user_id = $2;`,
		PokerID, UserID); err != nil {
		d.Logger.Error("delete poker_facilitator query error", zap.Error(err))
		return nil, errors.New("unable to delete facilitator")
	}

	leaderRows, leadersErr := d.DB.Query(`
		SELECT user_id FROM thunderdome.poker_facilitator WHERE poker_id = $1;
	`, PokerID)
	if leadersErr != nil {
		return leaders, nil
	}

	defer leaderRows.Close()
	for leaderRows.Next() {
		var leader string
		if err := leaderRows.Scan(
			&leader,
		); err != nil {
			d.Logger.Error("poker_facilitator query scan error", zap.Error(err))
		} else {
			leaders = append(leaders, leader)
		}
	}

	return leaders, nil
}

// ToggleSpectator changes a game users spectator status
func (d *PokerService) ToggleSpectator(PokerID string, UserID string, Spectator bool) ([]*thunderdome.PokerUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.poker_user SET spectator = $3 WHERE poker_id = $1 AND user_id = $2`, PokerID, UserID, Spectator); err != nil {
		d.Logger.Error("update poker user spectator error", zap.Error(err))
		return nil, err
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("error updating user last active timestamp", zap.Error(err))
	}

	users := d.GetUsers(PokerID)

	return users, nil
}

// DeleteGame removes all game associations and the game itself by PokerID
func (d *PokerService) DeleteGame(PokerID string) error {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.poker WHERE id = $1;`, PokerID); err != nil {
		d.Logger.Error("delete poker error", zap.Error(err))
		return err
	}

	return nil
}

// AddFacilitatorsByEmail adds additional game facilitators by email
func (d *PokerService) AddFacilitatorsByEmail(ctx context.Context, PokerID string, FacilitatorEmails []string) ([]string, error) {
	var leaders string
	var newLeaders []string

	for i, email := range FacilitatorEmails {
		FacilitatorEmails[i] = sanitizeEmail(email)
	}
	emails := strings.Join(FacilitatorEmails[:], ",")

	e := d.DB.QueryRowContext(ctx,
		`SELECT facilitators FROM thunderdome.poker_facilitator_add_by_email($1, $2);`,
		PokerID, emails,
	).Scan(&leaders)
	if e != nil {
		d.Logger.Error("poker_facilitator_add_by_email query error", zap.Error(e))
		return nil, errors.New("error adding poker facilitator by email")
	}

	_ = json.Unmarshal([]byte(leaders), &newLeaders)

	return newLeaders, nil
}

// GetGames gets a list of games
func (d *PokerService) GetGames(Limit int, Offset int) ([]*thunderdome.Poker, int, error) {
	var battles = make([]*thunderdome.Poker, 0)
	var Count int

	e := d.DB.QueryRow(
		"SELECT COUNT(*) FROM thunderdome.poker;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	battleRows, battlesErr := d.DB.Query(`
		SELECT b.id, b.name, b.voting_locked, b.active_story_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, b.created_date, b.updated_date,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS leaders
		FROM thunderdome.poker b
		LEFT JOIN thunderdome.poker_facilitator bl ON b.id = bl.poker_id
		GROUP BY b.id ORDER BY b.created_date DESC
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if battlesErr != nil {
		return nil, Count, battlesErr
	}

	defer battleRows.Close()
	for battleRows.Next() {
		var pv string
		var leaders string
		var ActivePlanID sql.NullString
		var b = &thunderdome.Poker{
			Users:              make([]*thunderdome.PokerUser, 0),
			Stories:            make([]*thunderdome.Story, 0),
			VotingLocked:       true,
			PointValuesAllowed: make([]string, 0),
			AutoFinishVoting:   true,
			Facilitators:       make([]string, 0),
		}
		if err := battleRows.Scan(
			&b.Id,
			&b.Name,
			&b.VotingLocked,
			&ActivePlanID,
			&pv,
			&b.AutoFinishVoting,
			&b.PointAverageRounding,
			&b.CreatedDate,
			&b.UpdatedDate,
			&leaders,
		); err != nil {
			d.Logger.Error("get poker games query error", zap.Error(err))
		} else {
			_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
			_ = json.Unmarshal([]byte(leaders), &b.Facilitators)
			b.ActiveStoryID = ActivePlanID.String
			battles = append(battles, b)
		}
	}

	return battles, Count, nil
}

// GetActiveGames gets a list of active games
func (d *PokerService) GetActiveGames(Limit int, Offset int) ([]*thunderdome.Poker, int, error) {
	var battles = make([]*thunderdome.Poker, 0)
	var Count int

	e := d.DB.QueryRow(
		"SELECT COUNT(DISTINCT pu.poker_id) FROM thunderdome.poker_user pu WHERE pu.active IS TRUE;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	battleRows, battlesErr := d.DB.Query(`
		SELECT b.id, b.name, b.voting_locked, b.active_story_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, b.created_date, b.updated_date,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS leaders
		FROM thunderdome.poker_user bu
		LEFT JOIN thunderdome.poker b ON b.id = bu.poker_id
		LEFT JOIN thunderdome.poker_facilitator bl ON b.id = bl.poker_id
		WHERE bu.active IS TRUE GROUP BY b.id
		LIMIT $1 OFFSET $2;
	`, Limit, Offset)
	if battlesErr != nil {
		return nil, Count, battlesErr
	}

	defer battleRows.Close()
	for battleRows.Next() {
		var pv string
		var leaders string
		var ActivePlanID sql.NullString
		var b = &thunderdome.Poker{
			Users:              make([]*thunderdome.PokerUser, 0),
			Stories:            make([]*thunderdome.Story, 0),
			VotingLocked:       true,
			PointValuesAllowed: make([]string, 0),
			AutoFinishVoting:   true,
			Facilitators:       make([]string, 0),
		}
		if err := battleRows.Scan(
			&b.Id,
			&b.Name,
			&b.VotingLocked,
			&ActivePlanID,
			&pv,
			&b.AutoFinishVoting,
			&b.PointAverageRounding,
			&b.CreatedDate,
			&b.UpdatedDate,
			&leaders,
		); err != nil {
			d.Logger.Error("get active poker games query error", zap.Error(err))
		} else {
			_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
			_ = json.Unmarshal([]byte(leaders), &b.Facilitators)
			b.ActiveStoryID = ActivePlanID.String
			battles = append(battles, b)
		}
	}

	return battles, Count, nil
}

// PurgeOldGames deletes games older than {DaysOld} days
func (d *PokerService) PurgeOldGames(ctx context.Context, DaysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.poker WHERE last_active < (NOW() - $1 * interval '1 day');`,
		DaysOld,
	); err != nil {
		return fmt.Errorf("error attempting to clean poker games: %v", err)
	}

	return nil
}
