package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// CorsHandler Cors returns a CORS handler.
func CorsHandler() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "sentry-trace", "baggage"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
