package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sapphirebroking.com/sapphire_mf/internal/config"
	"sapphirebroking.com/sapphire_mf/internal/consumer"
	"sapphirebroking.com/sapphire_mf/internal/processor"
	"sapphirebroking.com/sapphire_mf/internal/server"
	"sapphirebroking.com/sapphire_mf/internal/util"
)

func main(){
	logger := util.NewStandardLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config: %v", err)
	}
	logger.Info("Configuration loaded successfully.")

	// Initialize message processor
	msgProcessor := processor.NewProcessor(logger)
	logger.Info("Message processor initialized.")

	// Initialize Kafka consumer
	kafkaConsumer, err := consumer.NewConsumer(cfg.Kafka, msgProcessor, logger)
	if err != nil {
		logger.Fatal("Failed to create kafka consumer: %v", err)
	}
	logger.Info("Kafka consumer initialized.")

	// Initialize HTTP server
	httpServer := server.NewHTTPServer(cfg.Service, logger)
	logger.Info("HTTP server initialized.")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka consumer in goroutine
	go kafkaConsumer.Start(ctx)
	logger.Info("Kafka consumer started.")

	// Start HTTP server in goroutine
	go httpServer.Start()
	logger.Info("HTTP server started.")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Received shutdown signal, shutting down gracefully")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Cancel the main context to stop Kafka consumer
	cancel()

	// Shutdown Kafka consumer
	kafkaConsumer.Shutdown()
	logger.Info("Kafka consumer shutdown complete.")

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Failed to shutdown HTTP server: %v", err)
	}
	logger.Info("HTTP server shutdown complete.")

	logger.Info("Shutdown complete")

} 