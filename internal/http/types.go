package http

import (
	"context"
	"net/http"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/webhook/subscription"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

const (
	contextKeyUserID          contextKey = "userId"
	contextKeyUserType        contextKey = "userType"
	apiKeyHeaderName          string     = "X-API-Key"
	contextKeyUserTeamRoles   contextKey = "userTeamRoles"
	contextKeyOrgRole         contextKey = "orgRole"
	contextKeyDepartmentRole  contextKey = "departmentRole"
	contextKeyUserProjectRole contextKey = "userProjectRole"
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
	// Feature flag for Projects
	FeatureProject bool
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
	OIDCAuth   AuthProvider
	WebsocketConfig
}

type Service struct {
	Config               *Config
	Cookie               CookieManager
	UIConfig             thunderdome.UIConfig
	Handler              http.Handler
	Email                EmailService
	Logger               *otelzap.Logger
	UserDataSvc          UserDataSvc
	ApiKeyDataSvc        APIKeyDataSvc
	AlertDataSvc         AlertDataSvc
	AuthDataSvc          AuthDataSvc
	PokerDataSvc         PokerDataSvc
	CheckinDataSvc       CheckinDataSvc
	RetroDataSvc         RetroDataSvc
	StoryboardDataSvc    StoryboardDataSvc
	TeamDataSvc          TeamDataSvc
	OrganizationDataSvc  OrganizationDataSvc
	AdminDataSvc         AdminDataSvc
	JiraDataSvc          JiraDataSvc
	SubscriptionDataSvc  SubscriptionDataSvc
	RetroTemplateDataSvc RetroTemplateDataSvc
	SubscriptionSvc      *subscription.Service
	ProjectDataSvc       ProjectDataSvc
}

// standardJsonResponse structure used for all restful APIs response body
type standardJsonResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Data    any    `json:"data" swaggertype:"object"`
	Meta    any    `json:"meta" swaggertype:"object"`
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
	ListSupportTickets(ctx context.Context, limit, offset int) ([]*thunderdome.SupportTicket, int, error)
	GetSupportTicketByID(ctx context.Context, ticketID string) (*thunderdome.SupportTicket, error)
	UpdateSupportTicket(ctx context.Context, ticket *thunderdome.SupportTicket) error
	DeleteSupportTicket(ctx context.Context, ticketID string) error
	ListAdminUsers(ctx context.Context, limit int, offset int) ([]*thunderdome.User, int, error)
}

type AlertDataSvc interface {
	GetActiveAlerts(ctx context.Context) []any
	AlertsList(ctx context.Context, limit int, offset int) ([]*thunderdome.Alert, int, error)
	AlertsCreate(ctx context.Context, name string, alertType string, content string, active bool, allowDismiss bool, registeredOnly bool) error
	AlertsUpdate(ctx context.Context, alertID string, name string, alertType string, content string, active bool, allowDismiss bool, registeredOnly bool) error
	AlertDelete(ctx context.Context, alertID string) error
}

type APIKeyDataSvc interface {
	GenerateAPIKey(ctx context.Context, userID string, keyName string) (*thunderdome.APIKey, error)
	GetUserAPIKeys(ctx context.Context, userID string) ([]*thunderdome.APIKey, error)
	GetAPIKeyUser(ctx context.Context, apiKey string) (*thunderdome.User, error)
	GetAPIKeys(ctx context.Context, limit int, offset int) []*thunderdome.UserAPIKey
	UpdateUserAPIKey(ctx context.Context, userID string, keyID string, active bool) ([]*thunderdome.APIKey, error)
	DeleteUserAPIKey(ctx context.Context, userID string, keyID string) ([]*thunderdome.APIKey, error)
}

