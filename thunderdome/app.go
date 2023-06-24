package thunderdome

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
	FriendlyUIVerbs           bool
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
	FeaturePoker              bool
	FeatureRetro              bool
	FeatureStoryboard         bool
	RequireTeams              bool
}

type UIConfig struct {
	AnalyticsEnabled bool
	AnalyticsID      string
	AppConfig        AppConfig
	ActiveAlerts     []interface{}
}
