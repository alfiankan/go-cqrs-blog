package common

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SecureMiddleware() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection: "1; mode=block",
	})
}

func CompressMiddleware() echo.MiddlewareFunc {
	return middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	})
}

var MiddlewaresRegistry = []echo.MiddlewareFunc{
	SecureMiddleware(),
	CompressMiddleware(),
}
