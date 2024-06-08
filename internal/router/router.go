package router

import (
	"compress/gzip"
	_ "github.com/vlaship/book-catalog-go/api/docs" // swagger docs
	"github.com/vlaship/book-catalog-go/internal/app/controller"
	"github.com/vlaship/book-catalog-go/internal/authentication"
	"github.com/vlaship/book-catalog-go/internal/httphandling"
	"github.com/vlaship/book-catalog-go/internal/logger"
	mw "github.com/vlaship/book-catalog-go/internal/router/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	sw "github.com/swaggo/http-swagger/v2"
)

// Setup sets up the router for the application.
// @title Book Catalog API
// @version v1
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func Setup(
	controllers *controller.Controllers,
	log logger.Logger,
	userReader mw.UserReader,
	authenticator authentication.Authenticator,
	handler httphandling.HTTPErrorHandler,
) *chi.Mux {
	log.Trc().Msg("setup router")
	basePath := "/api"

	// init router
	r := chi.NewRouter()

	// init middleware
	l := log.New("router")
	r.Use(httplog.RequestLogger(l.Logger()))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(gzip.DefaultCompression))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.GetHead)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(mw.CorsHandler())
	r.Use(middleware.StripSlashes)
	r.Use(mw.NewContentTypeMiddleware(handler).AllowContentType("application/json"))

	r.Route(basePath, func(baseRouter chi.Router) {
		baseRouter.Group(func(authRouter chi.Router) {
			// auth validation
			authRouter.Use(mw.NewAuthMiddleware(authenticator, userReader, handler, log).Validation())

			// endpoints
			controllers.AuthorController.RegisterRoutes(authRouter)
			controllers.BookController.RegisterRoutes(authRouter)
			controllers.UserController.RegisterRoutes(authRouter)
		})
		// register auth
		controllers.AuthController.RegisterRoutes(baseRouter)
	})

	// register swagger
	r.Get("/swagger-ui*", sw.WrapHandler)

	// redirect from / to /swagger-ui
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger-ui/", http.StatusSeeOther)
	})

	// print routes
	printRoutes(r, log)

	return r
}

func printRoutes(routes chi.Routes, log logger.Logger) {
	log.Inf().Msg("routes")
	_ = chi.Walk(
		routes,
		func(
			method string,
			route string,
			_ http.Handler,
			_ ...func(http.Handler) http.Handler,
		) error {
			log.Inf().Msg("%s %s", method, route)
			return nil
		},
	)
}
