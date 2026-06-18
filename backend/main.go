package main

import (
	"log"
	"net/http"

	"booking-backend/config"
	database "booking-backend/db"
	"booking-backend/models"
	"booking-backend/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Slot{}, &models.Booking{}); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}
	if err := database.SeedSlots(db); err != nil {
		log.Fatalf("seed slots: %v", err)
	}

	router := routes.NewRouter(cfg, db)
	addr := ":" + cfg.Port

	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
