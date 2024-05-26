package email

import (
	"book-catalog/internal/apperr"
	"book-catalog/internal/config"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

const (
	mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

// SenderImpl is a send mail implementation.
type SenderImpl struct {
	auth smtp.Auth
	user string
	host string
	port uint16
}

// New sender
func New(cfg *config.Config) Sender {
	return &SenderImpl{
		user: cfg.SMTP.Username,
		host: cfg.SMTP.Host,
		port: cfg.SMTP.Port,
		auth: smtp.PlainAuth(
			"",
			cfg.SMTP.Username,
			cfg.SMTP.Password,
			cfg.SMTP.Host,
		),
	}
}

func (s *SenderImpl) Send(to []string, subj string, tmpl *template.Template, placeHolders any) error {
	body, err := s.prepareBody(subj, tmpl, placeHolders)
	if err != nil {
		return err
	}

	if err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.host, s.port),
		s.auth,
		s.user,
		to,
		body,
	); err != nil {
		return apperr.ErrSendMail.WithFunc(apperr.WithDetail(err.Error()))
	}
	return nil
}

func (s *SenderImpl) prepareBody(subj string, tmpl *template.Template, placeHolders any) ([]byte, error) {
	var body bytes.Buffer
	body.Write([]byte(subj))
	body.Write([]byte(mime))

	err := tmpl.Execute(&body, placeHolders)
	if err != nil {
		return nil, apperr.ErrExecuteTemplate.WithFunc(apperr.WithDetail(err.Error()))
	}

	return body.Bytes(), nil
}
