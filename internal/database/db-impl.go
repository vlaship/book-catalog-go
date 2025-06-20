package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/vlaship/book-catalog-go/internal/config"
	"github.com/vlaship/book-catalog-go/internal/logger"
	"time"

	zerologadapter "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

const (
	// FailedBeginTransaction Error messages
	FailedBeginTransaction = "failed to begin transaction: %v"
	// FailedCommitTransaction Error messages
	FailedCommitTransaction = "failed to commit transaction: %v"
)

// ConnPoolImpl is a struct that represents a connection pool to a Postgres database
type ConnPoolImpl struct {
	pool *pgxpool.Pool
}

// New creates a new connection pool to a Postgres database
func New(cfg *config.Config, log logger.Logger) (ConnPool, error) {
	log.Trc().Msg("parse config for database")

	// parse db config
	dbCfg, err := pgxpool.ParseConfig(cfg.ConnDB)
	if err != nil {
		log.Err(err).Msg("failed to parse config for database")
		return nil, err
	}

	// Optimize connection pool settings
	dbCfg.MaxConns = 25                    // Set max connections
	dbCfg.MinConns = 5                     // Set min connections
	dbCfg.MaxConnLifetime = 1 * time.Hour  // Max connection lifetime
	dbCfg.MaxConnIdleTime = 30 * time.Minute // Max idle time
	dbCfg.HealthCheckPeriod = 1 * time.Minute // Health check period

	logLevel, err := tracelog.LogLevelFromString(cfg.LogLevelDB)
	if err != nil {
		return nil, fmt.Errorf("invalid db log level: [%s]", cfg.LogLevelDB)
	}
	dbCfg.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   zerologadapter.NewLogger(log.New("pgx").Logger()),
		LogLevel: logLevel,
	}

	ctx := context.Background()
	log.Trc().Msg("connect to database")
	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		log.Err(err).Msg("failed to connect to database")
		return nil, err
	}

	log.Trc().Msg("ping database")
	if err = pool.Ping(ctx); err != nil {
		log.Err(err).Msg("failed to ping database")
		return nil, err
	}

	return &ConnPoolImpl{pool: pool}, nil
}

// Begin starts a transaction
func (cp *ConnPoolImpl) Begin(ctx context.Context) (pgx.Tx, error) {
	return cp.pool.Begin(ctx)
}

// Close closes the connection pool
func (cp *ConnPoolImpl) Close() {
	cp.pool.Close()
}
