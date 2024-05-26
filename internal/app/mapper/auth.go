package mapper

import (
	"book-catalog/internal/app/dto/request"
	"book-catalog/internal/app/dto/response"
	"book-catalog/internal/app/model"
)

// Auth is an interface for auth mapper
type Auth struct {
	Signin
	Signup
}

type Signin struct{}

type Signup struct{}

// Model creates new auth mapper
func (Signin) Model(in *request.Signin) model.User {
	return model.User{
		Username: in.Username,
		Password: in.Password,
	}
}

// Resp creates new auth mapper
func (Signin) Resp(out *model.Signin) response.Signin {
	return response.Signin{
		AccessToken:  out.AccessToken,
		Type:         "Bearer",
		ExpiresIn:    out.ExpiresIn,
		RefreshToken: out.RefreshToken,
	}
}

// Model creates new auth mapper
func (Signup) Model(in *request.Signup) model.User {
	return model.User{
		Username: in.Username,
		Password: in.Password,
		Data: model.UserData{
			FirstName: in.Firstname,
			LastName:  in.Lastname,
			Status:    model.UserStatusNonActive,
		},
	}
}
