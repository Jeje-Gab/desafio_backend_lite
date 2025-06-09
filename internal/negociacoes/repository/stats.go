package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"desafio_backend/internal/entity"
)

const (
	basePriceQuery = `
        SELECT preco_negocio
          FROM negociacoes.negociacoes
         WHERE codigo_instrumento = $1
         ORDER BY preco_negocio DESC
         LIMIT 1
    `
	basePriceQueryWithFrom = `
        SELECT preco_negocio
          FROM negociacoes.negociacoes
         WHERE codigo_instrumento = $1
           AND data_negocio >= $2
         ORDER BY preco_negocio DESC
         LIMIT 1
    `

	baseVolumeQuery = `
        SELECT daily_sum
          FROM negociacoes.negociacoes_daily_volume
         WHERE codigo_instrumento = $1
         ORDER BY daily_sum DESC
         LIMIT 1
    `

	// Busca do maior volume diário a partir de uma data mínima
	baseVolumeQueryWithFrom = `
        SELECT daily_sum
          FROM negociacoes.negociacoes_daily_volume
         WHERE codigo_instrumento = $1
           AND data_negocio       >= $2
         ORDER BY daily_sum DESC
         LIMIT 1
    `
)

// GetStatsRepo executa as queries de agregação
// mantém o *sql.DB pra fallback e cacheia os *sql.Stmt.
type GetStatsRepo struct {
	db                 *sql.DB
	mu                 sync.Mutex
	priceStmtNoFrom    *sql.Stmt
	priceStmtWithFrom  *sql.Stmt
	volumeStmtNoFrom   *sql.Stmt
	volumeStmtWithFrom *sql.Stmt
}

// NewGetStatsRepository continua sem error
func NewGetStatsRepository(db *sql.DB) *GetStatsRepo {
	return &GetStatsRepo{db: db}
}

// Execute retorna as estatísticas; faz o prepare on-demand
func (r *GetStatsRepo) Execute(ctx context.Context, ticker string, from *time.Time) (entity.Stats, error) {
	var stats entity.Stats

	// garante que as statements estão preparadas
	if err := r.prepareStatements(ctx); err != nil {
		return stats, fmt.Errorf("prepare statements: %w", err)
	}

	// escolhe o stmt certo
	var row *sql.Row
	if from != nil {
		row = r.priceStmtWithFrom.QueryRowContext(ctx, ticker, *from)
	} else {
		row = r.priceStmtNoFrom.QueryRowContext(ctx, ticker)
	}
	if err := row.Scan(&stats.MaxPrice); err != nil {
		return stats, fmt.Errorf("scan MaxPrice: %w", err)
	}

	if from != nil {
		row = r.volumeStmtWithFrom.QueryRowContext(ctx, ticker, *from)
	} else {
		row = r.volumeStmtNoFrom.QueryRowContext(ctx, ticker)
	}
	if err := row.Scan(&stats.MaxDailyVolume); err != nil {
		return stats, fmt.Errorf("scan MaxDailyVolume: %w", err)
	}

	return stats, nil
}

// prepareStatements faz PrepareContext só na primeira vez
func (r *GetStatsRepo) prepareStatements(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.priceStmtNoFrom != nil {
		return nil // já cacheado
	}

	var err error
	// Price
	r.priceStmtNoFrom, err = r.db.PrepareContext(ctx, basePriceQuery)
	if err != nil {
		return err
	}
	r.priceStmtWithFrom, err = r.db.PrepareContext(ctx, basePriceQueryWithFrom)
	if err != nil {
		return err
	}

	// Volume
	r.volumeStmtNoFrom, err = r.db.PrepareContext(ctx, baseVolumeQuery)
	if err != nil {
		return err
	}
	r.volumeStmtWithFrom, err = r.db.PrepareContext(ctx, baseVolumeQueryWithFrom)
	if err != nil {
		return err
	}

	return nil
}

// Close opcional, para liberar os statements se você quiser
func (r *GetStatsRepo) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.priceStmtNoFrom != nil {
		r.priceStmtNoFrom.Close()
		r.priceStmtWithFrom.Close()
		r.volumeStmtNoFrom.Close()
		r.volumeStmtWithFrom.Close()
	}
}
