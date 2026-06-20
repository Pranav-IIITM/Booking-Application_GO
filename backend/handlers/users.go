package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"booking-backend/middleware"
	"booking-backend/models"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type UsersHandler struct {
	Firestore *firestore.Client
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

	users := h.Firestore.Collection("users")
	iter := users.Where("firebase_uid", "==", firebaseUID).Limit(1).Documents(r.Context())
	snapshot, err := iter.Next()
	iter.Stop()
	if err != iterator.Done {
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not load user")
			return
		}

		var user models.User
		if err := snapshot.DataTo(&user); err != nil {
			writeError(w, http.StatusInternalServerError, "could not load user")
			return
		}
		user.ID = snapshot.Ref.ID
		writeJSON(w, http.StatusOK, map[string]any{"user": user})
		return
	}

	userRef := users.NewDoc()
	user := models.User{
		ID:          userRef.ID,
		FirebaseUID: firebaseUID,
		Name:        request.Name,
		Email:       request.Email,
		CreatedAt:   time.Now().UTC(),
	}
	if _, err := userRef.Set(r.Context(), user); err != nil {
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"user": user})
}
