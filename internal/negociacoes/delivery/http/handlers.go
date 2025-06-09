package http

import (
	"net/http"
	"time"

	"desafio_backend/internal/negociacoes"
	"github.com/labstack/echo/v4"
)

// StatsHandler trata as requisições de estatísticas.
type StatsHandler struct {
	service negociacoes.Service
}

// NewStatsHandler cria um novo handler.
func NewStatsHandler(s negociacoes.Service) *StatsHandler {
	return &StatsHandler{service: s}
}

// GetStats retorna max_price e max_daily_volume em JSON.
func (h *StatsHandler) GetStats(c echo.Context) error {
	ctx := c.Request().Context()

	ticker := c.QueryParam("ticker")
	fromStr := c.QueryParam("from")
	if ticker == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ticker é obrigatório"})
	}

	// Define o filtro 'from' corretamente em local time e sem componente horário
	var from *time.Time
	var err error
	if fromStr != "" {
		dataParam, err := time.ParseInLocation("2006-01-02", fromStr, time.Local)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "data inválida"})
		}
		from = &dataParam
	} else {
		from = nil
	}

	stats, err := h.service.GetStats(ctx, ticker, from)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, stats)
}
