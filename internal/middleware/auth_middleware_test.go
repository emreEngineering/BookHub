package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"BookHub/internal/requestcontext"
)

type fakeSessionService struct {
	userID int
	err    error
}

func (s *fakeSessionService) CreateSession(ctx context.Context, userID int) (string, error) {
	return "test-session", nil
}

func (s *fakeSessionService) GetUserID(ctx context.Context, sessionID string) (int, error) {
	if s.err != nil {
		return 0, s.err
	}
	return s.userID, nil
}

func (s *fakeSessionService) DeleteSession(ctx context.Context, sessionID string) error {
	return nil
}

func TestAuthMiddleware_RequireAuth_AllowsValidSession(t *testing.T) {
	middleware := NewAuthMiddleware(&fakeSessionService{userID: 42})
	handlerCalled := false

	handler := middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		userID, ok := requestcontext.UserID(r.Context())
		if !ok || userID != 42 {
			t.Fatalf("expected userID 42 in request context, got %d, %v", userID, ok)
		}
		w.WriteHeader(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.AddCookie(&http.Cookie{Name: "session_id", Value: "valid-session"})
	recorder := httptest.NewRecorder()

	handler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	if !handlerCalled {
		t.Fatal("expected wrapped handler to be called")
	}
}

func TestAuthMiddleware_RequireAuth_RejectsMissingCookie(t *testing.T) {
	middleware := NewAuthMiddleware(&fakeSessionService{userID: 42})
	handlerCalled := false

	handler := middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	recorder := httptest.NewRecorder()

	handler(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}

	if handlerCalled {
		t.Fatal("expected wrapped handler not to be called")
	}
}

func TestAuthMiddleware_RequireAuth_RejectsInvalidSession(t *testing.T) {
	middleware := NewAuthMiddleware(&fakeSessionService{err: errors.New("invalid session")})
	handlerCalled := false

	handler := middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.AddCookie(&http.Cookie{Name: "session_id", Value: "invalid-session"})
	recorder := httptest.NewRecorder()

	handler(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}

	if handlerCalled {
		t.Fatal("expected wrapped handler not to be called")
	}
}
