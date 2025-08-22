package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"sapphirebroking.com/sapphire_mf/internal/config"
	"sapphirebroking.com/sapphire_mf/internal/util"
)

type HTTPServer struct {
	logger util.Logger
	server *http.Server
}

const BaseUrl string = "/api/v1"

func NewHTTPServer(cfg *config.ServiceConfig, logger util.Logger) *HTTPServer {
	hs := &HTTPServer{
		logger: logger,
	}

	// Setup routes - pass logger instead of router
	router := SetupRoutes(logger)

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