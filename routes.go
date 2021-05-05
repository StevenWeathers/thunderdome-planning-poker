package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/spf13/viper"
)

//go:embed dist
var f embed.FS

func (s *server) routes() {
	fsys, err := fs.Sub(f, "dist")
	if err != nil {
		panic(err)
	}
	staticHandler := http.FileServer(http.FS(fsys))

	// static assets
	s.router.PathPrefix("/static/").Handler(http.StripPrefix(s.config.PathPrefix, staticHandler))
	s.router.PathPrefix("/img/").Handler(http.StripPrefix(s.config.PathPrefix, staticHandler))
	s.router.PathPrefix("/lang/").Handler(http.StripPrefix(s.config.PathPrefix, staticHandler))
	// user avatar generation
	if s.config.AvatarService == "goadorable" || s.config.AvatarService == "govatar" {
		s.router.PathPrefix("/avatar/{width}/{id}/{avatar}").Handler(s.handleUserAvatar()).Methods("GET")
		s.router.PathPrefix("/avatar/{width}/{id}").Handler(s.handleUserAvatar()).Methods("GET")
	}
	// api
	// user authentication, profile
	if viper.GetString("auth.method") == "ldap" {
		s.router.HandleFunc("/api/auth", s.handleLdapLogin()).Methods("POST")
	} else {
		s.router.HandleFunc("/api/auth", s.handleLogin()).Methods("POST")
		s.router.HandleFunc("/api/auth/forgot-password", s.handleForgotPassword()).Methods("POST")
		s.router.HandleFunc("/api/auth/reset-password", s.handleResetPassword()).Methods("POST")
		s.router.HandleFunc("/api/auth/update-password", s.userOnly(s.handleUpdatePassword())).Methods("POST")
		s.router.HandleFunc("/api/auth/verify", s.handleAccountVerification()).Methods("POST")
		s.router.HandleFunc("/api/enlist", s.handleUserEnlist()).Methods("POST")
	}
	s.router.HandleFunc("/api/warrior", s.handleUserRecruit()).Methods("POST")
	s.router.HandleFunc("/api/auth/logout", s.handleLogout()).Methods("POST")
	s.router.HandleFunc("/api/warrior/{id}/apikey/{keyID}", s.userOnly(s.handleUserAPIKeyUpdate())).Methods("PUT")
	s.router.HandleFunc("/api/warrior/{id}/apikey/{keyID}", s.userOnly(s.handleUserAPIKeyDelete())).Methods("DELETE")
	s.router.HandleFunc("/api/warrior/{id}/apikey", s.userOnly(s.handleAPIKeyGenerate())).Methods("POST")
	s.router.HandleFunc("/api/warrior/{id}/apikeys", s.userOnly(s.handleUserAPIKeys())).Methods("GET")
	s.router.HandleFunc("/api/warrior/{id}", s.userOnly(s.handleUserProfile())).Methods("GET")
	s.router.HandleFunc("/api/warrior/{id}", s.userOnly(s.handleUserProfileUpdate())).Methods("POST")
	s.router.HandleFunc("/api/warrior/{id}", s.userOnly(s.handleUserDelete())).Methods("DELETE")
	// battle(s)
	s.router.HandleFunc("/api/battle", s.userOnly(s.handleBattleCreate())).Methods("POST")
	s.router.HandleFunc("/api/battles", s.userOnly(s.handleBattlesGet())).Methods("GET")
	// organization(s)
	s.router.HandleFunc("/api/organizations/{limit}/{offset}", s.userOnly(s.handleGetOrganizationsByUser())).Methods("GET")
	s.router.HandleFunc("/api/organizations", s.userOnly(s.handleCreateOrganization())).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/departments/{limit}/{offset}", s.userOnly(s.orgUserOnly(s.handleGetOrganizationDepartments()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/departments", s.userOnly(s.orgAdminOnly(s.handleCreateDepartment()))).Methods("POST")
	// org departments(s)
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/teams/{limit}/{offset}", s.userOnly(s.departmentUserOnly(s.handleGetDepartmentTeams()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/teams", s.userOnly(s.departmentAdminOnly(s.handleCreateDepartmentTeam()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/users/{limit}/{offset}", s.userOnly(s.departmentUserOnly(s.handleGetDepartmentUsers()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/users", s.userOnly(s.departmentAdminOnly(s.handleDepartmentAddUser()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}/battles/{limit}/{offset}", s.userOnly(s.departmentTeamUserOnly(s.handleGetTeamBattles()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}/users/{limit}/{offset}", s.userOnly(s.departmentTeamUserOnly(s.handleGetTeamUsers()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}/users", s.userOnly(s.departmentTeamAdminOnly(s.handleTeamAddUser()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}/users", s.userOnly(s.departmentTeamAdminOnly(s.handleTeamRemoveUser()))).Methods("DELETE")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}", s.userOnly(s.departmentTeamUserOnly(s.handleDepartmentTeamByUser()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}", s.userOnly(s.departmentUserOnly(s.handleGetDepartmentByUser()))).Methods("GET")
	// org teams
	s.router.HandleFunc("/api/organization/{orgId}/teams/{limit}/{offset}", s.userOnly(s.orgUserOnly(s.handleGetOrganizationTeams()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/teams", s.userOnly(s.orgAdminOnly(s.handleCreateOrganizationTeam()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}/battles/{limit}/{offset}", s.userOnly(s.orgTeamOnly(s.handleGetTeamBattles()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}/users/{limit}/{offset}", s.userOnly(s.orgTeamOnly(s.handleGetTeamUsers()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}/users", s.userOnly(s.orgTeamAdminOnly(s.handleTeamAddUser()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}/users", s.userOnly(s.orgTeamAdminOnly(s.handleTeamRemoveUser()))).Methods("DELETE")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}", s.userOnly(s.orgTeamOnly(s.handleGetOrganizationTeamByUser()))).Methods("GET")
	// org users
	s.router.HandleFunc("/api/organization/{orgId}/users/{limit}/{offset}", s.userOnly(s.orgUserOnly(s.handleGetOrganizationUsers()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/users", s.userOnly(s.orgAdminOnly(s.handleOrganizationAddUser()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}", s.userOnly(s.orgUserOnly(s.handleGetOrganizationByUser()))).Methods("GET")
	// teams(s)
	s.router.HandleFunc("/api/teams/{limit}/{offset}", s.userOnly(s.handleGetTeamsByUser())).Methods("GET")
	s.router.HandleFunc("/api/teams", s.userOnly(s.handleCreateTeam())).Methods("POST")
	s.router.HandleFunc("/api/team/{teamId}/battles/{limit}/{offset}", s.userOnly(s.teamUserOnly(s.handleGetTeamBattles()))).Methods("GET")
	s.router.HandleFunc("/api/team/{teamId}/users/{limit}/{offset}", s.userOnly(s.teamUserOnly(s.handleGetTeamUsers()))).Methods("GET")
	s.router.HandleFunc("/api/team/{teamId}/users", s.userOnly(s.teamAdminOnly(s.handleTeamAddUser()))).Methods("POST")
	s.router.HandleFunc("/api/team/{teamId}/users", s.userOnly(s.teamAdminOnly(s.handleTeamRemoveUser()))).Methods("DELETE")
	s.router.HandleFunc("/api/team/{teamId}", s.userOnly(s.teamUserOnly(s.handleGetTeamByUser()))).Methods("GET")
	// teams(s)
	// admin routes
	s.router.HandleFunc("/api/admin/stats", s.adminOnly(s.handleAppStats())).Methods("GET")
	s.router.HandleFunc("/api/admin/warriors/{limit}/{offset}", s.adminOnly(s.handleGetRegisteredUsers())).Methods("GET")
	s.router.HandleFunc("/api/admin/warrior", s.adminOnly(s.handleUserCreate())).Methods("POST")
	s.router.HandleFunc("/api/admin/promote", s.adminOnly(s.handleUserPromote())).Methods("POST")
	s.router.HandleFunc("/api/admin/demote", s.adminOnly(s.handleUserDemote())).Methods("POST")
	s.router.HandleFunc("/api/admin/clean-battles", s.adminOnly(s.handleCleanBattles())).Methods("DELETE")
	s.router.HandleFunc("/api/admin/clean-guests", s.adminOnly(s.handleCleanGuests())).Methods("DELETE")
	s.router.HandleFunc("/api/admin/organizations/{limit}/{offset}", s.adminOnly(s.handleGetOrganizations())).Methods("GET")
	// websocket for battle
	s.router.HandleFunc("/api/arena/{id}", s.serveWs())
	// handle index.html
	s.router.PathPrefix("/").HandlerFunc(s.handleIndex())
}
