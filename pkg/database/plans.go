package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

// GetPlans retrieves plans for given battle from db
func (d *Database) GetPlans(BattleID string) []*Plan {
	var plans = make([]*Plan, 0)
	planRows, plansErr := d.db.Query(
		`SELECT
			id, name, type, reference_id, link, description, acceptance_criteria, points, active, skipped, votestart_time, voteend_time, votes
			FROM plans WHERE battle_id = $1 ORDER BY created_date
		`,
		BattleID,
	)
	if plansErr == nil {
		defer planRows.Close()
		for planRows.Next() {
			var v string
			var ReferenceID sql.NullString
			var Link sql.NullString
			var Description sql.NullString
			var AcceptanceCriteria sql.NullString
			var p = &Plan{PlanID: "",
				PlanName:           "",
				Type:               "",
				ReferenceID:        "",
				Link:               "",
				Description:        "",
				AcceptanceCriteria: "",
				Votes:              make([]*Vote, 0),
				Points:             "",
				PlanActive:         false,
				PlanSkipped:        false,
				VoteStartTime:      time.Now(),
				VoteEndTime:        time.Now(),
			}
			if err := planRows.Scan(
				&p.PlanID, &p.PlanName, &p.Type, &ReferenceID, &Link, &Description, &AcceptanceCriteria, &p.Points, &p.PlanActive, &p.PlanSkipped, &p.VoteStartTime, &p.VoteEndTime, &v,
			); err != nil {
				log.Println(err)
			} else {
				p.ReferenceID = ReferenceID.String
				p.Link = Link.String
				p.Description = Description.String
				p.AcceptanceCriteria = AcceptanceCriteria.String
				err = json.Unmarshal([]byte(v), &p.Votes)
				if err != nil {
					log.Println(err)
				}

				// don't send vote values to client, prevent sneaky devs from peaking at votes
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
func (d *Database) CreatePlan(BattleID string, warriorID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string) ([]*Plan, error) {
	err := d.ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("incorrect permissions")
	}

	// @TODO - refactor stored procedure to replace need for app generated uuid
	newID, _ := uuid.NewUUID()
	PlanID := newID.String()

	if _, err := d.db.Exec(
		`call create_plan($1, $2, $3, $4, $5, $6, $7, $8);`, BattleID, PlanID, PlanName, PlanType, ReferenceID, Link, Description, AcceptanceCriteria,
	); err != nil {
		log.Println(err)
	}

	plans := d.GetPlans(BattleID)

	return plans, nil
}

// ActivatePlanVoting sets the plan by ID to active, wipes any previous votes/points, and disables votingLock
func (d *Database) ActivatePlanVoting(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := d.ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call activate_plan_voting($1, $2);`, BattleID, PlanID,
	); err != nil {
		log.Println(err)
	}

	plans := d.GetPlans(BattleID)

	return plans, nil
}

// SetVote sets a warriors vote for the plan
func (d *Database) SetVote(BattleID string, WarriorID string, PlanID string, VoteValue string) (BattlePlans []*Plan, AllWarriorsVoted bool) {
	if _, err := d.db.Exec(
		`call set_warrior_vote($1, $2, $3);`, PlanID, WarriorID, VoteValue); err != nil {
		log.Println(err)
	}

	Plans := d.GetPlans(BattleID)
	ActiveWarriors := d.GetBattleActiveWarriors(BattleID)

	// determine if all active warriors have voted
	AllVoted := true
	for _, plan := range Plans {
		if plan.PlanID == PlanID {
			activePlanVoters := make(map[string]bool)

			for _, vote := range plan.Votes {
				var WarriorID string = vote.WarriorID
				activePlanVoters[WarriorID] = true
			}
			for _, war := range ActiveWarriors {
				_, warriorVoted := activePlanVoters[war.WarriorID]
				if warriorVoted == false {
					AllVoted = false
					break
				}
			}
			break
		}
	}

	return Plans, AllVoted
}

// RetractVote removes a warriors vote for the plan
func (d *Database) RetractVote(BattleID string, WarriorID string, PlanID string) []*Plan {
	if _, err := d.db.Exec(
		`call retract_warrior_vote($1, $2);`, PlanID, WarriorID); err != nil {
		log.Println(err)
	}

	plans := d.GetPlans(BattleID)

	return plans
}

// EndPlanVoting sets plan to active: false
func (d *Database) EndPlanVoting(BattleID string, warriorID string, PlanID string, AutoFinishVoting bool) ([]*Plan, error) {
	if AutoFinishVoting == false {
		err := d.ConfirmLeader(BattleID, warriorID)
		if err != nil {
			return nil, errors.New("incorrect permissions")
		}
	}

	if _, err := d.db.Exec(
		`call end_plan_voting($1, $2);`, BattleID, PlanID); err != nil {
		log.Println(err)
	}

	plans := d.GetPlans(BattleID)

	return plans, nil
}

// SkipPlan sets plan to active: false and unsets battle's activePlanId
func (d *Database) SkipPlan(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := d.ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call skip_plan_voting($1, $2);`, BattleID, PlanID); err != nil {
		log.Println(err)
	}

	plans := d.GetPlans(BattleID)

	return plans, nil
}

// RevisePlan updates the plan by ID
func (d *Database) RevisePlan(BattleID string, warriorID string, PlanID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string) ([]*Plan, error) {
	err := d.ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("incorrect permissions")
	}

	// set PlanID to true
	if _, err := d.db.Exec(
		`call revise_plan($1, $2, $3, $4, $5, $6, $7);`, PlanID, PlanName, PlanType, ReferenceID, Link, Description, AcceptanceCriteria); err != nil {
		log.Println(err)
	}

	plans := d.GetPlans(BattleID)

	return plans, nil
}

// BurnPlan removes a plan from the current battle by ID
func (d *Database) BurnPlan(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := d.ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call delete_plan($1, $2);`, BattleID, PlanID); err != nil {
		log.Println(err)
	}

	plans := d.GetPlans(BattleID)

	return plans, nil
}

// FinalizePlan sets plan to active: false
func (d *Database) FinalizePlan(BattleID string, warriorID string, PlanID string, PlanPoints string) ([]*Plan, error) {
	err := d.ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call finalize_plan($1, $2, $3);`, BattleID, PlanID, PlanPoints); err != nil {
		log.Println(err)
	}

	plans := d.GetPlans(BattleID)

	return plans, nil
}
