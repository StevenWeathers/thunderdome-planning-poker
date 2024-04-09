package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// OrganizationService represents a PostgreSQL implementation of thunderdome.OrganizationDataSvc.
type OrganizationService struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// OrganizationGet gets an organization
func (d *OrganizationService) OrganizationGet(ctx context.Context, OrgID string) (*thunderdome.Organization, error) {
	var org = &thunderdome.Organization{}

	err := d.DB.QueryRowContext(ctx,
		`SELECT o.id, o.name, o.created_date, o.updated_date
        FROM thunderdome.organization o
        WHERE o.id = $1;`,
		OrgID,
	).Scan(
		&org.Id,
		&org.Name,
		&org.CreatedDate,
		&org.UpdatedDate,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting organization: %v", err)
	}

	return org, nil
}

// OrganizationUserRole gets a user's role in organization
func (d *OrganizationService) OrganizationUserRole(ctx context.Context, UserID string, OrgID string) (string, error) {
	var role string

	err := d.DB.QueryRowContext(ctx,
		`SELECT ou.role
    FROM thunderdome.organization_user ou
    WHERE ou.organization_id = $2 AND ou.user_id = $1;`,
		UserID,
		OrgID,
	).Scan(
		&role,
	)
	if err != nil {
		return "", fmt.Errorf("error getting organization users role: %v", err)
	}

	return role, nil
}

// OrganizationListByUser gets a list of organizations the user is apart of
func (d *OrganizationService) OrganizationListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*thunderdome.UserOrganization {
	var organizations = make([]*thunderdome.UserOrganization, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT o.id, o.name, o.created_date, o.updated_date, ou.role
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.organization o ON ou.organization_id = o.id
        WHERE ou.user_id = $1
        ORDER BY created_date
		LIMIT $2
		OFFSET $3;`,
		UserID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var org thunderdome.UserOrganization

			if err := rows.Scan(
				&org.Id,
				&org.Name,
				&org.CreatedDate,
				&org.UpdatedDate,
				&org.Role,
			); err != nil {
				d.Logger.Ctx(ctx).Error("organization_list_by_user query scan error", zap.Error(err))
			} else {
				organizations = append(organizations, &org)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("organization_list_by_user query error", zap.Error(err))
	}

	return organizations
}

// OrganizationCreate creates an organization
func (d *OrganizationService) OrganizationCreate(ctx context.Context, UserID string, OrgName string) (*thunderdome.Organization, error) {
	o := &thunderdome.Organization{}

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM thunderdome.organization_create($1, $2);`,
		UserID,
		OrgName,
	).Scan(&o.Id, &o.Name, &o.CreatedDate, &o.UpdatedDate)

	if err != nil {
		return nil, fmt.Errorf("organization create query error :%v", err)
	}

	return o, nil
}

// OrganizationUpdate updates an organization
func (d *OrganizationService) OrganizationUpdate(ctx context.Context, OrgId string, OrgName string) (*thunderdome.Organization, error) {
	o := &thunderdome.Organization{}

	err := d.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.organization
		SET name = $1, updated_date = NOW()
		WHERE id = $2
		RETURNING id, name, created_date, updated_date;`,
		OrgName, OrgId,
	).Scan(&o.Id, &o.Name, &o.CreatedDate, &o.UpdatedDate)
	if err != nil {
		return nil, fmt.Errorf("organization update query error :%v", err)
	}

	return o, nil
}

// OrganizationUserList gets a list of organization users
func (d *OrganizationService) OrganizationUserList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.OrganizationUser {
	var users = make([]*thunderdome.OrganizationUser, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT u.id, u.name, COALESCE(u.email, ''), ou.role, u.avatar, COALESCE(u.picture_url, '')
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.users u ON ou.user_id = u.id
        WHERE ou.organization_id = $1
        ORDER BY ou.created_date
		LIMIT $2
		OFFSET $3;`,
		OrgID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr thunderdome.OrganizationUser

			if err := rows.Scan(
				&usr.Id,
				&usr.Name,
				&usr.Email,
				&usr.Role,
				&usr.Avatar,
				&usr.PictureURL,
			); err != nil {
				d.Logger.Ctx(ctx).Error("organization_user_list query scan error", zap.Error(err))
			} else {
				usr.GravatarHash = db.CreateGravatarHash(usr.Email)
				users = append(users, &usr)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("organization_user_list query error", zap.Error(err))
	}

	return users
}

// OrganizationAddUser adds a user to an organization
func (d *OrganizationService) OrganizationAddUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.organization_user (organization_id, user_id, role) VALUES ($1, $2, $3);`,
		OrgID,
		UserID,
		Role,
	)

	if err != nil {
		return "", fmt.Errorf("organization add user query error: %v", err)
	}

	return OrgID, nil
}

// OrganizationUpdateUser updates an organization user
func (d *OrganizationService) OrganizationUpdateUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.organization_user SET role = $3 WHERE organization_id = $1 AND user_id = $2;`,
		OrgID,
		UserID,
		Role,
	)

	if err != nil {
		return "", fmt.Errorf("organization update user query error: %v", err)
	}

	return OrgID, nil
}

