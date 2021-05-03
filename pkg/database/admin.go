package database

import (
	"errors"
	"log"
)

// ConfirmAdmin confirms whether the warrior is infact a GENERAL (ADMIN)
func (d *Database) ConfirmAdmin(AdminID string) error {
	var warriorRank string
	e := d.db.QueryRow("SELECT coalesce(rank, '') FROM warriors WHERE id = $1;", AdminID).Scan(&warriorRank)
	if e != nil {
		log.Println(e)
		return errors.New("could not find warriors rank")
	}

	if warriorRank != "GENERAL" {
		return errors.New(("warrior is not an admin"))
	}

	return nil
}

// ApplicationStats includes warrior, battle, and plan counts
type ApplicationStats struct {
	RegisteredCount   int `json:"registeredWarriorCount"`
	UnregisteredCount int `json:"unregisteredWarriorCount"`
	BattleCount       int `json:"battleCount"`
	PlanCount         int `json:"planCount"`
}

// GetAppStats gets counts of warriors (registered and unregistered), battles, and plans
func (d *Database) GetAppStats() (*ApplicationStats, error) {
	var Appstats ApplicationStats

	statsErr := d.db.QueryRow(`
		SELECT
			unregistered_warrior_count,
			registered_warrior_count,
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
	if statsErr != nil {
		log.Println("Unable to get application stats: ", statsErr)
		return nil, statsErr
	}

	return &Appstats, nil
}

// PromoteUser promotes a warrior to GENERAL (ADMIN) rank
func (d *Database) PromoteUser(WarriorID string) error {
	if _, err := d.db.Exec(
		`call promote_warrior($1);`,
		WarriorID,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to promote warrior to GENERAL")
	}

	return nil
}

// DemoteUser demotes a warrior to CORPORAL (Registered) rank
func (d *Database) DemoteUser(WarriorID string) error {
	if _, err := d.db.Exec(
		`call demote_warrior($1);`,
		WarriorID,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to demote warrior to CORPORAL")
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

// CleanGuests deletes guest warriors older than X days
func (d *Database) CleanGuests(DaysOld int) error {
	if _, err := d.db.Exec(
		`call clean_guest_warriors($1);`,
		DaysOld,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to clean Guest Warriors")
	}

	return nil
}
