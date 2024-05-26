package facade

import (
	"book-catalog/internal/app/dto/request"
	"book-catalog/internal/app/dto/response"
	"book-catalog/internal/app/mapper"
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/apperr"
	"book-catalog/internal/logger"
	"context"
	"errors"
)

// Auth interface
//
//go:generate mockgen -destination=../../../test/mock/facade/mock-auth.go -package=mock . Auth
type Auth interface {
	Signin(ctx context.Context, signin model.User) (*model.Signin, error)
	Signup(ctx context.Context, input model.User) (*model.User, error)
}

// MailSender is an interface for mail sender
//
//go:generate mockgen -destination=../../../test/mock/facade/mock-mail-sender.go -package=mock . MailSender
type MailSender interface {
	SendActivationMail(email types.Username, otp types.Token) error
	SendResetPasswordMail(email types.Username, otp types.Token) error
}

// TokenHandler is an interface for token generator
//
//go:generate mockgen -destination=../../../test/mock/facade/mock-token-generator.go -package=mock . TokenHandler
type TokenHandler interface {
	GenerateActivationOTP(ctx context.Context, username types.Username) types.Token
	GenerateResetPasswordOTP(ctx context.Context, username types.Username) types.Token
	ValidateActivationOTP(ctx context.Context, otp types.Token) (types.Username, error)
	ValidateResetPasswordOTP(ctx context.Context, otp types.Token) (types.Username, error)
}

// AuthFacade is a facade for authentication.
type AuthFacade struct {
	auth   Auth
	sender MailSender
	th     TokenHandler
	ur     UserReader
	uw     UserWriter
	m      mapper.Auth
	log    logger.Logger
}

// NewAuthFacade creates a new AuthFacade instance.
func NewAuthFacade(
	auth Auth,
	sender MailSender,
	th TokenHandler,
	ur UserReader,
	uw UserWriter,
	log logger.Logger,
) *AuthFacade {
	return &AuthFacade{
		auth:   auth,
		sender: sender,
		th:     th,
		ur:     ur,
		uw:     uw,
		m:      mapper.Auth{},
		log:    log.New("AuthFacade"),
	}
}

// Signin logging in
func (f *AuthFacade) Signin(ctx context.Context, req *request.Signin) (*response.Signin, error) {
	f.log.Dbg().Ctx(ctx).Values("req", req.String()).Msg("Signin")

	user := f.m.Signin.Model(req)
	out, err := f.auth.Signin(ctx, user)
	if err != nil {
		return nil, err
	}
	resp := f.m.Signin.Resp(out)

	return &resp, nil
}

// Signup signing up
func (f *AuthFacade) Signup(ctx context.Context, req *request.Signup) error {
	f.log.Dbg().Ctx(ctx).Values("req", req.String()).Msg("Signup")

	input := f.m.Signup.Model(req)
	user, err := f.auth.Signup(ctx, input)
	if err != nil {
		return err
	}

	return f.sendActivationMail(ctx, user.Username)
}

// Activate activating user
func (f *AuthFacade) Activate(ctx context.Context, req *request.Activation) error {
	f.log.Dbg().Ctx(ctx).Values("req", req).Msg("Activate")

	username, err := f.th.ValidateActivationOTP(ctx, req.OTP)
	if err != nil {
		return err
	}

	u := model.User{
		Username: username,
		Data:     model.UserData{Status: model.UserStatusActive},
	}

	return f.uw.Activate(ctx, u)
}

// Resend resending activation mail
func (f *AuthFacade) Resend(ctx context.Context, req *request.ResendActivation) error {
	f.log.Dbg().Ctx(ctx).Values("req", req).Msg("Resend")

	user, err := f.ur.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil
		}
		return err
	}

	return f.sendActivationMail(ctx, user.Username)
}

// Reset resetting password
func (f *AuthFacade) Reset(ctx context.Context, req *request.ResetPassword) error {
	f.log.Dbg().Ctx(ctx).Values("req", req.String()).Msg("Reset")

	user, err := f.ur.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil
		}
		return err
	}

	otp := f.th.GenerateResetPasswordOTP(ctx, user.Username)
	return f.sender.SendResetPasswordMail(user.Username, otp)
}

// Replace replacing password
func (f *AuthFacade) Replace(ctx context.Context, req *request.ReplacePassword) error {
	f.log.Dbg().Ctx(ctx).Values("req", req.String()).Msg("Replace")

	username, err := f.th.ValidateResetPasswordOTP(ctx, req.OTP)
	if err != nil {
		return err
	}

	u := model.User{
		Username: username,
		Password: req.NewPassword,
	}

	return f.uw.UpdatePassword(ctx, u)
}

func (f *AuthFacade) sendActivationMail(ctx context.Context, username types.Username) error {
	otp := f.th.GenerateActivationOTP(ctx, username)
	return f.sender.SendActivationMail(username, otp)
}
