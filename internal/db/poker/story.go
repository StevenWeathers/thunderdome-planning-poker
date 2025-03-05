package poker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"
)

// GetStories retrieves stories for given poker game
func (d *Service) GetStories(pokerID string, userID string) []*thunderdome.Story {
	// 尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("game:%s:stories", pokerID)
	if d.Redis != nil {
		if cachedData, err := d.Redis.Get(context.Background(), cacheKey).Result(); err == nil {
			var stories []*thunderdome.Story
			if err := json.Unmarshal([]byte(cachedData), &stories); err == nil {
				d.Logger.Debug("Stories cache hit", zap.String("game_id", pokerID))
				return stories
			}
		}
	}

	var stories = make([]*thunderdome.Story, 0)
	storyRows, storiesErr := d.DB.Query(
		`SELECT
			id, name, type, reference_id, link, description, acceptance_criteria, priority,
			points, active, skipped, votestart_time, voteend_time, votes,
			row_number() OVER (ORDER BY position ASC) as position
			FROM thunderdome.poker_story WHERE poker_id = $1 ORDER BY position
		`,
		pokerID,
	)
	if storiesErr == nil {
		defer storyRows.Close()
		for storyRows.Next() {
			var v string
			var referenceID sql.NullString
			var link sql.NullString
			var description sql.NullString
			var acceptanceCriteria sql.NullString
			var p = &thunderdome.Story{
				Votes:   make([]*thunderdome.Vote, 0),
				Active:  false,
				Skipped: false,
			}
			if err := storyRows.Scan(
				&p.ID,
				&p.Name,
				&p.Type,
				&referenceID,
				&link,
				&description,
				&acceptanceCriteria,
				&p.Priority,
				&p.Points,
				&p.Active,
				&p.Skipped,
				&p.VoteStartTime,
				&p.VoteEndTime,
				&v,
				&p.Position,
			); err != nil {
				d.Logger.Error("error getting poker stories", zap.Error(err))
			} else {
				p.ReferenceID = referenceID.String
				p.Link = link.String
				p.Description = description.String
				p.AcceptanceCriteria = acceptanceCriteria.String
				_ = json.Unmarshal([]byte(v), &p.Votes)
				stories = append(stories, p)
			}
		}
	}

	// 设置缓存
	if d.Redis != nil {
		if storiesJSON, err := json.Marshal(stories); err == nil {
			d.Redis.Set(context.Background(), cacheKey, storiesJSON, 1*time.Hour)
		}
	}

	return stories
}

// CreateStory adds a new story to the game
func (d *Service) CreateStory(pokerID string, name string, storyType string, referenceID string, link string, description string, acceptanceCriteria string, priority int32) ([]*thunderdome.Story, error) {
	sanitizedDescription := d.HTMLSanitizerPolicy.Sanitize(description)
	sanitizedAcceptanceCriteria := d.HTMLSanitizerPolicy.Sanitize(acceptanceCriteria)
	// default priority should be 99 for sort order purposes
	if priority == 0 {
		priority = 99
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
		pokerID, name, storyType, referenceID, link, sanitizedDescription, sanitizedAcceptanceCriteria, priority,
	); err != nil {
		d.Logger.Error("error creating poker story", zap.Error(err),
			zap.String("PokerID", pokerID), zap.String("Name", name))
	}

	// 清除缓存
	if d.Redis != nil {
		cacheKey := fmt.Sprintf("game:%s:stories", pokerID)
		d.Redis.Del(context.Background(), cacheKey)
	}

	stories := d.GetStories(pokerID, "")

	return stories, nil
}

// ActivateStoryVoting sets the story by ID to active, wipes any previous votes/points, and disables votingLock
func (d *Service) ActivateStoryVoting(pokerID string, storyID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_story_activate($1, $2);`, pokerID, storyID,
	); err != nil {
		d.Logger.Error("CALL thunderdome.poker_story_activate error", zap.Error(err),
			zap.String("PokerID", pokerID), zap.String("StoryID", storyID))
	}

	stories := d.GetStories(pokerID, "")

	return stories, nil
}

// SetVote sets a users vote for the story
func (d *Service) SetVote(pokerID string, userID string, storyID string, voteValue string) (Stories []*thunderdome.Story, allUsersVoted bool) {
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
		storyID, userID, voteValue); err != nil {
		d.Logger.Error("CALL thunderdome.poker_user_vote_set error", zap.Error(err),
			zap.String("PokerID", pokerID), zap.String("UserID", userID),
			zap.String("StoryID", storyID), zap.String("VoteValue", voteValue))
	}

	// 清除缓存
	if d.Redis != nil {
		cacheKey := fmt.Sprintf("game:%s:stories", pokerID)
		d.Redis.Del(context.Background(), cacheKey)
	}

	stories := d.GetStories(pokerID, "")
	activeUsers := d.GetActiveUsers(pokerID)

	// determine if all active users have voted
	allVoted := true
	for _, story := range stories {
		if story.ID == storyID {
			activePlanVoters := make(map[string]bool)

			for _, vote := range story.Votes {
				activePlanVoters[vote.UserID] = true
			}
			for _, war := range activeUsers {
				if _, UserVoted := activePlanVoters[war.ID]; !UserVoted && !war.Spectator {
					allVoted = false
					break
				}
			}
			break
		}
	}

	return stories, allVoted
}

// RetractVote removes a users vote for the story
func (d *Service) RetractVote(pokerID string, userID string, storyID string) ([]*thunderdome.Story, error) {
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
    `, storyID, userID); err != nil {
		d.Logger.Error("poker retract vote query error", zap.Error(err),
			zap.String("PokerID", pokerID), zap.String("UserID", userID), zap.String("StoryID", storyID))
		return nil, fmt.Errorf("poker retract vote query error: %v", err)
	}

	// 清除缓存
	if d.Redis != nil {
		cacheKey := fmt.Sprintf("game:%s:stories", pokerID)
		d.Redis.Del(context.Background(), cacheKey)
	}

	stories := d.GetStories(pokerID, "")

	return stories, nil
}

