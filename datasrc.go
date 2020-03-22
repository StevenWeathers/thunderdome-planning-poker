package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
)

var db *sql.DB

// Battle aka arena
type Battle struct {
	BattleID           string     `json:"id"`
	LeaderID           string     `json:"leaderId"`
	BattleName         string     `json:"name"`
	Warriors           []*Warrior `json:"warriors"`
	Plans              []*Plan    `json:"plans"`
	VotingLocked       bool       `json:"votingLocked"`
	ActivePlanID       string     `json:"activePlanId"`
	PointValuesAllowed []string   `json:"pointValuesAllowed"`
}

// Warrior aka user
type Warrior struct {
	WarriorID    string `json:"id"`
	WarriorName  string `json:"name"`
	WarriorEmail string `json:"email"`
	WarriorRank  string `json:"rank"`
	Active       bool   `json:"active"`
}

// Vote structure
type Vote struct {
	WarriorID string `json:"warriorId"`
	VoteValue string `json:"vote"`
}

// Plan aka Story structure
type Plan struct {
	PlanID        string    `json:"id"`
	PlanName      string    `json:"name"`
	Votes         []*Vote   `json:"votes"`
	Points        string    `json:"points"`
	PlanActive    bool      `json:"active"`
	PlanSkipped   bool      `json:"skipped"`
	VoteStartTime time.Time `json:"voteStartTime"`
	VoteEndTime   time.Time `json:"voteEndTime"`
}

// SetupDB runs db migrations, sets up a db connection pool
// and sets previously active warriors to false during startup
func SetupDB() {
	var (
		host     = GetEnv("DB_HOST", "db")
		port     = GetIntEnv("DB_PORT", 5432)
		user     = GetEnv("DB_USER", "thor")
		password = GetEnv("DB_PASS", "odinson")
		dbname   = GetEnv("DB_NAME", "thunderdome")
	)

	sqlFile, ioErr := pkger.Open("/schema.sql")
	if ioErr != nil {
		log.Println("Error reading schema.sql file required to migrate db")
		log.Fatal(ioErr)
	}
	sqlContent, ioErr := ioutil.ReadAll(sqlFile)
	if ioErr != nil {
		// this will hopefully only possibly panic during development as the file is already in memory otherwise
		log.Println("Error reading schema.sql file required to migrate db")
		log.Fatal(ioErr)
	}
	migrationSQL := string(sqlContent)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	if _, err := db.Exec(migrationSQL); err != nil {
		log.Fatal(err)
	}

	// on server start reset all warriors to active false for battles
	if _, err := db.Exec(
		`call deactivate_all_warriors();`); err != nil {
		log.Println(err)
	}
}

//CreateBattle adds a new battle to the db
func CreateBattle(LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*Plan) (*Battle, error) {
	var pointValuesJSON, _ = json.Marshal(PointValuesAllowed)

	var b = &Battle{
		BattleID:           "",
		LeaderID:           LeaderID,
		BattleName:         BattleName,
		Warriors:           make([]*Warrior, 0),
		Plans:              make([]*Plan, 0),
		VotingLocked:       true,
		ActivePlanID:       "",
		PointValuesAllowed: PointValuesAllowed,
	}

	e := db.QueryRow(
		`INSERT INTO battles (leader_id, name, point_values_allowed) VALUES ($1, $2, $3) RETURNING id`,
		LeaderID,
		BattleName,
		string(pointValuesJSON),
	).Scan(&b.BattleID)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Error Creating Battle")
	}

	for _, plan := range Plans {
		plan.Votes = make([]*Vote, 0)

		e := db.QueryRow(
			`INSERT INTO plans (battle_id, name) VALUES ($1, $2) RETURNING id`,
			b.BattleID,
			plan.PlanName,
		).Scan(&plan.PlanID)
		if e != nil {
			log.Println(e)
		}
	}

	b.Plans = Plans

	return b, nil
}

