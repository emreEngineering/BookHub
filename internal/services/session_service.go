package services

import "context"

type SessionService interface {
	CreateSession(ctx context.Context, userID int) (string, error)
	GetUserID(ctx context.Context, sessionID string) (int, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
