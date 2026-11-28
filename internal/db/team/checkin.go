package team

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/jackc/pgx/v5/pgconn"
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

const errCheckinAlreadyExists = "CHECKIN_ALREADY_EXISTS"

func isCheckinAlreadyExistsError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505" && pgErr.ConstraintName == "team_checkin_team_id_user_id_checkin_date_uidx"
}

// CheckinList gets a list of team checkins by day
func (d *CheckinService) CheckinList(ctx context.Context, teamID string, date string) ([]*thunderdome.TeamCheckin, error) {
	checkins := make([]*thunderdome.TeamCheckin, 0)

	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("checkin list invalid date: %v", err)
	}

	rows, err := d.DB.QueryContext(ctx, `SELECT
 		tc.id, u.id, u.name, u.email, u.avatar, COALESCE(u.picture, ''),
 		COALESCE(tc.yesterday, ''), COALESCE(tc.today, ''),
 		COALESCE(tc.blockers, ''), coalesce(tc.discuss, ''),
		tc.goals_met, tc.checkin_date, tc.created_date, tc.updated_date,
 		COALESCE(
			json_agg(tcc ORDER BY tcc.created_date) FILTER (WHERE tcc.id IS NOT NULL), '[]'
		) AS comments
		FROM thunderdome.team_checkin tc
		LEFT JOIN thunderdome.users u ON tc.user_id = u.id
		LEFT JOIN thunderdome.team_checkin_comment tcc ON tcc.checkin_id = tc.id
		WHERE tc.team_id = $1
		AND tc.checkin_date = $2
		GROUP BY tc.id, u.id
		ORDER BY tc.created_date ASC;
		`,
		teamID,
		targetDate,
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
				&checkin.CheckinDate,
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

// CheckinLastByUser gets the last checkin by a user before the requested date
func (d *CheckinService) CheckinLastByUser(ctx context.Context, teamID string, userID string, date string) (*thunderdome.TeamCheckin, error) {
	var checkin thunderdome.TeamCheckin

	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("checkin last by user invalid date: %v", err)
	}

	err = d.DB.QueryRowContext(ctx, `SELECT
	 		tc.id, COALESCE(tc.yesterday, ''), COALESCE(tc.today, ''),
	  		COALESCE(tc.blockers, ''), coalesce(tc.discuss, ''),
			tc.goals_met, tc.checkin_date, tc.created_date, tc.updated_date
		FROM thunderdome.team_checkin tc
		WHERE tc.team_id = $1
		AND tc.user_id = $2
		AND tc.checkin_date < $3
		ORDER BY tc.checkin_date DESC, tc.created_date DESC
		LIMIT 1;
		`,
		teamID,
		userID,
		targetDate,
	).Scan(
		&checkin.ID,
		&checkin.Yesterday,
		&checkin.Today,
		&checkin.Blockers,
		&checkin.Discuss,
		&checkin.GoalsMet,
		&checkin.CheckinDate,
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
	checkinDate string,
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

	if _, err := d.DB.ExecContext(ctx, `INSERT INTO thunderdome.team_checkin
		(team_id, user_id, yesterday, today, blockers, discuss, goals_met, checkin_date)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			COALESCE($8::date, CURRENT_DATE)
		);
		`,
		teamID,
		userID,
		sanitizedYesterday,
		sanitizedToday,
		sanitizedBlockers,
		sanitizedDiscuss,
		goalsMet,
		sql.NullString{String: checkinDate, Valid: checkinDate != ""},
	); err != nil {
		if isCheckinAlreadyExistsError(err) {
			return errors.New(errCheckinAlreadyExists)
		}

		return fmt.Errorf("checkin create error: %v", err)
	}

	return nil
}

// CheckinUpdate updates a team checkin on behalf of its owner or an admin
func (d *CheckinService) CheckinUpdate(
	ctx context.Context,
	teamID string,
	userID string,
	checkinID string,
	yesterday string, today string, blockers string, discuss string,
	goalsMet bool,
) error {
	sanitizedYesterday := d.HTMLSanitizerPolicy.Sanitize(yesterday)
	sanitizedToday := d.HTMLSanitizerPolicy.Sanitize(today)
	sanitizedBlockers := d.HTMLSanitizerPolicy.Sanitize(blockers)
	sanitizedDiscuss := d.HTMLSanitizerPolicy.Sanitize(discuss)

	result, err := d.DB.ExecContext(ctx, `
		UPDATE thunderdome.team_checkin c
		SET Yesterday = $2, today = $3, blockers = $4, discuss = $5, goals_met = $6
		WHERE c.id = $1 AND c.team_id = $7 AND `+db.AuthorOrTeamAdminCondition("c.user_id", "c.team_id", "$8")+`;
		`,
		checkinID,
		sanitizedYesterday,
		sanitizedToday,
		sanitizedBlockers,
		sanitizedDiscuss,
		goalsMet,
		teamID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("checkin update query error: %v", err)
	}

	return d.requireCheckinWrite(ctx, result, teamID, checkinID)
}

// CheckinDelete deletes a team checkin on behalf of its owner or an admin
func (d *CheckinService) CheckinDelete(ctx context.Context, teamID string, userID string, checkinID string) error {
	result, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team_checkin c
		WHERE c.id = $1 AND c.team_id = $2 AND `+db.AuthorOrTeamAdminCondition("c.user_id", "c.team_id", "$3")+`;`,
		checkinID,
		teamID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("checkin delete query error: %v", err)
	}

	return d.requireCheckinWrite(ctx, result, teamID, checkinID)
}

// requireCheckinWrite explains why a check-in write matched nothing
func (d *CheckinService) requireCheckinWrite(ctx context.Context, result sql.Result, teamID string, checkinID string) error {
	return d.requireOwnedWrite(ctx, result,
		`SELECT EXISTS (SELECT 1 FROM thunderdome.team_checkin WHERE id = $1 AND team_id = $2);`,
		checkinID, teamID, errCheckinNotFound, "REQUIRES_CHECKIN_OWNER")
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

	result, err := d.DB.ExecContext(ctx, `
		INSERT INTO thunderdome.team_checkin_comment (checkin_id, user_id, comment)
		SELECT id, $2::uuid, $3 FROM thunderdome.team_checkin WHERE id = $1 AND team_id = $4;
		`,
		checkinID,
		userID,
		comment,
		teamID,
	)
	if err != nil {
		return fmt.Errorf("checkin comment query error: %v", err)
	}

	return requireRow(result, errCheckinNotFound)
}

// CheckinCommentEdit edits a team checkin comment on behalf of its author or an admin
func (d *CheckinService) CheckinCommentEdit(ctx context.Context, teamID string, userID string, commentID string, comment string) error {
	result, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.team_checkin_comment cc SET comment = $4, updated_date = NOW()
		FROM thunderdome.team_checkin c
		WHERE cc.id = $1 AND c.id = cc.checkin_id AND c.team_id = $2 AND `+db.AuthorOrTeamAdminCondition("cc.user_id", "c.team_id", "$3")+`;`,
		commentID,
		teamID,
		userID,
		comment,
	)

	if err != nil {
		return fmt.Errorf("checkin edit comment query error: %v", err)
	}

	return d.requireCommentWrite(ctx, result, teamID, commentID)
}

