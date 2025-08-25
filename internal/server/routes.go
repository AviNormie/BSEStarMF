package server

import (
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"sapphirebroking.com/sapphire_mf/internal/server/handlers"
	"sapphirebroking.com/sapphire_mf/internal/server/services"
)

// SetupRoutes configures all routes for the MF application
func SetupRoutes(router chi.Router) {
	// Setup middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(cors.AllowAll().Handler)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Recoverer)

	// Create SOAP service instance
	soapService, err := services.NewSOAPClientService()
	if err != nil {
		// Log error but continue - handle gracefully in handlers
		panic("Failed to initialize SOAP service: " + err.Error())
	}

	// Add a root health check
	router.Get("/health", handlers.HealthHandler)

	// API v1 routes
	router.Route(BaseUrl, func(r chi.Router) {
		r.Get("/health", handlers.HealthHandler)
		
		// Authentication endpoint
		r.Post("/auth/getPassword", handlers.GetPasswordHandler)
		
		// Order entry endpoints
		r.Post("/order/sip", handlers.SIPHandler)
		r.Post("/order/xsip", handlers.XSIPHandler)
		r.Post("/order/lumpsum", handlers.LumpsumHandler(soapService))
		
		// Enhanced cancellation endpoints
		r.Post("/cancellation/enhanced-sip", handlers.EnhancedSIPCancellationHandler(soapService))
		r.Post("/cancellation/enhanced-xsip", handlers.EnhancedXSIPCancellationHandler(soapService))
	})
}
