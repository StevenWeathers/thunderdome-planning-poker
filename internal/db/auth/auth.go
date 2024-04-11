package auth

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

// Service represents a PostgreSQL implementation of thunderdome.AuthDataSvc.
type Service struct {
	DB         *sql.DB
	Logger     *otelzap.Logger
	AESHashkey string
}

// AuthUser authenticate the user
func (d *Service) AuthUser(ctx context.Context, UserEmail string, UserPassword string) (*thunderdome.User, *thunderdome.Credential, string, error) {
	var user thunderdome.User
	var cred thunderdome.Credential
	var passHash string
	sanitizedEmail := db.SanitizeEmail(UserEmail)

	err := d.DB.QueryRowContext(ctx,
		`SELECT u.id, u.name, c.email, u.type, c.password, u.avatar, c.verified, u.notifications_enabled,
 			COALESCE(u.locale, ''), u.disabled, c.mfa_enabled, u.theme, COALESCE(u.picture, '')
			FROM thunderdome.auth_credential c
			JOIN thunderdome.users u ON c.user_id = u.id 
			WHERE c.email = $1`,
		sanitizedEmail,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Type,
		&passHash,
		&user.Avatar,
		&cred.Verified,
		&user.NotificationsEnabled,
		&user.Locale,
		&user.Disabled,
		&cred.MFAEnabled,
		&user.Theme,
		&user.Picture,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			d.Logger.Ctx(ctx).Error("Unable to auth user not found", zap.Error(err), zap.String("email", sanitizedEmail))
			return nil, nil, "", errors.New("USER_NOT_FOUND")
		} else {
			return nil, nil, "", err
		}
	}

	if !db.ComparePasswords(passHash, UserPassword) {
		return nil, nil, "", errors.New("INVALID_PASSWORD")
	}

	if user.Disabled {
		return nil, nil, "", errors.New("USER_DISABLED")
	}

	// check to see if the bcrypt cost has been updated, if not do so
	if db.CheckPasswordCost(passHash) {
		hashedPassword, hashErr := db.HashSaltPassword(UserPassword)
		if hashErr == nil {
			_, updateErr := d.DB.Exec(`UPDATE thunderdome.auth_credential SET password = $2, updated_date = NOW() WHERE user_id = $1;`, user.Id, hashedPassword)
			if updateErr != nil {
				d.Logger.Error("Unable to update password cost", zap.Error(updateErr), zap.String("email", sanitizedEmail))
			}
		}
	}

	SessionId, sessErr := d.CreateSession(ctx, user.Id, !cred.MFAEnabled)
	if sessErr != nil {
		return nil, nil, "", sessErr
	}

	return &user, &cred, SessionId, nil
}

// UserResetRequest inserts a new user reset request
func (d *Service) UserResetRequest(ctx context.Context, UserEmail string) (resetID string, UserName string, resetErr error) {
	var ResetID sql.NullString
	var UserID sql.NullString
	var name sql.NullString

	err := d.DB.QueryRowContext(ctx, `
		SELECT resetId, userId, userName FROM thunderdome.user_reset_create($1);
		`,
		db.SanitizeEmail(UserEmail),
	).Scan(&ResetID, &UserID, &name)
	if err != nil {
		return "", "", fmt.Errorf("insert user reset request query error: %v", err)
	}

	return ResetID.String, name.String, nil
}

// UserResetPassword resets the user's password to a new password
func (d *Service) UserResetPassword(ctx context.Context, ResetID string, UserPassword string) (UserName string, UserEmail string, resetErr error) {
	var name sql.NullString
	var email sql.NullString

	hashedPassword, hashErr := db.HashSaltPassword(UserPassword)
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
		return "", "", fmt.Errorf("get user for password reset confirmation email error: %v", UserErr)
	}

	if _, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.user_password_reset($1, $2)`,
		ResetID, hashedPassword); err != nil {
		return "", "", fmt.Errorf("user password reset query error: %v", err)
	}

	return name.String, email.String, nil
}

// UserUpdatePassword updates a users password
func (d *Service) UserUpdatePassword(ctx context.Context, UserID string, UserPassword string) (Name string, Email string, resetErr error) {
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
		return "", "", fmt.Errorf("get user for password update query error: %v", UserErr)
	}

	hashedPassword, hashErr := db.HashSaltPassword(UserPassword)
	if hashErr != nil {
		return "", "", fmt.Errorf("hash password for password update error: %v", hashErr)
	}

	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.auth_credential SET password = $2, updated_date = NOW() WHERE user_id = $1;`,
		UserID, hashedPassword); err != nil {
		return "", "", fmt.Errorf("update password query error: %v", err)
	}

	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1;`,
		UserID, hashedPassword); err != nil {
		return "", "", fmt.Errorf("update user last_active query error: %v", err)
	}

	return UserName.String, UserEmail.String, nil
}

// UserVerifyRequest inserts a new user verify request
func (d *Service) UserVerifyRequest(ctx context.Context, UserId string) (*thunderdome.User, string, error) {
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
		return nil, "", fmt.Errorf("find user for verify request error: %v", e)
	}

	err := d.DB.QueryRowContext(ctx, `
		INSERT INTO thunderdome.user_verify (user_id) VALUES ($1) RETURNING verify_id;
		`,
		user.Id,
	).Scan(&VerifyId)
	if err != nil {
		return nil, VerifyId, fmt.Errorf("create user verify query error: %v", err)
	}

	return user, VerifyId, nil
}

// VerifyUserAccount updates a user account verified status
func (d *Service) VerifyUserAccount(ctx context.Context, VerifyID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.user_account_verify($1)`, VerifyID); err != nil {
		return fmt.Errorf("verify user acocunt query error: %v", err)
	}

	return nil
}

// MFASetupGenerate generates an MFA secret and QR code image base64
func (d *Service) MFASetupGenerate(email string) (string, string, error) {
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
func (d *Service) MFASetupValidate(ctx context.Context, UserID string, secret string, passcode string) error {
	if passcode == "" || secret == "" {
		return errors.New("MISSING_SECRET_OR_PASSCODE")
	}
	valid := totp.Validate(passcode, secret)

	if !valid {
		return errors.New("INVALID_AUTHENTICATOR_TOKEN")
	}

	encryptedSecret, secretErr := db.Encrypt(secret, d.AESHashkey)
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
func (d *Service) MFARemove(ctx context.Context, UserID string) error {
	if _, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.user_mfa_remove($1)`, UserID); err != nil {
		return fmt.Errorf("error removing user MFA: %w", err)
	}

	return nil
}

// MFATokenValidate validates the MFA secret and authenticator token for auth login
func (d *Service) MFATokenValidate(ctx context.Context, SessionId string, passcode string) error {
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
		return fmt.Errorf("find user mfa secret query error: %v", e)
	}

	if encryptedSecret == "" {
		return errors.New("no secret to validate against")
	}
	decryptedSecret, secretErr := db.Decrypt(encryptedSecret, d.AESHashkey)
	if secretErr != nil {
		return fmt.Errorf("unable to decode MFA secret: %v", secretErr)
	}

	valid := totp.Validate(passcode, decryptedSecret)
	if !valid {
		return errors.New("INVALID_AUTHENTICATOR_TOKEN")
	}

	err := d.EnableSession(ctx, SessionId)
	if err != nil {
		return fmt.Errorf("unable to enable user session: %v", err)
	}

	return nil
}
