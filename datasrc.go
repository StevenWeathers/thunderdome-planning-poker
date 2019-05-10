package main

import (
	"errors"

	"github.com/google/uuid"
)

// Battle aka arena
type Battle struct {
	BattleID   string     `json:"id"`
	LeaderID   string     `json:"leaderId"`
	BattleName string     `json:"name"`
	Warriors   []*Warrior `json:"warriors"`
	Plans	   []*Plan	  `json:"plans"`
	VotingLocked bool	  `json:"votingLocked"`
	ActivePlanID string	  `json:"activePlanId"`
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
	PlanID 		string 		`json:"id"`
	PlanName 	string 		`json:"name"`
	Votes   	[]*Vote 	`json:"votes"`
	Points		string 		`json:"points"`
	Active 		bool 		`json:"active"`
}

// Warriors stores all warriors in memory
var Warriors = make(map[string]*Warrior)

// Battles stores all battles in memory
var Battles = make(map[string]*Battle)

//CreateBattle adds a new battle to the map
func CreateBattle(LeaderID string, BattleName string) *Battle {
	newID, _ := uuid.NewUUID()
	id := newID.String()

	Battles[id] = &Battle{
		BattleID:   id,
		LeaderID:   LeaderID,
		BattleName: BattleName,
		Warriors:   make([]*Warrior, 0),
		Plans:      make([]*Plan, 0),
		VotingLocked: true,
		ActivePlanID: ""}

	return Battles[id]
}

// GetBattle gets a battle from the map by ID
func GetBattle(BattleID string) (*Battle, error) {
	if battle, ok := Battles[BattleID]; ok {
		return battle, nil
	} else {
		return nil, errors.New("Not found")
	}
}

// CreateWarrior adds a new warrior to the map
func CreateWarrior(WarriorName string) *Warrior {
	newID, _ := uuid.NewUUID()
	id := newID.String()

	Warriors[id] = &Warrior{WarriorID: id, WarriorName: WarriorName}

	return Warriors[id]
}

// GetWarrior gets a warrior from the map by ID
func GetWarrior(WarriorID string) (*Warrior, error) {
	if warrior, ok := Warriors[WarriorID]; ok {
		return warrior, nil
	} else {
		return nil, errors.New("Not found")
	}
}

func AddWarriorToBattle(BattleID string, WarriorID string) {
	Battles[BattleID].Warriors = append(Battles[BattleID].Warriors, Warriors[WarriorID])
}

// ReatreatWarrior removes a warrior from the current battle by ID
func RetreatWarrior(BattleID string, WarriorID string) {
	var warriorIndex int
	for i := range Battles[BattleID].Warriors {
		if Battles[BattleID].Warriors[i].WarriorID == WarriorID {
			warriorIndex = i
			break
		}
	}

	Battles[BattleID].Warriors = append(Battles[BattleID].Warriors[:warriorIndex], Battles[BattleID].Warriors[warriorIndex+1:]...)
}

// CreatePlan adds a new plan to a battle
func CreatePlan(BattleID string, PlanName string) *Battle {
	newID, _ := uuid.NewUUID()
	id := newID.String()

	newPlan := &Plan{PlanID: id,
		PlanName: PlanName,
		Votes: make([]*Vote, 0),
		Points: "",
		Active: false} 

	Battles[BattleID].Plans = append(Battles[BattleID].Plans, newPlan)

	return Battles[BattleID]
}

// ActivatePlanVoting sets the plan by ID to active and disables votingLock
// needs to disable other active plans...
func ActivatePlanVoting(BattleID string, PlanID string) *Battle {
	var planIndex int

	for i := range Battles[BattleID].Plans {
		if Battles[BattleID].Plans[i].PlanID == PlanID {
			planIndex = i
			break
		}
	}

	Battles[BattleID].Plans[planIndex].Active = true
	Battles[BattleID].VotingLocked = false
	Battles[BattleID].ActivePlanID = PlanID

	return Battles[BattleID]
}

// SetVote sets a warriors vote for the plan
func SetVote(BattleID string, WarriorID string, PlanID string, VoteValue string) *Battle {
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

	return Battles[BattleID]
}

// EndPlanVoting sets plan to active: false
func EndPlanVoting(BattleID string, PlanID string) *Battle {
	var planIndex int

	for i := range Battles[BattleID].Plans {
		if Battles[BattleID].Plans[i].PlanID == PlanID {
			planIndex = i
			break
		}
	}

	Battles[BattleID].Plans[planIndex].Active = false
	Battles[BattleID].VotingLocked = true

	return Battles[BattleID]
}