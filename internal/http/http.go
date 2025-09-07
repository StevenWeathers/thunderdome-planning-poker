// Package http provides restful API endpoints for Thunderdome webapp
package http

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/unrolled/secure"
	"github.com/unrolled/secure/cspbuilder"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/oauth"

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
//
//	@title						Thunderdome API
//	@description				Thunderdome Planning Poker API for both Internal and External use.
//	@description				WARNING: Currently not considered stable and is subject to change until 1.0 is released.
//	@contact.name				Steven Weathers
//	@contact.url				https://github.com/StevenWeathers/thunderdome-planning-poker
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@version					BETA
//	@query.collection.format	multi
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-Key
func New(apiService Service, FSS fs.FS, HFS http.FileSystem) *Service {
	staticHandler := http.FileServer(HFS)

	var a = &apiService
	authProviderConfigs := make([]thunderdome.AuthProviderConfig, 0)
	connectSrcCsp := []string{
		"'self'",
		getWebsocketConnectSrc(a.Config.SecureProtocol, a.Config.WebsocketSubdomain, a.Config.AppDomain),
	}

	// Content Security Policy
	cspBuilder := cspbuilder.Builder{
		Directives: map[string][]string{
			cspbuilder.DefaultSrc: {"'self'", fmt.Sprintf("*.%s", a.Config.AppDomain)},
			//	@TODO	- remove inline styles in svelte components to improve security by using nonce
			cspbuilder.StyleSrc:    {"'self'", "'unsafe-inline'", "https://fonts.googleapis.com"},
			cspbuilder.ScriptSrc:   {"$NONCE"},
			cspbuilder.FontSrc:     {"'self'", "https://fonts.gstatic.com"},
			cspbuilder.ImgSrc:      {"data:", "*"},
			cspbuilder.ConnectSrc:  connectSrcCsp,
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

	prefix := apiService.Config.PathPrefix

	router := http.NewServeMux()

	swaggerJsonPath := a.Config.PathPrefix + "/swagger/doc.json"
	swagger.SwaggerInfo.BasePath = a.Config.PathPrefix + "/api"
	// swagger docs for external API when enabled
	if a.Config.ExternalAPIEnabled {
		router.Handle(prefix+"/swagger/", httpSwagger.Handler(httpSwagger.URL(swaggerJsonPath)))
	}

	// user authentication, profile
	if a.Config.LdapEnabled {
		router.Handle("POST "+prefix+"/api/auth/ldap", a.handleLdapLogin())
	} else if a.Config.HeaderAuthEnabled {
		router.Handle("GET "+prefix+"/api/auth", a.handleHeaderLogin())
	} else if a.Config.OIDCAuth.Enabled {
		authProviderConfigs = append(authProviderConfigs, thunderdome.AuthProviderConfig{
			ProviderName:           a.Config.OIDCAuth.ProviderName,
			ProviderURL:            a.Config.OIDCAuth.ProviderURL,
			ClientID:               a.Config.OIDCAuth.ClientID,
			ClientSecret:           a.Config.OIDCAuth.ClientSecret,
			RequestedScopes:        a.Config.OIDCAuth.RequestedScopes,
			RequestedIDTokenClaims: a.Config.OIDCAuth.RequestedIDTokenClaims,
		})
	} else {
		if a.Config.GoogleAuth.Enabled {
			authProviderConfigs = append(authProviderConfigs, thunderdome.AuthProviderConfig{
				ProviderName:           a.Config.GoogleAuth.ProviderName,
				ProviderURL:            a.Config.GoogleAuth.ProviderURL,
				ClientID:               a.Config.GoogleAuth.ClientID,
				ClientSecret:           a.Config.GoogleAuth.ClientSecret,
				RequestedScopes:        a.Config.GoogleAuth.RequestedScopes,
				RequestedIDTokenClaims: a.Config.GoogleAuth.RequestedIDTokenClaims,
			})
		}
		router.Handle("POST "+prefix+"/api/auth", a.handleLogin())
		router.Handle("POST "+prefix+"/api/auth/forgot-password", a.handleForgotPassword())
		router.Handle("PATCH "+prefix+"/api/auth/reset-password", a.handleResetPassword())
		router.Handle("PATCH "+prefix+"/api/auth/update-password", a.userOnly(a.handleUpdatePassword()))
		router.Handle("PATCH "+prefix+"/api/auth/verify", a.handleAccountVerification())
		router.Handle("POST "+prefix+"/api/auth/register", a.handleUserRegistration())
		router.Handle("GET "+prefix+"/api/auth/invite/team/{inviteId}", a.handleGetTeamInviteByID())
		router.Handle("GET "+prefix+"/api/auth/invite/organization/{inviteId}", a.handleGetOrganizationInviteByID())
	}
	router.Handle("POST "+prefix+"/api/auth/mfa", a.handleMFALogin())
	router.Handle("DELETE "+prefix+"/api/auth/mfa", a.userOnly(a.registeredUserOnly(a.handleMFARemove())))
	router.Handle("POST "+prefix+"/api/auth/mfa/setup/generate", a.userOnly(a.registeredUserOnly(a.handleMFASetupGenerate())))
	router.Handle("POST "+prefix+"/api/auth/mfa/setup/validate", a.userOnly(a.registeredUserOnly(a.handleMFASetupValidate())))
	router.Handle("POST "+prefix+"/api/auth/guest", a.handleCreateGuestUser())
	router.Handle("GET "+prefix+"/api/auth/user", a.userOnly(a.handleSessionUserProfile()))
	router.Handle("DELETE "+prefix+"/api/auth/logout", a.handleLogout())
	// user(s)
	router.Handle("GET "+prefix+"/api/users/{userId}", a.userOnly(a.entityUserOnly(a.handleUserProfile())))
	router.Handle("PUT "+prefix+"/api/users/{userId}", a.userOnly(a.entityUserOnly(a.handleUserProfileUpdate())))
	router.Handle("DELETE "+prefix+"/api/users/{userId}", a.userOnly(a.entityUserOnly(a.handleUserDelete())))
	router.Handle("GET "+prefix+"/api/users/{userId}/credential", a.userOnly(a.entityUserOnly(a.handleUserCredential())))
	router.Handle("POST "+prefix+"/api/users/{userId}/request-verify", a.userOnly(a.entityUserOnly(a.handleVerifyRequest())))
	router.Handle("POST "+prefix+"/api/users/{userId}/email-change", a.userOnly(a.entityUserOnly(a.handleChangeEmailRequest())))
	router.Handle("POST "+prefix+"/api/users/{userId}/email-change/{changeId}", a.userOnly(a.entityUserOnly(a.handleChangeEmailAction())))
	router.Handle("POST "+prefix+"/api/users/{userId}/invite/team/{inviteId}", a.userOnly(a.registeredUserOnly(a.handleUserTeamInvite())))
	router.Handle("POST "+prefix+"/api/users/{userId}/invite/organization/{inviteId}", a.userOnly(a.registeredUserOnly(a.handleUserOrganizationInvite())))
	router.Handle("POST "+prefix+"/api/users/{userId}/invite/department/{inviteId}", a.userOnly(a.registeredUserOnly(a.handleUserDepartmentInvite())))
	router.Handle("GET "+prefix+"/api/users/{userId}/organizations", a.userOnly(a.entityUserOnly(a.handleGetOrganizationsByUser())))
	router.Handle("POST "+prefix+"/api/users/{userId}/organizations", a.userOnly(a.entityUserOnly(a.handleCreateOrganization())))
	router.Handle("GET "+prefix+"/api/users/{userId}/teams", a.userOnly(a.entityUserOnly(a.handleGetTeamsByUser())))
	router.Handle("POST "+prefix+"/api/users/{userId}/teams", a.userOnly(a.entityUserOnly(a.handleCreateTeam())))
	router.Handle("GET "+prefix+"/api/users/{userId}/teams-non-org", a.userOnly(a.entityUserOnly(a.handleGetTeamsByUserNonOrg())))
	if a.Config.SubscriptionsEnabled {
		router.Handle("GET "+prefix+"/api/users/{userId}/subscriptions", a.userOnly(a.entityUserOnly(a.handleGetEntityUserActiveSubs())))
		router.Handle("PATCH "+prefix+"/api/users/{userId}/subscriptions/{subscriptionId}", a.userOnly(a.entityUserOnly(a.handleEntityUserUpdateSubscription())))
	}
	router.Handle("GET "+prefix+"/api/users/{userId}/jira-instances", a.userOnly(a.entityUserOnly(a.subscribedEntityUserOnly(a.handleGetUserJiraInstances()))))
	router.Handle("POST "+prefix+"/api/users/{userId}/jira-instances", a.userOnly(a.entityUserOnly(a.subscribedEntityUserOnly(a.handleJiraInstanceCreate()))))
	router.Handle("PUT "+prefix+"/api/users/{userId}/jira-instances/{instanceId}", a.userOnly(a.entityUserOnly(a.subscribedEntityUserOnly(a.handleJiraInstanceUpdate()))))
	router.Handle("DELETE "+prefix+"/api/users/{userId}/jira-instances/{instanceId}", a.userOnly(a.entityUserOnly(a.subscribedEntityUserOnly(a.handleJiraInstanceDelete()))))
	router.Handle("POST "+prefix+"/api/users/{userId}/jira-instances/{instanceId}/jql-story-search", a.userOnly(a.entityUserOnly(a.subscribedEntityUserOnly(a.handleJiraStoryJQLSearch()))))

	if a.Config.ExternalAPIEnabled {
		router.Handle("GET "+prefix+"/api/users/{userId}/apikeys", a.userOnly(a.entityUserOnly(a.handleUserAPIKeys())))
		router.Handle("POST "+prefix+"/api/users/{userId}/apikeys", a.userOnly(a.verifiedUserOnly(a.handleAPIKeyGenerate())))
		router.Handle("PUT "+prefix+"/api/users/{userId}/apikeys/{keyID}", a.userOnly(a.entityUserOnly(a.handleUserAPIKeyUpdate())))
		router.Handle("DELETE "+prefix+"/api/users/{userId}/apikeys/{keyID}", a.userOnly(a.entityUserOnly(a.handleUserAPIKeyDelete())))
	}
	// country(s)
	if a.Config.ShowActiveCountries {
		router.Handle("GET "+prefix+"/api/active-countries", a.handleGetActiveCountries())
	}
	// org
	router.Handle("GET "+prefix+"/api/organizations/{orgId}", a.userOnly(a.orgUserOnly(a.handleGetOrganizationByUser())))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}", a.userOnly(a.orgAdminOnly(a.handleOrganizationUpdate())))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}", a.userOnly(a.orgAdminOnly(a.handleDeleteOrganization())))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/metrics", a.userOnly(a.orgUserOnly(a.handleOrganizationMetrics())))
	// org departments(s)
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments", a.userOnly(a.orgUserOnly(a.handleGetOrganizationDepartments())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments", a.userOnly(a.orgAdminOnly(a.handleCreateDepartment())))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentByUser())))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}", a.userOnly(a.orgAdminOnly(a.handleDepartmentUpdate())))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}", a.userOnly(a.orgAdminOnly(a.handleDeleteDepartment())))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/invites", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentUserInvites())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/invites", a.userOnly(a.departmentAdminOnly(a.handleDepartmentInviteUser())))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/invites/{inviteId}", a.userOnly(a.departmentAdminOnly(a.handleDeleteDepartmentUserInvite())))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/users", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentUsers())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/users", a.userOnly(a.departmentAdminOnly(a.handleDepartmentAddUser())))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/users/{userId}", a.userOnly(a.departmentAdminOnly(a.handleDepartmentUpdateUser())))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/users/{userId}", a.userOnly(a.departmentAdminOnly(a.handleDepartmentRemoveUser())))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams", a.userOnly(a.departmentUserOnly(a.handleGetDepartmentTeams())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams", a.userOnly(a.departmentAdminOnly(a.handleCreateDepartmentTeam())))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}", a.userOnly(a.teamUserOnly(a.handleDepartmentTeamByUser())))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}", a.userOnly(a.departmentAdminOnly(a.handleTeamUpdate())))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}", a.userOnly(a.departmentAdminOnly(a.handleDeleteTeam())))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleGetTeamUserInvites())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleTeamInviteUser())))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/invites/{inviteId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleDeleteTeamUserInvite()))))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users", a.userOnly(a.teamUserOnly(a.handleGetTeamUsers())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleDepartmentTeamAddUser()))))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamUpdateUser()))))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveUser()))))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinsGet())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinCreate(checkinSvc))))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/users/{userId}/last", a.userOnly(a.subscribedUserOnly(a.teamUserOnly(a.handleCheckinLastByUser()))))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinUpdate(checkinSvc))))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinDelete(checkinSvc))))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}/comments", a.userOnly(a.teamUserOnly(a.handleCheckinComment(checkinSvc))))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentEdit(checkinSvc))))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentDelete(checkinSvc))))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/metrics", a.userOnly(a.teamUserOnly(a.handleTeamMetrics())))
	// org teams
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams", a.userOnly(a.orgUserOnly(a.handleGetOrganizationTeams())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/teams", a.userOnly(a.orgAdminOnly(a.handleCreateOrganizationTeam())))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}", a.userOnly(a.teamUserOnly(a.handleGetOrganizationTeamByUser())))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/teams/{teamId}", a.userOnly(a.orgAdminOnly(a.handleTeamUpdate())))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/teams/{teamId}", a.userOnly(a.orgAdminOnly(a.handleDeleteTeam())))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleGetTeamUserInvites())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/teams/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleTeamInviteUser())))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/teams/{teamId}/invites/{inviteId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleDeleteTeamUserInvite()))))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/users", a.userOnly(a.teamUserOnly(a.handleGetTeamUsers())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/teams/{teamId}/users", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleOrganizationTeamAddUser()))))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/teams/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamUpdateUser()))))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/teams/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveUser()))))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinsGet())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/teams/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinCreate(checkinSvc))))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/checkins/users/{userId}/last", a.userOnly(a.subscribedUserOnly(a.teamUserOnly(a.handleCheckinLastByUser()))))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinUpdate(checkinSvc))))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinDelete(checkinSvc))))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/teams/{teamId}/checkins/{checkinId}/comments", a.userOnly(a.teamUserOnly(a.handleCheckinComment(checkinSvc))))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentEdit(checkinSvc))))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentDelete(checkinSvc))))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/metrics", a.userOnly(a.teamUserOnly(a.handleTeamMetrics())))
	// org users
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/users", a.userOnly(a.orgUserOnly(a.handleGetOrganizationUsers())))
	router.Handle("PUT "+prefix+"/api/organizations/{orgId}/users/{userId}", a.userOnly(a.orgAdminOnly(a.handleOrganizationUpdateUser())))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/users/{userId}", a.userOnly(a.orgAdminOnly(a.handleOrganizationRemoveUser())))
	router.Handle("GET "+prefix+"/api/organizations/{orgId}/invites", a.userOnly(a.orgUserOnly(a.handleGetOrganizationUserInvites())))
	router.Handle("POST "+prefix+"/api/organizations/{orgId}/invites", a.userOnly(a.orgAdminOnly(a.handleOrganizationInviteUser())))
	router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/invites/{inviteId}", a.userOnly(a.orgAdminOnly(a.handleDeleteOrganizationUserInvite())))
	// teams(s)
	router.Handle("GET "+prefix+"/api/teams/{teamId}", a.userOnly(a.teamUserOnly(a.handleGetTeamByUser())))
	router.Handle("PUT "+prefix+"/api/teams/{teamId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamUpdate()))))
	router.Handle("DELETE "+prefix+"/api/teams/{teamId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleDeleteTeam()))))
	router.Handle("GET "+prefix+"/api/teams/{teamId}/invites", a.userOnly(a.teamUserOnly(a.handleGetTeamUserInvites())))
	router.Handle("POST "+prefix+"/api/teams/{teamId}/invites", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamInviteUser()))))
	router.Handle("DELETE "+prefix+"/api/teams/{teamId}/invites/{inviteId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleDeleteTeamUserInvite()))))
	router.Handle("GET "+prefix+"/api/teams/{teamId}/users", a.userOnly(a.teamUserOnly(a.handleGetTeamUsers())))
	router.Handle("PUT "+prefix+"/api/teams/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamUpdateUser()))))
	router.Handle("DELETE "+prefix+"/api/teams/{teamId}/users/{userId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveUser()))))
	router.Handle(prefix+"/api/teams/{teamId}/checkin", checkinSvc.ServeWs())
	router.Handle("GET "+prefix+"/api/teams/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinsGet())))
	router.Handle("POST "+prefix+"/api/teams/{teamId}/checkins", a.userOnly(a.teamUserOnly(a.handleCheckinCreate(checkinSvc))))
	router.Handle("GET "+prefix+"/api/teams/{teamId}/checkins/users/{userId}/last", a.userOnly(a.subscribedUserOnly(a.teamUserOnly(a.handleCheckinLastByUser()))))
	router.Handle("PUT "+prefix+"/api/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinUpdate(checkinSvc))))
	router.Handle("DELETE "+prefix+"/api/teams/{teamId}/checkins/{checkinId}", a.userOnly(a.teamUserOnly(a.handleCheckinDelete(checkinSvc))))
	router.Handle("POST "+prefix+"/api/teams/{teamId}/checkins/{checkinId}/comments", a.userOnly(a.teamUserOnly(a.handleCheckinComment(checkinSvc))))
	router.Handle("PUT "+prefix+"/api/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentEdit(checkinSvc))))
	router.Handle("DELETE "+prefix+"/api/teams/{teamId}/checkins/{checkinId}/comments/{commentId}", a.userOnly(a.teamUserOnly(a.handleCheckinCommentDelete(checkinSvc))))
	router.Handle("GET "+prefix+"/api/teams/{teamId}/metrics", a.userOnly(a.teamUserOnly(a.handleTeamMetrics())))
	// admin
	router.Handle("GET "+prefix+"/api/admin/stats", a.userOnly(a.adminOnly(a.handleAppStats())))
	router.Handle("GET "+prefix+"/api/admin/users", a.userOnly(a.adminOnly(a.handleGetRegisteredUsers())))
	router.Handle("POST "+prefix+"/api/admin/users", a.userOnly(a.adminOnly(a.handleUserCreate())))
	router.Handle("PATCH "+prefix+"/api/admin/users/{userId}/promote", a.userOnly(a.adminOnly(a.handleUserPromote())))
	router.Handle("PATCH "+prefix+"/api/admin/users/{userId}/demote", a.userOnly(a.adminOnly(a.handleUserDemote())))
	router.Handle("PATCH "+prefix+"/api/admin/users/{userId}/disable", a.userOnly(a.adminOnly(a.handleUserDisable())))
	router.Handle("PATCH "+prefix+"/api/admin/users/{userId}/enable", a.userOnly(a.adminOnly(a.handleUserEnable())))
	router.Handle("PATCH "+prefix+"/api/admin/users/{userId}/password", a.userOnly(a.adminOnly(a.handleAdminUpdateUserPassword())))
	router.Handle("GET "+prefix+"/api/admin/organizations", a.userOnly(a.adminOnly(a.handleGetOrganizations())))
	router.Handle("GET "+prefix+"/api/admin/teams", a.userOnly(a.adminOnly(a.handleGetTeams())))
	router.Handle("GET "+prefix+"/api/admin/apikeys", a.userOnly(a.adminOnly(a.handleGetAPIKeys())))
	router.Handle("GET "+prefix+"/api/admin/search/users/email", a.userOnly(a.adminOnly(a.handleSearchRegisteredUsersByEmail())))
	// alert
	router.Handle("GET "+prefix+"/api/alerts", a.userOnly(a.adminOnly(a.handleGetAlerts())))
	router.Handle("POST "+prefix+"/api/alerts", a.userOnly(a.adminOnly(a.handleAlertCreate())))
	router.Handle("PUT "+prefix+"/api/alerts/{alertId}", a.userOnly(a.adminOnly(a.handleAlertUpdate())))
	router.Handle("DELETE "+prefix+"/api/alerts/{alertId}", a.userOnly(a.adminOnly(a.handleAlertDelete())))
	// maintenance
	router.Handle("DELETE "+prefix+"/api/maintenance/clean-guests", a.userOnly(a.adminOnly(a.handleCleanGuests())))
	// poker games(s)
	if a.Config.FeaturePoker {
		router.Handle("POST "+prefix+"/api/users/{userId}/battles", a.userOnly(a.entityUserOnly(a.handlePokerCreate())))
		router.Handle("GET "+prefix+"/api/users/{userId}/battles", a.userOnly(a.entityUserOnly(a.handleGetUserGames())))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/battles", a.userOnly(a.teamUserOnly(a.handleGetTeamPokerGames())))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/battles/{battleId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemovePokerGame()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/battles", a.userOnly(a.teamUserOnly(a.handlePokerCreate())))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/battles", a.userOnly(a.teamUserOnly(a.handleGetTeamPokerGames())))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/teams/{teamId}/battles/{battleId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemovePokerGame()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/teams/{teamId}/users/{userId}/battles", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handlePokerCreate()))))
		router.Handle("GET "+prefix+"/api/teams/{teamId}/battles", a.userOnly(a.teamUserOnly(a.handleGetTeamPokerGames())))
		router.Handle("DELETE "+prefix+"/api/teams/{teamId}/battles/{battleId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemovePokerGame()))))
		router.Handle("POST "+prefix+"/api/teams/{teamId}/users/{userId}/battles", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handlePokerCreate()))))
		router.Handle("DELETE "+prefix+"/api/maintenance/clean-battles", a.userOnly(a.adminOnly(a.handleCleanPokerGames())))
		router.Handle("GET "+prefix+"/api/battles", a.userOnly(a.adminOnly(a.handleGetPokerGames())))
		router.Handle("GET "+prefix+"/api/battles/{battleId}", a.userOnly(a.handleGetPokerGame()))
		router.Handle("DELETE "+prefix+"/api/battles/{battleId}", a.userOnly(a.handlePokerDelete(pokerSvc)))
		router.Handle("POST "+prefix+"/api/battles/{battleId}/plans", a.userOnly(a.handlePokerStoryAdd(pokerSvc)))
		router.Handle("PUT "+prefix+"/api/battles/{battleId}/plans/{planId}", a.userOnly(a.handlePokerStoryUpdate(pokerSvc)))
		router.Handle("DELETE "+prefix+"/api/battles/{battleId}/plans/{planId}", a.userOnly(a.handlePokerStoryDelete(pokerSvc)))
		router.Handle(prefix+"/api/arena/{battleId}", pokerSvc.ServeBattleWs())

		// estimation scales
		// Public estimation scale routes
		router.Handle("GET "+prefix+"/api/estimation-scales/public", a.userOnly(a.handleGetPublicEstimationScales()))
		router.Handle("GET "+prefix+"/api/estimation-scales/public/{scaleId}", a.userOnly(a.handleGetPublicEstimationScale()))

		// Organization-specific estimation scale routes
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.orgUserOnly(a.handleGetOrganizationEstimationScales()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationEstimationScaleCreate()))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationEstimationScaleUpdate()))))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationEstimationScaleDelete()))))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.departmentUserOnly(a.handleGetTeamEstimationScales()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleCreate())))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleUpdate())))))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleDelete())))))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.handleGetTeamEstimationScales()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/teams/{teamId}/estimation-scales", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleCreate())))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/teams/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleUpdate())))))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/teams/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleDelete())))))

		// Team-specific estimation scale routes
		router.Handle("GET "+prefix+"/api/teams/{teamId}/estimation-scales", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.handleGetTeamEstimationScales()))))
		router.Handle("POST "+prefix+"/api/teams/{teamId}/estimation-scales", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleCreate())))))
		router.Handle("PUT "+prefix+"/api/teams/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleUpdate())))))
		router.Handle("DELETE "+prefix+"/api/teams/{teamId}/estimation-scales/{scaleId}", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamEstimationScaleDelete())))))

		// Admin estimation scale routes
		router.Handle("GET "+prefix+"/api/admin/estimation-scales", a.userOnly(a.adminOnly(a.handleGetEstimationScales())))
		router.Handle("POST "+prefix+"/api/admin/estimation-scales", a.userOnly(a.adminOnly(a.handleEstimationScaleCreate())))
		router.Handle("PUT "+prefix+"/api/admin/estimation-scales/{scaleId}", a.userOnly(a.adminOnly(a.handleEstimationScaleUpdate())))
		router.Handle("DELETE "+prefix+"/api/admin/estimation-scales/{scaleId}", a.userOnly(a.adminOnly(a.handleEstimationScaleDelete())))

		// Organization-specific poker settings routes
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/poker-settings", a.userOnly(a.subscribedOrgOnly(a.orgUserOnly(a.handleGetOrganizationPokerSettings()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/poker-settings", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleCreateOrganizationPokerSettings()))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/poker-settings", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationPokerSettingsUpdate()))))

		// Department-specific poker settings routes
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/poker-settings", a.userOnly(a.subscribedOrgOnly(a.departmentUserOnly(a.handleGetDepartmentPokerSettings()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/poker-settings", a.userOnly(a.subscribedOrgOnly(a.departmentAdminOnly(a.handleCreateDepartmentPokerSettings()))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/poker-settings", a.userOnly(a.subscribedOrgOnly(a.departmentAdminOnly(a.handleDepartmentPokerSettingsUpdate()))))

		// Team-specific poker settings routes
		router.Handle("GET "+prefix+"/api/teams/{teamId}/poker-settings", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.handleGetTeamPokerSettings()))))
		router.Handle("POST "+prefix+"/api/teams/{teamId}/poker-settings", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleCreateTeamPokerSettings())))))
		router.Handle("PUT "+prefix+"/api/teams/{teamId}/poker-settings", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamPokerSettingsUpdate())))))

		// Admin poker settings routes
		router.Handle("GET "+prefix+"/api/admin/poker-settings/{id}", a.userOnly(a.adminOnly(a.handleGetPokerSettingsByID())))
		router.Handle("DELETE "+prefix+"/api/admin/poker-settings/{id}", a.userOnly(a.adminOnly(a.handleDeletePokerSettings())))
	}
	// retro(s)
	if a.Config.FeatureRetro {
		router.Handle("POST "+prefix+"/api/users/{userId}/retros", a.userOnly(a.entityUserOnly(a.handleRetroCreate())))
		router.Handle("GET "+prefix+"/api/users/{userId}/retros", a.userOnly(a.entityUserOnly(a.handleRetrosGetByUser())))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retros", a.userOnly(a.teamUserOnly(a.handleGetTeamRetros())))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retros/{retroId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveRetro()))))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retro-actions", a.userOnly(a.teamUserOnly(a.handleGetTeamRetroActions())))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/retros", a.userOnly(a.teamUserOnly(a.handleRetroCreate())))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/retros", a.userOnly(a.teamUserOnly(a.handleGetTeamRetros())))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/retro-actions", a.userOnly(a.teamUserOnly(a.handleGetTeamRetroActions())))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/teams/{teamId}/retros/{retroId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveRetro()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/teams/{teamId}/users/{userId}/retros", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handleRetroCreate()))))
		router.Handle("GET "+prefix+"/api/teams/{teamId}/retros", a.userOnly(a.teamUserOnly(a.handleGetTeamRetros())))
		router.Handle("DELETE "+prefix+"/api/teams/{teamId}/retros/{retroId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveRetro()))))
		router.Handle("GET "+prefix+"/api/teams/{teamId}/retro-actions", a.userOnly(a.teamUserOnly(a.handleGetTeamRetroActions())))
		router.Handle("POST "+prefix+"/api/teams/{teamId}/users/{userId}/retros", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handleRetroCreate()))))
		router.Handle("DELETE "+prefix+"/api/maintenance/clean-retros", a.userOnly(a.adminOnly(a.handleCleanRetros())))
		router.Handle("GET "+prefix+"/api/retros", a.userOnly(a.adminOnly(a.handleGetRetros())))
		router.Handle("GET "+prefix+"/api/retros/{retroId}", a.userOnly(a.handleRetroGet()))
		router.Handle("DELETE "+prefix+"/api/retros/{retroId}", a.userOnly(a.handleRetroDelete(retroSvc)))
		router.Handle("PUT "+prefix+"/api/retros/{retroId}/actions/{actionId}", a.userOnly(a.handleRetroActionUpdate(retroSvc)))
		router.Handle("DELETE "+prefix+"/api/retros/{retroId}/actions/{actionId}", a.userOnly(a.handleRetroActionDelete(retroSvc)))
		router.Handle("POST "+prefix+"/api/retros/{retroId}/actions/{actionId}/assignees", a.userOnly(a.handleRetroActionAssigneeAdd(retroSvc)))
		router.Handle("DELETE "+prefix+"/api/retros/{retroId}/actions/{actionId}/assignees", a.userOnly(a.handleRetroActionAssigneeRemove(retroSvc)))
		router.Handle("POST "+prefix+"/api/retros/{retroId}/actions/{actionId}/comments", a.userOnly(a.handleRetroActionCommentAdd()))
		router.Handle("PUT "+prefix+"/api/retros/{retroId}/actions/{actionId}/comments/{commentId}", a.userOnly(a.handleRetroActionCommentEdit()))
		router.Handle("DELETE "+prefix+"/api/retros/{retroId}/actions/{actionId}/comments/{commentId}", a.userOnly(a.handleRetroActionCommentDelete()))

		// Retro Templates
		router.Handle("GET "+prefix+"/api/retro-templates/public", a.userOnly(a.handleGetPublicRetroTemplates()))
		// Organization templates
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.orgUserOnly(a.handleGetOrganizationRetroTemplates()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationRetroTemplateCreate()))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationRetroTemplateUpdate()))))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationRetroTemplateDelete()))))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.departmentUserOnly(a.handleGetTeamRetroTemplates()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateCreate())))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateUpdate())))))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateDelete())))))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.handleGetTeamRetroTemplates()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/teams/{teamId}/retro-templates", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateCreate())))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/teams/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateUpdate())))))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/teams/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedOrgOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateDelete())))))

		// Team templates
		router.Handle("GET "+prefix+"/api/teams/{teamId}/retro-templates", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.handleGetTeamRetroTemplates()))))
		router.Handle("POST "+prefix+"/api/teams/{teamId}/retro-templates", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateCreate())))))
		router.Handle("PUT "+prefix+"/api/teams/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateUpdate())))))
		router.Handle("DELETE "+prefix+"/api/teams/{teamId}/retro-templates/{templateId}", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroTemplateDelete())))))

		// General template operations
		router.Handle("GET "+prefix+"/api/admin/retro-templates", a.userOnly(a.adminOnly(a.handleGetRetroTemplates())))
		router.Handle("GET "+prefix+"/api/admin/retro-templates/{templateId}", a.userOnly(a.adminOnly(a.handleGetRetroTemplateByID())))
		router.Handle("POST "+prefix+"/api/admin/retro-templates", a.userOnly(a.adminOnly(a.handleRetroTemplateCreate())))
		router.Handle("PUT "+prefix+"/api/admin/retro-templates/{templateId}", a.userOnly(a.adminOnly(a.handleRetroTemplateUpdate())))
		router.Handle("DELETE "+prefix+"/api/admin/retro-templates/{templateId}", a.userOnly(a.adminOnly(a.handleRetroTemplateDelete())))

		// Organization-specific retro settings routes
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/retro-settings", a.userOnly(a.subscribedOrgOnly(a.orgUserOnly(a.handleGetOrganizationRetroSettings()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/retro-settings", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleCreateOrganizationRetroSettings()))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/retro-settings", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationRetroSettingsUpdate()))))

		// Department-specific retro settings routes
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/retro-settings", a.userOnly(a.subscribedOrgOnly(a.departmentUserOnly(a.handleGetDepartmentRetroSettings()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/retro-settings", a.userOnly(a.subscribedOrgOnly(a.departmentAdminOnly(a.handleCreateDepartmentRetroSettings()))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/retro-settings", a.userOnly(a.subscribedOrgOnly(a.departmentAdminOnly(a.handleDepartmentRetroSettingsUpdate()))))

		// Team-specific retro settings routes
		router.Handle("GET "+prefix+"/api/teams/{teamId}/retro-settings", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.handleGetTeamRetroSettings()))))
		router.Handle("POST "+prefix+"/api/teams/{teamId}/retro-settings", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleCreateTeamRetroSettings())))))
		router.Handle("PUT "+prefix+"/api/teams/{teamId}/retro-settings", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRetroSettingsUpdate())))))

		// Admin retro settings routes
		router.Handle("GET "+prefix+"/api/admin/retro-settings/{id}", a.userOnly(a.adminOnly(a.handleGetRetroSettingsByID())))
		router.Handle("DELETE "+prefix+"/api/admin/retro-settings/{id}", a.userOnly(a.adminOnly(a.handleDeleteRetroSettings())))

		// Retro websocket
		router.Handle(prefix+"/api/retro/{retroId}", retroSvc.ServeWs())
	}
	// storyboard(s)
	if a.Config.FeatureStoryboard {
		router.Handle("POST "+prefix+"/api/users/{userId}/storyboards", a.userOnly(a.entityUserOnly(a.handleStoryboardCreate())))
		router.Handle("GET "+prefix+"/api/users/{userId}/storyboards", a.userOnly(a.entityUserOnly(a.handleGetUserStoryboards())))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/storyboards", a.userOnly(a.teamUserOnly(a.handleGetTeamStoryboards())))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/storyboards/{storyboardId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveStoryboard()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/storyboards", a.userOnly(a.teamUserOnly(a.handleStoryboardCreate())))
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/teams/{teamId}/storyboards", a.userOnly(a.teamUserOnly(a.handleGetTeamStoryboards())))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/teams/{teamId}/storyboards/{storyboardId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveStoryboard()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/teams/{teamId}/users/{userId}/storyboards", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handleStoryboardCreate()))))
		router.Handle("GET "+prefix+"/api/teams/{teamId}/storyboards", a.userOnly(a.teamUserOnly(a.handleGetTeamStoryboards())))
		router.Handle("DELETE "+prefix+"/api/teams/{teamId}/storyboards/{storyboardId}", a.userOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamRemoveStoryboard()))))
		router.Handle("POST "+prefix+"/api/teams/{teamId}/users/{userId}/storyboards", a.userOnly(a.teamUserOnly(a.entityUserOnly(a.handleStoryboardCreate()))))
		router.Handle("DELETE "+prefix+"/api/maintenance/clean-storyboards", a.userOnly(a.adminOnly(a.handleCleanStoryboards())))
		router.Handle("GET "+prefix+"/api/storyboards", a.userOnly(a.adminOnly(a.handleGetStoryboards())))
		router.Handle("GET "+prefix+"/api/storyboards/{storyboardId}", a.userOnly(a.handleStoryboardGet()))
		// Storyboard operations
		router.Handle("DELETE "+prefix+"/api/storyboards/{storyboardId}", a.userOnly(a.handleStoryboardDelete(storyboardSvc)))
		// Storyboard goal operations
		router.Handle("POST "+prefix+"/api/storyboards/{storyboardId}/goals", a.userOnly(a.handleStoryboardGoalAdd(storyboardSvc)))
		router.Handle("PUT "+prefix+"/api/storyboards/{storyboardId}/goals/{goalId}", a.userOnly(a.handleStoryboardGoalUpdate(storyboardSvc)))
		router.Handle("DELETE "+prefix+"/api/storyboards/{storyboardId}/goals/{goalId}", a.userOnly(a.handleStoryboardGoalDelete(storyboardSvc)))
		// Storyboard column operations
		router.Handle("POST "+prefix+"/api/storyboards/{storyboardId}/columns", a.userOnly(a.handleStoryboardColumnAdd(storyboardSvc)))
		router.Handle("PUT "+prefix+"/api/storyboards/{storyboardId}/columns/{columnId}", a.userOnly(a.handleStoryboardColumnUpdate(storyboardSvc)))
		router.Handle("DELETE "+prefix+"/api/storyboards/{storyboardId}/columns/{columnId}", a.userOnly(a.handleStoryboardColumnDelete(storyboardSvc)))
		// Storyboard story operations
		router.Handle("POST "+prefix+"/api/storyboards/{storyboardId}/stories", a.userOnly(a.handleStoryboardStoryAdd(storyboardSvc)))
		router.Handle("PUT "+prefix+"/api/storyboards/{storyboardId}/stories/{storyId}/name", a.userOnly(a.handleStoryboardStoryNameUpdate(storyboardSvc)))
		router.Handle("PUT "+prefix+"/api/storyboards/{storyboardId}/stories/{storyId}/content", a.userOnly(a.handleStoryboardStoryContentUpdate(storyboardSvc)))
		router.Handle("PUT "+prefix+"/api/storyboards/{storyboardId}/stories/{storyId}/color", a.userOnly(a.handleStoryboardStoryColorUpdate(storyboardSvc)))
		router.Handle("PUT "+prefix+"/api/storyboards/{storyboardId}/stories/{storyId}/points", a.userOnly(a.handleStoryboardStoryPointsUpdate(storyboardSvc)))
		router.Handle("PUT "+prefix+"/api/storyboards/{storyboardId}/stories/{storyId}/closed", a.userOnly(a.handleStoryboardStoryClosedUpdate(storyboardSvc)))
		router.Handle("PUT "+prefix+"/api/storyboards/{storyboardId}/stories/{storyId}/link", a.userOnly(a.handleStoryboardStoryLinkUpdate(storyboardSvc)))
		router.Handle("PUT "+prefix+"/api/storyboards/{storyboardId}/stories/{storyId}/move", a.userOnly(a.handleStoryboardStoryMove(storyboardSvc)))
		router.Handle("DELETE "+prefix+"/api/storyboards/{storyboardId}/stories/{storyId}", a.userOnly(a.handleStoryboardStoryDelete(storyboardSvc)))
		// Storyboard websocket
		router.Handle(""+prefix+"/api/storyboard/{storyboardId}", storyboardSvc.ServeWs())
	}

	if a.Config.FeatureProject {
		// Projects
		// Organization projects
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/projects", a.userOnly(a.subscribedOrgOnly(a.orgUserOnly(a.handleGetOrganizationProjects()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/projects", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationProjectCreate()))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/projects/{projectId}", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationProjectUpdate()))))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/projects/{projectId}", a.userOnly(a.subscribedOrgOnly(a.orgAdminOnly(a.handleOrganizationProjectDelete()))))

		// Department projects (nested under organization)
		router.Handle("GET "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/projects", a.userOnly(a.subscribedOrgOnly(a.departmentUserOnly(a.handleGetDepartmentProjects()))))
		router.Handle("POST "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/projects", a.userOnly(a.subscribedOrgOnly(a.departmentAdminOnly(a.handleDepartmentProjectCreate()))))
		router.Handle("PUT "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/projects/{projectId}", a.userOnly(a.subscribedOrgOnly(a.departmentAdminOnly(a.handleDepartmentProjectUpdate()))))
		router.Handle("DELETE "+prefix+"/api/organizations/{orgId}/departments/{departmentId}/projects/{projectId}", a.userOnly(a.subscribedOrgOnly(a.departmentAdminOnly(a.handleDepartmentProjectDelete()))))

		// Team projects (direct team routes only)
		router.Handle("GET "+prefix+"/api/teams/{teamId}/projects", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.handleGetTeamProjects()))))
		router.Handle("POST "+prefix+"/api/teams/{teamId}/projects", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamProjectCreate())))))
		router.Handle("PUT "+prefix+"/api/teams/{teamId}/projects/{projectId}", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamProjectUpdate())))))
		router.Handle("DELETE "+prefix+"/api/teams/{teamId}/projects/{projectId}", a.userOnly(a.subscribedTeamOnly(a.teamUserOnly(a.teamAdminOnly(a.handleTeamProjectDelete())))))

		// Project operations
		router.Handle("GET "+prefix+"/api/projects/{projectId}", a.userOnly(a.subscribedProjectOnly(a.projectUserOnly(a.handleGetProjectByID()))))

		// Project storyboards
		if a.Config.FeatureStoryboard {
			router.Handle("GET "+prefix+"/api/projects/{projectId}/storyboards", a.userOnly(a.subscribedProjectOnly(a.projectUserOnly(a.handleGetProjectStoryboards()))))
		}

		// Project retros
		if a.Config.FeatureRetro {
			router.Handle("GET "+prefix+"/api/projects/{projectId}/retros", a.userOnly(a.subscribedProjectOnly(a.projectUserOnly(a.handleGetProjectRetros()))))
		}

		// Project poker games
		if a.Config.FeaturePoker {
			router.Handle("GET "+prefix+"/api/projects/{projectId}/poker", a.userOnly(a.subscribedProjectOnly(a.projectUserOnly(a.handleGetProjectPokerGames()))))
		}

		// General project operations (admin)
		router.Handle("GET "+prefix+"/api/admin/projects", a.userOnly(a.adminOnly(a.handleAdminGetProjects())))
		router.Handle("GET "+prefix+"/api/admin/projects/{projectId}", a.userOnly(a.adminOnly(a.handleAdminGetProjectByID())))
		router.Handle("POST "+prefix+"/api/admin/projects", a.userOnly(a.adminOnly(a.handleAdminProjectCreate())))
		router.Handle("PUT "+prefix+"/api/admin/projects/{projectId}", a.userOnly(a.adminOnly(a.handleAdminProjectUpdate())))
		router.Handle("DELETE "+prefix+"/api/admin/projects/{projectId}", a.userOnly(a.adminOnly(a.handleAdminProjectDelete())))
	}

	// user avatar generation
	if a.Config.AvatarService == "goadorable" || a.Config.AvatarService == "govatar" {
		router.Handle("GET "+prefix+"/avatar/{width}/{id}/{avatar}", a.handleUserAvatar())
		router.Handle("GET "+prefix+"/avatar/{width}/{id}", a.handleUserAvatar())
	}

	if a.Config.SubscriptionsEnabled {
		router.Handle("GET "+prefix+"/api/subscriptions/{subscriptionId}", a.userOnly(a.adminOnly(a.handleSubscriptionGetByID())))
		router.Handle("PUT "+prefix+"/api/subscriptions/{subscriptionId}", a.userOnly(a.adminOnly(a.handleSubscriptionUpdate())))
		router.Handle("DELETE "+prefix+"/api/subscriptions/{subscriptionId}", a.userOnly(a.adminOnly(a.handleSubscriptionDelete())))
		router.Handle("GET "+prefix+"/api/subscriptions", a.userOnly(a.adminOnly(a.handleGetSubscriptions())))
		router.Handle("POST "+prefix+"/api/subscriptions", a.userOnly(a.adminOnly(a.handleSubscriptionCreate())))
		router.Handle("POST "+prefix+"/webhooks/subscriptions", a.SubscriptionSvc.HandleWebhook())
	}

	a.registerOauthProviderEndpoints(router, prefix, authProviderConfigs)

	// static assets
	router.Handle(prefix+"/assets/", http.StripPrefix(a.Config.PathPrefix, staticHandler))
	router.Handle(prefix+"/img/", http.StripPrefix(a.Config.PathPrefix, staticHandler))

	// health check for load balancers, k8s, etc...
	router.Handle(prefix+"/healthz", a.handleHealthCheck())
	router.Handle("GET "+prefix+"/api/{$}", a.handleApiIndex())

	router.Handle(prefix+"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't return the index for API or other path prefixes
		if strings.HasPrefix(r.URL.Path, prefix+"/api") ||
			strings.HasPrefix(r.URL.Path, prefix+"/static") ||
			strings.HasPrefix(r.URL.Path, prefix+"/img") ||
			strings.HasPrefix(r.URL.Path, prefix+"/swagger") ||
			strings.HasPrefix(r.URL.Path, prefix+"/avatar") ||
			r.URL.Path == prefix+"/healthz" {
			http.NotFound(w, r)
			return
		}

		a.handleIndex(FSS, a.UIConfig)(w, r)
	}))

	var handler http.Handler = router
	handler = a.panicRecovery(handler)
	handler = otelhttp.NewHandler(handler, "thunderdome")
	handler = cspMiddleware(handler, secureMiddleware, prefix, a.Config.ExternalAPIEnabled)

	a.Handler = handler

	return a
}

