package handlers

import (
	"net/http"

	rest_error "github.com/ilhamfi27/ddd-golang-template/internal/application/rest/errors"
	"github.com/labstack/echo/v4"
)

// Custom error handler
func customErrorHandler(err error, c echo.Context) {
	// Check for custom error type
	if appErr, ok := err.(*rest_error.HttpError); ok {
		// Respond with the appropriate HTTP status code and message
		c.JSON(appErr.Code, map[string]interface{}{
			"error": appErr.Message,
		})
		return
	}

	if echoErr, ok := err.(*echo.HTTPError); ok {
		// Respond with the appropriate HTTP status code and message
		c.JSON(echoErr.Code, map[string]interface{}{
			"error": echoErr.Message,
		})
		return
	}

	// Default response if error type is not matched
	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error": "Internal server error",
	})
}

func NewErrorHandler(h *Handler) {
	h.app.HTTPErrorHandler = customErrorHandler
}
