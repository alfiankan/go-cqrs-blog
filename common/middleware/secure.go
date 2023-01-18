package http_delivery

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SecureMiddleware() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection: "1; mode=block",
	})
}
