package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Battle aka arena
type Battle struct {
	BattleID     string     `json:"id"`
	LeaderID     string     `json:"leaderId"`
	BattleName   string     `json:"name"`
	Warriors     []*Warrior `json:"warriors"`
	Plans        []*Plan    `json:"plans"`
	VotingLocked bool       `json:"votingLocked"`
	ActivePlanID string     `json:"activePlanId"`
}

// Warrior aka user
type Warrior struct {
	WarriorID   string `json:"id"`
	WarriorName string `json:"name"`
}

// Vote structure
type Vote struct {
	WarriorID string `json:"warriorId"`
	VoteValue string `json:"vote"`
}

// Plan aka Story structure
type Plan struct {
	PlanID      string  `json:"id"`
	PlanName    string  `json:"name"`
	Votes       []*Vote `json:"votes"`
	Points      string  `json:"points"`
	PlanActive  bool    `json:"active"`
	PlanSkipped bool    `json:"skipped"`
}

// SetupDB runs db migrations, sets up a db connection pool
// and sets previously active warriors to false during startup
func SetupDB() {
	var (
		host     = GetEnv("DB_HOST", "db")
		port     = GetIntEnv("DB_PORT", 5432)
		user     = GetEnv("DB_USER", "thor")
		password = GetEnv("DB_PASS", "odinson")
		dbname   = GetEnv("DB_NAME", "thunderdome")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS battles (id UUID NOT NULL PRIMARY KEY, leader_id UUID, name VARCHAR(256), voting_locked BOOL DEFAULT true, active_plan_id UUID)"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS warriors (id UUID NOT NULL PRIMARY KEY, name VARCHAR(64))"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS plans (id UUID NOT NULL PRIMARY KEY, name VARCHAR(256), points VARCHAR(3) DEFAULT '', active BOOL DEFAULT false, battle_id UUID references battles(id) NOT NULL, votes JSONB DEFAULT '[]'::JSONB)"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS battles_warriors (battle_id UUID references battles NOT NULL, warrior_id UUID REFERENCES warriors NOT NULL, active BOOL DEFAULT false, PRIMARY KEY (battle_id, warrior_id))"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		"ALTER TABLE battles ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW()"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		"ALTER TABLE warriors ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW()"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		"ALTER TABLE plans ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW()"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		"ALTER TABLE warriors ADD COLUMN IF NOT EXISTS last_active TIMESTAMP DEFAULT NOW()"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		"ALTER TABLE plans ADD COLUMN IF NOT EXISTS updated_date TIMESTAMP DEFAULT NOW()"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		"ALTER TABLE battles ADD COLUMN IF NOT EXISTS updated_date TIMESTAMP DEFAULT NOW()"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		"ALTER TABLE plans ADD COLUMN IF NOT EXISTS skipped BOOL DEFAULT false"); err != nil {
		log.Fatal(err)
	}

	// on server start reset all warriors to active false for battles
	if _, err := db.Exec(
		`UPDATE battles_warriors SET active = false WHERE active = true`); err != nil {
		log.Println(err)
	}
}

//CreateBattle adds a new battle to the map
func CreateBattle(LeaderID string, BattleName string) (*Battle, error) {
	newID, _ := uuid.NewUUID()
	id := newID.String()

	var b = &Battle{
		BattleID:     id,
		LeaderID:     LeaderID,
		BattleName:   BattleName,
		Warriors:     make([]*Warrior, 0),
		Plans:        make([]*Plan, 0),
		VotingLocked: true,
		ActivePlanID: ""}

	e := db.QueryRow(`INSERT INTO battles (id, leader_id, name) VALUES ($1, $2, $3) RETURNING id`, id, LeaderID, BattleName).Scan(&b.BattleID)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Error Creating Battle")
	}

	return b, nil
}

