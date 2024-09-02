package config

import "github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

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
}

type Analytics struct {
	Enabled bool
	ID      string
}

type Admin struct {
	Email string
}

type Otel struct {
	Enabled      bool
	ServiceName  string `mapstructure:"service_name"`
	CollectorUrl string `mapstructure:"collector_url"`
	InsecureMode bool   `mapstructure:"insecure_mode"`
}

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

type AppConfig struct {
	AesHashkey                string   `mapstructure:"aes_hashkey"`
	AllowedPointValues        []string `mapstructure:"allowedPointValues"`
	DefaultPointValues        []string `mapstructure:"defaultPointValues"`
	ShowWarriorRank           bool     `mapstructure:"show_warrior_rank"`
	AvatarService             string   `mapstructure:"avatar_service"`
	ToastTimeout              int      `mapstructure:"toast_timeout"`
	AllowGuests               bool     `mapstructure:"allow_guests"`
	AllowRegistration         bool     `mapstructure:"allow_registration"`
	AllowJiraImport           bool     `mapstructure:"allow_jira_import"`
	AllowCsvImport            bool     `mapstructure:"allow_csv_import"`
	DefaultLocale             string   `mapstructure:"default_locale"`
	AllowExternalApi          bool     `mapstructure:"allow_external_api"`
	ExternalApiVerifyRequired bool     `mapstructure:"external_api_verify_required"`
	UserApikeyLimit           int      `mapstructure:"user_apikey_limit"`
	ShowActiveCountries       bool     `mapstructure:"show_active_countries"`
	CleanupBattlesDaysOld     int      `mapstructure:"cleanup_battles_days_old"`
	CleanupGuestsDaysOld      int      `mapstructure:"cleanup_guests_days_old"`
	CleanupRetrosDaysOld      int      `mapstructure:"cleanup_retros_days_old"`
	CleanupStoryboardsDaysOld int      `mapstructure:"cleanup_storyboards_days_old"`
	OrganizationsEnabled      bool     `mapstructure:"organizations_enabled"`
	RequireTeams              bool     `mapstructure:"require_teams"`
	SubscriptionsEnabled      bool     `mapstructure:"subscriptions_enabled"`
}

type Feature struct {
	Poker      bool
	Retro      bool
	Storyboard bool
}

type Google struct {
	Enabled      bool   `mapstructure:"enabled"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type Auth struct {
	Method string
	Ldap   AuthLdap
	Header AuthHeader
	Google
}

type AuthHeader struct {
	UsernameHeader string `mapstructure:"usernameHeader"`
	EmailHeader    string `mapstructure:"emailHeader"`
}

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
