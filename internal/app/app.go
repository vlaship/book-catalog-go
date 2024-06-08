package app

import (
	"github.com/vlaship/book-catalog-go/internal/app/controller"
	"github.com/vlaship/book-catalog-go/internal/app/facade"
	"github.com/vlaship/book-catalog-go/internal/app/repository"
	"github.com/vlaship/book-catalog-go/internal/app/service"
	"github.com/vlaship/book-catalog-go/internal/authentication"
	"github.com/vlaship/book-catalog-go/internal/cache"
	"github.com/vlaship/book-catalog-go/internal/config"
	"github.com/vlaship/book-catalog-go/internal/database"
	"github.com/vlaship/book-catalog-go/internal/email"
	"github.com/vlaship/book-catalog-go/internal/httphandling"
	"github.com/vlaship/book-catalog-go/internal/logger"
	"github.com/vlaship/book-catalog-go/internal/router"
	"github.com/vlaship/book-catalog-go/internal/snowflake"
	"github.com/vlaship/book-catalog-go/internal/template"
	"github.com/vlaship/book-catalog-go/internal/validation"

	"github.com/go-chi/chi/v5"
)

// App struct holds the dependencies for the application.
type App struct {
	DB     database.ConnPool
	Router *chi.Mux
}

// NewApp creates a new instance of the App with provided configurations.
func NewApp(cfg *config.Config, log logger.Logger) (*App, error) {
	// init db pool
	log.Trc().Msg("init database")
	pool, err := database.New(cfg, log)
	if err != nil {
		return nil, err
	}

	// init authenticator
	log.Trc().Msg("init authenticator")
	authenticator := authentication.New(cfg)

	// init template
	log.Trc().Msg("init templates")
	templates, err := template.NewTemplatesImpl()
	if err != nil {
		return nil, err
	}

	// init validator
	log.Trc().Msg("init validator")
	validator := validation.New()

	// init repositories
	log.Trc().Msg("init repositories")
	repos := repository.Wire(pool, log)

	// init sender
	log.Trc().Msg("init email sender")
	sender := email.New(cfg)

	// init cache
	log.Trc().Msg("init cache")
	caches := cache.New()

	// init ID generator
	log.Trc().Msg("init ID generator")
	idGen, err := snowflake.New(cfg.SnowflakeNode)
	if err != nil {
		return nil, err
	}

	// init services
	log.Trc().Msg("init services")
	services := service.Wire(cfg, repos, authenticator, templates, sender, caches, idGen, log)

	// init facades
	log.Trc().Msg("init facades")
	facades := facade.Wire(services, log)

	// init http error handler
	log.Trc().Msg("init http error handler")
	httpErrorHandler := httphandling.New(log)

	// init controllers
	log.Trc().Msg("init controllers")
	controllers := controller.Wire(facades, validator, httpErrorHandler, log)

	// init router
	log.Trc().Msg("init router")
	webRouter := router.Setup(controllers, log, repos.UserRepository, authenticator, httpErrorHandler)

	// create new App instance.
	app := &App{
		DB:     pool,
		Router: webRouter,
	}

	return app, nil
}
