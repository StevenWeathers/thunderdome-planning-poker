package main

import (
	"embed"
	"github.com/StevenWeathers/thunderdome-planning-poker/db"
	api "github.com/StevenWeathers/thunderdome-planning-poker/http"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/spf13/viper"
	"io/fs"
	"net/http"
	"os"
)

//go:embed dist
var f embed.FS

func (s *server) getFileSystem(useOS bool) (http.FileSystem, fs.FS) {
	if useOS {
		s.logger.Info("using live mode")
		return http.FS(os.DirFS("dist")), os.DirFS("dist")
	}

	fsys, err := fs.Sub(f, "dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys), fsys
}

func (s *server) routes() {
	HFS, FSS := s.getFileSystem(embedUseOS)

	httpConfig := &api.Config{
		AppDomain:                 s.config.AppDomain,
		FrontendCookieName:        s.config.FrontendCookieName,
		SecureCookieName:          viper.GetString("http.backend_cookie_name"),
		SecureCookieFlag:          viper.GetBool("http.secure_cookie"),
		SessionCookieName:         viper.GetString("http.session_cookie_name"),
		PathPrefix:                s.config.PathPrefix,
		ExternalAPIEnabled:        s.config.ExternalAPIEnabled,
		ExternalAPIVerifyRequired: viper.GetBool("config.external_api_verify_required"),
		UserAPIKeyLimit:           s.config.UserAPIKeyLimit,
		LdapEnabled:               s.config.LdapEnabled,
		HeaderAuthEnabled:         s.config.HeaderAuthEnabled,
		FeaturePoker:              viper.GetBool("feature.poker"),
		FeatureRetro:              viper.GetBool("feature.retro"),
		FeatureStoryboard:         viper.GetBool("feature.storyboard"),
		OrganizationsEnabled:      viper.GetBool("config.organizations_enabled"),
		AvatarService:             s.config.AvatarService,
		EmbedUseOS:                embedUseOS,
	}

	appConfig := thunderdome.AppConfig{
		AllowedPointValues:        viper.GetStringSlice("config.allowedPointValues"),
		DefaultPointValues:        viper.GetStringSlice("config.defaultPointValues"),
		ShowWarriorRank:           viper.GetBool("config.show_warrior_rank"),
		AvatarService:             viper.GetString("config.avatar_service"),
		ToastTimeout:              viper.GetInt("config.toast_timeout"),
		AllowGuests:               viper.GetBool("config.allow_guests"),
		AllowRegistration:         viper.GetBool("config.allow_registration") && viper.GetString("auth.method") == "normal",
		AllowJiraImport:           viper.GetBool("config.allow_jira_import"),
		AllowCsvImport:            viper.GetBool("config.allow_csv_import"),
		DefaultLocale:             viper.GetString("config.default_locale"),
		FriendlyUIVerbs:           viper.GetBool("config.friendly_ui_verbs"),
		OrganizationsEnabled:      viper.GetBool("config.organizations_enabled"),
		ExternalAPIEnabled:        s.config.ExternalAPIEnabled,
		UserAPIKeyLimit:           s.config.UserAPIKeyLimit,
		AppVersion:                s.config.Version,
		CookieName:                s.config.FrontendCookieName,
		PathPrefix:                s.config.PathPrefix,
		CleanupGuestsDaysOld:      viper.GetInt("config.cleanup_guests_days_old"),
		CleanupBattlesDaysOld:     viper.GetInt("config.cleanup_battles_days_old"),
		CleanupRetrosDaysOld:      viper.GetInt("config.cleanup_retros_days_old"),
		CleanupStoryboardsDaysOld: viper.GetInt("config.cleanup_storyboards_days_old"),
		ShowActiveCountries:       viper.GetBool("config.show_active_countries"),
		LdapEnabled:               s.config.LdapEnabled,
		HeaderAuthEnabled:         s.config.HeaderAuthEnabled,
		FeaturePoker:              viper.GetBool("feature.poker"),
		FeatureRetro:              viper.GetBool("feature.retro"),
		FeatureStoryboard:         viper.GetBool("feature.storyboard"),
		RequireTeams:              viper.GetBool("config.require_teams"),
	}

	uiConfig := thunderdome.UIConfig{
		AnalyticsEnabled: s.config.AnalyticsEnabled,
		AnalyticsID:      s.config.AnalyticsID,
		AppConfig:        appConfig,
	}

	// Create services.
	userService := &db.UserService{DB: s.db.DB, Logger: s.logger}
	apkService := &db.APIKeyService{DB: s.db.DB, Logger: s.logger}
	s.AlertService = &db.AlertService{DB: s.db.DB, Logger: s.logger}
	authService := &db.AuthService{DB: s.db.DB, Logger: s.logger, AESHashkey: s.db.Config.AESHashkey}
	battleService := &db.BattleService{
		DB: s.db.DB, Logger: s.logger, AESHashKey: s.db.Config.AESHashkey,
		HTMLSanitizerPolicy: s.db.HTMLSanitizerPolicy,
	}
	checkinService := &db.CheckinService{DB: s.db.DB, Logger: s.logger, HTMLSanitizerPolicy: s.db.HTMLSanitizerPolicy}
	retroService := &db.RetroService{DB: s.db.DB, Logger: s.logger, AESHashKey: s.db.Config.AESHashkey}
	storyboardService := &db.StoryboardService{DB: s.db.DB, Logger: s.logger, AESHashKey: s.db.Config.AESHashkey}
	teamService := &db.TeamService{DB: s.db.DB, Logger: s.logger}
	organizationService := &db.OrganizationService{DB: s.db.DB, Logger: s.logger}
	adminService := &db.AdminService{DB: s.db.DB, Logger: s.logger}

	a := api.Service{
		Config:              httpConfig,
		Router:              s.router,
		DB:                  s.db,
		Email:               s.email,
		Cookie:              s.cookie,
		Logger:              s.logger,
		UserService:         userService,
		APIKeyService:       apkService,
		AlertService:        s.AlertService,
		AuthService:         authService,
		BattleService:       battleService,
		CheckinService:      checkinService,
		RetroService:        retroService,
		StoryboardService:   storyboardService,
		TeamService:         teamService,
		OrganizationService: organizationService,
		AdminService:        adminService,
		UIConfig:            uiConfig,
	}

	api.Init(a, FSS, HFS)
}
