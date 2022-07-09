package db

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

// AuthUser authenticate the user
func (d *Database) AuthUser(UserEmail string, UserPassword string) (*model.User, string, error) {
	var user model.User
	var passHash string

	e := d.db.QueryRow(
		`SELECT id, name, email, type, password, avatar, verified, notifications_enabled, COALESCE(locale, ''), disabled, mfa_enabled FROM users WHERE email = $1`,
		UserEmail,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Type,
		&passHash,
		&user.Avatar,
		&user.Verified,
		&user.NotificationsEnabled,
		&user.Locale,
		&user.Disabled,
		&user.MFAEnabled,
	)
	if e != nil {
		d.logger.Error("Unable to auth user", zap.Error(e))
		return nil, "", errors.New("user not found")
	}

	if !comparePasswords(passHash, UserPassword) {
		return nil, "", errors.New("password invalid")
	}

	if user.Disabled {
		return nil, "", errors.New("user disabled")
	}

	// check to see if the bcrypt cost has been updated, if not do so
	if checkPasswordCost(passHash) == true {
		hashedPassword, hashErr := hashSaltPassword(UserPassword)
		if hashErr == nil {
			d.db.Exec(`call update_user_password($1, $2)`, user.Id, hashedPassword)
		}
	}

	SessionId, sessErr := d.CreateSession(user.Id)
	if sessErr != nil {
		return nil, "", sessErr
	}

	return &user, SessionId, nil
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
		d.logger.Error("Unable to reset user", zap.Error(e))
		return "", "", e
	}

	return ResetID.String, name.String, nil
}

// UserResetPassword resets the user's password to a new password
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
		d.logger.Error("Unable to get user for password reset confirmation email", zap.Error(UserErr))
		return "", "", UserErr
	}

	if _, err := d.db.Exec(
		`call reset_user_password($1, $2)`, ResetID, hashedPassword); err != nil {
		return "", "", err
	}

	return name.String, email.String, nil
}

// UserUpdatePassword updates a users password
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
		d.logger.Error("Unable to get user for password update", zap.Error(UserErr))
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

// UserVerifyRequest inserts a new user verify request
func (d *Database) UserVerifyRequest(UserId string) (*model.User, string, error) {
	var VerifyId string
	user := &model.User{
		Id: UserId,
	}

	e := d.db.QueryRow(
		`SELECT name, email FROM users WHERE id = $1`,
		user.Id,
	).Scan(
		&user.Name,
		&user.Email,
	)
	if e != nil {
		d.logger.Error("error finding user verify id", zap.Error(e))
		return nil, "", errors.New("user not found")
	}

	err := d.db.QueryRow(`
		INSERT INTO user_verify (user_id) VALUES ($1) RETURNING verify_id;
		`,
		user.Id,
	).Scan(&VerifyId)
	if err != nil {
		d.logger.Error("Unable to insert user verification", zap.Error(err))
		return nil, VerifyId, err
	}

	return user, VerifyId, nil
}

// VerifyUserAccount updates a user account verified status
func (d *Database) VerifyUserAccount(VerifyID string) error {
	if _, err := d.db.Exec(
		`call verify_user_account($1)`, VerifyID); err != nil {
		return err
	}

	return nil
}

// MFASetupGenerate generates an MFA secret and QR code image base64
func (d *Database) MFASetupGenerate(email string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Thunderdome.dev",
		AccountName: email,
	})
	if err != nil {
		return "", "", fmt.Errorf("error generating MFA TOTP key: %w", err)
	}

	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return "", "", fmt.Errorf("error converting MFA TOTP key to PNG: %w", err)
	}
	png.Encode(&buf, img)

	return key.Secret(), base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// MFASetupValidate validates the MFA secret and authenticator token
// if success enables the user mfa and stores the secret in db
func (d *Database) MFASetupValidate(UserID string, secret string, passcode string) error {
	valid := totp.Validate(passcode, secret)

	if !valid {
		return errors.New("INVALID_AUTHENTICATOR_TOKEN")
	}

	encryptedSecret, secretErr := encrypt(secret, d.config.AESHashkey)
	if secretErr != nil {
		return fmt.Errorf("error encrypting MFA secret: %w", secretErr)
	}

	if _, err := d.db.Exec(
		`call user_mfa_enable($1, $2)`, UserID, encryptedSecret); err != nil {
		return fmt.Errorf("error enabling user MFA: %w", err)
	}

	return nil
}

// MFARemove removes MFA requirement from user
func (d *Database) MFARemove(UserID string) error {
	if _, err := d.db.Exec(
		`call user_mfa_remove($1)`, UserID); err != nil {
		return fmt.Errorf("error removing user MFA: %w", err)
	}

	return nil
}

// MFATokenValidate validates the MFA secret and authenticator token for auth login
func (d *Database) MFATokenValidate(SessionId string, passcode string) error {
	var encryptedSecret string

	e := d.db.QueryRow(
		`SELECT um.secret FROM user_mfa um
 				LEFT JOIN user_session us ON us.user_id = um.user_id
 				WHERE us.session_id = $1`,
		SessionId,
	).Scan(
		&encryptedSecret,
	)
	if e != nil {
		d.logger.Error("error finding user MFA secret", zap.Error(e))
		return errors.New("user not found")
	}

	decryptedSecret, secretErr := decrypt(encryptedSecret, d.config.AESHashkey)
	if secretErr != nil {
		return errors.New("unable to decode MFA secret")
	}

	valid := totp.Validate(passcode, decryptedSecret)
	if !valid {
		return errors.New("INVALID_AUTHENTICATOR_TOKEN")
	}

	err := d.EnableSession(SessionId)
	if err != nil {
		return errors.New("unable to enable user session")
	}

	return nil
}
