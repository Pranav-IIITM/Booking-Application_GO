package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"booking-backend/middleware"
	"booking-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingsHandler struct {
	DB *gorm.DB
}

type createBookingRequest struct {
	SlotID json.RawMessage `json:"slotId"`
}

func (h *BookingsHandler) Create(w http.ResponseWriter, r *http.Request) {
	firebaseUID, ok := middleware.FirebaseUID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var request createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	slotID, err := parseUintField(request.SlotID, "slotId")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var booking models.Booking
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
			return err
		}

		var slot models.Slot
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&slot, slotID).Error; err != nil {
			return err
		}
		if slot.BookedCount >= slot.Capacity {
			return errSlotFull
		}

		booking = models.Booking{
			UserID: user.ID,
			SlotID: slot.ID,
			Status: "confirmed",
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		slot.BookedCount++
		if err := tx.Model(&slot).Update("booked_count", slot.BookedCount).Error; err != nil {
			return err
		}

		return tx.Preload("Slot").First(&booking, booking.ID).Error
	})

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			writeError(w, http.StatusNotFound, "user or slot not found")
		case errors.Is(err, errSlotFull):
			writeError(w, http.StatusConflict, "slot is already full")
		default:
			writeError(w, http.StatusInternalServerError, "could not create booking")
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"booking": booking})
}

func (h *BookingsHandler) List(w http.ResponseWriter, r *http.Request) {
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

	var bookings []models.Booking
	if err := h.DB.Where("user_id = ?", user.ID).Preload("Slot").Order("created_at desc").Find(&bookings).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "could not load bookings")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"bookings": bookings})
}

var errSlotFull = errors.New("slot full")

func parseUintField(raw json.RawMessage, fieldName string) (uint, error) {
	if len(raw) == 0 {
		return 0, fmt.Errorf("%s is required", fieldName)
	}

	var numeric uint64
	if err := json.Unmarshal(raw, &numeric); err == nil {
		if numeric == 0 {
			return 0, fmt.Errorf("%s must be greater than zero", fieldName)
		}
		return uint(numeric), nil
	}

	var text string
	if err := json.Unmarshal(raw, &text); err != nil {
		return 0, fmt.Errorf("%s must be a number", fieldName)
	}

	parsed, err := strconv.ParseUint(text, 10, 64)
	if err != nil || parsed == 0 {
		return 0, fmt.Errorf("%s must be a positive number", fieldName)
	}

	return uint(parsed), nil
}
