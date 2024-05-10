package tests

import (
	"net"
	"net/url"
	DB "p-system/db"
	"runtime"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

func StartDB(tb testing.TB) *sqlx.DB {
	tb.Helper()

	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("user", "password"),
		Path:   "testdb",
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()

	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("failed to create Docker pool: %v", err)
	}
	if err != nil {
		tb.Fatalf("failed to create Docker pool: %v", err)
	}

	pw, _ := pgURL.User.Password()
	env := []string{
		"POSTGRES_USER=" + pgURL.User.Username(),
		"POSTGRES_PASSWORD=" + pw,
		"POSTGRES_DB=" + pgURL.Path,
	}

	// Start the container.
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "alpine",
		Env:        env,
	})
	if err != nil {
		tb.Fatalf("failed to start postgres container: %v", err)
	}

	tb.Cleanup(func() {
		if err := pool.Purge(container); err != nil {
			tb.Fatalf("failed to purge container: %v", err)
		}
	})

	pgURL.Host = container.GetHostPort("5432/tcp")

	// Get the host.
	pgURL.Host = container.Container.NetworkSettings.IPAddress

	// On Mac, Docker runs in a VM.
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(container.GetBoundIP("5432/tcp"), container.GetPort("5432/tcp"))
	}

	// Retry until we can establish a connection to the database.
	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() (err error) {
		db, err := sqlx.Open("postgres", pgURL.String())
		if err != nil {
			return err
		}
		defer func() {
			cerr := db.Close()
			if err != nil {
				err = cerr
			}
		}()

		return db.Ping()
	})
	if err != nil {
		tb.Fatalf("Failed to start postgres: %v", err)
	}

	db, err := sqlx.Open("postgres", pgURL.String())
	if err != nil {
		tb.Fatalf("failed to open database: %v", err)
	}

	time.Sleep(10 * time.Second) // This is a hack. Must find a better way to wait for the db to establish a connection
	if err := DB.Migrate(db.DB); err != nil {
		tb.Fatalf("failed to run migrations: %v", err)
	}

	return db

}

const seeds = `INSERT INTO users (id,username, email, password)
VALUES ('d164e69d-26f5-448d-a18c-baeae517d9f2','john_doed', 'john@ample.com', 'hashed_password_here');
INSERT INTO wallets (id ,user_id, Balance)
VALUES ('d164e69d-26f5-448d-a18c-baeae517d991','d164e69d-26f5-448d-a18c-baeae517d9f2', 0);
INSERT INTO transactions (id,user_id, request_id, amount, type, status, reference)
VALUES ('d164e69d-26f5-448d-a18c-baeae517d9f5','d164e69d-26f5-448d-a18c-baeae517d9f2', 'unique_request_id', 1000, 'credit', 'completed', 'unique_reference');

`

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
