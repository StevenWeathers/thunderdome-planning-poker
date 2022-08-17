package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"net/http"
	"os"
	"time"

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
	"github.com/spf13/viper"
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
}

type server struct {
	config *Config
	router *mux.Router
	email  *email.Email
	cookie *securecookie.SecureCookie
	db     *db.Database
	logger *otelzap.Logger
}

func main() {
	zlog, _ := zap.NewProduction(
		zap.Fields(
			zap.String("version", version),
		),
	)
	defer zlog.Sync()
	logger := otelzap.New(zlog)

	embedUseOS = len(os.Args) > 1 && os.Args[1] == "live"

	InitConfig(logger)

	if viper.GetBool("otel.enabled") {
		cleanup := initTracer(
			logger,
			viper.GetString("otel.service_name"),
			viper.GetString("otel.collector_url"),
			viper.GetBool("otel.insecure_mode"),
		)
		defer cleanup(context.Background())
	}

	cookieHashkey := viper.GetString("http.cookie_hashkey")
	pathPrefix := viper.GetString("http.path_prefix")
	router := mux.NewRouter()

	if pathPrefix != "" {
		router = router.PathPrefix(pathPrefix).Subrouter()
	}

	router.Use(otelmux.Middleware("thunderdome"))

	s := &server{
		config: &Config{
			ListenPort:         viper.GetString("http.port"),
			AppDomain:          viper.GetString("http.domain"),
			AdminEmail:         viper.GetString("admin.email"),
			FrontendCookieName: viper.GetString("http.frontend_cookie_name"),
			AnalyticsEnabled:   viper.GetBool("analytics.enabled"),
			AnalyticsID:        viper.GetString("analytics.id"),
			Version:            version,
			AvatarService:      viper.GetString(("config.avatar_service")),
			PathPrefix:         pathPrefix,
			ExternalAPIEnabled: viper.GetBool("config.allow_external_api"),
			UserAPIKeyLimit:    viper.GetInt("config.user_apikey_limit"),
			LdapEnabled:        viper.GetString("auth.method") == "ldap",
		},
		router: router,
		cookie: securecookie.New([]byte(cookieHashkey), nil),
		logger: logger,
	}

	s.email = email.New(s.config.AppDomain, s.config.PathPrefix, s.logger)
	s.db = db.New(s.config.AdminEmail, &db.Config{
		Host:            viper.GetString("db.host"),
		Port:            viper.GetInt("db.port"),
		User:            viper.GetString("db.user"),
		Password:        viper.GetString("db.pass"),
		Name:            viper.GetString("db.name"),
		SSLMode:         viper.GetString("db.sslmode"),
		AESHashkey:      viper.GetString("config.aes_hashkey"),
		MaxIdleConns:    viper.GetInt("db.max_idle_conns"),
		MaxOpenConns:    viper.GetInt("db.max_open_conns"),
		ConnMaxLifetime: viper.GetInt("db.conn_max_lifetime"),
	}, s.logger)

	s.routes()

	srv := &http.Server{
		Handler:     s.router,
		Addr:        fmt.Sprintf(":%s", s.config.ListenPort),
		ReadTimeout: 15 * time.Second,
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
