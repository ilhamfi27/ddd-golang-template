package utils

import (
	"encoding/json"

	rest_error "github.com/ilhamfi27/ddd-golang-template/internal/application/rest/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func BindBody(body interface{}, e echo.Context) error {
	err := e.Bind(body)
	if err != nil {
		return rest_error.BadRequestError("Invalid request body")
	}
	return nil
}

func GetRequestBody(c echo.Context) map[string]interface{} {
	jsonMap := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
	if err != nil {
		log.Error("empty json body")
		return nil
	}
	return jsonMap
}
