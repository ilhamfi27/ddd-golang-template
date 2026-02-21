package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Respond formats and sends a response based on success or error
func Respond(c echo.Context, data interface{}, err error) error {
	// Check if there is an error and respond accordingly
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	// Successful response with data
	return c.JSON(http.StatusOK, data)
}
