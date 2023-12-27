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

func (s *Sessions) Create(sessionId, userId string) error {
	err := s.redis.Set(
		context.Background(),
		sessionId,
		userId,
		0,
	).Err()
	if err != nil {
		return fmt.Errorf("Error while creating new session: %w", err)
	}
	return nil
}

func (s *Sessions) Get(sessionId string) (string, error) {
	sessionValue, err := s.redis.Get(
		context.Background(),
		sessionId,
	).Result()
	if err != nil {
		if err == redis.Nil {
			return "", types.ErrSessionNotFound
		}
		return "", fmt.Errorf("Error while fetching session from redis: %w", err)
	}
	return sessionValue, nil
}

func (s *Sessions) Delete(sessionId string) error {
	err := s.redis.Del(
		context.Background(),
		sessionId,
	).Err()
	if err != nil {
		if err == redis.Nil {
			return types.ErrSessionNotFound
		}
		return fmt.Errorf("Error while delete session from redis: %w", err)
	}
	return nil
}
