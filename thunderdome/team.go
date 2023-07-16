package thunderdome

import (
	"context"
	"time"
)

type Team struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

type TeamUser struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Avatar       string `json:"avatar"`
	GravatarHash string `json:"gravatarHash"`
}

type TeamDataSvc interface {
	TeamUserRole(ctx context.Context, UserID string, TeamID string) (string, error)
	TeamGet(ctx context.Context, TeamID string) (*Team, error)
	TeamListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*Team
	TeamCreate(ctx context.Context, UserID string, TeamName string) (*Team, error)
	TeamAddUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error)
	TeamUserList(ctx context.Context, TeamID string, Limit int, Offset int) ([]*TeamUser, int, error)
	TeamRemoveUser(ctx context.Context, TeamID string, UserID string) error
	TeamBattleList(ctx context.Context, TeamID string, Limit int, Offset int) []*Poker
	TeamAddBattle(ctx context.Context, TeamID string, BattleID string) error
	TeamRemoveBattle(ctx context.Context, TeamID string, BattleID string) error
	TeamDelete(ctx context.Context, TeamID string) error
	TeamRetroList(ctx context.Context, TeamID string, Limit int, Offset int) []*Retro
	TeamAddRetro(ctx context.Context, TeamID string, RetroID string) error
	TeamRemoveRetro(ctx context.Context, TeamID string, RetroID string) error
	TeamStoryboardList(ctx context.Context, TeamID string, Limit int, Offset int) []*Storyboard
	TeamAddStoryboard(ctx context.Context, TeamID string, StoryboardID string) error
	TeamRemoveStoryboard(ctx context.Context, TeamID string, StoryboardID string) error
	TeamList(ctx context.Context, Limit int, Offset int) ([]*Team, int)
}
