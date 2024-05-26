package middleware

import (
	"book-catalog/internal/apperr"
	"book-catalog/internal/httphandling"
	"fmt"
	"net/http"
	"strings"
)

// AllowContentType enforces a whitelist of request Content-Types otherwise responds
// with a 415 Unsupported Media Type status.
func AllowContentType(contentTypes ...string) func(next http.Handler) http.Handler {
	allowedContentTypes := make(map[string]struct{}, len(contentTypes))
	for _, contentType := range contentTypes {
		allowedContentTypes[strings.TrimSpace(strings.ToLower(contentType))] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				// skip check for empty content body
				next.ServeHTTP(w, r)
				return
			}

			s := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
			if i := strings.Index(s, ";"); i > -1 {
				s = s[0:i]
			}

			if _, ok := allowedContentTypes[s]; ok {
				next.ServeHTTP(w, r)
				return
			}

			p := apperr.ErrUnsupportedMediaType.WithFunc(
				apperr.WithDetail(fmt.Sprintf("Content-Type must be one of: %s", strings.Join(contentTypes, ", "))),
			)
			httphandling.AppErrorResponse(w, r, p)
		}

		return http.HandlerFunc(fn)
	}
}