// GetBattle gets a battle from the map by ID
func GetBattle(BattleID string) (*Battle, error) {
	var b = &Battle{
		BattleID:     BattleID,
		LeaderID:     "",
		BattleName:   "",
		Warriors:     make([]*Warrior, 0),
		Plans:        make([]*Plan, 0),
		VotingLocked: true,
		ActivePlanID: ""}

	// get battle
	var ActivePlanID sql.NullString
	e := db.QueryRow("SELECT id, name, leader_id, voting_locked, active_plan_id FROM battles WHERE id = $1", BattleID).Scan(&b.BattleID, &b.BattleName, &b.LeaderID, &b.VotingLocked, &ActivePlanID)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Not found")
	}

	b.ActivePlanID = ActivePlanID.String
	b.Warriors = GetActiveWarriors(BattleID)
	b.Plans = GetPlans(BattleID)

	return b, nil
}

// ConfirmLeader confirms the warrior is infact leader of the battle
func ConfirmLeader(BattleID string, warriorID string) error {
	var leaderID string
	e := db.QueryRow("SELECT leader_id FROM battles WHERE id = $1", BattleID).Scan(&leaderID)
	if e != nil {
		log.Println(e)
		return errors.New("Battle Not found")
	}

	if leaderID != warriorID {
		return errors.New("Not Leader")
	}

	return nil
}

// CreateWarrior adds a new warrior to the db
func CreateWarrior(WarriorName string) (*Warrior, error) {
	newID, _ := uuid.NewUUID()
	id := newID.String()

	var WarriorID string
	e := db.QueryRow(`INSERT INTO warriors (id, name) VALUES ($1, $2) RETURNING id`, id, WarriorName).Scan(&WarriorID)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Unable to create new warrior")
	}

	return &Warrior{WarriorID: WarriorID, WarriorName: WarriorName}, nil
}

// GetWarrior gets a warrior from db by ID
func GetWarrior(WarriorID string) (*Warrior, error) {
	var w Warrior

	e := db.QueryRow("SELECT id, name FROM warriors WHERE id = $1", WarriorID).Scan(&w.WarriorID, &w.WarriorName)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Warrior Not found")
	}

	return &w, nil
}

// GetBattleWarrior gets a warrior from db by ID and checks battle active status
func GetBattleWarrior(BattleID string, WarriorID string) (*Warrior, error) {
	var active bool
	var w Warrior

	e := db.QueryRow("SELECT id, name FROM warriors WHERE id = $1", WarriorID).Scan(&w.WarriorID, &w.WarriorName)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Warrior Not found")
	}

	err := db.QueryRow("SELECT active FROM battles_warriors WHERE battle_id = $1 AND warrior_id = $2", BattleID, WarriorID).Scan(&active)
	if err != nil {
		log.Println(err)
	}

	if active {
		return nil, errors.New("Warrior Already Active in Battle")
	}

	return &w, nil
}

// GetActiveWarriors retrieves the active warriors for a given battle from db
func GetActiveWarriors(BattleID string) []*Warrior {
	var warriors = make([]*Warrior, 0)
	rows, err := db.Query("SELECT warriors.id, warriors.name FROM battles_warriors LEFT JOIN warriors ON battles_warriors.warrior_id = warriors.id where battles_warriors.battle_id = $1 AND battles_warriors.active = true ORDER BY warriors.name", BattleID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w Warrior
			if err := rows.Scan(&w.WarriorID, &w.WarriorName); err != nil {
				log.Println(err)
			} else {
				warriors = append(warriors, &w)
			}
		}
	}

	return warriors
}

// AddWarriorToBattle adds a warrior by ID to the battle by ID
func AddWarriorToBattle(BattleID string, WarriorID string) ([]*Warrior, error) {
	if _, err := db.Exec(
		`INSERT INTO battles_warriors (battle_id, warrior_id, active) VALUES ($1, $2, true) ON CONFLICT (battle_id, warrior_id) DO UPDATE SET active = true`, BattleID, WarriorID); err != nil {
		log.Println(err)
	}

	warriors := GetActiveWarriors(BattleID)

	return warriors, nil
}

// RetreatWarrior removes a warrior from the current battle by ID
func RetreatWarrior(BattleID string, WarriorID string) []*Warrior {
	if _, err := db.Exec(
		`UPDATE battles_warriors SET active = false WHERE battle_id = $1 AND warrior_id = $2`, BattleID, WarriorID); err != nil {
		log.Println(err)
	}

	if _, err := db.Exec(
		`UPDATE warriors SET last_active = NOW() WHERE id = $1`, WarriorID); err != nil {
		log.Println(err)
	}

	warriors := GetActiveWarriors(BattleID)

	return warriors
}

