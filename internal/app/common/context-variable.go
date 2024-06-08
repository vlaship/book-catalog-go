package common

import (
	"context"
	"github.com/vlaship/book-catalog-go/internal/app/model"
	"github.com/vlaship/book-catalog-go/internal/app/types"

	"github.com/go-chi/chi/v5/middleware"
)

func GetUser(ctx context.Context) *model.User {
	return ctx.Value(types.UserContextKey).(*model.User)
}

func GetRequestID(ctx context.Context) any {
	return ctx.Value(middleware.RequestIDKey)
}
