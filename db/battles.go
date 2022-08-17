package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

//CreateBattle creates a new story pointing session (battle)
func (d *Database) CreateBattle(ctx context.Context, LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*model.Plan, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string) (*model.Battle, error) {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, errors.New("unable to create battle join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if LeaderCode != "" {
		EncryptedCode, codeErr := encrypt(LeaderCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, errors.New("unable to create battle leader_code")
		}
		encryptedLeaderCode = EncryptedCode
	}

	var b = &model.Battle{
		Name:                 BattleName,
		Users:                make([]*model.BattleUser, 0),
		Plans:                make([]*model.Plan, 0),
		VotingLocked:         true,
		PointValuesAllowed:   PointValuesAllowed,
		AutoFinishVoting:     AutoFinishVoting,
		PointAverageRounding: PointAverageRounding,
		Leaders:              make([]string, 0),
		JoinCode:             JoinCode,
		LeaderCode:           LeaderCode,
	}
	b.Leaders = append(b.Leaders, LeaderID)

	e := d.db.QueryRowContext(ctx,
		`SELECT battleId FROM create_battle($1, $2, $3, $4, $5, $6, $7);`,
		LeaderID,
		BattleName,
		string(pointValuesJSON),
		AutoFinishVoting,
		PointAverageRounding,
		encryptedJoinCode,
		encryptedLeaderCode,
	).Scan(&b.Id)
	if e != nil {
		d.logger.Error("create_battle query error", zap.Error(e))
		return nil, errors.New("error creating battle")
	}

	for _, plan := range Plans {
		plan.Votes = make([]*model.Vote, 0)

		e := d.db.QueryRowContext(ctx,
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
			d.logger.Error("insert plans error", zap.Error(e))
		}
	}

	b.Plans = Plans

	return b, nil
}

//TeamCreateBattle creates a new story pointing session (battle) associated to a team
func (d *Database) TeamCreateBattle(ctx context.Context, TeamID string, LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*model.Plan, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string) (*model.Battle, error) {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, errors.New("unable to create battle join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if LeaderCode != "" {
		EncryptedCode, codeErr := encrypt(LeaderCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, errors.New("unable to create battle leader_code")
		}
		encryptedLeaderCode = EncryptedCode
	}

	var b = &model.Battle{
		Name:                 BattleName,
		Users:                make([]*model.BattleUser, 0),
		Plans:                make([]*model.Plan, 0),
		VotingLocked:         true,
		PointValuesAllowed:   PointValuesAllowed,
		AutoFinishVoting:     AutoFinishVoting,
		PointAverageRounding: PointAverageRounding,
		Leaders:              make([]string, 0),
		JoinCode:             JoinCode,
		LeaderCode:           LeaderCode,
	}
	b.Leaders = append(b.Leaders, LeaderID)

	e := d.db.QueryRowContext(ctx,
		`SELECT battleId FROM team_create_battle($1, $2, $3, $4, $5, $6, $7, $8);`,
		TeamID,
		LeaderID,
		BattleName,
		string(pointValuesJSON),
		AutoFinishVoting,
		PointAverageRounding,
		encryptedJoinCode,
		encryptedLeaderCode,
	).Scan(&b.Id)
	if e != nil {
		d.logger.Error("team_create_battle query error", zap.Error(e))
		return nil, errors.New("error creating battle")
	}

	for _, plan := range Plans {
		plan.Votes = make([]*model.Vote, 0)

		e := d.db.QueryRowContext(ctx,
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
			d.logger.Error("insert plans error", zap.Error(e))
		}
	}

	b.Plans = Plans

	return b, nil
}

