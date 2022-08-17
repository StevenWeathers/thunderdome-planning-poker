package storyboard

import (
	"context"
	"encoding/json"
	"errors"
)

// AddGoal handles adding a goal to storyboard
func (b *Service) AddGoal(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goals, err := b.db.CreateStoryboardGoal(StoryboardID, UserID, EventValue)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("goal_added", string(updatedGoals), "")

	return msg, nil, false
}

// ReviseGoal handles revising a storyboard goal
func (b *Service) ReviseGoal(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(EventValue), &goalObj)
	if err != nil {
		return nil, err, false
	}
	GoalID := goalObj["goalId"]
	GoalName := goalObj["name"]

	goals, err := b.db.ReviseGoalName(StoryboardID, UserID, GoalID, GoalName)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("goal_revised", string(updatedGoals), "")

	return msg, nil, false
}

// DeleteGoal handles deleting a storyboard goal
func (b *Service) DeleteGoal(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goals, err := b.db.DeleteStoryboardGoal(StoryboardID, UserID, EventValue)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("goal_deleted", string(updatedGoals), "")

	return msg, nil, false
}

// AddColumn handles adding a column to storyboard goal
func (b *Service) AddColumn(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(EventValue), &goalObj)
	if err != nil {
		return nil, err, false
	}
	GoalID := goalObj["goalId"]

	goals, err := b.db.CreateStoryboardColumn(StoryboardID, GoalID, UserID)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("column_added", string(updatedGoals), "")

	return msg, nil, false
}

// ReviseColumn handles revising a storyboard goal column
func (b *Service) ReviseColumn(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		ColumnID string `json:"id"`
		Name     string `json:"name"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	goals, err := b.db.ReviseStoryboardColumn(StoryboardID, UserID, rs.ColumnID, rs.Name)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("column_updated", string(updatedGoals), "")

	return msg, nil, false
}

// DeleteColumn handles deleting a storyboard goal column
func (b *Service) DeleteColumn(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goals, err := b.db.DeleteStoryboardColumn(StoryboardID, UserID, EventValue)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_deleted", string(updatedGoals), "")

	return msg, nil, false
}

// AddStory handles adding a story to storyboard
func (b *Service) AddStory(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(EventValue), &goalObj)
	if err != nil {
		return nil, err, false
	}
	GoalID := goalObj["goalId"]
	ColumnID := goalObj["columnId"]

	goals, err := b.db.CreateStoryboardStory(StoryboardID, GoalID, ColumnID, UserID)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_added", string(updatedGoals), "")

	return msg, nil, false
}

// UpdateStoryName handles revising a storyboard story name
func (b *Service) UpdateStoryName(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(EventValue), &goalObj)
	if err != nil {
		return nil, err, false
	}
	StoryID := goalObj["storyId"]
	StoryName := goalObj["name"]

	goals, err := b.db.ReviseStoryName(StoryboardID, UserID, StoryID, StoryName)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_updated", string(updatedGoals), "")

	return msg, nil, false
}

// UpdateStoryContent handles revising a storyboard story content
func (b *Service) UpdateStoryContent(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(EventValue), &goalObj)
	if err != nil {
		return nil, err, false
	}
	StoryID := goalObj["storyId"]
	StoryContent := goalObj["content"]

	goals, err := b.db.ReviseStoryContent(StoryboardID, UserID, StoryID, StoryContent)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_updated", string(updatedGoals), "")

	return msg, nil, false
}

// UpdateStoryColor handles revising a storyboard story color
func (b *Service) UpdateStoryColor(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(EventValue), &goalObj)
	if err != nil {
		return nil, err, false
	}
	StoryID := goalObj["storyId"]
	StoryColor := goalObj["color"]

	goals, err := b.db.ReviseStoryColor(StoryboardID, UserID, StoryID, StoryColor)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_updated", string(updatedGoals), "")

	return msg, nil, false
}

// UpdateStoryPoints handles revising a storyboard story points
func (b *Service) UpdateStoryPoints(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		StoryID string `json:"storyId"`
		Points  int    `json:"points"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	goals, err := b.db.ReviseStoryPoints(StoryboardID, UserID, rs.StoryID, rs.Points)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_updated", string(updatedGoals), "")

	return msg, nil, false
}

