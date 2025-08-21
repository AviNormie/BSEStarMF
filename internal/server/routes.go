package server

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"sapphirebroking.com/sapphire_mf/internal/server/handlers"
)

// SetupRoutes configures all routes for the MF application
func SetupRoutes(router *chi.Mux) {
	// Setup middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(cors.AllowAll().Handler)              // TODO configure CORS properly
	router.Use(middleware.Timeout(60 * time.Second)) // Set a timeout of 60 seconds
	router.Use(middleware.Recoverer)

	router.Route(BaseUrl, func(r chi.Router) {
		r.Get("/health", handlers.HealthHandler)
	})
}
