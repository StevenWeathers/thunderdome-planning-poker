package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// UserService represents a PostgreSQL implementation of thunderdome.UserService.
type UserService struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GetRegisteredUsers gets a list of registered users
func (d *UserService) GetRegisteredUsers(ctx context.Context, Limit int, Offset int) ([]*thunderdome.User, int, error) {
	var users = make([]*thunderdome.User, 0)
	var Count int

	err := d.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM users WHERE email IS NOT NULL;",
	).Scan(
		&Count,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("get registered users query error", zap.Error(err))
	}

	rows, err := d.DB.QueryContext(ctx,
		`
		SELECT id, name, email, type, avatar, verified, country, company, job_title, disabled
		FROM registered_users_list($1, $2);`,
		Limit,
		Offset,
	)
	if err != nil {
		return nil, Count, err
	}

	defer rows.Close()
	for rows.Next() {
		var w thunderdome.User

		if err := rows.Scan(
			&w.Id,
			&w.Name,
			&w.Email,
			&w.Type,
			&w.Avatar,
			&w.Verified,
			&w.Country,
			&w.Company,
			&w.JobTitle,
			&w.Disabled,
		); err != nil {
			d.Logger.Ctx(ctx).Error("registered_users_list query scan error", zap.Error(err))
		} else {
			w.GravatarHash = createGravatarHash(w.Email)
			users = append(users, &w)
		}
	}

	return users, Count, nil
}

// GetUser gets a user by ID
func (d *UserService) GetUser(ctx context.Context, UserID string) (*thunderdome.User, error) {
	var w thunderdome.User
	var UserEmail sql.NullString
	var UserCountry sql.NullString
	var UserLocale sql.NullString
	var UserCompany sql.NullString
	var UserJobTitle sql.NullString

	err := d.DB.QueryRowContext(ctx,
		`SELECT id, name, email, type, avatar, verified,
			notifications_enabled, country, locale, company, job_title,
			created_date, updated_date, last_active, disabled, mfa_enabled
			FROM users WHERE id = $1`,
		UserID,
	).Scan(
		&w.Id,
		&w.Name,
		&UserEmail,
		&w.Type,
		&w.Avatar,
		&w.Verified,
		&w.NotificationsEnabled,
		&UserCountry,
		&UserLocale,
		&UserCompany,
		&UserJobTitle,
		&w.CreatedDate,
		&w.UpdatedDate,
		&w.LastActive,
		&w.Disabled,
		&w.MFAEnabled,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("get user query error", zap.Error(err))
		return nil, errors.New("user not found")
	}

	w.Email = UserEmail.String
	w.Country = UserCountry.String
	w.Locale = UserLocale.String
	w.Company = UserCompany.String
	w.JobTitle = UserJobTitle.String
	if w.Email != "" {
		w.GravatarHash = createGravatarHash(w.Email)
	} else {
		w.GravatarHash = createGravatarHash(w.Id)
	}

	return &w, nil
}

// GetGuestUser gets a guest user by ID
func (d *UserService) GetGuestUser(ctx context.Context, UserID string) (*thunderdome.User, error) {
	var w thunderdome.User
	var UserEmail sql.NullString
	var UserCountry sql.NullString
	var UserLocale sql.NullString
	var UserCompany sql.NullString
	var UserJobTitle sql.NullString

	err := d.DB.QueryRowContext(ctx, `
SELECT id, name, email, type, avatar, verified, notifications_enabled, country, locale, company, job_title, created_date, updated_date, last_active
FROM users
WHERE id = $1 AND type = 'GUEST';
`,
		UserID,
	).Scan(
		&w.Id,
		&w.Name,
		&UserEmail,
		&w.Type,
		&w.Avatar,
		&w.Verified,
		&w.NotificationsEnabled,
		&UserCountry,
		&UserLocale,
		&UserCompany,
		&UserJobTitle,
		&w.CreatedDate,
		&w.UpdatedDate,
		&w.LastActive,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("get guest user query error", zap.Error(err))
		return nil, errors.New("user not found")
	}

	w.Email = UserEmail.String
	w.Country = UserCountry.String
	w.Locale = UserLocale.String
	w.Company = UserCompany.String
	w.JobTitle = UserJobTitle.String
	w.GravatarHash = createGravatarHash(w.Id)

	return &w, nil
}

