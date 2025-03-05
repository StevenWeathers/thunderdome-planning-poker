package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// Service represents the user database service
type Service struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GetRegisteredUsers gets a list of registered users
func (d *Service) GetRegisteredUsers(ctx context.Context, limit int, offset int) ([]*thunderdome.User, int, error) {
	var users = make([]*thunderdome.User, 0)
	var count int

	err := d.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM thunderdome.users WHERE type <> 'GUEST';",
	).Scan(
		&count,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("get registered users query error", zap.Error(err))
	}

	rows, err := d.DB.QueryContext(ctx,
		`
		SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''),
		 COALESCE(u.company, ''), COALESCE(u.job_title, ''), u.disabled, COALESCE(u.picture, '')
		FROM thunderdome.users u
		WHERE u.type <> 'GUEST'
		ORDER BY u.created_date
		LIMIT $1
		OFFSET $2;`,
		limit,
		offset,
	)
	if err != nil {
		return nil, count, fmt.Errorf("get registered users query error: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var w thunderdome.User

		if err := rows.Scan(
			&w.ID,
			&w.Name,
			&w.Email,
			&w.Type,
			&w.Avatar,
			&w.Verified,
			&w.Country,
			&w.Company,
			&w.JobTitle,
			&w.Disabled,
			&w.Picture,
		); err != nil {
			d.Logger.Ctx(ctx).Error("registered_users_list query scan error", zap.Error(err))
		} else {
			w.GravatarHash = db.CreateGravatarHash(w.Email)
			users = append(users, &w)
		}
	}

	return users, count, nil
}

// GetUserByID gets a user by ID
func (d *Service) GetUserByID(ctx context.Context, userID string) (*thunderdome.User, error) {
	d.Logger.Debug("Getting user from cache", zap.String("user_id", userID))
	var user thunderdome.User

	err := d.DB.QueryRowContext(ctx,
		`SELECT id, name, COALESCE(email, ''), type, avatar, verified,
			notifications_enabled, COALESCE(country, ''), COALESCE(locale, ''), COALESCE(company, ''),
			COALESCE(job_title, ''), created_date, updated_date, last_active, disabled, theme, COALESCE(picture, '')
			FROM thunderdome.users WHERE id = $1`,
		userID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Type,
		&user.Avatar,
		&user.Verified,
		&user.NotificationsEnabled,
		&user.Country,
		&user.Locale,
		&user.Company,
		&user.JobTitle,
		&user.CreatedDate,
		&user.UpdatedDate,
		&user.LastActive,
		&user.Disabled,
		&user.Theme,
		&user.Picture,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("get_user query error", zap.Error(err),
			zap.String("UserID", userID))
		d.Logger.Debug("User not found in cache", zap.String("user_id", userID))
		return nil, fmt.Errorf("get user query error: %v", err)
	}

	if user.Email != "" {
		user.GravatarHash = db.CreateGravatarHash(user.Email)
	} else {
		user.GravatarHash = db.CreateGravatarHash(user.ID)
	}

	return &user, nil
}

// GetGuestUserByID gets a guest user by ID
func (d *Service) GetGuestUserByID(ctx context.Context, userID string) (*thunderdome.User, error) {
	var user thunderdome.User

	err := d.DB.QueryRowContext(ctx, `
SELECT id, name, COALESCE(email, ''), type, avatar, verified, notifications_enabled,
 COALESCE(country, ''), COALESCE(locale, ''), COALESCE(company, ''), COALESCE(job_title, ''),
  created_date, updated_date, last_active, theme
FROM thunderdome.users
WHERE id = $1 AND type = 'GUEST';
`,
		userID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Type,
		&user.Avatar,
		&user.Verified,
		&user.NotificationsEnabled,
		&user.Country,
		&user.Locale,
		&user.Company,
		&user.JobTitle,
		&user.CreatedDate,
		&user.UpdatedDate,
		&user.LastActive,
		&user.Theme,
	)
	if err != nil {
		return nil, fmt.Errorf("get guest user query error: %v", err)
	}

	user.GravatarHash = db.CreateGravatarHash(user.ID)

	return &user, nil
}

// GetUserByEmail gets the user by email
func (d *Service) GetUserByEmail(ctx context.Context, userEmail string) (*thunderdome.User, error) {
	var user thunderdome.User

	err := d.DB.QueryRowContext(ctx,
		`SELECT u.id, u.name, u.email, u.type, c.verified, u.disabled
				FROM thunderdome.auth_credential c
				JOIN thunderdome.users u ON c.user_id = u.id
				WHERE c.email = $1`,
		db.SanitizeEmail(userEmail),
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Type,
		&user.Verified,
		&user.Disabled,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by email query error: %w", err)
	}

	user.GravatarHash = db.CreateGravatarHash(user.Email)

	return &user, nil
}