// UpdateStoryClosed handles revising a storyboard story closed status
func (b *Service) UpdateStoryClosed(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		StoryID string `json:"storyId"`
		Closed  bool   `json:"closed"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	goals, err := b.db.ReviseStoryClosed(StoryboardID, UserID, rs.StoryID, rs.Closed)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_updated", string(updatedGoals), "")

	return msg, nil, false
}

// UpdateStoryLink handles revising a storyboard story link
func (b *Service) UpdateStoryLink(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(EventValue), &goalObj)
	if err != nil {
		return nil, err, false
	}
	StoryID := goalObj["storyId"]
	Link := goalObj["link"]

	goals, err := b.db.ReviseStoryLink(StoryboardID, UserID, StoryID, Link)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_updated", string(updatedGoals), "")

	return msg, nil, false
}

// MoveStory handles moving a storyboard story between columns/goals
func (b *Service) MoveStory(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(EventValue), &goalObj)
	if err != nil {
		return nil, err, false
	}
	StoryID := goalObj["storyId"]
	GoalID := goalObj["goalId"]
	ColumnID := goalObj["columnId"]
	PlaceBefore := goalObj["placeBefore"]

	goals, err := b.db.MoveStoryboardStory(StoryboardID, UserID, StoryID, GoalID, ColumnID, PlaceBefore)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_moved", string(updatedGoals), "")

	return msg, nil, false
}

// DeleteStory handles deleting a storyboard story
func (b *Service) DeleteStory(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goals, err := b.db.DeleteStoryboardStory(StoryboardID, UserID, EventValue)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_deleted", string(updatedGoals), "")

	return msg, nil, false
}

// AddStoryComment handles adding a storyboard story comment
func (b *Service) AddStoryComment(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		StoryID string `json:"storyId"`
		Comment string `json:"comment"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	goals, err := b.db.AddStoryComment(StoryboardID, UserID, rs.StoryID, rs.Comment)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_updated", string(updatedGoals), "")

	return msg, nil, false
}

// EditStoryComment handles editing a storyboard story comment
func (b *Service) EditStoryComment(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		CommentID string `json:"commentId"`
		Comment   string `json:"comment"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	goals, err := b.db.EditStoryComment(StoryboardID, rs.CommentID, rs.Comment)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_updated", string(updatedGoals), "")

	return msg, nil, false
}

// DeleteStoryComment handles deleting a storyboard story comment
func (b *Service) DeleteStoryComment(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		CommentID string `json:"commentId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	goals, err := b.db.DeleteStoryComment(StoryboardID, rs.CommentID)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("story_updated", string(updatedGoals), "")

	return msg, nil, false
}

// AddPersona handles adding a storyboard persona
func (b *Service) AddPersona(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		Name        string `json:"name"`
		Role        string `json:"role"`
		Description string `json:"description"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	personas, err := b.db.AddStoryboardPersona(StoryboardID, UserID, rs.Name, rs.Role, rs.Description)
	if err != nil {
		return nil, err, false
	}
	updatedPersonas, _ := json.Marshal(personas)
	msg := createSocketEvent("personas_updated", string(updatedPersonas), "")

	return msg, nil, false
}

// UpdatePersona handles updating a storyboard persona
func (b *Service) UpdatePersona(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		PersonaID   string `json:"id"`
		Name        string `json:"name"`
		Role        string `json:"role"`
		Description string `json:"description"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	personas, err := b.db.UpdateStoryboardPersona(StoryboardID, UserID, rs.PersonaID, rs.Name, rs.Role, rs.Description)
	if err != nil {
		return nil, err, false
	}
	updatedPersonas, _ := json.Marshal(personas)
	msg := createSocketEvent("personas_updated", string(updatedPersonas), "")

	return msg, nil, false
}

