package handlers

import (
	"net/http"
	"sort"

	"booking-backend/models"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type SlotsHandler struct {
	Firestore *firestore.Client
}

func (h *SlotsHandler) List(w http.ResponseWriter, r *http.Request) {
	var slots []models.Slot
	iter := h.Firestore.Collection("slots").Documents(r.Context())
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not load slots")
			return
		}

		var slot models.Slot
		if err := snapshot.DataTo(&slot); err != nil {
			writeError(w, http.StatusInternalServerError, "could not load slots")
			return
		}
		slot.ID = snapshot.Ref.ID
		if slot.BookedCount < slot.Capacity {
			slots = append(slots, slot)
		}
	}

	sort.Slice(slots, func(i, j int) bool {
		if slots[i].Date == slots[j].Date {
			return slots[i].Time < slots[j].Time
		}
		return slots[i].Date < slots[j].Date
	})

	writeJSON(w, http.StatusOK, map[string]any{"slots": slots})
}
