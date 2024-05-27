package template

import (
	"embed"
	"html/template"
)

const (
	activationFilename = "templates/activation_email.html"
	resetFilename      = "templates/reset_password_email.html"
)

//go:embed templates/*.html
var templateFS embed.FS

// TemplatesImpl is a template
type TemplatesImpl struct {
	activation *template.Template
	reset      *template.Template
}

// NewTemplatesImpl creates new template
func NewTemplatesImpl() (Templates, error) {
	act, err := template.ParseFS(templateFS, activationFilename)
	if err != nil {
		return nil, err
	}
	rst, err := template.ParseFS(templateFS, resetFilename)
	if err != nil {
		return nil, err
	}

	return &TemplatesImpl{
		activation: act,
		reset:      rst,
	}, nil
}

// Activation returns activation template
func (p *TemplatesImpl) Activation() *template.Template {
	return p.activation
}

// ResetPassword returns reset password template
func (p *TemplatesImpl) ResetPassword() *template.Template {
	return p.reset
}
