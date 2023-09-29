package http

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/cookie"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

const (
	contextKeyUserID         contextKey = "userId"
	contextKeyUserType       contextKey = "userType"
	apiKeyHeaderName         string     = "X-API-Key"
	contextKeyOrgRole        contextKey = "orgRole"
	contextKeyDepartmentRole contextKey = "departmentRole"
	contextKeyTeamRole       contextKey = "teamRole"
	adminUserType            string     = "ADMIN"
	guestUserType            string     = "GUEST"
)

var validate *validator.Validate

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
	// Whether the external API is enabled
	ExternalAPIEnabled bool
	// Whether the external API requires user verified email
	ExternalAPIVerifyRequired bool
	// Number of API keys a user can create
	UserAPIKeyLimit int
	// Whether LDAP is enabled for authentication
	LdapEnabled bool
	// Whether header authentication is enabled
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
}

type Service struct {
	Config              *Config
	Cookie              *cookie.Cookie
	UIConfig            thunderdome.UIConfig
	Router              *mux.Router
	Email               thunderdome.EmailService
	Logger              *otelzap.Logger
	UserDataSvc         thunderdome.UserDataSvc
	ApiKeyDataSvc       thunderdome.APIKeyDataSvc
	AlertDataSvc        thunderdome.AlertDataSvc
	AuthDataSvc         thunderdome.AuthDataSvc
	PokerDataSvc        thunderdome.PokerDataSvc
	CheckinDataSvc      thunderdome.CheckinDataSvc
	RetroDataSvc        thunderdome.RetroDataSvc
	StoryboardDataSvc   thunderdome.StoryboardDataSvc
	TeamDataSvc         thunderdome.TeamDataSvc
	OrganizationDataSvc thunderdome.OrganizationDataSvc
	AdminDataSvc        thunderdome.AdminDataSvc
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
