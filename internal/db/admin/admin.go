package admin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

// Service represents the admin database service
type Service struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GetAppStats gets counts of common application metrics such as users and poker games
func (d *Service) GetAppStats(ctx context.Context) (*thunderdome.ApplicationStats, error) {
	var appStats thunderdome.ApplicationStats

	err := d.DB.QueryRowContext(ctx, `
		SELECT
    (SELECT COUNT(*) FROM thunderdome.users WHERE email IS NULL) AS unregistered_user_count,
    (SELECT COUNT(*) FROM thunderdome.users WHERE email IS NOT NULL) AS registered_user_count,
    (SELECT COUNT(*) FROM thunderdome.poker) AS poker_count,
    (SELECT COUNT(*) FROM thunderdome.poker_story) AS poker_story_count,
    (SELECT COUNT(*) FROM thunderdome.organization) AS organization_count,
    (SELECT COUNT(*) FROM thunderdome.organization_department) AS department_count,
    (SELECT COUNT(*) FROM thunderdome.team) AS team_count,
    (SELECT COUNT(*) FROM thunderdome.api_key) AS apikey_count,
    (SELECT COUNT(DISTINCT poker_id) FROM thunderdome.poker_user WHERE active IS true) AS active_poker_count,
    (SELECT COUNT(user_id) FROM thunderdome.poker_user WHERE active IS true) AS active_poker_user_count,
    (SELECT COUNT(*) FROM thunderdome.team_checkin) AS team_checkins_count,
    (SELECT COUNT(*) FROM thunderdome.retro) AS retro_count,
    (SELECT COUNT(DISTINCT retro_id) FROM thunderdome.retro_user WHERE active IS true) AS active_retro_count,
    (SELECT COUNT(user_id) FROM thunderdome.retro_user WHERE active IS true) AS active_retro_user_count,
    (SELECT COUNT(*) FROM thunderdome.retro_item) AS retro_item_count,
    (SELECT COUNT(*) FROM thunderdome.retro_action) AS retro_action_count,
    (SELECT COUNT(*) FROM thunderdome.storyboard) AS storyboard_count,
    (SELECT COUNT(DISTINCT storyboard_id) FROM thunderdome.storyboard_user WHERE active IS true) AS active_storyboard_count,
    (SELECT COUNT(user_id) FROM thunderdome.storyboard_user WHERE active IS true) AS active_storyboard_user_count,
    (SELECT COUNT(*) FROM thunderdome.storyboard_goal) AS storyboard_goal_count,
    (SELECT COUNT(*) FROM thunderdome.storyboard_column) AS storyboard_column_count,
    (SELECT COUNT(*) FROM thunderdome.storyboard_story) AS storyboard_story_count,
    (SELECT COUNT(*) FROM thunderdome.storyboard_persona) AS storyboard_persona_count,
    (SELECT COUNT(*) FROM thunderdome.estimation_scale) AS estimation_scale_count,
    (SELECT COUNT(*) FROM thunderdome.estimation_scale WHERE estimation_scale.is_public IS TRUE) AS public_estimation_scale_count,
    (SELECT COUNT(*) FROM thunderdome.estimation_scale WHERE organization_id IS NOT NULL) AS organization_estimation_scale_count,
    (SELECT COUNT(*) FROM thunderdome.estimation_scale WHERE team_id IS NOT NULL) AS team_estimation_scale_count,
    (SELECT COUNT(*) FROM thunderdome.subscription WHERE expires > CURRENT_TIMESTAMP AND active IS TRUE AND team_id IS NULL AND organization_id IS NULL) as user_sub_count,
    (SELECT COUNT(*) FROM thunderdome.subscription WHERE expires > CURRENT_TIMESTAMP AND active IS TRUE AND team_id IS NOT NULL) as team_sub_count,
    (SELECT COUNT(*) FROM thunderdome.subscription WHERE expires > CURRENT_TIMESTAMP AND active IS TRUE AND organization_id IS NOT NULL) as org_sub_count,
    (SELECT COUNT(*) FROM thunderdome.retro_template) AS retro_template_count,
    (SELECT COUNT(*) FROM thunderdome.retro_template WHERE retro_template.is_public IS TRUE) AS public_retro_template_count,
    (SELECT COUNT(*) FROM thunderdome.retro_template WHERE organization_id IS NOT NULL) AS organization_retro_template_count,
    (SELECT COUNT(*) FROM thunderdome.retro_template WHERE team_id IS NOT NULL) AS team_retro_template_count,
	(SELECT COUNT(*) FROM thunderdome.project) AS project_count,
	(SELECT COUNT(*) FROM thunderdome.support_ticket WHERE resolved_at IS NULL) AS open_support_ticket_count
		;`,
	).Scan(
		&appStats.UnregisteredCount,
		&appStats.RegisteredCount,
		&appStats.PokerCount,
		&appStats.PokerStoryCount,
		&appStats.OrganizationCount,
		&appStats.DepartmentCount,
		&appStats.TeamCount,
		&appStats.APIKeyCount,
		&appStats.ActivePokerCount,
		&appStats.ActivePokerUserCount,
		&appStats.TeamCheckinsCount,
		&appStats.RetroCount,
		&appStats.ActiveRetroCount,
		&appStats.ActiveRetroUserCount,
		&appStats.RetroItemCount,
		&appStats.RetroActionCount,
		&appStats.StoryboardCount,
		&appStats.ActiveStoryboardCount,
		&appStats.ActiveStoryboardUserCount,
		&appStats.StoryboardGoalCount,
		&appStats.StoryboardColumnCount,
		&appStats.StoryboardStoryCount,
		&appStats.StoryboardPersonaCount,
		&appStats.EstimationScaleCount,
		&appStats.PublicEstimationScaleCount,
		&appStats.OrganizationEstimationScaleCount,
		&appStats.TeamEstimationScaleCount,
		&appStats.UserSubscriptionActiveCount,
		&appStats.TeamSubscriptionActiveCount,
		&appStats.OrgSubscriptionActiveCount,
		&appStats.RetroTemplateCount,
		&appStats.PublicRetroTemplateCount,
		&appStats.OrganizationRetroTemplateCount,
		&appStats.TeamRetroTemplateCount,
		&appStats.ProjectCount,
		&appStats.OpenSupportTicketCount,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get application stats: %v", err)
	}

	return &appStats, nil
}

