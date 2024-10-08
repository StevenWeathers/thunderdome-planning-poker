package poker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
)

// UserNudge handles notifying user that they need to vote
func (b *Service) UserNudge(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	msg := wshub.CreateSocketEvent("jab_warrior", eventValue, userID)

	return msg, nil, false
}

// UserVote handles the participants vote event by setting their vote
// and checks if AutoFinishVoting && AllVoted if so ends voting
func (b *Service) UserVote(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	var msg []byte
	var wv struct {
		VoteValue        string `json:"voteValue"`
		StoryID          string `json:"planId"`
		AutoFinishVoting bool   `json:"autoFinishVoting"`
	}
	err := json.Unmarshal([]byte(eventValue), &wv)
	if err != nil {
		return nil, err, false
	}

	storys, allVoted := b.PokerService.SetVote(pokerID, userID, wv.StoryID, wv.VoteValue)

	updatedStorys, _ := json.Marshal(storys)
	msg = wshub.CreateSocketEvent("vote_activity", string(updatedStorys), userID)

	if allVoted && wv.AutoFinishVoting {
		plans, err := b.PokerService.EndStoryVoting(pokerID, wv.StoryID)
		if err != nil {
			return nil, err, false
		}
		updatedStorys, _ := json.Marshal(plans)
		msg = wshub.CreateSocketEvent("voting_ended", string(updatedStorys), "")
	}

	return msg, nil, false
}

// UserVoteRetract handles retracting a user vote
func (b *Service) UserVoteRetract(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	storyID := eventValue

	plans, err := b.PokerService.RetractVote(pokerID, userID, storyID)
	if err != nil {
		return nil, err, false
	}

	updatedStories, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("vote_retracted", string(updatedStories), userID)

	return msg, nil, false
}

// UserPromote handles promoting a user to a facilitator
func (b *Service) UserPromote(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	leaders, err := b.PokerService.AddFacilitator(pokerID, eventValue)
	if err != nil {
		return nil, err, false
	}
	leadersJson, _ := json.Marshal(leaders)

	msg := wshub.CreateSocketEvent("leaders_updated", string(leadersJson), "")

	return msg, nil, false
}

// UserDemote handles demoting a user from a facilitator
func (b *Service) UserDemote(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	leaders, err := b.PokerService.RemoveFacilitator(pokerID, eventValue)
	if err != nil {
		return nil, err, false
	}
	leadersJson, _ := json.Marshal(leaders)

	msg := wshub.CreateSocketEvent("leaders_updated", string(leadersJson), "")

	return msg, nil, false
}

// UserPromoteSelf handles self-promoting a user to a facilitator
func (b *Service) UserPromoteSelf(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	leaderCode, err := b.PokerService.GetFacilitatorCode(pokerID)
	if err != nil {
		return nil, err, false
	}

	if eventValue == leaderCode {
		leaders, err := b.PokerService.AddFacilitator(pokerID, userID)
		if err != nil {
			return nil, err, false
		}
		leadersJson, _ := json.Marshal(leaders)

		msg := wshub.CreateSocketEvent("leaders_updated", string(leadersJson), "")

		return msg, nil, false
	} else {
		return nil, errors.New("INCORRECT_LEADER_CODE"), false
	}
}

// UserSpectatorToggle handles toggling user spectator status
func (b *Service) UserSpectatorToggle(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	var st struct {
		Spectator bool `json:"spectator"`
	}
	err := json.Unmarshal([]byte(eventValue), &st)
	if err != nil {
		return nil, err, false
	}
	users, err := b.PokerService.ToggleSpectator(pokerID, userID, st.Spectator)
	if err != nil {
		return nil, err, false
	}
	usersJson, _ := json.Marshal(users)

	msg := wshub.CreateSocketEvent("users_updated", string(usersJson), "")

	return msg, nil, false
}

// StoryVoteEnd handles ending story voting
func (b *Service) StoryVoteEnd(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	plans, err := b.PokerService.EndStoryVoting(pokerID, eventValue)
	if err != nil {
		return nil, err, false
	}
	updatedStories, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("voting_ended", string(updatedStories), "")

	return msg, nil, false
}

// Revise handles editing the poker game settings
func (b *Service) Revise(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	var rb struct {
		BattleName           string   `json:"battleName"`
		PointValuesAllowed   []string `json:"pointValuesAllowed"`
		AutoFinishVoting     bool     `json:"autoFinishVoting"`
		PointAverageRounding string   `json:"pointAverageRounding"`
		HideVoterIdentity    bool     `json:"hideVoterIdentity"`
		JoinCode             string   `json:"joinCode"`
		LeaderCode           string   `json:"leaderCode"`
		TeamID               string   `json:"teamId"`
	}
	err := json.Unmarshal([]byte(eventValue), &rb)
	if err != nil {
		return nil, err, false
	}

	err = b.PokerService.UpdateGame(
		pokerID,
		rb.BattleName,
		rb.PointValuesAllowed,
		rb.AutoFinishVoting,
		rb.PointAverageRounding,
		rb.HideVoterIdentity,
		rb.JoinCode,
		rb.LeaderCode,
		rb.TeamID,
	)
	if err != nil {
		return nil, err, false
	}

	rb.LeaderCode = ""

	updatedBattle, _ := json.Marshal(rb)
	msg := wshub.CreateSocketEvent("battle_revised", string(updatedBattle), "")

	return msg, nil, false
}

