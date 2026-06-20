package models

import "time"

type Booking struct {
	ID        string    `firestore:"-" json:"id"`
	UserID    string    `firestore:"userID" json:"userId"`
	SlotID    string    `firestore:"slotID" json:"slotId"`
	Status    string    `firestore:"status" json:"status"`
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
	Slot      *Slot     `firestore:"-" json:"slot,omitempty"`
}
