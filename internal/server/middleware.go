package server

import (
	"github.com/Tihmmm/mr-decorator/pkg"
	"github.com/labstack/echo/v4"
	"net/http"
)

var apiKeyHash string

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		apiKeyIn := ctx.Request().Header.Get("Api-Key")
		if !pkg.CheckArgonHash(apiKeyIn, apiKeyHash) {
			return ctx.String(http.StatusUnauthorized, "API Key is invalid")
		}

		return next(ctx)
	}
}
