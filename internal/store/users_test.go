package store_test

import (
	"errors"
	"testing"
	"time"
	"visio/internal/store"
	"visio/internal/types"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestUsers_Insert(t *testing.T) {
	signupDate := time.Now()
	errFailToInsert := errors.New("fail to insert")

	testCases := []struct {
		name   string
		user   types.User
		mockFn func(sqlmock.Sqlmock)
		err    error
	}{
		{
			name: "success",
			user: types.User{
				Id:         "id",
				Email:      "email",
				Password:   "password",
				SignupDate: signupDate,
			},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectExec("insert into users").
					WithArgs("id", "email", "password", signupDate).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			err: nil,
		},
		{
			name: "fail to insert",
			user: types.User{
				Id:         "id",
				Email:      "email",
				Password:   "password",
				SignupDate: signupDate,
			},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectExec("insert into users").
					WithArgs("id", "email", "password", signupDate).
					WillReturnError(errFailToInsert)
			},
			err: errFailToInsert,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Fail to initiate sqlmock: %v", err)
			}
			defer db.Close()

			tc.mockFn(mock)

			sqlxDb := sqlx.NewDb(db, "sqlmock")
			usersStore := store.NewUsersStore(sqlxDb)
			err = usersStore.Insert(&tc.user)
			if !errors.Is(err, tc.err) {
				t.Errorf("Got %v, want %v", err, tc.err)
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Fail expectation: %v", err)
			}
		})
	}
}
