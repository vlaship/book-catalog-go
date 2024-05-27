package service

import (
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/apperr"
	"book-catalog/internal/logger"
	"book-catalog/internal/snowflake"
	"book-catalog/pkg/utils/mask"
	"context"
	"errors"
	"fmt"
)

// Authenticator interface
//
//go:generate mockgen -destination=../../../test/mock/service/mock-auth.go -package=mock . Auth
type Authenticator interface {
	GenerateAccessToken(userID types.UserID) (accessToken types.Token, expiresIn int64, err error)
}

// PasswordHandler interface
//
//go:generate mockgen -destination=../../../test/mock/service/mock-password-handler.go -package=mock . PasswordHandler
type PasswordHandler interface {
	Validate(password, hash types.Password) error
	Hash(password types.Password) (types.Password, error)
}

// AuthService is a service for authentication.
type AuthService struct {
	reader UserReader
	writer UserWriter
	auth   Authenticator
	pass   PasswordHandler
	idGen  snowflake.IDGenerator
	log    logger.Logger
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(
	reader UserReader,
	writer UserWriter,
	auth Authenticator,
	pass PasswordHandler,
	idGen snowflake.IDGenerator,
	log logger.Logger,
) *AuthService {
	return &AuthService{
		reader: reader,
		writer: writer,
		auth:   auth,
		pass:   pass,
		idGen:  idGen,
		log:    log.New("AuthService"),
	}
}

// Signin logging in
func (s *AuthService) Signin(ctx context.Context, req model.User) (*model.Signin, error) {
	s.log.Dbg().Ctx(ctx).Values("username", mask.String(string(req.Username))).Msg("Signin")

	user, err := s.reader.GetUserByUsername(ctx, req.Username)
	if err != nil {
		s.log.Wrn().Err(err).Ctx(ctx).Msg("GetUserByUsername")
		return nil, apperr.ErrUnauthorized
	}
	if err = s.pass.Validate(req.Password, user.Password); err != nil {
		s.log.Wrn().Err(err).Ctx(ctx).Msg("validatePassword")
		return nil, apperr.ErrUnauthorized
	}
	if user.Data.Status != model.UserStatusActive {
		return nil, user.GetAppError()
	}

	token, expiresIn, err := s.auth.GenerateAccessToken(user.ID)
	if err != nil {
		s.log.Err(err).Ctx(ctx).Msg("GenerateAccessToken")
		return nil, apperr.ErrInternalServerError
	}

	out := model.Signin{
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}

	return &out, nil
}

// Signup signing up
func (s *AuthService) Signup(ctx context.Context, input model.User) (*model.User, error) {
	s.log.Dbg().Ctx(ctx).Values("username", mask.String(string(input.Username))).Msg("Signup")

	hash, err := s.pass.Hash(input.Password)
	if err != nil {
		s.log.Err(err).Ctx(ctx).Msg("Failed to hash password")
		return nil, apperr.ErrInternalServerError
	}
	input.Password = hash

	input.ID = types.UserID(s.idGen.Generate())

	user, err := s.writer.Create(ctx, input)
	if err != nil {
		if errors.Is(err, apperr.ErrAlreadyExists) {
			return nil, apperr.ErrAlreadyExists.WithFunc(
				apperr.WithDetail(fmt.Sprintf("User [%s] already exists", input.Username)),
				apperr.WithTitle("User already exists"),
			)
		}
		s.log.Err(err).Ctx(ctx).Msg("failed to create user")
		return nil, apperr.ErrInternalServerError
	}

	return user, nil
}
