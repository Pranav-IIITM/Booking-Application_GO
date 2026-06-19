package routes

import (
	"encoding/json"
	"net/http"

	"booking-backend/config"
	"booking-backend/handlers"
	authmw "booking-backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:5500",
			"http://localhost:5500",
			"http://localhost:5501",
			"http://127.0.0.1:5501",
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	})

	slotsHandler := &handlers.SlotsHandler{DB: db}
	usersHandler := &handlers.UsersHandler{DB: db}
	bookingsHandler := &handlers.BookingsHandler{DB: db}
	meHandler := &handlers.MeHandler{DB: db}
	authMiddleware := authmw.FirebaseAuth(cfg.FirebaseAuth)

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	r.Get("/api/slots", slotsHandler.List)

	r.Group(func(protected chi.Router) {
		protected.Use(authMiddleware)
		protected.Get("/api/me", meHandler.Current)
		protected.Post("/api/users/sync", usersHandler.Sync)
		protected.Post("/api/book", bookingsHandler.Create)
		protected.Get("/api/bookings", bookingsHandler.List)
	})

	return r
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
