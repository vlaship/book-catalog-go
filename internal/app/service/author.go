package service

import (
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/logger"
	"context"
)

// AuthorReader is an interface for author reader
//
//go:generate mockgen -destination=../../../test/mock/service/mock-author-reader.go -package=mock . AuthorReader
type AuthorReader interface {
	GetAuthors(ctx context.Context) ([]model.Author, error)
	GetAuthor(ctx context.Context, authorID types.ID) (*model.Author, error)
}

// AuthorWriter is an interface for author writer
//
//go:generate mockgen -destination=../../../test/mock/service/mock-author-writer.go -package=mock . AuthorWriter
type AuthorWriter interface {
	CreateAuthor(ctx context.Context, author *model.Author) (*model.Author, error)
	UpdateAuthor(ctx context.Context, authorID types.ID, author *model.Author) error
	DeleteAuthor(ctx context.Context, authorID types.ID) error
}

// AuthorService is a service for author
type AuthorService struct {
	reader AuthorReader
	writer AuthorWriter
	log    logger.Logger
}

// NewAuthorService creates new author service
func NewAuthorService(
	reader AuthorReader,
	writer AuthorWriter,
	log logger.Logger,
) *AuthorService {
	return &AuthorService{
		reader: reader,
		writer: writer,
		log:    log.New("AuthorService"),
	}
}

// GetAuthors returns all authors
func (s *AuthorService) GetAuthors(ctx context.Context) ([]model.Author, error) {
	s.log.Trc().Ctx(ctx).Msg("GetAuthors")

	return s.reader.GetAuthors(ctx)
}

// GetAuthor returns author by id
func (s *AuthorService) GetAuthor(ctx context.Context, authorID types.ID) (*model.Author, error) {
	s.log.Dbg().Ctx(ctx).Values("authorID", authorID).Msg("GetAuthor")

	return s.reader.GetAuthor(ctx, authorID)
}

// CreateAuthor inserts new author
func (s *AuthorService) CreateAuthor(ctx context.Context, author *model.Author) (*model.Author, error) {
	s.log.Dbg().Ctx(ctx).Values("author", author).Msg("CreateAuthorReq")

	return s.writer.CreateAuthor(ctx, author)
}

// UpdateAuthor updates author by id
func (s *AuthorService) UpdateAuthor(ctx context.Context, authorID types.ID, author *model.Author) error {
	s.log.Dbg().Ctx(ctx).Values("authorID", authorID, "author", author).Msg("UpdateAuthorReq")

	return s.writer.UpdateAuthor(ctx, authorID, author)
}

// DeleteAuthor deletes author by id
func (s *AuthorService) DeleteAuthor(ctx context.Context, authorID types.ID) error {
	s.log.Dbg().Ctx(ctx).Values("authorID", authorID).Msg("DeleteAuthor")

	return s.writer.DeleteAuthor(ctx, authorID)
}
