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
func (d *OrganizationService) OrganizationGet(ctx context.Context, orgID string) (*thunderdome.Organization, error) {
	var org = &thunderdome.Organization{}

	err := d.DB.QueryRowContext(ctx,
		`SELECT o.id, o.name, o.created_date, o.updated_date,
 		CASE WHEN s.id IS NOT NULL AND s.expires > NOW() AND s.active = true THEN true ELSE false END AS is_subscribed
        FROM thunderdome.organization o
        LEFT JOIN thunderdome.subscription s ON o.id = s.organization_id
        WHERE o.id = $1;`,
		orgID,
	).Scan(
		&org.ID,
		&org.Name,
		&org.CreatedDate,
		&org.UpdatedDate,
		&org.Subscribed,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error getting organization: %v", err)
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("ORGANIZATION_NOT_FOUND")
	}

	return org, nil
}

// OrganizationUserRole gets a user's role in organization
func (d *OrganizationService) OrganizationUserRole(ctx context.Context, UserID string, orgID string) (string, error) {
	var role string

	err := d.DB.QueryRowContext(ctx,
		`SELECT ou.role
    FROM thunderdome.organization_user ou
    WHERE ou.organization_id = $2 AND ou.user_id = $1;`,
		UserID,
		orgID,
	).Scan(
		&role,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("error getting organization users role: %v", err)
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("USER_ROLE_NOT_FOUND")
	}

	return role, nil
}

// OrganizationListByUser gets a list of organizations the user is apart of
func (d *OrganizationService) OrganizationListByUser(ctx context.Context, userID string, limit int, offset int) []*thunderdome.UserOrganization {
	var organizations = make([]*thunderdome.UserOrganization, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT o.id, o.name, o.created_date, o.updated_date, ou.role
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.organization o ON ou.organization_id = o.id
        WHERE ou.user_id = $1
        ORDER BY created_date
		LIMIT $2
		OFFSET $3;`,
		userID,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var org thunderdome.UserOrganization

			if err := rows.Scan(
				&org.ID,
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
func (d *OrganizationService) OrganizationCreate(ctx context.Context, userID string, orgName string) (*thunderdome.Organization, error) {
	o := &thunderdome.Organization{}

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM thunderdome.organization_create($1, $2);`,
		userID,
		orgName,
	).Scan(&o.ID, &o.Name, &o.CreatedDate, &o.UpdatedDate)

	if err != nil {
		return nil, fmt.Errorf("organization create query error :%v", err)
	}

	return o, nil
}

// OrganizationUpdate updates an organization
func (d *OrganizationService) OrganizationUpdate(ctx context.Context, orgID string, orgName string) (*thunderdome.Organization, error) {
	o := &thunderdome.Organization{}

	err := d.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.organization
		SET name = $1, updated_date = NOW()
		WHERE id = $2
		RETURNING id, name, created_date, updated_date;`,
		orgName, orgID,
	).Scan(&o.ID, &o.Name, &o.CreatedDate, &o.UpdatedDate)
	if err != nil {
		return nil, fmt.Errorf("organization update query error :%v", err)
	}

	return o, nil
}

// OrganizationUserList gets a list of organization users
func (d *OrganizationService) OrganizationUserList(ctx context.Context, orgID string, limit int, offset int) []*thunderdome.OrganizationUser {
	var users = make([]*thunderdome.OrganizationUser, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT u.id, u.name, COALESCE(u.email, ''), ou.role, u.avatar, COALESCE(u.picture, '')
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.users u ON ou.user_id = u.id
        WHERE ou.organization_id = $1
        ORDER BY ou.created_date
		LIMIT $2
		OFFSET $3;`,
		orgID,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr thunderdome.OrganizationUser

			if err := rows.Scan(
				&usr.ID,
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
func (d *OrganizationService) OrganizationAddUser(ctx context.Context, orgID string, userID string, role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.organization_user (organization_id, user_id, role) VALUES ($1, $2, $3);`,
		orgID,
		userID,
		role,
	)

	if err != nil {
		return "", fmt.Errorf("organization add user query error: %v", err)
	}

	return orgID, nil
}

// OrganizationUpsertUser adds a user to an organization if not existing otherwise does nothing
func (d *OrganizationService) OrganizationUpsertUser(ctx context.Context, orgID string, userID string, role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.organization_user (organization_id, user_id, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`,
		orgID,
		userID,
		role,
	)

	if err != nil {
		return "", fmt.Errorf("organization upsert user query error: %v", err)
	}

	return orgID, nil
}

// OrganizationUpdateUser updates an organization user
func (d *OrganizationService) OrganizationUpdateUser(ctx context.Context, orgID string, userID string, role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.organization_user SET role = $3 WHERE organization_id = $1 AND user_id = $2;`,
		orgID,
		userID,
		role,
	)

	if err != nil {
		return "", fmt.Errorf("organization update user query error: %v", err)
	}

	return orgID, nil
}

// OrganizationRemoveUser removes a user from an organization
func (d *OrganizationService) OrganizationRemoveUser(ctx context.Context, orgID string, userID string) error {
	_, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.organization_user_remove($1, $2);`,
		orgID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("organization remove user query error: %v", err)
	}

	return nil
}

