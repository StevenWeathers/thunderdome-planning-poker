package poker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
)

// UserNudge handles notifying user that they need to vote
func (b *Service) UserNudge(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	msg := wshub.CreateSocketEvent("jab_warrior", EventValue, UserID)

	return msg, nil, false
}

// UserVote handles the participants vote event by setting their vote
// and checks if AutoFinishVoting && AllVoted if so ends voting
func (b *Service) UserVote(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	var msg []byte
	var wv struct {
		VoteValue        string `json:"voteValue"`
		StoryID          string `json:"planId"`
		AutoFinishVoting bool   `json:"autoFinishVoting"`
	}
	err := json.Unmarshal([]byte(EventValue), &wv)
	if err != nil {
		return nil, err, false
	}

	Storys, AllVoted := b.BattleService.SetVote(BattleID, UserID, wv.StoryID, wv.VoteValue)

	updatedStorys, _ := json.Marshal(Storys)
	msg = wshub.CreateSocketEvent("vote_activity", string(updatedStorys), UserID)

	if AllVoted && wv.AutoFinishVoting {
		plans, err := b.BattleService.EndStoryVoting(BattleID, wv.StoryID)
		if err != nil {
			return nil, err, false
		}
		updatedStorys, _ := json.Marshal(plans)
		msg = wshub.CreateSocketEvent("voting_ended", string(updatedStorys), "")
	}

	return msg, nil, false
}

// UserVoteRetract handles retracting a user vote
func (b *Service) UserVoteRetract(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	StoryID := EventValue

	plans, err := b.BattleService.RetractVote(BattleID, UserID, StoryID)
	if err != nil {
		return nil, err, false
	}

	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("vote_retracted", string(updatedStorys), UserID)

	return msg, nil, false
}

// UserPromote handles promoting a user to a leader
func (b *Service) UserPromote(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	leaders, err := b.BattleService.AddFacilitator(BattleID, EventValue)
	if err != nil {
		return nil, err, false
	}
	leadersJson, _ := json.Marshal(leaders)

	msg := wshub.CreateSocketEvent("leaders_updated", string(leadersJson), "")

	return msg, nil, false
}

// UserDemote handles demoting a user from a leader
func (b *Service) UserDemote(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	leaders, err := b.BattleService.RemoveFacilitator(BattleID, EventValue)
	if err != nil {
		return nil, err, false
	}
	leadersJson, _ := json.Marshal(leaders)

	msg := wshub.CreateSocketEvent("leaders_updated", string(leadersJson), "")

	return msg, nil, false
}

