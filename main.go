package main

import (
	_ "embed"
	"fmt"
	"github.com/StevenWeathers/thunderdome-planning-poker/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/email"
	"log"
	"net/http"
	"os"
	"time"

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
}

func main() {
	embedUseOS = len(os.Args) > 1 && os.Args[1] == "live"
	log.Println("Thunderdome version " + version)

	InitConfig()

	cookieHashkey := viper.GetString("http.cookie_hashkey")
	pathPrefix := viper.GetString("http.path_prefix")
	router := mux.NewRouter()

	if pathPrefix != "" {
		router = router.PathPrefix(pathPrefix).Subrouter()
	}

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
	}

	s.email = email.New(s.config.AppDomain, s.config.PathPrefix)
	s.db = db.New(s.config.AdminEmail)

	s.routes()

	srv := &http.Server{
		Handler: s.router,
		Addr:    fmt.Sprintf(":%s", s.config.ListenPort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Access the WebUI via 127.0.0.1:" + s.config.ListenPort)

	log.Fatal(srv.ListenAndServe())
}
