package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"golang.org/x/crypto/bcrypt"
)

// AuthUser attempts to authenticate the user
func (d *Database) AuthUser(UserEmail string, UserPassword string) (*model.User, error) {
	var w model.User
	var passHash string
	var UserLocale sql.NullString

	e := d.db.QueryRow(
		`SELECT id, name, email, type, password, avatar, verified, notifications_enabled, locale FROM users WHERE email = $1`,
		UserEmail,
	).Scan(
		&w.UserID,
		&w.UserName,
		&w.UserEmail,
		&w.UserType,
		&passHash,
		&w.UserAvatar,
		&w.Verified,
		&w.NotificationsEnabled,
		&UserLocale,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("user not found")
	}

	if !comparePasswords(passHash, UserPassword) {
		return nil, errors.New("password invalid")
	}

	// check to see if the bcrypt cost has been updated, if not do so
	if checkPasswordCost(passHash) == true {
		hashedPassword, hashErr := hashSaltPassword(UserPassword)
		if hashErr == nil {
			d.db.Exec(`call update_user_password($1, $2)`, w.UserID, hashedPassword)
		}
	}

	w.Locale = UserLocale.String

	return &w, nil
}

// UserResetRequest inserts a new user reset request
func (d *Database) UserResetRequest(UserEmail string) (resetID string, UserName string, resetErr error) {
	var ResetID sql.NullString
	var UserID sql.NullString
	var name sql.NullString

	e := d.db.QueryRow(`
		SELECT resetId, userId, userName FROM insert_user_reset($1);
		`,
		UserEmail,
	).Scan(&ResetID, &UserID, &name)
	if e != nil {
		log.Println("Unable to reset user: ", e)
		return "", "", e
	}

	return ResetID.String, name.String, nil
}

// UserResetPassword attempts to reset a users password
func (d *Database) UserResetPassword(ResetID string, UserPassword string) (UserName string, UserEmail string, resetErr error) {
	var name sql.NullString
	var email sql.NullString

	hashedPassword, hashErr := hashSaltPassword(UserPassword)
	if hashErr != nil {
		return "", "", hashErr
	}

	UserErr := d.db.QueryRow(`
		SELECT
			w.name, w.email
		FROM user_reset wr
		LEFT JOIN users w ON w.id = wr.user_id
		WHERE wr.reset_id = $1;
		`,
		ResetID,
	).Scan(&name, &email)
	if UserErr != nil {
		log.Println("Unable to get user for password reset confirmation email: ", UserErr)
		return "", "", UserErr
	}

	if _, err := d.db.Exec(
		`call reset_user_password($1, $2)`, ResetID, hashedPassword); err != nil {
		return "", "", err
	}

	return name.String, email.String, nil
}

// UserUpdatePassword attempts to update a users password
func (d *Database) UserUpdatePassword(UserID string, UserPassword string) (Name string, Email string, resetErr error) {
	var UserName sql.NullString
	var UserEmail sql.NullString

	UserErr := d.db.QueryRow(`
		SELECT
			w.name, w.email
		FROM users w
		WHERE w.id = $1;
		`,
		UserID,
	).Scan(&UserName, &UserEmail)
	if UserErr != nil {
		log.Println("Unable to get user for password update: ", UserErr)
		return "", "", UserErr
	}

	hashedPassword, hashErr := hashSaltPassword(UserPassword)
	if hashErr != nil {
		return "", "", hashErr
	}

	if _, err := d.db.Exec(
		`call update_user_password($1, $2)`, UserID, hashedPassword); err != nil {
		return "", "", err
	}

	return UserName.String, UserEmail.String, nil
}

// VerifyUserAccount attempts to verify a users account email
func (d *Database) VerifyUserAccount(VerifyID string) error {
	if _, err := d.db.Exec(
		`call verify_user_account($1)`, VerifyID); err != nil {
		return err
	}

	return nil
}

// hashSaltPassword takes a password byte and salt + hashes it
// returning a hash string to store in db
func hashSaltPassword(UserPassword string) (string, error) {
	pwd := []byte(UserPassword)
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

// comparePasswords takes a password hash and compares it to entered password
// returning true if matches false if not
func comparePasswords(hashedPwd string, password string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	SubmittedPassword := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, SubmittedPassword)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// checkPasswordCost checks the passwords stored hash for bcrypt cost
// if it does not match current cost then return true and let auth update the hash
func checkPasswordCost(hashedPwd string) bool {
	byteHash := []byte(hashedPwd)

	hashCost, costErr := bcrypt.Cost(byteHash)
	if costErr == nil && hashCost != bcrypt.DefaultCost {
		return true
	}

	return false
}
