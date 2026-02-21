package handlers

import (
	"log"
	"net/http"

	config "github.com/ilhamfi27/ddd-golang-template/internal/config/app"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	app *echo.Echo
	db  *gorm.DB
}

func NewHandler() *Handler {
	return &Handler{
		app: echo.New(),
	}
}

func (h *Handler) Start() error {
	port := config.GetEnvVars().Port

	dbHdl := NewDBHandler()
	h.db = dbHdl.DBInit()

	NewErrorHandler(h)
	NewMiddlewareHandler(h)
	NewRestHandler(h)
	err := h.app.Start(port)
	if err != nil && err != http.ErrServerClosed {
		// If an error occurs (other than a normal server closed error), log it and exit with a non-zero (failure) code
		log.Fatal(err)
	}
	return err
}

func (h *Handler) Migrate() {
	dbHdl := NewDBHandler()
	dbHdl.Migrate()
}
