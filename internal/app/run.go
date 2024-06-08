package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/vlaship/book-catalog-go/internal/config"
	"github.com/vlaship/book-catalog-go/internal/database"
	"github.com/vlaship/book-catalog-go/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Run starts the application
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

	// Start server
	server := &http.Server{
		Addr:         cfg.ServerProps.Port,
		Handler:      app.Router,
		ReadTimeout:  cfg.ServerProps.ReadTimeout,
		WriteTimeout: cfg.ServerProps.WriteTimeout,
		IdleTimeout:  cfg.ServerProps.IdleTimeout,
	}

	go func() {
		log.Inf().Msg("starting http-server...")
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Err(err).Msg("failed to start server: %w", err)
		}
	}()

	log.Inf().Msg("Book Catalog is started...")
	log.Inf().Msg("Port %v", cfg.ServerProps.Port)

	// wait for interrupt signal to gracefully shut down the server
	log.Dbg().Msg("waiting for interrupt signal...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Inf().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ServerProps.CancelContextTimeout)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to gracefully shut down server: %w", err)
	}

	return nil
}