// CheckinCommentDelete deletes a team checkin comment on behalf of its author or an admin
func (d *CheckinService) CheckinCommentDelete(ctx context.Context, teamID string, userID string, commentID string) error {
	result, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team_checkin_comment cc
		USING thunderdome.team_checkin c
		WHERE cc.id = $1 AND c.id = cc.checkin_id AND c.team_id = $2 AND `+db.AuthorOrTeamAdminCondition("cc.user_id", "c.team_id", "$3")+`;`,
		commentID,
		teamID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("checkin delete comment query error: %v", err)
	}

	return d.requireCommentWrite(ctx, result, teamID, commentID)
}

// requireCommentWrite explains why a comment write matched nothing
func (d *CheckinService) requireCommentWrite(ctx context.Context, result sql.Result, teamID string, commentID string) error {
	return d.requireOwnedWrite(ctx, result,
		`SELECT EXISTS (
			SELECT 1 FROM thunderdome.team_checkin_comment cc
			JOIN thunderdome.team_checkin c ON c.id = cc.checkin_id
			WHERE cc.id = $1 AND c.team_id = $2
		);`,
		commentID, teamID, errCheckinCommentNotFound, "REQUIRES_COMMENT_OWNER")
}

// requireOwnedWrite tells apart a missing record (notFound) from one the user may not change (forbidden)
func (d *CheckinService) requireOwnedWrite(
	ctx context.Context, result sql.Result, existsQuery string, id string, teamID string, notFound string, forbidden string,
) error {
	writeErr := requireRow(result, notFound)
	if writeErr == nil {
		return nil
	}

	var exists bool
	if err := d.DB.QueryRowContext(ctx, existsQuery, id, teamID).Scan(&exists); err != nil {
		return fmt.Errorf("ownership lookup query error: %v", err)
	}
	if exists {
		return errors.New(forbidden)
	}

	return writeErr
}

// requireRow returns a notFound error when a write scoped to its parent matched nothing
func requireRow(result sql.Result, notFound string) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %v", err)
	}
	if rows == 0 {
		return errors.New(notFound)
	}

	return nil
}

const (
	errCheckinNotFound        = "CHECKIN_NOT_FOUND"
	errCheckinCommentNotFound = "CHECKIN_COMMENT_NOT_FOUND"
)
