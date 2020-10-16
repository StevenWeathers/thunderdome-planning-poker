package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/email"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/spf13/viper"
)

var (
	version = "dev"
)

// ServerConfig holds server global config values
type ServerConfig struct {
	// port the application server will listen on
	ListenPort string
	// the domain of the application for cookie securing
	AppDomain string
	// name of the frontend cookie
	FrontendCookieName string
	// name of the warrior cookie
	SecureCookieName string
	// controls whether or not the cookie is set to secure, only works over HTTPS
	SecureCookieFlag bool
	// email to promote a warrior to GENERAL on app startup
	// the warrior should already be registered for this to work
	AdminEmail string
	// Whether or not to enable google analytics tracking
	AnalyticsEnabled bool
	// ID used for google analytics
	AnalyticsID string
	// the app version
	Version string
	// Which avatar service is utilized
	AvatarService string
}

type server struct {
	config   *ServerConfig
	router   *mux.Router
	email    *email.Email
	cookie   *securecookie.SecureCookie
	database *database.Database
}

func main() {
	fmt.Printf("Thunderdome version %s", version)

	InitConfig()

	var cookieHashkey = viper.GetString("http.cookie_hashkey")

	s := &server{
		config: &ServerConfig{
			ListenPort:         viper.GetString("http.port"),
			AppDomain:          viper.GetString("http.domain"),
			AdminEmail:         viper.GetString("admin.email"),
			FrontendCookieName: "warrior",
			SecureCookieName:   "warriorId",
			SecureCookieFlag:   viper.GetBool("http.secure_cookie"),
			AnalyticsEnabled:   viper.GetBool("analytics.enabled"),
			AnalyticsID:        viper.GetString("analytics.id"),
			Version:            version,
			AvatarService:      viper.GetString(("config.avatar_service")),
		},
		router: mux.NewRouter(),
		cookie: securecookie.New([]byte(cookieHashkey), nil),
	}
	s.email = email.New(s.config.AppDomain)
	s.database = database.New(s.config.AdminEmail)

	go h.run()

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
