package http

import (
	"desafio_backend/internal/negociacoes"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes associa endpoints do módulo negociações ao router.
func RegisterRoutes(g *echo.Group, svc negociacoes.Service) {
	handler := NewStatsHandler(svc)
	g.GET("/stats", handler.GetStats)
}
