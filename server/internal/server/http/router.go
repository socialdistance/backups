package http

import (
	"net/http"
	internalapp "server/internal/app"

	"github.com/labstack/echo/v4"
)

type Router struct {
	logger internalapp.Logger
	app    internalapp.App
}

func NewRouter(app internalapp.App, logger internalapp.Logger) *Router {
	return &Router{
		logger: logger,
		app:    app,
	}
}

func (r *Router) CommandHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello world")
}
