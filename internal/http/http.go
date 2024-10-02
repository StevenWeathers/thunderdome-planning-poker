// Package http provides restful API endpoints for Thunderdome webapp
package http

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"time"

	"github.com/unrolled/secure"
	"github.com/unrolled/secure/cspbuilder"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/oauth"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/storyboard"

	"github.com/StevenWeathers/thunderdome-planning-poker/docs/swagger"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/checkin"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/poker"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/retro"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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
	authProviderConfigs := make([]thunderdome.AuthProviderConfig, 0)
	// Content Security Policy
	cspBuilder := cspbuilder.Builder{
		Directives: map[string][]string{
			cspbuilder.DefaultSrc: {"self", fmt.Sprintf("*.%s", a.Config.AppDomain)},
			// @TODO - remove inline styles in svelte components to improve security by using nonce
			cspbuilder.StyleSrc:  {"'self'", "'unsafe-inline'", "https://fonts.googleapis.com"},
			cspbuilder.ScriptSrc: {"$NONCE"},
			cspbuilder.FontSrc:   {"'self'", "https://fonts.gstatic.com"},
			cspbuilder.ImgSrc:    {"data:", "*"},
			cspbuilder.ConnectSrc: {"'self'",
				getWebsocketConnectSrc(a.Config.SecureProtocol, a.Config.WebsocketSubdomain, a.Config.AppDomain),
				"https://*.google-analytics.com",
				"https://*.analytics.google.com",
				"https://*.googletagmanager.com",
				"https://*.google.com",
			},
			cspbuilder.ManifestSrc: {"'self'"},
		},
	}

	secureMiddleware := secure.New(secure.Options{
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		ForceSTSHeader:        true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: cspBuilder.MustBuild(),
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	})

	a.Router = mux.NewRouter()
	a.Router.Use(a.panicRecovery)

	if apiService.Config.PathPrefix != "" {
		a.Router = a.Router.PathPrefix(apiService.Config.PathPrefix).Subrouter()
	}

	swaggerJsonPath := "/" + a.Config.PathPrefix + "swagger/doc.json"
	swagger.SwaggerInfo.BasePath = a.Config.PathPrefix + "/api"
	// swagger docs for external API when enabled
	// has to come before csp policy as there is currently no way to configure csp nonce for swagger ui
	if a.Config.ExternalAPIEnabled {
		a.Router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL(swaggerJsonPath)))
	}

	router := a.Router.PathPrefix("/").Subrouter()
	router.Use(secureMiddleware.Handler)
	router.Use(otelmux.Middleware("thunderdome"))

	pokerSvc := poker.New(poker.Config{
		WriteWaitSec:       a.Config.WebsocketConfig.WriteWaitSec,
		PongWaitSec:        a.Config.WebsocketConfig.PongWaitSec,
		PingPeriodSec:      a.Config.WebsocketConfig.PingPeriodSec,
		AppDomain:          a.Config.AppDomain,
		WebsocketSubdomain: a.Config.WebsocketConfig.WebsocketSubdomain,
	}, a.Logger, a.Cookie.ValidateSessionCookie, a.Cookie.ValidateUserCookie, a.UserDataSvc, a.AuthDataSvc, a.PokerDataSvc)
	retroSvc := retro.New(retro.Config{
		WriteWaitSec:       a.Config.WebsocketConfig.WriteWaitSec,
		PongWaitSec:        a.Config.WebsocketConfig.PongWaitSec,
		PingPeriodSec:      a.Config.WebsocketConfig.PingPeriodSec,
		AppDomain:          a.Config.AppDomain,
		WebsocketSubdomain: a.Config.WebsocketConfig.WebsocketSubdomain,
	}, a.Logger, a.Cookie.ValidateSessionCookie, a.Cookie.ValidateUserCookie, a.UserDataSvc, a.AuthDataSvc,
		a.RetroDataSvc, a.RetroTemplateDataSvc, a.Email)
	storyboardSvc := storyboard.New(storyboard.Config{
		WriteWaitSec:       a.Config.WebsocketConfig.WriteWaitSec,
		PongWaitSec:        a.Config.WebsocketConfig.PongWaitSec,
		PingPeriodSec:      a.Config.WebsocketConfig.PingPeriodSec,
		AppDomain:          a.Config.AppDomain,
		WebsocketSubdomain: a.Config.WebsocketConfig.WebsocketSubdomain,
	}, a.Logger, a.Cookie.ValidateSessionCookie, a.Cookie.ValidateUserCookie, a.UserDataSvc, a.AuthDataSvc, a.StoryboardDataSvc)
	checkinSvc := checkin.New(checkin.Config{
		WriteWaitSec:       a.Config.WebsocketConfig.WriteWaitSec,
		PongWaitSec:        a.Config.WebsocketConfig.PongWaitSec,
		PingPeriodSec:      a.Config.WebsocketConfig.PingPeriodSec,
		AppDomain:          a.Config.AppDomain,
		WebsocketSubdomain: a.Config.WebsocketConfig.WebsocketSubdomain,
	}, a.Logger, a.Cookie.ValidateSessionCookie, a.Cookie.ValidateUserCookie, a.UserDataSvc, a.AuthDataSvc, a.CheckinDataSvc, a.TeamDataSvc)

	validate = validator.New()

	apiRouter := router.PathPrefix("/api").Subrouter()
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	orgRouter := apiRouter.PathPrefix("/organizations").Subrouter()
	teamRouter := apiRouter.PathPrefix("/teams").Subrouter()
	adminRouter := apiRouter.PathPrefix("/admin").Subrouter()

	apiRouter.HandleFunc("/", a.handleApiIndex()).Methods("GET")

	// user authentication, profile
	if a.Config.LdapEnabled {
		apiRouter.HandleFunc("/auth/ldap", a.handleLdapLogin()).Methods("POST")
	} else if a.Config.HeaderAuthEnabled {
		apiRouter.HandleFunc("/auth", a.handleHeaderLogin()).Methods("GET")
	} else {
		if a.Config.GoogleAuth.Enabled {
			authProviderConfigs = append(authProviderConfigs, thunderdome.AuthProviderConfig{
				ProviderName: a.Config.GoogleAuth.ProviderName,
				ProviderURL:  a.Config.GoogleAuth.ProviderURL,
				ClientID:     a.Config.GoogleAuth.ClientID,
				ClientSecret: a.Config.GoogleAuth.ClientSecret,
			})
		}
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
	userRouter.HandleFunc("/{userId}/credential", a.userOnly(a.entityUserOnly(a.handleUserCredential()))).Methods("GET")
	userRouter.HandleFunc("/{userId}/request-verify", a.userOnly(a.entityUserOnly(a.handleVerifyRequest()))).Methods("POST")
	userRouter.HandleFunc("/{userId}/invite/team/{inviteId}", a.userOnly(a.registeredUserOnly(a.handleUserTeamInvite()))).Methods("POST")
	userRouter.HandleFunc("/{userId}/invite/organization/{inviteId}", a.userOnly(a.registeredUserOnly(a.handleUserOrganizationInvite()))).Methods("POST")
	userRouter.HandleFunc("/{userId}/invite/department/{inviteId}", a.userOnly(a.registeredUserOnly(a.handleUserDepartmentInvite()))).Methods("POST")
	userRouter.HandleFunc("/{userId}/organizations", a.userOnly(a.entityUserOnly(a.handleGetOrganizationsByUser()))).Methods("GET")
	userRouter.HandleFunc("/{userId}/organizations", a.userOnly(a.entityUserOnly(a.handleCreateOrganization()))).Methods("POST")
	userRouter.HandleFunc("/{userId}/teams", a.userOnly(a.entityUserOnly(a.handleGetTeamsByUser()))).Methods("GET")
	userRouter.HandleFunc("/{userId}/teams", a.userOnly(a.entityUserOnly(a.handleCreateTeam()))).Methods("POST")
	userRouter.HandleFunc("/{userId}/teams-non-org", a.userOnly(a.entityUserOnly(a.handleGetTeamsByUserNonOrg()))).Methods("GET")
	if a.Config.SubscriptionsEnabled {
		userRouter.HandleFunc("/{userId}/subscriptions", a.userOnly(a.entityUserOnly(a.handleGetEntityUserActiveSubs()))).Methods("GET")
		userRouter.HandleFunc("/{userId}/subscriptions/{subscriptionId}", a.userOnly(a.entityUserOnly(a.handleEntityUserUpdateSubscription()))).Methods("PATCH")
	}
	userRouter.HandleFunc("/{userId}/jira-instances", a.userOnly(a.entityUserOnly(a.subscribedEntityUserOnly(a.handleGetUserJiraInstances())))).Methods("GET")
	userRouter.HandleFunc("/{userId}/jira-instances", a.userOnly(a.entityUserOnly(a.subscribedEntityUserOnly(a.handleJiraInstanceCreate())))).Methods("POST")
	userRouter.HandleFunc("/{userId}/jira-instances/{instanceId}", a.userOnly(a.entityUserOnly(a.subscribedEntityUserOnly(a.handleJiraInstanceUpdate())))).Methods("PUT")
	userRouter.HandleFunc("/{userId}/jira-instances/{instanceId}", a.userOnly(a.entityUserOnly(a.subscribedEntityUserOnly(a.handleJiraInstanceDelete())))).Methods("DELETE")
	userRouter.HandleFunc("/{userId}/jira-instances/{instanceId}/jql-story-search", a.userOnly(a.entityUserOnly(a.subscribedEntityUserOnly(a.handleJiraStoryJQLSearch())))).Methods("POST")

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
	orgRouter.HandleFunc("/{orgId}", a.userOnly(a.orgAdminOnly(a.handleOrganizationUpdate()))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}", a.userOnly(a.orgAdminOnly(a.handleDeleteOrganization()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/metrics", a.userOnly(a.orgUserOnly(a.handleOrganizationMetrics()))).Methods("GET")
	// org departments(s)
	orgRouter.HandleFunc("/{orgId}/departments", a.userOnly(a.orgUserOnly(a.handleGetOrganizationDepartments()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments", a.userOnly(a.orgAdminOnly(a.handleCreateDepartment()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}", a.userOnly(a.orgAdminOnly(a.handleDepartmentUpdate()))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}", a.userOnly(a.orgAdminOnly(a.handleDeleteDepartment()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/invites", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentUserInvites()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/invites", a.userOnly(a.departmentAdminOnly(a.handleDepartmentInviteUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/invites/{inviteId}", a.userOnly(a.departmentAdminOnly(a.handleDeleteDepartmentUserInvite()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users", a.userOnly(a.departmentAdminOnly(a.handleDepartmentAddUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users/{userId}", a.userOnly(a.departmentAdminOnly(a.handleDepartmentUpdateUser()))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/users/{userId}", a.userOnly(a.departmentAdminOnly(a.handleDepartmentRemoveUser()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentTeams()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams", a.userOnly(a.departmentAdminOnly(a.handleCreateDepartmentTeam()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}", a.userOnly(a.teamUserOnly(a.handleDepartmentTeamByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}", a.userOnly(a.departmentAdminOnly(a.handleTeamUpdate()))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}", a.userOnly(a.departmentAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleGetTeamUserInvites()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleTeamInviteUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/invites/{inviteId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleDeleteTeamUserInvite())))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users", a.userOnly(a.teamUserOnly(a.handleGetTeamUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleDepartmentTeamAddUser())))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamUpdateUser())))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveUser())))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinsGet()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinCreate(checkinSvc)))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/users/{userId}/last", a.userOnly(a.subscribedUserOnly(a.teamUserOnly(a.handleCheckinLastByUser())))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinUpdate(checkinSvc)))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinDelete(checkinSvc)))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}/comments", a.userOnly(a.teamUserOnly(a.handleCheckinComment(checkinSvc)))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentEdit(checkinSvc)))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentDelete(checkinSvc)))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/metrics", a.userOnly(a.teamUserOnly(a.handleTeamMetrics()))).Methods("GET")
	// org teams
	orgRouter.HandleFunc("/{orgId}/teams", a.userOnly(a.orgUserOnly(a.handleGetOrganizationTeams()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams", a.userOnly(a.orgAdminOnly(a.handleCreateOrganizationTeam()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}", a.userOnly(a.teamUserOnly(a.handleGetOrganizationTeamByUser()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}", a.userOnly(a.orgAdminOnly(a.handleTeamUpdate()))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}", a.userOnly(a.orgAdminOnly(a.handleDeleteTeam()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleGetTeamUserInvites()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleTeamInviteUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/invites/{inviteId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleDeleteTeamUserInvite())))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users", a.userOnly(a.teamUserOnly(a.handleGetTeamUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleOrganizationTeamAddUser())))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamUpdateUser())))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveUser())))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinsGet()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinCreate(checkinSvc)))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/users/{userId}/last", a.userOnly(a.subscribedUserOnly(a.teamUserOnly(a.handleCheckinLastByUser())))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinUpdate(checkinSvc)))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinDelete(checkinSvc)))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/{checkinId}/comments", a.userOnly(a.teamUserOnly(a.handleCheckinComment(checkinSvc)))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentEdit(checkinSvc)))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentDelete(checkinSvc)))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/teams/{teamId}/metrics", a.userOnly(a.teamUserOnly(a.handleTeamMetrics()))).Methods("GET")
	// org users
	orgRouter.HandleFunc("/{orgId}/users", a.userOnly(a.orgUserOnly(a.handleGetOrganizationUsers()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/users/{userId}", a.userOnly(a.orgAdminOnly(a.handleOrganizationUpdateUser()))).Methods("PUT")
	orgRouter.HandleFunc("/{orgId}/users/{userId}", a.userOnly(a.orgAdminOnly(a.handleOrganizationRemoveUser()))).Methods("DELETE")
	orgRouter.HandleFunc("/{orgId}/invites", a.userOnly(a.orgUserOnly(a.handleGetOrganizationUserInvites()))).Methods("GET")
	orgRouter.HandleFunc("/{orgId}/invites", a.userOnly(a.orgAdminOnly(a.handleOrganizationInviteUser()))).Methods("POST")
	orgRouter.HandleFunc("/{orgId}/invites/{inviteId}", a.userOnly(a.orgAdminOnly(a.handleDeleteOrganizationUserInvite()))).Methods("DELETE")
	// teams(s)
	teamRouter.HandleFunc("/{teamId}", a.userOnly(a.teamUserOnly(a.handleGetTeamByUser()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamUpdate())))).Methods("PUT")
	teamRouter.HandleFunc("/{teamId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleDeleteTeam())))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleGetTeamUserInvites()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}/invites", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamInviteUser())))).Methods("POST")
	teamRouter.HandleFunc("/{teamId}/invites/{inviteId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleDeleteTeamUserInvite())))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/users", a.userOnly(a.teamUserOnly(a.handleGetTeamUsers()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamUpdateUser())))).Methods("PUT")
	teamRouter.HandleFunc("/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveUser())))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/checkin", checkinSvc.ServeWs())
	teamRouter.HandleFunc("/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinsGet()))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinCreate(checkinSvc)))).Methods("POST")
	teamRouter.HandleFunc("/{teamId}/checkins/users/{userId}/last", a.userOnly(a.subscribedUserOnly(a.teamUserOnly(a.handleCheckinLastByUser())))).Methods("GET")
	teamRouter.HandleFunc("/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinUpdate(checkinSvc)))).Methods("PUT")
	teamRouter.HandleFunc("/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinDelete(checkinSvc)))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/checkins/{checkinId}/comments", a.userOnly(a.teamUserOnly(a.handleCheckinComment(checkinSvc)))).Methods("POST")
	teamRouter.HandleFunc("/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentEdit(checkinSvc)))).Methods("PUT")
	teamRouter.HandleFunc("/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentDelete(checkinSvc)))).Methods("DELETE")
	teamRouter.HandleFunc("/{teamId}/metrics", a.userOnly(a.teamUserOnly(a.handleTeamMetrics()))).Methods("GET")
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
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/battles", a.userOnly(a.teamUserOnly(a.handleGetTeamBattles()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/battles/{battleId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveBattle())))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/battles", a.userOnly(a.teamUserOnly(a.handlePokerCreate()))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/battles", a.userOnly(a.teamUserOnly(a.handleGetTeamBattles()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/battles/{battleId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveBattle())))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}/battles", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handlePokerCreate())))).Methods("POST")
		teamRouter.HandleFunc("/{teamId}/battles", a.userOnly(a.teamUserOnly(a.handleGetTeamBattles()))).Methods("GET")
		teamRouter.HandleFunc("/{teamId}/battles/{battleId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveBattle())))).Methods("DELETE")
		teamRouter.HandleFunc("/{teamId}/users/{userId}/battles", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handlePokerCreate())))).Methods("POST")
		apiRouter.HandleFunc("/maintenance/clean-battles", a.userOnly(a.adminOnly(a.handleCleanBattles()))).Methods("DELETE")
		apiRouter.HandleFunc("/battles", a.userOnly(a.adminOnly(a.handleGetPokerGames()))).Methods("GET")
		apiRouter.HandleFunc("/battles/{battleId}", a.userOnly(a.handleGetPokerGame())).Methods("GET")
		apiRouter.HandleFunc("/battles/{battleId}", a.userOnly(a.handlePokerDelete(pokerSvc))).Methods("DELETE")
		apiRouter.HandleFunc("/battles/{battleId}/plans", a.userOnly(a.handlePokerStoryAdd(pokerSvc))).Methods("POST")
		apiRouter.HandleFunc("/battles/{battleId}/plans/{planId}", a.userOnly(a.handlePokerStoryDelete(pokerSvc))).Methods("DELETE")
		apiRouter.HandleFunc("/arena/{battleId}", pokerSvc.ServeBattleWs())

		// estimation scales
		// Public estimation scale routes
		apiRouter.HandleFunc("/estimation-scales/public", a.userOnly(a.handleGetPublicEstimationScales())).Methods("GET")
		apiRouter.HandleFunc("/estimation-scales/public/{scaleId}", a.userOnly(a.handleGetPublicEstimationScale())).Methods("GET")

		// Organization-specific estimation scale routes
		orgRouter.HandleFunc("/{orgId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.orgUserOnly(a.handleGetOrganizationEstimationScales())))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationEstimationScaleCreate())))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationEstimationScaleUpdate())))).Methods("PUT")
		orgRouter.HandleFunc("/{orgId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationEstimationScaleDelete())))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.departmentUserOnly(a.handleGetTeamEstimationScales())))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleCreate()))))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleUpdate()))))).Methods("PUT")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleDelete()))))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.handleGetTeamEstimationScales())))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleCreate()))))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleUpdate()))))).Methods("PUT")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleDelete()))))).Methods("DELETE")

		// Team-specific estimation scale routes
		teamRouter.HandleFunc("/{teamId}/estimation-scales", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.handleGetTeamEstimationScales())))).Methods("GET")
		teamRouter.HandleFunc("/{teamId}/estimation-scales", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleCreate()))))).Methods("POST")
		teamRouter.HandleFunc("/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleUpdate()))))).Methods("PUT")
		teamRouter.HandleFunc("/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleDelete()))))).Methods("DELETE")

		// Admin estimation scale routes
		adminRouter.HandleFunc("/estimation-scales", a.userOnly(a.adminOnly(a.handleGetEstimationScales()))).Methods("GET")
		adminRouter.HandleFunc("/estimation-scales", a.userOnly(a.adminOnly(a.handleEstimationScaleCreate()))).Methods("POST")
		adminRouter.HandleFunc("/estimation-scales/{scaleId}", a.userOnly(a.adminOnly(a.handleEstimationScaleUpdate()))).Methods("PUT")
		adminRouter.HandleFunc("/estimation-scales/{scaleId}", a.userOnly(a.adminOnly(a.handleEstimationScaleDelete()))).Methods("DELETE")
	}
	// retro(s)
	if a.Config.FeatureRetro {
		userRouter.HandleFunc("/{userId}/retros", a.userOnly(a.entityUserOnly(a.handleRetroCreate()))).Methods("POST")
		userRouter.HandleFunc("/{userId}/retros", a.userOnly(a.entityUserOnly(a.handleRetrosGetByUser()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/retros", a.userOnly(a.teamUserOnly(a.handleGetTeamRetros()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/retros/{retroId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveRetro())))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/retro-actions", a.userOnly(a.teamUserOnly(a.handleGetTeamRetroActions()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/retros", a.userOnly(a.teamUserOnly(a.handleRetroCreate()))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/retros", a.userOnly(a.teamUserOnly(a.handleGetTeamRetros()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/retro-actions", a.userOnly(a.teamUserOnly(a.handleGetTeamRetroActions()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/retros/{retroId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveRetro())))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}/retros", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handleRetroCreate())))).Methods("POST")
		teamRouter.HandleFunc("/{teamId}/retros", a.userOnly(a.teamUserOnly(a.handleGetTeamRetros()))).Methods("GET")
		teamRouter.HandleFunc("/{teamId}/retros/{retroId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveRetro())))).Methods("DELETE")
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

		// Retro Templates
		apiRouter.HandleFunc("/retro-templates/public", a.userOnly(a.handleGetPublicRetroTemplates())).Methods("GET")
		// Organization templates
		orgRouter.HandleFunc("/{orgId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.orgUserOnly(a.handleGetOrganizationRetroTemplates())))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationRetroTemplateCreate())))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationRetroTemplateUpdate())))).Methods("PUT")
		orgRouter.HandleFunc("/{orgId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationRetroTemplateDelete())))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.departmentUserOnly(a.handleGetTeamRetroTemplates())))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateCreate()))))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateUpdate()))))).Methods("PUT")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateDelete()))))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.handleGetTeamRetroTemplates())))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateCreate()))))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateUpdate()))))).Methods("PUT")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateDelete()))))).Methods("DELETE")
		// Team templates
		teamRouter.HandleFunc("/{teamId}/retro-templates", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.handleGetTeamRetroTemplates())))).Methods("GET")
		teamRouter.HandleFunc("/{teamId}/retro-templates", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateCreate()))))).Methods("POST")
		teamRouter.HandleFunc("/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateUpdate()))))).Methods("PUT")
		teamRouter.HandleFunc("/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateDelete()))))).Methods("DELETE")
		// General template operations
		adminRouter.HandleFunc("/retro-templates", a.userOnly(a.adminOnly(a.handleGetRetroTemplates()))).Methods("GET")
		adminRouter.HandleFunc("/retro-templates/{templateId}", a.userOnly(a.adminOnly(a.handleGetRetroTemplateById()))).Methods("GET")
		adminRouter.HandleFunc("/retro-templates", a.userOnly(a.adminOnly(a.handleRetroTemplateCreate()))).Methods("POST")
		adminRouter.HandleFunc("/retro-templates/{templateId}", a.userOnly(a.adminOnly(a.handleRetroTemplateUpdate()))).Methods("PUT")
		adminRouter.HandleFunc("/retro-templates/{templateId}", a.userOnly(a.adminOnly(a.handleRetroTemplateDelete()))).Methods("DELETE")
		// Retro websocket
		apiRouter.HandleFunc("/retro/{retroId}", retroSvc.ServeWs())
	}
	// storyboard(s)
	if a.Config.FeatureStoryboard {
		userRouter.HandleFunc("/{userId}/storyboards", a.userOnly(a.entityUserOnly(a.handleStoryboardCreate()))).Methods("POST")
		userRouter.HandleFunc("/{userId}/storyboards", a.userOnly(a.entityUserOnly(a.handleGetUserStoryboards()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/storyboards", a.userOnly(a.teamUserOnly(a.handleGetTeamStoryboards()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/storyboards/{storyboardId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveStoryboard())))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/storyboards", a.userOnly(a.teamUserOnly(a.handleStoryboardCreate()))).Methods("POST")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/storyboards", a.userOnly(a.teamUserOnly(a.handleGetTeamStoryboards()))).Methods("GET")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/storyboards/{storyboardId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveStoryboard())))).Methods("DELETE")
		orgRouter.HandleFunc("/{orgId}/teams/{teamId}/users/{userId}/storyboards", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handleStoryboardCreate())))).Methods("POST")
		teamRouter.HandleFunc("/{teamId}/storyboards", a.userOnly(a.teamUserOnly(a.handleGetTeamStoryboards()))).Methods("GET")
		teamRouter.HandleFunc("/{teamId}/storyboards/{storyboardId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveStoryboard())))).Methods("DELETE")
		teamRouter.HandleFunc("/{teamId}/users/{userId}/storyboards", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handleStoryboardCreate())))).Methods("POST")
		apiRouter.HandleFunc("/maintenance/clean-storyboards", a.userOnly(a.adminOnly(a.handleCleanStoryboards()))).Methods("DELETE")
		apiRouter.HandleFunc("/storyboards", a.userOnly(a.adminOnly(a.handleGetStoryboards()))).Methods("GET")
		apiRouter.HandleFunc("/storyboards/{storyboardId}", a.userOnly(a.handleStoryboardGet())).Methods("GET")
		apiRouter.HandleFunc("/storyboards/{storyboardId}", a.userOnly(a.handleStoryboardDelete(storyboardSvc))).Methods("DELETE")
		apiRouter.HandleFunc("/storyboards/{storyboardId}/goals", a.userOnly(a.handleStoryboardGoalAdd(storyboardSvc))).Methods("POST")
		apiRouter.HandleFunc("/storyboards/{storyboardId}/columns", a.userOnly(a.handleStoryboardColumnAdd(storyboardSvc))).Methods("POST")
		apiRouter.HandleFunc("/storyboards/{storyboardId}/stories", a.userOnly(a.handleStoryboardStoryAdd(storyboardSvc))).Methods("POST")
		apiRouter.HandleFunc("/storyboards/{storyboardId}/stories/{storyId}/move", a.userOnly(a.handleStoryboardStoryMove(storyboardSvc))).Methods("PUT")
		apiRouter.HandleFunc("/storyboard/{storyboardId}", storyboardSvc.ServeWs())
	}

	// user avatar generation
	if a.Config.AvatarService == "goadorable" || a.Config.AvatarService == "govatar" {
		router.PathPrefix("/avatar/{width}/{id}/{avatar}").Handler(a.handleUserAvatar()).Methods("GET")
		router.PathPrefix("/avatar/{width}/{id}").Handler(a.handleUserAvatar()).Methods("GET")
	}

	if a.Config.SubscriptionsEnabled {
		apiRouter.PathPrefix("/subscriptions/{subscriptionId}").Handler(a.userOnly(a.adminOnly(a.handleSubscriptionGet()))).Methods("GET")
		apiRouter.PathPrefix("/subscriptions/{subscriptionId}").Handler(a.userOnly(a.adminOnly(a.handleSubscriptionUpdate()))).Methods("PUT")
		apiRouter.PathPrefix("/subscriptions/{subscriptionId}").Handler(a.userOnly(a.adminOnly(a.handleSubscriptionDelete()))).Methods("DELETE")
		apiRouter.PathPrefix("/subscriptions").Handler(a.userOnly(a.adminOnly(a.handleGetSubscriptions()))).Methods("GET")
		apiRouter.PathPrefix("/subscriptions").Handler(a.userOnly(a.adminOnly(a.handleSubscriptionCreate()))).Methods("POST")
		router.PathPrefix("/webhooks/subscriptions").Handler(a.SubscriptionSvc.HandleWebhook()).Methods("POST")
	}

	a.registerOauthProviderEndpoints(authProviderConfigs)

	// static assets
	router.PathPrefix("/static/").Handler(http.StripPrefix(a.Config.PathPrefix, staticHandler))
	router.PathPrefix("/img/").Handler(http.StripPrefix(a.Config.PathPrefix, staticHandler))

	// handle index.html
	router.PathPrefix("/").HandlerFunc(a.handleIndex(FSS, a.UIConfig))

	return a
}

