package service

import (
	"book-catalog/internal/app/model"
	"book-catalog/internal/logger"
	"context"
)

// TosReader is an interface for term of service reader
//
//go:generate mockgen -destination=../../../test/mock/service/mock-assumptions-reader.go -package=mock . TosReader
type TosReader interface {
	GetTos(ctx context.Context) (*model.TermOfService, error)
}

// TosService is a service for term of service
type TosService struct {
	reader TosReader
	log    logger.Logger
}

// NewTosService creates new term of service
func NewTosService(reader TosReader, log logger.Logger) *TosService {
	return &TosService{
		reader: reader,
		log:    log.New("TosService"),
	}
}

// GetTos get term of service
func (s *TosService) GetTos(ctx context.Context) (*model.TermOfService, error) {
	s.log.Trc().Ctx(ctx).Msg("GetTos")

	return s.reader.GetTos(ctx)
}
