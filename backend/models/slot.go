package models

type Slot struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Date        string `gorm:"not null" json:"date"`
	Time        string `gorm:"not null" json:"time"`
	Capacity    int    `gorm:"not null" json:"capacity"`
	BookedCount int    `gorm:"not null;default:0" json:"bookedCount"`
}
