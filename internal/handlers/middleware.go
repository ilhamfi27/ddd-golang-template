package handlers

import (
	"github.com/labstack/echo/v4/middleware"
)

func NewMiddlewareHandler(h *Handler) {
	h.app.Use(middleware.Logger())
	h.app.Use(middleware.CORS())
}
