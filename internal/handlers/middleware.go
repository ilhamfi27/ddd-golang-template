package handlers

import (
	config "github.com/ilhamfi27/ddd-golang-template/internal/config/app"
	"github.com/ilhamfi27/ddd-golang-template/internal/handlers/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func NewMiddlewareHandler(h *Handler) {
	// Initialize validator
	middleware.InitValidator()

	// Recovery middleware - recovers from panics
	h.app.Use(echomiddleware.Recover())

	// Request ID middleware - adds request ID for correlation
	h.app.Use(echomiddleware.RequestID())

	// Structured logging middleware
	h.app.Use(middleware.StructuredLoggingMiddleware())

	// CORS middleware
	h.app.Use(echomiddleware.CORS())

	// Tracing middleware - for distributed tracing
	h.app.Use(middleware.TracingMiddleware())

	// JWT middleware for protected routes
	env := config.GetEnvVars()
	jwtConfig := middleware.JWTConfig(env.JWTSecret)

	// Apply JWT only to /api routes
	h.app.Group("/api", echojwt.WithConfig(jwtConfig))
}
