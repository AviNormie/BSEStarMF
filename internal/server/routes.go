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
	router.Use(cors.AllowAll().Handler)              // TODO configure CORS properly
	router.Use(middleware.Timeout(60 * time.Second)) // Set a timeout of 60 seconds
	router.Use(middleware.Recoverer)

	// Create SOAP service instance
	soapService, err := services.NewSOAPClientService()
	if err != nil {
		// Handle error appropriately - for now, we'll continue without SOAP service
		// In production, you might want to panic or return an error
	}

	router.Route(BaseUrl, func(r chi.Router) {
		r.Get("/health", handlers.HealthHandler)
		
		// SOAP-based SIP/XSIP Registration and Cancellation endpoints
		r.Post("/sip/order", handlers.SIPHandler)
		r.Post("/xsip/order", handlers.XSIPHandler)
		
		// Lumpsum order endpoint
		r.Post("/lumpsum/order", handlers.LumpsumHandler(soapService))
		// Enhanced JSON-based SIP/XSIP Cancellation endpoints
		r.Post("/enhanced/sip/cancellation", handlers.EnhancedSIPCancellationHandler(soapService))
		r.Post("/enhanced/xsip/cancellation", handlers.EnhancedXSIPCancellationHandler(soapService))
	})
}
