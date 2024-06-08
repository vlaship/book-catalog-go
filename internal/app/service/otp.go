package service

import (
	"context"
	"github.com/vlaship/book-catalog-go/internal/app/types"
	"github.com/vlaship/book-catalog-go/internal/apperr"
	"github.com/vlaship/book-catalog-go/internal/cache"
	"github.com/vlaship/book-catalog-go/internal/logger"
	"github.com/vlaship/go-mask"
	"github.com/vlaship/go-otp"
	"time"
)

const (
	ttlActivation = 48 * time.Hour
	ttlReset      = 1 * time.Hour
)

// OTPService is a service for token.
type OTPService struct {
	cacher cache.Cache
	log    logger.Logger
}

// NewOTPService creates a new OTPService instance.
func NewOTPService(
	cacher cache.Cache,
	log logger.Logger,
) *OTPService {
	return &OTPService{
		cacher: cacher,
		log:    log.New("OTPService"),
	}
}

// GenerateActivationOTP generates activation otp
func (s *OTPService) GenerateActivationOTP(ctx context.Context, username types.Username) types.Token {
	s.log.Trc().Ctx(ctx).Values("username", mask.String(string(username))).Msg("GenerateActivationOTP")

	o := otp.Generate()
	s.cacher.Put(o, username, ttlActivation)

	return types.Token(o)
}

// ValidateActivationOTP validates activation otp
func (s *OTPService) ValidateActivationOTP(ctx context.Context, otp types.Token) (types.Username, error) {
	s.log.Trc().Ctx(ctx).Values("otp", otp).Msg("ValidateActivationOTP")

	username, ok := s.cacher.GetDel(string(otp))
	if !ok {
		return "", apperr.ErrInvalidOTP
	}

	return username.(types.Username), nil
}

func (s *OTPService) GenerateResetPasswordOTP(ctx context.Context, username types.Username) types.Token {
	s.log.Trc().Ctx(ctx).Values("username", mask.String(string(username))).Msg("GenerateResetPasswordOTP")

	o := otp.Generate()
	s.cacher.Put(o, username, ttlReset)

	return types.Token(o)
}

func (s *OTPService) ValidateResetPasswordOTP(ctx context.Context, otp types.Token) (types.Username, error) {
	s.log.Trc().Ctx(ctx).Values("otp", otp).Msg("ValidateResetPasswordOTP")

	username, ok := s.cacher.GetDel(string(otp))
	if !ok {
		return "", apperr.ErrInvalidOTP
	}

	return username.(types.Username), nil
}
