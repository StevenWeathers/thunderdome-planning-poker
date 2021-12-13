package db

import (
	"errors"
	"github.com/StevenWeathers/thunderdome-planning-poker/model"
)

// CheckinList gets a list of team checkins by day
func (d *Database) CheckinList(TeamId string, Date string, TimeZone string) ([]*model.TeamCheckin, error) {
	Checkins := make([]*model.TeamCheckin, 0)

	rows, err := d.db.Query(`SELECT
 		tc.id, u.id, u.name, u.avatar,
 		COALESCE(tc.yesterday, ''), COALESCE(tc.today, ''),
 		COALESCE(tc.blockers, ''), coalesce(tc.discuss, ''),
 		tc.goals_met, tc.created_date, tc.updated_date
		FROM team_checkin tc
		LEFT JOIN users u ON tc.user_id = u.id
		WHERE tc.team_id = $1
		AND date(tc.created_date AT TIME ZONE $3) = $2;
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

			if err := rows.Scan(
				&checkin.Id,
				&user.Id,
				&user.Name,
				&user.Avatar,
				&checkin.Yesterday,
				&checkin.Today,
				&checkin.Blockers,
				&checkin.Discuss,
				&checkin.GoalsMet,
				&checkin.CreatedDate,
				&checkin.UpdatedDate,
			); err != nil {
				return nil, err
			} else {
				checkin.User = &user
				Checkins = append(Checkins, &checkin)
			}
		}
	}

	return Checkins, err
}

// CheckinCreate creates a team checkin
func (d *Database) CheckinCreate(
	TeamId string, UserId string,
	Yesterday string, Today string, Blockers string, Discuss string,
	GoalsMet bool,
) error {
	var userCount int
	// target user must be on team to check in
	usrErr := d.db.QueryRow(`SELECT count(user_id) FROM team_user WHERE team_id = $1 AND user_id = $2;`,
		TeamId,
		UserId,
	).Scan(&userCount)
	if usrErr != nil {
		return usrErr
	}
	if userCount != 1 {
		return errors.New("REQUIRES_TEAM_USER")
	}

	if _, err := d.db.Exec(`INSERT INTO team_checkin
		(team_id, user_id, yesterday, today, blockers, discuss, goals_met)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
		`,
		TeamId,
		UserId,
		Yesterday,
		Today,
		Blockers,
		Discuss,
		GoalsMet,
	); err != nil {
		return err
	}

	return nil
}

// CheckinUpdate updates a team checkin
func (d *Database) CheckinUpdate(
	CheckinId string,
	Yesterday string, Today string, Blockers string, Discuss string,
	GoalsMet bool,
) error {
	if _, err := d.db.Exec(`
		UPDATE team_checkin
		SET Yesterday = $2, today = $3, blockers = $4, discuss = $5, goals_met = $6
		WHERE id = $1;
		`,
		CheckinId,
		Yesterday,
		Today,
		Blockers,
		Discuss,
		GoalsMet,
	); err != nil {
		return err
	}

	return nil
}

// CheckinDelete deletes a team checkin
func (d *Database) CheckinDelete(CheckinId string) error {
	_, err := d.db.Exec(
		`DELETE FROM team_checkin WHERE id = $1;`,
		CheckinId,
	)

	if err != nil {
		return err
	}

	return nil
}
