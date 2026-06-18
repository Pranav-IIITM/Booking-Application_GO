package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"booking-backend/middleware"
	"booking-backend/models"
	"gorm.io/gorm"
)

type UsersHandler struct {
	DB *gorm.DB
}

type syncUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UsersHandler) Sync(w http.ResponseWriter, r *http.Request) {
	firebaseUID, ok := middleware.FirebaseUID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var request syncUserRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}

	if request.Name == "" {
		request.Name = middleware.Name(r.Context())
	}
	if request.Email == "" {
		request.Email = middleware.Email(r.Context())
	}

	var user models.User
	err := h.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error
	if err == nil {
		updates := map[string]any{}
		if request.Name != "" && request.Name != user.Name {
			updates["name"] = request.Name
		}
		if request.Email != "" && request.Email != user.Email {
			updates["email"] = request.Email
		}
		if len(updates) > 0 {
			if err := h.DB.Model(&user).Updates(updates).Error; err != nil {
				writeError(w, http.StatusInternalServerError, "could not update user")
				return
			}
			if err := h.DB.First(&user, user.ID).Error; err != nil {
				writeError(w, http.StatusInternalServerError, "could not reload user")
				return
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{"user": user})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}

	user = models.User{
		FirebaseUID: firebaseUID,
		Name:        request.Name,
		Email:       request.Email,
	}
	if err := h.DB.Create(&user).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"user": user})
}
