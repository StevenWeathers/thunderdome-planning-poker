package db_test

import (
	"context"
	"testing"

	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/retro"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/storyboard"
)

// TestLiveCommentOwnership covers the retro item comments, reactions and storyboard story comments
// that are only changed over websockets.
func TestLiveCommentOwnership(t *testing.T) {
	d := testDB(t)
	ctx := context.Background()
	retros := &retro.Service{DB: d.DB, Logger: d.Logger}
	storyboards := &storyboard.Service{DB: d.DB, Logger: d.Logger, HTMLSanitizerPolicy: bluemonday.UGCPolicy()}
	exec := func(query string, args ...any) {
		t.Helper()
		_, err := d.DB.Exec(query, args...)
		require.NoError(t, err)
	}

	teamID := insertID(t, d, `INSERT INTO thunderdome.team (name) VALUES ('integration') RETURNING id`)
	teamMember := func(role string) string {
		userID := insertUser(t, d)
		exec(`INSERT INTO thunderdome.team_user (team_id, user_id, role) VALUES ($1, $2, $3)`, teamID, userID, role)
		return userID
	}
	author := teamMember("MEMBER")
	owner := insertUser(t, d)

	templateID := insertID(t, d, `SELECT id FROM thunderdome.retro_template LIMIT 1`)
	r, err := retros.CreateRetro(ctx, owner, "", "integration", "", "", 3, "visible",
		0, false, false, false, false, templateID)
	require.NoError(t, err)
	exec(`UPDATE thunderdome.retro SET team_id = $1 WHERE id = $2`, teamID, r.ID)
	_, err = retros.CreateRetroItem(r.ID, author, "worked", "integration")
	require.NoError(t, err)
	itemID := insertID(t, d, `SELECT id FROM thunderdome.retro_item WHERE retro_id = $1`, r.ID)

	boardID := insertID(t, d, `INSERT INTO thunderdome.storyboard (owner_id, name, team_id)
		VALUES ($1, 'integration', $2) RETURNING id`, owner, teamID)
	goalID := insertID(t, d, `INSERT INTO thunderdome.storyboard_goal (storyboard_id, name, display_order)
		VALUES ($1, 'integration', 'a') RETURNING id`, boardID)
	columnID := insertID(t, d, `INSERT INTO thunderdome.storyboard_column (storyboard_id, goal_id, name, display_order)
		VALUES ($1, $2, 'integration', 'a') RETURNING id`, boardID, goalID)
	storyID := insertID(t, d, `INSERT INTO thunderdome.storyboard_story (storyboard_id, goal_id, column_id, display_order)
		VALUES ($1, $2, $3, 'a') RETURNING id`, boardID, goalID, columnID)

	facilitator := insertUser(t, d)
	exec(`INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES ($1, $2)`, r.ID, facilitator)
	exec(`INSERT INTO thunderdome.storyboard_facilitator (storyboard_id, user_id) VALUES ($1, $2)`, boardID, facilitator)

	tests := []struct {
		name    string
		userID  string
		allowed bool // moderate storyboard story comments; retro content only ever belongs to its author
	}{
		{name: "author", userID: author, allowed: true},
		{name: "teammate", userID: teamMember("MEMBER")},
		{name: "outsider", userID: insertUser(t, d)},
		{name: "facilitator", userID: facilitator, allowed: true},
		{name: "team admin", userID: teamMember("ADMIN"), allowed: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := func(allowed bool, err error, refused string) {
				t.Helper()
				if allowed {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, refused)
				}
			}

			_, err := retros.ItemCommentAdd(r.ID, itemID, author, "original")
			require.NoError(t, err)
			itemCommentID := insertID(t, d, `SELECT id FROM thunderdome.retro_item_comment
				WHERE item_id = $1 ORDER BY created_date DESC LIMIT 1`, itemID)
			isAuthor := tt.userID == author
			_, err = retros.ItemCommentEdit(r.ID, tt.userID, itemCommentID, "edited")
			check(isAuthor, err, "COMMENT_CHANGE_NOT_ALLOWED")
			_, err = retros.ItemCommentDelete(r.ID, tt.userID, itemCommentID)
			check(isAuthor, err, "COMMENT_CHANGE_NOT_ALLOWED")

			_, err = storyboards.AddStoryComment(boardID, author, storyID, "original")
			require.NoError(t, err)
			storyCommentID := insertID(t, d, `SELECT id FROM thunderdome.storyboard_story_comment
				WHERE story_id = $1 ORDER BY created_date DESC LIMIT 1`, storyID)
			_, err = storyboards.EditStoryComment(boardID, tt.userID, storyCommentID, "edited")
			check(tt.allowed, err, "COMMENT_CHANGE_NOT_ALLOWED")
			_, err = storyboards.DeleteStoryComment(boardID, tt.userID, storyCommentID)
			check(tt.allowed, err, "COMMENT_CHANGE_NOT_ALLOWED")

			// reactions and feedback items belong to whoever added them, nobody moderates them
			exec(`DELETE FROM thunderdome.retro_item_reaction WHERE item_id = $1`, itemID)
			_, err = retros.ItemReactionAdd(r.ID, itemID, author, "👍")
			require.NoError(t, err)
			reactionID := insertID(t, d, `SELECT id FROM thunderdome.retro_item_reaction WHERE item_id = $1`, itemID)
			_, err = retros.ItemReactionDelete(r.ID, tt.userID, reactionID)
			check(isAuthor, err, "REACTION_CHANGE_NOT_ALLOWED")

			ownItemID := insertID(t, d, `INSERT INTO thunderdome.retro_item (retro_id, type, content, user_id)
				VALUES ($1, 'worked', 'deletable', $2) RETURNING id`, r.ID, author)
			_, err = retros.DeleteRetroItem(r.ID, tt.userID, "worked", ownItemID)
			check(isAuthor, err, "ITEM_CHANGE_NOT_ALLOWED")
			remaining := 1
			if isAuthor {
				remaining = 0
			}
			assert.Equal(t, remaining, count(t, d, `SELECT count(*) FROM thunderdome.retro_item WHERE id = $1`, ownItemID))
		})
	}
}
