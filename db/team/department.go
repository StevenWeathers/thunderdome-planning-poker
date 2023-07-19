package team

import (
	"context"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"
)

// DepartmentUserRole gets a users role in department (and organization)
func (d *OrganizationService) DepartmentUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string) (string, string, error) {
	var orgRole string
	var departmentRole string

	e := d.DB.QueryRowContext(ctx,
		`SELECT ou.role AS orgRole, COALESCE(du.role, '') AS departmentRole
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.department_user du ON du.user_id = $1 AND du.department_id = $3
        WHERE ou.organization_id = $2 AND ou.user_id = $1;`,
		UserID,
		OrgID,
		DepartmentID,
	).Scan(
		&orgRole,
		&departmentRole,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("department_get_user_role query error", zap.Error(e))
		return "", "", errors.New("error getting department users role")
	}

	return orgRole, departmentRole, nil
}

// DepartmentGet gets a department
func (d *OrganizationService) DepartmentGet(ctx context.Context, DepartmentID string) (*thunderdome.Department, error) {
	var org = &thunderdome.Department{}

	e := d.DB.QueryRowContext(ctx,
		`SELECT od.id, od.name, od.created_date, od.updated_date
        FROM thunderdome.organization_department od
        WHERE od.id = $1;`,
		DepartmentID,
	).Scan(
		&org.Id,
		&org.Name,
		&org.CreatedDate,
		&org.UpdatedDate,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("department_get_by_id query error", zap.Error(e))
		return nil, errors.New("department not found")
	}

	return org, nil
}

// OrganizationDepartmentList gets a list of organization departments
func (d *OrganizationService) OrganizationDepartmentList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.Department {
	var departments = make([]*thunderdome.Department, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT d.id, d.name, d.created_date, d.updated_date
        FROM thunderdome.organization_department d
        WHERE d.organization_id = $1
        ORDER BY d.created_date
		LIMIT $2
		OFFSET $3;`,
		OrgID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var department thunderdome.Department

			if err := rows.Scan(
				&department.Id,
				&department.Name,
				&department.CreatedDate,
				&department.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("department_list query scan error", zap.Error(err))
			} else {
				departments = append(departments, &department)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("department_list query error", zap.Error(err))
	}

	return departments
}

// DepartmentCreate creates an organization department
func (d *OrganizationService) DepartmentCreate(ctx context.Context, OrgID string, OrgName string) (*thunderdome.Department, error) {
	od := &thunderdome.Department{}

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM thunderdome.department_create($1, $2);`,
		OrgID,
		OrgName,
	).Scan(&od.Id, &od.Name, &od.CreatedDate, &od.UpdatedDate)

	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to create organization department", zap.Error(err))
		return nil, err
	}

	return od, nil
}

// DepartmentTeamList gets a list of department teams
func (d *OrganizationService) DepartmentTeamList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*thunderdome.Team {
	var teams = make([]*thunderdome.Team, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT t.id, t.name, t.created_date, t.updated_date
        FROM thunderdome.department_team dt
        LEFT JOIN thunderdome.team t ON dt.team_id = t.id
        WHERE dt.department_id = $1
        ORDER BY t.created_date
		LIMIT $2
		OFFSET $3;`,
		DepartmentID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team thunderdome.Team

			if err := rows.Scan(
				&team.Id,
				&team.Name,
				&team.CreatedDate,
				&team.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("department_team_list query scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("department_team_list query error", zap.Error(err))
	}

	return teams
}

// DepartmentTeamCreate creates a department team
func (d *OrganizationService) DepartmentTeamCreate(ctx context.Context, DepartmentID string, TeamName string) (*thunderdome.Team, error) {
	t := &thunderdome.Team{}

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM thunderdome.department_team_create($1, $2);`,
		DepartmentID,
		TeamName,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to create department tea", zap.Error(err))
		return nil, err
	}

	return t, nil
}

// DepartmentUserList gets a list of department users
func (d *OrganizationService) DepartmentUserList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*thunderdome.DepartmentUser {
	var users = make([]*thunderdome.DepartmentUser, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT u.id, u.name, COALESCE(u.email, ''), du.role, u.avatar
        FROM thunderdome.department_user du
        LEFT JOIN thunderdome.users u ON du.user_id = u.id
        WHERE du.department_id = $1
        ORDER BY du.created_date
		LIMIT $2
		OFFSET $3;`,
		DepartmentID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr thunderdome.DepartmentUser

			if err := rows.Scan(
				&usr.Id,
				&usr.Name,
				&usr.Email,
				&usr.Role,
				&usr.Avatar,
			); err != nil {
				d.Logger.Ctx(ctx).Error("department_user_list query scan error", zap.Error(err))
			} else {
				usr.GravatarHash = db.CreateGravatarHash(usr.Email)
				users = append(users, &usr)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("department_user_list query error", zap.Error(err))
	}

	return users
}

// DepartmentAddUser adds a user to an organization department
func (d *OrganizationService) DepartmentAddUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`SELECT thunderdome.department_user_add($1, $2, $3);`,
		DepartmentID,
		UserID,
		Role,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to add user to department", zap.Error(err))
		return "", err
	}

	return DepartmentID, nil
}

// DepartmentRemoveUser removes a user from a department (and department teams)
func (d *OrganizationService) DepartmentRemoveUser(ctx context.Context, DepartmentID string, UserID string) error {
	_, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.department_user_remove($1, $2);`,
		DepartmentID,
		UserID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to remove user from department", zap.Error(err))
		return err
	}

	return nil
}

// DepartmentTeamUserRole gets a users role in organization department team
func (d *OrganizationService) DepartmentTeamUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string, TeamID string) (string, string, string, error) {
	var orgRole string
	var departmentRole string
	var teamRole string

	e := d.DB.QueryRowContext(ctx,
		`SELECT ou.role AS orgRole, COALESCE(du.role, '') AS departmentRole, COALESCE(tu.role, '') AS teamRole
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.department_user du ON du.user_id = $1 AND du.department_id = $3
        LEFT JOIN thunderdome.team_user tu ON tu.user_id = $1 AND tu.team_id = $4
        WHERE ou.organization_id = $2 AND ou.user_id = $1;`,
		UserID,
		OrgID,
		DepartmentID,
		TeamID,
	).Scan(
		&orgRole,
		&departmentRole,
		&teamRole,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("department_team_user_role query error", zap.Error(e))
		return "", "", "", errors.New("error getting department team users role")
	}

	return orgRole, departmentRole, teamRole, nil
}

// DepartmentDelete deletes a department
func (d *OrganizationService) DepartmentDelete(ctx context.Context, DepartmentID string) error {
	_, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.department_delete($1);`,
		DepartmentID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("department_delete query error", zap.Error(err))
		return err
	}

	return nil
}
