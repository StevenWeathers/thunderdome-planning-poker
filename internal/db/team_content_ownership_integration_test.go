package db_test

import (
	"context"
	"testing"

	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/retro"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/team"
)

// TestTeamContentOwnership checks that only the author or an admin can change check-ins, check-in comments,
// kudos and retro action comments, and who may use a retro's actions.
func TestTeamContentOwnership(t *testing.T) {
	d := testDB(t)
	ctx := context.Background()
	checkins := &team.CheckinService{DB: d.DB, Logger: d.Logger, HTMLSanitizerPolicy: bluemonday.UGCPolicy()}
	retros := &retro.Service{DB: d.DB, Logger: d.Logger}
	exec := func(query string, args ...any) {
		t.Helper()
		_, err := d.DB.Exec(query, args...)
		require.NoError(t, err)
	}

	orgID := insertID(t, d, `INSERT INTO thunderdome.organization (name) VALUES ('integration') RETURNING id`)
	deptID := insertID(t, d, `INSERT INTO thunderdome.organization_department (name, organization_id)
		VALUES ('integration', $1) RETURNING id`, orgID)
	teamID := insertID(t, d, `INSERT INTO thunderdome.team (name, department_id)
		VALUES ('integration', $1) RETURNING id`, deptID)
	member := func(table string, column string, id string, role string) string {
		userID := insertUser(t, d)
		exec(`INSERT INTO thunderdome.`+table+` (`+column+`, user_id, role) VALUES ($1, $2, $3)`, id, userID, role)
		return userID
	}
	teamMember := func(role string) string { return member("team_user", "team_id", teamID, role) }

	author := teamMember("MEMBER")
	target := teamMember("MEMBER")
	// kudos read both users' email, which registered users always have
	exec(`UPDATE thunderdome.users SET email = id::text || '@example.com' WHERE id IN ($1, $2)`, author, target)
	siteAdmin := insertUser(t, d)
	exec(`UPDATE thunderdome.users SET type = 'ADMIN' WHERE id = $1`, siteAdmin)

	templateID := insertID(t, d, `SELECT id FROM thunderdome.retro_template LIMIT 1`)
	owner := insertUser(t, d)
	teamRetro, err := retros.CreateRetro(ctx, owner, "", "integration", "secret", "", 3, "visible",
		0, false, false, false, false, templateID)
	require.NoError(t, err)
	exec(`UPDATE thunderdome.retro SET team_id = $1 WHERE id = $2`, teamID, teamRetro.ID)
	_, err = retros.CreateRetroAction(teamRetro.ID, author, "integration")
	require.NoError(t, err)
	actionID := insertID(t, d, `SELECT id FROM thunderdome.retro_action WHERE retro_id = $1`, teamRetro.ID)

	joined := insertUser(t, d)
	exec(`INSERT INTO thunderdome.retro_user (retro_id, user_id) VALUES ($1, $2)`, teamRetro.ID, joined)

	tests := []struct {
		name        string
		userID      string
		mayChange   bool // content written by the author
		retroAccess bool
	}{
		{name: "author", userID: author, mayChange: true, retroAccess: true},
		{name: "teammate", userID: teamMember("MEMBER"), retroAccess: true},
		{name: "department member", userID: member("department_user", "department_id", deptID, "MEMBER")},
		{name: "organization member", userID: member("organization_user", "organization_id", orgID, "MEMBER")},
		{name: "outsider", userID: insertUser(t, d)},
		{name: "retro owner", userID: owner, retroAccess: true},
		{name: "retro participant", userID: joined, retroAccess: true},
		{name: "team admin", userID: teamMember("ADMIN"), mayChange: true, retroAccess: true},
		{name: "department admin", userID: member("department_user", "department_id", deptID, "ADMIN"),
			mayChange: true, retroAccess: true},
		{name: "organization admin", userID: member("organization_user", "organization_id", orgID, "ADMIN"),
			mayChange: true, retroAccess: true},
		{name: "site admin", userID: siteAdmin, mayChange: true, retroAccess: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := func(err error, forbidden string) {
				t.Helper()
				if tt.mayChange {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, forbidden)
				}
			}

			access, err := retros.CanUserAccessRetro(ctx, teamRetro.ID, tt.userID)
			require.NoError(t, err)
			assert.Equal(t, tt.retroAccess, access, "retro access")

			exec(`DELETE FROM thunderdome.team_checkin WHERE team_id = $1`, teamID)
			exec(`DELETE FROM thunderdome.team_kudos WHERE team_id = $1`, teamID)

			require.NoError(t, checkins.CheckinCreate(ctx, teamID, author, "", "y", "t", "", "", true))
			checkinID := insertID(t, d, `SELECT id FROM thunderdome.team_checkin WHERE team_id = $1`, teamID)
			require.NoError(t, checkins.CheckinComment(ctx, teamID, checkinID, author, "original"))
			commentID := insertID(t, d, `SELECT id FROM thunderdome.team_checkin_comment WHERE checkin_id = $1`, checkinID)

			check(checkins.CheckinCommentEdit(ctx, teamID, tt.userID, commentID, "edited"), "REQUIRES_COMMENT_OWNER")
			check(checkins.CheckinCommentDelete(ctx, teamID, tt.userID, commentID), "REQUIRES_COMMENT_OWNER")
			check(checkins.CheckinUpdate(ctx, teamID, tt.userID, checkinID, "y", "edited", "", "", true),
				"REQUIRES_CHECKIN_OWNER")
			check(checkins.CheckinDelete(ctx, teamID, tt.userID, checkinID), "REQUIRES_CHECKIN_OWNER")

			kudo, err := checkins.KudoCreate(ctx, teamID, author, target, "", "original")
			require.NoError(t, err)
			_, err = checkins.KudoUpdate(ctx, teamID, tt.userID, kudo.ID, target, "", "edited")
			check(err, "REQUIRES_KUDO_OWNER")
			check(checkins.KudoDelete(ctx, teamID, tt.userID, kudo.ID), "REQUIRES_KUDO_OWNER")

			_, err = retros.RetroActionCommentAdd(teamRetro.ID, actionID, author, "original")
			require.NoError(t, err)
			actionCommentID := insertID(t, d, `SELECT id FROM thunderdome.retro_action_comment
				WHERE action_id = $1 ORDER BY created_date DESC LIMIT 1`, actionID)
			_, err = retros.RetroActionCommentEdit(teamRetro.ID, actionID, actionCommentID, tt.userID, "edited")
			check(err, "REQUIRES_COMMENT_OWNER")
			_, err = retros.RetroActionCommentDelete(teamRetro.ID, actionID, actionCommentID, tt.userID)
			check(err, "REQUIRES_COMMENT_OWNER")
		})
	}
}
