package dto

import (
	. "github.com/vlaship/book-catalog-go/internal/app/dto/request" //nolint:revive,stylecheck // reduce number of symbols
)

type Request interface {
	Entity | Auth | UserData
}

type Entity interface {
	CreateBook | UpdateBook | CreateAuthor | UpdateAuthor
}

type Auth interface {
	Signin | Signup | Activation | ResendActivation | ResetPassword | ChangePassword | ReplacePassword
}