// CreateUserGuest adds a new guest user
func (d *Service) CreateUserGuest(ctx context.Context, userName string) (*thunderdome.User, error) {
	var userID string
	err := d.DB.QueryRowContext(ctx, `INSERT INTO thunderdome.users (name) VALUES ($1) RETURNING id`, userName).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("create guest user query error: %v", err)
	}

	return &thunderdome.User{ID: userID, Name: userName, Avatar: "robohash", NotificationsEnabled: true, Locale: "en", GravatarHash: db.CreateGravatarHash(userID), Type: thunderdome.GuestUserType}, nil
}

// CreateUserRegistered adds a new registered user
func (d *Service) CreateUserRegistered(ctx context.Context, userName string, userEmail string, userPassword string, activeUserID string) (newUser *thunderdome.User, verifyID string, registerErr error) {
	hashedPassword, hashErr := db.HashSaltPassword(userPassword)
	if hashErr != nil {
		return nil, "", fmt.Errorf("create registered user hash password error: %v", hashErr)
	}

	var verificationID string
	userType := thunderdome.RegisteredUserType
	avatar := "robohash"
	sanitizedEmail := db.SanitizeEmail(userEmail)
	user := &thunderdome.User{
		Name:         userName,
		Email:        sanitizedEmail,
		Type:         userType,
		Avatar:       avatar,
		GravatarHash: db.CreateGravatarHash(userEmail),
	}

	if activeUserID != "" {
		err := d.DB.QueryRowContext(ctx,
			`SELECT userId, verifyId FROM thunderdome.user_register_existing($1, $2, $3, $4, $5);`,
			activeUserID,
			userName,
			sanitizedEmail,
			hashedPassword,
			userType,
		).Scan(&user.ID, &verificationID)
		if err != nil {
			return nil, "", fmt.Errorf("create registered user from guest query error: %v", err)
		}
	} else {
		err := d.DB.QueryRow(
			`SELECT userId, verifyId FROM thunderdome.user_register($1, $2, $3, $4);`,
			userName,
			sanitizedEmail,
			hashedPassword,
			userType,
		).Scan(&user.ID, &verificationID)
		if err != nil {
			return nil, "", fmt.Errorf("create registered user query error: %v", err)
		}
	}

	return user, verificationID, nil
}

// CreateUser adds a new registered user
func (d *Service) CreateUser(ctx context.Context, name string, email string, password string) (newUser *thunderdome.User, VerifyID string, registerErr error) {
	hashedPassword, hashErr := db.HashSaltPassword(password)
	if hashErr != nil {
		return nil, "", hashErr
	}

	var verifyID string
	userType := thunderdome.RegisteredUserType
	avatar := "robohash"
	sanitizedEmail := db.SanitizeEmail(email)
	user := &thunderdome.User{
		Name:         name,
		Email:        sanitizedEmail,
		Type:         userType,
		Avatar:       avatar,
		GravatarHash: db.CreateGravatarHash(email),
	}

	err := d.DB.QueryRowContext(ctx,
		`SELECT userId, verifyId FROM thunderdome.user_register($1, $2, $3, $4);`,
		name,
		sanitizedEmail,
		hashedPassword,
		userType,
	).Scan(&user.ID, &verifyID)
	if err != nil {
		return nil, "", fmt.Errorf("create registered user query error: %v", err)
	}

	return user, verifyID, nil
}

// UpdateUserProfile updates the users profile (excludes: email, password)
func (d *Service) UpdateUserProfile(ctx context.Context, userID string, userName string, avatar string, notificationsEnabled bool, country string, locale string, company string, jobTitle string, theme string) error {
	if avatar == "" {
		avatar = "robohash"
	}
	if theme == "" {
		theme = "auto"
	}
	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.users
		SET
			name = $2,
			avatar = $3,
			notifications_enabled = $4,
			country = $5,
			locale = $6,
			company = $7,
			job_title = $8,
			theme = $9,
			last_active = NOW(),
			updated_date = NOW()
		WHERE id = $1;`,
		userID,
		userName,
		avatar,
		notificationsEnabled,
		country,
		locale,
		company,
		jobTitle,
		theme,
	); err != nil {
		return fmt.Errorf("update user profile query error: %v", err)
	}

	return nil
}

// UpdateUserProfileLdap updates the users profile (excludes: username, email, password)
func (d *Service) UpdateUserProfileLdap(ctx context.Context, userID string, avatar string, notificationsEnabled bool, country string, locale string, company string, jobTitle string, theme string) error {
	if avatar == "" {
		avatar = "robohash"
	}
	if theme == "" {
		theme = "auto"
	}
	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.users
			SET
				avatar = $2,
				notifications_enabled = $3,
				country = $4,
				locale = $5,
				company = $6,
				job_title = $7,
				theme = $8,
				last_active = NOW(),
				updated_date = NOW()
			WHERE id = $1;`,
		userID,
		avatar,
		notificationsEnabled,
		country,
		locale,
		company,
		jobTitle,
		theme,
	); err != nil {
		return fmt.Errorf("update ldap user profile query error: %v", err)
	}

	return nil
}

