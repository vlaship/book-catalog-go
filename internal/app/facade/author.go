package facade

import (
	"book-catalog/internal/app/dto/request"
	"book-catalog/internal/app/dto/response"
	"book-catalog/internal/app/mapper"
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/logger"
	"context"
)

// AuthorReader is an interface for author reader
//
//go:generate mockgen -destination=../../../test/mock/facade/mock-author-reader.go -package=mock . AuthorReader
type AuthorReader interface {
	GetAuthors(ctx context.Context) ([]model.Author, error)
	GetAuthor(ctx context.Context, authorID types.ID) (*model.Author, error)
}

// AuthorWriter is an interface for author writer
//
//go:generate mockgen -destination=../../../test/mock/facade/mock-author-writer.go -package=mock . AuthorWriter
type AuthorWriter interface {
	CreateAuthor(ctx context.Context, author *model.Author) (*model.Author, error)
	UpdateAuthor(ctx context.Context, authorID types.ID, author *model.Author) error
	DeleteAuthor(ctx context.Context, authorID types.ID) error
}

type AuthorFacade struct {
	reader AuthorReader
	writer AuthorWriter
	m      mapper.Author
	log    logger.Logger
}

// NewAuthorFacade creates new author facade
func NewAuthorFacade(reader AuthorReader, writer AuthorWriter, log logger.Logger) *AuthorFacade {
	return &AuthorFacade{
		reader: reader,
		writer: writer,
		m:      mapper.Author{},
		log:    log.New("AuthorFacade"),
	}
}

// GetAuthors returns all authors
func (f *AuthorFacade) GetAuthors(ctx context.Context) ([]response.ListAuthor, error) {
	f.log.Trc().Ctx(ctx).Msg("GetAuthors")

	authors, err := f.reader.GetAuthors(ctx)
	if err != nil {
		return nil, err
	}

	return f.m.AuthorsResp(authors), nil
}

// GetAuthor returns author by id
func (f *AuthorFacade) GetAuthor(ctx context.Context, authorID types.ID) (*response.Author, error) {
	f.log.Dbg().Ctx(ctx).Msg("GetAuthor")

	author, err := f.reader.GetAuthor(ctx, authorID)
	if err != nil {
		return nil, err
	}

	return f.m.AuthorResp(author), nil
}

// CreateAuthor inserts new author
func (f *AuthorFacade) CreateAuthor(ctx context.Context, author *request.CreateAuthor) (*response.CreateAuthor, error) {
	f.log.Dbg().Ctx(ctx).Values("author", author).Msg("CreateAuthorReq")

	req := f.m.CreateAuthorReq(author)

	resp, err := f.writer.CreateAuthor(ctx, req)
	if err != nil {
		return nil, err
	}

	return f.m.CreateAuthorResp(resp), nil
}

// UpdateAuthor updates author
func (f *AuthorFacade) UpdateAuthor(ctx context.Context, authorID types.ID, author *request.UpdateAuthor) error {
	f.log.Dbg().Ctx(ctx).Values("author", author).Msg("UpdateAuthorReq")

	return f.writer.UpdateAuthor(ctx, authorID, f.m.UpdateAuthorReq(author))
}

// DeleteAuthor deletes author
func (f *AuthorFacade) DeleteAuthor(ctx context.Context, authorID types.ID) error {
	f.log.Dbg().Ctx(ctx).Values("authorID", authorID).Msg("DeleteAuthor")

	return f.writer.DeleteAuthor(ctx, authorID)
}
