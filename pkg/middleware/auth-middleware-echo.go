package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (cre *Credential) AuthMiddlewareEcho(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", -1)

		err := cre.VerifyToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		if cre.TokenExpired(token) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}

		return next(c)
	}
}
