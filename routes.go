package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

//go:embed dist
var f embed.FS

func getFileSystem(useOS bool) (http.FileSystem, fs.FS) {
	if useOS {
		log.Print("using live mode")
		return http.FS(os.DirFS("dist")), fs.FS(os.DirFS("dist"))
	}

	fsys, err := fs.Sub(f, "dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys), fs.FS(fsys)
}

func (s *server) routes() {
	HFS, FSS := getFileSystem(embedUseOS)
	staticHandler := http.FileServer(HFS)

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
	// country(s)
	if viper.GetBool("config.show_active_countries") {
		s.router.HandleFunc("/api/active-countries", s.handleGetActiveCountries()).Methods("GET")
	}
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
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/user", s.userOnly(s.departmentAdminOnly(s.handleDepartmentRemoveUser()))).Methods("DELETE")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}/battles/{limit}/{offset}", s.userOnly(s.departmentTeamUserOnly(s.handleGetTeamBattles()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}/battle", s.userOnly(s.departmentTeamUserOnly(s.handleBattleCreate()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}/battle", s.userOnly(s.departmentTeamAdminOnly(s.handleTeamRemoveBattle()))).Methods("DELETE")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}/users/{limit}/{offset}", s.userOnly(s.departmentTeamUserOnly(s.handleGetTeamUsers()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}/users", s.userOnly(s.departmentTeamAdminOnly(s.handleDepartmentTeamAddUser()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}/user", s.userOnly(s.departmentTeamAdminOnly(s.handleTeamRemoveUser()))).Methods("DELETE")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team/{teamId}", s.userOnly(s.departmentTeamUserOnly(s.handleDepartmentTeamByUser()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}/team", s.userOnly(s.departmentAdminOnly(s.handleDeleteTeam()))).Methods("DELETE")
	s.router.HandleFunc("/api/organization/{orgId}/department/{departmentId}", s.userOnly(s.departmentUserOnly(s.handleGetDepartmentByUser()))).Methods("GET")
	// org teams
	s.router.HandleFunc("/api/organization/{orgId}/teams/{limit}/{offset}", s.userOnly(s.orgUserOnly(s.handleGetOrganizationTeams()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/teams", s.userOnly(s.orgAdminOnly(s.handleCreateOrganizationTeam()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}/battles/{limit}/{offset}", s.userOnly(s.orgTeamOnly(s.handleGetTeamBattles()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}/battle", s.userOnly(s.orgTeamOnly(s.handleBattleCreate()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}/battle", s.userOnly(s.orgTeamAdminOnly(s.handleTeamRemoveBattle()))).Methods("DELETE")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}/users/{limit}/{offset}", s.userOnly(s.orgTeamOnly(s.handleGetTeamUsers()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}/users", s.userOnly(s.orgTeamAdminOnly(s.handleOrganizationTeamAddUser()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}/user", s.userOnly(s.orgTeamAdminOnly(s.handleTeamRemoveUser()))).Methods("DELETE")
	s.router.HandleFunc("/api/organization/{orgId}/team/{teamId}", s.userOnly(s.orgTeamOnly(s.handleGetOrganizationTeamByUser()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/team", s.userOnly(s.orgAdminOnly(s.handleDeleteTeam()))).Methods("DELETE")
	// org users
	s.router.HandleFunc("/api/organization/{orgId}/users/{limit}/{offset}", s.userOnly(s.orgUserOnly(s.handleGetOrganizationUsers()))).Methods("GET")
	s.router.HandleFunc("/api/organization/{orgId}/users", s.userOnly(s.orgAdminOnly(s.handleOrganizationAddUser()))).Methods("POST")
	s.router.HandleFunc("/api/organization/{orgId}/user", s.userOnly(s.orgAdminOnly(s.handleOrganizationRemoveUser()))).Methods("DELETE")
	s.router.HandleFunc("/api/organization/{orgId}", s.userOnly(s.orgUserOnly(s.handleGetOrganizationByUser()))).Methods("GET")
	// teams(s)
	s.router.HandleFunc("/api/teams/{limit}/{offset}", s.userOnly(s.handleGetTeamsByUser())).Methods("GET")
	s.router.HandleFunc("/api/teams", s.userOnly(s.handleCreateTeam())).Methods("POST")
	s.router.HandleFunc("/api/team/{teamId}/battles/{limit}/{offset}", s.userOnly(s.teamUserOnly(s.handleGetTeamBattles()))).Methods("GET")
	s.router.HandleFunc("/api/team/{teamId}/battle", s.userOnly(s.teamUserOnly(s.handleBattleCreate()))).Methods("POST")
	s.router.HandleFunc("/api/team/{teamId}/battle", s.userOnly(s.teamAdminOnly(s.handleTeamRemoveBattle()))).Methods("DELETE")
	s.router.HandleFunc("/api/team/{teamId}/users/{limit}/{offset}", s.userOnly(s.teamUserOnly(s.handleGetTeamUsers()))).Methods("GET")
	s.router.HandleFunc("/api/team/{teamId}/users", s.userOnly(s.teamAdminOnly(s.handleTeamAddUser()))).Methods("POST")
	s.router.HandleFunc("/api/team/{teamId}/user", s.userOnly(s.teamAdminOnly(s.handleTeamRemoveUser()))).Methods("DELETE")
	s.router.HandleFunc("/api/team/{teamId}", s.userOnly(s.teamUserOnly(s.handleGetTeamByUser()))).Methods("GET")
	s.router.HandleFunc("/api/team", s.userOnly(s.teamAdminOnly(s.handleDeleteTeam()))).Methods("DELETE")
	// teams(s)
	// admin routes
	s.router.HandleFunc("/api/admin/stats", s.adminOnly(s.handleAppStats())).Methods("GET")
	s.router.HandleFunc("/api/admin/warriors/{limit}/{offset}", s.adminOnly(s.handleGetRegisteredUsers())).Methods("GET")
	s.router.HandleFunc("/api/admin/warrior", s.adminOnly(s.handleUserCreate())).Methods("POST")
	s.router.HandleFunc("/api/admin/user/{id}", s.adminOnly(s.handleAdminUserDelete())).Methods("DELETE")
	s.router.HandleFunc("/api/admin/promote", s.adminOnly(s.handleUserPromote())).Methods("POST")
	s.router.HandleFunc("/api/admin/demote", s.adminOnly(s.handleUserDemote())).Methods("POST")
	s.router.HandleFunc("/api/admin/maintenance/clean-battles", s.adminOnly(s.handleCleanBattles())).Methods("DELETE")
	s.router.HandleFunc("/api/admin/maintenance/clean-guests", s.adminOnly(s.handleCleanGuests())).Methods("DELETE")
	s.router.HandleFunc("/api/admin/maintenance/lowercase-emails", s.adminOnly(s.handleLowercaseUserEmails())).Methods("PUT")
	s.router.HandleFunc("/api/admin/organizations/{limit}/{offset}", s.adminOnly(s.handleGetOrganizations())).Methods("GET")
	s.router.HandleFunc("/api/admin/teams/{limit}/{offset}", s.adminOnly(s.handleGetTeams())).Methods("GET")
	s.router.HandleFunc("/api/admin/apikeys/{limit}/{offset}", s.adminOnly(s.handleGetAPIKeys())).Methods("GET")
	s.router.HandleFunc("/api/admin/alerts/{limit}/{offset}", s.adminOnly(s.handleGetAlerts())).Methods("GET")
	s.router.HandleFunc("/api/admin/alert/{id}", s.adminOnly(s.handleAlertUpdate())).Methods("PUT")
	s.router.HandleFunc("/api/admin/alert", s.adminOnly(s.handleAlertCreate())).Methods("POST")
	s.router.HandleFunc("/api/admin/alert", s.adminOnly(s.handleAlertDelete())).Methods("DELETE")
	// websocket for battle
	s.router.HandleFunc("/api/arena/{id}", s.serveWs())
	// handle index.html
	s.router.PathPrefix("/").HandlerFunc(s.handleIndex(FSS))
}
