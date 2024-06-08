package repository

import (
	"context"
	"github.com/vlaship/book-catalog-go/internal/app/model"
	"github.com/vlaship/book-catalog-go/internal/database"
	"github.com/vlaship/book-catalog-go/internal/logger"
)

// PropertyRepository is an interface for property repository
type PropertyRepository struct {
	pool database.ConnPool
	log  logger.Logger
}

// NewPropertyRepository creates new property repository
func NewPropertyRepository(pool database.ConnPool, log logger.Logger) *PropertyRepository {
	return &PropertyRepository{
		pool: pool,
		log:  log.New("PropertyRepository"),
	}
}

const getPropertiesQuery = `
	SELECT property_value FROM catalog.properties WHERE property_name = $1;
`

// Repo functions
func (r *PropertyRepository) l() logger.Logger {
	return r.log
}

func (r *PropertyRepository) p() database.ConnPool {
	return r.pool
}

// GetTos get term of service
func (r *PropertyRepository) GetTos(ctx context.Context) (*model.TermOfService, error) {
	return getProperties(ctx, r, getProperty[model.TermOfService]{
		propertyName: "tos",
	})
}

func getProperties[T model.Property](
	ctx context.Context,
	r *PropertyRepository,
	q getProperty[T],
) (*T, error) {

	req := entity[T]{
		query:        getPropertiesQuery,
		entityName:   q.propertyName,
		args:         []any{q.propertyName},
		destinations: func(property *T) []any { return []any{property} },
	}

	return getOne(ctx, r, req)
}
