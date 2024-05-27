package email

import "html/template"

// Sender interface
//
//go:generate mockgen -destination=../../test/mock/email/mock-sender.go -package=mock . Sender
type Sender interface {
	Send(to []string, subj string, template *template.Template, placeHolders any) error
}