type AuthDataSvc interface {
	AuthUser(ctx context.Context, email string, password string) (*thunderdome.User, *thunderdome.Credential, string, error)
	OauthCreateNonce(ctx context.Context) (string, error)
	OauthValidateNonce(ctx context.Context, nonceId string) error
	OauthAuthUser(ctx context.Context, provider string, sub string, email string, emailVerified bool, name string, pictureUrl string) (*thunderdome.User, string, error)
	OauthUpsertUser(ctx context.Context, provider string, sub string, email string, emailVerified bool, name string, pictureUrl string) (*thunderdome.User, string, error)
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
	GetSessionUserByID(ctx context.Context, sessionId string) (*thunderdome.User, error)
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
	FindInstancesByUserID(ctx context.Context, userId string) ([]thunderdome.JiraInstance, error)
	GetInstanceByID(ctx context.Context, instanceId string) (thunderdome.JiraInstance, error)
	CreateInstance(ctx context.Context, userId string, host string, clientMail string, accessToken string, jiraDataCenter bool) (thunderdome.JiraInstance, error)
	UpdateInstance(ctx context.Context, instanceId string, host string, clientMail string, accessToken string) (thunderdome.JiraInstance, error)
	DeleteInstance(ctx context.Context, instanceId string) error
}

type OrganizationDataSvc interface {
	OrganizationGetByID(ctx context.Context, orgID string) (*thunderdome.Organization, error)
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
	DepartmentGetByID(ctx context.Context, departmentID string) (*thunderdome.Department, error)
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
	TeamUserRoleByUserID(ctx context.Context, userID string, teamID string) (string, error)
	TeamGetByID(ctx context.Context, teamID string) (*thunderdome.Team, error)
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
	TeamPokerList(ctx context.Context, teamID string, limit int, offset int) ([]*thunderdome.Poker, int)
	TeamAddPoker(ctx context.Context, teamID string, pokerID string) error
	TeamRemovePoker(ctx context.Context, teamID string, pokerID string) error
	TeamDelete(ctx context.Context, teamID string) error
	TeamRetroList(ctx context.Context, teamID string, limit int, offset int) ([]*thunderdome.Retro, int)
	TeamAddRetro(ctx context.Context, teamID string, retroID string) error
	TeamRemoveRetro(ctx context.Context, teamID string, retroID string) error
	TeamStoryboardList(ctx context.Context, teamID string, limit int, offset int) ([]*thunderdome.Storyboard, int)
	TeamAddStoryboard(ctx context.Context, teamID string, storyboardID string) error
	TeamRemoveStoryboard(ctx context.Context, teamID string, storyboardID string) error
	TeamList(ctx context.Context, limit int, offset int) ([]*thunderdome.Team, int)
	TeamIsSubscribed(ctx context.Context, teamID string) (bool, error)
	GetTeamMetrics(ctx context.Context, teamID string) (*thunderdome.TeamMetrics, error)
	TeamUserRolesByUserID(ctx context.Context, userID string, teamID string) (*thunderdome.UserTeamRoleInfo, error)
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
	ProjectIsSubscribed(ctx context.Context, projectID string) (bool, error)
}

type UserDataSvc interface {
	GetUserByID(ctx context.Context, userID string) (*thunderdome.User, error)
	GetGuestUserByID(ctx context.Context, userID string) (*thunderdome.User, error)
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
	GetUserCredentialByUserID(ctx context.Context, userID string) (*thunderdome.Credential, error)
	RequestEmailChange(ctx context.Context, userId string) (string, error)
	ConfirmEmailChange(ctx context.Context, userId string, token string, newEmail string) error
	CreateSupportTicket(ctx context.Context, userId, fullName, email, inquiry string) (thunderdome.SupportTicket, error)
}

