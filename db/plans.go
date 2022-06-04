package db

import (
	"database/sql"
	"encoding/json"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// GetPlans retrieves plans for given battle
func (d *Database) GetPlans(BattleID string, UserID string) []*model.Plan {
	var plans = make([]*model.Plan, 0)
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
			var p = &model.Plan{
				Votes:   make([]*model.Vote, 0),
				Active:  false,
				Skipped: false,
			}
			if err := planRows.Scan(
				&p.Id, &p.Name, &p.Type, &ReferenceID, &Link, &Description, &AcceptanceCriteria, &p.Points, &p.Active, &p.Skipped, &p.VoteStartTime, &p.VoteEndTime, &v,
			); err != nil {
				d.logger.Error("get battle plans query error", zap.Error(err))
			} else {
				p.ReferenceId = ReferenceID.String
				p.Link = Link.String
				p.Description = Description.String
				p.AcceptanceCriteria = AcceptanceCriteria.String
				err = json.Unmarshal([]byte(v), &p.Votes)
				if err != nil {
					d.logger.Error("get battle plans query scan error", zap.Error(err))
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
func (d *Database) CreatePlan(BattleID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string) ([]*model.Plan, error) {
	SanitizedDescription := d.htmlSanitizerPolicy.Sanitize(Description)
	SanitizedAcceptanceCriteria := d.htmlSanitizerPolicy.Sanitize(AcceptanceCriteria)
	if _, err := d.db.Exec(
		`call create_plan($1, $2, $3, $4, $5, $6, $7);`, BattleID, PlanName, PlanType, ReferenceID, Link, SanitizedDescription, SanitizedAcceptanceCriteria,
	); err != nil {
		d.logger.Error("call create_plan error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// ActivatePlanVoting sets the plan by ID to active, wipes any previous votes/points, and disables votingLock
func (d *Database) ActivatePlanVoting(BattleID string, PlanID string) ([]*model.Plan, error) {
	if _, err := d.db.Exec(
		`call activate_plan_voting($1, $2);`, BattleID, PlanID,
	); err != nil {
		d.logger.Error("call activate_plan_voting error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// SetVote sets a users vote for the plan
func (d *Database) SetVote(BattleID string, UserID string, PlanID string, VoteValue string) (BattlePlans []*model.Plan, AllUsersVoted bool) {
	if _, err := d.db.Exec(
		`call set_user_vote($1, $2, $3);`, PlanID, UserID, VoteValue); err != nil {
		d.logger.Error("call set_user_vote error", zap.Error(err))
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
func (d *Database) RetractVote(BattleID string, UserID string, PlanID string) ([]*model.Plan, error) {
	if _, err := d.db.Exec(
		`call retract_user_vote($1, $2);`, PlanID, UserID); err != nil {
		d.logger.Error("call retract_user_vote error", zap.Error(err))
		return nil, err
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// EndPlanVoting sets plan to active: false
func (d *Database) EndPlanVoting(BattleID string, PlanID string) ([]*model.Plan, error) {
	if _, err := d.db.Exec(
		`call end_plan_voting($1, $2);`, BattleID, PlanID); err != nil {
		d.logger.Error("call end_plan_voting error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// SkipPlan sets plan to active: false and unsets battle's activePlanId
func (d *Database) SkipPlan(BattleID string, PlanID string) ([]*model.Plan, error) {
	if _, err := d.db.Exec(
		`call skip_plan_voting($1, $2);`, BattleID, PlanID); err != nil {
		d.logger.Error("call skip_plan_voting error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// RevisePlan updates the plan by ID
func (d *Database) RevisePlan(BattleID string, PlanID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string) ([]*model.Plan, error) {
	SanitizedDescription := d.htmlSanitizerPolicy.Sanitize(Description)
	SanitizedAcceptanceCriteria := d.htmlSanitizerPolicy.Sanitize(AcceptanceCriteria)
	// set PlanID to true
	if _, err := d.db.Exec(
		`call revise_plan($1, $2, $3, $4, $5, $6, $7);`, PlanID, PlanName, PlanType, ReferenceID, Link, SanitizedDescription, SanitizedAcceptanceCriteria); err != nil {
		d.logger.Error("call revise_plan error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// BurnPlan removes a plan from the current battle by ID
func (d *Database) BurnPlan(BattleID string, PlanID string) ([]*model.Plan, error) {
	if _, err := d.db.Exec(
		`call delete_plan($1, $2);`, BattleID, PlanID); err != nil {
		d.logger.Error("call delete_plan error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}

// FinalizePlan sets plan to active: false
func (d *Database) FinalizePlan(BattleID string, PlanID string, PlanPoints string) ([]*model.Plan, error) {
	if _, err := d.db.Exec(
		`call finalize_plan($1, $2, $3);`, BattleID, PlanID, PlanPoints); err != nil {
		d.logger.Error("call finalize_plan error", zap.Error(err))
	}

	plans := d.GetPlans(BattleID, "")

	return plans, nil
}
