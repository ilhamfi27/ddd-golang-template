package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

// StructuredLoggingMiddleware creates a middleware for structured logging
func StructuredLoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			// Log incoming request
			slog.InfoContext(req.Context(), "incoming request",
				slog.String("method", req.Method),
				slog.String("path", req.RequestURI),
				slog.String("ip", c.RealIP()),
				slog.String("user_agent", req.UserAgent()),
				slog.String("content_type", req.Header.Get("Content-Type")),
			)

			// Process request
			err := next(c)

			// Calculate duration
			duration := time.Since(start)

			// Log response
			slog.InfoContext(req.Context(), "request completed",
				slog.String("method", req.Method),
				slog.String("path", req.RequestURI),
				slog.Int("status", res.Status),
				slog.Duration("duration", duration),
				slog.Int64("bytes_sent", res.Size),
			)

			// Log errors if any
			if err != nil {
				slog.ErrorContext(req.Context(), "request error",
					slog.String("method", req.Method),
					slog.String("path", req.RequestURI),
					slog.String("error", err.Error()),
					slog.Duration("duration", duration),
				)
			}

			return err
		}
	}
}
