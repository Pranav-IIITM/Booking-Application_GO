package models

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FirebaseUID string    `gorm:"uniqueIndex;not null" json:"firebaseUid"`
	Name        string    `json:"name"`
	Email       string    `gorm:"index" json:"email"`
	CreatedAt   time.Time `json:"createdAt"`
}
