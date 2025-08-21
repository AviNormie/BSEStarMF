package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"sapphirebroking.com/sapphire_mf/internal/config"
	"sapphirebroking.com/sapphire_mf/internal/util"
)

type HTTPServer struct {
	logger util.Logger
	server *http.Server
}

const BaseUrl string = "/api/v1"

func NewHTTPServer(cfg *config.ServiceConfig, logger util.Logger) *HTTPServer {
	// Create Chi router
	router := chi.NewRouter()

	hs := &HTTPServer{
		logger: logger,
	}

	// Setup routes with middleware
	SetupRoutes(router)

	// Determine the address to bind to
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	if cfg.Host == "" {
		addr = fmt.Sprintf(":%d", cfg.Port)
	}

	hs.server = &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return hs
}

func (hs *HTTPServer) Start() {
	hs.logger.Info("Starting HTTP server on %s", hs.server.Addr)
	if err := hs.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		hs.logger.Fatal("Could not start HTTP server: %v", err)
	}
}

func (hs *HTTPServer) Shutdown(ctx context.Context) error {
	hs.logger.Info("Shutting down HTTP server...")
	return hs.server.Shutdown(ctx)
}

// SetupRoutes configures all routes for the MF application
func SetupRoutes(router *chi.Mux) {
	// Setup middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(cors.AllowAll().Handler)              // TODO configure CORS properly
	router.Use(middleware.Timeout(60 * time.Second)) // Set a timeout of 60 seconds
	router.Use(middleware.Recoverer)

	// Handle 404 errors
	router.NotFound(NotFoundHandler)
	router.MethodNotAllowed(MethodNotAllowedHandler)

	// API versioning
	router.Route(BaseUrl, func(r chi.Router) {
		// Health check routes
		r.Get("/health-check", HealthHandler)

		// MF endpoints - you can add your handlers here
		r.Get("/ping", PingHandler)
	})
}

// Basic handlers
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"sapphire-mf"}`))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"pong","service":"sapphire-mf"}`))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"endpoint not found"}`))
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(`{"error":"method not allowed"}`))
}