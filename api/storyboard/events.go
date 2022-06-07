package storyboard

import (
	"encoding/json"
	"errors"
)

// Delete handles deleting the storyboard
func (b *Service) Delete(StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	err := b.db.DeleteStoryboard(StoryboardID, UserID)
	if err != nil {
		return nil, err, false
	}
	msg := createSocketEvent("conceded", "", "")

	return msg, nil, false
}

// Abandon handles setting abandoned true so storyboard doesn't show up in users storyboard list, then leaves storyboard
func (b *Service) Abandon(StoryboardID string, UserID string, EventValue string) ([]byte, error, bool) {
	b.db.AbandonStoryboard(StoryboardID, UserID)

	return nil, errors.New("ABANDONED_RETRO"), true
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
