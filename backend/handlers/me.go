package handlers

import (
	"errors"
	"net/http"

	"booking-backend/middleware"
	"booking-backend/models"
	"gorm.io/gorm"
)

type MeHandler struct {
	DB *gorm.DB
}

func (h *MeHandler) Current(w http.ResponseWriter, r *http.Request) {
	firebaseUID, ok := middleware.FirebaseUID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var user models.User
	if err := h.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}
