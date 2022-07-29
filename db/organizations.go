package db

import (
	"context"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// OrganizationGet gets an organization
func (d *Database) OrganizationGet(ctx context.Context, OrgID string) (*model.Organization, error) {
	var org = &model.Organization{}

	e := d.db.QueryRowContext(ctx,
		`SELECT id, name, created_date, updated_date FROM organization_get_by_id($1)`,
		OrgID,
	).Scan(
		&org.Id,
		&org.Name,
		&org.CreatedDate,
		&org.UpdatedDate,
	)
	if e != nil {
		d.logger.Ctx(ctx).Error("organization_get_by_id query error", zap.Error(e))
		return nil, errors.New("error getting organization")
	}

	return org, nil
}

// OrganizationUserRole gets a users role in organization
func (d *Database) OrganizationUserRole(ctx context.Context, UserID string, OrgID string) (string, error) {
	var role string

	e := d.db.QueryRowContext(ctx,
		`SELECT role FROM organization_get_user_role($1, $2)`,
		UserID,
		OrgID,
	).Scan(
		&role,
	)
	if e != nil {
		d.logger.Ctx(ctx).Error("organization_get_user_role query error", zap.Error(e))
		return "", errors.New("error getting organization users role")
	}

	return role, nil
}

// OrganizationListByUser gets a list of organizations the user is apart of
func (d *Database) OrganizationListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*model.Organization {
	var organizations = make([]*model.Organization, 0)
	rows, err := d.db.QueryContext(ctx,
		`SELECT id, name, created_date, updated_date FROM organization_list_by_user($1, $2, $3);`,
		UserID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var org model.Organization

			if err := rows.Scan(
				&org.Id,
				&org.Name,
				&org.CreatedDate,
				&org.UpdatedDate,
			); err != nil {
				d.logger.Ctx(ctx).Error("organization_list_by_user query scan error", zap.Error(err))
			} else {
				organizations = append(organizations, &org)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("organization_list_by_user query error", zap.Error(err))
	}

	return organizations
}

// OrganizationCreate creates an organization
func (d *Database) OrganizationCreate(ctx context.Context, UserID string, OrgName string) (*model.Organization, error) {
	o := &model.Organization{}

	err := d.db.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM organization_create($1, $2);`,
		UserID,
		OrgName,
	).Scan(&o.Id, &o.Name, &o.CreatedDate, &o.UpdatedDate)

	if err != nil {
		d.logger.Ctx(ctx).Error("Unable to create organization", zap.Error(err))
		return nil, err
	}

	return o, nil
}

// OrganizationUserList gets a list of organization users
func (d *Database) OrganizationUserList(ctx context.Context, OrgID string, Limit int, Offset int) []*model.OrganizationUser {
	var users = make([]*model.OrganizationUser, 0)
	rows, err := d.db.QueryContext(ctx,
		`SELECT id, name, email, role, avatar FROM organization_user_list($1, $2, $3);`,
		OrgID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr model.OrganizationUser

			if err := rows.Scan(
				&usr.Id,
				&usr.Name,
				&usr.Email,
				&usr.Role,
				&usr.Avatar,
			); err != nil {
				d.logger.Ctx(ctx).Error("organization_user_list query scan error", zap.Error(err))
			} else {
				usr.GravatarHash = createGravatarHash(usr.Email)
				users = append(users, &usr)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("organization_user_list query error", zap.Error(err))
	}

	return users
}

// OrganizationAddUser adds a user to an organization
func (d *Database) OrganizationAddUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error) {
	_, err := d.db.ExecContext(ctx,
		`SELECT organization_user_add($1, $2, $3);`,
		OrgID,
		UserID,
		Role,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("Unable to add user to organization", zap.Error(err))
		return "", err
	}

	return OrgID, nil
}

// OrganizationRemoveUser removes a user from a organization
func (d *Database) OrganizationRemoveUser(ctx context.Context, OrganizationID string, UserID string) error {
	_, err := d.db.ExecContext(ctx,
		`CALL organization_user_remove($1, $2);`,
		OrganizationID,
		UserID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("Unable to remove user from organization", zap.Error(err))
		return err
	}

	return nil
}

// OrganizationTeamList gets a list of organization teams
func (d *Database) OrganizationTeamList(ctx context.Context, OrgID string, Limit int, Offset int) []*model.Team {
	var teams = make([]*model.Team, 0)
	rows, err := d.db.QueryContext(ctx,
		`SELECT id, name, created_date, updated_date FROM organization_team_list($1, $2, $3);`,
		OrgID,
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
				d.logger.Ctx(ctx).Error("organization_team_list query scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("organization_team_list query error", zap.Error(err))
	}

	return teams
}

// OrganizationTeamCreate creates an organization team
func (d *Database) OrganizationTeamCreate(ctx context.Context, OrgID string, TeamName string) (*model.Team, error) {
	t := &model.Team{}

	err := d.db.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM organization_team_create($1, $2);`,
		OrgID,
		TeamName,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		d.logger.Ctx(ctx).Error("Unable to create organization team", zap.Error(err))
		return nil, err
	}

	return t, nil
}

// OrganizationTeamUserRole gets a users role in organization team
func (d *Database) OrganizationTeamUserRole(ctx context.Context, UserID string, OrgID string, TeamID string) (string, string, error) {
	var orgRole string
	var teamRole string

	e := d.db.QueryRowContext(ctx,
		`SELECT orgRole, teamRole FROM organization_team_user_role($1, $2, $3)`,
		UserID,
		OrgID,
		TeamID,
	).Scan(
		&orgRole,
		&teamRole,
	)
	if e != nil {
		d.logger.Ctx(ctx).Error("organization_team_user_role query error", zap.Error(e))
		return "", "", errors.New("error getting organization team users role")
	}

	return orgRole, teamRole, nil
}

// OrganizationDelete deletes an organization
func (d *Database) OrganizationDelete(ctx context.Context, OrgID string) error {
	_, err := d.db.ExecContext(ctx,
		`CALL organization_delete($1);`,
		OrgID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("organization_delete query error", zap.Error(err))
		return err
	}

	return nil
}
