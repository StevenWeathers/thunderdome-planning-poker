package db

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// CheckinList gets a list of team checkins by day
func (d *Database) CheckinList(ctx context.Context, TeamId string, Date string, TimeZone string) ([]*model.TeamCheckin, error) {
	Checkins := make([]*model.TeamCheckin, 0)

	rows, err := d.db.QueryContext(ctx, `SELECT
 		tc.id, u.id, u.name, u.email, u.avatar,
 		COALESCE(tc.yesterday, ''), COALESCE(tc.today, ''),
 		COALESCE(tc.blockers, ''), coalesce(tc.discuss, ''),
 		tc.goals_met, tc.created_date, tc.updated_date,
 		COALESCE(
			json_agg(tcc ORDER BY tcc.created_date) FILTER (WHERE tcc.id IS NOT NULL), '[]'
		) AS comments
		FROM team_checkin tc
		LEFT JOIN users u ON tc.user_id = u.id
		LEFT JOIN team_checkin_comment tcc ON tcc.checkin_id = tc.id
		WHERE tc.team_id = $1
		AND date(tc.created_date AT TIME ZONE $3) = $2
		GROUP BY tc.id, u.id;
		`,
		TeamId,
		Date,
		TimeZone,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var checkin model.TeamCheckin
			var user model.TeamUser
			var comments string

			if err := rows.Scan(
				&checkin.Id,
				&user.Id,
				&user.Name,
				&user.GravatarHash,
				&user.Avatar,
				&checkin.Yesterday,
				&checkin.Today,
				&checkin.Blockers,
				&checkin.Discuss,
				&checkin.GoalsMet,
				&checkin.CreatedDate,
				&checkin.UpdatedDate,
				&comments,
			); err != nil {
				return nil, err
			} else {
				user.GravatarHash = createGravatarHash(user.GravatarHash)
				checkin.User = &user

				Comments := make([]*model.CheckinComment, 0)
				jsonErr := json.Unmarshal([]byte(comments), &Comments)
				if jsonErr != nil {
					d.logger.Ctx(ctx).Error("checkin comments json error", zap.Error(jsonErr))
				}
				checkin.Comments = Comments

				Checkins = append(Checkins, &checkin)
			}
		}
	}

	return Checkins, err
}

// CheckinCreate creates a team checkin
func (d *Database) CheckinCreate(
	ctx context.Context,
	TeamId string, UserId string,
	Yesterday string, Today string, Blockers string, Discuss string,
	GoalsMet bool,
) error {
	var userCount int
	// target user must be on team to check in
	usrErr := d.db.QueryRowContext(ctx, `SELECT count(user_id) FROM team_user WHERE team_id = $1 AND user_id = $2;`,
		TeamId,
		UserId,
	).Scan(&userCount)
	if usrErr != nil {
		return usrErr
	}
	if userCount != 1 {
		return errors.New("REQUIRES_TEAM_USER")
	}

	SanitizedYesterday := d.htmlSanitizerPolicy.Sanitize(Yesterday)
	SanitizedToday := d.htmlSanitizerPolicy.Sanitize(Today)
	SanitizedBlockers := d.htmlSanitizerPolicy.Sanitize(Blockers)
	SanitizedDiscuss := d.htmlSanitizerPolicy.Sanitize(Discuss)

	if _, err := d.db.Exec(`INSERT INTO team_checkin
		(team_id, user_id, yesterday, today, blockers, discuss, goals_met)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
		`,
		TeamId,
		UserId,
		SanitizedYesterday,
		SanitizedToday,
		SanitizedBlockers,
		SanitizedDiscuss,
		GoalsMet,
	); err != nil {
		return err
	}

	return nil
}

// CheckinUpdate updates a team checkin
func (d *Database) CheckinUpdate(
	ctx context.Context,
	CheckinId string,
	Yesterday string, Today string, Blockers string, Discuss string,
	GoalsMet bool,
) error {
	SanitizedYesterday := d.htmlSanitizerPolicy.Sanitize(Yesterday)
	SanitizedToday := d.htmlSanitizerPolicy.Sanitize(Today)
	SanitizedBlockers := d.htmlSanitizerPolicy.Sanitize(Blockers)
	SanitizedDiscuss := d.htmlSanitizerPolicy.Sanitize(Discuss)

	if _, err := d.db.ExecContext(ctx, `
		UPDATE team_checkin
		SET Yesterday = $2, today = $3, blockers = $4, discuss = $5, goals_met = $6
		WHERE id = $1;
		`,
		CheckinId,
		SanitizedYesterday,
		SanitizedToday,
		SanitizedBlockers,
		SanitizedDiscuss,
		GoalsMet,
	); err != nil {
		return err
	}

	return nil
}

// CheckinDelete deletes a team checkin
func (d *Database) CheckinDelete(ctx context.Context, CheckinId string) error {
	_, err := d.db.ExecContext(ctx,
		`DELETE FROM team_checkin WHERE id = $1;`,
		CheckinId,
	)

	if err != nil {
		return err
	}

	return nil
}

// CheckinComment comments on a team checkin
func (d *Database) CheckinComment(
	ctx context.Context,
	TeamId string,
	CheckinId string,
	UserId string,
	Comment string,
) error {
	var userCount int
	// target user must be on team to comment on checkin
	usrErr := d.db.QueryRowContext(ctx, `SELECT count(user_id) FROM team_user WHERE team_id = $1 AND user_id = $2;`,
		TeamId,
		UserId,
	).Scan(&userCount)
	if usrErr != nil {
		return usrErr
	}
	if userCount != 1 {
		return errors.New("REQUIRES_TEAM_USER")
	}

	if _, err := d.db.ExecContext(ctx, `
		INSERT INTO team_checkin_comment (checkin_id, user_id, comment) VALUES ($1, $2, $3);
		`,
		CheckinId,
		UserId,
		Comment,
	); err != nil {
		return err
	}

	return nil
}

// CheckinCommentEdit edits a team checkin comment
func (d *Database) CheckinCommentEdit(ctx context.Context, TeamId string, UserId string, CommentId string, Comment string) error {
	var userCount int
	// target user must be on team to comment on checkin
	usrErr := d.db.QueryRowContext(ctx, `SELECT count(user_id) FROM team_user WHERE team_id = $1 AND user_id = $2;`,
		TeamId,
		UserId,
	).Scan(&userCount)
	if usrErr != nil {
		return usrErr
	}
	if userCount != 1 {
		return errors.New("REQUIRES_TEAM_USER")
	}

	_, err := d.db.ExecContext(ctx,
		`UPDATE team_checkin_comment SET comment = $2, updated_date = NOW() WHERE id = $1;`,
		CommentId,
		Comment,
	)

	if err != nil {
		return err
	}

	return nil
}

// CheckinCommentDelete deletes a team checkin comment
func (d *Database) CheckinCommentDelete(ctx context.Context, CommentId string) error {
	_, err := d.db.ExecContext(ctx,
		`DELETE FROM team_checkin_comment WHERE id = $1;`,
		CommentId,
	)

	if err != nil {
		return err
	}

	return nil
}
