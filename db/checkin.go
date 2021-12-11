package db

import (
	"errors"
	"log"

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
				&checkin.Today,
				&checkin.Yesterday,
				&checkin.Blockers,
				&checkin.Discuss,
				&checkin.GoalsMet,
				&checkin.CreatedDate,
				&checkin.UpdatedDate,
			); err != nil {
				log.Println(err)
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
		log.Println(err)
		return errors.New("error attempting to checkin")
	}

	return nil
}
