package handlers

import (
	"encoding/json"
	"net/http"

	"BookHub/internal/responses"
	"BookHub/internal/services"
)

type AuthHandler struct {
	userService services.UserService
}

func NewAuthHandler(userService services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
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
