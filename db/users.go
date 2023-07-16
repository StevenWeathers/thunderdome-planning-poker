package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// UserService represents a PostgreSQL implementation of thunderdome.UserDataSvc.
type UserService struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GetRegisteredUsers gets a list of registered users
func (d *UserService) GetRegisteredUsers(ctx context.Context, Limit int, Offset int) ([]*thunderdome.User, int, error) {
	var users = make([]*thunderdome.User, 0)
	var Count int

	err := d.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM thunderdome.users WHERE email IS NOT NULL;",
	).Scan(
		&Count,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("get registered users query error", zap.Error(err))
	}

	rows, err := d.DB.QueryContext(ctx,
		`
		SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''), COALESCE(u.company, ''), COALESCE(u.job_title, ''), u.disabled
		FROM thunderdome.users u
		WHERE u.email IS NOT NULL
		ORDER BY u.created_date
		LIMIT $1
		OFFSET $2;`,
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
			FROM thunderdome.users WHERE id = $1`,
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
FROM thunderdome.users
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
		"SELECT id, name, email, type, verified, disabled FROM thunderdome.users WHERE LOWER(email) = $1",
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
	err := d.DB.QueryRowContext(ctx, `INSERT INTO thunderdome.users (name) VALUES ($1) RETURNING id`, UserName).Scan(&UserID)
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
			`SELECT userId, verifyId FROM thunderdome.user_register_existing($1, $2, $3, $4, $5);`,
			ActiveUserID,
			UserName,
			sanitizedEmail,
			hashedPassword,
			UserType,
		).Scan(&User.Id, &verifyID)
		if err != nil {
			d.Logger.Ctx(ctx).Error("user_register_existing query error", zap.Error(err))
			return nil, "", errors.New("a user with that email already exists")
		}
	} else {
		err := d.DB.QueryRow(
			`SELECT userId, verifyId FROM thunderdome.user_register($1, $2, $3, $4);`,
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
		`SELECT userId, verifyId FROM thunderdome.user_register($1, $2, $3, $4);`,
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
		`UPDATE thunderdome.users
		SET
			name = $2,
			avatar = $3,
			notifications_enabled = $4,
			country = $5,
			locale = $6,
			company = $7,
			job_title = $8,
			last_active = NOW(),
			updated_date = NOW()
		WHERE id = $1;`,
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
		`UPDATE thunderdome.users
			SET
				avatar = $2,
				notifications_enabled = $3,
				country = $4,
				locale = $5,
				company = $6,
				job_title = $7,
				last_active = NOW(),
				updated_date = NOW()
			WHERE id = $1;`,
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
				last_active = NOW(),
				updated_date = NOW()
			WHERE id = $1;`,
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
		`DELETE FROM thunderdome.users WHERE id = $1;`,
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
		FROM thunderdome.users_registered_email_search($1, $2, $3);`,
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
			d.Logger.Ctx(ctx).Error("users_registered_email_search query error", zap.Error(err))
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
		`UPDATE thunderdome.users SET type = 'ADMIN', updated_date = NOW() WHERE id = $1;`,
		UserID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("CALL thunderdome.promote_user error", zap.Error(err))
		return errors.New("error attempting to promote user to admin")
	}

	return nil
}

// DemoteUser demotes a user to registered type
func (d *UserService) DemoteUser(ctx context.Context, UserID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.users SET type = 'REGISTERED', updated_date = NOW() WHERE id = $1;`,
		UserID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("CALL thunderdome.demote_user error", zap.Error(err))
		return errors.New("error attempting to demote user to registered")
	}

	return nil
}

// DisableUser disables a user from logging in
func (d *UserService) DisableUser(ctx context.Context, UserID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.user_disable($1);`,
		UserID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("CALL thunderdome.user_disable error", zap.Error(err))
		return errors.New("error attempting to disable user")
	}

	return nil
}

// EnableUser enables a user allowing login
func (d *UserService) EnableUser(ctx context.Context, UserID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.users SET disabled = false, updated_date = NOW()
        WHERE id = $1;`,
		UserID,
	); err != nil {
		d.Logger.Ctx(ctx).Error("CALL thunderdome.user_enable error", zap.Error(err))
		return errors.New("error attempting to enable user")
	}

	return nil
}

// CleanGuests deletes guest users older than {DaysOld} days
func (d *UserService) CleanGuests(ctx context.Context, DaysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.users WHERE last_active < (NOW() - $1 * interval '1 day') AND type = 'GUEST';`,
		DaysOld,
	); err != nil {
		d.Logger.Ctx(ctx).Error("CALL thunderdome.clean_guest_users", zap.Error(err))
		return errors.New("error attempting to clean Guest Users")
	}

	return nil
}

// LowercaseUserEmails goes through and lower cases any user email that has uppercase letters
// returning the list of updated users
func (d *UserService) LowercaseUserEmails(ctx context.Context) ([]*thunderdome.User, error) {
	var users = make([]*thunderdome.User, 0)
	rows, err := d.DB.QueryContext(ctx,
		`UPDATE thunderdome.users u
        SET email = lower(u.email), updated_date = NOW()
        FROM (
            SELECT lower(su.email) AS email
            FROM thunderdome.users su
            WHERE su.email IS NOT NULL
            GROUP BY lower(su.email) HAVING count(su.*) = 1
        ) AS subquery
        WHERE lower(u.email) = subquery.email AND u.email ~ '[A-Z]' RETURNING u.name, u.email;`,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr thunderdome.User

			if err := rows.Scan(
				&usr.Name,
				&usr.Email,
			); err != nil {
				d.Logger.Ctx(ctx).Error("lowercase_unique_user_emails scan error", zap.Error(err))
				return nil, err
			} else {
				users = append(users, &usr)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("lowercase_unique_user_emails query error", zap.Error(err))
		return nil, err
	}

	return users, nil
}

// MergeDuplicateAccounts goes through and merges user accounts with duplicate emails that has uppercase letters
// returning the list of merged users
func (d *UserService) MergeDuplicateAccounts(ctx context.Context) ([]*thunderdome.User, error) {
	var users = make([]*thunderdome.User, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT name, email FROM thunderdome.user_merge_nonunique_accounts();`,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr thunderdome.User

			if err := rows.Scan(
				&usr.Name,
				&usr.Email,
			); err != nil {
				d.Logger.Ctx(ctx).Error("merge_nonunique_user_accounts scan error", zap.Error(err))
				return nil, err
			} else {
				users = append(users, &usr)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("merge_nonunique_user_accounts query error", zap.Error(err))
		return nil, err
	}

	return users, nil
}

// GetActiveCountries gets a list of user countries
func (d *UserService) GetActiveCountries(ctx context.Context) ([]string, error) {
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
		d.Logger.Ctx(ctx).Error("countries_active query error", zap.Error(err))
		return nil, errors.New("error attempting to get active countries")
	}

	return countries, nil
}
