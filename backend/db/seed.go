package db

import (
	"booking-backend/models"
	"gorm.io/gorm"
)

func SeedSlots(database *gorm.DB) error {
	var count int64
	if err := database.Model(&models.Slot{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	slots := []models.Slot{
		{Date: "2026-06-20", Time: "10:00 AM", Capacity: 10, BookedCount: 0},
		{Date: "2026-06-20", Time: "12:00 PM", Capacity: 10, BookedCount: 0},
		{Date: "2026-06-21", Time: "10:00 AM", Capacity: 5, BookedCount: 0},
		{Date: "2026-06-21", Time: "02:00 PM", Capacity: 8, BookedCount: 0},
		{Date: "2026-06-22", Time: "11:00 AM", Capacity: 6, BookedCount: 0},
	}

	return database.Create(&slots).Error
}
