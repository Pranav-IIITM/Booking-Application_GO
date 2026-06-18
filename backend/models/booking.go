package models

import "time"

type Booking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	SlotID    uint      `gorm:"not null;index" json:"slotId"`
	Status    string    `gorm:"not null;default:confirmed" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Slot      Slot      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"slot,omitempty"`
}
