package db

import (
	"errors"
	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"log"
)

// CreateSession creates a new user authenticated session
func (d *Database) CreateSession(UserId string) (string, error) {
	SessionId, err := randomBase64String(32)
	if err != nil {
		return "", err
	}

	if _, sessionErr := d.db.Exec(`
		INSERT INTO user_session (session_id, user_id) VALUES ($1, $2);
		`,
		SessionId,
		UserId,
	); sessionErr != nil {
		log.Println("Unable to create a user session: ", sessionErr)
		return "", sessionErr
	}

	return SessionId, nil
}

// GetSessionUser gets a user session by sessionId
func (d *Database) GetSessionUser(SessionId string) (*model.User, error) {
	User := &model.User{}

	e := d.db.QueryRow(`
		SELECT id, name, email, type, avatar, verified, notifications_enabled, country, locale, company, job_title, created_date, updated_date, last_active 
		FROM user_session_get($1);`,
		SessionId,
	).Scan(
		&User.Id,
		&User.Name,
		&User.Email,
		&User.Type,
		&User.Avatar,
		&User.Verified,
		&User.NotificationsEnabled,
		&User.Country,
		&User.Locale,
		&User.Company,
		&User.JobTitle,
		&User.CreatedDate,
		&User.UpdatedDate,
		&User.LastActive)
	if e != nil {
		log.Println(e)
		return nil, errors.New("active session match not found")
	}

	User.GravatarHash = createGravatarHash(User.Email)

	return User, nil
}

// DeleteSession deletes a user authenticated session
func (d *Database) DeleteSession(SessionId string) error {
	if _, sessionErr := d.db.Exec(`
		DELETE FROM user_session WHERE session_id = $1;
		`,
		SessionId,
	); sessionErr != nil {
		log.Println("Unable to delete user session: ", sessionErr)
		return sessionErr
	}

	return nil
}