type PokerDataSvc interface {
	// CreateGame creates a new poker game
	CreateGame(ctx context.Context, facilitatorID string, name string, estimationScaleID string, pointValuesAllowed []string, stories []*thunderdome.Story, autoFinishVoting bool, pointAverageRounding string, joinCode string, facilitatorCode string, hideVoterIdentity bool) (*thunderdome.Poker, error)
	// TeamCreateGame creates a new poker game for a team
	TeamCreateGame(ctx context.Context, teamID string, facilitatorID string, name string, estimationScaleID string, pointValuesAllowed []string, stories []*thunderdome.Story, autoFinishVoting bool, pointAverageRounding string, joinCode string, facilitatorCode string, hideVoterIdentity bool) (*thunderdome.Poker, error)
	// UpdateGame updates an existing poker game
	UpdateGame(pokerID string, name string, pointValuesAllowed []string, autoFinishVoting bool, pointAverageRounding string, hideVoterIdentity bool, joinCode string, facilitatorCode string, teamID string) error
	// GetFacilitatorCode retrieves the facilitator code for a poker game
	GetFacilitatorCode(pokerID string) (string, error)
	// GetGameByID retrieves a poker game by its ID
	GetGameByID(pokerID string, userID string) (*thunderdome.Poker, error)
	// GetGamesByUser retrieves a list of poker games for a user
	GetGamesByUser(userID string, limit int, offset int) ([]*thunderdome.Poker, int, error)
	// ConfirmFacilitator confirms a user as a facilitator for a poker game
	ConfirmFacilitator(pokerID string, userID string) error
	// GetUserActiveStatus retrieves the active status of a user in a poker game
	GetUserActiveStatus(pokerID string, userID string) error
	// GetUsers retrieves a list of users in a poker game
	GetUsers(pokerID string) []*thunderdome.PokerUser
	// GetActiveUsers retrieves a list of active users in a poker game
	GetActiveUsers(pokerID string) []*thunderdome.PokerUser
	// AddUser adds a user to a poker game
	AddUser(pokerID string, userID string) ([]*thunderdome.PokerUser, error)
	// RetreatUser sets a user as inactive in a poker game
	RetreatUser(pokerID string, userID string) []*thunderdome.PokerUser
	// AbandonGame sets a user as abandoned in a poker game
	AbandonGame(pokerID string, userID string) ([]*thunderdome.PokerUser, error)
	// AddFacilitator adds a facilitator to a poker game
	AddFacilitator(pokerID string, userID string) ([]string, error)
	// RemoveFacilitator removes a facilitator from a poker game
	RemoveFacilitator(pokerID string, userID string) ([]string, error)
	// ToggleSpectator toggles a user's spectator status in a poker game
	ToggleSpectator(pokerID string, userID string, spectator bool) ([]*thunderdome.PokerUser, error)
	// DeleteGame deletes a poker game
	DeleteGame(pokerID string) error
	// AddFacilitatorsByEmail adds facilitators to a poker game by email
	AddFacilitatorsByEmail(ctx context.Context, pokerID string, facilitatorEmails []string) ([]string, error)
	// GetGames retrieves a list of poker games
	GetGames(limit int, offset int) ([]*thunderdome.Poker, int, error)
	// GetActiveGames retrieves a list of active poker games
	GetActiveGames(limit int, offset int) ([]*thunderdome.Poker, int, error)
	// PurgeOldGames purges poker games older than a specified number of days
	PurgeOldGames(ctx context.Context, daysOld int) error
	// GetStories retrieves a list of stories in a poker game
	GetStories(pokerID string, userID string) []*thunderdome.Story
	// CreateStory creates a new story in a poker game
	CreateStory(pokerID string, name string, storyType string, referenceID string, link string, description string, acceptanceCriteria string, priority int32) ([]*thunderdome.Story, error)
	// ActivateStoryVoting activates voting for a story in a poker game
	ActivateStoryVoting(pokerID string, storyID string) ([]*thunderdome.Story, error)
	// SetVote sets a user's vote for a story in a poker game
	SetVote(pokerID string, userID string, storyID string, voteValue string) (stories []*thunderdome.Story, allUsersVoted bool)
	// RetractVote retracts a user's vote for a story in a poker game
	RetractVote(pokerID string, userID string, storyID string) ([]*thunderdome.Story, error)
	// EndStoryVoting ends voting for a story in a poker game
	EndStoryVoting(pokerID string, storyID string) ([]*thunderdome.Story, error)
	// SkipStory skips a story in a poker game
	SkipStory(pokerID string, storyID string) ([]*thunderdome.Story, error)
	// UpdateStory updates an existing story in a poker game
	UpdateStory(pokerID string, storyID string, name string, storyType string, referenceID string, link string, description string, acceptanceCriteria string, priority int32) ([]*thunderdome.Story, error)
	// DeleteStory deletes a story from a poker game
	DeleteStory(pokerID string, storyID string) ([]*thunderdome.Story, error)
	// ArrangeStory sets the position of the story relative to the story it's being placed before
	ArrangeStory(pokerID string, storyID string, beforeStoryID string) ([]*thunderdome.Story, error)
	// FinalizeStory finalizes the points for a story in a poker game
	FinalizeStory(pokerID string, storyID string, points string) ([]*thunderdome.Story, error)
	// GetEstimationScales retrieves a list of estimation scales
	GetEstimationScales(ctx context.Context, limit, offset int) ([]*thunderdome.EstimationScale, int, error)
	// GetPublicEstimationScales retrieves a list of public estimation scales
	GetPublicEstimationScales(ctx context.Context, limit, offset int) ([]*thunderdome.EstimationScale, int, error)
	// CreateEstimationScale creates a new estimation scale
	CreateEstimationScale(ctx context.Context, scale *thunderdome.EstimationScale) (*thunderdome.EstimationScale, error)
	// UpdateEstimationScale updates an existing estimation scale
	UpdateEstimationScale(ctx context.Context, scale *thunderdome.EstimationScale) (*thunderdome.EstimationScale, error)
	// DeleteEstimationScale deletes an estimation scale by its ID
	DeleteEstimationScale(ctx context.Context, scaleID string) error
	// GetDefaultEstimationScale retrieves the default estimation scale for an organization or team
	GetDefaultEstimationScale(ctx context.Context, organizationID, teamID string) (*thunderdome.EstimationScale, error)
	// GetDefaultPublicEstimationScale retrieves the default public estimation scale
	GetDefaultPublicEstimationScale(ctx context.Context) (*thunderdome.EstimationScale, error)
	// GetPublicEstimationScale retrieves a public estimation scale by its ID
	GetPublicEstimationScale(ctx context.Context, id string) (*thunderdome.EstimationScale, error)
	// GetOrganizationEstimationScales retrieves a list of estimation scales for an organization
	GetOrganizationEstimationScales(ctx context.Context, orgID string, limit, offset int) ([]*thunderdome.EstimationScale, int, error)
	// GetTeamEstimationScales retrieves a list of estimation scales for a team
	GetTeamEstimationScales(ctx context.Context, teamID string, limit, offset int) ([]*thunderdome.EstimationScale, int, error)
	// GetEstimationScale retrieves an estimation scale by its ID
	GetEstimationScale(ctx context.Context, scaleID string) (*thunderdome.EstimationScale, error)
	// DeleteOrganizationEstimationScale deletes an organization's estimation scale by its ID
	DeleteOrganizationEstimationScale(ctx context.Context, orgID string, scaleID string) error
	// DeleteTeamEstimationScale deletes a team's estimation scale by its ID
	DeleteTeamEstimationScale(ctx context.Context, teamID string, scaleID string) error
	// UpdateOrganizationEstimationScale updates an existing organization estimation scale
	UpdateOrganizationEstimationScale(ctx context.Context, scale *thunderdome.EstimationScale) (*thunderdome.EstimationScale, error)
	// UpdateTeamEstimationScale updates an existing team estimation scale
	UpdateTeamEstimationScale(ctx context.Context, scale *thunderdome.EstimationScale) (*thunderdome.EstimationScale, error)
	// GetSettingsByOrganization retrieves poker settings for an organization
	GetSettingsByOrganization(ctx context.Context, orgID string) (*thunderdome.PokerSettings, error)
	// GetSettingsByDepartment retrieves poker settings for a department
	GetSettingsByDepartment(ctx context.Context, deptID string) (*thunderdome.PokerSettings, error)
	// GetSettingsByTeam retrieves poker settings for a team
	GetSettingsByTeam(ctx context.Context, teamID string) (*thunderdome.PokerSettings, error)
	// CreateSettings creates poker settings for an organization or department or team
	CreateSettings(ctx context.Context, settings *thunderdome.PokerSettings) (*thunderdome.PokerSettings, error)
	// UpdateSettings updates existing poker settings
	UpdateSettings(ctx context.Context, settings *thunderdome.PokerSettings) (*thunderdome.PokerSettings, error)
	// UpdateOrganizationSettings updates existing organization poker settings
	UpdateOrganizationSettings(ctx context.Context, settings *thunderdome.PokerSettings) (*thunderdome.PokerSettings, error)
	// UpdateDepartmentSettings updates existing department poker settings
	UpdateDepartmentSettings(ctx context.Context, settings *thunderdome.PokerSettings) (*thunderdome.PokerSettings, error)
	// UpdateTeamSettings updates existing team poker settings
	UpdateTeamSettings(ctx context.Context, settings *thunderdome.PokerSettings) (*thunderdome.PokerSettings, error)
	// DeleteSettings deletes poker settings by its ID
	DeleteSettings(ctx context.Context, id string) error
	// GetSettingsByID retrieves poker settings by its ID
	GetSettingsByID(ctx context.Context, id string) (*thunderdome.PokerSettings, error)
	// EndGame ends a poker game with a specified reason
	EndGame(ctx context.Context, pokerID string, endReason string) (string, time.Time, error)
}

