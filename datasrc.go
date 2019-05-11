package main

import (
	"errors"

	"github.com/google/uuid"
)

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
	Active   bool    `json:"active"`
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
		BattleID:     id,
		LeaderID:     LeaderID,
		BattleName:   BattleName,
		Warriors:     make([]*Warrior, 0),
		Plans:        make([]*Plan, 0),
		VotingLocked: true,
		ActivePlanID: ""}

	return Battles[id]
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
	newID, _ := uuid.NewUUID()
	id := newID.String()

	Warriors[id] = &Warrior{WarriorID: id, WarriorName: WarriorName}

	return Warriors[id]
}

// GetWarrior gets a warrior from the map by ID
func GetWarrior(WarriorID string) (*Warrior, error) {
	if warrior, ok := Warriors[WarriorID]; ok {
		return warrior, nil
	}

	return nil, errors.New("Not found")
}

// AddWarriorToBattle adds a warrior by ID to the battle by ID
func AddWarriorToBattle(BattleID string, WarriorID string) []*Warrior {
	Battles[BattleID].Warriors = append(Battles[BattleID].Warriors, Warriors[WarriorID])

	return Battles[BattleID].Warriors
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
	newID, _ := uuid.NewUUID()
	id := newID.String()

	newPlan := &Plan{PlanID: id,
		PlanName: PlanName,
		Votes:    make([]*Vote, 0),
		Points:   "",
		Active:   false}

	Battles[BattleID].Plans = append(Battles[BattleID].Plans, newPlan)

	return Battles[BattleID].Plans
}

// ActivatePlanVoting sets the plan by ID to active and disables votingLock
func ActivatePlanVoting(BattleID string, PlanID string) []*Plan {
	var planIndex int
	var lastActivePlanIndex int
	var hasLastActivePlan = false

	for i := range Battles[BattleID].Plans {
		if Battles[BattleID].Plans[i].Active {
			hasLastActivePlan = true
			lastActivePlanIndex = i
		}
		if Battles[BattleID].Plans[i].PlanID == PlanID {
			planIndex = i
		}
	}

	// disable last active plan if one existed
	if hasLastActivePlan {
		Battles[BattleID].Plans[lastActivePlanIndex].Active = false
	}

	Battles[BattleID].Plans[planIndex].Active = true
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

	Battles[BattleID].Plans[planIndex].Active = false
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
	for i := range Battles[BattleID].Plans {
		if Battles[BattleID].Plans[i].PlanID == PlanID {
			planIndex = i
			break
		}
	}

	Battles[BattleID].Plans = append(Battles[BattleID].Plans[:planIndex], Battles[BattleID].Plans[planIndex+1:]...)

	return Battles[BattleID].Plans
}