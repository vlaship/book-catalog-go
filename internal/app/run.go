package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vlaship/book-catalog-go/internal/config"
	"github.com/vlaship/book-catalog-go/internal/database"
	"github.com/vlaship/book-catalog-go/internal/logger"
)

// Run starts the application
// Ensure all resources are properly closed during shutdown
func Run(cfg *config.Config, log logger.Logger) error {
	log.Inf().Msg("Book Catalog is initializing...")

	// start db migration
	log.Trc().Msg("start db migration...")
	if err := database.GooseMigrateDatabase(cfg, log); err != nil {
		return fmt.Errorf("failed db migration: %w", err)
	}

	// init app
	log.Trc().Msg("init app...")
	app, err := NewApp(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	defer app.DB.Close()

	// Start server with context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := &http.Server{
		Addr:         cfg.ServerProps.Port,
		Handler:      app.Router,
		ReadTimeout:  cfg.ServerProps.ReadTimeout,
		WriteTimeout: cfg.ServerProps.WriteTimeout,
		IdleTimeout:  cfg.ServerProps.IdleTimeout,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		log.Inf().Msg("starting http-server...")
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Err(err).Msg("failed to start server")
			cancel()
		}
	}()

	log.Inf().Msg("Book Catalog is started...")
	log.Inf().Msg(fmt.Sprintf("Port %v", cfg.ServerProps.Port))

	// wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Inf().Msg("shutting down server...")
	case <-ctx.Done():
		// Ensure all necessary cleanup tasks are performed when a shutdown signal is received
		log.Inf().Msg("server shutdown requested...")
	}

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ServerProps.CancelContextTimeout)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err = server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to gracefully shut down server: %w", err)
	}

	return nil
}
