package repository

import (
	"context"
	"database/sql"
	"desafio_backend/internal/entity"
	"time"

	"desafio_backend/internal/negociacoes"
)

// postgresRepo implementa negociacoes.Repository usando PostgreSQL.
type postgresRepo struct {
	getStats *GetStatsRepo
}

// NewRepository retorna um repository conectado ao banco.
func NewRepository(db *sql.DB) negociacoes.Repository {
	return &postgresRepo{getStats: NewGetStatsRepository(db)}
}

func (r *postgresRepo) GetStats(ctx context.Context, ticker string, from *time.Time) (entity.Stats, error) {
	return r.getStats.Execute(ctx, ticker, from)
}