type RetroDataSvc interface {
	CreateRetro(ctx context.Context, ownerID, teamID string, retroName, joinCode, facilitatorCode string, maxVotes int, brainstormVisibility string, phaseTimeLimitMin int, phaseAutoAdvance bool, allowCumulativeVoting bool, hideVotesDuringVoting bool, templateID string) (*thunderdome.Retro, error)
	EditRetro(retroID string, retroName string, joinCode string, facilitatorCode string, maxVotes int, brainstormVisibility string, phaseAutoAdvance bool, hideVotesDuringVoting bool) error
	RetroGetByID(retroID string, userID string) (*thunderdome.Retro, error)
	RetroGetByUser(userID string, limit int, offset int) ([]*thunderdome.Retro, int, error)
	RetroConfirmFacilitator(retroID string, userID string) error
	RetroGetUsers(retroID string) []*thunderdome.RetroUser
	GetRetroFacilitators(retroID string) []string
	RetroAddUser(retroID string, userID string) ([]*thunderdome.RetroUser, error)
	RetroFacilitatorAdd(retroID string, userID string) ([]string, error)
	RetroFacilitatorRemove(retroID string, userID string) ([]string, error)
	RetroRetreatUser(retroID string, userID string) []*thunderdome.RetroUser
	RetroAbandon(retroID string, userID string) ([]*thunderdome.RetroUser, error)
	RetroAdvancePhase(retroID string, phase string) (*thunderdome.Retro, error)
	RetroDelete(retroID string) error
	GetRetroUserActiveStatus(retroID string, userID string) error
	GetRetros(limit int, offset int) ([]*thunderdome.Retro, int, error)
	GetActiveRetros(limit int, offset int) ([]*thunderdome.Retro, int, error)
	GetRetroFacilitatorCode(retroID string) (string, error)
	CleanRetros(ctx context.Context, daysOld int) error
	MarkUserReady(retroID string, userID string) ([]string, error)
	UnmarkUserReady(retroID string, userID string) ([]string, error)

	CreateRetroAction(retroID string, userID string, content string) ([]*thunderdome.RetroAction, error)
	UpdateRetroAction(retroID string, actionID string, content string, completed bool) (Actions []*thunderdome.RetroAction, DeleteError error)
	DeleteRetroAction(retroID string, userID string, actionID string) ([]*thunderdome.RetroAction, error)
	GetRetroActions(retroID string) []*thunderdome.RetroAction
	GetTeamRetroActions(teamID string, limit int, offset int, completed bool) ([]*thunderdome.RetroAction, int, error)
	RetroActionCommentAdd(retroID string, actionID string, userID string, comment string) ([]*thunderdome.RetroAction, error)
	RetroActionCommentEdit(retroID string, actionID string, commentID string, comment string) ([]*thunderdome.RetroAction, error)
	RetroActionCommentDelete(retroID string, actionID string, commentID string) ([]*thunderdome.RetroAction, error)
	RetroActionAssigneeAdd(retroID string, actionID string, userID string) ([]*thunderdome.RetroAction, error)
	RetroActionAssigneeDelete(retroID string, actionID string, userID string) ([]*thunderdome.RetroAction, error)

	CreateRetroItem(retroID string, userID string, itemType string, content string) ([]*thunderdome.RetroItem, error)
	GroupRetroItem(retroID string, itemId string, groupId string) (thunderdome.RetroItem, error)
	DeleteRetroItem(retroID string, userID string, itemType string, itemID string) ([]*thunderdome.RetroItem, error)
	GetRetroItems(retroID string) []*thunderdome.RetroItem
	GetRetroGroups(retroID string) []*thunderdome.RetroGroup
	GroupNameChange(retroID string, groupID string, name string) (thunderdome.RetroGroup, error)
	GetRetroVotes(retroID string) []*thunderdome.RetroVote
	GroupUserVote(retroID string, groupID string, userID string) ([]*thunderdome.RetroVote, error)
	GroupUserSubtractVote(retroID string, groupID string, userID string) ([]*thunderdome.RetroVote, error)
	ItemCommentAdd(retroID string, itemID string, userID string, comment string) ([]*thunderdome.RetroItem, error)
	ItemCommentEdit(retroID string, commentID string, comment string) ([]*thunderdome.RetroItem, error)
	ItemCommentDelete(retroID string, commentID string) ([]*thunderdome.RetroItem, error)

	GetSettingsByOrganization(ctx context.Context, orgID string) (*thunderdome.RetroSettings, error)
	GetSettingsByDepartment(ctx context.Context, deptID string) (*thunderdome.RetroSettings, error)
	GetSettingsByTeam(ctx context.Context, teamID string) (*thunderdome.RetroSettings, error)
	CreateSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error)
	UpdateSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error)
	UpdateOrganizationSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error)
	UpdateDepartmentSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error)
	UpdateTeamSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error)
	DeleteSettings(ctx context.Context, id string) error
	GetSettingsByID(ctx context.Context, id string) (*thunderdome.RetroSettings, error)
}

