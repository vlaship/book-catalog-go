//go:build wireinject
// +build wireinject

package service

import (
	"book-catalog/internal/app/repository"
	"book-catalog/internal/cache"
	"book-catalog/internal/config"
	"book-catalog/internal/email"
	"book-catalog/internal/logger"
	"book-catalog/internal/snowflake"
	"book-catalog/internal/template"
	"github.com/google/wire"
)

func Wire(
	cfg *config.Config,
	repos *repository.Repositories,
	auth Authenticator,
	templates template.Templates,
	sender email.Sender,
	cacher cache.Cache,
	idGen snowflake.IDGenerator,
	log logger.Logger,
) *Services {
	wire.Build(
		NewBookService,
		NewAuthorService,
		NewTosService,
		NewAuthService,
		NewSendMailService,
		NewUserService,
		NewOTPService,
		NewPasswordService,

		BookReaderProvider,
		BookWriterProvider,
		AuthorReaderProvider,
		AuthorWriterProvider,

		UserWriterProvider,
		UserReaderProvider,
		TosReaderProvider,
		PasswordHandlerProvider,
		wire.Struct(new(Services), "*"),
	)
	return &Services{}
}

// PasswordHandlerProvider is a provider for PasswordHandler
func PasswordHandlerProvider() PasswordHandler {
	return NewPasswordService()
}

// UserReaderProvider is a provider for UserReader
func UserReaderProvider(repos *repository.Repositories) UserReader {
	return repos.UserRepository
}

// TosReaderProvider is a provider for TosReader
func TosReaderProvider(repos *repository.Repositories) TosReader {
	return repos.PropertyRepository
}

// UserWriterProvider is a provider for UserWriter
func UserWriterProvider(repos *repository.Repositories) UserWriter {
	return repos.UserRepository
}

// BookReaderProvider is a provider for BookReader
func BookReaderProvider(repos *repository.Repositories) BookReader {
	return repos.BookRepository
}

// BookWriterProvider is a provider for BookWriter
func BookWriterProvider(repos *repository.Repositories) BookWriter {
	return repos.BookRepository
}

// AuthorReaderProvider is a provider for AuthorReader
func AuthorReaderProvider(repos *repository.Repositories) AuthorReader {
	return repos.AuthorRepository
}

// AuthorWriterProvider is a provider for AuthorWriter
func AuthorWriterProvider(repos *repository.Repositories) AuthorWriter {
	return repos.AuthorRepository
}