// GetUserByEmail gets the user by email
func (d *UserService) GetUserByEmail(ctx context.Context, UserEmail string) (*thunderdome.User, error) {
	var w thunderdome.User

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, name, email, type, verified, disabled FROM users WHERE LOWER(email) = $1",
		sanitizeEmail(UserEmail),
	).Scan(
		&w.Id,
		&w.Name,
		&w.Email,
		&w.Type,
		&w.Verified,
		&w.Disabled,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("get user by email query error", zap.Error(err))
		return nil, errors.New("user email not found")
	}

	w.GravatarHash = createGravatarHash(w.Email)

	return &w, nil
}

// CreateUserGuest adds a new guest user
func (d *UserService) CreateUserGuest(ctx context.Context, UserName string) (*thunderdome.User, error) {
	var UserID string
	err := d.DB.QueryRowContext(ctx, `INSERT INTO users (name) VALUES ($1) RETURNING id`, UserName).Scan(&UserID)
	if err != nil {
		d.Logger.Ctx(ctx).Error("create guest user query error", zap.Error(err))
		return nil, errors.New("unable to create new user")
	}

	return &thunderdome.User{Id: UserID, Name: UserName, Avatar: "robohash", NotificationsEnabled: true, Locale: "en", GravatarHash: createGravatarHash(UserID)}, nil
}

// CreateUserRegistered adds a new registered user
func (d *UserService) CreateUserRegistered(ctx context.Context, UserName string, UserEmail string, UserPassword string, ActiveUserID string) (NewUser *thunderdome.User, VerifyID string, RegisterErr error) {
	hashedPassword, hashErr := hashSaltPassword(UserPassword)
	if hashErr != nil {
		return nil, "", hashErr
	}

	var verifyID string
	UserType := "REGISTERED"
	UserAvatar := "robohash"
	sanitizedEmail := sanitizeEmail(UserEmail)
	User := &thunderdome.User{
		Name:         UserName,
		Email:        sanitizedEmail,
		Type:         UserType,
		Avatar:       UserAvatar,
		GravatarHash: createGravatarHash(UserEmail),
	}

	if ActiveUserID != "" {
		err := d.DB.QueryRowContext(ctx,
			`SELECT userId, verifyId FROM register_existing_user($1, $2, $3, $4, $5);`,
			ActiveUserID,
			UserName,
			sanitizedEmail,
			hashedPassword,
			UserType,
		).Scan(&User.Id, &verifyID)
		if err != nil {
			d.Logger.Ctx(ctx).Error("register_existing_user query error", zap.Error(err))
			return nil, "", errors.New("a user with that email already exists")
		}
	} else {
		err := d.DB.QueryRow(
			`SELECT userId, verifyId FROM register_user($1, $2, $3, $4);`,
			UserName,
			sanitizedEmail,
			hashedPassword,
			UserType,
		).Scan(&User.Id, &verifyID)
		if err != nil {
			d.Logger.Ctx(ctx).Error("register_user query error", zap.Error(err))
			return nil, "", errors.New("a user with that email already exists")
		}
	}

	return User, verifyID, nil
}

// CreateUser adds a new registered user
func (d *UserService) CreateUser(ctx context.Context, UserName string, UserEmail string, UserPassword string) (NewUser *thunderdome.User, VerifyID string, RegisterErr error) {
	hashedPassword, hashErr := hashSaltPassword(UserPassword)
	if hashErr != nil {
		return nil, "", hashErr
	}

	var verifyID string
	UserType := "REGISTERED"
	UserAvatar := "robohash"
	sanitizedEmail := sanitizeEmail(UserEmail)
	User := &thunderdome.User{
		Name:         UserName,
		Email:        sanitizedEmail,
		Type:         UserType,
		Avatar:       UserAvatar,
		GravatarHash: createGravatarHash(UserEmail),
	}

	err := d.DB.QueryRowContext(ctx,
		`SELECT userId, verifyId FROM register_user($1, $2, $3, $4);`,
		UserName,
		sanitizedEmail,
		hashedPassword,
		UserType,
	).Scan(&User.Id, &verifyID)
	if err != nil {
		d.Logger.Ctx(ctx).Error("register_user query error", zap.Error(err))
		return nil, "", errors.New("a user with that email already exists")
	}

	return User, verifyID, nil
}

