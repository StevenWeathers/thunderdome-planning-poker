package database

import (
	"errors"
	"log"
)

// ConfirmAdmin confirms whether the users is infact an admin
func (d *Database) ConfirmAdmin(UserID string) error {
	var UserType string
	e := d.db.QueryRow("SELECT coalesce(type, '') FROM users WHERE id = $1;", UserID).Scan(&UserType)
	if e != nil {
		log.Println(e)
		return errors.New("could not find users type")
	}

	if UserType != "GENERAL" {
		return errors.New(("user is not an admin"))
	}

	return nil
}

// ApplicationStats includes user, battle, and plan counts
type ApplicationStats struct {
	RegisteredCount   int `json:"registeredUserCount"`
	UnregisteredCount int `json:"unregisteredUserCount"`
	BattleCount       int `json:"battleCount"`
	PlanCount         int `json:"planCount"`
}

// GetAppStats gets counts of users (registered and unregistered), battles, and plans
func (d *Database) GetAppStats() (*ApplicationStats, error) {
	var Appstats ApplicationStats

	err := d.db.QueryRow(`
		SELECT
			unregistered_user_count,
			registered_user_count,
			battle_count,
			plan_count
		FROM get_app_stats();
		`,
	).Scan(
		&Appstats.UnregisteredCount,
		&Appstats.RegisteredCount,
		&Appstats.BattleCount,
		&Appstats.PlanCount,
	)
	if err != nil {
		log.Println("Unable to get application stats: ", err)
		return nil, err
	}

	return &Appstats, nil
}

// PromoteUser promotes a user to admin type
func (d *Database) PromoteUser(UserID string) error {
	if _, err := d.db.Exec(
		`call promote_user($1);`,
		UserID,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to promote user to admin")
	}

	return nil
}

// DemoteUser demotes a user to registered type
func (d *Database) DemoteUser(UserID string) error {
	if _, err := d.db.Exec(
		`call demote_user($1);`,
		UserID,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to demote user to registered")
	}

	return nil
}

// CleanBattles deletes battles older than X days
func (d *Database) CleanBattles(DaysOld int) error {
	if _, err := d.db.Exec(
		`call clean_battles($1);`,
		DaysOld,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to clean battles")
	}

	return nil
}

// CleanGuests deletes guest users older than X days
func (d *Database) CleanGuests(DaysOld int) error {
	if _, err := d.db.Exec(
		`call clean_guest_users($1);`,
		DaysOld,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to clean Guest Users")
	}

	return nil
}
