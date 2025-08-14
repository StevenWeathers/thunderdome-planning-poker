package thunderdome

// ApplicationStats includes counts of different data points of the application
type ApplicationStats struct {
	UnregisteredCount                int `json:"unregisteredUserCount"`
	RegisteredCount                  int `json:"registeredUserCount"`
	APIKeyCount                      int `json:"apikeyCount"`
	PokerCount                       int `json:"battleCount"`
	ActivePokerCount                 int `json:"activeBattleCount"`
	ActivePokerUserCount             int `json:"activeBattleUserCount"`
	PokerStoryCount                  int `json:"planCount"`
	OrganizationCount                int `json:"organizationCount"`
	DepartmentCount                  int `json:"departmentCount"`
	TeamCount                        int `json:"teamCount"`
	TeamCheckinsCount                int `json:"teamCheckinsCount"`
	RetroCount                       int `json:"retroCount"`
	ActiveRetroCount                 int `json:"activeRetroCount"`
	ActiveRetroUserCount             int `json:"activeRetroUserCount"`
	RetroItemCount                   int `json:"retroItemCount"`
	RetroActionCount                 int `json:"retroActionCount"`
	StoryboardCount                  int `json:"storyboardCount"`
	ActiveStoryboardCount            int `json:"activeStoryboardCount"`
	ActiveStoryboardUserCount        int `json:"activeStoryboardUserCount"`
	StoryboardGoalCount              int `json:"storyboardGoalCount"`
	StoryboardColumnCount            int `json:"storyboardColumnCount"`
	StoryboardStoryCount             int `json:"storyboardStoryCount"`
	StoryboardPersonaCount           int `json:"storyboardPersonaCount"`
	EstimationScaleCount             int `json:"estimationScaleCount"`
	PublicEstimationScaleCount       int `json:"publicEstimationScaleCount"`
	OrganizationEstimationScaleCount int `json:"organizationEstimationScaleCount"`
	TeamEstimationScaleCount         int `json:"teamEstimationScaleCount"`
	UserSubscriptionActiveCount      int `json:"userSubscriptionActiveCount"`
	TeamSubscriptionActiveCount      int `json:"teamSubscriptionActiveCount"`
	OrgSubscriptionActiveCount       int `json:"orgSubscriptionActiveCount"`
	RetroTemplateCount               int `json:"retroTemplateCount"`
	OrganizationRetroTemplateCount   int `json:"organizationRetroTemplateCount"`
	TeamRetroTemplateCount           int `json:"teamRetroTemplateCount"`
	PublicRetroTemplateCount         int `json:"publicRetroTemplateCount"`
	ProjectCount                     int `json:"projectCount"`
}
