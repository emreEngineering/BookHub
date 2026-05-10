package handlers

import (
	"encoding/json"
	"net/http"

	"BookHub/internal/responses"
	"BookHub/internal/services"
)

type AuthHandler struct {
	userService    services.UserService
	sessionService services.SessionService
}

func NewAuthHandler(userService services.UserService, sessionService services.SessionService) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		sessionService: sessionService,
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
