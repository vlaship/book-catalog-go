package service

import (
	"book-catalog/internal/app/types"
	"book-catalog/internal/apperr"
	"book-catalog/internal/config"
	"book-catalog/internal/email"
	"book-catalog/internal/logger"
	"book-catalog/internal/template"
	"book-catalog/pkg/utils/mask"
)

const (
	subjActivationMail    = "Subject: Activate Your Account\n"
	subjResetPasswordMail = "Subject: Reset Password\n" //nolint:gosec // just subject
)

// SendMailService is an interface for send mail service
type SendMailService struct {
	sender            email.Sender
	templates         template.Templates
	userActivationURL string
	resetPasswordURL  string
	log               logger.Logger
}

// NewSendMailService creates new send mail service
func NewSendMailService(
	sender email.Sender,
	templates template.Templates,
	cfg *config.Config,
	log logger.Logger,
) *SendMailService {
	return &SendMailService{
		sender:            sender,
		templates:         templates,
		userActivationURL: cfg.Domain + "/auth/activate",
		log:               log.New("SendMailService"),
	}
}

type tmpl struct {
	URL string
	OTP string
}

// SendActivationMail sends activation mail
func (s *SendMailService) SendActivationMail(to types.Username, otp types.Token) error {
	s.log.Dbg().Values("to", mask.String(string(to))).Msg("SendActivationMail")

	t := tmpl{
		URL: s.userActivationURL,
		OTP: string(otp),
	}
	err := s.sender.Send([]string{string(to)}, subjActivationMail, s.templates.Activation(), t)
	if err != nil {
		s.log.Err(err).Msg("SendActivationMail")
		return apperr.ErrSendMail
	}

	return nil
}

// SendResetPasswordMail sends reset password mail
func (s *SendMailService) SendResetPasswordMail(to types.Username, otp types.Token) error {
	s.log.Dbg().Values("to", mask.String(string(to))).Msg("SendResetPasswordMail")

	t := tmpl{
		URL: s.resetPasswordURL,
		OTP: string(otp),
	}
	err := s.sender.Send([]string{string(to)}, subjResetPasswordMail, s.templates.ResetPassword(), t)
	if err != nil {
		s.log.Err(err).Msg("SendResetPasswordMail")
		return apperr.ErrSendMail
	}

	return nil
}
