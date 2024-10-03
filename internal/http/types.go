package http

import (
	"context"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/webhook/subscription"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

const (
	contextKeyUserID         contextKey = "userId"
	contextKeyUserType       contextKey = "userType"
	apiKeyHeaderName         string     = "X-API-Key"
	contextKeyUserTeamRoles  contextKey = "userTeamRoles"
	contextKeyOrgRole        contextKey = "orgRole"
	contextKeyDepartmentRole contextKey = "departmentRole"
)

var validate *validator.Validate

type WebsocketConfig struct {
	// Time allowed to write a message to the peer.
	WriteWaitSec int

	// Time allowed to read the next pong message from the peer.
	PongWaitSec int

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriodSec int

	// Websocket subdomain (allow websockets to be routed via a subdomain)
	WebsocketSubdomain string
}

type AuthProvider struct {
	Enabled bool
	thunderdome.AuthProviderConfig
}

// Config contains configuration values used by the APIs
type Config struct {
	Port                  string
	HttpWriteTimeout      int
	HttpReadTimeout       int
	HttpIdleTimeout       int
	HttpReadHeaderTimeout int
	// the domain of the application for cookie securing
	AppDomain string
	// PathPrefix allows the application to be run on a shared domain
	PathPrefix string
	// SecureProtocol whether the application is accessed through HTTPS
	SecureProtocol bool
	// Whether the external API is enabled
	ExternalAPIEnabled bool
	// Whether the external API requires user verified email
	ExternalAPIVerifyRequired bool
	// Number of API keys a user can create
	UserAPIKeyLimit int
	// Whether LDAP authentication is enabled for self-hosted
	LdapEnabled bool
	// Whether header authentication is enabled for self-hosted
	HeaderAuthEnabled bool
	// Feature flag for Poker Planning
	FeaturePoker bool
	// Feature flag for Retrospectives
	FeatureRetro bool
	// Feature flag for Storyboards
	FeatureStoryboard bool
	// Whether Organizations (and Departments) feature is enabled
	OrganizationsEnabled bool
	// Which avatar service is utilized
	AvatarService string
	// ID of default template to select for Retro creation
	RetroDefaultTemplateID string
	// Whether to use the OS filesystem or embedded
	EmbedUseOS                bool
	CleanupBattlesDaysOld     int
	CleanupRetrosDaysOld      int
	CleanupStoryboardsDaysOld int
	CleanupGuestsDaysOld      int
	RequireTeams              bool
	AuthLdapUrl               string
	AuthLdapUseTls            bool
	AuthLdapBindname          string
	AuthLdapBindpass          string
	AuthLdapBasedn            string
	AuthLdapFilter            string
	AuthLdapMailAttr          string
	AuthLdapCnAttr            string
	AuthHeaderUsernameHeader  string
	AuthHeaderEmailHeader     string
	AllowGuests               bool
	AllowRegistration         bool
	ShowActiveCountries       bool
	SubscriptionsEnabled      bool

	GoogleAuth AuthProvider
	WebsocketConfig
}

type Service struct {
	Config               *Config
	Cookie               CookieManager
	UIConfig             thunderdome.UIConfig
	Router               *mux.Router
	Email                thunderdome.EmailService
	Logger               *otelzap.Logger
	UserDataSvc          UserDataSvc
	ApiKeyDataSvc        APIKeyDataSvc
	AlertDataSvc         AlertDataSvc
	AuthDataSvc          AuthDataSvc
	PokerDataSvc         thunderdome.PokerDataSvc
	CheckinDataSvc       CheckinDataSvc
	RetroDataSvc         thunderdome.RetroDataSvc
	StoryboardDataSvc    thunderdome.StoryboardDataSvc
	TeamDataSvc          TeamDataSvc
	OrganizationDataSvc  OrganizationDataSvc
	AdminDataSvc         AdminDataSvc
	JiraDataSvc          JiraDataSvc
	SubscriptionDataSvc  SubscriptionDataSvc
	RetroTemplateDataSvc thunderdome.RetroTemplateDataSvc
	SubscriptionSvc      *subscription.Service
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

type CookieManager interface {
	CreateUserCookie(w http.ResponseWriter, userID string) error
	CreateSessionCookie(w http.ResponseWriter, sessionID string) error
	CreateUserUICookie(w http.ResponseWriter, userUiCookie thunderdome.UserUICookie) error
	ClearUserCookies(w http.ResponseWriter)
	ValidateUserCookie(w http.ResponseWriter, r *http.Request) (string, error)
	ValidateSessionCookie(w http.ResponseWriter, r *http.Request) (string, error)
	CreateCookie(w http.ResponseWriter, cookieName string, value string, maxAge int) error
	GetCookie(w http.ResponseWriter, r *http.Request, cookieName string) (string, error)
	DeleteCookie(w http.ResponseWriter, cookieName string)
	CreateAuthStateCookie(w http.ResponseWriter, state string) error
	ValidateAuthStateCookie(w http.ResponseWriter, r *http.Request, state string) error
	DeleteAuthStateCookie(w http.ResponseWriter) error
}

type AdminDataSvc interface {
	GetAppStats(ctx context.Context) (*thunderdome.ApplicationStats, error)
}

type AlertDataSvc interface {
	GetActiveAlerts(ctx context.Context) []interface{}
	AlertsList(ctx context.Context, limit int, offset int) ([]*thunderdome.Alert, int, error)
	AlertsCreate(ctx context.Context, name string, alertType string, content string, active bool, allowDismiss bool, registeredOnly bool) error
	AlertsUpdate(ctx context.Context, alertID string, name string, alertType string, content string, active bool, allowDismiss bool, registeredOnly bool) error
	AlertDelete(ctx context.Context, alertID string) error
}

type APIKeyDataSvc interface {
	GenerateApiKey(ctx context.Context, userID string, keyName string) (*thunderdome.APIKey, error)
	GetUserApiKeys(ctx context.Context, userID string) ([]*thunderdome.APIKey, error)
	GetApiKeyUser(ctx context.Context, apiKey string) (*thunderdome.User, error)
	GetAPIKeys(ctx context.Context, limit int, offset int) []*thunderdome.UserAPIKey
	UpdateUserApiKey(ctx context.Context, userID string, keyID string, active bool) ([]*thunderdome.APIKey, error)
	DeleteUserApiKey(ctx context.Context, userID string, keyID string) ([]*thunderdome.APIKey, error)
}

type AuthDataSvc interface {
	AuthUser(ctx context.Context, email string, password string) (*thunderdome.User, *thunderdome.Credential, string, error)
	OauthCreateNonce(ctx context.Context) (string, error)
	OauthValidateNonce(ctx context.Context, nonceId string) error
	OauthAuthUser(ctx context.Context, provider string, sub string, email string, emailVerified bool, name string, pictureUrl string) (*thunderdome.User, string, error)
	UserResetRequest(ctx context.Context, email string) (resetID string, userName string, resetErr error)
	UserResetPassword(ctx context.Context, resetID string, password string) (userName string, email string, resetErr error)
	UserUpdatePassword(ctx context.Context, userID string, password string) (name string, email string, resetErr error)
	UserVerifyRequest(ctx context.Context, userId string) (*thunderdome.User, string, error)
	VerifyUserAccount(ctx context.Context, verifyID string) error
	MFASetupGenerate(email string) (string, string, error)
	MFASetupValidate(ctx context.Context, userID string, secret string, passcode string) error
	MFARemove(ctx context.Context, userID string) error
	MFATokenValidate(ctx context.Context, sessionId string, passcode string) error
	CreateSession(ctx context.Context, userId string, enabled bool) (string, error)
	EnableSession(ctx context.Context, sessionId string) error
	GetSessionUser(ctx context.Context, sessionId string) (*thunderdome.User, error)
	DeleteSession(ctx context.Context, sessionId string) error
}

type CheckinDataSvc interface {
	CheckinList(ctx context.Context, teamID string, Date string, TimeZone string) ([]*thunderdome.TeamCheckin, error)
	CheckinCreate(ctx context.Context, teamID string, userId string, yesterday string, today string, blockers string, discuss string, goalsMet bool) error
	CheckinUpdate(ctx context.Context, checkinId string, yesterday string, today string, blockers string, discuss string, goalsMet bool) error
	CheckinDelete(ctx context.Context, checkinId string) error
	CheckinComment(ctx context.Context, teamID string, checkinId string, userId string, comment string) error
	CheckinCommentEdit(ctx context.Context, teamID string, userId string, commentId string, comment string) error
	CheckinCommentDelete(ctx context.Context, commentId string) error
	CheckinLastByUser(ctx context.Context, teamID string, userId string) (*thunderdome.TeamCheckin, error)
}

type JiraDataSvc interface {
	FindInstancesByUserId(ctx context.Context, userId string) ([]thunderdome.JiraInstance, error)
	GetInstanceById(ctx context.Context, instanceId string) (thunderdome.JiraInstance, error)
	CreateInstance(ctx context.Context, userId string, host string, clientMail string, accessToken string) (thunderdome.JiraInstance, error)
	UpdateInstance(ctx context.Context, instanceId string, host string, clientMail string, accessToken string) (thunderdome.JiraInstance, error)
	DeleteInstance(ctx context.Context, instanceId string) error
}

type OrganizationDataSvc interface {
	OrganizationGet(ctx context.Context, orgID string) (*thunderdome.Organization, error)
	OrganizationUserRole(ctx context.Context, userID string, orgID string) (string, error)
	OrganizationListByUser(ctx context.Context, userID string, limit int, offset int) []*thunderdome.UserOrganization
	OrganizationCreate(ctx context.Context, userID string, orgName string) (*thunderdome.Organization, error)
	OrganizationUpdate(ctx context.Context, orgID string, orgName string) (*thunderdome.Organization, error)
	OrganizationUserList(ctx context.Context, orgID string, limit int, offset int) []*thunderdome.OrganizationUser
	OrganizationAddUser(ctx context.Context, orgID string, userID string, Role string) (string, error)
	OrganizationUpsertUser(ctx context.Context, orgID string, userID string, Role string) (string, error)
	OrganizationUpdateUser(ctx context.Context, orgID string, userID string, Role string) (string, error)
	OrganizationRemoveUser(ctx context.Context, organizationID string, userID string) error
	OrganizationInviteUser(ctx context.Context, orgID string, Email string, Role string) (string, error)
	OrganizationUserGetInviteByID(ctx context.Context, inviteID string) (thunderdome.OrganizationUserInvite, error)
	OrganizationDeleteUserInvite(ctx context.Context, inviteID string) error
	OrganizationGetUserInvites(ctx context.Context, orgID string) ([]thunderdome.OrganizationUserInvite, error)
	OrganizationTeamList(ctx context.Context, orgID string, limit int, offset int) []*thunderdome.Team
	OrganizationTeamCreate(ctx context.Context, orgID string, teamName string) (*thunderdome.Team, error)
	OrganizationTeamUserRole(ctx context.Context, userID string, orgID string, teamID string) (string, string, error)
	OrganizationDelete(ctx context.Context, orgID string) error
	OrganizationList(ctx context.Context, limit int, offset int) []*thunderdome.Organization
	OrganizationIsSubscribed(ctx context.Context, orgID string) (bool, error)
	GetOrganizationMetrics(ctx context.Context, organizationID string) (*thunderdome.OrganizationMetrics, error)

	DepartmentUserRole(ctx context.Context, userID string, orgID string, departmentID string) (string, string, error)
	DepartmentGet(ctx context.Context, departmentID string) (*thunderdome.Department, error)
	OrganizationDepartmentList(ctx context.Context, orgID string, limit int, offset int) []*thunderdome.Department
	DepartmentCreate(ctx context.Context, orgID string, orgName string) (*thunderdome.Department, error)
	DepartmentUpdate(ctx context.Context, deptID string, deptName string) (*thunderdome.Department, error)
	DepartmentTeamList(ctx context.Context, departmentID string, limit int, offset int) []*thunderdome.Team
	DepartmentTeamCreate(ctx context.Context, departmentID string, teamName string) (*thunderdome.Team, error)
	DepartmentUserList(ctx context.Context, departmentID string, limit int, offset int) []*thunderdome.DepartmentUser
	DepartmentAddUser(ctx context.Context, departmentID string, userID string, Role string) (string, error)
	DepartmentUpsertUser(ctx context.Context, departmentID string, userID string, Role string) (string, error)
	DepartmentUpdateUser(ctx context.Context, departmentID string, userID string, Role string) (string, error)
	DepartmentRemoveUser(ctx context.Context, departmentID string, userID string) error
	DepartmentTeamUserRole(ctx context.Context, userID string, orgID string, departmentID string, teamID string) (string, string, string, error)
	DepartmentDelete(ctx context.Context, departmentID string) error
	DepartmentInviteUser(ctx context.Context, deptID string, email string, role string) (string, error)
	DepartmentUserGetInviteByID(ctx context.Context, inviteID string) (thunderdome.DepartmentUserInvite, error)
	DepartmentDeleteUserInvite(ctx context.Context, inviteID string) error
	DepartmentGetUserInvites(ctx context.Context, deptID string) ([]thunderdome.DepartmentUserInvite, error)
}

type TeamDataSvc interface {
	TeamUserRole(ctx context.Context, userID string, teamID string) (string, error)
	TeamGet(ctx context.Context, teamID string) (*thunderdome.Team, error)
	TeamListByUser(ctx context.Context, userID string, limit int, offset int) []*thunderdome.UserTeam
	TeamListByUserNonOrg(ctx context.Context, userID string, limit int, offset int) []*thunderdome.UserTeam
	TeamCreate(ctx context.Context, userID string, teamName string) (*thunderdome.Team, error)
	TeamUpdate(ctx context.Context, teamID string, teamName string) (*thunderdome.Team, error)
	TeamAddUser(ctx context.Context, teamID string, userID string, role string) (string, error)
	TeamUserList(ctx context.Context, teamID string, limit int, offset int) ([]*thunderdome.TeamUser, int, error)
	TeamUpdateUser(ctx context.Context, teamID string, userID string, role string) (string, error)
	TeamRemoveUser(ctx context.Context, teamID string, userID string) error
	TeamInviteUser(ctx context.Context, teamID string, Email string, role string) (string, error)
	TeamUserGetInviteByID(ctx context.Context, inviteID string) (thunderdome.TeamUserInvite, error)
	TeamDeleteUserInvite(ctx context.Context, inviteID string) error
	TeamGetUserInvites(ctx context.Context, teamId string) ([]thunderdome.TeamUserInvite, error)
	TeamPokerList(ctx context.Context, teamID string, limit int, offset int) []*thunderdome.Poker
	TeamAddPoker(ctx context.Context, teamID string, pokerID string) error
	TeamRemovePoker(ctx context.Context, teamID string, pokerID string) error
	TeamDelete(ctx context.Context, teamID string) error
	TeamRetroList(ctx context.Context, teamID string, limit int, offset int) []*thunderdome.Retro
	TeamAddRetro(ctx context.Context, teamID string, retroID string) error
	TeamRemoveRetro(ctx context.Context, teamID string, retroID string) error
	TeamStoryboardList(ctx context.Context, teamID string, limit int, offset int) []*thunderdome.Storyboard
	TeamAddStoryboard(ctx context.Context, teamID string, storyboardID string) error
	TeamRemoveStoryboard(ctx context.Context, teamID string, storyboardID string) error
	TeamList(ctx context.Context, limit int, offset int) ([]*thunderdome.Team, int)
	TeamIsSubscribed(ctx context.Context, teamID string) (bool, error)
	GetTeamMetrics(ctx context.Context, teamID string) (*thunderdome.TeamMetrics, error)
	TeamUserRoles(ctx context.Context, userID string, teamID string) (*thunderdome.UserTeamRoleInfo, error)
}

type SubscriptionDataSvc interface {
	CheckActiveSubscriber(ctx context.Context, userID string) error
	GetSubscriptionByID(ctx context.Context, subscriptionID string) (thunderdome.Subscription, error)
	GetSubscriptionBySubscriptionID(ctx context.Context, subscriptionID string) (thunderdome.Subscription, error)
	GetActiveSubscriptionsByUserID(ctx context.Context, userID string) ([]thunderdome.Subscription, error)
	CreateSubscription(ctx context.Context, subscription thunderdome.Subscription) (thunderdome.Subscription, error)
	UpdateSubscription(ctx context.Context, subscriptionID string, subscription thunderdome.Subscription) (thunderdome.Subscription, error)
	GetSubscriptions(ctx context.Context, limit int, offset int) ([]thunderdome.Subscription, int, error)
	DeleteSubscription(ctx context.Context, subscriptionID string) error
}

type UserDataSvc interface {
	GetUser(ctx context.Context, userID string) (*thunderdome.User, error)
	GetGuestUser(ctx context.Context, userID string) (*thunderdome.User, error)
	GetUserByEmail(ctx context.Context, email string) (*thunderdome.User, error)
	GetRegisteredUsers(ctx context.Context, limit int, offset int) ([]*thunderdome.User, int, error)
	SearchRegisteredUsersByEmail(ctx context.Context, email string, limit int, offset int) ([]*thunderdome.User, int, error)
	CreateUser(ctx context.Context, userName string, email string, userPassword string) (newUser *thunderdome.User, verifyID string, registerErr error)
	CreateUserGuest(ctx context.Context, userName string) (*thunderdome.User, error)
	CreateUserRegistered(ctx context.Context, userName string, email string, userPassword string, activeuserID string) (newUser *thunderdome.User, verifyID string, registerErr error)
	UpdateUserAccount(ctx context.Context, userID string, userName string, email string, userAvatar string, notificationsEnabled bool, country string, locale string, company string, jobTitle string, theme string) error
	UpdateUserProfile(ctx context.Context, userID string, userName string, userAvatar string, notificationsEnabled bool, country string, locale string, company string, jobTitle string, theme string) error
	UpdateUserProfileLdap(ctx context.Context, userID string, userAvatar string, notificationsEnabled bool, country string, locale string, company string, jobTitle string, theme string) error
	PromoteUser(ctx context.Context, userID string) error
	DemoteUser(ctx context.Context, userID string) error
	DisableUser(ctx context.Context, userID string) error
	EnableUser(ctx context.Context, userID string) error
	DeleteUser(ctx context.Context, userID string) error
	CleanGuests(ctx context.Context, daysOld int) error
	GetActiveCountries(ctx context.Context) ([]string, error)
	GetUserCredential(ctx context.Context, userID string) (*thunderdome.Credential, error)
}
