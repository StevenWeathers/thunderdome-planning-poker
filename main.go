package main

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/db/admin"
	"github.com/StevenWeathers/thunderdome-planning-poker/db/alert"
	"github.com/StevenWeathers/thunderdome-planning-poker/db/apikey"
	"github.com/StevenWeathers/thunderdome-planning-poker/db/auth"
	"github.com/StevenWeathers/thunderdome-planning-poker/db/poker"
	"github.com/StevenWeathers/thunderdome-planning-poker/db/retro"
	"github.com/StevenWeathers/thunderdome-planning-poker/db/storyboard"
	"github.com/StevenWeathers/thunderdome-planning-poker/db/team"
	"github.com/StevenWeathers/thunderdome-planning-poker/db/user"
	api "github.com/StevenWeathers/thunderdome-planning-poker/http"
	"github.com/StevenWeathers/thunderdome-planning-poker/ui"

	"github.com/StevenWeathers/thunderdome-planning-poker/config"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"

	"github.com/StevenWeathers/thunderdome-planning-poker/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/email"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var embedUseOS bool
var (
	version = "dev"
)

// Config holds server global config values
type Config struct {
	// port the application server will listen on
	ListenPort string
	// the domain of the application for cookie securing
	AppDomain string
	// name of the cookie used exclusively by the UI
	FrontendCookieName string
	// email to promote a user to Admin type on app startup
	// the user should already be registered for this to work
	AdminEmail string
	// Whether to enable Google Analytics tracking
	AnalyticsEnabled bool
	// ID used for Google Analytics
	AnalyticsID string
	// the app version
	Version string
	// Which avatar service is utilized
	AvatarService string
	// PathPrefix allows the application to be run on a shared domain
	PathPrefix string
	// Whether the external API is enabled
	ExternalAPIEnabled bool
	// Number of API keys a user can create
	UserAPIKeyLimit int
	// Whether LDAP is enabled for authentication
	LdapEnabled bool
	// Whether header authentication is enabled
	HeaderAuthEnabled bool
}

type server struct {
	config       *Config
	router       *mux.Router
	email        thunderdome.EmailService
	cookie       *securecookie.SecureCookie
	db           *db.Service
	logger       *otelzap.Logger
	AlertService thunderdome.AlertDataSvc
}