type RetroTemplateDataSvc interface {
	// GetPublicTemplates retrieves all public retro templates
	GetPublicTemplates(ctx context.Context) ([]*thunderdome.RetroTemplate, error)
	// GetTemplatesByOrganization retrieves all templates for a specific organization
	GetTemplatesByOrganization(ctx context.Context, organizationID string) ([]*thunderdome.RetroTemplate, error)
	// GetTemplatesByTeam retrieves all templates for a specific team
	GetTemplatesByTeam(ctx context.Context, teamID string) ([]*thunderdome.RetroTemplate, error)
	// GetTemplateByID retrieves a specific template by its ID
	GetTemplateByID(ctx context.Context, templateID string) (*thunderdome.RetroTemplate, error)
	// CreateTemplate creates a new retro template
	CreateTemplate(ctx context.Context, template *thunderdome.RetroTemplate) error
	// UpdateTemplate updates an existing retro template
	UpdateTemplate(ctx context.Context, template *thunderdome.RetroTemplate) error
	// DeleteTemplate deletes a retro template by its ID
	DeleteTemplate(ctx context.Context, templateID string) error
	// ListTemplates retrieves a paginated list of templates
	ListTemplates(ctx context.Context, limit int, offset int) ([]*thunderdome.RetroTemplate, int, error)
	// GetDefaultPublicTemplate retrieves the default public template
	GetDefaultPublicTemplate(ctx context.Context) (*thunderdome.RetroTemplate, error)
	// GetDefaultTeamTemplate retrieves the default template for a given team
	GetDefaultTeamTemplate(ctx context.Context, teamID string) (*thunderdome.RetroTemplate, error)
	// GetDefaultOrganizationTemplate retrieves the default template for a given organization
	GetDefaultOrganizationTemplate(ctx context.Context, organizationID string) (*thunderdome.RetroTemplate, error)
	// UpdateTeamTemplate updates an existing team retro template
	UpdateTeamTemplate(ctx context.Context, template *thunderdome.RetroTemplate) error
	// UpdateOrganizationTemplate updates an existing organization retro template
	UpdateOrganizationTemplate(ctx context.Context, template *thunderdome.RetroTemplate) error
	// DeleteOrganizationTemplate deletes an organization retro template by its ID
	DeleteOrganizationTemplate(ctx context.Context, orgID string, templateID string) error
	// DeleteTeamTemplate deletes a team retro template by its ID
	DeleteTeamTemplate(ctx context.Context, teamID string, templateID string) error
}

