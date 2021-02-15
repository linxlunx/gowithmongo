package wrapper

import (
	"github.com/labstack/echo/v4"
)

// Props is base struct for response
type Props struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error is called when there is an error logic or null value
func Error(code int, message string, c echo.Context) error {
	props := &Props{
		Code:    code,
		Message: message,
	}
	return c.JSON(code, props)
}
