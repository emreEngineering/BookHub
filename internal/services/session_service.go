package services

import (
	"encoding/hex"
	"errors"
	"math/rand"
)

type SessionService interface {
	CreateSession(userID int) (string, error)
	GetUserID(sessionID string) (int, error)
	DeleteSession(sessionID string) error
}

type MemorySessionService struct {
	session map[string]int
}

func NewSessionService() *MemorySessionService {
	return &MemorySessionService{
		session: make(map[string]int),
	}
}

func (s *MemorySessionService) CreateSession(userID int) (string, error) {
	sessionID, err := generateSessionID()

	if err != nil {
		return "", err
	}
	s.session[sessionID] = userID
	return sessionID, nil
}

func (s *MemorySessionService) GetUserID(sessionID string) (int, error) {
	userID, exists := s.session[sessionID]

	if !exists {
		return 0, errors.New("geçersiz session")
	}
	return userID, nil
}

func (s *MemorySessionService) DeleteSession(sessionID string) error {
	delete(s.session, sessionID)
	return nil
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
