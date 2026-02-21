package handlers

import (
	"github.com/ilhamfi27/ddd-golang-template/internal/application/rest/controllers"
	"github.com/ilhamfi27/ddd-golang-template/internal/domains"
	repo "github.com/ilhamfi27/ddd-golang-template/internal/infrastructure/repositories"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	ExamplePath = "/example"
)

func NewRestHandler(h *Handler) {

	baseController := controllers.NewBaseController()

	exampleRepo := repo.NewGormExampleRepository(h.db)
	exampleService := domains.NewExampleService(exampleRepo)
	exampleController := controllers.NewExampleController(exampleService)

	h.app.GET("/", baseController.GetBase)
	h.app.GET("/healthcheck", baseController.Healthcheck)

	// Example
	h.app.GET(ExamplePath, exampleController.ExampleGetMtd)
	h.app.POST(ExamplePath, exampleController.ExamplePostMtd)

	h.app.GET("/swagger/*", echoSwagger.WrapHandler)
}