func (s *Service) registerOauthProviderEndpoints(router *http.ServeMux, prefix string, providers []thunderdome.AuthProviderConfig) {
	ctx := context.Background()
	var redirectBaseURL string
	var port string

	// redirect with port for localhost
	if s.Config.AppDomain == "localhost" {
		port = fmt.Sprintf(":%s", s.Config.Port)
	}

	if s.Config.SecureProtocol {
		redirectBaseURL = fmt.Sprintf("https://%s%s%s", s.Config.AppDomain, port, s.Config.PathPrefix)
	} else {
		redirectBaseURL = fmt.Sprintf("http://%s%s%s", s.Config.AppDomain, port, s.Config.PathPrefix)
	}

	for _, c := range providers {
		providerNameUrlPath := strings.ToLower(c.ProviderName)
		oauthPathPrefix := "/oauth/"
		oauthLoginPathPrefix, _ := url.JoinPath(oauthPathPrefix, providerNameUrlPath, "/login")
		oauthCallbackPathPrefix, _ := url.JoinPath(oauthPathPrefix, providerNameUrlPath, "/callback")
		callbackRedirectURL, _ := url.JoinPath(redirectBaseURL, oauthCallbackPathPrefix)
		authProvider, err := oauth.New(oauth.Config{
			AuthProviderConfig:  c,
			CallbackRedirectURL: callbackRedirectURL,
			UIRedirectURL:       fmt.Sprintf("%s/", s.Config.PathPrefix),
			InternalOnlyOidc:    s.Config.OIDCAuth.Enabled,
		}, s.Cookie, s.Logger, s.AuthDataSvc, s.SubscriptionDataSvc, ctx)
		if err != nil {
			panic(err)
		}
		router.Handle("GET "+prefix+oauthLoginPathPrefix, authProvider.HandleOAuth2Redirect())
		router.Handle("GET "+prefix+oauthCallbackPathPrefix, authProvider.HandleOAuth2Callback())
	}
}

func (s *Service) ListenAndServe() error {
	srv := &http.Server{
		Handler:           s.Handler,
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
func (s *Service) handleIndex(filesystem fs.FS, uiConfig thunderdome.UIConfig) http.HandlerFunc {
	tmpl := s.getIndexTemplate(filesystem)

	ActiveAlerts = s.AlertDataSvc.GetActiveAlerts(context.Background()) // prime the active alerts cache

	return func(w http.ResponseWriter, r *http.Request) {
		uiConfig.ActiveAlerts = ActiveAlerts // get the latest alerts from memory
		nonce := secure.CSPNonce(r.Context())

		if s.Config.EmbedUseOS {
			tmpl = s.getIndexTemplate(filesystem)
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

func (s *Service) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
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

func cspMiddleware(next http.Handler, secureMiddleware *secure.Secure, prefix string, externalAPIEnabled bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if externalAPIEnabled && strings.HasPrefix(r.URL.Path, prefix+"/swagger") {
			// Skip CSP for swagger docs
			next.ServeHTTP(w, r)
			return
		}

		// For all other requests, apply CSP
		secureMiddleware.Handler(next).ServeHTTP(w, r)
	})
}