// GetPlans retrieves plans for given battle from db
func GetPlans(BattleID string) []*Plan {
	var plans = make([]*Plan, 0)
	planRows, plansErr := db.Query("SELECT id, name, points, active, skipped, votes FROM plans WHERE battle_id = $1 ORDER BY created_date", BattleID)
	if plansErr == nil {
		defer planRows.Close()
		for planRows.Next() {
			var v string
			var p = &Plan{PlanID: "",
				PlanName:    "",
				Votes:       make([]*Vote, 0),
				Points:      "",
				PlanActive:  false,
				PlanSkipped: false,
			}
			if err := planRows.Scan(&p.PlanID, &p.PlanName, &p.Points, &p.PlanActive, &p.PlanSkipped, &v); err != nil {
				log.Println(err)
			} else {
				err = json.Unmarshal([]byte(v), &p.Votes)
				if err != nil {
					log.Println(err)
				}

				for i := range p.Votes {
					vote := p.Votes[i]
					if p.PlanActive {
						vote.VoteValue = ""
					}
				}

				plans = append(plans, p)
			}
		}
	}

	return plans
}

// CreatePlan adds a new plan to a battle
func CreatePlan(BattleID string, warriorID string, PlanName string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	newID, _ := uuid.NewUUID()
	id := newID.String()

	var PlanID string
	e := db.QueryRow(`INSERT INTO plans (id, battle_id, name) VALUES ($1, $2, $3) RETURNING id`, id, BattleID, PlanName).Scan(&PlanID)
	if e != nil {
		log.Println(e)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// ActivatePlanVoting sets the plan by ID to active, wipes any previous votes/points, and disables votingLock
func ActivatePlanVoting(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	// set current to false
	if _, err := db.Exec(`UPDATE plans SET updated_date = NOW(), active = false WHERE battle_id = $1`, BattleID); err != nil {
		log.Println(err)
	}

	// set PlanID to true
	if _, err := db.Exec(
		`UPDATE plans SET updated_date = NOW(), active = true, skipped = false, points = '', votes = '[]'::jsonb WHERE id = $1`, PlanID); err != nil {
		log.Println(err)
	}

	// set battle VotingLocked and ActivePlanID
	if _, err := db.Exec(
		`UPDATE battles SET updated_date = NOW(), voting_locked = false, active_plan_id = $1 WHERE id = $2`, PlanID, BattleID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// SetVote sets a warriors vote for the plan
func SetVote(BattleID string, WarriorID string, PlanID string, VoteValue string) []*Plan {
	// get plan
	var v string
	e := db.QueryRow("SELECT votes FROM plans WHERE id = $1", PlanID).Scan(&v)
	if e != nil {
		log.Println(e)
		// return nil, errors.New("Plan Not found")
	}
	var votes []*Vote
	err := json.Unmarshal([]byte(v), &votes)
	if err != nil {
		log.Println(err)
	}

	var voteIndex int
	var voteFound bool

	// find vote index
	for vi := range votes {
		if votes[vi].WarriorID == WarriorID {
			voteFound = true
			voteIndex = vi
			break
		}
	}

	if voteFound {
		votes[voteIndex].VoteValue = VoteValue
	} else {
		newVote := &Vote{WarriorID: WarriorID,
			VoteValue: VoteValue}

		votes = append(votes, newVote)
	}

	// update votes on Plan
	var votesJSON, _ = json.Marshal(votes)
	if _, err := db.Exec(
		`UPDATE plans SET updated_date = NOW(), votes = $1 WHERE id = $2`, string(votesJSON), PlanID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans
}

// RetractVote removes a warriors vote for the plan
func RetractVote(BattleID string, WarriorID string, PlanID string) []*Plan {
	// get plan
	var v string
	e := db.QueryRow("SELECT votes FROM plans WHERE id = $1", PlanID).Scan(&v)
	if e != nil {
		log.Println(e)
		// return nil, errors.New("Plan Not found")
	}
	var votes []*Vote
	err := json.Unmarshal([]byte(v), &votes)
	if err != nil {
		log.Println(err)
	}

	var voteIndex int
	var voteFound bool

	// find vote index
	for vi := range votes {
		if votes[vi].WarriorID == WarriorID {
			voteFound = true
			voteIndex = vi
			break
		}
	}

	if voteFound {
		votes = append(votes[:voteIndex], votes[voteIndex+1:]...) 
	}

	// update votes on Plan
	var votesJSON, _ = json.Marshal(votes)
	if _, err := db.Exec(
		`UPDATE plans SET updated_date = NOW(), votes = $1 WHERE id = $2`, string(votesJSON), PlanID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans
}

// EndPlanVoting sets plan to active: false
func EndPlanVoting(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	// set current to false
	if _, err := db.Exec(`UPDATE plans SET updated_date = NOW(), active = false WHERE battle_id = $1`, BattleID); err != nil {
		log.Println(err)
	}

	// set battle VotingLocked
	if _, err := db.Exec(
		`UPDATE battles SET updated_date = NOW(), voting_locked = true WHERE id = $1`, BattleID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// SkipPlan sets plan to active: false and unsets battle's activePlanId
func SkipPlan(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	// set current to false
	if _, err := db.Exec(`UPDATE plans SET updated_date = NOW(), active = false, skipped = true WHERE battle_id = $1`, BattleID); err != nil {
		log.Println(err)
	}

	// set battle VotingLocked and activePlanId to null
	if _, err := db.Exec(
		`UPDATE battles SET updated_date = NOW(), voting_locked = true, active_plan_id = null WHERE id = $1`, BattleID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// RevisePlanName updates the plan name by ID
func RevisePlanName(BattleID string, warriorID string, PlanID string, PlanName string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	// set PlanID to true
	if _, err := db.Exec(
		`UPDATE plans SET updated_date = NOW(), name = $1 WHERE id = $2`, PlanName, PlanID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// BurnPlan removes a plan from the current battle by ID
func BurnPlan(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	var isActivePlan bool

	// delete plan
	e := db.QueryRow("DELETE FROM plans WHERE id = $1 RETURNING active", PlanID).Scan(&isActivePlan)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Plan Not found")
	}

	if isActivePlan {
		if _, err := db.Exec(
			`UPDATE battles SET updated_date = NOW(), voting_locked = true, active_plan_id = null WHERE id = $1`, BattleID); err != nil {
			log.Println(err)
		}
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// FinalizePlan sets plan to active: false
func FinalizePlan(BattleID string, warriorID string, PlanID string, PlanPoints string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	// set PlanID to true
	if _, err := db.Exec(
		`UPDATE plans SET updated_date = NOW(), active = false, points = $1 WHERE id = $2`, PlanPoints, PlanID); err != nil {
		log.Println(err)
	}

	// set battle ActivePlanID
	if _, err := db.Exec(
		`UPDATE battles SET updated_date = NOW(), active_plan_id = null WHERE id = $1`, BattleID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// SetBattleLeader sets the leaderId for the battle
func SetBattleLeader(BattleID string, warriorID string, LeaderID string) (*Battle, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	// set battle VotingLocked
	if _, err := db.Exec(
		`UPDATE battles SET updated_date = NOW(), leader_id = $1 WHERE id = $2`, LeaderID, BattleID); err != nil {
		log.Println(err)
	}

	battle, err := GetBattle(BattleID)
	if err != nil {
		return nil, errors.New("Unable to promote leader")
	}

	return battle, nil
}

// DeleteBattle removes all battle associations and the battle itself from DB by BattleID
func DeleteBattle(BattleID string, warriorID string) error {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return errors.New("Incorrect permissions")
	}

	// delete plans associated with battle
	if _, err := db.Exec(
		`DELETE FROM plans WHERE battle_id = $1`, BattleID); err != nil {
		log.Println(err)
	}

	// delete battle_warriors associations
	if _, err := db.Exec(
		`DELETE FROM battles_warriors WHERE battle_id = $1`, BattleID); err != nil {
		log.Println(err)
	}

	// delete battle itself
	if _, err := db.Exec(
		`DELETE FROM battles WHERE id = $1`, BattleID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
