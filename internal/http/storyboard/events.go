package storyboard

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
)

// AddGoal handles adding a goal to storyboard
func (s *Service) AddGoal(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goal, err := s.StoryboardService.CreateStoryboardGoal(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	goals := s.StoryboardService.GetStoryboardGoals(storyboardID)
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("goal_added", string(updatedGoals), "")

	return goal, msg, nil, false
}

// ReviseGoal handles revising a storyboard goal
func (s *Service) ReviseGoal(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	goalID := goalObj["goalId"]
	goalName := goalObj["name"]

	goals, err := s.StoryboardService.ReviseGoalName(storyboardID, userID, goalID, goalName)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("goal_revised", string(updatedGoals), "")

	return nil, msg, nil, false
}

// DeleteGoal handles deleting a storyboard goal
func (s *Service) DeleteGoal(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goals, err := s.StoryboardService.DeleteStoryboardGoal(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("goal_deleted", string(updatedGoals), "")

	return nil, msg, nil, false
}

// AddColumn handles adding a column to storyboard goal
func (s *Service) AddColumn(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	goalID := goalObj["goalId"]

	column, err := s.StoryboardService.CreateStoryboardColumn(storyboardID, goalID, userID)
	if err != nil {
		return nil, nil, err, false
	}
	goals := s.StoryboardService.GetStoryboardGoals(storyboardID)
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("column_added", string(updatedGoals), "")

	return column, msg, nil, false
}

// ReviseColumn handles revising a storyboard goal column
func (s *Service) ReviseColumn(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		ColumnID string `json:"id"`
		Name     string `json:"name"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := s.StoryboardService.ReviseStoryboardColumn(storyboardID, userID, rs.ColumnID, rs.Name)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("column_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// DeleteColumn handles deleting a storyboard goal column
func (s *Service) DeleteColumn(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goals, err := s.StoryboardService.DeleteStoryboardColumn(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_deleted", string(updatedGoals), "")

	return nil, msg, nil, false
}

// MoveColumn handles moving a storyboard column between goals
func (s *Service) MoveColumn(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var moveColumnInput struct {
		ColumnID    string `json:"columnId"`
		GoalID      string `json:"goalId"`
		PlaceBefore string `json:"placeBefore"`
	}
	err := json.Unmarshal([]byte(eventValue), &moveColumnInput)
	if err != nil {
		return nil, nil, err, false
	}

	err = s.StoryboardService.MoveStoryboardColumn(storyboardID, userID, moveColumnInput.ColumnID, moveColumnInput.GoalID, moveColumnInput.PlaceBefore)
	if err != nil {
		return nil, nil, err, false
	}

	goal, err := s.StoryboardService.GetStoryboardGoal(storyboardID, moveColumnInput.GoalID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedGoal, _ := json.Marshal(goal)
	msg := wshub.CreateSocketEvent("column_moved", string(updatedGoal), "")

	return nil, msg, nil, false
}

// ColumnPersonaAdd handles adding a persona to a storyboard goal column
func (s *Service) ColumnPersonaAdd(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		ColumnID  string `json:"column_id"`
		PersonaID string `json:"persona_id"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := s.StoryboardService.ColumnPersonaAdd(storyboardID, rs.ColumnID, rs.PersonaID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("column_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// ColumnPersonaRemove handles removing a persona from a storyboard goal column
func (s *Service) ColumnPersonaRemove(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		ColumnID  string `json:"column_id"`
		PersonaID string `json:"persona_id"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := s.StoryboardService.ColumnPersonaRemove(storyboardID, rs.ColumnID, rs.PersonaID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("column_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// AddStory handles adding a story to storyboard
func (s *Service) AddStory(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var ns struct {
		GoalID   string `json:"goalId"`
		ColumnID string `json:"columnId"`
	}
	err := json.Unmarshal([]byte(eventValue), &ns)
	if err != nil {
		return nil, nil, err, false
	}

	story, err := s.StoryboardService.CreateStoryboardStory(storyboardID, ns.GoalID, ns.ColumnID, userID)
	if err != nil {
		return nil, nil, err, false
	}
	goals := s.StoryboardService.GetStoryboardGoals(storyboardID)
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_added", string(updatedGoals), "")

	return story, msg, nil, false
}

// UpdateStoryName handles revising a storyboard story name
func (s *Service) UpdateStoryName(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	storyID := goalObj["storyId"]
	storyName := goalObj["name"]

	goals, err := s.StoryboardService.ReviseStoryName(storyboardID, userID, storyID, storyName)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// UpdateStoryContent handles revising a storyboard story content
func (s *Service) UpdateStoryContent(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	storyID := goalObj["storyId"]
	storyContent := goalObj["content"]

	goals, err := s.StoryboardService.ReviseStoryContent(storyboardID, userID, storyID, storyContent)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// UpdateStoryColor handles revising a storyboard story color
func (s *Service) UpdateStoryColor(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	storyID := goalObj["storyId"]
	storyColor := goalObj["color"]

	goals, err := s.StoryboardService.ReviseStoryColor(storyboardID, userID, storyID, storyColor)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// UpdateStoryPoints handles revising a storyboard story points
func (s *Service) UpdateStoryPoints(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		StoryID string `json:"storyId"`
		Points  int    `json:"points"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := s.StoryboardService.ReviseStoryPoints(storyboardID, userID, rs.StoryID, rs.Points)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// UpdateStoryClosed handles revising a storyboard story closed status
func (s *Service) UpdateStoryClosed(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		StoryID string `json:"storyId"`
		Closed  bool   `json:"closed"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := s.StoryboardService.ReviseStoryClosed(storyboardID, userID, rs.StoryID, rs.Closed)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// UpdateStoryLink handles revising a storyboard story link
func (s *Service) UpdateStoryLink(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	storyID := goalObj["storyId"]
	link := goalObj["link"]

	goals, err := s.StoryboardService.ReviseStoryLink(storyboardID, userID, storyID, link)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// MoveStory handles moving a storyboard story between columns/goals
func (s *Service) MoveStory(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goalObj := make(map[string]string)
	err := json.Unmarshal([]byte(eventValue), &goalObj)
	if err != nil {
		return nil, nil, err, false
	}
	storyID := goalObj["storyId"]
	goalID := goalObj["goalId"]
	columnID := goalObj["columnId"]
	placeBefore := goalObj["placeBefore"]

	goals, err := s.StoryboardService.MoveStoryboardStory(storyboardID, userID, storyID, goalID, columnID, placeBefore)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_moved", string(updatedGoals), "")

	return nil, msg, nil, false
}

// DeleteStory handles deleting a storyboard story
func (s *Service) DeleteStory(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goals, err := s.StoryboardService.DeleteStoryboardStory(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_deleted", string(updatedGoals), "")

	return nil, msg, nil, false
}

// AddStoryComment handles adding a storyboard story comment
func (s *Service) AddStoryComment(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		StoryID string `json:"storyId"`
		Comment string `json:"comment"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := s.StoryboardService.AddStoryComment(storyboardID, userID, rs.StoryID, rs.Comment)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// EditStoryComment handles editing a storyboard story comment
func (s *Service) EditStoryComment(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		CommentID string `json:"commentId"`
		Comment   string `json:"comment"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := s.StoryboardService.EditStoryComment(storyboardID, rs.CommentID, rs.Comment)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// DeleteStoryComment handles deleting a storyboard story comment
func (s *Service) DeleteStoryComment(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		CommentID string `json:"commentId"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	goals, err := s.StoryboardService.DeleteStoryComment(storyboardID, rs.CommentID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("story_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// AddPersona handles adding a storyboard persona
func (s *Service) AddPersona(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		Name        string `json:"name"`
		Role        string `json:"role"`
		Description string `json:"description"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	personas, err := s.StoryboardService.AddStoryboardPersona(storyboardID, userID, rs.Name, rs.Role, rs.Description)
	if err != nil {
		return nil, nil, err, false
	}
	updatedPersonas, _ := json.Marshal(personas)
	msg := wshub.CreateSocketEvent("personas_updated", string(updatedPersonas), "")

	return nil, msg, nil, false
}

// UpdatePersona handles updating a storyboard persona
func (s *Service) UpdatePersona(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
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

	personas, err := s.StoryboardService.UpdateStoryboardPersona(storyboardID, userID, rs.PersonaID, rs.Name, rs.Role, rs.Description)
	if err != nil {
		return nil, nil, err, false
	}
	updatedPersonas, _ := json.Marshal(personas)
	msg := wshub.CreateSocketEvent("personas_updated", string(updatedPersonas), "")

	return nil, msg, nil, false
}

// DeletePersona handles deleting a storyboard persona
func (s *Service) DeletePersona(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	goals, err := s.StoryboardService.DeleteStoryboardPersona(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	updatedGoals, _ := json.Marshal(goals)
	msg := wshub.CreateSocketEvent("personas_updated", string(updatedGoals), "")

	return nil, msg, nil, false
}

// FacilitatorAdd handles adding a storyboard facilitator
func (s *Service) FacilitatorAdd(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	storyboard, err := s.StoryboardService.StoryboardFacilitatorAdd(storyboardID, rs.UserID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedStoryboard, _ := json.Marshal(storyboard)
	msg := wshub.CreateSocketEvent("storyboard_updated", string(updatedStoryboard), "")

	return nil, msg, nil, false
}

// FacilitatorRemove handles removing a storyboard facilitator
func (s *Service) FacilitatorRemove(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	err := json.Unmarshal([]byte(eventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	storyboard, err := s.StoryboardService.StoryboardFacilitatorRemove(storyboardID, rs.UserID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedStoryboard, _ := json.Marshal(storyboard)
	msg := wshub.CreateSocketEvent("storyboard_updated", string(updatedStoryboard), "")

	return nil, msg, nil, false
}

// FacilitatorSelf handles self-promoting a user to a facilitator
func (s *Service) FacilitatorSelf(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	facilitatorCode, err := s.StoryboardService.GetStoryboardFacilitatorCode(storyboardID)
	if err != nil {
		return nil, nil, err, false
	}

	if eventValue == facilitatorCode {
		storyboard, err := s.StoryboardService.StoryboardFacilitatorAdd(storyboardID, userID)
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
func (s *Service) ReviseColorLegend(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	storyboard, err := s.StoryboardService.StoryboardReviseColorLegend(storyboardID, userID, eventValue)
	if err != nil {
		return nil, nil, err, false
	}
	updatedStoryboard, _ := json.Marshal(storyboard)
	msg := wshub.CreateSocketEvent("storyboard_updated", string(updatedStoryboard), "")

	return nil, msg, nil, false
}

// EditStoryboard handles editing the storyboard settings
func (s *Service) EditStoryboard(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	var rb struct {
		Name            string `json:"storyboardName"`
		JoinCode        string `json:"joinCode"`
		FacilitatorCode string `json:"facilitatorCode"`
	}
	err := json.Unmarshal([]byte(eventValue), &rb)
	if err != nil {
		return nil, nil, err, false
	}

	err = s.StoryboardService.EditStoryboard(
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
func (s *Service) Delete(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	err := s.StoryboardService.DeleteStoryboard(storyboardID, userID)
	if err != nil {
		return nil, nil, err, false
	}
	msg := wshub.CreateSocketEvent("storyboard_conceded", "", "")

	return nil, msg, nil, false
}

// Abandon handles setting abandoned true so storyboard doesn't show up in users storyboard list, then leaves storyboard
func (s *Service) Abandon(ctx context.Context, storyboardID string, userID string, eventValue string) (any, []byte, error, bool) {
	_, err := s.StoryboardService.AbandonStoryboard(storyboardID, userID)
	if err != nil {
		return nil, nil, err, false
	}

	return nil, nil, errors.New("ABANDONED_STORYBOARD"), true
}
