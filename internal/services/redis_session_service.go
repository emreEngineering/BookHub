package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSessionService struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisSessionService(client *redis.Client) *RedisSessionService {
	return &RedisSessionService{
		client: client,
		ttl:    time.Hour,
	}
}

func (s *RedisSessionService) CreateSession(userID int) (string, error) {
	sessionID, err := generateRedisSessionID()
	if err != nil {
		return "", err
	}

	err = s.client.Set(context.Background(), sessionKey(sessionID), strconv.Itoa(userID), s.ttl).Err()
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *RedisSessionService) GetUserID(sessionID string) (int, error) {
	userIDValue, err := s.client.Get(context.Background(), sessionKey(sessionID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, errors.New("geçersiz session")
		}
		return 0, err
	}

	userID, err := strconv.Atoi(userIDValue)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *RedisSessionService) DeleteSession(sessionID string) error {
	return s.client.Del(context.Background(), sessionKey(sessionID)).Err()
}

func sessionKey(sessionID string) string {
	return "session:" + sessionID
}

func generateRedisSessionID() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
