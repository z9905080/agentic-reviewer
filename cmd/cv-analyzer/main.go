package main

import (
	"context" // Added for future graceful shutdown
	// "log" // Standard log removed, using zerolog
	"net/http" // Added for http.Server
	"os"
	"os/signal" // Added for graceful shutdown
	"syscall"   // Added for graceful shutdown
	"time"      // Added for graceful shutdown

	// "github.com/gin-gonic/gin" // Gin engine comes from initializeAPI, direct import not needed here
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log" // Alias for global logger
	// _ "cv-analyzer/api/docs" // Preserving Swagger docs import if it was there. Script didn't include it.
)

// @title CV Analyzer API
// @version 1.0
// @description This is a server for analyzing CVs against Job Descriptions.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
	// zerolog setup
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs // More performant
	// Use ConsoleWriter for human-readable output during development
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	zlog.Logger = zlog.Output(consoleWriter)
	// Set global log level (e.g., from env var later)
	zerolog.SetGlobalLevel(zerolog.InfoLevel) // Default to Info
	// Example: Read from ENV VAR
	// logLevel := os.Getenv("LOG_LEVEL")
	// if level, err := zerolog.ParseLevel(logLevel); err == nil {
	// 	zerolog.SetGlobalLevel(level)
	// }


	// Initialize API using Wire
	// The wire_gen.go initializeAPI() returns (*gin.Engine, error)
	engine, err := initializeAPI()
	if err != nil {
		// Use standard log here as zerolog might not be fully set up if init fails early
		// Or ensure zlog is safe to use (it should be after global setup)
		zlog.Fatal().Err(err).Msg("Failed to initialize API")
	}
	zlog.Info().Msg("API initialized successfully")


	// Server configuration (example, can be from config later)
	serverAddr := ":8080" // TODO: Make this configurable
	zlog.Info().Str("address", serverAddr).Msg("Server starting")

	srv := &http.Server{
		Addr:    serverAddr,
		Handler: engine, // Gin engine as handler
		// TODO: Add ReadTimeout, WriteTimeout, IdleTimeout for production robustness
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Fatal().Err(err).Msg("Server ListenAndServe failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit // Block until a signal is received
	zlog.Info().Str("signal", sig.String()).Msg("Received shutdown signal, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 30-second timeout for shutdown
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zlog.Fatal().Err(err).Msg("Server forced to shutdown due to error during graceful shutdown")
	}

	zlog.Info().Msg("Server exited gracefully")
}
