package handlers

import (
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

	user, err := h.userService.Register(request)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logActivity("user_registered", "Kullanıcı kayıt oldu", &user.ID, map[string]interface{}{
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	})

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

	user, err := h.userService.Login(request)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	sessionID, err := h.sessionService.CreateSession(user.ID)
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

	h.logActivity("user_login", "Kullanıcı giriş yaptı", &user.ID, map[string]interface{}{
		"email": user.Email,
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
	userID, err := h.sessionService.GetUserID(cookie.Value)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "Giriş yapmalısınız")
		return
	}
	user, err := h.userService.GetUserByID(userID)
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

	userID, _ := h.sessionService.GetUserID(cookie.Value)

	err = h.sessionService.DeleteSession(cookie.Value)
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
	h.logActivity("user_logout", "Kullanıcı çıkış yaptı", userIDPtr, nil)

	responses.Success(w, http.StatusOK, "Çıkış başarılı", nil)
}

func (h *AuthHandler) logActivity(eventType string, message string, userID *int, metadata map[string]interface{}) {
	if h.activityLogger == nil {
		return
	}

	_ = h.activityLogger.Log(eventType, message, userID, metadata)
}
