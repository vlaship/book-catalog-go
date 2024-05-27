package mapper

import (
	"book-catalog/internal/app/dto/request"
	"book-catalog/internal/app/dto/response"
	"book-catalog/internal/app/model"
)

// User mapper
type User struct{}

// Model creates a new user model
func (m *User) Model(req *request.UserData) model.User {
	return model.User{
		Data: model.UserData{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		},
	}
}

// Resp creates a new user response
func (m *User) Resp(out *model.User) response.User {
	return response.User{
		Username: out.Username,
		Info: response.UserInfo{
			FirstName: out.Data.FirstName,
			LastName:  out.Data.LastName,
			Email:     out.Data.Email,
		},
	}
}
