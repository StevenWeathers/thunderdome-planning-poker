package db

import (
	"context"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// DepartmentUserRole gets a users role in department (and organization)
func (d *Database) DepartmentUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string) (string, string, error) {
	var orgRole string
	var departmentRole string

	e := d.db.QueryRowContext(ctx,
		`SELECT orgRole, departmentRole FROM department_get_user_role($1, $2, $3)`,
		UserID,
		OrgID,
		DepartmentID,
	).Scan(
		&orgRole,
		&departmentRole,
	)
	if e != nil {
		d.logger.Ctx(ctx).Error("department_get_user_role query error", zap.Error(e))
		return "", "", errors.New("error getting department users role")
	}

	return orgRole, departmentRole, nil
}

// DepartmentGet gets a department
func (d *Database) DepartmentGet(ctx context.Context, DepartmentID string) (*model.Department, error) {
	var org = &model.Department{}

	e := d.db.QueryRowContext(ctx,
		`SELECT id, name, created_date, updated_date FROM department_get_by_id($1)`,
		DepartmentID,
	).Scan(
		&org.Id,
		&org.Name,
		&org.CreatedDate,
		&org.UpdatedDate,
	)
	if e != nil {
		d.logger.Ctx(ctx).Error("department_get_by_id query error", zap.Error(e))
		return nil, errors.New("department not found")
	}

	return org, nil
}

// OrganizationDepartmentList gets a list of organization departments
func (d *Database) OrganizationDepartmentList(ctx context.Context, OrgID string, Limit int, Offset int) []*model.Department {
	var departments = make([]*model.Department, 0)
	rows, err := d.db.QueryContext(ctx,
		`SELECT id, name, created_date, updated_date FROM department_list($1, $2, $3);`,
		OrgID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var department model.Department

			if err := rows.Scan(
				&department.Id,
				&department.Name,
				&department.CreatedDate,
				&department.UpdatedDate,
			); err != nil {
				d.logger.Ctx(ctx).Error("department_list query scan error", zap.Error(err))
			} else {
				departments = append(departments, &department)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("department_list query error", zap.Error(err))
	}

	return departments
}

// DepartmentCreate creates an organization department
func (d *Database) DepartmentCreate(ctx context.Context, OrgID string, OrgName string) (*model.Department, error) {
	od := &model.Department{}

	err := d.db.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM department_create($1, $2);`,
		OrgID,
		OrgName,
	).Scan(&od.Id, &od.Name, &od.CreatedDate, &od.UpdatedDate)

	if err != nil {
		d.logger.Ctx(ctx).Error("Unable to create organization department", zap.Error(err))
		return nil, err
	}

	return od, nil
}

// DepartmentTeamList gets a list of department teams
func (d *Database) DepartmentTeamList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*model.Team {
	var teams = make([]*model.Team, 0)
	rows, err := d.db.QueryContext(ctx,
		`SELECT id, name, created_date, updated_date FROM department_team_list($1, $2, $3);`,
		DepartmentID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team model.Team

			if err := rows.Scan(
				&team.Id,
				&team.Name,
				&team.CreatedDate,
				&team.UpdatedDate,
			); err != nil {
				d.logger.Ctx(ctx).Error("department_team_list query scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("department_team_list query error", zap.Error(err))
	}

	return teams
}

// DepartmentTeamCreate creates a department team
func (d *Database) DepartmentTeamCreate(ctx context.Context, DepartmentID string, TeamName string) (*model.Team, error) {
	t := &model.Team{}

	err := d.db.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM department_team_create($1, $2);`,
		DepartmentID,
		TeamName,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		d.logger.Ctx(ctx).Error("Unable to create department tea", zap.Error(err))
		return nil, err
	}

	return t, nil
}

// DepartmentUserList gets a list of department users
func (d *Database) DepartmentUserList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*model.DepartmentUser {
	var users = make([]*model.DepartmentUser, 0)
	rows, err := d.db.QueryContext(ctx,
		`SELECT id, name, email, role, avatar FROM department_user_list($1, $2, $3);`,
		DepartmentID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr model.DepartmentUser

			if err := rows.Scan(
				&usr.Id,
				&usr.Name,
				&usr.Email,
				&usr.Role,
				&usr.Avatar,
			); err != nil {
				d.logger.Ctx(ctx).Error("department_user_list query scan error", zap.Error(err))
			} else {
				usr.GravatarHash = createGravatarHash(usr.Email)
				users = append(users, &usr)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("department_user_list query error", zap.Error(err))
	}

	return users
}

// DepartmentAddUser adds a user to an organization department
func (d *Database) DepartmentAddUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error) {
	_, err := d.db.ExecContext(ctx,
		`SELECT department_user_add($1, $2, $3);`,
		DepartmentID,
		UserID,
		Role,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("Unable to add user to department", zap.Error(err))
		return "", err
	}

	return DepartmentID, nil
}

// DepartmentRemoveUser removes a user from a department (and department teams)
func (d *Database) DepartmentRemoveUser(ctx context.Context, DepartmentID string, UserID string) error {
	_, err := d.db.ExecContext(ctx,
		`CALL department_user_remove($1, $2);`,
		DepartmentID,
		UserID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("Unable to remove user from department", zap.Error(err))
		return err
	}

	return nil
}

// DepartmentTeamUserRole gets a users role in organization department team
func (d *Database) DepartmentTeamUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string, TeamID string) (string, string, string, error) {
	var orgRole string
	var departmentRole string
	var teamRole string

	e := d.db.QueryRowContext(ctx,
		`SELECT orgRole, departmentRole, teamRole FROM department_team_user_role($1, $2, $3, $4)`,
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
		d.logger.Ctx(ctx).Error("department_team_user_role query error", zap.Error(e))
		return "", "", "", errors.New("error getting department team users role")
	}

	return orgRole, departmentRole, teamRole, nil
}

// DepartmentDelete deletes a department
func (d *Database) DepartmentDelete(ctx context.Context, DepartmentID string) error {
	_, err := d.db.ExecContext(ctx,
		`CALL department_delete($1);`,
		DepartmentID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("department_delete query error", zap.Error(err))
		return err
	}

	return nil
}
