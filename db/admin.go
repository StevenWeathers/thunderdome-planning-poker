package db

import (
	"context"
	"database/sql"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

// AdminService represents a PostgreSQL implementation of thunderdome.AdminDataSvc.
type AdminService struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GetAppStats gets counts of common application metrics such as users and battles
func (d *AdminService) GetAppStats(ctx context.Context) (*thunderdome.ApplicationStats, error) {
	var Appstats thunderdome.ApplicationStats

	err := d.DB.QueryRowContext(ctx, `
		SELECT
			unregistered_user_count,
			registered_user_count,
			poker_count,
			poker_story_count,
			organization_count,
			department_count,
			team_count,
			apikey_count,
			active_poker_count,
			active_poker_user_count,
			team_checkins_count,
			retro_count,
			active_retro_count,
			active_retro_user_count,
			retro_item_count,
			retro_action_count,
			storyboard_count,
			active_storyboard_count,
			active_storyboard_user_count,
			storyboard_goal_count,
			storyboard_column_count,
			storyboard_story_count,
			storyboard_persona_count
		FROM thunderdome.appstats_get();
		`,
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
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to get application stats", zap.Error(err))
		return nil, err
	}

	return &Appstats, nil
}
