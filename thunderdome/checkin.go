package thunderdome

import "context"

type TeamCheckin struct {
	Id          string            `json:"id"`
	User        *TeamUser         `json:"user"`
	Yesterday   string            `json:"yesterday"`
	Today       string            `json:"today"`
	Blockers    string            `json:"blockers"`
	Discuss     string            `json:"discuss"`
	GoalsMet    bool              `json:"goalsMet"`
	CreatedDate string            `json:"createdDate"`
	UpdatedDate string            `json:"updatedDate"`
	Comments    []*CheckinComment `json:"comments"`
}

// CheckinComment A checkin comment by a user
type CheckinComment struct {
	ID          string `json:"id"`
	CheckinID   string `json:"checkin_id"`
	UserID      string `json:"user_id"`
	Comment     string `json:"comment"`
	CreateDate  string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

type CheckinService interface {
	CheckinList(ctx context.Context, TeamId string, Date string, TimeZone string) ([]*TeamCheckin, error)
	CheckinCreate(ctx context.Context, TeamId string, UserId string, Yesterday string, Today string, Blockers string, Discuss string, GoalsMet bool) error
	CheckinUpdate(ctx context.Context, CheckinId string, Yesterday string, Today string, Blockers string, Discuss string, GoalsMet bool) error
	CheckinDelete(ctx context.Context, CheckinId string) error
	CheckinComment(ctx context.Context, TeamId string, CheckinId string, UserId string, Comment string) error
	CheckinCommentEdit(ctx context.Context, TeamId string, UserId string, CommentId string, Comment string) error
	CheckinCommentDelete(ctx context.Context, CommentId string) error
}
