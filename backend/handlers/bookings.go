package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"booking-backend/middleware"
	"booking-backend/models"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingsHandler struct {
	Firestore *firestore.Client
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

	slotID, err := parseIDField(request.SlotID, "slotId")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var booking models.Booking
	err = h.Firestore.RunTransaction(r.Context(), func(ctx context.Context, tx *firestore.Transaction) error {
		userQuery := h.Firestore.Collection("users").Where("firebase_uid", "==", firebaseUID).Limit(1)
		userSnapshots, err := tx.Documents(userQuery).GetAll()
		if err != nil {
			return err
		}
		if len(userSnapshots) == 0 {
			return errUserNotFound
		}
		userSnapshot := userSnapshots[0]
		var user models.User
		if err := userSnapshot.DataTo(&user); err != nil {
			return err
		}
		user.ID = userSnapshot.Ref.ID

		slotRef := h.Firestore.Collection("slots").Doc(slotID)
		slotSnapshot, err := tx.Get(slotRef)
		if err != nil {
			return err
		}
		var slot models.Slot
		if err := slotSnapshot.DataTo(&slot); err != nil {
			return err
		}
		slot.ID = slotSnapshot.Ref.ID
		if slot.BookedCount >= slot.Capacity {
			return errSlotFull
		}

		bookingRef := h.Firestore.Collection("bookings").NewDoc()
		booking = models.Booking{
			ID:        bookingRef.ID,
			UserID:    user.ID,
			SlotID:    slot.ID,
			Status:    "confirmed",
			CreatedAt: time.Now().UTC(),
			Slot:      &slot,
		}
		if err := tx.Set(bookingRef, booking); err != nil {
			return err
		}

		return tx.Update(slotRef, []firestore.Update{
			{Path: "bookedCount", Value: firestore.Increment(1)},
		})
	})

	if err != nil {
		switch {
		case status.Code(err) == codes.NotFound, errors.Is(err, errUserNotFound):
			writeError(w, http.StatusNotFound, "user or slot not found")
		case errors.Is(err, errSlotFull):
			writeError(w, http.StatusBadRequest, "slot is already full")
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

	userIter := h.Firestore.Collection("users").Where("firebase_uid", "==", firebaseUID).Limit(1).Documents(r.Context())
	userSnapshot, err := userIter.Next()
	userIter.Stop()
	if err == iterator.Done {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}
	var user models.User
	if err := userSnapshot.DataTo(&user); err != nil {
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}
	user.ID = userSnapshot.Ref.ID

	var bookings []models.Booking
	iter := h.Firestore.Collection("bookings").Where("userID", "==", user.ID).Documents(r.Context())
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not load bookings")
			return
		}

		var booking models.Booking
		if err := snapshot.DataTo(&booking); err != nil {
			writeError(w, http.StatusInternalServerError, "could not load bookings")
			return
		}
		booking.ID = snapshot.Ref.ID

		slotSnapshot, err := h.Firestore.Collection("slots").Doc(booking.SlotID).Get(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not load bookings")
			return
		}
		var slot models.Slot
		if err := slotSnapshot.DataTo(&slot); err != nil {
			writeError(w, http.StatusInternalServerError, "could not load bookings")
			return
		}
		slot.ID = slotSnapshot.Ref.ID
		booking.Slot = &slot
		bookings = append(bookings, booking)
	}

	sort.Slice(bookings, func(i, j int) bool {
		return bookings[i].CreatedAt.After(bookings[j].CreatedAt)
	})

	writeJSON(w, http.StatusOK, map[string]any{"bookings": bookings})
}

var errSlotFull = errors.New("slot full")
var errUserNotFound = errors.New("user not found")

func parseIDField(raw json.RawMessage, fieldName string) (string, error) {
	if len(raw) == 0 {
		return "", fmt.Errorf("%s is required", fieldName)
	}

	var numeric uint64
	if err := json.Unmarshal(raw, &numeric); err == nil {
		if numeric == 0 {
			return "", fmt.Errorf("%s must be greater than zero", fieldName)
		}
		return fmt.Sprintf("%d", numeric), nil
	}

	var text string
	if err := json.Unmarshal(raw, &text); err != nil {
		return "", fmt.Errorf("%s must be a number", fieldName)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return "", fmt.Errorf("%s is required", fieldName)
	}

	return text, nil
}
