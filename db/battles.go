package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/microcosm-cc/bluemonday"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"strings"

	"go.uber.org/zap"
)

// BattleService represents a PostgreSQL implementation of thunderdome.BattleService.
type BattleService struct {
	DB                  *sql.DB
	Logger              *otelzap.Logger
	AESHashKey          string
	HTMLSanitizerPolicy *bluemonday.Policy
}

//CreateBattle creates a new story pointing session (battle)
func (d *BattleService) CreateBattle(ctx context.Context, LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*thunderdome.Plan, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string, HideVoterIdentity bool) (*thunderdome.Battle, error) {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to create battle join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if LeaderCode != "" {
		EncryptedCode, codeErr := encrypt(LeaderCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to create battle leader_code")
		}
		encryptedLeaderCode = EncryptedCode
	}

	var b = &thunderdome.Battle{
		Name:                 BattleName,
		Users:                make([]*thunderdome.BattleUser, 0),
		Plans:                make([]*thunderdome.Plan, 0),
		VotingLocked:         true,
		PointValuesAllowed:   PointValuesAllowed,
		AutoFinishVoting:     AutoFinishVoting,
		PointAverageRounding: PointAverageRounding,
		HideVoterIdentity:    HideVoterIdentity,
		Leaders:              make([]string, 0),
		JoinCode:             JoinCode,
		LeaderCode:           LeaderCode,
	}
	b.Leaders = append(b.Leaders, LeaderID)

	e := d.DB.QueryRowContext(ctx,
		`SELECT battleId FROM create_battle($1, $2, $3, $4, $5, $6, $7, $8);`,
		LeaderID,
		BattleName,
		string(pointValuesJSON),
		AutoFinishVoting,
		PointAverageRounding,
		HideVoterIdentity,
		encryptedJoinCode,
		encryptedLeaderCode,
	).Scan(&b.Id)
	if e != nil {
		d.Logger.Error("create_battle query error", zap.Error(e))
		return nil, errors.New("error creating battle")
	}

	for _, plan := range Plans {
		plan.Votes = make([]*thunderdome.Vote, 0)

		e := d.DB.QueryRowContext(ctx,
			`INSERT INTO plans (battle_id, name, type, reference_id, link, description, acceptance_criteria) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			b.Id,
			plan.Name,
			plan.Type,
			plan.ReferenceId,
			plan.Link,
			plan.Description,
			plan.AcceptanceCriteria,
		).Scan(&plan.Id)
		if e != nil {
			d.Logger.Error("insert plans error", zap.Error(e))
		}
	}

	b.Plans = Plans

	return b, nil
}

//TeamCreateBattle creates a new story pointing session (battle) associated to a team
func (d *BattleService) TeamCreateBattle(ctx context.Context, TeamID string, LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*thunderdome.Plan, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string, HideVoterIdentity bool) (*thunderdome.Battle, error) {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to create battle join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if LeaderCode != "" {
		EncryptedCode, codeErr := encrypt(LeaderCode, d.AESHashKey)
		if codeErr != nil {
			return nil, errors.New("unable to create battle leader_code")
		}
		encryptedLeaderCode = EncryptedCode
	}

	var b = &thunderdome.Battle{
		Name:                 BattleName,
		Users:                make([]*thunderdome.BattleUser, 0),
		Plans:                make([]*thunderdome.Plan, 0),
		VotingLocked:         true,
		PointValuesAllowed:   PointValuesAllowed,
		AutoFinishVoting:     AutoFinishVoting,
		PointAverageRounding: PointAverageRounding,
		HideVoterIdentity:    HideVoterIdentity,
		Leaders:              make([]string, 0),
		JoinCode:             JoinCode,
		LeaderCode:           LeaderCode,
	}
	b.Leaders = append(b.Leaders, LeaderID)

	e := d.DB.QueryRowContext(ctx,
		`SELECT battleId FROM team_create_battle($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		TeamID,
		LeaderID,
		BattleName,
		string(pointValuesJSON),
		AutoFinishVoting,
		PointAverageRounding,
		HideVoterIdentity,
		encryptedJoinCode,
		encryptedLeaderCode,
	).Scan(&b.Id)
	if e != nil {
		d.Logger.Error("team_create_battle query error", zap.Error(e))
		return nil, errors.New("error creating battle")
	}

	for _, plan := range Plans {
		plan.Votes = make([]*thunderdome.Vote, 0)

		e := d.DB.QueryRowContext(ctx,
			`INSERT INTO plans (battle_id, name, type, reference_id, link, description, acceptance_criteria) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			b.Id,
			plan.Name,
			plan.Type,
			plan.ReferenceId,
			plan.Link,
			plan.Description,
			plan.AcceptanceCriteria,
		).Scan(&plan.Id)
		if e != nil {
			d.Logger.Error("insert plans error", zap.Error(e))
		}
	}

	b.Plans = Plans

	return b, nil
}

// ReviseBattle updates the battle by ID
func (d *BattleService) ReviseBattle(BattleID string, BattleName string, PointValuesAllowed []string, AutoFinishVoting bool, PointAverageRounding string, HideVoterIdentity bool, JoinCode string, LeaderCode string) error {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.AESHashKey)
		if codeErr != nil {
			return errors.New("unable to revise battle join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if LeaderCode != "" {
		EncryptedCode, codeErr := encrypt(LeaderCode, d.AESHashKey)
		if codeErr != nil {
			return errors.New("unable to revise battle leadercode")
		}
		encryptedLeaderCode = EncryptedCode
	}

	if _, err := d.DB.Exec(`
		UPDATE battles
		SET name = $2, point_values_allowed = $3, auto_finish_voting = $4, point_average_rounding = $5, hide_voter_identity = $6, join_code = $7, leader_code = $8, updated_date = NOW()
		WHERE id = $1`,
		BattleID, BattleName, string(pointValuesJSON), AutoFinishVoting, PointAverageRounding, HideVoterIdentity, encryptedJoinCode, encryptedLeaderCode,
	); err != nil {
		d.Logger.Error("update battle error", zap.Error(err))
		return errors.New("unable to revise battle")
	}

	return nil
}

// GetBattleLeaderCode retrieve the battle leader_code
func (d *BattleService) GetBattleLeaderCode(BattleID string) (string, error) {
	var EncryptedLeaderCode string

	if err := d.DB.QueryRow(`
		SELECT COALESCE(leader_code, '') FROM battles
		WHERE id = $1`,
		BattleID,
	).Scan(&EncryptedLeaderCode); err != nil {
		d.Logger.Error("get battle leadercode error", zap.Error(err))
		return "", errors.New("unable to retrieve battle leader_code")
	}

	if EncryptedLeaderCode == "" {
		return "", errors.New("unable to retrieve battle leader_code")
	}
	DecryptedCode, codeErr := decrypt(EncryptedLeaderCode, d.AESHashKey)
	if codeErr != nil {
		return "", errors.New("unable to retrieve battle leader_code")
	}

	return DecryptedCode, nil
}

// GetBattle gets a battle by ID
func (d *BattleService) GetBattle(BattleID string, UserID string) (*thunderdome.Battle, error) {
	var b = &thunderdome.Battle{
		Id:                 BattleID,
		Users:              make([]*thunderdome.BattleUser, 0),
		Plans:              make([]*thunderdome.Plan, 0),
		VotingLocked:       true,
		PointValuesAllowed: make([]string, 0),
		AutoFinishVoting:   true,
		Leaders:            make([]string, 0),
	}

	// get battle
	var ActivePlanID sql.NullString
	var pv string
	var leaders string
	var JoinCode string
	var LeaderCode string
	e := d.DB.QueryRow(
		`
		SELECT b.id, b.name, b.voting_locked, b.active_plan_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, b.hide_voter_identity, COALESCE(b.join_code, ''), COALESCE(b.leader_code, ''), b.created_date, b.updated_date,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS leaders
		FROM battles b
		LEFT JOIN battles_leaders bl ON b.id = bl.battle_id
		WHERE b.id = $1
		GROUP BY b.id`,
		BattleID,
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
		d.Logger.Error("error getting battle", zap.Error(e))
		return nil, errors.New("not found")
	}

	_ = json.Unmarshal([]byte(leaders), &b.Leaders)
	_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
	b.ActivePlanID = ActivePlanID.String

	isBattleLeader := contains(b.Leaders, UserID)

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
		b.LeaderCode = DecryptedCode
	}

	b.Users = d.GetBattleUsers(BattleID)
	b.Plans = d.GetPlans(BattleID, UserID)

	return b, nil
}

// GetBattlesByUser gets a list of battles by UserID
func (d *BattleService) GetBattlesByUser(UserID string, Limit int, Offset int) ([]*thunderdome.Battle, int, error) {
	var Count int
	var battles = make([]*thunderdome.Battle, 0)

	e := d.DB.QueryRow(`
		SELECT COUNT(*) FROM battles b
		LEFT JOIN battles_users bw ON b.id = bw.battle_id
		WHERE bw.user_id = $1 AND bw.abandoned = false;
	`, UserID).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	battleRows, battlesErr := d.DB.Query(`
		SELECT b.id, b.name, b.voting_locked, b.active_plan_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, b.created_date, b.updated_date,
		CASE WHEN COUNT(p) = 0 THEN '[]'::json ELSE array_to_json(array_agg(row_to_json(p))) END AS plans,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS leaders
		FROM battles b
		LEFT JOIN plans p ON b.id = p.battle_id
		LEFT JOIN battles_leaders bl ON b.id = bl.battle_id
		LEFT JOIN battles_users bw ON b.id = bw.battle_id
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
		var b = &thunderdome.Battle{
			Users:              make([]*thunderdome.BattleUser, 0),
			Plans:              make([]*thunderdome.Plan, 0),
			VotingLocked:       true,
			PointValuesAllowed: make([]string, 0),
			AutoFinishVoting:   true,
			Leaders:            make([]string, 0),
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
			d.Logger.Error("error getting battle by user", zap.Error(e))
		} else {
			_ = json.Unmarshal([]byte(plans), &b.Plans)
			_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
			_ = json.Unmarshal([]byte(leaders), &b.Leaders)
			b.ActivePlanID = ActivePlanID.String
			battles = append(battles, b)
		}
	}

	return battles, Count, nil
}

// ConfirmLeader confirms the user is a leader of the battle
func (d *BattleService) ConfirmLeader(BattleID string, UserID string) error {
	var leaderID string
	var role string
	err := d.DB.QueryRow("SELECT type FROM users WHERE id = $1", UserID).Scan(&role)
	if err != nil {
		d.Logger.Error("error getting user role", zap.Error(err))
		return errors.New("unable to get user role")
	}

	e := d.DB.QueryRow("SELECT user_id FROM battles_leaders WHERE battle_id = $1 AND user_id = $2", BattleID, UserID).Scan(&leaderID)
	if e != nil && role != "ADMIN" {
		d.Logger.Error("error confirming battle leader", zap.Error(e))
		return errors.New("not a battle leader")
	}

	return nil
}

// GetBattleUserActiveStatus checks battle active status of User for given battle
func (d *BattleService) GetBattleUserActiveStatus(BattleID string, UserID string) error {
	var active bool

	e := d.DB.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM battles_users
		WHERE user_id = $2 AND battle_id = $1;`,
		BattleID,
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

// GetBattleUsers retrieves the users for a given battle
func (d *BattleService) GetBattleUsers(BattleID string) []*thunderdome.BattleUser {
	var users = make([]*thunderdome.BattleUser, 0)
	rows, err := d.DB.Query(
		`SELECT
			w.id, w.name, w.type, w.avatar, bw.active, bw.spectator, COALESCE(w.email, '')
		FROM battles_users bw
		LEFT JOIN users w ON bw.user_id = w.id
		WHERE bw.battle_id = $1
		ORDER BY w.name`,
		BattleID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w thunderdome.BattleUser
			if err := rows.Scan(&w.Id, &w.Name, &w.Type, &w.Avatar, &w.Active, &w.Spectator, &w.GravatarHash); err != nil {
				d.Logger.Error("error getting battle users", zap.Error(err))
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

// GetBattleActiveUsers retrieves the active users for a given battle
func (d *BattleService) GetBattleActiveUsers(BattleID string) []*thunderdome.BattleUser {
	var users = make([]*thunderdome.BattleUser, 0)
	rows, err := d.DB.Query(
		`SELECT
			w.id, w.name, w.type, w.avatar, bw.active, bw.spectator, COALESCE(w.email, '')
		FROM battles_users bw
		LEFT JOIN users w ON bw.user_id = w.id
		WHERE bw.battle_id = $1 AND bw.active = true
		ORDER BY w.name`,
		BattleID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w thunderdome.BattleUser
			if err := rows.Scan(&w.Id, &w.Name, &w.Type, &w.Avatar, &w.Active, &w.Spectator, &w.GravatarHash); err != nil {
				d.Logger.Error("error getting active battle users", zap.Error(err))
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

// AddUserToBattle adds a user by ID to the battle by ID
func (d *BattleService) AddUserToBattle(BattleID string, UserID string) ([]*thunderdome.BattleUser, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO battles_users (battle_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (battle_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		BattleID,
		UserID,
	); err != nil {
		d.Logger.Error("error adding user to battle", zap.Error(err))
	}

	users := d.GetBattleUsers(BattleID)

	return users, nil
}

// RetreatUser removes a user from the current battle by ID
func (d *BattleService) RetreatUser(BattleID string, UserID string) []*thunderdome.BattleUser {
	if _, err := d.DB.Exec(
		`UPDATE battles_users SET active = false WHERE battle_id = $1 AND user_id = $2`, BattleID, UserID); err != nil {
		d.Logger.Error("error updating battle user to active false", zap.Error(err))
	}

	if _, err := d.DB.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("error updating user last active timestamp", zap.Error(err))
	}

	users := d.GetBattleUsers(BattleID)

	return users
}

// AbandonBattle removes a user from the current battle by ID and sets abandoned true
func (d *BattleService) AbandonBattle(BattleID string, UserID string) ([]*thunderdome.BattleUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE battles_users SET active = false, abandoned = true WHERE battle_id = $1 AND user_id = $2`, BattleID, UserID); err != nil {
		d.Logger.Error("error updating battle user to abandoned", zap.Error(err))
		return nil, err
	}

	if _, err := d.DB.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("error updating user last active timestamp", zap.Error(err))
		return nil, err
	}

	users := d.GetBattleUsers(BattleID)

	return users, nil
}

// SetBattleLeader sets the leaderId for the battle
func (d *BattleService) SetBattleLeader(BattleID string, LeaderID string) ([]string, error) {
	leaders := make([]string, 0)

	// set battle leader
	if _, err := d.DB.Exec(
		`call set_battle_leader($1, $2);`, BattleID, LeaderID); err != nil {
		d.Logger.Error("call set_battle_leader query error", zap.Error(err))
		return nil, errors.New("unable to promote leader")
	}

	leaderRows, leadersErr := d.DB.Query(`
		SELECT user_id FROM battles_leaders WHERE battle_id = $1;
	`, BattleID)
	if leadersErr != nil {
		return leaders, nil
	}

	defer leaderRows.Close()
	for leaderRows.Next() {
		var leader string
		if err := leaderRows.Scan(
			&leader,
		); err != nil {
			d.Logger.Error("battles_leaders query scan error", zap.Error(err))
		} else {
			leaders = append(leaders, leader)
		}
	}

	return leaders, nil
}

// DemoteBattleLeader removes a user from battle leaders
func (d *BattleService) DemoteBattleLeader(BattleID string, LeaderID string) ([]string, error) {
	leaders := make([]string, 0)

	// set battle leader
	if _, err := d.DB.Exec(
		`call demote_battle_leader($1, $2);`, BattleID, LeaderID); err != nil {
		d.Logger.Error("call demote_battle_leader query error", zap.Error(err))
		return nil, errors.New("unable to demote leader")
	}

	leaderRows, leadersErr := d.DB.Query(`
		SELECT user_id FROM battles_leaders WHERE battle_id = $1;
	`, BattleID)
	if leadersErr != nil {
		return leaders, nil
	}

	defer leaderRows.Close()
	for leaderRows.Next() {
		var leader string
		if err := leaderRows.Scan(
			&leader,
		); err != nil {
			d.Logger.Error("battles_leaders query scan error", zap.Error(err))
		} else {
			leaders = append(leaders, leader)
		}
	}

	return leaders, nil
}

// ToggleSpectator changes a battle users spectator status
func (d *BattleService) ToggleSpectator(BattleID string, UserID string, Spectator bool) ([]*thunderdome.BattleUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE battles_users SET spectator = $3 WHERE battle_id = $1 AND user_id = $2`, BattleID, UserID, Spectator); err != nil {
		d.Logger.Error("update battle user spectator error", zap.Error(err))
		return nil, err
	}

	if _, err := d.DB.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("error updating user last active timestamp", zap.Error(err))
	}

	users := d.GetBattleUsers(BattleID)

	return users, nil
}

// DeleteBattle removes all battle associations and the battle itself by BattleID
func (d *BattleService) DeleteBattle(BattleID string) error {
	if _, err := d.DB.Exec(
		`call delete_battle($1);`, BattleID); err != nil {
		d.Logger.Error("delete battle error", zap.Error(err))
		return err
	}

	return nil
}

// AddBattleLeadersByEmail adds additional battle leaders using provided emails for matches
func (d *BattleService) AddBattleLeadersByEmail(ctx context.Context, BattleID string, LeaderEmails []string) ([]string, error) {
	var leaders string
	var newLeaders []string

	for i, email := range LeaderEmails {
		LeaderEmails[i] = sanitizeEmail(email)
	}
	emails := strings.Join(LeaderEmails[:], ",")

	e := d.DB.QueryRowContext(ctx,
		`select leaders FROM add_battle_leaders_by_email($1, $2);`, BattleID, emails,
	).Scan(&leaders)
	if e != nil {
		d.Logger.Error("add_battle_leaders_by_email query error", zap.Error(e))
		return nil, errors.New("error creating battle")
	}

	_ = json.Unmarshal([]byte(leaders), &newLeaders)

	return newLeaders, nil
}

// GetBattles gets a list of battles
func (d *BattleService) GetBattles(Limit int, Offset int) ([]*thunderdome.Battle, int, error) {
	var battles = make([]*thunderdome.Battle, 0)
	var Count int

	e := d.DB.QueryRow(
		"SELECT COUNT(*) FROM battles;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	battleRows, battlesErr := d.DB.Query(`
		SELECT b.id, b.name, b.voting_locked, b.active_plan_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, b.created_date, b.updated_date,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS leaders
		FROM battles b
		LEFT JOIN battles_leaders bl ON b.id = bl.battle_id
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
		var b = &thunderdome.Battle{
			Users:              make([]*thunderdome.BattleUser, 0),
			Plans:              make([]*thunderdome.Plan, 0),
			VotingLocked:       true,
			PointValuesAllowed: make([]string, 0),
			AutoFinishVoting:   true,
			Leaders:            make([]string, 0),
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
			d.Logger.Error("get battles query error", zap.Error(err))
		} else {
			_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
			_ = json.Unmarshal([]byte(leaders), &b.Leaders)
			b.ActivePlanID = ActivePlanID.String
			battles = append(battles, b)
		}
	}

	return battles, Count, nil
}