// Delete handles deleting the poker game
func (b *Service) Delete(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	err := b.PokerService.DeleteGame(pokerID)
	if err != nil {
		return nil, err, false
	}
	msg := wshub.CreateSocketEvent("battle_conceded", "", "")

	return msg, nil, false
}

// StoryAdd adds a new story to the poker game
func (b *Service) StoryAdd(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	var p struct {
		Name               string `json:"planName"`
		Type               string `json:"type"`
		ReferenceID        string `json:"referenceId"`
		Link               string `json:"link"`
		Description        string `json:"description"`
		AcceptanceCriteria string `json:"acceptanceCriteria"`
		Priority           int32  `json:"priority"`
	}
	err := json.Unmarshal([]byte(eventValue), &p)
	if err != nil {
		return nil, err, false
	}

	plans, err := b.PokerService.CreateStory(pokerID, p.Name, p.Type, p.ReferenceID, p.Link, p.Description, p.AcceptanceCriteria, p.Priority)
	if err != nil {
		return nil, err, false
	}
	updatedStories, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_added", string(updatedStories), "")

	return msg, nil, false
}

// StoryRevise handles editing a poker story
func (b *Service) StoryRevise(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	var p struct {
		ID                 string `json:"planId"`
		Name               string `json:"planName"`
		Type               string `json:"type"`
		ReferenceID        string `json:"referenceId"`
		Link               string `json:"link"`
		Description        string `json:"description"`
		AcceptanceCriteria string `json:"acceptanceCriteria"`
		Priority           int32  `json:"priority"`
	}
	err := json.Unmarshal([]byte(eventValue), &p)
	if err != nil {
		return nil, err, false
	}

	stories, err := b.PokerService.UpdateStory(pokerID, p.ID, p.Name, p.Type, p.ReferenceID, p.Link, p.Description, p.AcceptanceCriteria, p.Priority)
	if err != nil {
		return nil, err, false
	}
	updatedStories, _ := json.Marshal(stories)
	msg := wshub.CreateSocketEvent("plan_revised", string(updatedStories), "")

	return msg, nil, false
}

// StoryDelete handles deleting a story
func (b *Service) StoryDelete(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	plans, err := b.PokerService.DeleteStory(pokerID, eventValue)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_burned", string(updatedStorys), "")

	return msg, nil, false
}

// StoryArrange sets the position of the story relative to the beforeStory
func (b *Service) StoryArrange(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	var p struct {
		StoryID       string `json:"story_id"`
		BeforeStoryID string `json:"before_story_id"`
	}
	err := json.Unmarshal([]byte(eventValue), &p)
	if err != nil {
		return nil, err, false
	}

	plans, err := b.PokerService.ArrangeStory(pokerID, p.StoryID, p.BeforeStoryID)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("story_arranged", string(updatedStorys), "")

	return msg, nil, false
}

// StoryActivate handles activating a story for voting
func (b *Service) StoryActivate(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	plans, err := b.PokerService.ActivateStoryVoting(pokerID, eventValue)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_activated", string(updatedStorys), "")

	return msg, nil, false
}

// StorySkip handles skipping a story voting
func (b *Service) StorySkip(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	plans, err := b.PokerService.SkipStory(pokerID, eventValue)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_skipped", string(updatedStorys), "")

	return msg, nil, false
}

// StoryFinalize handles setting a story point value
func (b *Service) StoryFinalize(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	var p struct {
		ID     string `json:"planId"`
		Points string `json:"planPoints"`
	}
	err := json.Unmarshal([]byte(eventValue), &p)
	if err != nil {
		return nil, err, false
	}

	plans, err := b.PokerService.FinalizeStory(pokerID, p.ID, p.Points)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_finalized", string(updatedStorys), "")

	return msg, nil, false
}

// Abandon handles setting abandoned true so game doesn't show up in users poker game list, then leaves game
func (b *Service) Abandon(ctx context.Context, pokerID string, userID string, eventValue string) ([]byte, error, bool) {
	_, err := b.PokerService.AbandonGame(pokerID, userID)
	if err != nil {
		return nil, err, false
	}

	return nil, errors.New("ABANDONED_BATTLE"), true
}
