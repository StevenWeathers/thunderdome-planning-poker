package admin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// Service represents a PostgreSQL implementation of thunderdome.AdminDataSvc.
type Service struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GetAppStats gets counts of common application metrics such as users and poker games
func (d *Service) GetAppStats(ctx context.Context) (*thunderdome.ApplicationStats, error) {
	var Appstats thunderdome.ApplicationStats

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
    (SELECT COUNT(*) FROM thunderdome.retro_template WHERE team_id IS NOT NULL) AS team_retro_template_count
		;`,
	).Scan(
		&Appstats.UnregisteredCount,
		&Appstats.RegisteredCount,
		&Appstats.PokerCount,
		&Appstats.PokerStoryCount,
		&Appstats.OrganizationCount,
		&Appstats.DepartmentCount,
		&Appstats.TeamCount,
		&Appstats.APIKeyCount,
		&Appstats.ActivePokerCount,
		&Appstats.ActivePokerUserCount,
		&Appstats.TeamCheckinsCount,
		&Appstats.RetroCount,
		&Appstats.ActiveRetroCount,
		&Appstats.ActiveRetroUserCount,
		&Appstats.RetroItemCount,
		&Appstats.RetroActionCount,
		&Appstats.StoryboardCount,
		&Appstats.ActiveStoryboardCount,
		&Appstats.ActiveStoryboardUserCount,
		&Appstats.StoryboardGoalCount,
		&Appstats.StoryboardColumnCount,
		&Appstats.StoryboardStoryCount,
		&Appstats.StoryboardPersonaCount,
		&Appstats.EstimationScaleCount,
		&Appstats.PublicEstimationScaleCount,
		&Appstats.OrganizationEstimationScaleCount,
		&Appstats.TeamEstimationScaleCount,
		&Appstats.UserSubscriptionActiveCount,
		&Appstats.TeamSubscriptionActiveCount,
		&Appstats.OrgSubscriptionActiveCount,
		&Appstats.RetroTemplateCount,
		&Appstats.PublicRetroTemplateCount,
		&Appstats.OrganizationRetroTemplateCount,
		&Appstats.TeamRetroTemplateCount,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get application stats: %v", err)
	}

	return &Appstats, nil
}
