package negociacoes

import "github.com/labstack/echo/v4"

// Delivery expõe o registro de rotas do módulo.
type Delivery interface {
	RegisterRoutes(g *echo.Group, svc Service)
}
