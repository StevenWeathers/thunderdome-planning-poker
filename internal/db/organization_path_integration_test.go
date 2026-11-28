package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/team"
)

func TestOrganizationPathExists(t *testing.T) {
	d := testDB(t)
	ctx := context.Background()
	orgs := &team.OrganizationService{DB: d.DB, Logger: d.Logger}

	organization := func(name string) string {
		return insertID(t, d, `INSERT INTO thunderdome.organization (name) VALUES ($1) RETURNING id`, name)
	}
	orgID := organization("mine")
	otherOrgID := organization("other")
	department := func(name string) string {
		return insertID(t, d, `INSERT INTO thunderdome.organization_department (name, organization_id)
			VALUES ($1, $2) RETURNING id`, name, orgID)
	}
	deptID := department("first")
	otherDeptID := department("second")
	orgTeamID := insertID(t, d, `INSERT INTO thunderdome.team (name, organization_id)
		VALUES ('integration', $1) RETURNING id`, orgID)
	deptTeamID := insertID(t, d, `INSERT INTO thunderdome.team (name, department_id)
		VALUES ('integration', $1) RETURNING id`, deptID)

	tests := []struct {
		name   string
		orgID  string
		deptID string
		teamID string
		want   bool
	}{
		{name: "department in organization", orgID: orgID, deptID: deptID, want: true},
		{name: "department in another organization", orgID: otherOrgID, deptID: deptID},
		{name: "organization team", orgID: orgID, teamID: orgTeamID, want: true},
		{name: "organization team of another organization", orgID: otherOrgID, teamID: orgTeamID},
		{name: "department team via organization", orgID: orgID, teamID: deptTeamID, want: true},
		{name: "department team via department", orgID: orgID, deptID: deptID, teamID: deptTeamID, want: true},
		{name: "department team of another organization", orgID: otherOrgID, deptID: deptID, teamID: deptTeamID},
		{name: "team outside the department", orgID: orgID, deptID: otherDeptID, teamID: deptTeamID},
		{name: "organization team via department", orgID: orgID, deptID: deptID, teamID: orgTeamID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orgs.OrganizationPathExists(ctx, tt.orgID, tt.deptID, tt.teamID)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