func (s *Service) registerOauthProviderEndpoints(providers []thunderdome.AuthProviderConfig) {
	ctx := context.Background()
	var redirectBaseURL string
	var port string

	// redirect with port for localhost
	if s.Config.AppDomain == "localhost" {
		port = fmt.Sprintf(":%s", s.Config.Port)
	}

	if s.Config.SecureProtocol {
		redirectBaseURL = fmt.Sprintf("https://%s%s", s.Config.AppDomain, port)
	} else {
		redirectBaseURL = fmt.Sprintf("http://%s%s%s", s.Config.AppDomain, port, s.Config.PathPrefix)
	}

	for _, c := range providers {
		oauthLoginPathPrefix, _ := url.JoinPath("/oauth/", c.ProviderName, "/login")
		oauthCallbackPathPrefix, _ := url.JoinPath("/oauth/", c.ProviderName, "/callback")
		callbackRedirectURL, _ := url.JoinPath(redirectBaseURL, oauthCallbackPathPrefix)
		authProvider, err := oauth.New(oauth.Config{
			AuthProviderConfig:  c,
			CallbackRedirectURL: callbackRedirectURL,
			UIRedirectURL:       fmt.Sprintf("%s/", s.Config.PathPrefix),
		}, s.Cookie, s.Logger, s.AuthDataSvc, s.SubscriptionDataSvc, ctx)
		if err != nil {
			panic(err)
		}
		s.Router.HandleFunc(oauthLoginPathPrefix, authProvider.HandleOAuth2Redirect()).Methods("GET")
		s.Router.HandleFunc(oauthCallbackPathPrefix, authProvider.HandleOAuth2Callback()).Methods("GET")
	}
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
		nonce := secure.CSPNonce(r.Context())

		if s.Config.EmbedUseOS {
			tmpl = s.getIndexTemplate(FSS)
		}

		type templateData struct {
			UIConfig thunderdome.UIConfig
			Nonce    string
		}
		td := templateData{
			UIConfig: uiConfig,
			Nonce:    nonce,
		}

		err := tmpl.Execute(w, td)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

// handleApiIndex returns a handler for the API index route
func (s *Service) handleApiIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}
}
