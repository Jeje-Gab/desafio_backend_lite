package usecase

import (
	"context"
	"desafio_backend/internal/entity"
	"desafio_backend/internal/negociacoes"
	"time"
)

// service implementa negociacoes.Service.
type service struct {
	getStats *GetStatsUC
}

func NewService(negociosRepo negociacoes.Repository) negociacoes.Service {
	return &service{
		getStats: NewGetStatsUCsitory(negociosRepo),
	}
}

func (s *service) GetStats(ctx context.Context, ticker string, from *time.Time) (entity.Stats, error) {
	return s.getStats.Execute(ctx, ticker, from)
}
