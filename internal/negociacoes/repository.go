package negociacoes

import (
	"context"
	"desafio_backend/internal/entity"
	"time"
)

// Repository define o contrato de persistência.
type Repository interface {
	GetStats(ctx context.Context, ticker string, from *time.Time) (entity.Stats, error)
}
