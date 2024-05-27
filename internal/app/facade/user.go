package facade

import (
	"book-catalog/internal/app/common"
	"book-catalog/internal/app/dto/request"
	"book-catalog/internal/app/dto/response"
	"book-catalog/internal/app/mapper"
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/logger"
	"context"
)

// UserReader is an interface for user reader
//
//go:generate mockgen -destination=../../../test/mock/facade/mock-user-reader.go -package=mock . UserReader
type UserReader interface {
	GetUserByUsername(ctx context.Context, username types.Username) (*model.User, error)
}

// UserWriter is an interface for user writer
//
//go:generate mockgen -destination=../../../test/mock/facade/mock-user-writer.go -package=mock . UserWriter
type UserWriter interface {
	Activate(ctx context.Context, user model.User) error
	UpdatePassword(ctx context.Context, user model.User) error
	UpdateInfo(ctx context.Context, user model.User) error
}

// UserFacade is an interface for user facade
type UserFacade struct {
	reader UserReader
	writer UserWriter
	m      mapper.User
	log    logger.Logger
}

// NewUserFacade creates a new UserFacade instance.
func NewUserFacade(
	reader UserReader,
	writer UserWriter,
	log logger.Logger,
) *UserFacade {
	return &UserFacade{
		reader: reader,
		writer: writer,
		m:      mapper.User{},
		log:    log.New("UserFacade"),
	}
}

// GetUser returns user
func (f *UserFacade) GetUser(ctx context.Context) response.User {
	f.log.Trc().Ctx(ctx).Msg("GetUser")

	user := common.GetUser(ctx)
	return f.m.Resp(user)
}

// UpdateInfo updates user
func (f *UserFacade) UpdateInfo(ctx context.Context, req *request.UserData) error {
	f.log.Dbg().Ctx(ctx).Values("user", req).Msg("UpdateInfo")

	user := f.m.Model(req)
	return f.writer.UpdateInfo(ctx, user)
}
