package middleware

import (
	"BookHub/internal/requestcontext"
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
		userID, err := m.userIDFromRequest(r)
		if err != nil {
			responses.Error(w, http.StatusUnauthorized, "Giriş yapmalısınız")
			return
		}

		next(w, r.WithContext(requestcontext.WithUserID(r.Context(), userID)))
	}
}

func (m *AuthMiddleware) RequireWebAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := m.userIDFromRequest(r)
		if err != nil {
			http.Redirect(w, r, "/web/login", http.StatusSeeOther)
			return
		}

		next(w, r.WithContext(requestcontext.WithUserID(r.Context(), userID)))
	}
}

func (m *AuthMiddleware) RequireAuthForMethods(next http.HandlerFunc, methods ...string) http.HandlerFunc {
	protectedMethods := make(map[string]struct{}, len(methods))
	for _, method := range methods {
		protectedMethods[method] = struct{}{}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if _, protected := protectedMethods[r.Method]; !protected {
			next(w, r)
			return
		}

		m.RequireAuth(next)(w, r)
	}
}

func (m *AuthMiddleware) userIDFromRequest(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, err
	}

	return m.sessionService.GetUserID(r.Context(), cookie.Value)
}