// OrganizationRemoveUser removes a user from an organization
func (d *OrganizationService) OrganizationRemoveUser(ctx context.Context, OrganizationID string, UserID string) error {
	_, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.organization_user_remove($1, $2);`,
		OrganizationID,
		UserID,
	)

	if err != nil {
		return fmt.Errorf("organization remove user query error: %v", err)
	}

	return nil
}

// OrganizationInviteUser invites a user to an organization
func (d *OrganizationService) OrganizationInviteUser(ctx context.Context, OrgID string, Email string, Role string) (string, error) {
	var inviteId string
	err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.organization_user_invite (organization_id, email, role) VALUES ($1, $2, $3) RETURNING invite_id;`,
		OrgID,
		Email,
		Role,
	).Scan(&inviteId)

	if err != nil {
		return "", fmt.Errorf("organization invite user query error: %v", err)
	}

	return inviteId, nil
}

// OrganizationUserGetInviteByID gets a organization user invite
func (d *OrganizationService) OrganizationUserGetInviteByID(ctx context.Context, InviteID string) (thunderdome.OrganizationUserInvite, error) {
	oui := thunderdome.OrganizationUserInvite{}
	err := d.DB.QueryRowContext(ctx,
		`SELECT invite_id, organization_id, email, role, created_date, expire_date
 				FROM thunderdome.organization_user_invite WHERE invite_id = $1;`,
		InviteID,
	).Scan(&oui.InviteId, &oui.OrganizationId, &oui.Email, &oui.Role, &oui.CreatedDate, &oui.ExpireDate)

	if err != nil {
		return oui, fmt.Errorf("organization get user invite query error: %v", err)
	}

	return oui, nil
}

// OrganizationDeleteUserInvite deletes a user organization invite
func (d *OrganizationService) OrganizationDeleteUserInvite(ctx context.Context, InviteID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.organization_user_invite where invite_id = $1;`,
		InviteID,
	)

	if err != nil {
		return fmt.Errorf("organization delete user invite query error: %v", err)
	}

	return nil
}

// OrganizationGetUserInvites gets organizations user invites
func (d *OrganizationService) OrganizationGetUserInvites(ctx context.Context, orgId string) ([]thunderdome.OrganizationUserInvite, error) {
	var invites = make([]thunderdome.OrganizationUserInvite, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT invite_id, organization_id, email, role, created_date, expire_date
 				FROM thunderdome.organization_user_invite WHERE organization_id = $1;`,
		orgId,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var invite thunderdome.OrganizationUserInvite

			if err := rows.Scan(
				&invite.InviteId,
				&invite.OrganizationId,
				&invite.Email,
				&invite.Role,
				&invite.CreatedDate,
				&invite.ExpireDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("OrganizationGetUserInvites query scan error", zap.Error(err))
			} else {
				invites = append(invites, invite)
			}
		}
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			d.Logger.Ctx(ctx).Error("OrganizationGetUserInvites query error", zap.Error(err))
		}
	}

	return invites, nil
}

// OrganizationTeamList gets a list of organization teams
func (d *OrganizationService) OrganizationTeamList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.Team {
	var teams = make([]*thunderdome.Team, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT t.id, t.name, t.created_date, t.updated_date
        FROM thunderdome.team t
        WHERE t.organization_id = $1
        ORDER BY t.created_date
		LIMIT $2
		OFFSET $3;`,
		OrgID,
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
				d.Logger.Ctx(ctx).Error("organization_team_list query scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("organization_team_list query error", zap.Error(err))
	}

	return teams
}

// OrganizationTeamCreate creates an organization team
func (d *OrganizationService) OrganizationTeamCreate(ctx context.Context, OrgID string, TeamName string) (*thunderdome.Team, error) {
	t := &thunderdome.Team{}

	err := d.DB.QueryRowContext(ctx, `
		INSERT INTO thunderdome.team (name, organization_id) 
		VALUES ($1, $2) RETURNING id, name, created_date, updated_date;`,
		TeamName,
		OrgID,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		return nil, fmt.Errorf("organization create team query error: %v", err)
	}

	return t, nil
}

// OrganizationTeamUserRole gets a user's role in organization team
func (d *OrganizationService) OrganizationTeamUserRole(ctx context.Context, UserID string, OrgID string, TeamID string) (string, string, error) {
	var orgRole string
	var teamRole string

	err := d.DB.QueryRowContext(ctx,
		`SELECT ou.role AS orgRole, COALESCE(tu.role, '') AS teamRole
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.team_user tu ON tu.user_id = $1 AND tu.team_id = $3
        WHERE ou.organization_id = $2 AND ou.user_id = $1;`,
		UserID,
		OrgID,
		TeamID,
	).Scan(
		&orgRole,
		&teamRole,
	)
	if err != nil {
		return "", "", fmt.Errorf("error getting organization team users role: %v", err)
	}

	return orgRole, teamRole, nil
}

// OrganizationDelete deletes an organization
func (d *OrganizationService) OrganizationDelete(ctx context.Context, OrgID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.organization WHERE id = $1;`,
		OrgID,
	)

	if err != nil {
		return fmt.Errorf("organization delete query error: %v", err)
	}

	return nil
}

// OrganizationList gets a list of organizations
func (d *OrganizationService) OrganizationList(ctx context.Context, Limit int, Offset int) []*thunderdome.Organization {
	var organizations = make([]*thunderdome.Organization, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT o.id, o.name, o.created_date, o.updated_date
        FROM thunderdome.organization o
        ORDER BY o.created_date
		LIMIT $1
		OFFSET $2;`,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var org thunderdome.Organization

			if err := rows.Scan(
				&org.Id,
				&org.Name,
				&org.CreatedDate,
				&org.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("organization_list scan error", zap.Error(err))
			} else {
				organizations = append(organizations, &org)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("organization_list query error", zap.Error(err))
	}

	return organizations
}
