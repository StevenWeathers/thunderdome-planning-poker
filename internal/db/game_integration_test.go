package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/poker"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db/retro"
)

// testDB creates a temporary database on a Postgres server, migrates it and drops it afterwards.
// Tests are skipped unless TEST_DB_HOST is set, TEST_DB_PORT defaults to 5432 and the credentials
// match docker-compose.yml.
func testDB(t *testing.T) *db.Service {
	t.Helper()

	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		t.Skip("TEST_DB_HOST is not set")
	}
	port := 5432
	if p := os.Getenv("TEST_DB_PORT"); p != "" {
		var err error
		port, err = strconv.Atoi(p)
		require.NoError(t, err)
	}

	admin, err := sql.Open("pgx", fmt.Sprintf(
		"host=%s port=%d user=thor password=odinson dbname=thunderdome sslmode=disable", host, port))
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })

	name := fmt.Sprintf("thunderdome_test_%d", time.Now().UnixNano())
	_, err = admin.Exec(`CREATE DATABASE ` + name)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := admin.Exec(`DROP DATABASE ` + name + ` WITH (FORCE)`)
		assert.NoError(t, err)
	})

	d := db.New("", &db.Config{
		Host: host, Port: port, User: "thor", Password: "odinson", Name: name, SSLMode: "disable",
		MaxOpenConns: 20, MaxIdleConns: 20, ConnMaxLifetime: 5,
	}, otelzap.New(zap.NewNop()), true)
	t.Cleanup(func() { _ = d.DB.Close() })

	return d
}

func insertID(t *testing.T, d *db.Service, query string, args ...any) string {
	t.Helper()

	var id string
	require.NoError(t, d.DB.QueryRow(query, args...).Scan(&id))

	return id
}

func insertUser(t *testing.T, d *db.Service) string {
	return insertID(t, d, `INSERT INTO thunderdome.users (name, type) VALUES ('integration', 'GUEST') RETURNING id`)
}

func TestConcurrentGameCreation(t *testing.T) {
	d := testDB(t)
	ctx := context.Background()
	userID := insertUser(t, d)
	teamID := insertID(t, d, `INSERT INTO thunderdome.team (name) VALUES ('integration') RETURNING id`)
	scaleID := insertID(t, d, `SELECT id FROM thunderdome.estimation_scale WHERE name = 'Thunderdome Default'`)
	templateID := insertID(t, d, `SELECT id FROM thunderdome.retro_template LIMIT 1`)
	retros := &retro.Service{DB: d.DB, Logger: d.Logger}
	games := &poker.Service{DB: d.DB, Logger: d.Logger}

	tests := []struct {
		name   string
		table  string
		create func() error
	}{
		{"retro", "retro", func() error {
			_, err := retros.CreateRetro(ctx, userID, "", "integration", "", "", 3, "visible",
				0, false, false, false, false, templateID)
			return err
		}},
		{"poker game", "poker", func() error {
			_, err := games.CreateGame(ctx, userID, "integration", scaleID, []string{"1"}, nil,
				false, "ceil", "", "", false)
			return err
		}},
		{"team poker game", "poker", func() error {
			_, err := games.TeamCreateGame(ctx, teamID, userID, "integration", scaleID, []string{"1"}, nil,
				false, "ceil", "", "", false)
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// conflicts are most likely while the table is small enough to be scanned
			_, err := d.DB.Exec(`VACUUM ANALYZE thunderdome.` + tt.table)
			require.NoError(t, err)

			const workers, perWorker = 8, 25
			errs := make(chan error, workers*perWorker)
			var wg sync.WaitGroup
			for range workers {
				wg.Go(func() {
					for range perWorker {
						if err := tt.create(); err != nil {
							errs <- err
						}
					}
				})
			}
			wg.Wait()
			close(errs)

			assert.Empty(t, errs, "first error: %v", <-errs)
		})
	}
}

func TestLastActiveTriggers(t *testing.T) {
	d := testDB(t)

	for _, game := range []string{"poker", "retro", "storyboard"} {
		t.Run(game, func(t *testing.T) {
			ownerID := insertUser(t, d)
			joinerID := insertUser(t, d)
			gameID := insertID(t, d,
				`INSERT INTO thunderdome.`+game+` (owner_id, name) VALUES ($1, 'integration') RETURNING id`, ownerID)

			exec := func(query string, args ...any) {
				t.Helper()
				_, err := d.DB.Exec(query, args...)
				require.NoError(t, err)
			}
			setAge := func(age string) {
				exec(`UPDATE thunderdome.`+game+` SET last_active = NOW() - $2::interval WHERE id = $1`, gameID, age)
			}
			age := func() time.Duration {
				var lastActive time.Time
				require.NoError(t, d.DB.QueryRow(
					`SELECT last_active FROM thunderdome.`+game+` WHERE id = $1`, gameID).Scan(&lastActive))
				return time.Since(lastActive)
			}
			userTable := `thunderdome.` + game + `_user`
			join := `INSERT INTO ` + userTable + ` (` + game + `_id, user_id) VALUES ($1, $2)`

			setAge("2 days")
			exec(join, gameID, joinerID)
			assert.Less(t, age(), time.Minute, "insert is activity")

			setAge("2 days")
			exec(`UPDATE `+userTable+` SET active = true WHERE user_id = $1`, joinerID)
			assert.Less(t, age(), time.Minute, "update is activity")

			setAge("10 minutes")
			exec(`UPDATE `+userTable+` SET active = false WHERE user_id = $1`, joinerID)
			assert.Greater(t, age(), 9*time.Minute, "recent activity is not rewritten")

			// guest cleanup cascades to these rows and must not revive old games
			setAge("2 days")
			exec(`DELETE FROM thunderdome.users WHERE id = $1`, joinerID)
			assert.Greater(t, age(), 47*time.Hour, "delete is not activity")

			// deleting a game with children must still work
			exec(join, gameID, ownerID)
			exec(`DELETE FROM thunderdome.`+game+` WHERE id = $1`, gameID)
		})
	}

	t.Run("every trigger runs for each row", func(t *testing.T) {
		var total, perRow int
		require.NoError(t, d.DB.QueryRow(`
			SELECT count(*), count(*) FILTER (WHERE tgtype & 1 = 1)
			FROM pg_trigger WHERE tgname LIKE '%\_last\_active'`).Scan(&total, &perRow))

		assert.Equal(t, 11, total)
		assert.Equal(t, total, perRow)
	})
}
