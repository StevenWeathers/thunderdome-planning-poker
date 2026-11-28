package db_test

import (
	"context"
	"testing"

	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/jira"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/poker"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/retro"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/team"
)

func count(t *testing.T, d *db.Service, query string, args ...any) int {
	t.Helper()

	var n int
	require.NoError(t, d.DB.QueryRow(query, args...).Scan(&n))

	return n
}

// TestWritesAreScopedToParent checks that a child ID combined with another parent's ID changes nothing.
func TestWritesAreScopedToParent(t *testing.T) {
	d := testDB(t)
	ctx := context.Background()
	userID := insertUser(t, d)

	newTeam := func() string {
		id := insertID(t, d, `INSERT INTO thunderdome.team (name) VALUES ('integration') RETURNING id`)
		_, err := d.DB.Exec(`INSERT INTO thunderdome.team_user (team_id, user_id, role) VALUES ($1, $2, 'MEMBER')`, id, userID)
		require.NoError(t, err)
		return id
	}
	teamID, otherTeamID := newTeam(), newTeam()

	t.Run("check-ins", func(t *testing.T) {
		svc := &team.CheckinService{DB: d.DB, Logger: d.Logger, HTMLSanitizerPolicy: bluemonday.UGCPolicy()}
		require.NoError(t, svc.CheckinCreate(ctx, teamID, userID, "", "y", "mine", "", "", true))
		checkinID := insertID(t, d, `SELECT id FROM thunderdome.team_checkin WHERE team_id = $1`, teamID)
		today := func() int {
			return count(t, d, `SELECT count(*) FROM thunderdome.team_checkin WHERE id = $1 AND today = 'mine'`, checkinID)
		}

		assert.Error(t, svc.CheckinUpdate(ctx, otherTeamID, userID, checkinID, "y", "theirs", "", "", true))
		assert.Equal(t, 1, today())
		assert.Error(t, svc.CheckinComment(ctx, otherTeamID, checkinID, userID, "theirs"))
		require.NoError(t, svc.CheckinComment(ctx, teamID, checkinID, userID, "mine"))
		commentID := insertID(t, d, `SELECT id FROM thunderdome.team_checkin_comment WHERE checkin_id = $1`, checkinID)
		assert.Error(t, svc.CheckinCommentEdit(ctx, otherTeamID, userID, commentID, "theirs"))
		assert.NoError(t, svc.CheckinCommentEdit(ctx, teamID, userID, commentID, "edited"))
		assert.Error(t, svc.CheckinCommentDelete(ctx, otherTeamID, userID, commentID))
		assert.NoError(t, svc.CheckinCommentDelete(ctx, teamID, userID, commentID))
		assert.Error(t, svc.CheckinDelete(ctx, otherTeamID, userID, checkinID))
		assert.Equal(t, 1, today())
		assert.NoError(t, svc.CheckinUpdate(ctx, teamID, userID, checkinID, "y", "updated", "", "", true))
		assert.NoError(t, svc.CheckinDelete(ctx, teamID, userID, checkinID))
	})

	orgs := &team.OrganizationService{DB: d.DB, Logger: d.Logger}
	orgID := insertID(t, d, `INSERT INTO thunderdome.organization (name) VALUES ('mine') RETURNING id`)
	otherOrgID := insertID(t, d, `INSERT INTO thunderdome.organization (name) VALUES ('other') RETURNING id`)
	deptID := insertID(t, d, `INSERT INTO thunderdome.organization_department (name, organization_id)
		VALUES ('mine', $1) RETURNING id`, orgID)
	otherDeptID := insertID(t, d, `INSERT INTO thunderdome.organization_department (name, organization_id)
		VALUES ('other', $1) RETURNING id`, otherOrgID)

	t.Run("invites", func(t *testing.T) {
		teams := &team.Service{DB: d.DB, Logger: d.Logger}
		tests := []struct {
			name   string
			parent string
			other  string
			table  string
			invite func(string) (string, error)
			remove func(string, string) error
		}{
			{"team", teamID, otherTeamID, "team_user_invite",
				func(id string) (string, error) { return teams.TeamInviteUser(ctx, id, "a@example.com", "MEMBER") },
				func(parent, id string) error { return teams.TeamDeleteUserInvite(ctx, parent, id) }},
			{"organization", orgID, otherOrgID, "organization_user_invite",
				func(id string) (string, error) {
					return orgs.OrganizationInviteUser(ctx, id, "a@example.com", "MEMBER")
				},
				func(parent, id string) error { return orgs.OrganizationDeleteUserInvite(ctx, parent, id) }},
			{"department", deptID, otherDeptID, "department_user_invite",
				func(id string) (string, error) { return orgs.DepartmentInviteUser(ctx, id, "a@example.com", "MEMBER") },
				func(parent, id string) error { return orgs.DepartmentDeleteUserInvite(ctx, parent, id) }},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				inviteID, err := tt.invite(tt.parent)
				require.NoError(t, err)
				remaining := `SELECT count(*) FROM thunderdome.` + tt.table + ` WHERE invite_id = $1`

				require.NoError(t, tt.remove(tt.other, inviteID))
				assert.Equal(t, 1, count(t, d, remaining, inviteID))
				require.NoError(t, tt.remove(tt.parent, inviteID))
				assert.Equal(t, 0, count(t, d, remaining, inviteID))
			})
		}
	})

	t.Run("department members must belong to the organization", func(t *testing.T) {
		_, err := orgs.OrganizationAddUser(ctx, otherOrgID, userID, "MEMBER")
		require.NoError(t, err)

		_, err = orgs.DepartmentAddUser(ctx, deptID, userID, "MEMBER")
		assert.Error(t, err)
		_, err = orgs.DepartmentAddUser(ctx, otherDeptID, userID, "MEMBER")
		assert.NoError(t, err)
	})

	t.Run("jira instances", func(t *testing.T) {
		svc := &jira.Service{DB: d.DB, Logger: d.Logger, AESHashKey: "integration"}
		instance, err := svc.CreateInstance(ctx, userID, "https://example.com", "a@example.com", "token", false)
		require.NoError(t, err)
		otherUserID := insertUser(t, d)

		_, err = svc.UpdateInstance(ctx, otherUserID, instance.ID, "https://evil.example.com", "b@example.com", "x")
		assert.Error(t, err)
		assert.Error(t, svc.DeleteInstance(ctx, otherUserID, instance.ID))
		_, err = svc.UpdateInstance(ctx, userID, instance.ID, "https://example.org", "a@example.com", "token")
		assert.NoError(t, err)
		assert.NoError(t, svc.DeleteInstance(ctx, userID, instance.ID))
	})

	t.Run("retros", func(t *testing.T) {
		svc := &retro.Service{DB: d.DB, Logger: d.Logger}
		templateID := insertID(t, d, `SELECT id FROM thunderdome.retro_template LIMIT 1`)
		newRetro := func() string {
			r, err := svc.CreateRetro(ctx, userID, "", "integration", "", "", 3, "visible",
				0, false, false, false, false, templateID)
			require.NoError(t, err)
			return r.ID
		}
		retroID, otherRetroID := newRetro(), newRetro()

		_, err := svc.CreateRetroItem(retroID, userID, "worked", "mine")
		require.NoError(t, err)
		itemID := insertID(t, d, `SELECT id FROM thunderdome.retro_item WHERE retro_id = $1`, retroID)
		_, err = svc.CreateRetroAction(retroID, userID, "mine")
		require.NoError(t, err)
		actionID := insertID(t, d, `SELECT id FROM thunderdome.retro_action WHERE retro_id = $1`, retroID)

		// each step runs against the wrong retro first and must change nothing, then against the right one
		step := func(name string, query string, id string, want int, run func(retroID string)) {
			t.Run(name, func(t *testing.T) {
				before := count(t, d, query, id)
				run(otherRetroID)
				assert.Equal(t, before, count(t, d, query, id), "changed through another retro")
				run(retroID)
				assert.Equal(t, want, count(t, d, query, id))
			})
		}
		itemComments := `SELECT count(*) FROM thunderdome.retro_item_comment WHERE item_id = $1`
		step("item comment add", itemComments, itemID, 1, func(r string) { _, _ = svc.ItemCommentAdd(r, itemID, userID, "mine") })
		itemCommentID := insertID(t, d, `SELECT id FROM thunderdome.retro_item_comment WHERE item_id = $1`, itemID)
		step("item comment edit", `SELECT count(*) FROM thunderdome.retro_item_comment WHERE id = $1 AND comment = 'edited'`,
			itemCommentID, 1, func(r string) { _, _ = svc.ItemCommentEdit(r, userID, itemCommentID, "edited") })
		step("item comment delete", itemComments, itemID, 0, func(r string) { _, _ = svc.ItemCommentDelete(r, userID, itemCommentID) })

		reactions := `SELECT count(*) FROM thunderdome.retro_item_reaction WHERE item_id = $1`
		step("reaction add", reactions, itemID, 1, func(r string) { _, _ = svc.ItemReactionAdd(r, itemID, userID, "👍") })
		reactionID := insertID(t, d, `SELECT id FROM thunderdome.retro_item_reaction WHERE item_id = $1`, itemID)
		step("reaction delete", reactions, itemID, 0, func(r string) { _, _ = svc.ItemReactionDelete(r, userID, reactionID) })

		step("action update", `SELECT count(*) FROM thunderdome.retro_action WHERE id = $1 AND completed`,
			actionID, 1, func(r string) { _, _ = svc.UpdateRetroAction(r, actionID, "mine", true) })
		assignees := `SELECT count(*) FROM thunderdome.retro_action_assignee WHERE action_id = $1`
		step("assignee add", assignees, actionID, 1, func(r string) { _, _ = svc.RetroActionAssigneeAdd(r, actionID, userID) })
		step("assignee delete", assignees, actionID, 0, func(r string) { _, _ = svc.RetroActionAssigneeDelete(r, actionID, userID) })

		comments := `SELECT count(*) FROM thunderdome.retro_action_comment WHERE action_id = $1`
		step("action comment add", comments, actionID, 1, func(r string) { _, _ = svc.RetroActionCommentAdd(r, actionID, userID, "mine") })
		commentID := insertID(t, d, `SELECT id FROM thunderdome.retro_action_comment WHERE action_id = $1`, actionID)
		step("action comment edit", `SELECT count(*) FROM thunderdome.retro_action_comment WHERE id = $1 AND comment = 'edited'`,
			commentID, 1, func(r string) { _, _ = svc.RetroActionCommentEdit(r, actionID, commentID, userID, "edited") })
		step("action comment delete", comments, actionID, 0, func(r string) { _, _ = svc.RetroActionCommentDelete(r, actionID, commentID, userID) })

		step("action delete", `SELECT count(*) FROM thunderdome.retro_action WHERE id = $1`,
			actionID, 0, func(r string) { _, _ = svc.DeleteRetroAction(r, userID, actionID) })
		step("item delete", `SELECT count(*) FROM thunderdome.retro_item WHERE id = $1`,
			itemID, 0, func(r string) { _, _ = svc.DeleteRetroItem(r, userID, "worked", itemID) })
	})

	t.Run("poker stories", func(t *testing.T) {
		svc := &poker.Service{DB: d.DB, Logger: d.Logger, HTMLSanitizerPolicy: bluemonday.UGCPolicy()}
		scaleID := insertID(t, d, `SELECT id FROM thunderdome.estimation_scale WHERE name = 'Thunderdome Default'`)
		newGame := func() string {
			g, err := svc.CreateGame(ctx, userID, "integration", scaleID, []string{"1"}, nil, false, "ceil", "", "", false)
			require.NoError(t, err)
			return g.ID
		}
		gameID, otherGameID := newGame(), newGame()
		_, err := svc.CreateStory(gameID, "mine", "story", "", "", "", "", 1)
		require.NoError(t, err)
		storyID := insertID(t, d, `SELECT id FROM thunderdome.poker_story WHERE poker_id = $1`, gameID)
		untouched := func() int {
			return count(t, d, `SELECT count(*) FROM thunderdome.poker_story
				WHERE id = $1 AND name = 'mine' AND NOT active AND points = '' AND votes = '[]'::jsonb`, storyID)
		}

		_, err = svc.UpdateStory(otherGameID, storyID, "theirs", "story", "", "", "", "", 1)
		require.NoError(t, err)
		svc.SetVote(otherGameID, userID, storyID, "1")
		_, err = svc.RetractVote(otherGameID, userID, storyID)
		require.NoError(t, err)
		_, err = svc.ArrangeStory(otherGameID, storyID, "")
		require.NoError(t, err)
		_, err = svc.ActivateStoryVoting(otherGameID, storyID)
		assert.Error(t, err)
		_, err = svc.FinalizeStory(otherGameID, storyID, "1")
		assert.Error(t, err)
		_, err = svc.DeleteStory(otherGameID, storyID)
		assert.Error(t, err)
		assert.Equal(t, 1, untouched())

		_, err = svc.ActivateStoryVoting(gameID, storyID)
		require.NoError(t, err)
		svc.SetVote(gameID, userID, storyID, "1")
		assert.Equal(t, 1, count(t, d, `SELECT jsonb_array_length(votes) FROM thunderdome.poker_story WHERE id = $1`, storyID))
		_, err = svc.RetractVote(gameID, userID, storyID)
		require.NoError(t, err)
		_, err = svc.UpdateStory(gameID, storyID, "renamed", "story", "", "", "", "", 1)
		require.NoError(t, err)
		_, err = svc.FinalizeStory(gameID, storyID, "1")
		require.NoError(t, err)
		assert.Equal(t, 1, count(t, d, `SELECT count(*) FROM thunderdome.poker_story
			WHERE id = $1 AND name = 'renamed' AND points = '1'`, storyID))
		_, err = svc.DeleteStory(gameID, storyID)
		require.NoError(t, err)
		assert.Equal(t, 0, count(t, d, `SELECT count(*) FROM thunderdome.poker_story WHERE id = $1`, storyID))
	})
}
