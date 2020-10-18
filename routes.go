package main

import (
	"net/http"

	"github.com/markbates/pkger"
	"github.com/spf13/viper"
)

func (s *server) routes() {
	staticHandler := http.FileServer(pkger.Dir("/dist"))
	// static assets
	s.router.PathPrefix("/static/").Handler(http.StripPrefix(s.config.PathPrefix, staticHandler))
	s.router.PathPrefix("/img/").Handler(http.StripPrefix(s.config.PathPrefix, staticHandler))
	s.router.PathPrefix("/lang/").Handler(http.StripPrefix(s.config.PathPrefix, staticHandler))
	// warrior avatar generation
	if s.config.AvatarService == "goadorable" || s.config.AvatarService == "govatar" {
		s.router.PathPrefix("/avatar/{width}/{id}/{avatar}").Handler(s.handleWarriorAvatar()).Methods("GET")
		s.router.PathPrefix("/avatar/{width}/{id}").Handler(s.handleWarriorAvatar()).Methods("GET")
	}
	// api (currently internal to UI application)
	// warrior authentication, profile
	if viper.GetString("auth.method") == "ldap" {
		s.router.HandleFunc("/api/auth", s.handleLdapLogin()).Methods("POST")
	} else {
		s.router.HandleFunc("/api/auth", s.handleLogin()).Methods("POST")
		s.router.HandleFunc("/api/auth/forgot-password", s.handleForgotPassword()).Methods("POST")
		s.router.HandleFunc("/api/auth/reset-password", s.handleResetPassword()).Methods("POST")
		s.router.HandleFunc("/api/auth/update-password", s.handleUpdatePassword()).Methods("POST")
		s.router.HandleFunc("/api/auth/verify", s.handleAccountVerification()).Methods("POST")
		s.router.HandleFunc("/api/enlist", s.handleWarriorEnlist()).Methods("POST")
	}
	s.router.HandleFunc("/api/warrior", s.handleWarriorRecruit()).Methods("POST")
	s.router.HandleFunc("/api/auth/logout", s.handleLogout()).Methods("POST")
	s.router.HandleFunc("/api/warrior/{id}", s.handleWarriorProfile()).Methods("GET")
	s.router.HandleFunc("/api/warrior/{id}", s.handleWarriorProfileUpdate()).Methods("POST")
	// battle(s)
	s.router.HandleFunc("/api/battle", s.handleBattleCreate()).Methods("POST")
	s.router.HandleFunc("/api/battles", s.handleBattlesGet())
	// admin routes
	s.router.HandleFunc("/api/admin/stats", s.adminOnly(s.handleAppStats()))
	s.router.HandleFunc("/api/admin/warriors", s.adminOnly(s.handleGetRegisteredWarriors()))
	s.router.HandleFunc("/api/admin/warrior", s.adminOnly(s.handleWarriorCreate())).Methods("POST")
	// websocket for battle
	s.router.HandleFunc("/api/arena/{id}", s.serveWs())
	// handle index.html
	s.router.PathPrefix("/").HandlerFunc(s.handleIndex())
}
