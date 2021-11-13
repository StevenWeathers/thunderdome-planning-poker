// Package api provides restful API endpoints for Thunderdome webapp
package api

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/email"
	"github.com/StevenWeathers/thunderdome-planning-poker/swaggerdocs"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Config contains configuration values used by the APIs
type Config struct {
	// the domain of the application for cookie securing
	AppDomain string
	// PathPrefix allows the application to be run on a shared domain
	PathPrefix string
	// Whether the external API is enabled
	ExternalAPIEnabled bool
	// Number of API keys a user can create
	UserAPIKeyLimit int
	// name of the cookie used exclusively by the UI
	FrontendCookieName string
	// name of the user cookie
	SecureCookieName string
	// controls whether the cookie is set to secure, only works over HTTPS
	SecureCookieFlag bool
	// Whether LDAP is enabled for authentication
	LdapEnabled bool
}

type api struct {
	config *Config
	router *mux.Router
	email  *email.Email
	cookie *securecookie.SecureCookie
	db     *database.Database
}

// standardJsonResponse structure used for all restful APIs response body
type standardJsonResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data" swaggertype:"object"`
	Meta    interface{} `json:"meta" swaggertype:"object"`
}

// pagination meta structure for query result pagination
type pagination struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type contextKey string

const (
	contextKeyUserID         contextKey = "userId"
	contextKeyUserType       contextKey = "userType"
	apiKeyHeaderName         string     = "X-API-Key"
	contextKeyOrgRole        contextKey = "orgRole"
	contextKeyDepartmentRole contextKey = "departmentRole"
	contextKeyTeamRole       contextKey = "teamRole"
	adminUserType            string     = "GENERAL"
)

