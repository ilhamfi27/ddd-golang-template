package middleware

import (
	"context"

	config "github.com/ilhamfi27/ddd-golang-template/internal/config/app"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware creates a middleware for distributed tracing
func TracingMiddleware() echo.MiddlewareFunc {
	tracer := otel.Tracer(config.GetEnvVars().ServiceName)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			// Start a span for this request
			ctx, span := tracer.Start(req.Context(), req.URL.Path,
				trace.WithAttributes(
					attribute.String("http.method", req.Method),
					attribute.String("http.url", req.RequestURI),
					attribute.String("http.client_ip", c.RealIP()),
					attribute.String("http.user_agent", req.UserAgent()),
				),
			)
			defer span.End()

			// Update request context with trace context
			c.SetRequest(req.WithContext(ctx))

			// Process request
			err := next(c)

			// Record response attributes
			span.SetAttributes(
				attribute.Int("http.status_code", c.Response().Status),
				attribute.Int64("http.response_size", c.Response().Size),
			)

			// Record error if any
			if err != nil {
				span.RecordError(err)
			}

			return err
		}
	}
}

// CreateChildSpan creates a child span within the current trace context
func CreateChildSpan(c echo.Context, operationName string) (context.Context, trace.Span) {
	tracer := otel.Tracer("ddd-golang-template")
	return tracer.Start(c.Request().Context(), operationName)
}
