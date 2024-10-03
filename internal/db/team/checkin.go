package team

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/microcosm-cc/bluemonday"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// CheckinService represents the database service for team checkins
type CheckinService struct {
	DB                  *sql.DB
	Logger              *otelzap.Logger
	HTMLSanitizerPolicy *bluemonday.Policy
}

// CheckinList gets a list of team checkins by day
func (d *CheckinService) CheckinList(ctx context.Context, teamID string, date string, timeZone string) ([]*thunderdome.TeamCheckin, error) {
	checkins := make([]*thunderdome.TeamCheckin, 0)

	rows, err := d.DB.QueryContext(ctx, `SELECT
 		tc.id, u.id, u.name, u.email, u.avatar, COALESCE(u.picture, ''),
 		COALESCE(tc.yesterday, ''), COALESCE(tc.today, ''),
 		COALESCE(tc.blockers, ''), coalesce(tc.discuss, ''),
 		tc.goals_met, tc.created_date, tc.updated_date,
 		COALESCE(
			json_agg(tcc ORDER BY tcc.created_date) FILTER (WHERE tcc.id IS NOT NULL), '[]'
		) AS comments
		FROM thunderdome.team_checkin tc
		LEFT JOIN thunderdome.users u ON tc.user_id = u.id
		LEFT JOIN thunderdome.team_checkin_comment tcc ON tcc.checkin_id = tc.id
		WHERE tc.team_id = $1
		AND date(tc.created_date AT TIME ZONE $3) = $2
		GROUP BY tc.id, u.id;
		`,
		teamID,
		date,
		timeZone,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var checkin thunderdome.TeamCheckin
			var user thunderdome.TeamUser
			var commentsVal string

			if err := rows.Scan(
				&checkin.ID,
				&user.ID,
				&user.Name,
				&user.GravatarHash,
				&user.Avatar,
				&user.PictureURL,
				&checkin.Yesterday,
				&checkin.Today,
				&checkin.Blockers,
				&checkin.Discuss,
				&checkin.GoalsMet,
				&checkin.CreatedDate,
				&checkin.UpdatedDate,
				&commentsVal,
			); err != nil {
				return nil, err
			} else {
				user.GravatarHash = db.CreateGravatarHash(user.GravatarHash)
				checkin.User = &user

				comments := make([]*thunderdome.CheckinComment, 0)
				jsonErr := json.Unmarshal([]byte(commentsVal), &comments)
				if jsonErr != nil {
					d.Logger.Ctx(ctx).Error("checkin comments json error", zap.Error(jsonErr))
				}
				checkin.Comments = comments

				checkins = append(checkins, &checkin)
			}
		}
	}

	return checkins, err
}

// CheckinLastByUser gets the last checkin by a user
func (d *CheckinService) CheckinLastByUser(ctx context.Context, teamID string, userID string) (*thunderdome.TeamCheckin, error) {
	var checkin thunderdome.TeamCheckin

	err := d.DB.QueryRowContext(ctx, `SELECT
 		tc.id, COALESCE(tc.yesterday, ''), COALESCE(tc.today, ''),
 		COALESCE(tc.blockers, ''), coalesce(tc.discuss, ''),
 		tc.goals_met, tc.created_date, tc.updated_date
		FROM thunderdome.team_checkin tc
		WHERE tc.team_id = $1 AND tc.user_id = $2
		ORDER BY tc.created_date DESC LIMIT 1;
		`,
		teamID,
		userID,
	).Scan(
		&checkin.ID,
		&checkin.Yesterday,
		&checkin.Today,
		&checkin.Blockers,
		&checkin.Discuss,
		&checkin.GoalsMet,
		&checkin.CreatedDate,
		&checkin.UpdatedDate)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("NO_LAST_CHECKIN")
	} else if err != nil {
		return nil, err
	}

	return &checkin, err
}