// ReviseBattle updates the battle by ID
func (d *Database) ReviseBattle(BattleID string, BattleName string, PointValuesAllowed []string, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string) error {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)
	var encryptedJoinCode string
	var encryptedLeaderCode string

	if JoinCode != "" {
		EncryptedCode, codeErr := encrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return errors.New("unable to revise battle join_code")
		}
		encryptedJoinCode = EncryptedCode
	}

	if LeaderCode != "" {
		EncryptedCode, codeErr := encrypt(LeaderCode, d.config.AESHashkey)
		if codeErr != nil {
			return errors.New("unable to revise battle leadercode")
		}
		encryptedLeaderCode = EncryptedCode
	}

	if _, err := d.db.Exec(`
		UPDATE battles
		SET name = $2, point_values_allowed = $3, auto_finish_voting = $4, point_average_rounding = $5, join_code = $6, leader_code = $7, updated_date = NOW()
		WHERE id = $1`,
		BattleID, BattleName, string(pointValuesJSON), AutoFinishVoting, PointAverageRounding, encryptedJoinCode, encryptedLeaderCode,
	); err != nil {
		d.logger.Error("update battle error", zap.Error(err))
		return errors.New("unable to revise battle")
	}

	return nil
}

// GetBattleLeaderCode retrieve the battle leader_code
func (d *Database) GetBattleLeaderCode(BattleID string) (string, error) {
	var EncryptedLeaderCode string

	if err := d.db.QueryRow(`
		SELECT COALESCE(leader_code, '') FROM battles
		WHERE id = $1`,
		BattleID,
	).Scan(&EncryptedLeaderCode); err != nil {
		d.logger.Error("get battle leadercode error", zap.Error(err))
		return "", errors.New("unable to retrieve battle leader_code")
	}

	if EncryptedLeaderCode == "" {
		return "", errors.New("unable to retrieve battle leader_code")
	}
	DecryptedCode, codeErr := decrypt(EncryptedLeaderCode, d.config.AESHashkey)
	if codeErr != nil {
		return "", errors.New("unable to retrieve battle leader_code")
	}

	return DecryptedCode, nil
}

