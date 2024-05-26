package template

import (
	"html/template"
)

// Templates is an interface for template
//
//go:generate mockgen -destination=../../test/mock/template/mock-template.go -package=mock . Template
type Templates interface {
	Activation() *template.Template
	ResetPassword() *template.Template
}
