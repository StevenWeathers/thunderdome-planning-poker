package db

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

// AuthService represents a PostgreSQL implementation of thunderdome.AuthService.
type AuthService struct {
	DB         *sql.DB
	Logger     *otelzap.Logger
	AESHashkey string
}

// AuthUser authenticate the user
func (d *AuthService) AuthUser(ctx context.Context, UserEmail string, UserPassword string) (*thunderdome.User, string, error) {
	var user thunderdome.User
	var passHash string
	sanitizedEmail := sanitizeEmail(UserEmail)

	err := d.DB.QueryRowContext(ctx,
		`SELECT id, name, email, type, password, avatar, verified, notifications_enabled, COALESCE(locale, ''), disabled, mfa_enabled FROM thunderdome.users WHERE LOWER(email) = $1`,
		sanitizedEmail,
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			d.Logger.Ctx(ctx).Error("Unable to auth user not found", zap.Error(err), zap.String("email", sanitizedEmail))
			return nil, "", errors.New("USER_NOT_FOUND")
		} else {
			return nil, "", err
		}
	}

	if !comparePasswords(passHash, UserPassword) {
		return nil, "", errors.New("INVALID_PASSWORD")
	}

	if user.Disabled {
		return nil, "", errors.New("USER_DISABLED")
	}

	// check to see if the bcrypt cost has been updated, if not do so
	if checkPasswordCost(passHash) {
		hashedPassword, hashErr := hashSaltPassword(UserPassword)
		if hashErr == nil {
			_, updateErr := d.DB.Exec(`UPDATE thunderdome.users SET password = $2, last_active = NOW(), updated_date = NOW() WHERE id = $1;`, user.Id, hashedPassword)
			if updateErr != nil {
				d.Logger.Error("Unable to update password cost", zap.Error(updateErr), zap.String("email", sanitizedEmail))
			}
		}
	}

	SessionId, sessErr := d.CreateSession(ctx, user.Id)
	if sessErr != nil {
		return nil, "", sessErr
	}

	return &user, SessionId, nil
}

// UserResetRequest inserts a new user reset request
func (d *AuthService) UserResetRequest(ctx context.Context, UserEmail string) (resetID string, UserName string, resetErr error) {
	var ResetID sql.NullString
	var UserID sql.NullString
	var name sql.NullString

	e := d.DB.QueryRowContext(ctx, `
		SELECT resetId, userId, userName FROM thunderdome.user_reset_create($1);
		`,
		sanitizeEmail(UserEmail),
	).Scan(&ResetID, &UserID, &name)
	if e != nil {
		d.Logger.Ctx(ctx).Error("Unable to reset user", zap.Error(e), zap.String("email", UserEmail))
		return "", "", e
	}

	return ResetID.String, name.String, nil
}

// UserResetPassword resets the user's password to a new password
func (d *AuthService) UserResetPassword(ctx context.Context, ResetID string, UserPassword string) (UserName string, UserEmail string, resetErr error) {
	var name sql.NullString
	var email sql.NullString

	hashedPassword, hashErr := hashSaltPassword(UserPassword)
	if hashErr != nil {
		return "", "", hashErr
	}

	UserErr := d.DB.QueryRowContext(ctx, `
		SELECT
			w.name, w.email
		FROM thunderdome.user_reset wr
		LEFT JOIN thunderdome.users w ON w.id = wr.user_id
		WHERE wr.reset_id = $1;
		`,
		ResetID,
	).Scan(&name, &email)
	if UserErr != nil {
		d.Logger.Ctx(ctx).Error("Unable to get user for password reset confirmation email", zap.Error(UserErr))
		return "", "", UserErr
	}

	if _, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.user_password_reset($1, $2)`,
		ResetID, hashedPassword); err != nil {
		return "", "", err
	}

	return name.String, email.String, nil
}

// UserUpdatePassword updates a users password
func (d *AuthService) UserUpdatePassword(ctx context.Context, UserID string, UserPassword string) (Name string, Email string, resetErr error) {
	var UserName sql.NullString
	var UserEmail sql.NullString

	UserErr := d.DB.QueryRowContext(ctx, `
		SELECT
			w.name, w.email
		FROM thunderdome.users w
		WHERE w.id = $1;
		`,
		UserID,
	).Scan(&UserName, &UserEmail)
	if UserErr != nil {
		d.Logger.Ctx(ctx).Error("Unable to get user for password update", zap.Error(UserErr))
		return "", "", UserErr
	}

	hashedPassword, hashErr := hashSaltPassword(UserPassword)
	if hashErr != nil {
		return "", "", hashErr
	}

	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.users SET password = $2, last_active = NOW(), updated_date = NOW() WHERE id = $1;`,
		UserID, hashedPassword); err != nil {
		return "", "", err
	}

	return UserName.String, UserEmail.String, nil
}

