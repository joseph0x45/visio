package tests

import (
	"errors"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/redis/go-redis/v9"
	"testing"
	"visio/internal/store"
	"visio/internal/types"
)

type testCase struct {
	Name         string
	SessionId    string
	SessionValue string
	TestFunc     func(sessions *store.Sessions, tc *testCase) error
	ExpectedErr  error
}

func (t *testCase) Run(sessions *store.Sessions) error {
	return t.TestFunc(sessions, t)
}

var testCases = []testCase{
	{
		Name:         "Create session",
		SessionId:    "Session1",
		SessionValue: "Session1Value",
		TestFunc: func(sessions *store.Sessions, tc *testCase) error {
			return sessions.Create(tc.SessionId, tc.SessionValue)
		},
		ExpectedErr: nil,
	},
	{
		Name:         "Retrieve session using it's id",
		SessionId:    "Session1",
		SessionValue: "Session1Value",
		TestFunc: func(sessions *store.Sessions, tc *testCase) error {
			sessionValue, err := sessions.Get(tc.SessionId)
			if err != nil {
				return err
			}
			if sessionValue != tc.SessionValue {
				return fmt.Errorf("Invalid session value returned: Wanted %s got %s", tc.SessionValue, sessionValue)
			}
			return nil
		},
		ExpectedErr: nil,
	},
	{
		Name:      "Retrieve session using a non-existent id",
		SessionId: "RandomId",
		TestFunc: func(sessions *store.Sessions, tc *testCase) error {
			_, err := sessions.Get(tc.SessionId)
			return err
		},
		ExpectedErr: types.ErrSessionNotFound,
	},
	{
		Name:      "Delete session",
		SessionId: "Session1",
		TestFunc: func(sessions *store.Sessions, tc *testCase) error {
			return sessions.Delete(tc.SessionId)
		},
		ExpectedErr: nil,
	},
	{
		Name:      "Delete session using a non-existent id",
		SessionId: "RandomId",
		TestFunc: func(sessions *store.Sessions, tc *testCase) error {
			return sessions.Delete(tc.SessionId)
		},
		ExpectedErr: nil,
	},
}

func TestSession(t *testing.T) {
	var redisClient *redis.Client
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to create pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		t.Fatalf("Failed to connect to Docker: %s", err)
	}
	resource, err := pool.Run("redis", "latest", nil)
	if err != nil {
		t.Fatalf("Failed to run Redis container: %s", err)
	}
	if err = pool.Retry(func() error {
		redisClient = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		})
		sessions := store.NewSessionsStore(redisClient)
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				err = tc.Run(sessions)
				if !errors.Is(err, tc.ExpectedErr) {
					t.Fatalf("Wanted %s got %s ", tc.ExpectedErr, err)
				}
			})
		}
		return nil
	}); err != nil {
		t.Fatalf("Failed to connect to Redis: %s", err)
	}

	if err = pool.Purge(resource); err != nil {
		t.Fatalf("Failed to purge resource: %s", err)
	}
}