// OrganizationInviteUser invites a user to an organization
func (d *OrganizationService) OrganizationInviteUser(ctx context.Context, orgID string, email string, role string) (string, error) {
	var inviteID string
	err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.organization_user_invite (organization_id, email, role) VALUES ($1, $2, $3) RETURNING invite_id;`,
		orgID,
		email,
		role,
	).Scan(&inviteID)

	if err != nil {
		return "", fmt.Errorf("organization invite user query error: %v", err)
	}

	return inviteID, nil
}

// OrganizationUserGetInviteByID gets an organization user invite
func (d *OrganizationService) OrganizationUserGetInviteByID(ctx context.Context, inviteID string) (thunderdome.OrganizationUserInvite, error) {
	oui := thunderdome.OrganizationUserInvite{}
	err := d.DB.QueryRowContext(ctx,
		`SELECT invite_id, organization_id, email, role, created_date, expire_date
 				FROM thunderdome.organization_user_invite WHERE invite_id = $1 AND expire_date > CURRENT_TIMESTAMP;`,
		inviteID,
	).Scan(&oui.InviteID, &oui.OrganizationID, &oui.Email, &oui.Role, &oui.CreatedDate, &oui.ExpireDate)

	if err != nil {
		d.Logger.Ctx(ctx).Error("OrganizationUserGetInviteByID query error	", zap.Error(err),
			zap.String("invite_id", inviteID))
		return oui, fmt.Errorf("organization get user invite query error: %v", err)
	}

	return oui, nil
}

// OrganizationDeleteUserInvite deletes a user organization invite
func (d *OrganizationService) OrganizationDeleteUserInvite(ctx context.Context, inviteID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.organization_user_invite where invite_id = $1;`,
		inviteID,
	)

	if err != nil {
		return fmt.Errorf("organization delete user invite query error: %v", err)
	}

	return nil
}

// OrganizationGetUserInvites gets organizations user invites
func (d *OrganizationService) OrganizationGetUserInvites(ctx context.Context, orgID string) ([]thunderdome.OrganizationUserInvite, error) {
	var invites = make([]thunderdome.OrganizationUserInvite, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT invite_id, organization_id, email, role, created_date, expire_date
 				FROM thunderdome.organization_user_invite WHERE organization_id = $1;`,
		orgID,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var invite thunderdome.OrganizationUserInvite

			if err := rows.Scan(
				&invite.InviteID,
				&invite.OrganizationID,
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
func (d *OrganizationService) OrganizationTeamList(ctx context.Context, orgID string, limit int, offset int) []*thunderdome.Team {
	var teams = make([]*thunderdome.Team, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT t.id, t.name, t.created_date, t.updated_date
        FROM thunderdome.team t
        WHERE t.organization_id = $1
        ORDER BY t.created_date
		LIMIT $2
		OFFSET $3;`,
		orgID,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team thunderdome.Team

			if err := rows.Scan(
				&team.ID,
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
func (d *OrganizationService) OrganizationTeamCreate(ctx context.Context, orgID string, teamName string) (*thunderdome.Team, error) {
	t := &thunderdome.Team{}

	err := d.DB.QueryRowContext(ctx, `
		INSERT INTO thunderdome.team (name, organization_id)
		VALUES ($1, $2) RETURNING id, name, created_date, updated_date;`,
		teamName,
		orgID,
	).Scan(&t.ID, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		return nil, fmt.Errorf("organization create team query error: %v", err)
	}

	return t, nil
}

// OrganizationTeamUserRole gets a user's role in organization team
func (d *OrganizationService) OrganizationTeamUserRole(ctx context.Context, userID string, orgID string, teamID string) (string, string, error) {
	var orgRole string
	var teamRole string

	err := d.DB.QueryRowContext(ctx,
		`SELECT ou.role AS orgRole, COALESCE(tu.role, '') AS teamRole
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.team_user tu ON tu.user_id = $1 AND tu.team_id = $3
        WHERE ou.organization_id = $2 AND ou.user_id = $1;`,
		userID,
		orgID,
		teamID,
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
func (d *OrganizationService) OrganizationDelete(ctx context.Context, orgID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.organization WHERE id = $1;`,
		orgID,
	)

	if err != nil {
		return fmt.Errorf("organization delete query error: %v", err)
	}

	return nil
}

// OrganizationList gets a list of organizations
func (d *OrganizationService) OrganizationList(ctx context.Context, limit int, offset int) []*thunderdome.Organization {
	var organizations = make([]*thunderdome.Organization, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT o.id, o.name, o.created_date, o.updated_date
        FROM thunderdome.organization o
        ORDER BY o.created_date
		LIMIT $1
		OFFSET $2;`,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var org thunderdome.Organization

			if err := rows.Scan(
				&org.ID,
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

func (d *OrganizationService) OrganizationIsSubscribed(ctx context.Context, orgID string) (bool, error) {
	var subscribed bool

	err := d.DB.QueryRowContext(ctx,
		`SELECT
    COALESCE(
        (SELECT TRUE
         FROM thunderdome.subscription
         WHERE organization_id = $1
           AND active = TRUE
           AND expires > CURRENT_TIMESTAMP
         LIMIT 1),
        FALSE
    ) AS is_subscribed;`,
		orgID,
	).Scan(
		&subscribed,
	)
	if err != nil {
		return false, fmt.Errorf("error getting organization subscription: %v", err)
	}

	return subscribed, nil
}
