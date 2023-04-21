package db

import (
	"context"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
)

var _ Repository = &Database{}

type Repository interface {
	// Admin
	GetAppStats(ctx context.Context) (*model.ApplicationStats, error)
	PromoteUser(ctx context.Context, UserID string) error
	DemoteUser(ctx context.Context, UserID string) error
	DisableUser(ctx context.Context, UserID string) error
	EnableUser(ctx context.Context, UserID string) error
	CleanBattles(ctx context.Context, DaysOld int) error
	CleanRetros(ctx context.Context, DaysOld int) error
	CleanStoryboards(ctx context.Context, DaysOld int) error
	CleanGuests(ctx context.Context, DaysOld int) error
	LowercaseUserEmails(ctx context.Context) ([]*model.User, error)
	MergeDuplicateAccounts(ctx context.Context) ([]*model.User, error)
	OrganizationList(ctx context.Context, Limit int, Offset int) []*model.Organization
	TeamList(ctx context.Context, Limit int, Offset int) ([]*model.Team, int)
	GetAPIKeys(ctx context.Context, Limit int, Offset int) []*model.UserAPIKey

	// Alerts
	GetActiveAlerts(ctx context.Context) []interface{}
	AlertsList(ctx context.Context, Limit int, Offset int) ([]*model.Alert, int, error)
	AlertsCreate(ctx context.Context, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error
	AlertsUpdate(ctx context.Context, ID string, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error
	AlertDelete(ctx context.Context, AlertID string) error

	// Apikey
	GenerateApiKey(ctx context.Context, UserID string, KeyName string) (*model.APIKey, error)
	GetUserApiKeys(ctx context.Context, UserID string) ([]*model.APIKey, error)
	UpdateUserApiKey(ctx context.Context, UserID string, KeyID string, Active bool) ([]*model.APIKey, error)
	DeleteUserApiKey(ctx context.Context, UserID string, KeyID string) ([]*model.APIKey, error)
	GetApiKeyUser(ctx context.Context, APK string) (*model.User, error)

	// Auth
	AuthUser(ctx context.Context, UserEmail string, UserPassword string) (*model.User, string, error)
	UserResetRequest(ctx context.Context, UserEmail string) (resetID string, UserName string, resetErr error)
	UserResetPassword(ctx context.Context, ResetID string, UserPassword string) (UserName string, UserEmail string, resetErr error)
	UserUpdatePassword(ctx context.Context, UserID string, UserPassword string) (Name string, Email string, resetErr error)
	UserVerifyRequest(ctx context.Context, UserId string) (*model.User, string, error)
	VerifyUserAccount(ctx context.Context, VerifyID string) error
	MFASetupGenerate(email string) (string, string, error)
	MFASetupValidate(ctx context.Context, UserID string, secret string, passcode string) error
	MFARemove(ctx context.Context, UserID string) error
	MFATokenValidate(ctx context.Context, SessionId string, passcode string) error

	// Battles
	CreateBattle(ctx context.Context, LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*model.Plan, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string, HideVoterIdentity bool) (*model.Battle, error)
	TeamCreateBattle(ctx context.Context, TeamID string, LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*model.Plan, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string, HideVoterIdentity bool) (*model.Battle, error)
	ReviseBattle(BattleID string, BattleName string, PointValuesAllowed []string, AutoFinishVoting bool, PointAverageRounding string, HideVoterIdentity bool, JoinCode string, LeaderCode string) error
	GetBattleLeaderCode(BattleID string) (string, error)
	GetBattle(BattleID string, UserID string) (*model.Battle, error)
	GetBattlesByUser(UserID string, Limit int, Offset int) ([]*model.Battle, int, error)
	ConfirmLeader(BattleID string, UserID string) error
	GetBattleUserActiveStatus(BattleID string, UserID string) error
	GetBattleUsers(BattleID string) []*model.BattleUser
	GetBattleActiveUsers(BattleID string) []*model.BattleUser
	AddUserToBattle(BattleID string, UserID string) ([]*model.BattleUser, error)
	RetreatUser(BattleID string, UserID string) []*model.BattleUser
	AbandonBattle(BattleID string, UserID string) ([]*model.BattleUser, error)
	SetBattleLeader(BattleID string, LeaderID string) ([]string, error)
	DemoteBattleLeader(BattleID string, LeaderID string) ([]string, error)
	ToggleSpectator(BattleID string, UserID string, Spectator bool) ([]*model.BattleUser, error)
	DeleteBattle(BattleID string) error
	AddBattleLeadersByEmail(ctx context.Context, BattleID string, LeaderEmails []string) ([]string, error)
	GetBattles(Limit int, Offset int) ([]*model.Battle, int, error)
	GetActiveBattles(Limit int, Offset int) ([]*model.Battle, int, error)

	// Checkin
	CheckinList(ctx context.Context, TeamId string, Date string, TimeZone string) ([]*model.TeamCheckin, error)
	CheckinCreate(ctx context.Context, TeamId string, UserId string, Yesterday string, Today string, Blockers string, Discuss string, GoalsMet bool) error
	CheckinUpdate(ctx context.Context, CheckinId string, Yesterday string, Today string, Blockers string, Discuss string, GoalsMet bool) error
	CheckinDelete(ctx context.Context, CheckinId string) error
	CheckinComment(ctx context.Context, TeamId string, CheckinId string, UserId string, Comment string) error
	CheckinCommentEdit(ctx context.Context, TeamId string, UserId string, CommentId string, Comment string) error
	CheckinCommentDelete(ctx context.Context, CommentId string) error

	// Departments
	DepartmentUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string) (string, string, error)
	DepartmentGet(ctx context.Context, DepartmentID string) (*model.Department, error)
	OrganizationDepartmentList(ctx context.Context, OrgID string, Limit int, Offset int) []*model.Department
	DepartmentCreate(ctx context.Context, OrgID string, OrgName string) (*model.Department, error)
	DepartmentTeamList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*model.Team
	DepartmentTeamCreate(ctx context.Context, DepartmentID string, TeamName string) (*model.Team, error)
	DepartmentUserList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*model.DepartmentUser
	DepartmentAddUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error)
	DepartmentRemoveUser(ctx context.Context, DepartmentID string, UserID string) error
	DepartmentTeamUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string, TeamID string) (string, string, string, error)
	DepartmentDelete(ctx context.Context, DepartmentID string) error

	// Organization
	OrganizationGet(ctx context.Context, OrgID string) (*model.Organization, error)
	OrganizationUserRole(ctx context.Context, UserID string, OrgID string) (string, error)
	OrganizationListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*model.Organization
	OrganizationCreate(ctx context.Context, UserID string, OrgName string) (*model.Organization, error)
	OrganizationUserList(ctx context.Context, OrgID string, Limit int, Offset int) []*model.OrganizationUser
	OrganizationAddUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error)
	OrganizationRemoveUser(ctx context.Context, OrganizationID string, UserID string) error
	OrganizationTeamList(ctx context.Context, OrgID string, Limit int, Offset int) []*model.Team
	OrganizationTeamCreate(ctx context.Context, OrgID string, TeamName string) (*model.Team, error)
	OrganizationTeamUserRole(ctx context.Context, UserID string, OrgID string, TeamID string) (string, string, error)
	OrganizationDelete(ctx context.Context, OrgID string) error

	// Plans
	GetPlans(BattleID string, UserID string) []*model.Plan
	CreatePlan(BattleID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*model.Plan, error)
	ActivatePlanVoting(BattleID string, PlanID string) ([]*model.Plan, error)
	SetVote(BattleID string, UserID string, PlanID string, VoteValue string) (BattlePlans []*model.Plan, AllUsersVoted bool)
	RetractVote(BattleID string, UserID string, PlanID string) ([]*model.Plan, error)
	EndPlanVoting(BattleID string, PlanID string) ([]*model.Plan, error)
	SkipPlan(BattleID string, PlanID string) ([]*model.Plan, error)
	RevisePlan(BattleID string, PlanID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*model.Plan, error)
	BurnPlan(BattleID string, PlanID string) ([]*model.Plan, error)
	FinalizePlan(BattleID string, PlanID string, PlanPoints string) ([]*model.Plan, error)

	// Retro
	RetroCreate(OwnerID string, RetroName string, Format string, JoinCode string, FacilitatorCode string, MaxVotes int, BrainstormVisibility string) (*model.Retro, error)
	TeamRetroCreate(ctx context.Context, TeamID string, OwnerID string, RetroName string, Format string, JoinCode string, FacilitatorCode string, MaxVotes int, BrainstormVisibility string) (*model.Retro, error)
	EditRetro(RetroID string, RetroName string, JoinCode string, FacilitatorCode string, maxVotes int, brainstormVisibility string) error
	RetroGet(RetroID string, UserID string) (*model.Retro, error)
	RetroGetByUser(UserID string) ([]*model.Retro, error)
	RetroConfirmFacilitator(RetroID string, userID string) error
	RetroGetUsers(RetroID string) []*model.RetroUser
	GetRetroFacilitators(RetroID string) []string
	RetroAddUser(RetroID string, UserID string) ([]*model.RetroUser, error)
	RetroFacilitatorAdd(RetroID string, UserID string) ([]string, error)
	RetroFacilitatorRemove(RetroID string, UserID string) ([]string, error)
	RetroRetreatUser(RetroID string, UserID string) []*model.RetroUser
	RetroAbandon(RetroID string, UserID string) ([]*model.RetroUser, error)
	RetroAdvancePhase(RetroID string, Phase string) (*model.Retro, error)
	RetroDelete(RetroID string) error
	GetRetroUserActiveStatus(RetroID string, UserID string) error
	GetRetros(Limit int, Offset int) ([]*model.Retro, int, error)
	GetActiveRetros(Limit int, Offset int) ([]*model.Retro, int, error)
	GetRetroFacilitatorCode(RetroID string) (string, error)

	// Retro actions
	CreateRetroAction(RetroID string, UserID string, Content string) ([]*model.RetroAction, error)
	UpdateRetroAction(RetroID string, ActionID string, Content string, Completed bool) (Actions []*model.RetroAction, DeleteError error)
	DeleteRetroAction(RetroID string, userID string, ActionID string) ([]*model.RetroAction, error)
	GetRetroActions(RetroID string) []*model.RetroAction
	GetTeamRetroActions(TeamID string, Limit int, Offset int, Completed bool) ([]*model.RetroAction, int, error)
	RetroActionCommentAdd(RetroID string, ActionID string, UserID string, Comment string) ([]*model.RetroAction, error)
	RetroActionCommentEdit(RetroID string, ActionID string, CommentID string, Comment string) ([]*model.RetroAction, error)
	RetroActionCommentDelete(RetroID string, ActionID string, CommentID string) ([]*model.RetroAction, error)
	RetroActionAssigneeAdd(RetroID string, ActionID string, UserID string) ([]*model.RetroAction, error)
	RetroActionAssigneeDelete(RetroID string, ActionID string, UserID string) ([]*model.RetroAction, error)

	// Retro items
	CreateRetroItem(RetroID string, UserID string, ItemType string, Content string) ([]*model.RetroItem, error)
	GroupRetroItem(RetroID string, ItemId string, GroupId string) ([]*model.RetroItem, error)
	DeleteRetroItem(RetroID string, userID string, Type string, ItemID string) ([]*model.RetroItem, error)
	GetRetroItems(RetroID string) []*model.RetroItem
	GetRetroGroups(RetroID string) []*model.RetroGroup
	GroupNameChange(RetroID string, GroupId string, Name string) ([]*model.RetroGroup, error)
	GetRetroVotes(RetroID string) []*model.RetroVote
	GroupUserVote(RetroID string, GroupID string, UserID string) ([]*model.RetroVote, error)
	GroupUserSubtractVote(RetroID string, GroupID string, UserID string) ([]*model.RetroVote, error)

	// Session
	CreateSession(ctx context.Context, UserId string) (string, error)
	EnableSession(ctx context.Context, SessionId string) error
	GetSessionUser(ctx context.Context, SessionId string) (*model.User, error)
	DeleteSession(ctx context.Context, SessionId string) error

	// Storyboard
	CreateStoryboardColumn(StoryboardID string, GoalID string, userID string) ([]*model.StoryboardGoal, error)
	ReviseStoryboardColumn(StoryboardID string, UserID string, ColumnID string, ColumnName string) ([]*model.StoryboardGoal, error)
	DeleteStoryboardColumn(StoryboardID string, userID string, ColumnID string) ([]*model.StoryboardGoal, error)
	CreateStoryboardGoal(StoryboardID string, userID string, GoalName string) ([]*model.StoryboardGoal, error)
	ReviseGoalName(StoryboardID string, userID string, GoalID string, GoalName string) ([]*model.StoryboardGoal, error)
	DeleteStoryboardGoal(StoryboardID string, userID string, GoalID string) ([]*model.StoryboardGoal, error)
	GetStoryboardGoals(StoryboardID string) []*model.StoryboardGoal
	CreateStoryboardStory(StoryboardID string, GoalID string, ColumnID string, userID string) ([]*model.StoryboardGoal, error)
	ReviseStoryName(StoryboardID string, userID string, StoryID string, StoryName string) ([]*model.StoryboardGoal, error)
	ReviseStoryContent(StoryboardID string, userID string, StoryID string, StoryContent string) ([]*model.StoryboardGoal, error)
	ReviseStoryColor(StoryboardID string, userID string, StoryID string, StoryColor string) ([]*model.StoryboardGoal, error)
	ReviseStoryPoints(StoryboardID string, userID string, StoryID string, Points int) ([]*model.StoryboardGoal, error)
	ReviseStoryClosed(StoryboardID string, userID string, StoryID string, Closed bool) ([]*model.StoryboardGoal, error)
	ReviseStoryLink(StoryboardID string, userID string, StoryID string, Link string) ([]*model.StoryboardGoal, error)
	MoveStoryboardStory(StoryboardID string, userID string, StoryID string, GoalID string, ColumnID string, PlaceBefore string) ([]*model.StoryboardGoal, error)
	DeleteStoryboardStory(StoryboardID string, userID string, StoryID string) ([]*model.StoryboardGoal, error)
	AddStoryComment(StoryboardID string, UserID string, StoryID string, Comment string) ([]*model.StoryboardGoal, error)
	EditStoryComment(StoryboardID string, CommentID string, Comment string) ([]*model.StoryboardGoal, error)
	DeleteStoryComment(StoryboardID string, CommentID string) ([]*model.StoryboardGoal, error)
	CreateStoryboard(ctx context.Context, OwnerID string, StoryboardName string, JoinCode string, FacilitatorCode string) (*model.Storyboard, error)
	TeamCreateStoryboard(ctx context.Context, TeamID string, OwnerID string, StoryboardName string, JoinCode string, FacilitatorCode string) (*model.Storyboard, error)
	EditStoryboard(StoryboardID string, StoryboardName string, JoinCode string, FacilitatorCode string) error
	GetStoryboard(StoryboardID string, UserID string) (*model.Storyboard, error)
	GetStoryboardsByUser(UserID string) ([]*model.Storyboard, int, error)
	ConfirmStoryboardFacilitator(StoryboardID string, UserID string) error
	GetStoryboardUsers(StoryboardID string) []*model.StoryboardUser
	GetStoryboardPersonas(StoryboardID string) []*model.StoryboardPersona
	AddUserToStoryboard(StoryboardID string, UserID string) ([]*model.StoryboardUser, error)
	RetreatStoryboardUser(StoryboardID string, UserID string) []*model.StoryboardUser
	GetStoryboardUserActiveStatus(StoryboardID string, UserID string) error
	AbandonStoryboard(StoryboardID string, UserID string) ([]*model.StoryboardUser, error)
	SetStoryboardOwner(StoryboardID string, userID string, OwnerID string) (*model.Storyboard, error)
	StoryboardReviseColorLegend(StoryboardID string, UserID string, ColorLegend string) (*model.Storyboard, error)
	DeleteStoryboard(StoryboardID string, userID string) error
	AddStoryboardPersona(StoryboardID string, UserID string, Name string, Role string, Description string) ([]*model.StoryboardPersona, error)
	UpdateStoryboardPersona(StoryboardID string, UserID string, PersonaID string, Name string, Role string, Description string) ([]*model.StoryboardPersona, error)
	DeleteStoryboardPersona(StoryboardID string, UserID string, PersonaID string) ([]*model.StoryboardPersona, error)
	GetStoryboards(Limit int, Offset int) ([]*model.Storyboard, int, error)
	GetActiveStoryboards(Limit int, Offset int) ([]*model.Storyboard, int, error)
	StoryboardFacilitatorAdd(StoryboardId string, UserID string) (*model.Storyboard, error)
	StoryboardFacilitatorRemove(StoryboardId string, UserID string) (*model.Storyboard, error)
	GetStoryboardFacilitatorCode(StoryboardID string) (string, error)

	// Teams
	TeamUserRole(ctx context.Context, UserID string, TeamID string) (string, error)
	TeamGet(ctx context.Context, TeamID string) (*model.Team, error)
	TeamListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*model.Team
	TeamCreate(ctx context.Context, UserID string, TeamName string) (*model.Team, error)
	TeamAddUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error)
	TeamUserList(ctx context.Context, TeamID string, Limit int, Offset int) ([]*model.TeamUser, int, error)
	TeamRemoveUser(ctx context.Context, TeamID string, UserID string) error
	TeamBattleList(ctx context.Context, TeamID string, Limit int, Offset int) []*model.Battle
	TeamAddBattle(ctx context.Context, TeamID string, BattleID string) error
	TeamRemoveBattle(ctx context.Context, TeamID string, BattleID string) error
	TeamDelete(ctx context.Context, TeamID string) error
	TeamRetroList(ctx context.Context, TeamID string, Limit int, Offset int) []*model.Retro
	TeamAddRetro(ctx context.Context, TeamID string, RetroID string) error
	TeamRemoveRetro(ctx context.Context, TeamID string, RetroID string) error
	TeamStoryboardList(ctx context.Context, TeamID string, Limit int, Offset int) []*model.Storyboard
	TeamAddStoryboard(ctx context.Context, TeamID string, StoryboardID string) error
	TeamRemoveStoryboard(ctx context.Context, TeamID string, StoryboardID string) error

	// Users
	GetRegisteredUsers(ctx context.Context, Limit int, Offset int) ([]*model.User, int, error)
	GetUser(ctx context.Context, UserID string) (*model.User, error)
	GetGuestUser(ctx context.Context, UserID string) (*model.User, error)
	GetUserByEmail(ctx context.Context, UserEmail string) (*model.User, error)
	CreateUserGuest(ctx context.Context, UserName string) (*model.User, error)
	CreateUserRegistered(ctx context.Context, UserName string, UserEmail string, UserPassword string, ActiveUserID string) (NewUser *model.User, VerifyID string, SessionID string, RegisterErr error)
	CreateUser(ctx context.Context, UserName string, UserEmail string, UserPassword string) (NewUser *model.User, VerifyID string, RegisterErr error)
	UpdateUserProfile(ctx context.Context, UserID string, UserName string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string) error
	UpdateUserProfileLdap(ctx context.Context, UserID string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string) error
	UpdateUserAccount(ctx context.Context, UserID string, UserName string, UserEmail string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string) error
	DeleteUser(ctx context.Context, UserID string) error
	GetActiveCountries(ctx context.Context) ([]string, error)
	SearchRegisteredUsersByEmail(ctx context.Context, Email string, Limit int, Offset int) ([]*model.User, int, error)
}
