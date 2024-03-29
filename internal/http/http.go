// Package http provides restful API endpoints for Thunderdome webapp
package http

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/storyboard"

	"github.com/StevenWeathers/thunderdome-planning-poker/docs/swagger"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/checkin"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/poker"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/retro"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
)

// New initializes the http handlers
// @title                       Thunderdome API
// @description                 Thunderdome Planning Poker API for both Internal and External use.
// @description                 WARNING: Currently not considered stable and is subject to change until 1.0 is released.
// @contact.name                Steven Weathers
// @contact.url                 https://github.com/StevenWeathers/thunderdome-planning-poker
// @license.name                Apache 2.0
// @license.url                 http://www.apache.org/licenses/LICENSE-2.0.html
// @version                     BETA
// @query.collection.format     multi
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        X-API-Key
func New(apiService Service, FSS fs.FS, HFS http.FileSystem) *Service {
	staticHandler := http.FileServer(HFS)
	var a = &apiService
	a.Router = mux.NewRouter()
	a.Router.Use(a.panicRecovery)

	if apiService.Config.PathPrefix != "" {
		a.Router = a.Router.PathPrefix(apiService.Config.PathPrefix).Subrouter()
	}

	a.Router.Use(otelmux.Middleware("thunderdome"))

	pokerSvc := poker.New(poker.Config{
		WriteWaitSec:  a.Config.WebsocketConfig.WriteWaitSec,
		PongWaitSec:   a.Config.WebsocketConfig.PongWaitSec,
		PingPeriodSec: a.Config.WebsocketConfig.PingPeriodSec,
		AppDomain:     a.Config.AppDomain,
	}, a.Logger, a.Cookie.ValidateSessionCookie, a.Cookie.ValidateUserCookie, a.UserDataSvc, a.AuthDataSvc, a.PokerDataSvc)
	retroSvc := retro.New(retro.Config{
		WriteWaitSec:  a.Config.WebsocketConfig.WriteWaitSec,
		PongWaitSec:   a.Config.WebsocketConfig.PongWaitSec,
		PingPeriodSec: a.Config.WebsocketConfig.PingPeriodSec,
	}, a.Logger, a.Cookie.ValidateSessionCookie, a.Cookie.ValidateUserCookie, a.UserDataSvc, a.AuthDataSvc, a.RetroDataSvc, a.Email)
	storyboardSvc := storyboard.New(storyboard.Config{
		WriteWaitSec:  a.Config.WebsocketConfig.WriteWaitSec,
		PongWaitSec:   a.Config.WebsocketConfig.PongWaitSec,
		PingPeriodSec: a.Config.WebsocketConfig.PingPeriodSec,
	}, a.Logger, a.Cookie.ValidateSessionCookie, a.Cookie.ValidateUserCookie, a.UserDataSvc, a.AuthDataSvc, a.StoryboardDataSvc)
	checkinSvc := checkin.New(checkin.Config{
		WriteWaitSec:  a.Config.WebsocketConfig.WriteWaitSec,
		PongWaitSec:   a.Config.WebsocketConfig.PongWaitSec,
		PingPeriodSec: a.Config.WebsocketConfig.PingPeriodSec,
	}, a.Logger, a.Cookie.ValidateSessionCookie, a.Cookie.ValidateUserCookie, a.UserDataSvc, a.AuthDataSvc, a.CheckinDataSvc, a.TeamDataSvc)
	swaggerJsonPath := "/" + a.Config.PathPrefix + "swagger/doc.json"
	validate = validator.New()

	swagger.SwaggerInfo.BasePath = a.Config.PathPrefix + "/api"
	// swagger docs for external API when enabled
	if a.Config.ExternalAPIEnabled {
		a.Router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL(swaggerJsonPath)))
	}

	apiRouter := a.Router.PathPrefix("/api").Subrouter()
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	orgRouter := apiRouter.PathPrefix("/organizations").Subrouter()
	teamRouter := apiRouter.PathPrefix("/teams").Subrouter()
	adminRouter := apiRouter.PathPrefix("/admin").Subrouter()

	// user authentication, profile
	if a.Config.LdapEnabled {
		apiRouter.HandleFunc("/auth/ldap", a.handleLdapLogin()).Methods("POST")
	} else if a.Config.HeaderAuthEnabled {
		apiRouter.HandleFunc("/auth", a.handleHeaderLogin()).Methods("GET")
	} else {
		apiRouter.HandleFunc("/auth", a.handleLogin()).Methods("POST")
		apiRouter.HandleFunc("/auth/forgot-password", a.handleForgotPassword()).Methods("POST")
		apiRouter.HandleFunc("/auth/reset-password", a.handleResetPassword()).Methods("PATCH")
		apiRouter.HandleFunc("/auth/update-password", a.userOnly(a.handleUpdatePassword())).Methods("PATCH")
		apiRouter.HandleFunc("/auth/verify", a.handleAccountVerification()).Methods("PATCH")
		apiRouter.HandleFunc("/auth/register", a.handleUserRegistration()).Methods("POST")
		apiRouter.HandleFunc("/auth/invite/team/{inviteId}", a.handleGetTeamInvite()).Methods("GET")
		apiRouter.HandleFunc("/auth/invite/organization/{inviteId}", a.handleGetOrganizationInvite()).Methods("GET")
	}
	apiRouter.HandleFunc("/auth/mfa", a.handleMFALogin()).Methods("POST")
	apiRouter.HandleFunc("/auth/mfa", a.userOnly(a.registeredUserOnly(a.handleMFARemove()))).Methods("DELETE")
	apiRouter.HandleFunc("/auth/mfa/setup/generate", a.userOnly(a.registeredUserOnly(a.handleMFASetupGenerate()))).Methods("POST")
	apiRouter.HandleFunc("/auth/mfa/setup/validate", a.userOnly(a.registeredUserOnly(a.handleMFASetupValidate()))).Methods("POST")
	apiRouter.HandleFunc("/auth/guest", a.handleCreateGuestUser()).Methods("POST")
	apiRouter.HandleFunc("/auth/user", a.userOnly(a.handleSessionUserProfile())).Methods("GET")
	apiRouter.HandleFunc("/auth/logout", a.handleLogout()).Methods("DELETE")
	// user(s)
	userRouter.HandleFunc("/{userId}", a.userOnly(a.entityUserOnly(a.handleUserProfile()))).Methods("GET")
	userRouter.HandleFunc("/{userId}", a.userOnly(a.entityUserOnly(a.handleUserProfileUpdate()))).Methods("PUT")
	userRouter.HandleFunc("/{userId}", a.userOnly(a.entityUserOnly(a.handleUserDelete()))).Methods("DELETE")
	userRouter.HandleFunc("/{userId}/request-verify", a.userOnly(a.entityUserOnly(a.handleVerifyRequest()))).Methods("POST")
	userRouter.HandleFunc("/{userId}/organizations", a.userOnly(a.entityUserOnly(a.handleGetOrganizationsByUser()))).Methods("GET")
	userRouter.HandleFunc("/{userId}/organizations", a.userOnly(a.entityUserOnly(a.handleCreateOrganization()))).Methods("POST")
	userRouter.HandleFunc("/{userId}/teams", a.userOnly(a.entityUserOnly(a.handleGetTeamsByUser()))).Methods("GET")
	userRouter.HandleFunc("/{userId}/teams", a.userOnly(a.entityUserOnly(a.handleCreateTeam()))).Methods("POST")
	userRouter.HandleFunc("/{userId}/jira-instances", a.userOnly(a.entityUserOnly(a.subscribedUserOnly(a.handleGetUserJiraInstances())))).Methods("GET")
	userRouter.HandleFunc("/{userId}/jira-instances", a.userOnly(a.entityUserOnly(a.subscribedUserOnly(a.handleJiraInstanceCreate())))).Methods("POST")
	userRouter.HandleFunc("/{userId}/jira-instances/{instanceId}", a.userOnly(a.entityUserOnly(a.subscribedUserOnly(a.handleJiraInstanceUpdate())))).Methods("PUT")
	userRouter.HandleFunc("/{userId}/jira-instances/{instanceId}", a.userOnly(a.entityUserOnly(a.subscribedUserOnly(a.handleJiraInstanceDelete())))).Methods("DELETE")
	userRouter.HandleFunc("/{userId}/jira-instances/{instanceId}/jql-story-search", a.userOnly(a.entityUserOnly(a.subscribedUserOnly(a.handleJiraStoryJQLSearch())))).Methods("POST")

	if a.Config.ExternalAPIEnabled {
		userRouter.HandleFunc("/{userId}/apikeys", a.userOnly(a.entityUserOnly(a.handleUserAPIKeys()))).Methods("GET")
		userRouter.HandleFunc("/{userId}/apikeys", a.userOnly(a.verifiedUserOnly(a.handleAPIKeyGenerate()))).Methods("POST")
		userRouter.HandleFunc("/{userId}/apikeys/{keyID}", a.userOnly(a.entityUserOnly(a.handleUserAPIKeyUpdate()))).Methods("PUT")
		userRouter.HandleFunc("/{userId}/apikeys/{keyID}", a.userOnly(a.entityUserOnly(a.handleUserAPIKeyDelete()))).Methods("DELETE")
	}
	// country(s)
	if a.Config.ShowActiveCountries {
		apiRouter.HandleFunc("/active-countries", a.handleGetActiveCountries()).Methods("GET")
	}
	// org
	orgRouter.HandleFunc("/{orgId}", a.userOnly(a.orgUserOnly(a.handleGetOrganizationByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}", a.userOnly(a.orgAdminOnly(a.handleDeleteOrganization()))).Methods("DELETE")
	// org departments(s)
	orgRouter.HandleFunc("/{orgId}/departments", a.userOnly(a.orgUserOnly(a.handleGetOrganizationDepartments()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments", a.userOnly(a.orgAdminOnly(a.handleCreateDepartment()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}", a.userOnly(a.orgAdminOnly(a.handleDeleteDepartment()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users", a.userOnly(a.departmentAdminOnly(a.handleDepartmentAddUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users/{userId}", a.userOnly(a.departmentAdminOnly(a.handleDepartmentUpdateUser()))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users/{userId}", a.userOnly(a.departmentAdminOnly(a.handleDepartmentRemoveUser()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentTeams()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams", a.userOnly(a.departmentAdminOnly(a.handleCreateDepartmentTeam()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}", a.userOnly(a.departmentTeamUserOnly(a.handleDepartmentTeamByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}", a.userOnly(a.departmentAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/invites", a.userOnly(a.departmentTeamUserOnly(a.handleGetTeamUserInvites()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/invites/{inviteId}", a.userOnly(a.departmentTeamAdminOnly(a.handleDeleteTeamUserInvite()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users", a.userOnly(a.departmentTeamUserOnly(a.handleGetTeamUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users", a.userOnly(a.departmentTeamAdminOnly(a.handleDepartmentTeamAddUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}", a.userOnly(a.departmentTeamAdminOnly(a.handleTeamUpdateUser()))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}", a.userOnly(a.departmentTeamAdminOnly(a.handleTeamRemoveUser()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins", a.userOnly(a.departmentTeamUserOnly(a.handleCheckinsGet()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins", a.userOnly(a.departmentTeamUserOnly(a.handleCheckinCreate(checkinSvc)))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.departmentTeamUserOnly(a.handleCheckinUpdate(checkinSvc)))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.departmentTeamUserOnly(a.handleCheckinDelete(checkinSvc)))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}/comments", a.userOnly(a.departmentTeamUserOnly(a.handleCheckinComment(checkinSvc)))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.departmentTeamUserOnly(a.handleCheckinCommentEdit(checkinSvc)))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.departmentTeamUserOnly(a.handleCheckinCommentDelete(checkinSvc)))).Methods("DELETE")
	// org teams
	orgRouter.HandleFunc("/{orgId}/teams", a.userOnly(a.orgUserOnly(a.handleGetOrganizationTeams()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams", a.userOnly(a.orgAdminOnly(a.handleCreateOrganizationTeam()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}", a.userOnly(a.orgTeamOnly(a.handleGetOrganizationTeamByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}", a.userOnly(a.orgAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/invites", a.userOnly(a.orgTeamOnly(a.handleGetTeamUserInvites()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/invites/{inviteId}", a.userOnly(a.orgTeamAdminOnly(a.handleDeleteTeamUserInvite()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users", a.userOnly(a.orgTeamOnly(a.handleGetTeamUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users", a.userOnly(a.orgTeamAdminOnly(a.handleOrganizationTeamAddUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}", a.userOnly(a.orgTeamAdminOnly(a.handleTeamUpdateUser()))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}", a.userOnly(a.orgTeamAdminOnly(a.handleTeamRemoveUser()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins", a.userOnly(a.orgTeamOnly(a.handleCheckinsGet()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins", a.userOnly(a.orgTeamOnly(a.handleCheckinCreate(checkinSvc)))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.orgTeamOnly(a.handleCheckinUpdate(checkinSvc)))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.orgTeamOnly(a.handleCheckinDelete(checkinSvc)))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/{checkinId}/comments", a.userOnly(a.orgTeamOnly(a.handleCheckinComment(checkinSvc)))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.orgTeamOnly(a.handleCheckinCommentEdit(checkinSvc)))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.orgTeamOnly(a.handleCheckinCommentDelete(checkinSvc)))).Methods("DELETE")
	// org users
	orgRouter.HandleFunc("/{orgId}/users", a.userOnly(a.orgUserOnly(a.handleGetOrganizationUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/users", a.userOnly(a.orgAdminOnly(a.handleOrganizationAddUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/users/{userId}", a.userOnly(a.orgAdminOnly(a.handleOrganizationUpdateUser()))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/users/{userId}", a.userOnly(a.orgAdminOnly(a.handleOrganizationRemoveUser()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/invites", a.userOnly(a.orgUserOnly(a.handleGetOrganizationUserInvites()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/invites/{inviteId}", a.userOnly(a.orgAdminOnly(a.handleDeleteOrganizationUserInvite()))).Methods("DELETE")
	// teams(s)
	teamRouter.HandleFunc("/{teamId}", a.userOnly(a.teamUserOnly(a.handleGetTeamByUser()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}", a.userOnly(a.teamAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleGetTeamUserInvites()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}/invites/{inviteId}", a.userOnly(a.teamAdminOnly(a.handleDeleteTeamUserInvite()))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/users", a.userOnly(a.teamUserOnly(a.handleGetTeamUsers()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}/users", a.userOnly(a.teamAdminOnly(a.handleTeamAddUser()))).Methods("POST")
	teamRouter.HandleFunc("/{teamId}/users/{userId}", a.userOnly(a.teamAdminOnly(a.handleTeamUpdateUser()))).Methods("PUT")
	teamRouter.HandleFunc("/{teamId}/users/{userId}", a.userOnly(a.teamAdminOnly(a.handleTeamRemoveUser()))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/checkin", checkinSvc.ServeWs())
	teamRouter.HandleFunc("/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinsGet()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinCreate(checkinSvc)))).Methods("POST")
	teamRouter.HandleFunc("/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinUpdate(checkinSvc)))).Methods("PUT")
	teamRouter.HandleFunc("/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinDelete(checkinSvc)))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/checkins/{checkinId}/comments", a.userOnly(a.teamUserOnly(a.handleCheckinComment(checkinSvc)))).Methods("POST")
	teamRouter.HandleFunc("/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentEdit(checkinSvc)))).Methods("PUT")
	teamRouter.HandleFunc("/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentDelete(checkinSvc)))).Methods("DELETE")
	// admin
	adminRouter.HandleFunc("/stats", a.userOnly(a.adminOnly(a.handleAppStats()))).Methods("GET")
	adminRouter.HandleFunc("/users", a.userOnly(a.adminOnly(a.handleGetRegisteredUsers()))).Methods("GET")
	adminRouter.HandleFunc("/users", a.userOnly(a.adminOnly(a.handleUserCreate()))).Methods("POST")
	adminRouter.HandleFunc("/users/{userId}/promote", a.userOnly(a.adminOnly(a.handleUserPromote()))).Methods("PATCH")
	adminRouter.HandleFunc("/users/{userId}/demote", a.userOnly(a.adminOnly(a.handleUserDemote()))).Methods("PATCH")
	adminRouter.HandleFunc("/users/{userId}/disable", a.userOnly(a.adminOnly(a.handleUserDisable()))).Methods("PATCH")
	adminRouter.HandleFunc("/users/{userId}/enable", a.userOnly(a.adminOnly(a.handleUserEnable()))).Methods("PATCH")
	adminRouter.HandleFunc("/users/{userId}/password", a.userOnly(a.adminOnly(a.handleAdminUpdateUserPassword()))).Methods("PATCH")
	adminRouter.HandleFunc("/organizations", a.userOnly(a.adminOnly(a.handleGetOrganizations()))).Methods("GET")
	adminRouter.HandleFunc("/teams", a.userOnly(a.adminOnly(a.handleGetTeams()))).Methods("GET")
	adminRouter.HandleFunc("/apikeys", a.userOnly(a.adminOnly(a.handleGetAPIKeys()))).Methods("GET")
	adminRouter.HandleFunc("/search/users/email", a.userOnly(a.adminOnly(a.handleSearchRegisteredUsersByEmail()))).Methods("GET")
	// alert
	apiRouter.HandleFunc("/alerts", a.userOnly(a.adminOnly(a.handleGetAlerts()))).Methods("GET")
	apiRouter.HandleFunc("/alerts", a.userOnly(a.adminOnly(a.handleAlertCreate()))).Methods("POST")
	apiRouter.HandleFunc("/alerts/{alertId}", a.userOnly(a.adminOnly(a.handleAlertUpdate()))).Methods("PUT")
	apiRouter.HandleFunc("/alerts/{alertId}", a.userOnly(a.adminOnly(a.handleAlertDelete()))).Methods("DELETE")
	// maintenance
	apiRouter.HandleFunc("/maintenance/clean-guests", a.userOnly(a.adminOnly(a.handleCleanGuests()))).Methods("DELETE")
	// poker games(s)
	if a.Config.FeaturePoker {
		userRouter.HandleFunc("/{userId}/battles", a.userOnly(a.entityUserOnly(a.handlePokerCreate()))).Methods("POST")
		userRouter.HandleFunc("/{userId}/battles", a.userOnly(a.entityUserOnly(a.handleGetUserGames()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/battles", a.userOnly(a.departmentTeamUserOnly(a.handleGetTeamBattles()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/battles/{battleId}", a.userOnly(a.departmentTeamAdminOnly(a.handleTeamRemoveBattle()))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/battles", a.userOnly(a.departmentTeamUserOnly(a.handlePokerCreate()))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/battles", a.userOnly(a.orgTeamOnly(a.handleGetTeamBattles()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/battles/{battleId}", a.userOnly(a.orgTeamAdminOnly(a.handleTeamRemoveBattle()))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}/battles", a.userOnly(a.orgTeamOnly(a.entityUserOnly(a.handlePokerCreate())))).Methods("POST")
		teamRouter.HandleFunc("/{teamId}/battles", a.userOnly(a.teamUserOnly(a.handleGetTeamBattles()))).Methods("GET")
		teamRouter.HandleFunc("/{teamId}/battles/{battleId}", a.userOnly(a.teamAdminOnly(a.handleTeamRemoveBattle()))).Methods("DELETE")
		teamRouter.HandleFunc("/{teamId}/users/{userId}/battles", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handlePokerCreate())))).Methods("POST")
		apiRouter.HandleFunc("/maintenance/clean-battles", a.userOnly(a.adminOnly(a.handleCleanBattles()))).Methods("DELETE")
		apiRouter.HandleFunc("/battles", a.userOnly(a.adminOnly(a.handleGetPokerGames()))).Methods("GET")
		apiRouter.HandleFunc("/battles/{battleId}", a.userOnly(a.handleGetPokerGame())).Methods("GET")
		apiRouter.HandleFunc("/battles/{battleId}", a.userOnly(a.handlePokerDelete(pokerSvc))).Methods("DELETE")
		apiRouter.HandleFunc("/battles/{battleId}/plans", a.userOnly(a.handlePokerStoryAdd(pokerSvc))).Methods("POST")
		apiRouter.HandleFunc("/battles/{battleId}/plans/{planId}", a.userOnly(a.handlePokerStoryDelete(pokerSvc))).Methods("DELETE")
		apiRouter.HandleFunc("/arena/{battleId}", pokerSvc.ServeBattleWs())
	}
	// retro(s)
	if a.Config.FeatureRetro {
		userRouter.HandleFunc("/{userId}/retros", a.userOnly(a.entityUserOnly(a.handleRetroCreate()))).Methods("POST")
		userRouter.HandleFunc("/{userId}/retros", a.userOnly(a.entityUserOnly(a.handleRetrosGetByUser()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/retros", a.userOnly(a.departmentTeamUserOnly(a.handleGetTeamRetros()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/retros/{retroId}", a.userOnly(a.departmentTeamAdminOnly(a.handleTeamRemoveRetro()))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/retro-actions", a.userOnly(a.departmentTeamUserOnly(a.handleGetTeamRetroActions()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/retros", a.userOnly(a.departmentTeamUserOnly(a.handleRetroCreate()))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/retros", a.userOnly(a.orgTeamOnly(a.handleGetTeamRetros()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/retro-actions", a.userOnly(a.orgTeamOnly(a.handleGetTeamRetroActions()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/retros/{retroId}", a.userOnly(a.orgTeamAdminOnly(a.handleTeamRemoveRetro()))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}/retros", a.userOnly(a.orgTeamOnly(a.entityUserOnly(a.handleRetroCreate())))).Methods("POST")
		teamRouter.HandleFunc("/{teamId}/retros", a.userOnly(a.teamUserOnly(a.handleGetTeamRetros()))).Methods("GET")
		teamRouter.HandleFunc("/{teamId}/retros/{retroId}", a.userOnly(a.teamAdminOnly(a.handleTeamRemoveRetro()))).Methods("DELETE")
		teamRouter.HandleFunc("/{teamId}/retro-actions", a.userOnly(a.teamUserOnly(a.handleGetTeamRetroActions()))).Methods("GET")
		teamRouter.HandleFunc("/{teamId}/users/{userId}/retros", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handleRetroCreate())))).Methods("POST")
		apiRouter.HandleFunc("/maintenance/clean-retros", a.userOnly(a.adminOnly(a.handleCleanRetros()))).Methods("DELETE")
		apiRouter.HandleFunc("/retros", a.userOnly(a.adminOnly(a.handleGetRetros()))).Methods("GET")
		apiRouter.HandleFunc("/retros/{retroId}", a.userOnly(a.handleRetroGet())).Methods("GET")
		apiRouter.HandleFunc("/retros/{retroId}", a.userOnly(a.handleRetroDelete(retroSvc))).Methods("DELETE")
		apiRouter.HandleFunc("/retros/{retroId}/actions/{actionId}", a.userOnly(a.handleRetroActionUpdate(retroSvc))).Methods("PUT")
		apiRouter.HandleFunc("/retros/{retroId}/actions/{actionId}", a.userOnly(a.handleRetroActionDelete(retroSvc))).Methods("DELETE")
		apiRouter.HandleFunc("/retros/{retroId}/actions/{actionId}/assignees", a.userOnly(a.handleRetroActionAssigneeAdd(retroSvc))).Methods("POST")
		apiRouter.HandleFunc("/retros/{retroId}/actions/{actionId}/assignees", a.userOnly(a.handleRetroActionAssigneeRemove(retroSvc))).Methods("DELETE")
		apiRouter.HandleFunc("/retros/{retroId}/actions/{actionId}/comments", a.userOnly(a.handleRetroActionCommentAdd())).Methods("POST")
		apiRouter.HandleFunc("/retros/{retroId}/actions/{actionId}/comments/{commentId}", a.userOnly(a.handleRetroActionCommentEdit())).Methods("PUT")
		apiRouter.HandleFunc("/retros/{retroId}/actions/{actionId}/comments/{commentId}", a.userOnly(a.handleRetroActionCommentDelete())).Methods("DELETE")
		apiRouter.HandleFunc("/retro/{retroId}", retroSvc.ServeWs())
	}
	// storyboard(s)
	if a.Config.FeatureRetro {
		userRouter.HandleFunc("/{userId}/storyboards", a.userOnly(a.entityUserOnly(a.handleStoryboardCreate()))).Methods("POST")
		userRouter.HandleFunc("/{userId}/storyboards", a.userOnly(a.entityUserOnly(a.handleGetUserStoryboards()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/storyboards", a.userOnly(a.departmentTeamUserOnly(a.handleGetTeamStoryboards()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/storyboards/{storyboardId}", a.userOnly(a.departmentTeamAdminOnly(a.handleTeamRemoveStoryboard()))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/storyboards", a.userOnly(a.departmentTeamUserOnly(a.handleStoryboardCreate()))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/storyboards", a.userOnly(a.orgTeamOnly(a.handleGetTeamStoryboards()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/storyboards/{storyboardId}", a.userOnly(a.orgTeamAdminOnly(a.handleTeamRemoveStoryboard()))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}/storyboards", a.userOnly(a.orgTeamOnly(a.entityUserOnly(a.handleStoryboardCreate())))).Methods("POST")
		teamRouter.HandleFunc("/{teamId}/storyboards", a.userOnly(a.teamUserOnly(a.handleGetTeamStoryboards()))).Methods("GET")
		teamRouter.HandleFunc("/{teamId}/storyboards/{storyboardId}", a.userOnly(a.teamAdminOnly(a.handleTeamRemoveStoryboard()))).Methods("DELETE")
		teamRouter.HandleFunc("/{teamId}/users/{userId}/storyboards", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handleStoryboardCreate())))).Methods("POST")
		apiRouter.HandleFunc("/maintenance/clean-storyboards", a.userOnly(a.adminOnly(a.handleCleanStoryboards()))).Methods("DELETE")
		apiRouter.HandleFunc("/storyboards", a.userOnly(a.adminOnly(a.handleGetStoryboards()))).Methods("GET")
		apiRouter.HandleFunc("/storyboards/{storyboardId}", a.userOnly(a.handleStoryboardGet())).Methods("GET")
		apiRouter.HandleFunc("/storyboards/{storyboardId}", a.userOnly(a.handleStoryboardDelete(storyboardSvc))).Methods("DELETE")
		apiRouter.HandleFunc("/storyboard/{storyboardId}", storyboardSvc.ServeWs())
	}

	// user avatar generation
	if a.Config.AvatarService == "goadorable" || a.Config.AvatarService == "govatar" {
		a.Router.PathPrefix("/avatar/{width}/{id}/{avatar}").Handler(a.handleUserAvatar()).Methods("GET")
		a.Router.PathPrefix("/avatar/{width}/{id}").Handler(a.handleUserAvatar()).Methods("GET")
	}

	if a.Config.SubscriptionsEnabled {
		apiRouter.PathPrefix("/subscriptions/{subscriptionId}").Handler(a.userOnly(a.adminOnly(a.handleSubscriptionGet()))).Methods("GET")
		apiRouter.PathPrefix("/subscriptions/{subscriptionId}").Handler(a.userOnly(a.adminOnly(a.handleSubscriptionUpdate()))).Methods("PUT")
		apiRouter.PathPrefix("/subscriptions/{subscriptionId}").Handler(a.userOnly(a.adminOnly(a.handleSubscriptionDelete()))).Methods("DELETE")
		apiRouter.PathPrefix("/subscriptions").Handler(a.userOnly(a.adminOnly(a.handleGetSubscriptions()))).Methods("GET")
		apiRouter.PathPrefix("/subscriptions").Handler(a.userOnly(a.adminOnly(a.handleSubscriptionCreate()))).Methods("POST")
		a.Router.PathPrefix("/webhooks/subscriptions").Handler(a.SubscriptionSvc.HandleWebhook()).Methods("POST")
	}

	// static assets
	a.Router.PathPrefix("/static/").Handler(http.StripPrefix(a.Config.PathPrefix, staticHandler))
	a.Router.PathPrefix("/img/").Handler(http.StripPrefix(a.Config.PathPrefix, staticHandler))

	// handle index.html
	a.Router.PathPrefix("/").HandlerFunc(a.handleIndex(FSS, a.UIConfig))

	return a
}

func (s *Service) ListenAndServe() error {
	srv := &http.Server{
		Handler:           s.Router,
		Addr:              fmt.Sprintf(":%s", s.Config.Port),
		WriteTimeout:      time.Duration(s.Config.HttpWriteTimeout) * time.Second,
		ReadTimeout:       time.Duration(s.Config.HttpReadTimeout) * time.Second,
		IdleTimeout:       time.Duration(s.Config.HttpIdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(s.Config.HttpReadHeaderTimeout) * time.Second,
	}

	s.Logger.Info("Access the WebUI via 127.0.0.1:" + s.Config.Port)

	return srv.ListenAndServe()
}

// handleIndex parses the index html file, injecting any relevant data
func (s *Service) handleIndex(FSS fs.FS, uiConfig thunderdome.UIConfig) http.HandlerFunc {
	tmpl := s.getIndexTemplate(FSS)

	ActiveAlerts = s.AlertDataSvc.GetActiveAlerts(context.Background()) // prime the active alerts cache

	return func(w http.ResponseWriter, r *http.Request) {
		uiConfig.ActiveAlerts = ActiveAlerts // get the latest alerts from memory

		if s.Config.EmbedUseOS {
			tmpl = s.getIndexTemplate(FSS)
		}

		err := tmpl.Execute(w, uiConfig)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}
