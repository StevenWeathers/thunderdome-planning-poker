package thunderdome

type SubscriptionPlanConfig struct {
	Enabled           bool   `mapstructure:"enabled"`
	MonthPrice        string `mapstructure:"month_price"`
	YearPrice         string `mapstructure:"year_price"`
	MonthCheckoutLink string `mapstructure:"month_checkout_link"`
	YearCheckoutLink  string `mapstructure:"year_checkout_link"`
}

type SubscriptionConfig struct {
	ManageLink    string                 `mapstructure:"manage_link"`
	AccountSecret string                 `mapstructure:"account_secret" json:"-"`
	WebhookSecret string                 `mapstructure:"webhook_secret" json:"-"`
	Individual    SubscriptionPlanConfig `mapstructure:"individual"`
	Team          SubscriptionPlanConfig `mapstructure:"team"`
	Organization  SubscriptionPlanConfig `mapstructure:"organization"`
}

type AppConfig struct {
	AllowedPointValues        []string
	DefaultPointValues        []string
	ShowWarriorRank           bool
	AvatarService             string
	ToastTimeout              int
	AllowGuests               bool
	AllowRegistration         bool
	AllowJiraImport           bool
	AllowCsvImport            bool
	DefaultLocale             string
	OrganizationsEnabled      bool
	AppVersion                string
	CookieName                string
	PathPrefix                string
	ExternalAPIEnabled        bool
	UserAPIKeyLimit           int
	CleanupGuestsDaysOld      int
	CleanupBattlesDaysOld     int
	CleanupRetrosDaysOld      int
	CleanupStoryboardsDaysOld int
	ShowActiveCountries       bool
	LdapEnabled               bool
	HeaderAuthEnabled         bool
	GoogleAuthEnabled         bool
	FeaturePoker              bool
	FeatureRetro              bool
	FeatureStoryboard         bool
	RequireTeams              bool
	RepoURL                   string
	SubscriptionsEnabled      bool
	Subscription              SubscriptionConfig
}

type UIConfig struct {
	AnalyticsEnabled bool
	AnalyticsID      string
	AppConfig        AppConfig
	ActiveAlerts     []interface{}
}
