package controllers

import "github.com/labstack/echo/v4"

type BaseController struct {
}

func NewBaseController() *BaseController {
	return &BaseController{}
}

// @Summary Show the root path.
// @Description get the root path.
// @Tags Root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func (b *BaseController) GetBase(h echo.Context) error {
	return h.JSON(200, map[string]string{"message": "Hello, World!"})
}

// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /healthcheck [get]
func (b *BaseController) Healthcheck(h echo.Context) error {
	return h.JSON(200, map[string]string{"message": "I'm alive!"})
}
