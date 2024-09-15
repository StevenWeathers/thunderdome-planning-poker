package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// TeamInviteUser invites a user to a team
func (d *Service) TeamInviteUser(ctx context.Context, TeamID string, Email string, Role string) (string, error) {
	var inviteId string
	err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.team_user_invite (team_id, email, role) VALUES ($1, $2, $3) RETURNING invite_id;`,
		TeamID,
		Email,
		Role,
	).Scan(&inviteId)

	if err != nil {
		return "", fmt.Errorf("team invite user query error: %v", err)
	}

	return inviteId, nil
}

// TeamUserGetInviteByID gets a team user invite
func (d *Service) TeamUserGetInviteByID(ctx context.Context, InviteID string) (thunderdome.TeamUserInvite, error) {
	tui := thunderdome.TeamUserInvite{}
	err := d.DB.QueryRowContext(ctx,
		`SELECT invite_id, team_id, email, role, created_date, expire_date
 				FROM thunderdome.team_user_invite WHERE invite_id = $1;`,
		InviteID,
	).Scan(&tui.InviteId, &tui.TeamId, &tui.Email, &tui.Role, &tui.CreatedDate, &tui.ExpireDate)

	if err != nil {
		return tui, fmt.Errorf("team get user invite query error: %v", err)
	}

	return tui, nil
}

// TeamDeleteUserInvite deletes a user team invite
func (d *Service) TeamDeleteUserInvite(ctx context.Context, InviteID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team_user_invite where invite_id = $1;`,
		InviteID,
	)

	if err != nil {
		return fmt.Errorf("team delete user invite query error: %v", err)
	}

	return nil
}

// TeamGetUserInvites gets teams user invites
func (d *Service) TeamGetUserInvites(ctx context.Context, teamId string) ([]thunderdome.TeamUserInvite, error) {
	var invites = make([]thunderdome.TeamUserInvite, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT invite_id, team_id, email, role, created_date, expire_date
 				FROM thunderdome.team_user_invite WHERE team_id = $1;`,
		teamId,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var invite thunderdome.TeamUserInvite

			if err := rows.Scan(
				&invite.InviteId,
				&invite.TeamId,
				&invite.Email,
				&invite.Role,
				&invite.CreatedDate,
				&invite.ExpireDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("TeamGetUserInvites query scan error", zap.Error(err))
			} else {
				invites = append(invites, invite)
			}
		}
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			d.Logger.Ctx(ctx).Error("TeamGetUserInvites query error", zap.Error(err))
		}
	}

	return invites, nil
}
