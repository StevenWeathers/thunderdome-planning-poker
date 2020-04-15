package main

import (
	"net/http"

	"github.com/markbates/pkger"
)

func (s *server) routes() {
	staticHandler := http.FileServer(pkger.Dir("/dist"))
	// static assets
	s.router.PathPrefix("/css/").Handler(staticHandler)
	s.router.PathPrefix("/js/").Handler(staticHandler)
	s.router.PathPrefix("/img/").Handler(staticHandler)
	// api (currently internal to UI application)
	// warrior authentication, profile
	s.router.HandleFunc("/api/auth", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/api/auth/logout", s.handleLogout()).Methods("POST")
	s.router.HandleFunc("/api/auth/forgot-password", s.handleForgotPassword()).Methods("POST")
	s.router.HandleFunc("/api/auth/reset-password", s.handleResetPassword()).Methods("POST")
	s.router.HandleFunc("/api/auth/update-password", s.handleUpdatePassword()).Methods("POST")
	s.router.HandleFunc("/api/auth/verify", s.handleAccountVerification()).Methods("POST")
	s.router.HandleFunc("/api/warrior", s.handleWarriorRecruit()).Methods("POST")
	s.router.HandleFunc("/api/enlist", s.handleWarriorEnlist()).Methods("POST")
	s.router.HandleFunc("/api/warrior/{id}", s.handleWarriorProfile()).Methods("GET")
	s.router.HandleFunc("/api/warrior/{id}", s.handleWarriorProfileUpdate()).Methods("POST")
	// battle(s)
	s.router.HandleFunc("/api/battle", s.handleBattleCreate()).Methods("POST")
	s.router.HandleFunc("/api/battle/{id}", s.handleBattleGet())
	s.router.HandleFunc("/api/battles", s.handleBattlesGet())
	// admin routes
	s.router.HandleFunc("/api/admin/stats", s.adminOnly(s.handleAppStats()))
	// websocket for battle
	s.router.HandleFunc("/api/arena/{id}", s.serveWs())
	// handle index.html
	s.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/"
		staticHandler.ServeHTTP(w, r)
	})
}