// GetActiveBattles gets a list of active battles
func (d *BattleService) GetActiveBattles(Limit int, Offset int) ([]*thunderdome.Battle, int, error) {
	var battles = make([]*thunderdome.Battle, 0)
	var Count int

	e := d.DB.QueryRow(
		"SELECT COUNT(DISTINCT bu.battle_id) FROM battles_users bu WHERE bu.active IS TRUE;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	battleRows, battlesErr := d.DB.Query(`
		SELECT b.id, b.name, b.voting_locked, b.active_plan_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, b.created_date, b.updated_date,
		CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS leaders
		FROM battles_users bu
		LEFT JOIN battles b ON b.id = bu.battle_id
		LEFT JOIN battles_leaders bl ON b.id = bl.battle_id
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
		var b = &thunderdome.Battle{
			Users:              make([]*thunderdome.BattleUser, 0),
			Plans:              make([]*thunderdome.Plan, 0),
			VotingLocked:       true,
			PointValuesAllowed: make([]string, 0),
			AutoFinishVoting:   true,
			Leaders:            make([]string, 0),
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
			d.Logger.Error("get active battles query error", zap.Error(err))
		} else {
			_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
			_ = json.Unmarshal([]byte(leaders), &b.Leaders)
			b.ActivePlanID = ActivePlanID.String
			battles = append(battles, b)
		}
	}

	return battles, Count, nil
}

// CleanBattles deletes battles older than {DaysOld} days
func (d *BattleService) CleanBattles(ctx context.Context, DaysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`call clean_battles($1);`,
		DaysOld,
	); err != nil {
		d.Logger.Ctx(ctx).Error("call clean_battles", zap.Error(err))
		return errors.New("error attempting to clean battles")
	}

	return nil
}
