package main

import (
	"database/sql"
	"errors"
    "log"

	_ "github.com/lib/pq"
)

var sqlConnectionHost string = "db-1"
var sqlConnectionPort string = "26257"
var sqlConnectionUser string = "thor"
var sqlConnectionDB string = "thunderdome"
var sqlConnectionPath string = "postgresql://thor@db-1:26257/thunderdome?sslmode=disable"

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
	PlanID   string  `json:"id"`
	PlanName string  `json:"name"`
	Votes    []*Vote `json:"votes"`
	Points   string  `json:"points"`
	PlanActive   bool    `json:"active"`
}

// Battles stores all battles in memory
var Battles = make(map[string]*Battle)

func SetupDB() {
	db, err := sql.Open("postgres", "postgresql://thor@db-1:26257/thunderdome?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }

    if _, err := db.Exec(
        "CREATE TABLE IF NOT EXISTS battles (id UUID DEFAULT uuid_v4()::UUID PRIMARY KEY, leader_id UUID, name STRING(256), voting_locked BOOL DEFAULT true, active_plan_id UUID)"); err != nil {
        log.Fatal(err)
	}
	
	if _, err := db.Exec(
        "CREATE TABLE IF NOT EXISTS warriors (id UUID DEFAULT uuid_v4()::UUID PRIMARY KEY, name STRING(64))"); err != nil {
        log.Fatal(err)
	}
	
	if _, err := db.Exec(
        "CREATE TABLE IF NOT EXISTS plans (id UUID DEFAULT uuid_v4()::UUID PRIMARY KEY, name STRING(256), points STRING(3) DEFAULT '', active BOOL DEFAULT false, battle_id UUID)"); err != nil {
        log.Fatal(err)
	}
	
	if _, err := db.Exec(
        "CREATE TABLE IF NOT EXISTS battles_warriors (battle_id UUID references battles, warrior_id UUID REFERENCES warriors, active BOOL DEFAULT false, PRIMARY KEY (battle_id, warrior_id))"); err != nil {
        log.Fatal(err)
	}
}

//CreateBattle adds a new battle to the map
func CreateBattle(LeaderID string, BattleName string) (*Battle, error) {
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}

	var b = &Battle{
		BattleID:     "",
		LeaderID:     LeaderID,
		BattleName:   BattleName,
		Warriors:     make([]*Warrior, 0),
		Plans:        make([]*Plan, 0),
		VotingLocked: true,
		ActivePlanID: ""}

	e := db.QueryRow(`INSERT INTO battles (leader_id, name) VALUES ($1, $2) RETURNING id`, LeaderID, BattleName).Scan(&b.BattleID)
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
	
    db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
    }

	// get battle
	// @TODO - solve nil active_plan_id issue
	e := db.QueryRow("SELECT id, name, leader_id, voting_locked FROM battles WHERE id = '"+BattleID+"'::UUID").Scan(&b.BattleID, &b.BattleName, &b.LeaderID, &b.VotingLocked)
    if e != nil {
		log.Println(e)
		return nil, errors.New("Not found")
	}

	b.Warriors = GetActiveWarriors(BattleID)
	b.Plans = GetPlans(BattleID)
	
	return b, nil
}

// CreateWarrior adds a new warrior to the db
func CreateWarrior(WarriorName string) *Warrior {
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}

	var WarriorID string
	e := db.QueryRow(`INSERT INTO warriors (name) VALUES ($1) RETURNING id`, WarriorName).Scan(&WarriorID)
	if e != nil {
        log.Println(e)
	}

	return &Warrior{WarriorID: WarriorID, WarriorName: WarriorName}
}

// GetWarrior gets a warrior from db by ID
func GetWarrior(WarriorID string) (*Warrior, error) {
	var w Warrior
	
    db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
    }

	e := db.QueryRow("SELECT id, name FROM warriors WHERE id = '"+WarriorID+"'::UUID").Scan(&w.WarriorID, &w.WarriorName)
    if e != nil {
		log.Println(e)
		return nil, errors.New("Not found")
	}

	return &w, nil
}

// GetActiveWarriors retrieves the active warriors for a given battle from db
func GetActiveWarriors(BattleID string) []*Warrior {
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}
	
	var warriors = make([]*Warrior, 0)
	rows, err := db.Query("SELECT warriors.id, warriors.name FROM battles_warriors LEFT JOIN warriors ON battles_warriors.warrior_id = warriors.id where battles_warriors.battle_id = $1 AND battles_warriors.active = true", BattleID)
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
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}

	if _, err := db.Exec(
        `UPSERT INTO battles_warriors (battle_id, warrior_id, active) VALUES ($1, $2, true)`, BattleID, WarriorID); err != nil {
        log.Println(err)
	}

	warriors := GetActiveWarriors(BattleID)

	return warriors, nil
}

// RetreatWarrior removes a warrior from the current battle by ID
func RetreatWarrior(BattleID string, WarriorID string) []*Warrior {
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}

	if _, err := db.Exec(
        `UPSERT INTO battles_warriors (battle_id, warrior_id, active) VALUES ($1, $2, false)`, BattleID, WarriorID); err != nil {
        log.Println(err)
	}

	warriors := GetActiveWarriors(BattleID)

	return warriors
}

