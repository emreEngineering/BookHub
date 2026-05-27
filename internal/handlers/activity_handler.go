package handlers

import (
	"net/http"
	"strconv"

	"BookHub/internal/activity"
	"BookHub/internal/responses"
)

type ActivityHandler struct {
	activityLogger activity.ActivityLogger
}

func NewActivityHandler(activityLogger activity.ActivityLogger) *ActivityHandler {
	return &ActivityHandler{
		activityLogger: activityLogger,
	}
}

func (h *ActivityHandler) ListActivityLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece GET destekler")
		return
	}

	limit := int64(20)
	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		parsedLimit, err := strconv.ParseInt(limitParam, 10, 64)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, "Geçersiz limit")
			return
		}
		limit = parsedLimit
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	logs, err := h.activityLogger.FindLatest(r.Context(), limit)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Activity loglar alınamadı")
		return
	}

	responses.Success(w, http.StatusOK, "Activity loglar getirildi", logs)
}
