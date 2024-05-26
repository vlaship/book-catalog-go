package service

import (
	"book-catalog/internal/app/common"
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/apperr"
	"book-catalog/internal/logger"
	"book-catalog/pkg/utils/mask"
	"context"
)

// UserReader interface
//
//go:generate mockgen -destination=../../../test/mock/service/mock-user_reader.go -package=mock . UserReader
type UserReader interface {
	GetUserByID(ctx context.Context, userID types.UserID) (*model.User, error)
	GetUserByUsername(ctx context.Context, username types.Username) (*model.User, error)
}

// UserWriter interface
//
//go:generate mockgen -destination=../../../test/mock/service/mock-user_writer.go -package=mock . UserWriter
type UserWriter interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	UpdateInfo(ctx context.Context, user model.User, userID types.UserID) error
	UpdateStatus(ctx context.Context, user model.User) error
	UpdatePassword(ctx context.Context, user model.User) error
}

// UserService is a service for user.
type UserService struct {
	reader UserReader
	writer UserWriter
	pass   PasswordHandler
	log    logger.Logger
}

// NewUserService creates a new UserService instance.
func NewUserService(
	reader UserReader,
	writer UserWriter,
	pass PasswordHandler,
	log logger.Logger,
) *UserService {
	return &UserService{
		reader: reader,
		writer: writer,
		pass:   pass,
		log:    log.New("UserService"),
	}
}

// GetUserByID gets user by ID.
func (s *UserService) GetUserByID(ctx context.Context, userID types.UserID) (*model.User, error) {
	s.log.Dbg().Ctx(ctx).Values("userID", userID).Msg("GetUserByID")

	return s.reader.GetUserByID(ctx, userID)
}

// GetUserByUsername gets user by username.
func (s *UserService) GetUserByUsername(ctx context.Context, username types.Username) (*model.User, error) {
	s.log.Dbg().Ctx(ctx).Values("username", mask.String(string(username))).Msg("GetUserByUsername")

	return s.reader.GetUserByUsername(ctx, username)
}

// Activate activating user
func (s *UserService) Activate(ctx context.Context, user model.User) error {
	s.log.Dbg().Ctx(ctx).Values("username", mask.String(string(user.Username))).Msg("Activate")

	err := s.writer.UpdateStatus(ctx, user)
	if err != nil {
		s.log.Wrn().Err(err).Ctx(ctx).Msg("UpdateStatus")
		return apperr.ErrInternalServerError
	}

	return nil
}

// UpdatePassword replaces user password.
func (s *UserService) UpdatePassword(ctx context.Context, user model.User) error {
	s.log.Dbg().Ctx(ctx).Values("username", mask.String(string(user.Username))).Msg("UpdatePassword")

	hash, err := s.pass.Hash(user.Password)
	if err != nil {
		s.log.Wrn().Err(err).Ctx(ctx).Msg("Hash")
		return apperr.ErrInternalServerError
	}
	user.Password = hash

	return s.writer.UpdatePassword(ctx, user)
}

// UpdateInfo updates user info.
func (s *UserService) UpdateInfo(ctx context.Context, user model.User) error {
	userID := common.GetUser(ctx).ID
	s.log.Dbg().Ctx(ctx).Values("user", user, "userID", userID).Msg("UpdateInfo")

	return s.writer.UpdateInfo(ctx, user, userID)
}