// ListAdminUsers gets a list of all admin users
func (d *Service) ListAdminUsers(ctx context.Context, limit int, offset int) ([]*thunderdome.User, int, error) {
	var users = make([]*thunderdome.User, 0)
	var count int

	err := d.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM thunderdome.users WHERE type = 'ADMIN';",
	).Scan(
		&count,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("list admin users query error", zap.Error(err))
	}

	rows, err := d.DB.QueryContext(ctx,
		`
		SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''),
		 COALESCE(u.company, ''), COALESCE(u.job_title, ''), u.disabled, COALESCE(u.picture, '')
		FROM thunderdome.users u
		WHERE u.type = 'ADMIN'
		ORDER BY u.created_date
		LIMIT $1
		OFFSET $2;`,
		limit,
		offset,
	)
	if err != nil {
		return nil, count, fmt.Errorf("list admin users query error: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var w thunderdome.User

		if err := rows.Scan(
			&w.ID,
			&w.Name,
			&w.Email,
			&w.Type,
			&w.Avatar,
			&w.Verified,
			&w.Country,
			&w.Company,
			&w.JobTitle,
			&w.Disabled,
			&w.Picture,
		); err != nil {
			d.Logger.Ctx(ctx).Error("list admin users query scan error", zap.Error(err))
		} else {
			w.GravatarHash = db.CreateGravatarHash(w.Email)
			users = append(users, &w)
		}
	}

	return users, count, nil
}
