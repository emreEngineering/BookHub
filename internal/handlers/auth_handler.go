package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"BookHub/internal/activity"
	"BookHub/internal/responses"
	"BookHub/internal/services"
)

type AuthHandler struct {
	userService    services.UserService
	sessionService services.SessionService
	activityLogger activity.ActivityLogger
}

func NewAuthHandler(userService services.UserService, sessionService services.SessionService, activityLogger activity.ActivityLogger) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		sessionService: sessionService,
		activityLogger: activityLogger,
	}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece POST destekler")
		return
	}

	var request services.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "Geçersiz JSON")
		return
	}

	user, err := h.userService.Register(r.Context(), request)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	responses.Success(w, http.StatusCreated, "Kullanıcı oluşturuldu", user)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece POST destekler")
		return
	}

	var request services.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "Geçersiz JSON")
		return
	}

	user, err := h.userService.Login(r.Context(), request)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	sessionID, err := h.sessionService.CreateSession(r.Context(), user.ID)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	})

	responses.Success(w, http.StatusOK, "Giriş başarılı", user)
}

func (h *AuthHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, "Buendpoint sadece GET destekler")
		return
	}
	cookie, err := r.Cookie("session_id")

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "Giriş yapmalısnız")
		return
	}
	userID, err := h.sessionService.GetUserID(r.Context(), cookie.Value)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "Giriş yapmalısınız")
		return
	}
	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "Giriş yapmalısınız")
		return
	}

	responses.Success(w, http.StatusOK, "Kullanıcı getirildi", user)
}
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece post destekler")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "zaten giriş yapılmamış")
		return
	}

	userID, _ := h.sessionService.GetUserID(r.Context(), cookie.Value)

	err = h.sessionService.DeleteSession(r.Context(), cookie.Value)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Çıkış yapılamadı")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	var userIDPtr *int
	if userID != 0 {
		userIDPtr = &userID
	}
	h.logActivity(r.Context(), "user_logout", "Kullanıcı çıkış yaptı", userIDPtr, nil)

	responses.Success(w, http.StatusOK, "Çıkış başarılı", nil)
}

func (h *AuthHandler) logActivity(ctx context.Context, eventType string, message string, userID *int, metadata map[string]interface{}) {
	if h.activityLogger == nil {
		return
	}

	_ = h.activityLogger.Log(ctx, eventType, message, userID, metadata)
}
