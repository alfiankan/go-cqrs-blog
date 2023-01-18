package common

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SecureMiddleware xss protection
func SecureMiddleware() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection: "1; mode=block",
	})
}

// CompressMiddleware compressing http
func CompressMiddleware() echo.MiddlewareFunc {
	return middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	})
}

// ServerLog log request response
func ServerLog() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339} method=${method}, uri=${uri}, status=${status} error=${error} latency=${latency_human}\n",
	})
}

// MiddlewaresRegistry hold runnable middleware
var MiddlewaresRegistry = []echo.MiddlewareFunc{
	SecureMiddleware(),
	CompressMiddleware(),
	ServerLog(),
}
