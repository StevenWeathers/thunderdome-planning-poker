package poker

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"
)

// GetStories retrieves stories for given poker game
func (d *Service) GetStories(PokerID string, UserID string) []*thunderdome.Story {
	var plans = make([]*thunderdome.Story, 0)
	planRows, plansErr := d.DB.Query(
		`SELECT
			id, name, type, reference_id, link, description, acceptance_criteria, priority, 
			points, active, skipped, votestart_time, voteend_time, votes, 
			row_number() OVER (ORDER BY position ASC) as position
			FROM thunderdome.poker_story WHERE poker_id = $1 ORDER BY position
		`,
		PokerID,
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
				&p.Id, &p.Name, &p.Type, &ReferenceID, &Link, &Description, &AcceptanceCriteria, &p.Priority,
				&p.Points, &p.Active, &p.Skipped, &p.VoteStartTime, &p.VoteEndTime, &v, &p.Position,
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

// CreateStory adds a new story to the game
func (d *Service) CreateStory(PokerID string, Name string, Type string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*thunderdome.Story, error) {
	SanitizedDescription := d.HTMLSanitizerPolicy.Sanitize(Description)
	SanitizedAcceptanceCriteria := d.HTMLSanitizerPolicy.Sanitize(AcceptanceCriteria)
	// default priority should be 99 for sort order purposes
	if Priority == 0 {
		Priority = 99
	}
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.poker_story (
		poker_id, name, type, reference_id, link, description, acceptance_criteria, priority, position)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, (
      coalesce(
        (select max(position) from thunderdome.poker_story where poker_id = $1),
        -1
      ) + 1
    ));`,
		PokerID, Name, Type, ReferenceID, Link, SanitizedDescription, SanitizedAcceptanceCriteria, Priority,
	); err != nil {
		d.Logger.Error("error creating poker story", zap.Error(err))
	}

	plans := d.GetStories(PokerID, "")

	return plans, nil
}

// ActivateStoryVoting sets the story by ID to active, wipes any previous votes/points, and disables votingLock
func (d *Service) ActivateStoryVoting(PokerID string, StoryID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_story_activate($1, $2);`, PokerID, StoryID,
	); err != nil {
		d.Logger.Error("CALL thunderdome.poker_story_activate error", zap.Error(err))
	}

	plans := d.GetStories(PokerID, "")

	return plans, nil
}

// SetVote sets a users vote for the story
func (d *Service) SetVote(PokerID string, UserID string, StoryID string, VoteValue string) (Stories []*thunderdome.Story, AllUsersVoted bool) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.poker_story p1
		SET votes = (
			SELECT json_agg(data)
			FROM (
				SELECT coalesce(newVote."warriorId", oldVote."warriorId") AS "warriorId", coalesce(newVote.vote, oldVote.vote) AS vote
				FROM jsonb_populate_recordset(null::thunderdome.UsersVote,p1.votes) AS oldVote
				FULL JOIN jsonb_populate_recordset(null::thunderdome.UsersVote,
					('[{"warriorId":"'|| $2::TEXT ||'", "vote":"'|| $3 ||'"}]')::JSONB
				) AS newVote
				ON newVote."warriorId" = oldVote."warriorId"
			) data
		)
		WHERE p1.id = $1;`,
		StoryID, UserID, VoteValue); err != nil {
		d.Logger.Error("CALL thunderdome.poker_user_vote_set error", zap.Error(err))
	}

	Plans := d.GetStories(PokerID, "")
	ActiveUsers := d.GetActiveUsers(PokerID)

	// determine if all active users have voted
	AllVoted := true
	for _, plan := range Plans {
		if plan.Id == StoryID {
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

// RetractVote removes a users vote for the story
func (d *Service) RetractVote(PokerID string, UserID string, StoryID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.poker_story p1
		SET votes = (
			SELECT coalesce(json_agg(data), '[]'::JSON)
			FROM (
				SELECT coalesce(oldVote."warriorId") AS "warriorId", coalesce(oldVote.vote) AS vote
				FROM jsonb_populate_recordset(null::thunderdome.UsersVote,p1.votes) AS oldVote
				WHERE oldVote."warriorId" != $2
			) data
		)
		WHERE p1.id = $1;
    `, StoryID, UserID); err != nil {
		return nil, fmt.Errorf("poker retract vote query error: %v", err)
	}

	plans := d.GetStories(PokerID, "")

	return plans, nil
}

