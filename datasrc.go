package main

import (
	"database/sql"
	"errors"
	// "fmt"
    "log"

	// "github.com/google/uuid"
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
	// Connect to the "bank" database.
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
        "CREATE TABLE IF NOT EXISTS plans (id UUID DEFAULT uuid_v4()::UUID PRIMARY KEY, name STRING(256), points STRING(3), active BOOL DEFAULT false)"); err != nil {
        log.Fatal(err)
    }

	// if _, err := db.Exec(
    //     "INSERT INTO battles (leaderId, name) VALUES ('abf03a57-71e7-11e9-82e5-0242ac110009', 'Asgard')"); err != nil {
    //     log.Fatal(err)
	// }
}

//CreateBattle adds a new battle to the map
func CreateBattle(LeaderID string, BattleName string) *Battle {
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}

	var BattleID string
	e := db.QueryRow(`INSERT INTO battles (leader_id, name) VALUES ($1, $2) RETURNING id`, LeaderID, BattleName).Scan(&BattleID)
	if e != nil {
        log.Fatal(e)
	}

	Battles[BattleID] = &Battle{
		BattleID:     BattleID,
		LeaderID:     LeaderID,
		BattleName:   BattleName,
		Warriors:     make([]*Warrior, 0),
		Plans:        make([]*Plan, 0),
		VotingLocked: true,
		ActivePlanID: ""}

	return Battles[BattleID]
}

// GetBattle gets a battle from the map by ID
func GetBattle(BattleID string) (*Battle, error) {    
	if battle, ok := Battles[BattleID]; ok {
		return battle, nil
	}
	return nil, errors.New("Not found")
}

// CreateWarrior adds a new warrior to the map
func CreateWarrior(WarriorName string) *Warrior {
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}

	var WarriorID string
	e := db.QueryRow(`INSERT INTO warriors (name) VALUES ($1) RETURNING id`, WarriorName).Scan(&WarriorID)
	if e != nil {
        log.Fatal(e)
	}

	return &Warrior{WarriorID: WarriorID, WarriorName: WarriorName}
}

// GetWarrior gets a warrior from the map by ID
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

// AddWarriorToBattle adds a warrior by ID to the battle by ID
func AddWarriorToBattle(BattleID string, WarriorID string) ([]*Warrior, error) {
	w, err := GetWarrior(WarriorID)
	if (err != nil) {
		log.Println(err)
		return nil, errors.New("Not found")
	}
	Battles[BattleID].Warriors = append(Battles[BattleID].Warriors, w)

	return Battles[BattleID].Warriors, nil
}

// RetreatWarrior removes a warrior from the current battle by ID
func RetreatWarrior(BattleID string, WarriorID string) []*Warrior {
	var warriorIndex int
	for i := range Battles[BattleID].Warriors {
		if Battles[BattleID].Warriors[i].WarriorID == WarriorID {
			warriorIndex = i
			break
		}
	}

	Battles[BattleID].Warriors = append(Battles[BattleID].Warriors[:warriorIndex], Battles[BattleID].Warriors[warriorIndex+1:]...)

	return Battles[BattleID].Warriors
}

// CreatePlan adds a new plan to a battle
func CreatePlan(BattleID string, PlanName string) []*Plan {
	db, err := sql.Open("postgres", sqlConnectionPath)
    if err != nil {
        log.Println("error connecting to the database: ", err)
	}

	var PlanID string
	e := db.QueryRow(`INSERT INTO plans (name) VALUES ($1) RETURNING id`, PlanName).Scan(&PlanID)
	if e != nil {
        log.Fatal(e)
	}

	newPlan := &Plan{PlanID: PlanID,
		PlanName: PlanName,
		Votes:    make([]*Vote, 0),
		Points:   "",
		PlanActive:   false}

	Battles[BattleID].Plans = append(Battles[BattleID].Plans, newPlan)

	return Battles[BattleID].Plans
}

// ActivatePlanVoting sets the plan by ID to active and disables votingLock
func ActivatePlanVoting(BattleID string, PlanID string) []*Plan {
	var planIndex int
	var lastActivePlanIndex int
	var hasLastActivePlan = false

	for i := range Battles[BattleID].Plans {
		if Battles[BattleID].Plans[i].PlanActive {
			hasLastActivePlan = true
			lastActivePlanIndex = i
		}
		if Battles[BattleID].Plans[i].PlanID == PlanID {
			planIndex = i
		}
	}

	// disable last active plan if one existed
	if hasLastActivePlan {
		Battles[BattleID].Plans[lastActivePlanIndex].PlanActive = false
	}

	Battles[BattleID].Plans[planIndex].PlanActive = true
	Battles[BattleID].VotingLocked = false
	Battles[BattleID].ActivePlanID = PlanID

	return Battles[BattleID].Plans
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