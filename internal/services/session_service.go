package services

type SessionService interface {
	CreateSession(userID int) (string, error)
	GetUserID(sessionID string) (int, error)
	DeleteSession(sessionID string) error
}
