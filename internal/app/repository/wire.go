//go:build wireinject
// +build wireinject

package repository

import (
	"github.com/vlaship/book-catalog-go/internal/database"
	"github.com/vlaship/book-catalog-go/internal/logger"

	"github.com/google/wire"
)

// Wire creates the repository instances
func Wire(pool database.ConnPool, log logger.Logger) *Repositories {
	wire.Build(
		NewBookRepository,
		NewAuthorRepository,
		NewPropertyRepository,
		NewUserRepository,
		wire.Struct(new(Repositories), "*"),
	)
	return &Repositories{}
}
