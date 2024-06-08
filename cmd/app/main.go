package main

import (
	"github.com/vlaship/book-catalog-go/internal/app"
	"github.com/vlaship/book-catalog-go/internal/config"
	"github.com/vlaship/book-catalog-go/internal/logger"
	"log/slog"
	"os"
)

func main() {
	slog.Info("Book Catalog is waking up...")

	// load config
	cfg := config.MustGet()

	// create logger
	log := logger.NewLogger(cfg)

	if err := app.Run(cfg, log); err != nil {
		log.Err(err).Msg("failed to start server")
		os.Exit(1)
	}

	log.Inf().Msg("server shut down")
}