// UpdateUserAccount updates the users profile including email (excludes: password)
func (d *Service) UpdateUserAccount(ctx context.Context, userID string, userName string, email string, avatar string, notificationsEnabled bool, country string, locale string, company string, jobTitle string, theme string) error {
	if avatar == "" {
		avatar = "robohash"
	}
	if theme == "" {
		theme = "auto"
	}
	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.users
			SET
				name = $2,
				email = $3,
				avatar = $4,
				notifications_enabled = $5,
				country = $6,
				locale = $7,
				company = $8,
				job_title = $9,
				theme = $10,
				last_active = NOW(),
				updated_date = NOW()
			WHERE id = $1;`,
		userID,
		userName,
		db.SanitizeEmail(email),
		avatar,
		notificationsEnabled,
		country,
		locale,
		company,
		jobTitle,
		theme,
	); err != nil {
		return fmt.Errorf("update user account query error: %v", err)
	}

	return nil
}

// DeleteUser deletes a user
func (d *Service) DeleteUser(ctx context.Context, userID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.users WHERE id = $1;`,
		userID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("delete_user query error", zap.Error(err))
		return fmt.Errorf("delete user query error: %v", err)
	}

	return nil
}

// SearchRegisteredUsersByEmail retrieves the registered users filtered by email likeness
func (d *Service) SearchRegisteredUsersByEmail(ctx context.Context, email string, limit int, offset int) ([]*thunderdome.User, int, error) {
	var users = make([]*thunderdome.User, 0)
	var count int

	rows, err := d.DB.QueryContext(ctx,
		`
		SELECT id, name, email, type, avatar, verified, country, company, job_title, count
		FROM thunderdome.users_registered_email_search($1, $2, $3);`,
		db.SanitizeEmail(email),
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("search registered users by email query error: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var user thunderdome.User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Type,
			&user.Avatar,
			&user.Verified,
			&user.Country,
			&user.Company,
			&user.JobTitle,
			&count,
		); err != nil {
			d.Logger.Ctx(ctx).Error("users_registered_email_search query error", zap.Error(err))
		} else {
			user.GravatarHash = db.CreateGravatarHash(user.Email)
			users = append(users, &user)
		}
	}

	return users, count, nil
}

// PromoteUser promotes a user to admin type
func (d *Service) PromoteUser(ctx context.Context, userID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.users SET type = 'ADMIN', updated_date = NOW() WHERE id = $1;`,
		userID,
	); err != nil {
		return fmt.Errorf("promote user to admin query error: %v", err)
	}

	return nil
}

// DemoteUser demotes a user to registered type
func (d *Service) DemoteUser(ctx context.Context, userID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.users SET type = 'REGISTERED', updated_date = NOW() WHERE id = $1;`,
		userID,
	); err != nil {
		return fmt.Errorf("demote admin user query error: %v", err)
	}

	return nil
}

// DisableUser disables a user from logging in
func (d *Service) DisableUser(ctx context.Context, userID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.user_disable($1);`,
		userID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("CALL thunderdome.user_disable error", zap.Error(err))
		return fmt.Errorf("disable user query error: %v", err)
	}

	return nil
}

// EnableUser enables a user allowing login
func (d *Service) EnableUser(ctx context.Context, userID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.users SET disabled = false, updated_date = NOW()
        WHERE id = $1;`,
		userID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("CALL thunderdome.user_enable error", zap.Error(err))
		return fmt.Errorf("enable user query error: %v", err)
	}

	return nil
}

// CleanGuests deletes guest users older than {DaysOld} days
func (d *Service) CleanGuests(ctx context.Context, daysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.users WHERE last_active < (NOW() - $1 * interval '1 day') AND type = 'GUEST';`,
		daysOld,
	); err != nil {
		return fmt.Errorf("error attempting to delete guest users older than %d days: %v", daysOld, err)
	}

	return nil
}

// GetActiveCountries gets a list of user countries
func (d *Service) GetActiveCountries(ctx context.Context) ([]string, error) {
	var countries = make([]string, 0)

	rows, err := d.DB.QueryContext(ctx, `SELECT ac.country FROM thunderdome.active_countries ac;`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var country sql.NullString
			if err := rows.Scan(
				&country,
			); err != nil {
				d.Logger.Ctx(ctx).Error("countries_active query scan error", zap.Error(err))
			} else {
				if country.String != "" {
					countries = append(countries, country.String)
				}
			}
		}
	} else {
		return nil, fmt.Errorf("get active countries query error: %v", err)
	}

	return countries, nil
}

// GetUserCredentialByUserID gets the user credential by user ID if they have one
func (d *Service) GetUserCredentialByUserID(ctx context.Context, userID string) (*thunderdome.Credential, error) {
	var c thunderdome.Credential

	err := d.DB.QueryRowContext(ctx,
		`SELECT user_id, email, verified, mfa_enabled
			FROM thunderdome.auth_credential
			WHERE user_id = $1`,
		userID,
	).Scan(
		&c.UserID,
		&c.Email,
		&c.Verified,
		&c.MFAEnabled,
	)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get user credential query error: %v", err)
	}

	return &c, nil
}
