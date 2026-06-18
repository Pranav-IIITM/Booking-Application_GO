package handlers

import (
	"net/http"

	"booking-backend/models"
	"gorm.io/gorm"
)

type SlotsHandler struct {
	DB *gorm.DB
}

func (h *SlotsHandler) List(w http.ResponseWriter, r *http.Request) {
	var slots []models.Slot
	if err := h.DB.Where("booked_count < capacity").Order("date asc, time asc").Find(&slots).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "could not load slots")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"slots": slots})
}
