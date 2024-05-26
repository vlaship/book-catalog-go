package repository

import (
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/database"
	"book-catalog/internal/logger"
	"context"
)

// AuthorRepository is an interface for author repository
type AuthorRepository struct {
	pool database.ConnPool
	log  logger.Logger
}

// NewAuthorRepository creates new author repository
func NewAuthorRepository(pool database.ConnPool, log logger.Logger) *AuthorRepository {
	return &AuthorRepository{
		pool: pool,
		log:  log.New("AuthorRepository"),
	}
}

func (r *AuthorRepository) l() logger.Logger {
	return r.log
}

func (r *AuthorRepository) p() database.ConnPool {
	return r.pool
}

const entityNameAuthor = "author"

const (
	getAuthors = `
	SELECT author_id, author_name, author_dob
	FROM catalog.authors WHERE deleted = FALSE;
`
	getAuthorByID = `
	SELECT author_id, author_name, author_dob
	FROM catalog.authors WHERE author_id = $1 AND deleted = FALSE;
`
	insertAuthor = `
	INSERT INTO catalog.authors (author_name, author_dob)
	VALUES ($1, $2)
	RETURNING author_id;
`
	updateAuthor = `
	UPDATE catalog.authors
	SET author_name = $2, author_dob = $3, updated_at = now()
	WHERE author_id = $1 AND deleted = FALSE;
`
	deleteAuthor = `
	UPDATE catalog.authors SET deleted = TRUE WHERE author_id = $1;
`
)

// GetAuthors returns all authors
func (r *AuthorRepository) GetAuthors(ctx context.Context) ([]model.Author, error) {
	r.log.Trc().Ctx(ctx).Msg("GetAuthors")

	req := getEntity[model.Author]{
		query:      getAuthors,
		entityName: entityNameAuthor,
		destinations: func(author *model.Author) []any {
			return []any{
				&author.ID,
				&author.Name,
				&author.Dob,
			}
		},
	}

	return getAll(ctx, r, req)
}

// GetAuthor returns author by id
func (r *AuthorRepository) GetAuthor(ctx context.Context, authorID types.ID) (*model.Author, error) {
	r.log.Dbg().Ctx(ctx).Values("authorID", authorID).Msg("GetAuthor")

	req := getEntity[model.Author]{
		query:      getAuthorByID,
		entityName: entityNameAuthor,
		args:       []any{authorID},
		destinations: func(author *model.Author) []any {
			return []any{
				&author.ID,
				&author.Name,
				&author.Dob,
			}
		},
	}

	return getOne(ctx, r, req)
}

// CreateAuthor inserts new author
func (r *AuthorRepository) CreateAuthor(ctx context.Context, author *model.Author) (*model.Author, error) {
	r.log.Dbg().Ctx(ctx).Values("Author", author).Msg("CreateAuthor")

	req := createEntity[model.Author]{
		getEntity: getEntity[model.Author]{
			query:        insertAuthor,
			entityName:   entityNameAuthor,
			args:         []any{author.Name, author.Dob},
			destinations: func(author *model.Author) []any { return []any{&author.ID} },
		},
	}

	return create(ctx, r, req)
}

// UpdateAuthor updates author
func (r *AuthorRepository) UpdateAuthor(
	ctx context.Context,
	authorID types.ID,
	author *model.Author,
) error {
	r.log.Dbg().Ctx(ctx).Values("author", author).Msg("UpdateAuthor")

	req := execRequest{
		query:      updateAuthor,
		entityName: entityNameAuthor,
		args: []any{
			authorID,
			author.Name,
			author.Dob,
		},
	}

	return exec(ctx, r, req)
}

// DeleteAuthor deletes author
func (r *AuthorRepository) DeleteAuthor(ctx context.Context, authorID types.ID) error {
	r.log.Dbg().Ctx(ctx).Values("authorID", authorID).Msg("DeleteAuthor")

	req := execRequest{
		query:      deleteAuthor,
		entityName: entityNameAuthor,
		args:       []any{authorID},
	}

	return exec(ctx, r, req)
}
