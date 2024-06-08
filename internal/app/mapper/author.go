package mapper

import (
	"github.com/vlaship/book-catalog-go/internal/app/dto/request"
	"github.com/vlaship/book-catalog-go/internal/app/dto/response"
	"github.com/vlaship/book-catalog-go/internal/app/model"
	"github.com/vlaship/book-catalog-go/internal/app/types"
)

// Author is a mapper for author
type Author struct{}

// CreateAuthorReq creates a new author model
func (m *Author) CreateAuthorReq(req *request.CreateAuthor) *model.Author {
	return &model.Author{
		Name: req.Name,
		Dob:  req.Dob.Time,
	}
}

// UpdateAuthorReq updates a author model
func (m *Author) UpdateAuthorReq(req *request.UpdateAuthor) *model.Author {
	return &model.Author{
		Name: req.Name,
		Dob:  req.Dob.Time,
	}
}

// CreateAuthorResp creates a new author response
func (m *Author) CreateAuthorResp(out *model.Author) *response.CreateAuthor {
	return &response.CreateAuthor{
		ID: out.ID,
	}
}

// AuthorResp creates a new author response
func (m *Author) AuthorResp(out *model.Author) *response.Author {
	return &response.Author{
		ID:   out.ID,
		Name: out.Name,
		Dob:  types.DateDay{Time: out.Dob},
	}
}

// AuthorsResp creates a new list of author response
func (m *Author) AuthorsResp(out []model.Author) []response.ListAuthor {
	authors := make([]response.ListAuthor, 0, len(out))
	for i := range out {
		authors = append(authors, response.ListAuthor{
			ID:   out[i].ID,
			Name: out[i].Name,
		})
	}
	return authors
}