// GetBattle gets a battle by ID
func GetBattle(BattleID string) (*Battle, error) {
	var b = &Battle{
		BattleID:           BattleID,
		LeaderID:           "",
		BattleName:         "",
		Warriors:           make([]*Warrior, 0),
		Plans:              make([]*Plan, 0),
		VotingLocked:       true,
		ActivePlanID:       "",
		PointValuesAllowed: make([]string, 0),
	}

	// get battle
	var ActivePlanID sql.NullString
	var pv string
	e := db.QueryRow(
		"SELECT id, name, leader_id, voting_locked, active_plan_id, point_values_allowed FROM battles WHERE id = $1",
		BattleID,
	).Scan(
		&b.BattleID,
		&b.BattleName,
		&b.LeaderID,
		&b.VotingLocked,
		&ActivePlanID,
		&pv,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Not found")
	}

	_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
	b.ActivePlanID = ActivePlanID.String
	b.Warriors = GetBattleWarriors(BattleID)
	b.Plans = GetPlans(BattleID)

	return b, nil
}

// GetBattlesByWarrior gets a list of battles by WarriorID
func GetBattlesByWarrior(WarriorID string) ([]*Battle, error) {
	var battles = make([]*Battle, 0)
	battleRows, battlesErr := db.Query(`
		SELECT b.id, b.name, b.leader_id, b.voting_locked, b.active_plan_id, b.point_values_allowed,
		CASE WHEN COUNT(p) = 0 THEN '[]'::json ELSE array_to_json(array_agg(row_to_json(p))) END AS plans
		FROM battles b
		LEFT JOIN plans p ON b.id = p.battle_id
		LEFT JOIN battles_warriors bw ON b.id = bw.battle_id WHERE bw.warrior_id = $1
		GROUP BY b.id ORDER BY b.created_date DESC
	`, WarriorID)
	if battlesErr != nil {
		return nil, errors.New("Not found")
	}

	defer battleRows.Close()
	for battleRows.Next() {
		var points string
		var pv string
		var ActivePlanID sql.NullString
		var b = &Battle{
			BattleID:           "",
			LeaderID:           "",
			BattleName:         "",
			Warriors:           make([]*Warrior, 0),
			Plans:              make([]*Plan, 0),
			VotingLocked:       true,
			ActivePlanID:       "",
			PointValuesAllowed: make([]string, 0),
		}
		if err := battleRows.Scan(
			&b.BattleID,
			&b.BattleName,
			&b.LeaderID,
			&b.VotingLocked,
			&ActivePlanID,
			&pv,
			&points,
		); err != nil {
			log.Println(err)
		} else {
			_ = json.Unmarshal([]byte(points), &b.Plans)
			_ = json.Unmarshal([]byte(pv), &b.PointValuesAllowed)
			b.ActivePlanID = ActivePlanID.String
			battles = append(battles, b)
		}
	}

	return battles, nil
}

// ConfirmLeader confirms the warrior is infact leader of the battle
func ConfirmLeader(BattleID string, warriorID string) error {
	var leaderID string
	e := db.QueryRow("SELECT leader_id FROM battles WHERE id = $1", BattleID).Scan(&leaderID)
	if e != nil {
		log.Println(e)
		return errors.New("Battle Not found")
	}

	if leaderID != warriorID {
		return errors.New("Not Leader")
	}

	return nil
}

// GetBattleWarrior gets a warrior from db by ID and checks battle active status
func GetBattleWarrior(BattleID string, WarriorID string) (*Warrior, error) {
	var active bool
	var w Warrior

	e := db.QueryRow(
		`SELECT
			w.id, w.name, coalesce(w.email, ''), w.rank, coalesce(bw.active, FALSE)
		FROM warriors w
		LEFT JOIN battles_warriors bw ON bw.warrior_id = w.id AND bw.battle_id = $1
		WHERE id = $2`,
		BattleID,
		WarriorID,
	).Scan(
		&w.WarriorID,
		&w.WarriorName,
		&w.WarriorEmail,
		&w.WarriorRank,
		&active,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Warrior Not found")
	}

	if active {
		return nil, errors.New("Warrior Already Active in Battle")
	}

	return &w, nil
}

