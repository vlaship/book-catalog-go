package repository

import (
	"book-catalog/internal/app/model"
	"book-catalog/internal/database"
	"book-catalog/internal/logger"
	"context"
	"github.com/jackc/pgx/v5"
)

// Repo is an interface for repositories
//
//go:generate mockgen -destination=../../../test/mock/repository/mock-repo.go -package=mocks . Repo
type Repo interface {
	l() logger.Logger
	p() database.ConnPool
}

type execRequest struct {
	query      string
	entityName string
	args       []any
}

// getEntity is a struct for get one/all request
type getEntity[T model.Entity] struct {
	//entity       T
	query        string
	entityName   string
	args         []any
	destinations func(t *T) []any
}

// getProperty is a struct for get property request
type getProperty[T model.Property] struct {
	propertyName string
}

type createEntity[T model.Entity] struct {
	getEntity[T]
}

func getOne[T model.Entity](
	ctx context.Context,
	r Repo,
	req getEntity[T],
) (*T, error) {
	tx, err := r.p().Begin(ctx)
	if err != nil {
		r.l().Err(err).Ctx(ctx).Msg(database.FailedBeginTransaction)
		return nil, database.GetErrorByCode(err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, req.query, req.args...)
	if err != nil {
		r.l().Wrn().Err(err).Ctx(ctx).Msg("failed to query %s", req.entityName)
		return nil, database.GetErrorByCode(err)
	}

	res, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (T, error) {
		var t T
		if err := row.Scan(req.destinations(&t)...); err != nil {
			return t, err
		}
		return t, nil
	})
	if err != nil {
		r.l().Err(err).Ctx(ctx).Msg("failed to scan %s", req.entityName)
		return nil, database.GetErrorByCode(err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.l().Err(err).Ctx(ctx).Msg(database.FailedCommitTransaction)
		return nil, database.GetErrorByCode(err)
	}

	return &res, nil
}

func getAll[T model.Entity](
	ctx context.Context,
	r Repo,
	req getEntity[T],
) ([]T, error) {
	tx, err := r.p().Begin(ctx)
	if err != nil {
		r.l().Err(err).Ctx(ctx).Msg(database.FailedBeginTransaction)
		return nil, database.GetErrorByCode(err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, req.query, req.args...)
	if err != nil {
		r.l().Wrn().Err(err).Ctx(ctx).Msg("failed to query %s", req.entityName)
		return nil, database.GetErrorByCode(err)
	}

	entities, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (T, error) {
		var t T
		if err := row.Scan(req.destinations(&t)...); err != nil {
			return t, err
		}
		return t, nil
	})
	if err != nil {
		r.l().Err(err).Ctx(ctx).Msg("failed to scan %s", req.entityName)
		return nil, database.GetErrorByCode(err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.l().Err(err).Ctx(ctx).Msg(database.FailedCommitTransaction)
		return nil, database.GetErrorByCode(err)
	}

	return entities, nil
}

func create[T model.Entity](
	ctx context.Context,
	r Repo,
	req createEntity[T],
) (*T, error) {
	tx, err := r.p().Begin(ctx)
	if err != nil {
		r.l().Err(err).Ctx(ctx).Msg(database.FailedBeginTransaction)
		return nil, database.GetErrorByCode(err)
	}
	defer tx.Rollback(ctx)

	var t T
	if err := tx.QueryRow(ctx, req.query, req.args...).Scan(req.destinations(&t)...); err != nil {
		r.l().Wrn().Err(err).Ctx(ctx).Msg("failed to create %s", req.entityName)
		return nil, database.GetErrorByCode(err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.l().Err(err).Ctx(ctx).Msg(database.FailedCommitTransaction)
		return nil, database.GetErrorByCode(err)
	}

	return &t, nil
}

func exec(
	ctx context.Context,
	r Repo,
	req execRequest,
) error {
	tx, err := r.p().Begin(ctx)
	if err != nil {
		r.l().Err(err).Ctx(ctx).Msg(database.FailedBeginTransaction)
		return database.GetErrorByCode(err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(
		ctx,
		req.query,
		req.args...,
	)
	if err != nil {
		r.l().Err(err).Ctx(ctx).Msg("failed to upsert %s", req.entityName)
		return database.GetErrorByCode(err)
	}

	if err := database.CheckAffectedRows(tag); err != nil {
		return database.GetErrorByCode(err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.l().Err(err).Ctx(ctx).Msg(database.FailedCommitTransaction)
		return database.GetErrorByCode(err)
	}

	return nil
}
