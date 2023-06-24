package thunderdome

import (
	"context"
	"time"
)

// User aka user
type User struct {
	Id                   string    `json:"id"`
	Name                 string    `json:"name"`
	Email                string    `json:"email"`
	Type                 string    `json:"rank"`
	Avatar               string    `json:"avatar"`
	Verified             bool      `json:"verified"`
	NotificationsEnabled bool      `json:"notificationsEnabled"`
	Country              string    `json:"country"`
	Locale               string    `json:"locale"`
	Company              string    `json:"company"`
	JobTitle             string    `json:"jobTitle"`
	GravatarHash         string    `json:"gravatarHash"`
	CreatedDate          time.Time `json:"createdDate"`
	UpdatedDate          time.Time `json:"updatedDate"`
	LastActive           time.Time `json:"lastActive"`
	Disabled             bool      `json:"disabled"`
	MFAEnabled           bool      `json:"mfaEnabled"`
}

type UserService interface {
	GetUser(ctx context.Context, UserID string) (*User, error)
	GetGuestUser(ctx context.Context, UserID string) (*User, error)
	GetUserByEmail(ctx context.Context, UserEmail string) (*User, error)
	GetRegisteredUsers(ctx context.Context, Limit int, Offset int) ([]*User, int, error)
	SearchRegisteredUsersByEmail(ctx context.Context, Email string, Limit int, Offset int) ([]*User, int, error)
	CreateUser(ctx context.Context, UserName string, UserEmail string, UserPassword string) (NewUser *User, VerifyID string, RegisterErr error)
	CreateUserGuest(ctx context.Context, UserName string) (*User, error)
	CreateUserRegistered(ctx context.Context, UserName string, UserEmail string, UserPassword string, ActiveUserID string) (NewUser *User, VerifyID string, RegisterErr error)
	UpdateUserAccount(ctx context.Context, UserID string, UserName string, UserEmail string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string) error
	UpdateUserProfile(ctx context.Context, UserID string, UserName string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string) error
	UpdateUserProfileLdap(ctx context.Context, UserID string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string) error
	PromoteUser(ctx context.Context, UserID string) error
	DemoteUser(ctx context.Context, UserID string) error
	DisableUser(ctx context.Context, UserID string) error
	EnableUser(ctx context.Context, UserID string) error
	DeleteUser(ctx context.Context, UserID string) error
	CleanGuests(ctx context.Context, DaysOld int) error
	MergeDuplicateAccounts(ctx context.Context) ([]*User, error)
	LowercaseUserEmails(ctx context.Context) ([]*User, error)
}