// GetBattleWarriors retrieves the warriors for a given battle from db
func GetBattleWarriors(BattleID string) []*Warrior {
	var warriors = make([]*Warrior, 0)
	rows, err := db.Query(
		`SELECT
			w.id, w.name, w.email, w.rank, bw.active
		FROM battles_warriors bw
		LEFT JOIN warriors w ON bw.warrior_id = w.id
		WHERE bw.battle_id = $1
		ORDER BY w.name`,
		BattleID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w Warrior
			var warriorEmail sql.NullString
			if err := rows.Scan(&w.WarriorID, &w.WarriorName, &warriorEmail, &w.WarriorRank, &w.Active); err != nil {
				log.Println(err)
			} else {
				w.WarriorEmail = warriorEmail.String
				warriors = append(warriors, &w)
			}
		}
	}

	return warriors
}

// AddWarriorToBattle adds a warrior by ID to the battle by ID
func AddWarriorToBattle(BattleID string, WarriorID string) ([]*Warrior, error) {
	if _, err := db.Exec(
		`INSERT INTO battles_warriors (battle_id, warrior_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (battle_id, warrior_id) DO UPDATE SET active = true`,
		BattleID,
		WarriorID,
	); err != nil {
		log.Println(err)
	}

	warriors := GetBattleWarriors(BattleID)

	return warriors, nil
}

// RetreatWarrior removes a warrior from the current battle by ID
func RetreatWarrior(BattleID string, WarriorID string) []*Warrior {
	if _, err := db.Exec(
		`UPDATE battles_warriors SET active = false WHERE battle_id = $1 AND warrior_id = $2`, BattleID, WarriorID); err != nil {
		log.Println(err)
	}

	if _, err := db.Exec(
		`UPDATE warriors SET last_active = NOW() WHERE id = $1`, WarriorID); err != nil {
		log.Println(err)
	}

	warriors := GetBattleWarriors(BattleID)

	return warriors
}

// GetPlans retrieves plans for given battle from db
func GetPlans(BattleID string) []*Plan {
	var plans = make([]*Plan, 0)
	planRows, plansErr := db.Query(
		"SELECT id, name, points, active, skipped, votestart_time, voteend_time, votes FROM plans WHERE battle_id = $1 ORDER BY created_date",
		BattleID,
	)
	if plansErr == nil {
		defer planRows.Close()
		for planRows.Next() {
			var v string
			var p = &Plan{PlanID: "",
				PlanName:      "",
				Votes:         make([]*Vote, 0),
				Points:        "",
				PlanActive:    false,
				PlanSkipped:   false,
				VoteStartTime: time.Now(),
				VoteEndTime:   time.Now(),
			}
			if err := planRows.Scan(&p.PlanID, &p.PlanName, &p.Points, &p.PlanActive, &p.PlanSkipped, &p.VoteStartTime, &p.VoteEndTime, &v); err != nil {
				log.Println(err)
			} else {
				err = json.Unmarshal([]byte(v), &p.Votes)
				if err != nil {
					log.Println(err)
				}

				// don't send vote values to client, prevent sneaky devs from peaking at votes
				for i := range p.Votes {
					vote := p.Votes[i]
					if p.PlanActive {
						vote.VoteValue = ""
					}
				}

				plans = append(plans, p)
			}
		}
	}

	return plans
}