func main() {
	zlog, _ := zap.NewProduction(
		zap.Fields(
			zap.String("version", version),
		),
	)
	defer func() {
		_ = zlog.Sync()
	}()
	logger := otelzap.New(zlog)

	embedUseOS = len(os.Args) > 1 && os.Args[1] == "live"

	c := config.InitConfig(logger)

	if c.Otel.Enabled {
		cleanup := initTracer(
			logger,
			c.Otel.ServiceName,
			c.Otel.CollectorUrl,
			c.Otel.InsecureMode,
		)
		defer func() {
			_ = cleanup(context.Background())
		}()
	}

	cookieHashKey := c.Http.CookieHashkey
	pathPrefix := c.Http.PathPrefix
	router := mux.NewRouter()

	if pathPrefix != "" {
		router = router.PathPrefix(pathPrefix).Subrouter()
	}

	router.Use(otelmux.Middleware("thunderdome"))

	s := &server{
		config: &Config{
			ListenPort:         c.Http.Port,
			AppDomain:          c.Http.Domain,
			AdminEmail:         c.Admin.Email,
			FrontendCookieName: c.Http.FrontendCookieName,
			AnalyticsEnabled:   c.Analytics.Enabled,
			AnalyticsID:        c.Analytics.ID,
			Version:            version,
			AvatarService:      c.Config.AvatarService,
			PathPrefix:         pathPrefix,
			ExternalAPIEnabled: c.Config.AllowExternalApi,
			UserAPIKeyLimit:    c.Config.UserApikeyLimit,
			LdapEnabled:        c.Auth.Method == "ldap",
			HeaderAuthEnabled:  c.Auth.Method == "header",
		},
		router: router,
		cookie: securecookie.New([]byte(cookieHashKey), nil),
		logger: logger,
	}

	s.email = email.New(&email.Config{
		AppURL:       "https://" + c.Http.Domain + c.Http.PathPrefix + "/",
		SenderName:   "Thunderdome",
		SmtpEnabled:  c.Smtp.Enabled,
		SmtpHost:     c.Smtp.Host,
		SmtpPort:     c.Smtp.Port,
		SmtpSecure:   c.Smtp.Secure,
		SmtpIdentity: c.Smtp.Identity,
		SmtpUser:     c.Smtp.User,
		SmtpPass:     c.Smtp.Pass,
		SmtpSender:   c.Smtp.Sender,
	}, s.logger)
	s.db = db.New(s.config.AdminEmail, &db.Config{
		Host:            c.Db.Host,
		Port:            c.Db.Port,
		User:            c.Db.User,
		Password:        c.Db.Pass,
		Name:            c.Db.Name,
		SSLMode:         c.Db.Sslmode,
		AESHashkey:      c.Config.AesHashkey,
		MaxIdleConns:    c.Db.MaxIdleConns,
		MaxOpenConns:    c.Db.MaxOpenConns,
		ConnMaxLifetime: c.Db.ConnMaxLifetime,
	}, s.logger)

	HFS, FSS := ui.New(embedUseOS)

	httpConfig := &api.Config{
		AppDomain:                 s.config.AppDomain,
		FrontendCookieName:        s.config.FrontendCookieName,
		SecureCookieName:          c.Http.BackendCookieName,
		SecureCookieFlag:          c.Http.SecureCookie,
		SessionCookieName:         c.Http.SessionCookieName,
		PathPrefix:                s.config.PathPrefix,
		ExternalAPIEnabled:        s.config.ExternalAPIEnabled,
		ExternalAPIVerifyRequired: c.Config.ExternalApiVerifyRequired,
		UserAPIKeyLimit:           s.config.UserAPIKeyLimit,
		LdapEnabled:               s.config.LdapEnabled,
		HeaderAuthEnabled:         s.config.HeaderAuthEnabled,
		FeaturePoker:              c.Feature.Poker,
		FeatureRetro:              c.Feature.Retro,
		FeatureStoryboard:         c.Feature.Storyboard,
		OrganizationsEnabled:      c.Config.OrganizationsEnabled,
		AvatarService:             s.config.AvatarService,
		EmbedUseOS:                embedUseOS,
	}

	appConfig := thunderdome.AppConfig{
		AllowedPointValues:        c.Config.AllowedPointValues,
		DefaultPointValues:        c.Config.DefaultPointValues,
		ShowWarriorRank:           c.Config.ShowWarriorRank,
		AvatarService:             c.Config.AvatarService,
		ToastTimeout:              c.Config.ToastTimeout,
		AllowGuests:               c.Config.AllowGuests,
		AllowRegistration:         c.Config.AllowRegistration && c.Auth.Method == "normal",
		AllowJiraImport:           c.Config.AllowJiraImport,
		AllowCsvImport:            c.Config.AllowCsvImport,
		DefaultLocale:             c.Config.DefaultLocale,
		FriendlyUIVerbs:           c.Config.FriendlyUiVerbs,
		OrganizationsEnabled:      c.Config.OrganizationsEnabled,
		ExternalAPIEnabled:        s.config.ExternalAPIEnabled,
		UserAPIKeyLimit:           s.config.UserAPIKeyLimit,
		AppVersion:                s.config.Version,
		CookieName:                s.config.FrontendCookieName,
		PathPrefix:                s.config.PathPrefix,
		CleanupGuestsDaysOld:      c.Config.CleanupGuestsDaysOld,
		CleanupBattlesDaysOld:     c.Config.CleanupBattlesDaysOld,
		CleanupRetrosDaysOld:      c.Config.CleanupRetrosDaysOld,
		CleanupStoryboardsDaysOld: c.Config.CleanupStoryboardsDaysOld,
		ShowActiveCountries:       c.Config.ShowActiveCountries,
		LdapEnabled:               s.config.LdapEnabled,
		HeaderAuthEnabled:         s.config.HeaderAuthEnabled,
		FeaturePoker:              c.Feature.Poker,
		FeatureRetro:              c.Feature.Retro,
		FeatureStoryboard:         c.Feature.Storyboard,
		RequireTeams:              c.Config.RequireTeams,
	}

	uiConfig := thunderdome.UIConfig{
		AnalyticsEnabled: s.config.AnalyticsEnabled,
		AnalyticsID:      s.config.AnalyticsID,
		AppConfig:        appConfig,
	}

	// Create services.
	userService := &user.Service{DB: s.db.DB, Logger: s.logger}
	apkService := &apikey.Service{DB: s.db.DB, Logger: s.logger}
	s.AlertService = &alert.Service{DB: s.db.DB, Logger: s.logger}
	authService := &auth.Service{DB: s.db.DB, Logger: s.logger, AESHashkey: s.db.Config.AESHashkey}
	battleService := &poker.Service{
		DB: s.db.DB, Logger: s.logger, AESHashKey: s.db.Config.AESHashkey,
		HTMLSanitizerPolicy: s.db.HTMLSanitizerPolicy,
	}
	checkinService := &team.CheckinService{DB: s.db.DB, Logger: s.logger, HTMLSanitizerPolicy: s.db.HTMLSanitizerPolicy}
	retroService := &retro.Service{DB: s.db.DB, Logger: s.logger, AESHashKey: s.db.Config.AESHashkey}
	storyboardService := &storyboard.Service{DB: s.db.DB, Logger: s.logger, AESHashKey: s.db.Config.AESHashkey}
	teamService := &team.Service{DB: s.db.DB, Logger: s.logger}
	organizationService := &team.OrganizationService{DB: s.db.DB, Logger: s.logger}
	adminService := &admin.Service{DB: s.db.DB, Logger: s.logger}

	a := api.Service{
		Config:              httpConfig,
		Router:              s.router,
		Email:               s.email,
		Cookie:              s.cookie,
		Logger:              s.logger,
		UserDataSvc:         userService,
		ApiKeyDataSvc:       apkService,
		AlertDataSvc:        s.AlertService,
		AuthDataSvc:         authService,
		PokerDataSvc:        battleService,
		CheckinDataSvc:      checkinService,
		RetroDataSvc:        retroService,
		StoryboardDataSvc:   storyboardService,
		TeamDataSvc:         teamService,
		OrganizationDataSvc: organizationService,
		AdminDataSvc:        adminService,
		UIConfig:            uiConfig,
	}

	api.Init(a, FSS, HFS)

	srv := &http.Server{
		Handler:           s.router,
		Addr:              fmt.Sprintf(":%s", s.config.ListenPort),
		WriteTimeout:      time.Duration(c.Http.WriteTimeout) * time.Second,
		ReadTimeout:       time.Duration(c.Http.ReadTimeout) * time.Second,
		IdleTimeout:       time.Duration(c.Http.IdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(c.Http.ReadHeaderTimeout) * time.Second,
	}

	s.logger.Info("Access the WebUI via 127.0.0.1:" + s.config.ListenPort)

	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Fatal(err.Error())
	}
}

func initTracer(logger *otelzap.Logger, serviceName string, collectorURL string, insecure bool) func(context.Context) error {
	logger.Ctx(context.Background()).Info("initializing open telemetry")
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		logger.Ctx(context.Background()).Fatal("error initializing tracer", zap.Error(err))
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Ctx(context.Background()).Error("Could not set resources: ", zap.Error(err))
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}
