package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"sapphirebroking.com/sapphire_mf/internal/server/handlers"
	"sapphirebroking.com/sapphire_mf/internal/server/services"
	"sapphirebroking.com/sapphire_mf/internal/util"
)

// SetupRoutes configures all API routes
// Move SOAP routes outside the /api/v1 group:
func SetupRoutes(logger util.Logger) http.Handler {
    r := chi.NewRouter()
    
    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)
    r.Use(middleware.Timeout(60 * time.Second))

    // CORS configuration
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    // Initialize services
    elogService := services.NewELOGClientService(logger)

    // Initialize ELOG handler
    elogHandler := handlers.NewELOGHandler(elogService, logger)

    // API routes
    r.Route("/api/v1", func(r chi.Router) {
        // Health check
        r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{"status":"ok","timestamp":` + fmt.Sprintf("%d", time.Now().Unix()) + `}`))
        })

        // ELOG endpoints
        r.Post("/elog/request", elogHandler.ELOGRequestHandler)
        r.Get("/elog/callback", elogHandler.ELOGCallbackHandler)

        // Authentication endpoints
        r.Post("/auth/getPassword", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{"success":true,"message":"Authentication endpoint"}`))
        })

        // REAL SIP/XSIP HANDLERS - REPLACE PLACEHOLDERS
        r.Post("/sip/order", handlers.SIPHandler)
        r.Post("/xsip/order", handlers.XSIPHandler)
        r.Post("/enhanced/sip/cancellation", handlers.EnhancedSIPCancellationHandler)
        r.Post("/enhanced/xsip/cancellation", handlers.EnhancedXSIPCancellationHandler)

        // Add these SOAP-specific routes
        // Replace lines 71-74 with:
        r.Route("/soap", func(r chi.Router) {
            // Use existing SIP handler (already supports SOAP XML)
            r.Post("/SipCancellation", handlers.SIPHandler)
            
            // Use existing XSIP handler (already supports SOAP XML)
            r.Post("/XsipCancellation", handlers.XSIPHandler)
            
            // Comment out WSDL endpoints (handlers don't exist)
            // r.Get("/SipCancellation?wsdl", handlers.SIPCancellationWSDLHandler)
            // r.Get("/XsipCancellation?wsdl", handlers.XSIPCancellationWSDLHandler)
        })
    })
    
    // SOAP routes (separate from /api/v1)
    r.Route("/soap", func(r chi.Router) {
        r.Post("/SipCancellation", handlers.SIPHandler)
        r.Post("/XsipCancellation", handlers.XSIPHandler)
    })
    
    return r
}
