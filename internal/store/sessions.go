package store

import (
	"context"
	"fmt"
	"visio/internal/types"
	"github.com/redis/go-redis/v9"
)

type Sessions struct {
	redis *redis.Client
}

func NewSessionsStore(redis *redis.Client) *Sessions {
	return &Sessions{
		redis: redis,
	}
}

func (s *Sessions) Create(id, value string) error {
	err := s.redis.Set(
		context.Background(),
		id,
		value,
		0,
	).Err()
	if err != nil {
		return fmt.Errorf("Error while creating new session: %w", err)
	}
	return nil
}

func (s *Sessions) Get(id string) (string, error) {
	sessionValue, err := s.redis.Get(
		context.Background(),
		id,
	).Result()
	if err != nil {
		if err == redis.Nil {
			return "", types.ErrSessionNotFound
		}
		return "", fmt.Errorf("Error while getting session: %w", err)
	}
	return sessionValue, nil
}