// @title Thunderdome API
// @description Thunderdome Planning Poker API for both Internal and External use.
// @description WARNING: Currently not considered stable and is subject to change until 1.0 is released.
// @contact.name Steven Weathers
// @contact.url https://github.com/StevenWeathers/thunderdome-planning-poker
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @query.collection.format multi
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @version BETA
func Init(config *Config, router *mux.Router, database *database.Database, email *email.Email, cookie *securecookie.SecureCookie) *api {
	var a = &api{
		config: config,
		router: router,
		db:     database,
		email:  email,
		cookie: cookie,
	}
	swaggerJsonPath := "/" + a.config.PathPrefix + "swagger/doc.json"

	swaggerdocs.SwaggerInfo.BasePath = a.config.PathPrefix + "/api"
	// swagger docs for external API when enabled
	if a.config.ExternalAPIEnabled {
		a.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL(swaggerJsonPath)))
	}

	go h.run()

	apiRouter := a.router.PathPrefix("/api").Subrouter()
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	orgRouter := apiRouter.PathPrefix("/organizations").Subrouter()
	teamRouter := apiRouter.PathPrefix("/teams").Subrouter()
	adminRouter := apiRouter.PathPrefix("/admin").Subrouter()

	// user authentication, profile
	if a.config.LdapEnabled == true {
		apiRouter.HandleFunc("/auth/ldap", a.handleLdapLogin()).Methods("POST")
	} else {
		apiRouter.HandleFunc("/auth", a.handleLogin()).Methods("POST")
		apiRouter.HandleFunc("/auth/forgot-password", a.handleForgotPassword()).Methods("POST")
		apiRouter.HandleFunc("/auth/reset-password", a.handleResetPassword()).Methods("PATCH")
		apiRouter.HandleFunc("/auth/update-password", a.userOnly(a.handleUpdatePassword())).Methods("PATCH")
		apiRouter.HandleFunc("/auth/verify", a.handleAccountVerification()).Methods("PATCH")
		apiRouter.HandleFunc("/auth/register", a.handleUserRegistration()).Methods("POST")
	}
	apiRouter.HandleFunc("/auth/guest", a.handleCreateGuestUser()).Methods("POST")
	apiRouter.HandleFunc("/auth/logout", a.handleLogout()).Methods("DELETE")
	// user(s)
	userRouter.HandleFunc("/{userId}", a.entityUserOnly(a.handleUserProfile())).Methods("GET")
	userRouter.HandleFunc("/{userId}", a.entityUserOnly(a.handleUserProfileUpdate())).Methods("PUT")
	userRouter.HandleFunc("/{userId}", a.entityUserOnly(a.handleUserDelete())).Methods("DELETE")
	userRouter.HandleFunc("/{userId}/battles", a.entityUserOnly(a.handleBattleCreate())).Methods("POST")
	userRouter.HandleFunc("/{userId}/battles", a.entityUserOnly(a.handleGetUserBattles())).Methods("GET")
	userRouter.HandleFunc("/{userId}/organizations", a.entityUserOnly(a.handleGetOrganizationsByUser())).Methods("GET")
	userRouter.HandleFunc("/{userId}/organizations", a.entityUserOnly(a.handleCreateOrganization())).Methods("POST")
	userRouter.HandleFunc("/{userId}/teams", a.entityUserOnly(a.handleGetTeamsByUser())).Methods("GET")
	userRouter.HandleFunc("/{userId}/teams", a.entityUserOnly(a.handleCreateTeam())).Methods("POST")
	if a.config.ExternalAPIEnabled {
		userRouter.HandleFunc("/{userId}/apikeys", a.entityUserOnly(a.handleUserAPIKeys())).Methods("GET")
		userRouter.HandleFunc("/{userId}/apikeys", a.verifiedUserOnly(a.handleAPIKeyGenerate())).Methods("POST")
		userRouter.HandleFunc("/{userId}/apikeys/{keyID}", a.entityUserOnly(a.handleUserAPIKeyUpdate())).Methods("PUT")
		userRouter.HandleFunc("/{userId}/apikeys/{keyID}", a.entityUserOnly(a.handleUserAPIKeyDelete())).Methods("DELETE")
	}
	// country(s)
	if viper.GetBool("config.show_active_countries") {
		apiRouter.HandleFunc("/active-countries", a.handleGetActiveCountries()).Methods("GET")
	}
	// org departments(s)
	orgRouter.HandleFunc("/{orgId}/departments", a.userOnly(a.orgUserOnly(a.handleGetOrganizationDepartments()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments", a.userOnly(a.orgAdminOnly(a.handleCreateDepartment()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users", a.userOnly(a.departmentAdminOnly(a.handleDepartmentAddUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users/{userId}", a.userOnly(a.departmentAdminOnly(a.handleDepartmentRemoveUser()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentTeams()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams", a.userOnly(a.departmentAdminOnly(a.handleCreateDepartmentTeam()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}", a.userOnly(a.departmentTeamUserOnly(a.handleDepartmentTeamByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}", a.userOnly(a.departmentAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/battles", a.userOnly(a.departmentTeamUserOnly(a.handleGetTeamBattles()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/battles/{battleId}", a.userOnly(a.departmentTeamAdminOnly(a.handleTeamRemoveBattle()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users", a.userOnly(a.departmentTeamUserOnly(a.handleGetTeamUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users", a.userOnly(a.departmentTeamAdminOnly(a.handleDepartmentTeamAddUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}", a.userOnly(a.departmentTeamAdminOnly(a.handleTeamRemoveUser()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/battles", a.userOnly(a.departmentTeamUserOnly(a.handleBattleCreate()))).Methods("POST")
	// org teams
	orgRouter.HandleFunc("/{orgId}/teams", a.userOnly(a.orgUserOnly(a.handleGetOrganizationTeams()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams", a.userOnly(a.orgAdminOnly(a.handleCreateOrganizationTeam()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}", a.userOnly(a.orgTeamOnly(a.handleGetOrganizationTeamByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}", a.userOnly(a.orgAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/battles", a.userOnly(a.orgTeamOnly(a.handleGetTeamBattles()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/battles/{battleId}", a.userOnly(a.orgTeamAdminOnly(a.handleTeamRemoveBattle()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users", a.userOnly(a.orgTeamOnly(a.handleGetTeamUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users", a.userOnly(a.orgTeamAdminOnly(a.handleOrganizationTeamAddUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}", a.userOnly(a.orgTeamAdminOnly(a.handleTeamRemoveUser()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}/battles", a.userOnly(a.orgTeamOnly(a.handleBattleCreate()))).Methods("POST")
	// org users
	orgRouter.HandleFunc("/{orgId}", a.userOnly(a.orgUserOnly(a.handleGetOrganizationByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/users", a.userOnly(a.orgUserOnly(a.handleGetOrganizationUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/users", a.userOnly(a.orgAdminOnly(a.handleOrganizationAddUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/users/{userId}", a.userOnly(a.orgAdminOnly(a.handleOrganizationRemoveUser()))).Methods("DELETE")
	// teams(s)
	teamRouter.HandleFunc("/{teamId}", a.userOnly(a.teamUserOnly(a.handleGetTeamByUser()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}", a.userOnly(a.teamAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/battles", a.userOnly(a.teamUserOnly(a.handleGetTeamBattles()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}/battles/{battleId}", a.userOnly(a.teamAdminOnly(a.handleTeamRemoveBattle()))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/users", a.userOnly(a.teamUserOnly(a.handleGetTeamUsers()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}/users", a.userOnly(a.teamAdminOnly(a.handleTeamAddUser()))).Methods("POST")
	teamRouter.HandleFunc("/{teamId}/users/{userId}", a.userOnly(a.teamAdminOnly(a.handleTeamRemoveUser()))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/users/{userId}/battles", a.userOnly(a.teamUserOnly(a.handleBattleCreate()))).Methods("POST")
	// admin
	adminRouter.HandleFunc("/stats", a.adminOnly(a.handleAppStats())).Methods("GET")
	adminRouter.HandleFunc("/users", a.adminOnly(a.handleGetRegisteredUsers())).Methods("GET")
	adminRouter.HandleFunc("/users", a.adminOnly(a.handleUserCreate())).Methods("POST")
	adminRouter.HandleFunc("/users/{userId}/promote", a.adminOnly(a.handleUserPromote())).Methods("PATCH")
	adminRouter.HandleFunc("/users/{userId}/demote", a.adminOnly(a.handleUserDemote())).Methods("PATCH")
	adminRouter.HandleFunc("/organizations", a.adminOnly(a.handleGetOrganizations())).Methods("GET")
	adminRouter.HandleFunc("/teams", a.adminOnly(a.handleGetTeams())).Methods("GET")
	adminRouter.HandleFunc("/apikeys", a.adminOnly(a.handleGetAPIKeys())).Methods("GET")
	adminRouter.HandleFunc("/search/users/email", a.adminOnly(a.handleSearchRegisteredUsersByEmail())).Methods("GET")
	// alert
	apiRouter.HandleFunc("/alerts", a.adminOnly(a.handleGetAlerts())).Methods("GET")
	apiRouter.HandleFunc("/alerts", a.adminOnly(a.handleAlertCreate())).Methods("POST")
	apiRouter.HandleFunc("/alerts/{alertId}", a.adminOnly(a.handleAlertUpdate())).Methods("PUT")
	apiRouter.HandleFunc("/alerts/{alertId}", a.adminOnly(a.handleAlertDelete())).Methods("DELETE")
	// maintenance
	apiRouter.HandleFunc("/maintenance/clean-battles", a.adminOnly(a.handleCleanBattles())).Methods("DELETE")
	apiRouter.HandleFunc("/maintenance/clean-guests", a.adminOnly(a.handleCleanGuests())).Methods("DELETE")
	apiRouter.HandleFunc("/maintenance/lowercase-emails", a.adminOnly(a.handleLowercaseUserEmails())).Methods("PATCH")
	// battle
	apiRouter.HandleFunc("/battles", a.adminOnly(a.handleGetBattles())).Methods("GET")
	apiRouter.HandleFunc("/battles/{battleId}", a.userOnly(a.handleGetBattle())).Methods("GET")
	// websocket for battle
	apiRouter.HandleFunc("/arena/{battleId}", a.serveWs())

	return a
}
