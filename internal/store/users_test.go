package store

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
	"visio/internal/types"
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

func migrateTestDB(db *sqlx.DB) error {
	q := `
create table if not exists users (
  id text not null primary key,
  email text not null unique,
  password_hash text not null,
  signup_date text not null
);

create table if not exists keys (
  id text not null primary key,
  user_id text not null references users(id) on delete cascade,
  prefix text not null unique,
  key_hash text not null,
  creation_date text not null
);

create table if not exists faces (
  id text not null primary key,
  label text not null,
  user_id text not null references users(id) on delete cascade,
  descriptor text not null,
  unique (label, user_id)
);
`
	_, err := db.Exec(q)
	return err
}

func TestUsers_Insert(t *testing.T) {
	s := NewUsersStore(testDB)

	t.Run("duplicate email", func(t *testing.T) {
		err := s.Insert(&types.User{
			Id:           "1",
			Email:        "foo@gmail.com",
			PasswordHash: "password1",
			SignupDate:   time.Now().UTC().Format("January, 2 2006"),
		})
		require.NoError(t, err)

		err = s.Insert(&types.User{
			Id:           "2",
			Email:        "foo@gmail.com",
			PasswordHash: "password2",
			SignupDate:   time.Now().UTC().Format("January, 2 2006"),
		})
		require.Error(t, err)
		var pqErr *pq.Error
		require.True(t, errors.As(err, &pqErr))
		// 23505 - unique_violation
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		require.Equal(t, "unique_violation", pqErr.Code.Name())

		testDB.Exec("delete from users cascade;")
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	t.Run("success", func(t *testing.T) {
		err := s.Insert(&types.User{
			Id:           "1",
			Email:        "foo@gmail.com",
			PasswordHash: "password1",
			SignupDate:   time.Now().UTC().Format("January, 2 2006"),
		})
		require.NoError(t, err)

		err = s.Insert(&types.User{
			Id:           "2",
			Email:        "bar@gmail.com",
			PasswordHash: "password2",
			SignupDate:   time.Now().UTC().Format("January, 2 2006"),
		})
		require.NoError(t, err)

		testDB.Exec("delete from users cascade;")
	})
}

func TestUsers_GetById(t *testing.T) {
	s := NewUsersStore(testDB)

	t.Run("user not found", func(t *testing.T) {
		user, err := s.GetById("1")
		require.Nil(t, user)
		require.Equal(t, err, types.ErrUserNotFound)

		testDB.Exec("delete from users cascade;")
	})

	t.Run("user exists", func(t *testing.T) {
		existingUser := types.User{
			Id:           "1",
			Email:        "emailz",
			PasswordHash: "passwordz",
			SignupDate:   time.Now().Truncate(time.Hour).UTC().Format("January, 2 2006"),
		}

		require.NoError(t, s.Insert(&existingUser))

		user, err := s.GetById("1")
		require.Equal(t, user.Id, existingUser.Id)
		require.Equal(t, user.Email, existingUser.Email)
		require.Equal(t, user.PasswordHash, existingUser.PasswordHash)
		require.Equal(t, user.SignupDate, existingUser.SignupDate)
		require.NoError(t, err)

		testDB.Exec("delete from users cascade;")
	})
}

func TestUsers_GetByEmail(t *testing.T) {
	s := NewUsersStore(testDB)

	t.Run("user not found", func(t *testing.T) {
		user, err := s.GetByEmail("emailz")
		require.Nil(t, user)
		require.Equal(t, err, types.ErrUserNotFound)

		testDB.Exec("delete from users cascade;")
	})

	t.Run("user exists", func(t *testing.T) {
		existingUser := types.User{
			Id:           "1",
			Email:        "emailz",
			PasswordHash: "passwordz",
			SignupDate:   time.Now().Truncate(time.Hour).UTC().Format("January, 2 2006"),
		}

		require.NoError(t, s.Insert(&existingUser))

		user, err := s.GetByEmail("emailz")
		require.Equal(t, user.Id, existingUser.Id)
		require.Equal(t, user.Email, existingUser.Email)
		require.Equal(t, user.PasswordHash, existingUser.PasswordHash)
		require.Equal(t, user.SignupDate, existingUser.SignupDate)
		require.NoError(t, err)

		testDB.Exec("delete from users cascade;")
	})
}
