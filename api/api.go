package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/email"
	"github.com/StevenWeathers/thunderdome-planning-poker/swaggerdocs"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type ApiConfig struct {
	// the domain of the application for cookie securing
	AppDomain string
	// PathPrefix allows the application to be run on a shared domain
	PathPrefix string
	// Whether or not the external API is enabled
	ExternalAPIEnabled bool
	// name of the cookie used exclusively by the UI
	FrontendCookieName string
	// name of the user cookie
	SecureCookieName string
	// controls whether or not the cookie is set to secure, only works over HTTPS
	SecureCookieFlag bool
}

type api struct {
	config *ApiConfig
	router *mux.Router
	email  *email.Email
	cookie *securecookie.SecureCookie
	db     *database.Database
}

type contextKey string

var (
	contextKeyUserID         contextKey = "userId"
	apiKeyHeaderName         string     = "X-API-Key"
	contextKeyOrgRole        contextKey = "orgRole"
	contextKeyDepartmentRole contextKey = "departmentRole"
	contextKeyTeamRole       contextKey = "teamRole"
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
func Init(config *ApiConfig, router *mux.Router, database *database.Database, email *email.Email, cookie *securecookie.SecureCookie) *api {
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

	// user authentication, profile
	if viper.GetString("auth.method") == "ldap" {
		apiRouter.HandleFunc("/auth", a.handleLdapLogin()).Methods("POST")
	} else {
		apiRouter.HandleFunc("/auth", a.handleLogin()).Methods("POST")
		apiRouter.HandleFunc("/auth/forgot-password", a.handleForgotPassword()).Methods("POST")
		apiRouter.HandleFunc("/auth/reset-password", a.handleResetPassword()).Methods("POST")
		apiRouter.HandleFunc("/auth/update-password", a.userOnly(a.handleUpdatePassword())).Methods("POST")
		apiRouter.HandleFunc("/auth/verify", a.handleAccountVerification()).Methods("POST")
		apiRouter.HandleFunc("/enlist", a.handleUserEnlist()).Methods("POST")
	}
	apiRouter.HandleFunc("/warrior", a.handleUserRecruit()).Methods("POST")
	apiRouter.HandleFunc("/auth/logout", a.handleLogout()).Methods("POST")
	if a.config.ExternalAPIEnabled {
		apiRouter.HandleFunc("/warrior/{id}/apikey/{keyID}", a.userOnly(a.handleUserAPIKeyUpdate())).Methods("PUT")
		apiRouter.HandleFunc("/warrior/{id}/apikey/{keyID}", a.userOnly(a.handleUserAPIKeyDelete())).Methods("DELETE")
		apiRouter.HandleFunc("/warrior/{id}/apikey", a.userOnly(a.handleAPIKeyGenerate())).Methods("POST")
		apiRouter.HandleFunc("/warrior/{id}/apikeys", a.userOnly(a.handleUserAPIKeys())).Methods("GET")
	}
	apiRouter.HandleFunc("/warrior/{id}", a.userOnly(a.handleUserProfile())).Methods("GET")
	apiRouter.HandleFunc("/warrior/{id}", a.userOnly(a.handleUserProfileUpdate())).Methods("POST")
	apiRouter.HandleFunc("/warrior/{id}", a.userOnly(a.handleUserDelete())).Methods("DELETE")
	// battle(s)
	apiRouter.HandleFunc("/battle", a.userOnly(a.handleBattleCreate())).Methods("POST")
	apiRouter.HandleFunc("/battles", a.userOnly(a.handleBattlesGet())).Methods("GET")
	// country(s)
	if viper.GetBool("config.show_active_countries") {
		apiRouter.HandleFunc("/active-countries", a.handleGetActiveCountries()).Methods("GET")
	}
	// organization(s)
	apiRouter.HandleFunc("/organizations/{limit}/{offset}", a.userOnly(a.handleGetOrganizationsByUser())).Methods("GET")
	apiRouter.HandleFunc("/organizations", a.userOnly(a.handleCreateOrganization())).Methods("POST")
	apiRouter.HandleFunc("/organization/{orgId}/departments/{limit}/{offset}", a.userOnly(a.orgUserOnly(a.handleGetOrganizationDepartments()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/departments", a.userOnly(a.orgAdminOnly(a.handleCreateDepartment()))).Methods("POST")
	// org departments(s)
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/teams/{limit}/{offset}", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentTeams()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/teams", a.userOnly(a.departmentAdminOnly(a.handleCreateDepartmentTeam()))).Methods("POST")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/users/{limit}/{offset}", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentUsers()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/users", a.userOnly(a.departmentAdminOnly(a.handleDepartmentAddUser()))).Methods("POST")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/user", a.userOnly(a.departmentAdminOnly(a.handleDepartmentRemoveUser()))).Methods("DELETE")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/team/{teamId}/battles/{limit}/{offset}", a.userOnly(a.departmentTeamUserOnly(a.handleGetTeamBattles()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/team/{teamId}/battle", a.userOnly(a.departmentTeamUserOnly(a.handleBattleCreate()))).Methods("POST")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/team/{teamId}/battle", a.userOnly(a.departmentTeamAdminOnly(a.handleTeamRemoveBattle()))).Methods("DELETE")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/team/{teamId}/users/{limit}/{offset}", a.userOnly(a.departmentTeamUserOnly(a.handleGetTeamUsers()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/team/{teamId}/users", a.userOnly(a.departmentTeamAdminOnly(a.handleDepartmentTeamAddUser()))).Methods("POST")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/team/{teamId}/user", a.userOnly(a.departmentTeamAdminOnly(a.handleTeamRemoveUser()))).Methods("DELETE")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/team/{teamId}", a.userOnly(a.departmentTeamUserOnly(a.handleDepartmentTeamByUser()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}/team", a.userOnly(a.departmentAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	apiRouter.HandleFunc("/organization/{orgId}/department/{departmentId}", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentByUser()))).Methods("GET")
	// org teams
	apiRouter.HandleFunc("/organization/{orgId}/teams/{limit}/{offset}", a.userOnly(a.orgUserOnly(a.handleGetOrganizationTeams()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/teams", a.userOnly(a.orgAdminOnly(a.handleCreateOrganizationTeam()))).Methods("POST")
	apiRouter.HandleFunc("/organization/{orgId}/team/{teamId}/battles/{limit}/{offset}", a.userOnly(a.orgTeamOnly(a.handleGetTeamBattles()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/team/{teamId}/battle", a.userOnly(a.orgTeamOnly(a.handleBattleCreate()))).Methods("POST")
	apiRouter.HandleFunc("/organization/{orgId}/team/{teamId}/battle", a.userOnly(a.orgTeamAdminOnly(a.handleTeamRemoveBattle()))).Methods("DELETE")
	apiRouter.HandleFunc("/organization/{orgId}/team/{teamId}/users/{limit}/{offset}", a.userOnly(a.orgTeamOnly(a.handleGetTeamUsers()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/team/{teamId}/users", a.userOnly(a.orgTeamAdminOnly(a.handleOrganizationTeamAddUser()))).Methods("POST")
	apiRouter.HandleFunc("/organization/{orgId}/team/{teamId}/user", a.userOnly(a.orgTeamAdminOnly(a.handleTeamRemoveUser()))).Methods("DELETE")
	apiRouter.HandleFunc("/organization/{orgId}/team/{teamId}", a.userOnly(a.orgTeamOnly(a.handleGetOrganizationTeamByUser()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/team", a.userOnly(a.orgAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	// org users
	apiRouter.HandleFunc("/organization/{orgId}/users/{limit}/{offset}", a.userOnly(a.orgUserOnly(a.handleGetOrganizationUsers()))).Methods("GET")
	apiRouter.HandleFunc("/organization/{orgId}/users", a.userOnly(a.orgAdminOnly(a.handleOrganizationAddUser()))).Methods("POST")
	apiRouter.HandleFunc("/organization/{orgId}/user", a.userOnly(a.orgAdminOnly(a.handleOrganizationRemoveUser()))).Methods("DELETE")
	apiRouter.HandleFunc("/organization/{orgId}", a.userOnly(a.orgUserOnly(a.handleGetOrganizationByUser()))).Methods("GET")
	// teams(s)
	apiRouter.HandleFunc("/teams/{limit}/{offset}", a.userOnly(a.handleGetTeamsByUser())).Methods("GET")
	apiRouter.HandleFunc("/teams", a.userOnly(a.handleCreateTeam())).Methods("POST")
	apiRouter.HandleFunc("/team/{teamId}/battles/{limit}/{offset}", a.userOnly(a.teamUserOnly(a.handleGetTeamBattles()))).Methods("GET")
	apiRouter.HandleFunc("/team/{teamId}/battle", a.userOnly(a.teamUserOnly(a.handleBattleCreate()))).Methods("POST")
	apiRouter.HandleFunc("/team/{teamId}/battle", a.userOnly(a.teamAdminOnly(a.handleTeamRemoveBattle()))).Methods("DELETE")
	apiRouter.HandleFunc("/team/{teamId}/users/{limit}/{offset}", a.userOnly(a.teamUserOnly(a.handleGetTeamUsers()))).Methods("GET")
	apiRouter.HandleFunc("/team/{teamId}/users", a.userOnly(a.teamAdminOnly(a.handleTeamAddUser()))).Methods("POST")
	apiRouter.HandleFunc("/team/{teamId}/user", a.userOnly(a.teamAdminOnly(a.handleTeamRemoveUser()))).Methods("DELETE")
	apiRouter.HandleFunc("/team/{teamId}", a.userOnly(a.teamUserOnly(a.handleGetTeamByUser()))).Methods("GET")
	apiRouter.HandleFunc("/team", a.userOnly(a.teamAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	// admin
	apiRouter.HandleFunc("/admin/stats", a.adminOnly(a.handleAppStats())).Methods("GET")
	apiRouter.HandleFunc("/admin/warriors/{limit}/{offset}", a.adminOnly(a.handleGetRegisteredUsers())).Methods("GET")
	apiRouter.HandleFunc("/admin/warrior", a.adminOnly(a.handleUserCreate())).Methods("POST")
	apiRouter.HandleFunc("/admin/user/{id}", a.adminOnly(a.handleAdminUserDelete())).Methods("DELETE")
	apiRouter.HandleFunc("/admin/promote", a.adminOnly(a.handleUserPromote())).Methods("POST")
	apiRouter.HandleFunc("/admin/demote", a.adminOnly(a.handleUserDemote())).Methods("POST")
	apiRouter.HandleFunc("/admin/organizations/{limit}/{offset}", a.adminOnly(a.handleGetOrganizations())).Methods("GET")
	apiRouter.HandleFunc("/admin/teams/{limit}/{offset}", a.adminOnly(a.handleGetTeams())).Methods("GET")
	apiRouter.HandleFunc("/admin/apikeys/{limit}/{offset}", a.adminOnly(a.handleGetAPIKeys())).Methods("GET")
	// alert
	apiRouter.HandleFunc("/alerts", a.adminOnly(a.handleGetAlerts())).Methods("GET")
	apiRouter.HandleFunc("/alerts", a.adminOnly(a.handleAlertCreate())).Methods("POST")
	apiRouter.HandleFunc("/alerts/{id}", a.adminOnly(a.handleAlertUpdate())).Methods("PUT")
	apiRouter.HandleFunc("/alerts/{id}", a.adminOnly(a.handleAlertDelete())).Methods("DELETE")
	// maintenance
	apiRouter.HandleFunc("/maintenance/clean-battles", a.adminOnly(a.handleCleanBattles())).Methods("DELETE")
	apiRouter.HandleFunc("/maintenance/clean-guests", a.adminOnly(a.handleCleanGuests())).Methods("DELETE")
	apiRouter.HandleFunc("/maintenance/lowercase-emails", a.adminOnly(a.handleLowercaseUserEmails())).Methods("PUT")
	// websocket for battle
	apiRouter.HandleFunc("/arena/{id}", a.serveWs())

	return a
}

// getJSONRequestBody gets a JSON request body broken into a key/value map
func (a *api) getJSONRequestBody(r *http.Request, w http.ResponseWriter) map[string]interface{} {
	body, _ := ioutil.ReadAll(r.Body) // check for errors
	keyVal := make(map[string]interface{})
	jsonErr := json.Unmarshal(body, &keyVal) // check for errors

	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	return keyVal
}

// respondWithJSON takes a payload and writes the response
func (a *api) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