// EndStoryVoting sets story to active: false
func (d *Service) EndStoryVoting(pokerID string, storyID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_plan_voting_stop($1, $2);`, pokerID, storyID); err != nil {
		d.Logger.Error("CALL thunderdome.poker_plan_voting_stop error", zap.Error(err),
			zap.String("PokerID", pokerID), zap.String("StoryID", storyID))
	}

	// 清除缓存
	if d.Redis != nil {
		cacheKey := fmt.Sprintf("game:%s:stories", pokerID)
		d.Redis.Del(context.Background(), cacheKey)
	}

	stories := d.GetStories(pokerID, "")

	return stories, nil
}

// SkipStory sets story to active: false and unsets games activeStoryId
func (d *Service) SkipStory(pokerID string, storyID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_vote_skip($1, $2);`, pokerID, storyID); err != nil {
		d.Logger.Error("CALL thunderdome.poker_vote_skip error", zap.Error(err),
			zap.String("PokerID", pokerID), zap.String("StoryID", storyID))
	}

	// 清除缓存
	if d.Redis != nil {
		cacheKey := fmt.Sprintf("game:%s:stories", pokerID)
		d.Redis.Del(context.Background(), cacheKey)
	}

	stories := d.GetStories(pokerID, "")

	return stories, nil
}

// UpdateStory updates the story by ID
func (d *Service) UpdateStory(pokerID string, storyID string, name string, storyType string, referenceID string, link string, description string, acceptanceCriteria string, priority int32) ([]*thunderdome.Story, error) {
	sanitizedDescription := d.HTMLSanitizerPolicy.Sanitize(description)
	sanitizedAcceptanceCriteria := d.HTMLSanitizerPolicy.Sanitize(acceptanceCriteria)
	// default priority should be 99 for sort order purposes
	if priority == 0 {
		priority = 99
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
		storyID, name, storyType, referenceID, link, sanitizedDescription, sanitizedAcceptanceCriteria, priority); err != nil {
		d.Logger.Error("error getting poker story", zap.Error(err),
			zap.String("PokerID", pokerID), zap.String("StoryID", storyID))
	}

	// 清除缓存
	if d.Redis != nil {
		cacheKey := fmt.Sprintf("game:%s:stories", pokerID)
		d.Redis.Del(context.Background(), cacheKey)
	}

	stories := d.GetStories(pokerID, "")

	return stories, nil
}

// DeleteStory removes a story from the current game by ID
func (d *Service) DeleteStory(pokerID string, storyID string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_story_delete($1, $2);`, pokerID, storyID); err != nil {
		d.Logger.Error("CALL thunderdome.poker_story_delete error", zap.Error(err),
			zap.String("PokerID", pokerID), zap.String("StoryID", storyID))
	}

	// 清除缓存
	if d.Redis != nil {
		cacheKey := fmt.Sprintf("game:%s:stories", pokerID)
		d.Redis.Del(context.Background(), cacheKey)
	}

	stories := d.GetStories(pokerID, "")

	return stories, nil
}

// ArrangeStory sets the position of the story relative to the story it's being placed before
func (d *Service) ArrangeStory(pokerID string, storyID string, beforeStoryID string) ([]*thunderdome.Story, error) {
	if beforeStoryID == "" {
		_, err := d.DB.Exec(`UPDATE thunderdome.poker_story SET
			position = (SELECT max(position) FROM thunderdome.poker_story WHERE poker_id = $1) + 1
			WHERE id = $2;`,
			pokerID, storyID)
		if err != nil {
			d.Logger.Error("poker ArrangeStory get beforeStoryId error", zap.Error(err),
				zap.String("PokerID", pokerID), zap.String("StoryID", storyID))
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
			pokerID, storyID, beforeStoryID)
		if err != nil {
			d.Logger.Error("poker ArrangeStory error", zap.Error(err),
				zap.String("PokerID", pokerID), zap.String("StoryID", storyID),
				zap.String("BeforeStoryID", beforeStoryID))
		}
	}

	stories := d.GetStories(pokerID, "")

	return stories, nil
}

// FinalizeStory sets story to active: false and updates the points
func (d *Service) FinalizeStory(pokerID string, storyID string, points string) ([]*thunderdome.Story, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.poker_story_finalize($1, $2, $3);`, pokerID, storyID, points); err != nil {
		d.Logger.Error("CALL thunderdome.poker_story_finalize error", zap.Error(err),
			zap.String("PokerID", pokerID),
			zap.String("StoryID", storyID),
			zap.String("Points", points))
	}

	// 清除缓存
	if d.Redis != nil {
		cacheKey := fmt.Sprintf("game:%s:stories", pokerID)
		d.Redis.Del(context.Background(), cacheKey)
	}

	stories := d.GetStories(pokerID, "")

	return stories, nil
}
