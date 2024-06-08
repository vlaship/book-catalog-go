//go:build wireinject
// +build wireinject

package controller

import (
	"github.com/google/wire"
	"github.com/vlaship/book-catalog-go/internal/app/facade"
	"github.com/vlaship/book-catalog-go/internal/httphandling"
	"github.com/vlaship/book-catalog-go/internal/logger"
	"github.com/vlaship/book-catalog-go/internal/validation"
)

func Wire(
	facades *facade.Facades,
	validator validation.Validator,
	handler httphandling.HTTPErrorHandler,
	log logger.Logger,
) *Controllers {
	wire.Build(
		NewAuthController,
		NewUserController,
		NewBookController,
		NewAuthorController,
		AuthProvider,
		UserReaderProvider,
		UserWriterProvider,
		ActivatorProvider,
		PasswordResetHandlerProvider,
		BookReaderProvider,
		BookWriterProvider,
		AuthorReaderProvider,
		AuthorWriterProvider,
		wire.Struct(new(Controllers), "*"),
	)
	return &Controllers{}
}

// AuthProvider is a provider for Auth
func AuthProvider(facades *facade.Facades) Auth {
	return facades.AuthFacade
}

// UserReaderProvider is a provider for UserReader
func UserReaderProvider(facades *facade.Facades) UserReader {
	return facades.UserFacade
}

// UserWriterProvider is a provider for UserWriter
func UserWriterProvider(facades *facade.Facades) UserWriter {
	return facades.UserFacade
}

// ActivatorProvider is a provider for Activator
func ActivatorProvider(facades *facade.Facades) Activator {
	return facades.AuthFacade
}

// PasswordResetHandlerProvider is a provider for PasswordResetHandler
func PasswordResetHandlerProvider(facades *facade.Facades) PasswordResetHandler {
	return facades.AuthFacade
}

// BookReaderProvider is a provider for BookReader
func BookReaderProvider(facades *facade.Facades) BookReader {
	return facades.BookFacade
}

// BookWriterProvider is a provider for BookWriter
func BookWriterProvider(facades *facade.Facades) BookWriter {
	return facades.BookFacade
}

// AuthorReaderProvider is a provider for AuthorReader
func AuthorReaderProvider(facades *facade.Facades) AuthorReader {
	return facades.AuthorFacade
}

// AuthorWriterProvider is a provider for AuthorWriter
func AuthorWriterProvider(facades *facade.Facades) AuthorWriter {
	return facades.AuthorFacade
}