// UserPromoteSelf handles self-promoting a user to a leader
func (b *Service) UserPromoteSelf(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	leaderCode, err := b.BattleService.GetFacilitatorCode(BattleID)
	if err != nil {
		return nil, err, false
	}

	if EventValue == leaderCode {
		leaders, err := b.BattleService.AddFacilitator(BattleID, UserID)
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
func (b *Service) UserSpectatorToggle(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	var st struct {
		Spectator bool `json:"spectator"`
	}
	err := json.Unmarshal([]byte(EventValue), &st)
	if err != nil {
		return nil, err, false
	}
	users, err := b.BattleService.ToggleSpectator(BattleID, UserID, st.Spectator)
	if err != nil {
		return nil, err, false
	}
	usersJson, _ := json.Marshal(users)

	msg := wshub.CreateSocketEvent("users_updated", string(usersJson), "")

	return msg, nil, false
}

// StoryVoteEnd handles ending plan voting
func (b *Service) StoryVoteEnd(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	plans, err := b.BattleService.EndStoryVoting(BattleID, EventValue)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("voting_ended", string(updatedStorys), "")

	return msg, nil, false
}

// Revise handles editing the battle settings
func (b *Service) Revise(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
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
	err := json.Unmarshal([]byte(EventValue), &rb)
	if err != nil {
		return nil, err, false
	}

	err = b.BattleService.UpdateGame(
		BattleID,
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

// Delete handles deleting the battle
func (b *Service) Delete(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	err := b.BattleService.DeleteGame(BattleID)
	if err != nil {
		return nil, err, false
	}
	msg := wshub.CreateSocketEvent("battle_conceded", "", "")

	return msg, nil, false
}

// StoryAdd adds a new plan to the battle
func (b *Service) StoryAdd(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	var p struct {
		Name               string `json:"planName"`
		Type               string `json:"type"`
		ReferenceId        string `json:"referenceId"`
		Link               string `json:"link"`
		Description        string `json:"description"`
		AcceptanceCriteria string `json:"acceptanceCriteria"`
		Priority           int32  `json:"priority"`
	}
	err := json.Unmarshal([]byte(EventValue), &p)
	if err != nil {
		return nil, err, false
	}

	plans, err := b.BattleService.CreateStory(BattleID, p.Name, p.Type, p.ReferenceId, p.Link, p.Description, p.AcceptanceCriteria, p.Priority)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_added", string(updatedStorys), "")

	return msg, nil, false
}

// StoryRevise handles editing a battle plan
func (b *Service) StoryRevise(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	var p struct {
		Id                 string `json:"planId"`
		Name               string `json:"planName"`
		Type               string `json:"type"`
		ReferenceId        string `json:"referenceId"`
		Link               string `json:"link"`
		Description        string `json:"description"`
		AcceptanceCriteria string `json:"acceptanceCriteria"`
		Priority           int32  `json:"priority"`
	}
	err := json.Unmarshal([]byte(EventValue), &p)
	if err != nil {
		return nil, err, false
	}

	plans, err := b.BattleService.UpdateStory(BattleID, p.Id, p.Name, p.Type, p.ReferenceId, p.Link, p.Description, p.AcceptanceCriteria, p.Priority)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_revised", string(updatedStorys), "")

	return msg, nil, false
}

// StoryDelete handles deleting a plan
func (b *Service) StoryDelete(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	plans, err := b.BattleService.DeleteStory(BattleID, EventValue)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_burned", string(updatedStorys), "")

	return msg, nil, false
}

// StoryArrange sets the position of the story relative to the beforeStory
func (b *Service) StoryArrange(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	var p struct {
		StoryID       string `json:"story_id"`
		BeforeStoryID string `json:"before_story_id"`
	}
	err := json.Unmarshal([]byte(EventValue), &p)
	if err != nil {
		return nil, err, false
	}

	plans, err := b.BattleService.ArrangeStory(BattleID, p.StoryID, p.BeforeStoryID)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("story_arranged", string(updatedStorys), "")

	return msg, nil, false
}

// StoryActivate handles activating a plan for voting
func (b *Service) StoryActivate(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	plans, err := b.BattleService.ActivateStoryVoting(BattleID, EventValue)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_activated", string(updatedStorys), "")

	return msg, nil, false
}

// StorySkip handles skipping a plan voting
func (b *Service) StorySkip(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	plans, err := b.BattleService.SkipStory(BattleID, EventValue)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_skipped", string(updatedStorys), "")

	return msg, nil, false
}

// StoryFinalize handles setting a plan point value
func (b *Service) StoryFinalize(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	var p struct {
		Id     string `json:"planId"`
		Points string `json:"planPoints"`
	}
	err := json.Unmarshal([]byte(EventValue), &p)
	if err != nil {
		return nil, err, false
	}

	plans, err := b.BattleService.FinalizeStory(BattleID, p.Id, p.Points)
	if err != nil {
		return nil, err, false
	}
	updatedStorys, _ := json.Marshal(plans)
	msg := wshub.CreateSocketEvent("plan_finalized", string(updatedStorys), "")

	return msg, nil, false
}

// Abandon handles setting abandoned true so battle doesn't show up in users battle list, then leaves battle
func (b *Service) Abandon(ctx context.Context, BattleID string, UserID string, EventValue string) ([]byte, error, bool) {
	_, err := b.BattleService.AbandonGame(BattleID, UserID)
	if err != nil {
		return nil, err, false
	}

	return nil, errors.New("ABANDONED_BATTLE"), true
}
