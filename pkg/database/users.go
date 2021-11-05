package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
)

// GetRegisteredUsers retrieves the registered users from db
func (d *Database) GetRegisteredUsers(Limit int, Offset int) []*model.User {
	var users = make([]*model.User, 0)
	rows, err := d.db.Query(
		`
		SELECT id, name, email, type, avatar, verified, country, company, job_title
		FROM registered_users_list($1, $2);`,
		Limit,
		Offset,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w model.User

			if err := rows.Scan(
				&w.UserID,
				&w.UserName,
				&w.UserEmail,
				&w.UserType,
				&w.UserAvatar,
				&w.Verified,
				&w.Country,
				&w.Company,
				&w.JobTitle,
			); err != nil {
				log.Println(err)
			} else {
				users = append(users, &w)
			}
		}
	} else {
		log.Println(err)
	}

	return users
}

// GetUser gets a user from db by ID
func (d *Database) GetUser(UserID string) (*model.User, error) {
	var w model.User
	var UserEmail sql.NullString
	var UserCountry sql.NullString
	var UserLocale sql.NullString
	var UserCompany sql.NullString
	var UserJobTitle sql.NullString

	e := d.db.QueryRow(
		"SELECT id, name, email, type, avatar, verified, notifications_enabled, country, locale, company, job_title FROM users WHERE id = $1",
		UserID,
	).Scan(
		&w.UserID,
		&w.UserName,
		&UserEmail,
		&w.UserType,
		&w.UserAvatar,
		&w.Verified,
		&w.NotificationsEnabled,
		&UserCountry,
		&UserLocale,
		&UserCompany,
		&UserJobTitle,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("user not found")
	}

	w.UserEmail = UserEmail.String
	w.Country = UserCountry.String
	w.Locale = UserLocale.String
	w.Company = UserCompany.String
	w.JobTitle = UserJobTitle.String

	return &w, nil
}

// GetUserByEmail gets the user by email
func (d *Database) GetUserByEmail(UserEmail string) (*model.User, error) {
	var w model.User
	e := d.db.QueryRow(
		"SELECT id, name, email, type, verified FROM users WHERE email = $1",
		UserEmail,
	).Scan(
		&w.UserID,
		&w.UserName,
		&w.UserEmail,
		&w.UserType,
		&w.Verified,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("user email not found")
	}

	return &w, nil
}

// CreateUserGuest adds a new guest user to the db
func (d *Database) CreateUserGuest(UserName string) (*model.User, error) {
	var UserID string
	e := d.db.QueryRow(`INSERT INTO users (name) VALUES ($1) RETURNING id`, UserName).Scan(&UserID)
	if e != nil {
		log.Println(e)
		return nil, errors.New("unable to create new user")
	}

	return &model.User{UserID: UserID, UserName: UserName, UserAvatar: "identicon", NotificationsEnabled: true, Locale: "en"}, nil
}

// CreateUserRegistered adds a new registered user to the db
func (d *Database) CreateUserRegistered(UserName string, UserEmail string, UserPassword string, ActiveUserID string) (NewUser *model.User, VerifyID string, RegisterErr error) {
	hashedPassword, hashErr := hashSaltPassword(UserPassword)
	if hashErr != nil {
		return nil, "", hashErr
	}

	var UserID string
	var verifyID string
	UserType := "CORPORAL"
	UserAvatar := "identicon"

	if ActiveUserID != "" {
		e := d.db.QueryRow(
			`SELECT userId, verifyId FROM register_existing_user($1, $2, $3, $4, $5);`,
			ActiveUserID,
			UserName,
			UserEmail,
			hashedPassword,
			UserType,
		).Scan(&UserID, &verifyID)
		if e != nil {
			log.Println(e)
			return nil, "", errors.New("a user with that email already exists")
		}
	} else {
		e := d.db.QueryRow(
			`SELECT userId, verifyId FROM register_user($1, $2, $3, $4);`,
			UserName,
			UserEmail,
			hashedPassword,
			UserType,
		).Scan(&UserID, &verifyID)
		if e != nil {
			log.Println(e)
			return nil, "", errors.New("a user with that email already exists")
		}
	}

	return &model.User{UserID: UserID, UserName: UserName, UserEmail: UserEmail, UserType: UserType, UserAvatar: UserAvatar}, verifyID, nil
}

// UpdateUserProfile attempts to update the users profile
func (d *Database) UpdateUserProfile(UserID string, UserName string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string) error {
	if UserAvatar == "" {
		UserAvatar = "identicon"
	}
	if _, err := d.db.Exec(
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
		log.Println(err)
		return errors.New("error attempting to update users profile")
	}

	return nil
}

// UpdateUserProfile attempts to delete a user
func (d *Database) DeleteUser(UserID string) error {
	if _, err := d.db.Exec(
		`call delete_user($1);`,
		UserID,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to delete user")
	}

	return nil
}

// GetActiveCountries gets a list of user countries
func (d *Database) GetActiveCountries() ([]string, error) {
	var countries = make([]string, 0)

	rows, err := d.db.Query(`SELECT * FROM countries_active();`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var country sql.NullString
			if err := rows.Scan(
				&country,
			); err != nil {
				log.Println(err)
			} else {
				if country.String != "" {
					countries = append(countries, country.String)
				}
			}
		}
	} else {
		log.Println(err)
		return nil, errors.New("error attempting to get active countries")
	}

	return countries, nil
}