// CreatePlan adds a new plan to a battle
func CreatePlan(BattleID string, warriorID string, PlanName string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	// @TODO - refactor stored procedure to replace need for app generated uuid
	newID, _ := uuid.NewUUID()
	PlanID := newID.String()

	if _, err := db.Exec(
		`call create_plan($1, $2, $3);`, BattleID, PlanID, PlanName,
	); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// ActivatePlanVoting sets the plan by ID to active, wipes any previous votes/points, and disables votingLock
func ActivatePlanVoting(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := db.Exec(
		`call activate_plan_voting($1, $2);`, BattleID, PlanID,
	); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// SetVote sets a warriors vote for the plan
func SetVote(BattleID string, WarriorID string, PlanID string, VoteValue string) []*Plan {
	if _, err := db.Exec(
		`call set_warrior_vote($1, $2, $3);`, PlanID, WarriorID, VoteValue); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans
}

// RetractVote removes a warriors vote for the plan
func RetractVote(BattleID string, WarriorID string, PlanID string) []*Plan {
	if _, err := db.Exec(
		`call retract_warrior_vote($1, $2);`, PlanID, WarriorID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans
}

// EndPlanVoting sets plan to active: false
func EndPlanVoting(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := db.Exec(
		`call end_plan_voting($1, $2);`, BattleID, PlanID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// SkipPlan sets plan to active: false and unsets battle's activePlanId
func SkipPlan(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := db.Exec(
		`call skip_plan_voting($1, $2);`, BattleID, PlanID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// RevisePlanName updates the plan name by ID
func RevisePlanName(BattleID string, warriorID string, PlanID string, PlanName string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	// set PlanID to true
	if _, err := db.Exec(
		`call revise_plan_name($2, $1);`, PlanName, PlanID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// BurnPlan removes a plan from the current battle by ID
func BurnPlan(BattleID string, warriorID string, PlanID string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := db.Exec(
		`call delete_plan($1, $2);`, BattleID, PlanID); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// FinalizePlan sets plan to active: false
func FinalizePlan(BattleID string, warriorID string, PlanID string, PlanPoints string) ([]*Plan, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := db.Exec(
		`call finalize_plan($1, $2, $3);`, BattleID, PlanID, PlanPoints); err != nil {
		log.Println(err)
	}

	plans := GetPlans(BattleID)

	return plans, nil
}

// SetBattleLeader sets the leaderId for the battle
func SetBattleLeader(BattleID string, warriorID string, LeaderID string) (*Battle, error) {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	// set battle VotingLocked
	if _, err := db.Exec(
		`call set_battle_leader($1, $2);`, BattleID, LeaderID); err != nil {
		log.Println(err)
	}

	battle, err := GetBattle(BattleID)
	if err != nil {
		return nil, errors.New("Unable to promote leader")
	}

	return battle, nil
}

// DeleteBattle removes all battle associations and the battle itself from DB by BattleID
func DeleteBattle(BattleID string, warriorID string) error {
	err := ConfirmLeader(BattleID, warriorID)
	if err != nil {
		return errors.New("Incorrect permissions")
	}

	if _, err := db.Exec(
		`call delete_battle($1);`, BattleID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

/*
	Warrior
*/

// GetWarrior gets a warrior from db by ID
func GetWarrior(WarriorID string) (*Warrior, error) {
	var w Warrior
	var warriorEmail sql.NullString

	e := db.QueryRow(
		"SELECT id, name, email, rank FROM warriors WHERE id = $1",
		WarriorID,
	).Scan(
		&w.WarriorID,
		&w.WarriorName,
		&warriorEmail,
		&w.WarriorRank,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Warrior Not found")
	}

	w.WarriorEmail = warriorEmail.String

	return &w, nil
}

// AuthWarrior attempts to authenticate the warrior
func AuthWarrior(WarriorEmail string, WarriorPassword string) (*Warrior, error) {
	var w Warrior
	var passHash string

	e := db.QueryRow(
		`SELECT id, name, email, rank, password FROM warriors WHERE email = $1`,
		WarriorEmail,
	).Scan(
		&w.WarriorID,
		&w.WarriorName,
		&w.WarriorEmail,
		&w.WarriorRank,
		&passHash,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Warrior Not found")
	}

	if ComparePasswords(passHash, []byte(WarriorPassword)) == false {
		return nil, errors.New("Password invalid")
	}

	return &w, nil
}

// CreateWarriorPrivate adds a new warrior private (guest) to the db
func CreateWarriorPrivate(WarriorName string) (*Warrior, error) {
	var WarriorID string
	e := db.QueryRow(`INSERT INTO warriors (name) VALUES ($1) RETURNING id`, WarriorName).Scan(&WarriorID)
	if e != nil {
		log.Println(e)
		return nil, errors.New("Unable to create new warrior")
	}

	return &Warrior{WarriorID: WarriorID, WarriorName: WarriorName}, nil
}

// CreateWarriorCorporal adds a new warrior corporal (registered) to the db
func CreateWarriorCorporal(WarriorName string, WarriorEmail string, WarriorPassword string) (*Warrior, error) {
	hashedPassword, hashErr := HashAndSalt([]byte(WarriorPassword))
	if hashErr != nil {
		return nil, hashErr
	}

	var WarriorID string
	WarriorRank := "CORPORAL"

	e := db.QueryRow(
		`INSERT INTO warriors (name, email, password, rank) VALUES ($1, $2, $3, $4) RETURNING id`,
		WarriorName,
		WarriorEmail,
		hashedPassword,
		WarriorRank,
	).Scan(&WarriorID)
	if e != nil {
		log.Println(e)
		return nil, errors.New("a Warrior with that email already exists")
	}

	return &Warrior{WarriorID: WarriorID, WarriorName: WarriorName, WarriorEmail: WarriorEmail, WarriorRank: WarriorRank}, nil
}

// UpdateWarriorProfile attempts to update the warriors profile
func UpdateWarriorProfile(WarriorID string, WarriorName string) error {
	if _, err := db.Exec(
		`UPDATE warriors SET name = $2 WHERE id = $1;`,
		WarriorID,
		WarriorName,
	); err != nil {
		log.Println(err)
		return errors.New("Error attempting to update warriors profile")
	}

	return nil
}

// WarriorResetRequest inserts a new warrior reset request
func WarriorResetRequest(WarriorEmail string) (resetID string, warriorName string, resetErr error) {
	var ResetID sql.NullString
	var WarriorID sql.NullString
	var WarriorName sql.NullString

	warErr := db.QueryRow(`
		SELECT w.id, w.name FROM warriors w WHERE w.email = $1;
		`,
		WarriorEmail,
	).Scan(&WarriorID, &WarriorName)
	if warErr != nil {
		log.Println("Unable to get warrior for reset email: ", warErr)
		// we don't want to alert the user that the email isn't valid
		return "", "", nil
	}

	e := db.QueryRow(`
		INSERT INTO warrior_reset (warrior_id)
		VALUES ($1)
		RETURNING reset_id;
		`,
		WarriorID.String,
	).Scan(&ResetID)
	if e != nil {
		log.Println("Unable to reset warrior: ", e)
		// we don't want to alert the user that the email isn't valid
		return "", "", nil
	}

	return ResetID.String, WarriorName.String, nil
}

// WarriorResetPassword attempts to reset a warriors password
func WarriorResetPassword(ResetID string, WarriorPassword string) (warriorName string, warriorEmail string, resetErr error) {
	var WarriorName sql.NullString
	var WarriorEmail sql.NullString

	hashedPassword, hashErr := HashAndSalt([]byte(WarriorPassword))
	if hashErr != nil {
		return "", "", hashErr
	}

	warErr := db.QueryRow(`
		SELECT
			w.name, w.email
		FROM warrior_reset wr
		LEFT JOIN warriors w ON w.id = wr.warrior_id
		WHERE wr.reset_id = $1;
		`,
		ResetID,
	).Scan(&WarriorName, &WarriorEmail)
	if warErr != nil {
		log.Println("Unable to get warrior for password reset confirmation email: ", warErr)
		return "", "", warErr
	}

	if _, err := db.Exec(
		`call reset_warrior_password($1, $2)`, ResetID, hashedPassword); err != nil {
		return "", "", err
	}

	return WarriorName.String, WarriorEmail.String, nil
}
