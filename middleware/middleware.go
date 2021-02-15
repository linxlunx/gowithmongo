package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AppMiddleware receive app name from main router
type AppMiddleware struct {
	appName string
}

// CORS function to set rule for what we allow and what we forbid
// Allow CORS *
func (am *AppMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Server", am.appName)
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Response().Header().Set("Accept", "application/json")

		if c.Request().Method == "OPTIONS" {
			return c.String(http.StatusOK, "")
		}

		return next(c)
	}
}

// InitAppMiddleware init app for middleware
func InitAppMiddleware(appName string) *AppMiddleware {
	return &AppMiddleware{
		appName: appName,
	}
}
