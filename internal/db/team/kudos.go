package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// KudoList gets a list of team kudos by day.
func (d *CheckinService) KudoList(ctx context.Context, teamID string, date string) ([]*thunderdome.TeamKudo, error) {
	kudos := make([]*thunderdome.TeamKudo, 0)

	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("kudo list invalid date: %v", err)
	}

	rows, err := d.DB.QueryContext(ctx, `SELECT
		tk.id, tk.team_id,
		u.id, u.name, u.email, u.avatar, COALESCE(u.picture, ''),
		tu.id, tu.name, tu.email, tu.avatar, COALESCE(tu.picture, ''),
		COALESCE(tk.comment, ''), tk.kudos_date, tk.created_date, tk.updated_date
		FROM thunderdome.team_kudos tk
		JOIN thunderdome.users u ON tk.user_id = u.id
		JOIN thunderdome.users tu ON tk.target_user_id = tu.id
		WHERE tk.team_id = $1
		AND tk.kudos_date = $2
		ORDER BY tk.created_date ASC;`,
		teamID,
		targetDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		kudo, scanErr := scanTeamKudo(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		kudos = append(kudos, kudo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return kudos, nil
}

// KudoGet gets a single team kudo by ID.
func (d *CheckinService) KudoGet(ctx context.Context, teamID string, kudoID string) (*thunderdome.TeamKudo, error) {
	row := d.DB.QueryRowContext(ctx, `SELECT
		tk.id, tk.team_id,
		u.id, u.name, u.email, u.avatar, COALESCE(u.picture, ''),
		tu.id, tu.name, tu.email, tu.avatar, COALESCE(tu.picture, ''),
		COALESCE(tk.comment, ''), tk.kudos_date, tk.created_date, tk.updated_date
		FROM thunderdome.team_kudos tk
		JOIN thunderdome.users u ON tk.user_id = u.id
		JOIN thunderdome.users tu ON tk.target_user_id = tu.id
		WHERE tk.team_id = $1
		AND tk.id = $2;`,
		teamID,
		kudoID,
	)

	kudo, err := scanTeamKudo(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("NO_TEAM_KUDO")
		}
		return nil, fmt.Errorf("kudo get query error: %v", err)
	}

	return kudo, nil
}

// KudoCreate creates a team kudo.
func (d *CheckinService) KudoCreate(
	ctx context.Context,
	teamID string,
	userID string,
	targetUserID string,
	kudosDate string,
	comment string,
) (*thunderdome.TeamKudo, error) {
	if err := d.requireTeamUser(ctx, teamID, userID, "kudo create get submitting team user error"); err != nil {
		return nil, err
	}
	if err := d.requireTeamUser(ctx, teamID, targetUserID, "kudo create get target team user error"); err != nil {
		return nil, err
	}

	normalizedDate, err := normalizeDateInput(kudosDate)
	if err != nil {
		return nil, fmt.Errorf("kudo create invalid date: %v", err)
	}

	sanitizedComment := d.HTMLSanitizerPolicy.Sanitize(comment)

	var kudoID string
	err = d.DB.QueryRowContext(ctx, `INSERT INTO thunderdome.team_kudos
		(team_id, user_id, target_user_id, comment, kudos_date)
		VALUES ($1, $2, $3, $4, COALESCE($5::date, CURRENT_DATE))
		ON CONFLICT (team_id, user_id, target_user_id, kudos_date) DO NOTHING
		RETURNING id;`,
		teamID,
		userID,
		targetUserID,
		sanitizedComment,
		normalizedDate,
	).Scan(&kudoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("KUDO_ALREADY_EXISTS")
		}
		return nil, fmt.Errorf("kudo create query error: %v", err)
	}

	return d.KudoGet(ctx, teamID, kudoID)
}

