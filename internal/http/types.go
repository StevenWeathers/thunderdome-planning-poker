package http

import (
	"context"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/webhook/subscription"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"net/http"
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
	UserDataSvc          thunderdome.UserDataSvc
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
	SubscriptionDataSvc  thunderdome.SubscriptionDataSvc
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
	CreateUserCookie(w http.ResponseWriter, UserID string) error
	CreateSessionCookie(w http.ResponseWriter, SessionID string) error
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
	AlertsList(ctx context.Context, Limit int, Offset int) ([]*thunderdome.Alert, int, error)
	AlertsCreate(ctx context.Context, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error
	AlertsUpdate(ctx context.Context, ID string, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error
	AlertDelete(ctx context.Context, AlertID string) error
}

type APIKeyDataSvc interface {
	GenerateApiKey(ctx context.Context, UserID string, KeyName string) (*thunderdome.APIKey, error)
	GetUserApiKeys(ctx context.Context, UserID string) ([]*thunderdome.APIKey, error)
	GetApiKeyUser(ctx context.Context, APK string) (*thunderdome.User, error)
	GetAPIKeys(ctx context.Context, Limit int, Offset int) []*thunderdome.UserAPIKey
	UpdateUserApiKey(ctx context.Context, UserID string, KeyID string, Active bool) ([]*thunderdome.APIKey, error)
	DeleteUserApiKey(ctx context.Context, UserID string, KeyID string) ([]*thunderdome.APIKey, error)
}

type AuthDataSvc interface {
	AuthUser(ctx context.Context, UserEmail string, UserPassword string) (*thunderdome.User, *thunderdome.Credential, string, error)
	OauthCreateNonce(ctx context.Context) (string, error)
	OauthValidateNonce(ctx context.Context, nonceId string) error
	OauthAuthUser(ctx context.Context, provider string, sub string, email string, emailVerified bool, name string, pictureUrl string) (*thunderdome.User, string, error)
	UserResetRequest(ctx context.Context, UserEmail string) (resetID string, UserName string, resetErr error)
	UserResetPassword(ctx context.Context, ResetID string, UserPassword string) (UserName string, UserEmail string, resetErr error)
	UserUpdatePassword(ctx context.Context, UserID string, UserPassword string) (Name string, Email string, resetErr error)
	UserVerifyRequest(ctx context.Context, UserId string) (*thunderdome.User, string, error)
	VerifyUserAccount(ctx context.Context, VerifyID string) error
	MFASetupGenerate(email string) (string, string, error)
	MFASetupValidate(ctx context.Context, UserID string, secret string, passcode string) error
	MFARemove(ctx context.Context, UserID string) error
	MFATokenValidate(ctx context.Context, SessionId string, passcode string) error
	CreateSession(ctx context.Context, UserId string, enabled bool) (string, error)
	EnableSession(ctx context.Context, SessionId string) error
	GetSessionUser(ctx context.Context, SessionId string) (*thunderdome.User, error)
	DeleteSession(ctx context.Context, SessionId string) error
}

type CheckinDataSvc interface {
	CheckinList(ctx context.Context, TeamId string, Date string, TimeZone string) ([]*thunderdome.TeamCheckin, error)
	CheckinCreate(ctx context.Context, TeamId string, UserId string, Yesterday string, Today string, Blockers string, Discuss string, GoalsMet bool) error
	CheckinUpdate(ctx context.Context, CheckinId string, Yesterday string, Today string, Blockers string, Discuss string, GoalsMet bool) error
	CheckinDelete(ctx context.Context, CheckinId string) error
	CheckinComment(ctx context.Context, TeamId string, CheckinId string, UserId string, Comment string) error
	CheckinCommentEdit(ctx context.Context, TeamId string, UserId string, CommentId string, Comment string) error
	CheckinCommentDelete(ctx context.Context, CommentId string) error
	CheckinLastByUser(ctx context.Context, TeamId string, UserId string) (*thunderdome.TeamCheckin, error)
}

type JiraDataSvc interface {
	FindInstancesByUserId(ctx context.Context, userId string) ([]thunderdome.JiraInstance, error)
	GetInstanceById(ctx context.Context, instanceId string) (thunderdome.JiraInstance, error)
	CreateInstance(ctx context.Context, userId string, host string, clientMail string, accessToken string) (thunderdome.JiraInstance, error)
	UpdateInstance(ctx context.Context, instanceId string, host string, clientMail string, accessToken string) (thunderdome.JiraInstance, error)
	DeleteInstance(ctx context.Context, instanceId string) error
}

type OrganizationDataSvc interface {
	OrganizationGet(ctx context.Context, OrgID string) (*thunderdome.Organization, error)
	OrganizationUserRole(ctx context.Context, UserID string, OrgID string) (string, error)
	OrganizationListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*thunderdome.UserOrganization
	OrganizationCreate(ctx context.Context, UserID string, OrgName string) (*thunderdome.Organization, error)
	OrganizationUpdate(ctx context.Context, OrgId string, OrgName string) (*thunderdome.Organization, error)
	OrganizationUserList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.OrganizationUser
	OrganizationAddUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error)
	OrganizationUpsertUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error)
	OrganizationUpdateUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error)
	OrganizationRemoveUser(ctx context.Context, OrganizationID string, UserID string) error
	OrganizationInviteUser(ctx context.Context, OrgID string, Email string, Role string) (string, error)
	OrganizationUserGetInviteByID(ctx context.Context, InviteID string) (thunderdome.OrganizationUserInvite, error)
	OrganizationDeleteUserInvite(ctx context.Context, InviteID string) error
	OrganizationGetUserInvites(ctx context.Context, orgId string) ([]thunderdome.OrganizationUserInvite, error)
	OrganizationTeamList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.Team
	OrganizationTeamCreate(ctx context.Context, OrgID string, TeamName string) (*thunderdome.Team, error)
	OrganizationTeamUserRole(ctx context.Context, UserID string, OrgID string, TeamID string) (string, string, error)
	OrganizationDelete(ctx context.Context, OrgID string) error
	OrganizationList(ctx context.Context, Limit int, Offset int) []*thunderdome.Organization
	OrganizationIsSubscribed(ctx context.Context, OrgID string) (bool, error)
	GetOrganizationMetrics(ctx context.Context, organizationID string) (*thunderdome.OrganizationMetrics, error)

	DepartmentUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string) (string, string, error)
	DepartmentGet(ctx context.Context, DepartmentID string) (*thunderdome.Department, error)
	OrganizationDepartmentList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.Department
	DepartmentCreate(ctx context.Context, OrgID string, OrgName string) (*thunderdome.Department, error)
	DepartmentUpdate(ctx context.Context, DeptId string, DeptName string) (*thunderdome.Department, error)
	DepartmentTeamList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*thunderdome.Team
	DepartmentTeamCreate(ctx context.Context, DepartmentID string, TeamName string) (*thunderdome.Team, error)
	DepartmentUserList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*thunderdome.DepartmentUser
	DepartmentAddUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error)
	DepartmentUpsertUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error)
	DepartmentUpdateUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error)
	DepartmentRemoveUser(ctx context.Context, DepartmentID string, UserID string) error
	DepartmentTeamUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string, TeamID string) (string, string, string, error)
	DepartmentDelete(ctx context.Context, DepartmentID string) error
	DepartmentInviteUser(ctx context.Context, DeptID string, Email string, Role string) (string, error)
	DepartmentUserGetInviteByID(ctx context.Context, InviteID string) (thunderdome.DepartmentUserInvite, error)
	DepartmentDeleteUserInvite(ctx context.Context, InviteID string) error
	DepartmentGetUserInvites(ctx context.Context, deptId string) ([]thunderdome.DepartmentUserInvite, error)
}

