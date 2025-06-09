package negociacoes

import (
	"context"
	"desafio_backend/internal/entity"
	"time"
)

// Service define a lógica de negócio de negociações.
type Service interface {
	GetStats(ctx context.Context, ticker string, from *time.Time) (entity.Stats, error)
}
