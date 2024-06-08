//go:build wireinject
// +build wireinject

package facade

import (
	"github.com/google/wire"
	"github.com/vlaship/book-catalog-go/internal/app/service"
	"github.com/vlaship/book-catalog-go/internal/logger"
)

func Wire(services *service.Services, log logger.Logger) *Facades {
	wire.Build(
		NewBookFacade,
		NewAuthorFacade,
		NewAuthFacade,
		NewUserFacade,
		AuthProvider,
		BookReaderProvider,
		BookWriterProvider,
		AuthorReaderProvider,
		AuthorWriterProvider,
		MailSenderProvider,
		TokenHandlerProvider,
		UserReaderProvider,
		UserWriterProvider,
		wire.Struct(new(Facades), "*"),
	)
	return &Facades{}
}

// AuthProvider is a provider for Auth
func AuthProvider(services *service.Services) Auth {
	return services.AuthService
}

// MailSenderProvider is a provider for MailSender
func MailSenderProvider(services *service.Services) MailSender {
	return services.SendMailService
}

// TokenHandlerProvider is a provider for TokenHandler
func TokenHandlerProvider(services *service.Services) TokenHandler {
	return services.OTPService
}

// UserReaderProvider is a provider for UserReader
func UserReaderProvider(services *service.Services) UserReader {
	return services.UserService
}

// UserWriterProvider is a provider for UserWriter
func UserWriterProvider(services *service.Services) UserWriter {
	return services.UserService
}

// BookReaderProvider is a provider for BookReader
func BookReaderProvider(services *service.Services) BookReader {
	return services.BookService
}

// BookWriterProvider is a provider for BookWriter
func BookWriterProvider(services *service.Services) BookWriter {
	return services.BookService
}

// AuthorReaderProvider is a provider for AuthorReader
func AuthorReaderProvider(services *service.Services) AuthorReader {
	return services.AuthorService
}

// AuthorWriterProvider is a provider for AuthorWriter
func AuthorWriterProvider(services *service.Services) AuthorWriter {
	return services.AuthorService
}