// UserVerifyRequest inserts a new user verify request
func (d *AuthService) UserVerifyRequest(ctx context.Context, UserId string) (*thunderdome.User, string, error) {
	var VerifyId string
	user := &thunderdome.User{
		Id: UserId,
	}

	e := d.DB.QueryRowContext(ctx,
		`SELECT name, email FROM thunderdome.users WHERE id = $1`,
		user.Id,
	).Scan(
		&user.Name,
		&user.Email,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("error finding user verify id", zap.Error(e))
		return nil, "", errors.New("user not found")
	}

	err := d.DB.QueryRowContext(ctx, `
		INSERT INTO thunderdome.user_verify (user_id) VALUES ($1) RETURNING verify_id;
		`,
		user.Id,
	).Scan(&VerifyId)
	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to insert user verification", zap.Error(err))
		return nil, VerifyId, err
	}

	return user, VerifyId, nil
}

// VerifyUserAccount updates a user account verified status
func (d *AuthService) VerifyUserAccount(ctx context.Context, VerifyID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.user_account_verify($1)`, VerifyID); err != nil {
		return err
	}

	return nil
}

// MFASetupGenerate generates an MFA secret and QR code image base64
func (d *AuthService) MFASetupGenerate(email string) (string, string, error) {
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
	err = png.Encode(&buf, img)
	if err != nil {
		return "", "", fmt.Errorf("error encoding MFA PNG: %w", err)
	}

	return key.Secret(), base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// MFASetupValidate validates the MFA secret and authenticator token
// if success enables the user mfa and stores the secret in db
func (d *AuthService) MFASetupValidate(ctx context.Context, UserID string, secret string, passcode string) error {
	if passcode == "" || secret == "" {
		return errors.New("MISSING_SECRET_OR_PASSCODE")
	}
	valid := totp.Validate(passcode, secret)

	if !valid {
		return errors.New("INVALID_AUTHENTICATOR_TOKEN")
	}

	encryptedSecret, secretErr := encrypt(secret, d.AESHashkey)
	if secretErr != nil {
		return fmt.Errorf("error encrypting MFA secret: %w", secretErr)
	}

	if _, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.user_mfa_enable($1, $2)`, UserID, encryptedSecret); err != nil {
		return fmt.Errorf("error enabling user MFA: %w", err)
	}

	return nil
}

// MFARemove removes MFA requirement from user
func (d *AuthService) MFARemove(ctx context.Context, UserID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.user_mfa_remove($1)`, UserID); err != nil {
		return fmt.Errorf("error removing user MFA: %w", err)
	}

	return nil
}

// MFATokenValidate validates the MFA secret and authenticator token for auth login
func (d *AuthService) MFATokenValidate(ctx context.Context, SessionId string, passcode string) error {
	var encryptedSecret string

	e := d.DB.QueryRowContext(ctx,
		`SELECT COALESCE(um.secret, '') FROM thunderdome.user_mfa um
 				LEFT JOIN thunderdome.user_session us ON us.user_id = um.user_id
 				WHERE us.session_id = $1`,
		SessionId,
	).Scan(
		&encryptedSecret,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("error finding user MFA secret", zap.Error(e))
		return errors.New("user not found")
	}

	if encryptedSecret == "" {
		return errors.New("no secret to validate against")
	}
	decryptedSecret, secretErr := decrypt(encryptedSecret, d.AESHashkey)
	if secretErr != nil {
		return errors.New("unable to decode MFA secret")
	}

	valid := totp.Validate(passcode, decryptedSecret)
	if !valid {
		return errors.New("INVALID_AUTHENTICATOR_TOKEN")
	}

	err := d.EnableSession(ctx, SessionId)
	if err != nil {
		return errors.New("unable to enable user session")
	}

	return nil
}
