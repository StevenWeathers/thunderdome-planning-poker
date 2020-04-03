package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/markbates/pkger"
)

// AppDomain is the domain of the application for cookie securing
var AppDomain string

// SecureCookieHashkey is used to hash the secure cookie
var SecureCookieHashkey []byte

// SecureCookieName is obviously the name of the secure cookie
var SecureCookieName = "warriorId"

// SecureCookieFlag controls whether or not the cookie is set to secure, only works over HTTPS
var SecureCookieFlag bool

// Sc is the secure cookie instance with secret hash
var Sc = securecookie.New([]byte("some-secret"), nil)

// AdminEmail is used to promote a warrior to GENERAL on app startup
// the warrior should already be registered for this to work
var AdminEmail string

func main() {
	var listenPort = fmt.Sprintf(":%s", GetEnv("PORT", "8080"))

	AppDomain = GetEnv("APP_DOMAIN", "thunderdome.dev")
	AdminEmail = GetEnv("ADMIN_EMAIL", "")
	SecureCookieHashkey = []byte(GetEnv("COOKIE_HASHKEY", "strongest-avenger"))
	SecureCookieFlag = GetBoolEnv("COOKIE_SECURE", true)
	Sc = securecookie.New(SecureCookieHashkey, nil)

	SetupDB() // Sets up DB Connection, and if necessary Tables

	GetMailserverConfig()

	go h.run()

	staticHandler := http.FileServer(pkger.Dir("/dist"))

	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(staticHandler)
	router.PathPrefix("/js/").Handler(staticHandler)
	router.PathPrefix("/img/").Handler(staticHandler)
	router.HandleFunc("/api/auth", LoginHandler).Methods("POST")
	router.HandleFunc("/api/auth/logout", LogoutHandler).Methods("POST")
	router.HandleFunc("/api/auth/forgot-password", ForgotPasswordHandler).Methods("POST")
	router.HandleFunc("/api/auth/reset-password", ResetPasswordHandler).Methods("POST")
	router.HandleFunc("/api/auth/update-password", UpdatePasswordHandler).Methods("POST")
	router.HandleFunc("/api/auth/verify", VerifyAccountHandler).Methods("POST")
	router.HandleFunc("/api/warrior", RecruitWarriorHandler).Methods("POST")
	router.HandleFunc("/api/enlist", EnlistWarriorHandler).Methods("POST")
	router.HandleFunc("/api/warrior/{id}", WarriorProfileHandler).Methods("GET")
	router.HandleFunc("/api/warrior/{id}", WarriorUpdateProfileHandler).Methods("POST")
	router.HandleFunc("/api/battle", CreateBattleHandler).Methods("POST")
	router.HandleFunc("/api/battle/{id}", GetBattleHandler)
	router.HandleFunc("/api/battles", GetBattlesHandler)
	router.HandleFunc("/api/admin/stats", GetAppStatsHandler)
	router.HandleFunc("/api/arena/{id}", serveWs)
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/"
		staticHandler.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    listenPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Access the WebUI via 127.0.0.1" + listenPort)

	log.Fatal(srv.ListenAndServe())
}