// CheckinCreate creates a team checkin
func (d *CheckinService) CheckinCreate(
	ctx context.Context,
	teamID string, userID string,
	yesterday string, today string, blockers string, discuss string,
	goalsMet bool,
) error {
	var userCount int
	// target user must be on team to check in
	usrErr := d.DB.QueryRowContext(ctx, `SELECT count(user_id) FROM thunderdome.team_user WHERE team_id = $1 AND user_id = $2;`,
		teamID,
		userID,
	).Scan(&userCount)
	if usrErr != nil {
		return fmt.Errorf("checkin create get team user error: %v", usrErr)
	}
	if userCount != 1 {
		return errors.New("REQUIRES_TEAM_USER")
	}

	sanitizedYesterday := d.HTMLSanitizerPolicy.Sanitize(yesterday)
	sanitizedToday := d.HTMLSanitizerPolicy.Sanitize(today)
	sanitizedBlockers := d.HTMLSanitizerPolicy.Sanitize(blockers)
	sanitizedDiscuss := d.HTMLSanitizerPolicy.Sanitize(discuss)

	if _, err := d.DB.Exec(`INSERT INTO thunderdome.team_checkin
		(team_id, user_id, yesterday, today, blockers, discuss, goals_met)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
		`,
		teamID,
		userID,
		sanitizedYesterday,
		sanitizedToday,
		sanitizedBlockers,
		sanitizedDiscuss,
		goalsMet,
	); err != nil {
		return fmt.Errorf("checkin create error: %v", err)
	}

	return nil
}

// CheckinUpdate updates a team checkin
func (d *CheckinService) CheckinUpdate(
	ctx context.Context,
	checkinID string,
	yesterday string, today string, blockers string, discuss string,
	goalsMet bool,
) error {
	sanitizedYesterday := d.HTMLSanitizerPolicy.Sanitize(yesterday)
	sanitizedToday := d.HTMLSanitizerPolicy.Sanitize(today)
	sanitizedBlockers := d.HTMLSanitizerPolicy.Sanitize(blockers)
	sanitizedDiscuss := d.HTMLSanitizerPolicy.Sanitize(discuss)

	if _, err := d.DB.ExecContext(ctx, `
		UPDATE thunderdome.team_checkin
		SET Yesterday = $2, today = $3, blockers = $4, discuss = $5, goals_met = $6
		WHERE id = $1;
		`,
		checkinID,
		sanitizedYesterday,
		sanitizedToday,
		sanitizedBlockers,
		sanitizedDiscuss,
		goalsMet,
	); err != nil {
		return fmt.Errorf("checkin update query error: %v", err)
	}

	return nil
}

// CheckinDelete deletes a team checkin
func (d *CheckinService) CheckinDelete(ctx context.Context, checkinID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team_checkin WHERE id = $1;`,
		checkinID,
	)

	if err != nil {
		return fmt.Errorf("checkin delete query error: %v", err)
	}

	return nil
}

// CheckinComment comments on a team checkin
func (d *CheckinService) CheckinComment(
	ctx context.Context,
	teamID string,
	checkinID string,
	userID string,
	comment string,
) error {
	var userCount int
	// target user must be on team to comment on checkin
	usrErr := d.DB.QueryRowContext(ctx, `SELECT count(user_id) FROM thunderdome.team_user WHERE team_id = $1 AND user_id = $2;`,
		teamID,
		userID,
	).Scan(&userCount)
	if usrErr != nil {
		return fmt.Errorf("checkin comment get team user error: %v", usrErr)
	}
	if userCount != 1 {
		return errors.New("REQUIRES_TEAM_USER")
	}

	if _, err := d.DB.ExecContext(ctx, `
		INSERT INTO thunderdome.team_checkin_comment (checkin_id, user_id, comment) VALUES ($1, $2, $3);
		`,
		checkinID,
		userID,
		comment,
	); err != nil {
		return fmt.Errorf("checkin comment query error: %v", err)
	}

	return nil
}

// CheckinCommentEdit edits a team checkin comment
func (d *CheckinService) CheckinCommentEdit(ctx context.Context, teamID string, userID string, commentID string, comment string) error {
	var userCount int
	// target user must be on team to comment on checkin
	usrErr := d.DB.QueryRowContext(ctx, `SELECT count(user_id) FROM thunderdome.team_user WHERE team_id = $1 AND user_id = $2;`,
		teamID,
		userID,
	).Scan(&userCount)
	if usrErr != nil {
		return fmt.Errorf("checkin edit comment get team user error: %v", usrErr)
	}
	if userCount != 1 {
		return errors.New("REQUIRES_TEAM_USER")
	}

	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.team_checkin_comment SET comment = $2, updated_date = NOW() WHERE id = $1;`,
		commentID,
		comment,
	)

	if err != nil {
		return fmt.Errorf("checkin edit comment query error: %v", err)
	}

	return nil
}

// CheckinCommentDelete deletes a team checkin comment
func (d *CheckinService) CheckinCommentDelete(ctx context.Context, commentID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team_checkin_comment WHERE id = $1;`,
		commentID,
	)

	if err != nil {
		return fmt.Errorf("checkin delete comment query error: %v", err)
	}

	return nil
}
