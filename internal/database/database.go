package database

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vlaship/book-catalog-go/internal/apperr"
)

// GetErrorByCode returns error by pgErr.Code
func GetErrorByCode(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return apperr.ErrAlreadyExists
		}
		if pgErr.Code == "23503" {
			return apperr.ErrNoFoundForeignKey
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return apperr.ErrNotFound
	}

	return err
}

// CheckAffectedRows checks affected rows
func CheckAffectedRows(tag pgconn.CommandTag) error {
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	if tag.RowsAffected() > 1 {
		return apperr.ErrAffectedMoreThanOneRow
	}

	return nil
}
