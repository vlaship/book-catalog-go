package common

import (
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"context"

	"github.com/go-chi/chi/v5/middleware"
)

func GetUser(ctx context.Context) *model.User {
	return ctx.Value(types.UserContextKey).(*model.User)
}

func GetRequestID(ctx context.Context) any {
	return ctx.Value(middleware.RequestIDKey)
}
