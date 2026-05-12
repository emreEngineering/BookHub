package middleware

import (
	"BookHub/internal/responses"
	"BookHub/internal/services"
	"net/http"
)

type AuthMiddleware struct {
	sessionService services.SessionService
}

func NewAuthMiddleware(sessionService services.SessionService) *AuthMiddleware {
	return &AuthMiddleware{
		sessionService: sessionService,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			responses.Error(w, http.StatusUnauthorized, "Giriş yapmalısınız")
			return
		}
		_, err = m.sessionService.GetUserID(cookie.Value)

		if err != nil {
			responses.Error(w, http.StatusUnauthorized, "Giriş yapmalısınız")
			return
		}

		next(w, r)
	}
}
