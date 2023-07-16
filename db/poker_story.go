package db

import (
	"database/sql"
	"encoding/json"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"
)

// GetPlans retrieves plans for given battle
func (d *PokerService) GetPlans(BattleID string, UserID string) []*thunderdome.Story {
	var plans = make([]*thunderdome.Story, 0)
	planRows, plansErr := d.DB.Query(
		`SELECT
			id, name, type, reference_id, link, description, acceptance_criteria, priority, points, active, skipped, votestart_time, voteend_time, votes
			FROM thunderdome.poker_story WHERE poker_id = $1 ORDER BY created_date
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
			var p = &thunderdome.Story{
				Votes:   make([]*thunderdome.Vote, 0),
				Active:  false,
				Skipped: false,
			}
			if err := planRows.Scan(
				&p.Id, &p.Name, &p.Type, &ReferenceID, &Link, &Description, &AcceptanceCriteria, &p.Priority, &p.Points, &p.Active, &p.Skipped, &p.VoteStartTime, &p.VoteEndTime, &v,
			); err != nil {
				d.Logger.Error("get poker stories query error", zap.Error(err))
			} else {
				p.ReferenceId = ReferenceID.String
				p.Link = Link.String
				p.Description = Description.String
				p.AcceptanceCriteria = AcceptanceCriteria.String
				err = json.Unmarshal([]byte(v), &p.Votes)
				if err != nil {
					d.Logger.Error("get poker stories query scan error", zap.Error(err))
				}

				// don't send others vote values to client, prevent sneaky devs from peaking at votes
				for i := range p.Votes {
					if p.Active && p.Votes[i].UserId != UserID {
						p.Votes[i].VoteValue = ""
					}
				}

				plans = append(plans, p)
			}
		}
	}

	return plans
}

// CreatePlan adds a new plan to a battle
func (d *PokerService) CreatePlan(BattleID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*thunderdome.Story, error) {
	SanitizedDescription := d.HTMLSanitizerPolicy.Sanitize(Description)
	SanitizedAcceptanceCriteria := d.HTMLSanitizerPolicy.Sanitize(AcceptanceCriteria)
	// default priority should be 99 for sort order purposes
	if Priority == 0 {
		Priority = 99
	}
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.poker_story (poker_id, name, type, reference_id, link, description, acceptance_criteria, priority)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		BattleID, PlanName, PlanType, ReferenceID, Link, SanitizedDescription, SanitizedAcceptanceCriteria, Priority,
	); err != nil {
		d.Logger.Error("error creating poker story", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// ActivatePlanVoting sets the plan by ID to active, wipes any previous votes/points, and disables votingLock
func (d *PokerService) ActivatePlanVoting(BattleID string, PlanID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_story_activate($1, $2);`, BattleID, PlanID,
	); err != nil {
		d.Logger.Error("CALL thunderdome.poker_story_activate error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// SetVote sets a users vote for the plan
func (d *PokerService) SetVote(BattleID string, UserID string, PlanID string, VoteValue string) (BattlePlans []*thunderdome.Story, AllUsersVoted bool) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_user_vote_set($1, $2, $3);`, PlanID, UserID, VoteValue); err != nil {
		d.Logger.Error("CALL thunderdome.poker_user_vote_set error", zap.Error(err))
	}

	Plans := d.GetPlans(BattleID, "")
	ActiveUsers := d.GetBattleActiveUsers(BattleID)

	// determine if all active users have voted
	AllVoted := true
	for _, plan := range Plans {
		if plan.Id == PlanID {
			activePlanVoters := make(map[string]bool)

			for _, vote := range plan.Votes {
				var UserID string = vote.UserId
				activePlanVoters[UserID] = true
			}
			for _, war := range ActiveUsers {
				if _, UserVoted := activePlanVoters[war.Id]; !UserVoted && !war.Spectator {
					AllVoted = false
					break
				}
			}
			break
		}
	}

	return Plans, AllVoted
}

// RetractVote removes a users vote for the plan
func (d *PokerService) RetractVote(BattleID string, UserID string, PlanID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_user_vote_retract($1, $2);`, PlanID, UserID); err != nil {
		d.Logger.Error("CALL thunderdome.poker_vote_retract error", zap.Error(err))
		return nil, err
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// EndPlanVoting sets plan to active: false
func (d *PokerService) EndPlanVoting(BattleID string, PlanID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_plan_voting_stop($1, $2);`, BattleID, PlanID); err != nil {
		d.Logger.Error("CALL thunderdome.poker_plan_voting_stop error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// SkipPlan sets plan to active: false and unsets battle's activePlanId
func (d *PokerService) SkipPlan(BattleID string, PlanID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_vote_skip($1, $2);`, BattleID, PlanID); err != nil {
		d.Logger.Error("CALL thunderdome.poker_vote_skip error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// RevisePlan updates the plan by ID
func (d *PokerService) RevisePlan(BattleID string, PlanID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*thunderdome.Story, error) {
	SanitizedDescription := d.HTMLSanitizerPolicy.Sanitize(Description)
	SanitizedAcceptanceCriteria := d.HTMLSanitizerPolicy.Sanitize(AcceptanceCriteria)
	// default priority should be 99 for sort order purposes
	if Priority == 0 {
		Priority = 99
	}
	// set PlanID to true
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.poker_story
    SET
        updated_date = NOW(),
        name = $2,
        type = $3,
        reference_id = $4,
        link = $5,
        description = $6,
        acceptance_criteria = $7,
        priority = $8
    WHERE id = $1;`,
		PlanID, PlanName, PlanType, ReferenceID, Link, SanitizedDescription, SanitizedAcceptanceCriteria, Priority); err != nil {
		d.Logger.Error("error getting poker story", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// BurnPlan removes a plan from the current battle by ID
func (d *PokerService) BurnPlan(pokerID string, storyID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_story_delete($1, $2);`, pokerID, storyID); err != nil {
		d.Logger.Error("CALL thunderdome.poker_story_delete error", zap.Error(err))
	}

	plans := d.GetPlans(pokerID, "")

	return plans, nil
}

// FinalizePlan sets plan to active: false
func (d *PokerService) FinalizePlan(BattleID string, PlanID string, PlanPoints string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_story_finalize($1, $2, $3);`, BattleID, PlanID, PlanPoints); err != nil {
		d.Logger.Error("CALL thunderdome.poker_story_finalize error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}