type StoryboardDataSvc interface {
	CreateStoryboard(ctx context.Context, ownerID string, storyboardName string, joinCode string, facilitatorCode string) (*thunderdome.Storyboard, error)
	TeamCreateStoryboard(ctx context.Context, TeamID string, ownerID string, storyboardName string, joinCode string, facilitatorCode string) (*thunderdome.Storyboard, error)
	EditStoryboard(storyboardID string, storyboardName string, joinCode string, facilitatorCode string) error
	GetStoryboardByID(storyboardID string, userID string) (*thunderdome.Storyboard, error)
	GetStoryboardsByUser(userID string, limit int, offset int) ([]*thunderdome.Storyboard, int, error)
	ConfirmStoryboardFacilitator(storyboardID string, userID string) error
	GetStoryboardUsers(storyboardID string) []*thunderdome.StoryboardUser
	GetStoryboardPersonas(storyboardID string) []*thunderdome.StoryboardPersona
	GetStoryboards(limit int, offset int) ([]*thunderdome.Storyboard, int, error)
	GetActiveStoryboards(limit int, offset int) ([]*thunderdome.Storyboard, int, error)
	AddUserToStoryboard(storyboardID string, userID string) ([]*thunderdome.StoryboardUser, error)
	RetreatStoryboardUser(storyboardID string, userID string) []*thunderdome.StoryboardUser
	GetStoryboardUserActiveStatus(storyboardID string, userID string) error
	AbandonStoryboard(storyboardID string, userID string) ([]*thunderdome.StoryboardUser, error)
	StoryboardFacilitatorAdd(StoryboardId string, userID string) (*thunderdome.Storyboard, error)
	StoryboardFacilitatorRemove(StoryboardId string, userID string) (*thunderdome.Storyboard, error)
	GetStoryboardFacilitatorCode(storyboardID string) (string, error)
	StoryboardReviseColorLegend(storyboardID string, userID string, colorLegend string) (*thunderdome.Storyboard, error)
	DeleteStoryboard(storyboardID string, userID string) error
	CleanStoryboards(ctx context.Context, daysOld int) error

	AddStoryboardPersona(storyboardID string, userID string, name string, role string, description string) ([]*thunderdome.StoryboardPersona, error)
	UpdateStoryboardPersona(storyboardID string, userID string, personaID string, name string, role string, description string) ([]*thunderdome.StoryboardPersona, error)
	DeleteStoryboardPersona(storyboardID string, userID string, personaID string) ([]*thunderdome.StoryboardPersona, error)

	CreateStoryboardGoal(storyboardID string, userID string, goalName string) (*thunderdome.StoryboardGoal, error)
	ReviseGoalName(storyboardID string, userID string, goalID string, goalName string) ([]*thunderdome.StoryboardGoal, error)
	DeleteStoryboardGoal(storyboardID string, userID string, goalID string) ([]*thunderdome.StoryboardGoal, error)
	GetStoryboardGoals(storyboardID string) []*thunderdome.StoryboardGoal
	GetStoryboardGoal(storyboardID string, goalID string) (*thunderdome.StoryboardGoal, error)

	CreateStoryboardColumn(storyboardID string, goalID string, userID string) (*thunderdome.StoryboardColumn, error)
	ReviseStoryboardColumn(storyboardID string, userID string, columnID string, columnName string) ([]*thunderdome.StoryboardGoal, error)
	DeleteStoryboardColumn(storyboardID string, userID string, columnID string) ([]*thunderdome.StoryboardGoal, error)
	ColumnPersonaAdd(storyboardID string, columnID string, personaID string) ([]*thunderdome.StoryboardGoal, error)
	ColumnPersonaRemove(storyboardID string, columnID string, personaID string) ([]*thunderdome.StoryboardGoal, error)
	MoveStoryboardColumn(storyboardID string, userID string, columnID string, goalID string, placeBeforeID string) error

	CreateStoryboardStory(storyboardID string, goalID string, columnID string, userID string) (*thunderdome.StoryboardStory, error)
	ReviseStoryName(storyboardID string, userID string, storyID string, storyName string) ([]*thunderdome.StoryboardGoal, error)
	ReviseStoryContent(storyboardID string, userID string, storyID string, storyContent string) ([]*thunderdome.StoryboardGoal, error)
	ReviseStoryColor(storyboardID string, userID string, storyID string, storyColor string) ([]*thunderdome.StoryboardGoal, error)
	ReviseStoryPoints(storyboardID string, userID string, storyID string, points int) ([]*thunderdome.StoryboardGoal, error)
	ReviseStoryClosed(storyboardID string, userID string, storyID string, closed bool) ([]*thunderdome.StoryboardGoal, error)
	ReviseStoryLink(storyboardID string, userID string, storyID string, link string) ([]*thunderdome.StoryboardGoal, error)
	MoveStoryboardStory(storyboardID string, userID string, storyID string, goalID string, columnID string, placeBefore string) ([]*thunderdome.StoryboardGoal, error)
	DeleteStoryboardStory(storyboardID string, userID string, storyID string) ([]*thunderdome.StoryboardGoal, error)
	AddStoryComment(storyboardID string, userID string, storyID string, comment string) ([]*thunderdome.StoryboardGoal, error)
	EditStoryComment(storyboardID string, commentID string, comment string) ([]*thunderdome.StoryboardGoal, error)
	DeleteStoryComment(storyboardID string, commentID string) ([]*thunderdome.StoryboardGoal, error)
}