// DeletePersona handles deleting a storyboard persona
func (b *Service) DeletePersona(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	goals, err := b.db.DeleteStoryboardPersona(StoryboardID, UserID, EventValue)
	if err != nil {
		return nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := createSocketEvent("personas_updated", string(updatedGoals), "")

	return msg, nil, false
}

// FacilitatorAdd handles adding a storyboard facilitator
func (b *Service) FacilitatorAdd(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	storyboard, err := b.db.StoryboardFacilitatorAdd(StoryboardID, rs.UserID)
	if err != nil {
		return nil, err, false
	}
	updatedStoryboard, _ := json.Marshal(storyboard)
	msg := createSocketEvent("storyboard_updated", string(updatedStoryboard), "")

	return msg, nil, false
}

// FacilitatorRemove handles removing a storyboard facilitator
func (b *Service) FacilitatorRemove(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	storyboard, err := b.db.StoryboardFacilitatorRemove(StoryboardID, rs.UserID)
	if err != nil {
		return nil, err, false
	}
	updatedStoryboard, _ := json.Marshal(storyboard)
	msg := createSocketEvent("storyboard_updated", string(updatedStoryboard), "")

	return msg, nil, false
}

// FacilitatorSelf handles self-promoting a user to a facilitator
func (b *Service) FacilitatorSelf(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	facilitatorCode, err := b.db.GetStoryboardFacilitatorCode(StoryboardID)
	if err != nil {
		return nil, err, false
	}

	if EventValue == facilitatorCode {
		storyboard, err := b.db.StoryboardFacilitatorAdd(StoryboardID, UserID)
		if err != nil {
			return nil, err, false
		}
		updatedStoryboard, _ := json.Marshal(storyboard)

		msg := createSocketEvent("storyboard_updated", string(updatedStoryboard), "")

		return msg, nil, false
	} else {
		return nil, errors.New("INCORRECT_FACILITATOR_CODE"), false
	}
}

// ReviseColorLegend handles revising a storyboard color legend
func (b *Service) ReviseColorLegend(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	storyboard, err := b.db.StoryboardReviseColorLegend(StoryboardID, UserID, EventValue)
	if err != nil {
		return nil, err, false
	}
	updatedStoryboard, _ := json.Marshal(storyboard)
	msg := createSocketEvent("storyboard_updated", string(updatedStoryboard), "")

	return msg, nil, false
}

// EditStoryboard handles editing the storyboard settings
func (b *Service) EditStoryboard(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rb struct {
		Name            string `json:"storyboardName"`
		JoinCode        string `json:"joinCode"`
		FacilitatorCode string `json:"facilitatorCode"`
	}
	err := json.Unmarshal([]byte(EventValue), &rb)
	if err != nil {
		return nil, err, false
	}

	err = b.db.EditStoryboard(
		StoryboardID,
		rb.Name,
		rb.JoinCode,
		rb.FacilitatorCode,
	)
	if err != nil {
		return nil, err, false
	}

	updatedStoryboard, _ := json.Marshal(rb)
	msg := createSocketEvent("storyboard_edited", string(updatedStoryboard), "")

	return msg, nil, false
}

// Delete handles deleting the storyboard
func (b *Service) Delete(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	err := b.db.DeleteStoryboard(StoryboardID, UserID)
	if err != nil {
		return nil, err, false
	}
	msg := createSocketEvent("storyboard_conceded", "", "")

	return msg, nil, false
}

// Abandon handles setting abandoned true so storyboard doesn't show up in users storyboard list, then leaves storyboard
func (b *Service) Abandon(ctx context.Context, StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	_, err := b.db.AbandonStoryboard(StoryboardID, UserID)
	if err != nil {
		return nil, err, false
	}

	return nil, errors.New("ABANDONED_STORYBOARD"), true
}

// socketEvent is the event structure used for socket messages
type socketEvent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	User  string `json:"userId"`
}

func createSocketEvent(Type string, Value string, User string) []byte {
	newEvent := &socketEvent{
		Type:  Type,
		Value: Value,
		User:  User,
	}

	event, _ := json.Marshal(newEvent)

	return event
}