// GetBattle gets a battle by ID
func (d *Database) GetBattle(BattleID string, UserID string) (*model.Battle, error) {
	var b = &model.Battle{
		Id:                 BattleID,
		Users:              make([]*model.BattleUser, 0),
		Plans:              make([]*model.Plan, 0),
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
	e := d.db.QueryRow(
		`
		SELECT b.id, b.name, b.voting_locked, b.active_plan_id, b.point_values_allowed, b.auto_finish_voting, b.point_average_rounding, COALESCE(b.join_code, ''), COALESCE(b.leader_code, ''), b.created_date, b.updated_date,
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
		&JoinCode,
		&LeaderCode,
		&b.CreatedDate,
		&b.UpdatedDate,
		&leaders,
	)
	if e != nil {
		d.logger.Error("error getting battle", zap.Error(e))
		return nil, errors.New("not found")
	}

	_ = json.Unmarshal([]byte(leaders), &b.Leaders)
	_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
	b.ActivePlanID = ActivePlanID.String

	isBattleLeader := contains(b.Leaders, UserID)

	if JoinCode != "" {
		DecryptedCode, codeErr := decrypt(JoinCode, d.config.AESHashkey)
		if codeErr != nil {
			return nil, errors.New("unable to decode join_code")
		}
		b.JoinCode = DecryptedCode
	}

	if LeaderCode != "" && isBattleLeader {
		DecryptedCode, codeErr := decrypt(LeaderCode, d.config.AESHashkey)
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
func (d *Database) GetBattlesByUser(UserID string, Limit int, Offset int) ([]*model.Battle, int, error) {
	var Count int
	var battles = make([]*model.Battle, 0)

	e := d.db.QueryRow(`
		SELECT COUNT(*) FROM battles b
		LEFT JOIN battles_users bw ON b.id = bw.battle_id
		WHERE bw.user_id = $1 AND bw.abandoned = false;
	`, UserID).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	battleRows, battlesErr := d.db.Query(`
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
		var b = &model.Battle{
			Users:              make([]*model.BattleUser, 0),
			Plans:              make([]*model.Plan, 0),
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
			d.logger.Error("error getting battle by user", zap.Error(e))
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
func (d *Database) ConfirmLeader(BattleID string, UserID string) error {
	var leaderID string
	e := d.db.QueryRow("SELECT user_id FROM battles_leaders WHERE battle_id = $1 AND user_id = $2", BattleID, UserID).Scan(&leaderID)
	if e != nil {
		d.logger.Error("error confirming battle leader", zap.Error(e))
		return errors.New("not a battle leader")
	}

	return nil
}

// GetBattleUserActiveStatus checks battle active status of User for given battle
func (d *Database) GetBattleUserActiveStatus(BattleID string, UserID string) error {
	var active bool

	e := d.db.QueryRow(`
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
func (d *Database) GetBattleUsers(BattleID string) []*model.BattleUser {
	var users = make([]*model.BattleUser, 0)
	rows, err := d.db.Query(
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
			var w model.BattleUser
			if err := rows.Scan(&w.Id, &w.Name, &w.Type, &w.Avatar, &w.Active, &w.Spectator, &w.GravatarHash); err != nil {
				d.logger.Error("error getting battle users", zap.Error(err))
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
func (d *Database) GetBattleActiveUsers(BattleID string) []*model.BattleUser {
	var users = make([]*model.BattleUser, 0)
	rows, err := d.db.Query(
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
			var w model.BattleUser
			if err := rows.Scan(&w.Id, &w.Name, &w.Type, &w.Avatar, &w.Active, &w.Spectator, &w.GravatarHash); err != nil {
				d.logger.Error("error getting active battle users", zap.Error(err))
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
func (d *Database) AddUserToBattle(BattleID string, UserID string) ([]*model.BattleUser, error) {
	if _, err := d.db.Exec(
		`INSERT INTO battles_users (battle_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (battle_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		BattleID,
		UserID,
	); err != nil {
		d.logger.Error("error adding user to battle", zap.Error(err))
	}

	users := d.GetBattleUsers(BattleID)

	return users, nil
}

// RetreatUser removes a user from the current battle by ID
func (d *Database) RetreatUser(BattleID string, UserID string) []*model.BattleUser {
	if _, err := d.db.Exec(
		`UPDATE battles_users SET active = false WHERE battle_id = $1 AND user_id = $2`, BattleID, UserID); err != nil {
		d.logger.Error("error updating battle user to active false", zap.Error(err))
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.logger.Error("error updating user last active timestamp", zap.Error(err))
	}

	users := d.GetBattleUsers(BattleID)

	return users
}

// AbandonBattle removes a user from the current battle by ID and sets abandoned true
func (d *Database) AbandonBattle(BattleID string, UserID string) ([]*model.BattleUser, error) {
	if _, err := d.db.Exec(
		`UPDATE battles_users SET active = false, abandoned = true WHERE battle_id = $1 AND user_id = $2`, BattleID, UserID); err != nil {
		d.logger.Error("error updating battle user to abandoned", zap.Error(err))
		return nil, err
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.logger.Error("error updating user last active timestamp", zap.Error(err))
		return nil, err
	}

	users := d.GetBattleUsers(BattleID)

	return users, nil
}

// SetBattleLeader sets the leaderId for the battle
func (d *Database) SetBattleLeader(BattleID string, LeaderID string) ([]string, error) {
	leaders := make([]string, 0)

	// set battle leader
	if _, err := d.db.Exec(
		`call set_battle_leader($1, $2);`, BattleID, LeaderID); err != nil {
		d.logger.Error("call set_battle_leader query error", zap.Error(err))
		return nil, errors.New("unable to promote leader")
	}

	leaderRows, leadersErr := d.db.Query(`
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
			d.logger.Error("battles_leaders query scan error", zap.Error(err))
		} else {
			leaders = append(leaders, leader)
		}
	}

	return leaders, nil
}

// DemoteBattleLeader removes a user from battle leaders
func (d *Database) DemoteBattleLeader(BattleID string, LeaderID string) ([]string, error) {
	leaders := make([]string, 0)

	// set battle leader
	if _, err := d.db.Exec(
		`call demote_battle_leader($1, $2);`, BattleID, LeaderID); err != nil {
		d.logger.Error("call demote_battle_leader query error", zap.Error(err))
		return nil, errors.New("unable to demote leader")
	}

	leaderRows, leadersErr := d.db.Query(`
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
			d.logger.Error("battles_leaders query scan error", zap.Error(err))
		} else {
			leaders = append(leaders, leader)
		}
	}

	return leaders, nil
}

// ToggleSpectator changes a battle users spectator status
func (d *Database) ToggleSpectator(BattleID string, UserID string, Spectator bool) ([]*model.BattleUser, error) {
	if _, err := d.db.Exec(
		`UPDATE battles_users SET spectator = $3 WHERE battle_id = $1 AND user_id = $2`, BattleID, UserID, Spectator); err != nil {
		d.logger.Error("update battle user spectator error", zap.Error(err))
		return nil, err
	}

	if _, err := d.db.Exec(
		`UPDATE users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.logger.Error("error updating user last active timestamp", zap.Error(err))
	}

	users := d.GetBattleUsers(BattleID)

	return users, nil
}

// DeleteBattle removes all battle associations and the battle itself by BattleID
func (d *Database) DeleteBattle(BattleID string) error {
	if _, err := d.db.Exec(
		`call delete_battle($1);`, BattleID); err != nil {
		d.logger.Error("delete battle error", zap.Error(err))
		return err
	}

	return nil
}

// AddBattleLeadersByEmail adds additional battle leaders using provided emails for matches
func (d *Database) AddBattleLeadersByEmail(ctx context.Context, BattleID string, LeaderEmails []string) ([]string, error) {
	var leaders string
	var newLeaders []string

	for i, email := range LeaderEmails {
		LeaderEmails[i] = sanitizeEmail(email)
	}
	emails := strings.Join(LeaderEmails[:], ",")

	e := d.db.QueryRowContext(ctx,
		`select leaders FROM add_battle_leaders_by_email($1, $2);`, BattleID, emails,
	).Scan(&leaders)
	if e != nil {
		d.logger.Error("add_battle_leaders_by_email query error", zap.Error(e))
		return nil, errors.New("error creating battle")
	}

	_ = json.Unmarshal([]byte(leaders), &newLeaders)

	return newLeaders, nil
}

// GetBattles gets a list of battles
func (d *Database) GetBattles(Limit int, Offset int) ([]*model.Battle, int, error) {
	var battles = make([]*model.Battle, 0)
	var Count int

	e := d.db.QueryRow(
		"SELECT COUNT(*) FROM battles;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	battleRows, battlesErr := d.db.Query(`
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
		var b = &model.Battle{
			Users:              make([]*model.BattleUser, 0),
			Plans:              make([]*model.Plan, 0),
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
			d.logger.Error("get battles query error", zap.Error(err))
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
func (d *Database) GetActiveBattles(Limit int, Offset int) ([]*model.Battle, int, error) {
	var battles = make([]*model.Battle, 0)
	var Count int

	e := d.db.QueryRow(
		"SELECT COUNT(DISTINCT bu.battle_id) FROM battles_users bu WHERE bu.active IS TRUE;",
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	battleRows, battlesErr := d.db.Query(`
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
		var b = &model.Battle{
			Users:              make([]*model.BattleUser, 0),
			Plans:              make([]*model.Plan, 0),
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
			d.logger.Error("get active battles query error", zap.Error(err))
		} else {
			_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
			_ = json.Unmarshal([]byte(leaders), &b.Leaders)
			b.ActivePlanID = ActivePlanID.String
			battles = append(battles, b)
		}
	}

	return battles, Count, nil
}
