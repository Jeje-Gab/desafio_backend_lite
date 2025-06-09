package usecase

import (
	"context"
	"desafio_backend/internal/negociacoes"
	"time"

	"desafio_backend/internal/entity"
)

type GetStatsUC struct {
	negociosRepo negociacoes.Repository
}

func NewGetStatsUCsitory(negociosRepo negociacoes.Repository) *GetStatsUC {
	return &GetStatsUC{negociosRepo: negociosRepo}
}

func (r *GetStatsUC) Execute(ctx context.Context, ticker string, from *time.Time) (entity.Stats, error) {
	var (
		stats entity.Stats
		err   error
	)
	stats, err = r.negociosRepo.GetStats(ctx, ticker, from)
	if err != nil {
		return stats, err
	}
	return stats, nil
}
