package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
)

//CreateBattle adds a new battle to the db
func (d *Database) CreateBattle(LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*Plan, AutoFinishVoting bool) (*Battle, error) {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)

	var b = &Battle{
		BattleID:           "",
		LeaderID:           LeaderID,
		BattleName:         BattleName,
		Warriors:           make([]*BattleWarrior, 0),
		Plans:              make([]*Plan, 0),
		VotingLocked:       true,
		ActivePlanID:       "",
		PointValuesAllowed: PointValuesAllowed,
		AutoFinishVoting:   AutoFinishVoting,
	}

	e := d.db.QueryRow(
		`INSERT INTO battles (leader_id, name, point_values_allowed, auto_finish_voting) VALUES ($1, $2, $3, $4) RETURNING id`,
		LeaderID,
		BattleName,
		string(pointValuesJSON),
		AutoFinishVoting,
	).Scan(&b.BattleID)
	if e != nil {
		log.Println(e)
		return nil, errors.New("error creating battle")
	}

	for _, plan := range Plans {
		plan.Votes = make([]*Vote, 0)

		e := d.db.QueryRow(
			`INSERT INTO plans (battle_id, name, type, reference_id, link, description, acceptance_criteria) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			b.BattleID,
			plan.PlanName,
			plan.Type,
			plan.ReferenceID,
			plan.Link,
			plan.Description,
			plan.AcceptanceCriteria,
		).Scan(&plan.PlanID)
		if e != nil {
			log.Println(e)
		}
	}

	b.Plans = Plans

	return b, nil
}

// ReviseBattle updates the battle by ID
func (d *Database) ReviseBattle(BattleID string, warriorID string, BattleName string, PointValuesAllowed []string, AutoFinishVoting bool) error {
	err := d.ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return errors.New("incorrect permissions")
	}

	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)
	if _, err := d.db.Exec(
		`UPDATE battles SET name = $2, point_values_allowed = $3, auto_finish_voting = $4 WHERE id = $1`, BattleID, BattleName, string(pointValuesJSON), AutoFinishVoting); err != nil {
		log.Println(err)
		return errors.New("unable to revise battle")
	}

	return nil
}

// GetBattle gets a battle by ID
func (d *Database) GetBattle(BattleID string, WarriorID string) (*Battle, error) {
	var b = &Battle{
		BattleID:           BattleID,
		LeaderID:           "",
		BattleName:         "",
		Warriors:           make([]*BattleWarrior, 0),
		Plans:              make([]*Plan, 0),
		VotingLocked:       true,
		ActivePlanID:       "",
		PointValuesAllowed: make([]string, 0),
		AutoFinishVoting:   true,
	}

	// get battle
	var ActivePlanID sql.NullString
	var pv string
	e := d.db.QueryRow(
		"SELECT id, name, leader_id, voting_locked, active_plan_id, point_values_allowed, auto_finish_voting FROM battles WHERE id = $1",
		BattleID,
	).Scan(
		&b.BattleID,
		&b.BattleName,
		&b.LeaderID,
		&b.VotingLocked,
		&ActivePlanID,
		&pv,
		&b.AutoFinishVoting,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("not found")
	}

	_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
	b.ActivePlanID = ActivePlanID.String
	b.Warriors = d.GetBattleWarriors(BattleID)
	b.Plans = d.GetPlans(BattleID, WarriorID)

	return b, nil
}

// GetBattlesByWarrior gets a list of battles by WarriorID
func (d *Database) GetBattlesByWarrior(WarriorID string) ([]*Battle, error) {
	var battles = make([]*Battle, 0)
	battleRows, battlesErr := d.db.Query(`
		SELECT b.id, b.name, b.leader_id, b.voting_locked, b.active_plan_id, b.point_values_allowed, b.auto_finish_voting,
		CASE WHEN COUNT(p) = 0 THEN '[]'::json ELSE array_to_json(array_agg(row_to_json(p))) END AS plans
		FROM battles b
		LEFT JOIN plans p ON b.id = p.battle_id
		LEFT JOIN battles_warriors bw ON b.id = bw.battle_id WHERE bw.warrior_id = $1 AND bw.abandoned = false
		GROUP BY b.id ORDER BY b.created_date DESC
	`, WarriorID)
	if battlesErr != nil {
		return nil, errors.New("not found")
	}

	defer battleRows.Close()
	for battleRows.Next() {
		var plans string
		var pv string
		var ActivePlanID sql.NullString
		var b = &Battle{
			BattleID:           "",
			LeaderID:           "",
			BattleName:         "",
			Warriors:           make([]*BattleWarrior, 0),
			Plans:              make([]*Plan, 0),
			VotingLocked:       true,
			ActivePlanID:       "",
			PointValuesAllowed: make([]string, 0),
			AutoFinishVoting:   true,
		}
		if err := battleRows.Scan(
			&b.BattleID,
			&b.BattleName,
			&b.LeaderID,
			&b.VotingLocked,
			&ActivePlanID,
			&pv,
			&b.AutoFinishVoting,
			&plans,
		); err != nil {
			log.Println(err)
		} else {
			_ = json.Unmarshal([]byte(plans), &b.Plans)
			_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
			b.ActivePlanID = ActivePlanID.String
			battles = append(battles, b)
		}
	}

	return battles, nil
}

// ConfirmLeader confirms the warrior is infact leader of the battle
func (d *Database) ConfirmLeader(BattleID string, warriorID string) error {
	var leaderID string
	e := d.db.QueryRow("SELECT leader_id FROM battles WHERE id = $1", BattleID).Scan(&leaderID)
	if e != nil {
		log.Println(e)
		return errors.New("battle not found")
	}

	if leaderID != warriorID {
		return errors.New("not leader")
	}

	return nil
}

// GetBattleWarrior gets a warrior from db by ID and checks battle active status
func (d *Database) GetBattleWarrior(BattleID string, WarriorID string) (*BattleWarrior, error) {
	var active bool
	var w BattleWarrior

	e := d.db.QueryRow(
		`SELECT
			w.id, w.name, w.rank, w.avatar, coalesce(bw.active, FALSE)
		FROM warriors w
		LEFT JOIN battles_warriors bw ON bw.warrior_id = w.id AND bw.battle_id = $1
		WHERE id = $2`,
		BattleID,
		WarriorID,
	).Scan(
		&w.WarriorID,
		&w.WarriorName,
		&w.WarriorRank,
		&w.WarriorAvatar,
		&active,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("warrior not found")
	}

	if active {
		return nil, errors.New("warrior already active in battle")
	}

	return &w, nil
}

// GetBattleWarriors retrieves the warriors for a given battle from db
func (d *Database) GetBattleWarriors(BattleID string) []*BattleWarrior {
	var warriors = make([]*BattleWarrior, 0)
	rows, err := d.db.Query(
		`SELECT
			w.id, w.name, w.rank, w.avatar, bw.active
		FROM battles_warriors bw
		LEFT JOIN warriors w ON bw.warrior_id = w.id
		WHERE bw.battle_id = $1
		ORDER BY w.name`,
		BattleID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w BattleWarrior
			if err := rows.Scan(&w.WarriorID, &w.WarriorName, &w.WarriorRank, &w.WarriorAvatar, &w.Active); err != nil {
				log.Println(err)
			} else {
				warriors = append(warriors, &w)
			}
		}
	}

	return warriors
}

// GetBattleActiveWarriors retrieves the active warriors for a given battle from db
func (d *Database) GetBattleActiveWarriors(BattleID string) []*BattleWarrior {
	var warriors = make([]*BattleWarrior, 0)
	rows, err := d.db.Query(
		`SELECT
			w.id, w.name, w.rank, w.avatar, bw.active
		FROM battles_warriors bw
		LEFT JOIN warriors w ON bw.warrior_id = w.id
		WHERE bw.battle_id = $1 AND bw.active = true
		ORDER BY w.name`,
		BattleID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w BattleWarrior
			if err := rows.Scan(&w.WarriorID, &w.WarriorName, &w.WarriorRank, &w.WarriorAvatar, &w.Active); err != nil {
				log.Println(err)
			} else {
				warriors = append(warriors, &w)
			}
		}
	}

	return warriors
}

// AddWarriorToBattle adds a warrior by ID to the battle by ID
func (d *Database) AddWarriorToBattle(BattleID string, WarriorID string) ([]*BattleWarrior, error) {
	if _, err := d.db.Exec(
		`INSERT INTO battles_warriors (battle_id, warrior_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (battle_id, warrior_id) DO UPDATE SET active = true, abandoned = false`,
		BattleID,
		WarriorID,
	); err != nil {
		log.Println(err)
	}

	warriors := d.GetBattleWarriors(BattleID)

	return warriors, nil
}

// RetreatWarrior removes a warrior from the current battle by ID
func (d *Database) RetreatWarrior(BattleID string, WarriorID string) []*BattleWarrior {
	if _, err := d.db.Exec(
		`UPDATE battles_warriors SET active = false WHERE battle_id = $1 AND warrior_id = $2`, BattleID, WarriorID); err != nil {
		log.Println(err)
	}

	if _, err := d.db.Exec(
		`UPDATE warriors SET last_active = NOW() WHERE id = $1`, WarriorID); err != nil {
		log.Println(err)
	}

	warriors := d.GetBattleWarriors(BattleID)

	return warriors
}

// AbandonBattle removes a warrior from the current battle by ID and sets abandoned true
func (d *Database) AbandonBattle(BattleID string, WarriorID string) ([]*BattleWarrior, error) {
	if _, err := d.db.Exec(
		`UPDATE battles_warriors SET active = false, abandoned = true WHERE battle_id = $1 AND warrior_id = $2`, BattleID, WarriorID); err != nil {
		log.Println(err)
		return nil, err
	}

	if _, err := d.db.Exec(
		`UPDATE warriors SET last_active = NOW() WHERE id = $1`, WarriorID); err != nil {
		log.Println(err)
		return nil, err
	}

	warriors := d.GetBattleWarriors(BattleID)

	return warriors, nil
}

// SetBattleLeader sets the leaderId for the battle
func (d *Database) SetBattleLeader(BattleID string, warriorID string, LeaderID string) error {
	err := d.ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return errors.New("incorrect permissions")
	}

	// set battle VotingLocked
	if _, err := d.db.Exec(
		`call set_battle_leader($1, $2);`, BattleID, LeaderID); err != nil {
		log.Println(err)
		return errors.New("unable to promote leader")
	}

	return nil
}

// DeleteBattle removes all battle associations and the battle itself from DB by BattleID
func (d *Database) DeleteBattle(BattleID string, warriorID string) error {
	err := d.ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return errors.New("incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call delete_battle($1);`, BattleID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
