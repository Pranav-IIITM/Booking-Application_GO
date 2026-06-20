package main

import (
	"context"
	"log"
	"net/http"

	"booking-backend/config"
	"booking-backend/db"
	"booking-backend/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	authClient, firestoreClient, err := config.InitFirebase()
	if err != nil {
		log.Fatalf("init firebase: %v", err)
	}
	defer firestoreClient.Close()
	cfg.FirebaseAuth = authClient

	ctx := context.Background()
	if err := db.SeedSlots(ctx, firestoreClient); err != nil {
		log.Printf("seed slots: %v", err)
	}

	router := routes.NewRouter(cfg, firestoreClient)
	addr := ":" + cfg.Port

	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
