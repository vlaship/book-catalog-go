package database

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// ConnPool is a connection pool interface.
//
//go:generate mockgen -destination=../../test/mock/database/mock-conn_pool.go -package=mock . ConnPool
type ConnPool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}
