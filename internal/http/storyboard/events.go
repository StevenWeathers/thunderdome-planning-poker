package storyboard

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
)

// AddGoal handles adding a goal to storyboard
func (b *Service) AddGoal(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goal, err := b.StoryboardService.CreateStoryboardGoal(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	goals := b.StoryboardService.GetStoryboardGoals(storyboardID)
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("goal_added", string(updatedGoals), "")

	return goal, msg, nil, false
}

// ReviseGoal handles revising a storyboard goal
func (b *Service) ReviseGoal(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	goalID := goalObj["goalId"]
	goalName := goalObj["name"]

	goals, err := b.StoryboardService.ReviseGoalName(storyboardID, userID, goalID, goalName)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("goal_revised", string(updatedGoals), "")

	return nil, msg, nil, false
}

// DeleteGoal handles deleting a storyboard goal
func (b *Service) DeleteGoal(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goals, err := b.StoryboardService.DeleteStoryboardGoal(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("goal_deleted", string(updatedGoals), "")

	return nil, msg, nil, false
}

// AddColumn handles adding a column to storyboard goal
func (b *Service) AddColumn(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	goalID := goalObj["goalId"]

	column, err := b.StoryboardService.CreateStoryboardColumn(storyboardID, goalID, userID)
	if err != nil {
		return nil, nil, err, false
	}
	goals := b.StoryboardService.GetStoryboardGoals(storyboardID)
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("column_added", string(updatedGoals), "")

	return column, msg, nil, false
}

// ReviseColumn handles revising a storyboard goal column
func (b *Service) ReviseColumn(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		ColumnID string `json:"id"`
		Name     string `json:"name"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := b.StoryboardService.ReviseStoryboardColumn(storyboardID, userID, rs.ColumnID, rs.Name)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("column_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// DeleteColumn handles deleting a storyboard goal column
func (b *Service) DeleteColumn(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goals, err := b.StoryboardService.DeleteStoryboardColumn(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_deleted", string(updatedGoals), "")

	return nil, msg, nil, false
}

// MoveColumn handles moving a storyboard column between goals
func (b *Service) MoveColumn(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var moveColumnInput struct {
		ColumnID    string `json:"columnId"`
		GoalID      string `json:"goalId"`
		PlaceBefore string `json:"placeBefore"`
	}
	err := json.Unmarshal([]byte(eventValue), &moveColumnInput)
	if err != nil {
		return nil, nil, err, false
	}

	err = b.StoryboardService.MoveStoryboardColumn(storyboardID, userID, moveColumnInput.ColumnID, moveColumnInput.GoalID, moveColumnInput.PlaceBefore)
	if err != nil {
		return nil, nil, err, false
	}

	goal, err := b.StoryboardService.GetStoryboardGoal(storyboardID, moveColumnInput.GoalID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedGoal, _ := json.Marshal(goal)
	msg := wshub.CreateSocketEvent("column_moved", string(updatedGoal), "")

	return nil, msg, nil, false
}

// ColumnPersonaAdd handles adding a persona to a storyboard goal column
func (b *Service) ColumnPersonaAdd(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		ColumnID  string `json:"column_id"`
		PersonaID string `json:"persona_id"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := b.StoryboardService.ColumnPersonaAdd(storyboardID, rs.ColumnID, rs.PersonaID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("column_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// ColumnPersonaRemove handles removing a persona from a storyboard goal column
func (b *Service) ColumnPersonaRemove(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		ColumnID  string `json:"column_id"`
		PersonaID string `json:"persona_id"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := b.StoryboardService.ColumnPersonaRemove(storyboardID, rs.ColumnID, rs.PersonaID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("column_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// AddStory handles adding a story to storyboard
func (b *Service) AddStory(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var ns struct {
		GoalID   string `json:"goalId"`
		ColumnID string `json:"columnId"`
	}
	err := json.Unmarshal([]byte(eventValue), &ns)
	if err != nil {
		return nil, nil, err, false
	}

	story, err := b.StoryboardService.CreateStoryboardStory(storyboardID, ns.GoalID, ns.ColumnID, userID)
	if err != nil {
		return nil, nil, err, false
	}
	goals := b.StoryboardService.GetStoryboardGoals(storyboardID)
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_added", string(updatedGoals), "")

	return story, msg, nil, false
}

// UpdateStoryName handles revising a storyboard story name
func (b *Service) UpdateStoryName(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	storyID := goalObj["storyId"]
	storyName := goalObj["name"]

	goals, err := b.StoryboardService.ReviseStoryName(storyboardID, userID, storyID, storyName)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// UpdateStoryContent handles revising a storyboard story content
func (b *Service) UpdateStoryContent(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	storyID := goalObj["storyId"]
	storyContent := goalObj["content"]

	goals, err := b.StoryboardService.ReviseStoryContent(storyboardID, userID, storyID, storyContent)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// UpdateStoryColor handles revising a storyboard story color
func (b *Service) UpdateStoryColor(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	storyID := goalObj["storyId"]
	storyColor := goalObj["color"]

	goals, err := b.StoryboardService.ReviseStoryColor(storyboardID, userID, storyID, storyColor)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// UpdateStoryPoints handles revising a storyboard story points
func (b *Service) UpdateStoryPoints(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		StoryID string `json:"storyId"`
		Points  int    `json:"points"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := b.StoryboardService.ReviseStoryPoints(storyboardID, userID, rs.StoryID, rs.Points)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// UpdateStoryClosed handles revising a storyboard story closed status
func (b *Service) UpdateStoryClosed(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		StoryID string `json:"storyId"`
		Closed  bool   `json:"closed"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := b.StoryboardService.ReviseStoryClosed(storyboardID, userID, rs.StoryID, rs.Closed)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// UpdateStoryLink handles revising a storyboard story link
func (b *Service) UpdateStoryLink(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	storyID := goalObj["storyId"]
	link := goalObj["link"]

	goals, err := b.StoryboardService.ReviseStoryLink(storyboardID, userID, storyID, link)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// MoveStory handles moving a storyboard story between columns/goals
func (b *Service) MoveStory(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	storyID := goalObj["storyId"]
	goalID := goalObj["goalId"]
	columnID := goalObj["columnId"]
	placeBefore := goalObj["placeBefore"]

	goals, err := b.StoryboardService.MoveStoryboardStory(storyboardID, userID, storyID, goalID, columnID, placeBefore)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_moved", string(updatedGoals), "")

	return nil, msg, nil, false
}

// DeleteStory handles deleting a storyboard story
func (b *Service) DeleteStory(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goals, err := b.StoryboardService.DeleteStoryboardStory(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_deleted", string(updatedGoals), "")

	return nil, msg, nil, false
}

// AddStoryComment handles adding a storyboard story comment
func (b *Service) AddStoryComment(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		StoryID string `json:"storyId"`
		Comment string `json:"comment"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := b.StoryboardService.AddStoryComment(storyboardID, userID, rs.StoryID, rs.Comment)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// EditStoryComment handles editing a storyboard story comment
func (b *Service) EditStoryComment(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		CommentID string `json:"commentId"`
		Comment   string `json:"comment"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := b.StoryboardService.EditStoryComment(storyboardID, rs.CommentID, rs.Comment)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// DeleteStoryComment handles deleting a storyboard story comment
func (b *Service) DeleteStoryComment(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		CommentID string `json:"commentId"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := b.StoryboardService.DeleteStoryComment(storyboardID, rs.CommentID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// AddPersona handles adding a storyboard persona
func (b *Service) AddPersona(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		Name        string `json:"name"`
		Role        string `json:"role"`
		Description string `json:"description"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	personas, err := b.StoryboardService.AddStoryboardPersona(storyboardID, userID, rs.Name, rs.Role, rs.Description)
	if err != nil {
		return nil, nil, err, false
	}
	updatedPersonas, _ := json.Marshal(personas)
	msg := wshub.CreateSocketEvent("personas_updated", string(updatedPersonas), "")

	return nil, msg, nil, false
}

// UpdatePersona handles updating a storyboard persona
func (b *Service) UpdatePersona(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		PersonaID   string `json:"id"`
		Name        string `json:"name"`
		Role        string `json:"role"`
		Description string `json:"description"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	personas, err := b.StoryboardService.UpdateStoryboardPersona(storyboardID, userID, rs.PersonaID, rs.Name, rs.Role, rs.Description)
	if err != nil {
		return nil, nil, err, false
	}
	updatedPersonas, _ := json.Marshal(personas)
	msg := wshub.CreateSocketEvent("personas_updated", string(updatedPersonas), "")

	return nil, msg, nil, false
}

// DeletePersona handles deleting a storyboard persona
func (b *Service) DeletePersona(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goals, err := b.StoryboardService.DeleteStoryboardPersona(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("personas_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// FacilitatorAdd handles adding a storyboard facilitator
func (b *Service) FacilitatorAdd(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	storyboard, err := b.StoryboardService.StoryboardFacilitatorAdd(storyboardID, rs.UserID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedStoryboard, _ := json.Marshal(storyboard)
	msg := wshub.CreateSocketEvent("storyboard_updated", string(updatedStoryboard), "")

	return nil, msg, nil, false
}

// FacilitatorRemove handles removing a storyboard facilitator
func (b *Service) FacilitatorRemove(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	storyboard, err := b.StoryboardService.StoryboardFacilitatorRemove(storyboardID, rs.UserID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedStoryboard, _ := json.Marshal(storyboard)
	msg := wshub.CreateSocketEvent("storyboard_updated", string(updatedStoryboard), "")

	return nil, msg, nil, false
}

// FacilitatorSelf handles self-promoting a user to a facilitator
func (b *Service) FacilitatorSelf(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	facilitatorCode, err := b.StoryboardService.GetStoryboardFacilitatorCode(storyboardID)
	if err != nil {
		return nil, nil, err, false
	}

	if eventValue == facilitatorCode {
		storyboard, err := b.StoryboardService.StoryboardFacilitatorAdd(storyboardID, userID)
		if err != nil {
			return nil, nil, err, false
		}
		updatedStoryboard, _ := json.Marshal(storyboard)

		msg := wshub.CreateSocketEvent("storyboard_updated", string(updatedStoryboard), "")

		return nil, msg, nil, false
	} else {
		return nil, nil, errors.New("INCORRECT_FACILITATOR_CODE"), false
	}
}

// ReviseColorLegend handles revising a storyboard color legend
func (b *Service) ReviseColorLegend(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	storyboard, err := b.StoryboardService.StoryboardReviseColorLegend(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	updatedStoryboard, _ := json.Marshal(storyboard)
	msg := wshub.CreateSocketEvent("storyboard_updated", string(updatedStoryboard), "")

	return nil, msg, nil, false
}

// EditStoryboard handles editing the storyboard settings
func (b *Service) EditStoryboard(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rb struct {
		Name            string `json:"storyboardName"`
		JoinCode        string `json:"joinCode"`
		FacilitatorCode string `json:"facilitatorCode"`
	}
	err := json.Unmarshal([]byte(eventValue), &rb)
	if err != nil {
		return nil, nil, err, false
	}

	err = b.StoryboardService.EditStoryboard(
		storyboardID,
		rb.Name,
		rb.JoinCode,
		rb.FacilitatorCode,
	)
	if err != nil {
		return nil, nil, err, false
	}

	updatedStoryboard, _ := json.Marshal(rb)
	msg := wshub.CreateSocketEvent("storyboard_edited", string(updatedStoryboard), "")

	return nil, msg, nil, false
}

// Delete handles deleting the storyboard
func (b *Service) Delete(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	err := b.StoryboardService.DeleteStoryboard(storyboardID, userID)
	if err != nil {
		return nil, nil, err, false
	}
	msg := wshub.CreateSocketEvent("storyboard_conceded", "", "")

	return nil, msg, nil, false
}

// Abandon handles setting abandoned true so storyboard doesn't show up in users storyboard list, then leaves storyboard
func (b *Service) Abandon(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	_, err := b.StoryboardService.AbandonStoryboard(storyboardID, userID)
	if err != nil {
		return nil, nil, err, false
	}

	return nil, nil, errors.New("ABANDONED_STORYBOARD"), true
}
