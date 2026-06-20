package models

import "time"

type User struct {
	ID          string    `firestore:"-" json:"id"`
	FirebaseUID string    `firestore:"firebase_uid" json:"firebaseUid"`
	Name        string    `firestore:"name" json:"name"`
	Email       string    `firestore:"email" json:"email"`
	CreatedAt   time.Time `firestore:"createdAt" json:"createdAt"`
}
