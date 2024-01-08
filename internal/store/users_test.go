package store

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	"visio/internal/types"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

const (
	testDBUser     = "testuser"
	testDBPassword = "testpassword"
	testDBName     = "testdb"
)

var testDB *sqlx.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "11",
			Env: []string{
				"POSTGRES_PASSWORD=" + testDBPassword,
				"POSTGRES_USER=" + testDBUser,
				"POSTGRES_DB=" + testDBName,
				"listen_addresses = '*'",
			},
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		testDBUser,
		testDBPassword,
		resource.GetHostPort("5432/tcp"),
		testDBName,
	)

	if err = pool.Retry(func() error {
		testDB, err = sqlx.Connect("postgres", dbURL)
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if err := migrateTestDB(testDB); err != nil {
		log.Fatalf("Fail to migrate test DB: %v", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

// TODO: As an improvement, we might consider to use a migration folder
//
// github.com/golang-migrate/migrate
func migrateTestDB(db *sqlx.DB) error {
	q := `
create table if not exists users (
	id text not null primary key,
	email text not null unique,
	password text not null,
	signup_date timestamp not null
);`
	_, err := db.Exec(q)
	return err
}

func TestUsers_Insert(t *testing.T) {
	s := NewUsersStore(testDB)

	t.Run("duplicate email", func(t *testing.T) {
		err := s.Insert(&types.User{
			Id:         "1",
			Email:      "foo@gmail.com",
			Password:   "password1",
			SignupDate: time.Now(),
		})
		require.NoError(t, err)

		err = s.Insert(&types.User{
			Id:         "2",
			Email:      "foo@gmail.com",
			Password:   "password2",
			SignupDate: time.Now(),
		})
		require.Error(t, err)
		var pqErr *pq.Error
		require.True(t, errors.As(err, &pqErr))
		// 23505 - unique_violation
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		require.Equal(t, "unique_violation", pqErr.Code.Name())

		testDB.Exec("TRUNCATE users;")
	})

	t.Run("success", func(t *testing.T) {
		err := s.Insert(&types.User{
			Id:         "1",
			Email:      "foo@gmail.com",
			Password:   "password1",
			SignupDate: time.Now(),
		})
		require.NoError(t, err)

		err = s.Insert(&types.User{
			Id:         "2",
			Email:      "bar@gmail.com",
			Password:   "password2",
			SignupDate: time.Now(),
		})
		require.NoError(t, err)

		testDB.Exec("TRUNCATE users;")
	})
}
