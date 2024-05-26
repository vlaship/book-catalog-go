//go:build wireinject
// +build wireinject

package repository

import (
	"book-catalog/internal/database"
	"book-catalog/internal/logger"

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
