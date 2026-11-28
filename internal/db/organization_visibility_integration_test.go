package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/team"
)

func TestOrganizationListsByUser(t *testing.T) {
	d := testDB(t)
	ctx := context.Background()
	orgs := &team.OrganizationService{DB: d.DB, Logger: d.Logger}

	userID := insertUser(t, d)
	orgID := insertID(t, d, `INSERT INTO thunderdome.organization (name) VALUES ('integration') RETURNING id`)
	department := func(name string) string {
		return insertID(t, d, `INSERT INTO thunderdome.organization_department (name, organization_id)
			VALUES ($1, $2) RETURNING id`, name, orgID)
	}
	orgTeam := func() string {
		return insertID(t, d, `INSERT INTO thunderdome.team (name, organization_id)
			VALUES ('integration', $1) RETURNING id`, orgID)
	}
	deptTeam := func(deptID string) string {
		return insertID(t, d, `INSERT INTO thunderdome.team (name, department_id)
			VALUES ('integration', $1) RETURNING id`, deptID)
	}
	join := func(query string, id string) {
		_, err := d.DB.Exec(query, id, userID)
		require.NoError(t, err)
	}

	myDept := department("mine")
	department("other")
	myOrgTeam := orgTeam()
	orgTeam()
	myDeptTeam := deptTeam(myDept)
	deptTeam(myDept)

	join(`INSERT INTO thunderdome.department_user (department_id, user_id, role) VALUES ($1, $2, 'MEMBER')`, myDept)
	join(`INSERT INTO thunderdome.team_user (team_id, user_id, role) VALUES ($1, $2, 'MEMBER')`, myOrgTeam)
	join(`INSERT INTO thunderdome.team_user (team_id, user_id, role) VALUES ($1, $2, 'MEMBER')`, myDeptTeam)

	departments := orgs.OrganizationDepartmentListByUser(ctx, orgID, userID, 10, 0)
	require.Len(t, departments, 1)
	assert.Equal(t, myDept, departments[0].ID)

	teams := orgs.OrganizationTeamListByUser(ctx, orgID, userID, 10, 0)
	require.Len(t, teams, 1)
	assert.Equal(t, myOrgTeam, teams[0].ID)

	teams = orgs.DepartmentTeamListByUser(ctx, myDept, userID, 10, 0)
	require.Len(t, teams, 1)
	assert.Equal(t, myDeptTeam, teams[0].ID)
}