// GetPlans retrieves plans for given battle from db
func GetPlans(BattleID string) []*Plan {
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}
	
	var plans = make([]*Plan, 0)
	planRows, plansErr := db.Query("SELECT id, name, points, active FROM plans WHERE battle_id = $1", BattleID)
	if plansErr == nil {
        defer planRows.Close()
		for planRows.Next() {
			var p = &Plan{PlanID: "",
				PlanName: "",
				Votes: make([]*Vote,0),
				Points: "",
				PlanActive: false}
			if err := planRows.Scan(&p.PlanID, &p.PlanName, &p.Points, &p.PlanActive); err != nil {
				log.Println(err)
			} else {
				plans = append(plans, p)
			}
		}
	}

	return plans
}

// CreatePlan adds a new plan to a battle
func CreatePlan(BattleID string, PlanName string) []*Plan {
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}

	var PlanID string
	e := db.QueryRow(`INSERT INTO plans (battle_id, name) VALUES ($1, $2) RETURNING id`, BattleID, PlanName).Scan(&PlanID)
	if e != nil {
        log.Println(e)
	}

	plans := GetPlans(BattleID)

	return plans
}

// ActivatePlanVoting sets the plan by ID to active and disables votingLock
func ActivatePlanVoting(BattleID string, PlanID string) []*Plan {
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}

	// set current to false
	if _, err := db.Exec(`UPDATE plans SET active = false WHERE battle_id = $1`, BattleID); err != nil {
        log.Println(err)
	}

	// set PlanID to true
	if _, err := db.Exec(
        `UPDATE plans SET active = true WHERE id = $1`, PlanID); err != nil {
        log.Println(err)
	}

	// set battle VotingLocked and ActivePlanID
	if _, err := db.Exec(
        `UPDATE battles SET voting_locked = false, active_plan_id = $1 WHERE id = $2`, PlanID, BattleID); err != nil {
        log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans
}

// SetVote sets a warriors vote for the plan
func SetVote(BattleID string, WarriorID string, PlanID string, VoteValue string) []*Plan {
	var planIndex int
	var voteIndex int
	var voteFound bool

	// find plan index
	for pi := range Battles[BattleID].Plans {
		if Battles[BattleID].Plans[pi].PlanID == PlanID {
			planIndex = pi

			// find vote index
			for vi := range Battles[BattleID].Plans[planIndex].Votes {
				if Battles[BattleID].Plans[pi].Votes[vi].WarriorID == WarriorID {
					voteFound = true
					voteIndex = vi
					break
				}
			}
			break
		}
	}

	if voteFound {
		Battles[BattleID].Plans[planIndex].Votes[voteIndex].VoteValue = VoteValue
	} else {
		newVote := &Vote{WarriorID: WarriorID,
			VoteValue: VoteValue}

		Battles[BattleID].Plans[planIndex].Votes = append(Battles[BattleID].Plans[planIndex].Votes, newVote)
	}

	return Battles[BattleID].Plans
}

// EndPlanVoting sets plan to active: false
func EndPlanVoting(BattleID string, PlanID string) []*Plan {
	var planIndex int

	for i := range Battles[BattleID].Plans {
		if Battles[BattleID].Plans[i].PlanID == PlanID {
			planIndex = i
			break
		}
	}

	Battles[BattleID].Plans[planIndex].PlanActive = false
	Battles[BattleID].VotingLocked = true

	return Battles[BattleID].Plans
}

// RevisePlanName updates the plan name by ID
func RevisePlanName(BattleID string, PlanID string, PlanName string) []*Plan {
	var planIndex int

	for i := range Battles[BattleID].Plans {
		if Battles[BattleID].Plans[i].PlanID == PlanID {
			planIndex = i
			break
		}
	}

	Battles[BattleID].Plans[planIndex].PlanName = PlanName

	return Battles[BattleID].Plans
}

// BurnPlan removes a plan from the current battle by ID
func BurnPlan(BattleID string, PlanID string) []*Plan {
	var planIndex int
	var isActivePlan bool

	for i := range Battles[BattleID].Plans {
		if Battles[BattleID].Plans[i].PlanID == PlanID {
			planIndex = i
			if Battles[BattleID].ActivePlanID == PlanID {
				isActivePlan = true
			}
			break
		}
	}

	if (isActivePlan) {
		Battles[BattleID].ActivePlanID = ""
		Battles[BattleID].VotingLocked = true
	}

	Battles[BattleID].Plans = append(Battles[BattleID].Plans[:planIndex], Battles[BattleID].Plans[planIndex+1:]...)

	return Battles[BattleID].Plans
}

// FinalizePlan sets plan to active: false
func FinalizePlan(BattleID string, PlanID string, PlanPoints string) []*Plan {
	var planIndex int

	for i := range Battles[BattleID].Plans {
		if Battles[BattleID].Plans[i].PlanID == PlanID {
			planIndex = i
			break
		}
	}

	Battles[BattleID].Plans[planIndex].PlanActive = false
	Battles[BattleID].Plans[planIndex].Points = PlanPoints
	Battles[BattleID].ActivePlanID = ""

	return Battles[BattleID].Plans
}