// EndStoryVoting sets story to active: false
func (d *Service) EndStoryVoting(PokerID string, StoryID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_plan_voting_stop($1, $2);`, PokerID, StoryID); err != nil {
		d.Logger.Error("CALL thunderdome.poker_plan_voting_stop error", zap.Error(err))
	}

	plans := d.GetStories(PokerID, "")

	return plans, nil
}

// SkipStory sets story to active: false and unsets games activeStoryId
func (d *Service) SkipStory(PokerID string, StoryID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_vote_skip($1, $2);`, PokerID, StoryID); err != nil {
		d.Logger.Error("CALL thunderdome.poker_vote_skip error", zap.Error(err))
	}

	plans := d.GetStories(PokerID, "")

	return plans, nil
}

// UpdateStory updates the story by ID
func (d *Service) UpdateStory(PokerID string, StoryID string, Name string, Type string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*thunderdome.Story, error) {
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
		StoryID, Name, Type, ReferenceID, Link, SanitizedDescription, SanitizedAcceptanceCriteria, Priority); err != nil {
		d.Logger.Error("error getting poker story", zap.Error(err))
	}

	plans := d.GetStories(PokerID, "")

	return plans, nil
}

// DeleteStory removes a story from the current game by ID
func (d *Service) DeleteStory(PokerID string, StoryID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_story_delete($1, $2);`, PokerID, StoryID); err != nil {
		d.Logger.Error("CALL thunderdome.poker_story_delete error", zap.Error(err))
	}

	plans := d.GetStories(PokerID, "")

	return plans, nil
}

// ArrangeStory sets the position of the story relative to the story it's being placed before
func (d *Service) ArrangeStory(PokerID string, StoryID string, BeforeStoryID string) ([]*thunderdome.Story, error) {
	if BeforeStoryID == "" {
		_, err := d.DB.Exec(`UPDATE thunderdome.poker_story SET 
			position = (SELECT max(position) FROM thunderdome.poker_story WHERE poker_id = $1) + 1
			WHERE id = $2;`,
			PokerID, StoryID)
		if err != nil {
			d.Logger.Error("poker ArrangeStory get beforeStoryId error", zap.Error(err))
		}
	} else {
		_, err := d.DB.Exec(
			`UPDATE thunderdome.poker_story SET position = (
			  -- find position of item referenced in before argument (default to 0)
			  with "before_position" as (
				select coalesce(
				  (select "position" from thunderdome.poker_story where id = $3),
				  -- in case item was not found, use last item in list and add 1 to add item to end of list
				  (select max("position") + 1 from thunderdome.poker_story where poker_id = $1),
				  -- in case no item exists, use 0
				  0
				) as "position"
			  ),
			  -- find position of previous item relative to "before item"
			  "before_prev_position" as (
				select coalesce(
				  (
					select w."position"
					from thunderdome.poker_story w, "before_position" b
					where
					  w.poker_id = $1 and
					  -- positions may not be integers, so we cannot simply deduct 1
					  -- this is why we find the first item with a smaller position
					  w."position" < b."position"
					order by "position" desc limit 1
				  ),
				  -- in case previous position does not exist (before item was first in list), simply deduct 1
				  (select b.position - 1 from "before_position" b)
				) as "position"
			  )
			  -- average both positions to fit new row into gap
			  select (b.position + p.position) / 2
			  from "before_position" b, "before_prev_position" p
			) WHERE id = $2;`,
			PokerID, StoryID, BeforeStoryID)
		if err != nil {
			d.Logger.Error("poker ArrangeStory error", zap.Error(err))
		}
	}

	plans := d.GetStories(PokerID, "")

	return plans, nil
}

// FinalizeStory sets story to active: false and updates the points
func (d *Service) FinalizeStory(PokerID string, StoryID string, Points string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_story_finalize($1, $2, $3);`, PokerID, StoryID, Points); err != nil {
		d.Logger.Error("CALL thunderdome.poker_story_finalize error", zap.Error(err))
	}

	plans := d.GetStories(PokerID, "")

	return plans, nil
}