type EmailService interface {
	SendWelcome(userName string, userEmail string, verifyID string) error
	SendEmailVerification(userName string, userEmail string, verifyID string) error
	SendForgotPassword(userName string, userEmail string, resetID string) error
	SendPasswordReset(userName string, userEmail string) error
	SendPasswordUpdate(userName string, userEmail string) error
	SendDeleteConfirmation(userName string, userEmail string) error
	SendTeamInvite(TeamName string, userEmail string, inviteID string) error
	SendOrganizationInvite(organizationName string, userEmail string, inviteID string) error
	SendDepartmentInvite(organizationName string, departmentName string, userEmail string, inviteID string) error
	// SendRetroOverview sends the retro overview (items, action items) email to attendees
	SendRetroOverview(retro *thunderdome.Retro, template *thunderdome.RetroTemplate, userName string, userEmail string) error
	SendEmailChangeRequest(userName string, userEmail string, changeId string) error
	SendEmailChangeConfirmation(userName string, userEmail string, newEmail string) error
	SendNewTicketToAdmins(adminUser thunderdome.User, ticketID string) error
}

// ProjectDataSvc represents the interface for project data operations
type ProjectDataSvc interface {
	// Basic CRUD operations
	CreateProject(ctx context.Context, project *thunderdome.Project) error
	GetProjectByID(ctx context.Context, projectID string) (*thunderdome.Project, error)
	UpdateProject(ctx context.Context, project *thunderdome.Project) error
	DeleteProject(ctx context.Context, projectID string) error
	ListProjects(ctx context.Context, limit int, offset int) ([]*thunderdome.Project, int, error)

	// Member-related operations
	IsUserProjectMember(ctx context.Context, userID, projectID string) (bool, string, error)

	// Organization-scoped operations
	GetProjectsByOrganization(ctx context.Context, organizationID string) ([]*thunderdome.Project, error)
	UpdateOrganizationProject(ctx context.Context, project *thunderdome.Project) error
	DeleteOrganizationProject(ctx context.Context, orgID string, projectID string) error

	// Department-scoped operations
	GetProjectsByDepartment(ctx context.Context, departmentID string) ([]*thunderdome.Project, error)
	UpdateDepartmentProject(ctx context.Context, project *thunderdome.Project) error
	DeleteDepartmentProject(ctx context.Context, deptID string, projectID string) error

	// Team-scoped operations
	GetProjectsByTeam(ctx context.Context, teamID string) ([]*thunderdome.Project, error)
	UpdateTeamProject(ctx context.Context, project *thunderdome.Project) error
	DeleteTeamProject(ctx context.Context, teamID string, projectID string) error

	// Project associated Storyboards
	AssociateStoryboard(ctx context.Context, projectID string, storyboardID string) error
	ListStoryboards(ctx context.Context, projectId string, limit int, offset int) ([]*thunderdome.Storyboard, error)
	RemoveStoryboard(ctx context.Context, projectID string, storyboardID string) error

	// Project associated Retros
	AssociateRetro(ctx context.Context, projectID string, retroID string) error
	ListRetros(ctx context.Context, projectId string, limit int, offset int) ([]*thunderdome.Retro, error)
	RemoveRetro(ctx context.Context, projectID string, retroID string) error

	// Project associated Poker Games
	AssociatePoker(ctx context.Context, projectID string, pokerID string) error
	ListPokerGames(ctx context.Context, projectId string, limit int, offset int) ([]*thunderdome.Poker, error)
	RemovePokerGame(ctx context.Context, projectID string, pokerID string) error
}
