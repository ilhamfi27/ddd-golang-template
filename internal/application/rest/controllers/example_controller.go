package controllers

import (
	"strings"

	"github.com/ilhamfi27/ddd-golang-template/internal/application/dto"
	"github.com/ilhamfi27/ddd-golang-template/internal/domains"
	"github.com/ilhamfi27/ddd-golang-template/internal/utils"
	"github.com/labstack/echo/v4"
)

type ExampleController struct {
	ExampleService *domains.ExampleService
}

func NewExampleController(d *domains.ExampleService) *ExampleController {
	return &ExampleController{
		ExampleService: d,
	}
}

// @Summary Get data
// @Description Get data using GET method.
// @Tags Example
// @Accept application/json
// @Produce json
// @Param example query string true "Example"
// @Success 200 {object} map[string]interface{}
// @Router /example [get]
func (c *ExampleController) ExampleGetMtd(h echo.Context) error {
	name := h.QueryParam("name")
	var data dto.ParseExampleDto
	data.Name = strings.TrimSpace(name)
	result, err := c.ExampleService.Hello(h.Request().Context(), data)
	return utils.Respond(h, result, err)
}

// @Summary Get data
// @Description Get data using POST method.
// @Tags Example
// @Accept json
// @Produce json
// @Param body body dto.ParseExampleDto true "Parse Example Body"
// @Success 200 {object} map[string]interface{}
// @Router /example [post]
func (c *ExampleController) ExamplePostMtd(h echo.Context) error {
	var data dto.ParseExampleDto
	err := utils.BindBody(&data, h)
	if err != nil {
		return utils.Respond(h, nil, err)
	}
	result, err := c.ExampleService.Hello(h.Request().Context(), data)
	return utils.Respond(h, result, err)
}
