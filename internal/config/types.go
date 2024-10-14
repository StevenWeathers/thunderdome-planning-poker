package config

import "github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

// Config is the main application configuration
type Config struct {
	Http
	Analytics
	Admin
	Otel
	Db
	Smtp
	Config AppConfig
	Feature
	Auth
	Subscription thunderdome.SubscriptionConfig
}

// Http is the application HTTP server configuration
type Http struct {
	Port                   string
	SecureCookie           bool   `mapstructure:"secure_cookie"`
	BackendCookieName      string `mapstructure:"backend_cookie_name"`
	SessionCookieName      string `mapstructure:"session_cookie_name"`
	FrontendCookieName     string `mapstructure:"frontend_cookie_name"`
	AuthStateCookieName    string `mapstructure:"auth_state_cookie_name"`
	Domain                 string
	PathPrefix             string `mapstructure:"path_prefix"`
	SecureProtocol         bool   `mapstructure:"secure_protocol"`
	WriteTimeout           int    `mapstructure:"write_timeout"`
	ReadTimeout            int    `mapstructure:"read_timeout"`
	IdleTimeout            int    `mapstructure:"idle_timeout"`
	ReadHeaderTimeout      int    `mapstructure:"read_header_timeout"`
	CookieHashkey          string `mapstructure:"cookie_hashkey"`
	WebsocketWriteWaitSec  int    `mapstructure:"websocket_write_wait_sec"`
	WebsocketPingPeriodSec int    `mapstructure:"websocket_ping_period_sec"`
	WebsocketPongWaitSec   int    `mapstructure:"websocket_pong_wait_sec"`
	WebsocketSubdomain     string `mapstructure:"websocket_subdomain"`
}

// Analytics is the application analytics configuration
type Analytics struct {
	Enabled bool
	ID      string
}

// Admin is the application admin configuration
type Admin struct {
	Email string
}

// Otel is the application OpenTelemetry configuration
type Otel struct {
	Enabled      bool
	ServiceName  string `mapstructure:"service_name"`
	CollectorUrl string `mapstructure:"collector_url"`
	InsecureMode bool   `mapstructure:"insecure_mode"`
}

// Db is the application database configuration
type Db struct {
	Host            string
	Port            int
	User            string
	Pass            string
	Name            string
	Sslmode         string
	MaxOpenConns    int `mapstructure:"max_open_conns"`
	MaxIdleConns    int `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`
}

// Smtp is the application SMTP configuration
type Smtp struct {
	Enabled       bool
	Host          string
	Port          int
	Secure        bool
	Sender        string
	User          string
	Pass          string
	SkipTLSVerify bool `mapstructure:"skip_tls_verify"`
	Auth          string
}

// AppConfig is the application configuration
type AppConfig struct {
	AesHashkey                  string   `mapstructure:"aes_hashkey"`
	AllowedPointValues          []string `mapstructure:"allowedPointValues"`
	DefaultPointValues          []string `mapstructure:"defaultPointValues"`
	ShowWarriorRank             bool     `mapstructure:"show_warrior_rank"`
	AvatarService               string   `mapstructure:"avatar_service"`
	ToastTimeout                int      `mapstructure:"toast_timeout"`
	AllowGuests                 bool     `mapstructure:"allow_guests"`
	AllowRegistration           bool     `mapstructure:"allow_registration"`
	AllowJiraImport             bool     `mapstructure:"allow_jira_import"`
	AllowCsvImport              bool     `mapstructure:"allow_csv_import"`
	DefaultLocale               string   `mapstructure:"default_locale"`
	AllowExternalApi            bool     `mapstructure:"allow_external_api"`
	ExternalApiVerifyRequired   bool     `mapstructure:"external_api_verify_required"`
	UserApikeyLimit             int      `mapstructure:"user_apikey_limit"`
	ShowActiveCountries         bool     `mapstructure:"show_active_countries"`
	CleanupBattlesDaysOld       int      `mapstructure:"cleanup_battles_days_old"`
	CleanupGuestsDaysOld        int      `mapstructure:"cleanup_guests_days_old"`
	CleanupRetrosDaysOld        int      `mapstructure:"cleanup_retros_days_old"`
	CleanupStoryboardsDaysOld   int      `mapstructure:"cleanup_storyboards_days_old"`
	OrganizationsEnabled        bool     `mapstructure:"organizations_enabled"`
	RequireTeams                bool     `mapstructure:"require_teams"`
	SubscriptionsEnabled        bool     `mapstructure:"subscriptions_enabled"`
	RetroDefaultTemplateID      string   `mapstructure:"retro_default_template_id"`
	DefaultPointAverageRounding string   `mapstructure:"default_point_average_rounding"`
}

// Feature is the application feature enablement configuration
type Feature struct {
	Poker      bool
	Retro      bool
	Storyboard bool
}

// Google is the application Google OAuth configuration
type Google struct {
	Enabled      bool   `mapstructure:"enabled"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// Auth is the application authentication configuration
type Auth struct {
	Method string
	Ldap   AuthLdap
	Header AuthHeader
	Google
}

// AuthHeader is the application authentication header configuration
type AuthHeader struct {
	UsernameHeader string `mapstructure:"usernameHeader"`
	EmailHeader    string `mapstructure:"emailHeader"`
}

// AuthLdap is the application LDAP authentication configuration
type AuthLdap struct {
	Url      string
	UseTls   bool `mapstructure:"use_tls"`
	Bindname string
	Bindpass string
	Basedn   string
	Filter   string
	MailAttr string `mapstructure:"mail_attr"`
	CnAttr   string `mapstructure:"cn_attr"`
}
