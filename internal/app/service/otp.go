package service

import (
	"book-catalog/internal/app/types"
	"book-catalog/internal/apperr"
	"book-catalog/internal/cache"
	"book-catalog/internal/logger"
	"book-catalog/pkg/utils/generate"
	"book-catalog/pkg/utils/mask"
	"context"
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

	otp := generate.OTP()
	s.cacher.Put(otp, username, ttlActivation)

	return types.Token(otp)
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

	otp := generate.OTP()
	s.cacher.Put(otp, username, ttlReset)

	return types.Token(otp)
}

func (s *OTPService) ValidateResetPasswordOTP(ctx context.Context, otp types.Token) (types.Username, error) {
	s.log.Trc().Ctx(ctx).Values("otp", otp).Msg("ValidateResetPasswordOTP")

	username, ok := s.cacher.GetDel(string(otp))
	if !ok {
		return "", apperr.ErrInvalidOTP
	}

	return username.(types.Username), nil
}