// UpdateUserProfile updates the users profile (excludes: email, password)
func (d *UserService) UpdateUserProfile(ctx context.Context, UserID string, UserName string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string) error {
	if UserAvatar == "" {
		UserAvatar = "robohash"
	}
	if _, err := d.DB.ExecContext(ctx,
		`call user_profile_update($1, $2, $3, $4, $5, $6, $7, $8);`,
		UserID,
		UserName,
		UserAvatar,
		NotificationsEnabled,
		Country,
		Locale,
		Company,
		JobTitle,
	); err != nil {
		d.Logger.Ctx(ctx).Error("user_profile_update query error", zap.Error(err))
		return errors.New("error attempting to update users profile")
	}

	return nil
}

// UpdateUserProfileLdap updates the users profile (excludes: username, email, password)
func (d *UserService) UpdateUserProfileLdap(ctx context.Context, UserID string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string) error {
	if UserAvatar == "" {
		UserAvatar = "robohash"
	}
	if _, err := d.DB.ExecContext(ctx,
		`call user_profile_ldap_update($1, $2, $3, $4, $5, $6, $7);`,
		UserID,
		UserAvatar,
		NotificationsEnabled,
		Country,
		Locale,
		Company,
		JobTitle,
	); err != nil {
		d.Logger.Ctx(ctx).Error("user_profile_ldap_update query error", zap.Error(err))
		return errors.New("error attempting to update users profile")
	}

	return nil
}

// UpdateUserAccount updates the users profile including email (excludes: password)
func (d *UserService) UpdateUserAccount(ctx context.Context, UserID string, UserName string, UserEmail string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string) error {
	if UserAvatar == "" {
		UserAvatar = "robohash"
	}
	if _, err := d.DB.ExecContext(ctx,
		`call user_account_update($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		UserID,
		UserName,
		sanitizeEmail(UserEmail),
		UserAvatar,
		NotificationsEnabled,
		Country,
		Locale,
		Company,
		JobTitle,
	); err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user
func (d *UserService) DeleteUser(ctx context.Context, UserID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`call delete_user($1);`,
		UserID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("delete_user query error", zap.Error(err))
		return errors.New("error attempting to delete user")
	}

	return nil
}

// SearchRegisteredUsersByEmail retrieves the registered users filtered by email likeness
func (d *UserService) SearchRegisteredUsersByEmail(ctx context.Context, Email string, Limit int, Offset int) ([]*thunderdome.User, int, error) {
	var users = make([]*thunderdome.User, 0)
	var count int

	rows, err := d.DB.QueryContext(ctx,
		`
		SELECT id, name, email, type, avatar, verified, country, company, job_title, count
		FROM registered_users_email_search($1, $2, $3);`,
		sanitizeEmail(Email),
		Limit,
		Offset,
	)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	for rows.Next() {
		var w thunderdome.User

		if err := rows.Scan(
			&w.Id,
			&w.Name,
			&w.Email,
			&w.Type,
			&w.Avatar,
			&w.Verified,
			&w.Country,
			&w.Company,
			&w.JobTitle,
			&count,
		); err != nil {
			d.Logger.Ctx(ctx).Error("registered_users_email_search query error", zap.Error(err))
		} else {
			w.GravatarHash = createGravatarHash(w.Email)
			users = append(users, &w)
		}
	}

	return users, count, nil
}

// PromoteUser promotes a user to admin type
func (d *UserService) PromoteUser(ctx context.Context, UserID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`call promote_user($1);`,
		UserID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("call promote_user error", zap.Error(err))
		return errors.New("error attempting to promote user to admin")
	}

	return nil
}

// DemoteUser demotes a user to registered type
func (d *UserService) DemoteUser(ctx context.Context, UserID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`call demote_user($1);`,
		UserID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("call demote_user error", zap.Error(err))
		return errors.New("error attempting to demote user to registered")
	}

	return nil
}

// DisableUser disables a user from logging in
func (d *UserService) DisableUser(ctx context.Context, UserID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`call user_disable($1);`,
		UserID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("call user_disable error", zap.Error(err))
		return errors.New("error attempting to disable user")
	}

	return nil
}

// EnableUser enables a user allowing login
func (d *UserService) EnableUser(ctx context.Context, UserID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`call user_enable($1);`,
		UserID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("call user_enable error", zap.Error(err))
		return errors.New("error attempting to enable user")
	}

	return nil
}