type TeamDataSvc interface {
	TeamUserRole(ctx context.Context, UserID string, TeamID string) (string, error)
	TeamGet(ctx context.Context, TeamID string) (*thunderdome.Team, error)
	TeamListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*thunderdome.UserTeam
	TeamListByUserNonOrg(ctx context.Context, UserID string, Limit int, Offset int) []*thunderdome.UserTeam
	TeamCreate(ctx context.Context, UserID string, TeamName string) (*thunderdome.Team, error)
	TeamUpdate(ctx context.Context, TeamId string, TeamName string) (*thunderdome.Team, error)
	TeamAddUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error)
	TeamUserList(ctx context.Context, TeamID string, Limit int, Offset int) ([]*thunderdome.TeamUser, int, error)
	TeamUpdateUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error)
	TeamRemoveUser(ctx context.Context, TeamID string, UserID string) error
	TeamInviteUser(ctx context.Context, TeamID string, Email string, Role string) (string, error)
	TeamUserGetInviteByID(ctx context.Context, InviteID string) (thunderdome.TeamUserInvite, error)
	TeamDeleteUserInvite(ctx context.Context, InviteID string) error
	TeamGetUserInvites(ctx context.Context, teamId string) ([]thunderdome.TeamUserInvite, error)
	TeamPokerList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Poker
	TeamAddPoker(ctx context.Context, TeamID string, PokerID string) error
	TeamRemovePoker(ctx context.Context, TeamID string, PokerID string) error
	TeamDelete(ctx context.Context, TeamID string) error
	TeamRetroList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Retro
	TeamAddRetro(ctx context.Context, TeamID string, RetroID string) error
	TeamRemoveRetro(ctx context.Context, TeamID string, RetroID string) error
	TeamStoryboardList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Storyboard
	TeamAddStoryboard(ctx context.Context, TeamID string, StoryboardID string) error
	TeamRemoveStoryboard(ctx context.Context, TeamID string, StoryboardID string) error
	TeamList(ctx context.Context, Limit int, Offset int) ([]*thunderdome.Team, int)
	TeamIsSubscribed(ctx context.Context, TeamID string) (bool, error)
	GetTeamMetrics(ctx context.Context, teamID string) (*thunderdome.TeamMetrics, error)
	TeamUserRoles(ctx context.Context, UserID string, TeamID string) (*thunderdome.UserTeamRoleInfo, error)
}