// KudoUpdate updates a team kudo.
func (d *CheckinService) KudoUpdate(
	ctx context.Context,
	teamID string,
	kudoID string,
	targetUserID string,
	kudosDate string,
	comment string,
) (*thunderdome.TeamKudo, error) {
	var userID string
	err := d.DB.QueryRowContext(ctx,
		`SELECT user_id FROM thunderdome.team_kudos WHERE id = $1 AND team_id = $2;`,
		kudoID,
		teamID,
	).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("NO_TEAM_KUDO")
		}
		return nil, fmt.Errorf("kudo update load existing query error: %v", err)
	}

	if err := d.requireTeamUser(ctx, teamID, targetUserID, "kudo update get target team user error"); err != nil {
		return nil, err
	}

	normalizedDate, err := normalizeDateInput(kudosDate)
	if err != nil {
		return nil, fmt.Errorf("kudo update invalid date: %v", err)
	}

	var duplicateCount int
	err = d.DB.QueryRowContext(ctx, `SELECT count(id)
		FROM thunderdome.team_kudos
		WHERE team_id = $1
		AND user_id = $2
		AND target_user_id = $3
		AND kudos_date = COALESCE($4::date, CURRENT_DATE)
		AND id <> $5;`,
		teamID,
		userID,
		targetUserID,
		normalizedDate,
		kudoID,
	).Scan(&duplicateCount)
	if err != nil {
		return nil, fmt.Errorf("kudo update uniqueness query error: %v", err)
	}
	if duplicateCount > 0 {
		return nil, fmt.Errorf("KUDO_ALREADY_EXISTS")
	}

	sanitizedComment := d.HTMLSanitizerPolicy.Sanitize(comment)

	result, err := d.DB.ExecContext(ctx, `UPDATE thunderdome.team_kudos
		SET target_user_id = $3,
			comment = $4,
			kudos_date = COALESCE($5::date, CURRENT_DATE),
			updated_date = NOW()
		WHERE id = $1
		AND team_id = $2;`,
		kudoID,
		teamID,
		targetUserID,
		sanitizedComment,
		normalizedDate,
	)
	if err != nil {
		return nil, fmt.Errorf("kudo update query error: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("kudo update rows affected error: %v", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("NO_TEAM_KUDO")
	}

	return d.KudoGet(ctx, teamID, kudoID)
}

// KudoDelete deletes a team kudo.
func (d *CheckinService) KudoDelete(ctx context.Context, teamID string, kudoID string) error {
	result, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team_kudos WHERE id = $1 AND team_id = $2;`,
		kudoID,
		teamID,
	)
	if err != nil {
		return fmt.Errorf("kudo delete query error: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("kudo delete rows affected error: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("NO_TEAM_KUDO")
	}

	return nil
}

func (d *CheckinService) requireTeamUser(ctx context.Context, teamID string, userID string, queryErr string) error {
	var userCount int
	err := d.DB.QueryRowContext(ctx,
		`SELECT count(user_id) FROM thunderdome.team_user WHERE team_id = $1 AND user_id = $2;`,
		teamID,
		userID,
	).Scan(&userCount)
	if err != nil {
		return fmt.Errorf("%s: %v", queryErr, err)
	}
	if userCount != 1 {
		return errors.New("REQUIRES_TEAM_USER")
	}

	return nil
}

func normalizeDateInput(date string) (sql.NullString, error) {
	if date == "" {
		return sql.NullString{}, nil
	}

	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return sql.NullString{}, err
	}

	return sql.NullString{String: targetDate.Format("2006-01-02"), Valid: true}, nil
}

type teamKudoScanner interface {
	Scan(dest ...any) error
}

func scanTeamKudo(scanner teamKudoScanner) (*thunderdome.TeamKudo, error) {
	var kudo thunderdome.TeamKudo
	var user thunderdome.TeamUser
	var targetUser thunderdome.TeamUser

	err := scanner.Scan(
		&kudo.ID,
		&kudo.TeamID,
		&user.ID,
		&user.Name,
		&user.GravatarHash,
		&user.Avatar,
		&user.PictureURL,
		&targetUser.ID,
		&targetUser.Name,
		&targetUser.GravatarHash,
		&targetUser.Avatar,
		&targetUser.PictureURL,
		&kudo.Comment,
		&kudo.KudosDate,
		&kudo.CreatedDate,
		&kudo.UpdatedDate,
	)
	if err != nil {
		return nil, err
	}

	user.GravatarHash = db.CreateGravatarHash(user.GravatarHash)
	targetUser.GravatarHash = db.CreateGravatarHash(targetUser.GravatarHash)
	kudo.User = &user
	kudo.TargetUser = &targetUser

	return &kudo, nil
}
