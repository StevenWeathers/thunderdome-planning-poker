package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
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
}

type server struct {
	config *ServerConfig
	router *mux.Router
	email  *Email
	cookie *securecookie.SecureCookie
}

func main() {
	var cookieHashkey = GetEnv("COOKIE_HASHKEY", "strongest-avenger")

	s := &server{
		config: &ServerConfig{
			ListenPort:         GetEnv("PORT", "8080"),
			AppDomain:          GetEnv("APP_DOMAIN", "thunderdome.dev"),
			AdminEmail:         GetEnv("ADMIN_EMAIL", ""),
			FrontendCookieName: "warrior",
			SecureCookieName:   "warriorId",
			SecureCookieFlag:   GetBoolEnv("COOKIE_SECURE", true),
			AnalyticsEnabled:   GetBoolEnv("ANALYTICS_ENABLED", true),
			AnalyticsID:        GetEnv("ANALYTICS_ID", "UA-140245309-1"),
		},
		router: mux.NewRouter(),
		cookie: securecookie.New([]byte(cookieHashkey), nil),
	}
	s.email = NewEmail(s.config.AppDomain)

	SetupDB(s.config.AdminEmail)

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
