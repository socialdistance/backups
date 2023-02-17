package http

import (
	"encoding/json"
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
	var task TaskDTO

	dto, err := task.GetModelTask()
	if err != nil {
		r.logger.Error(err.Error())
	}

	response, err := json.Marshal(dto)
	if err != nil {
		r.logger.Error(err.Error())
	}

	return c.JSON(http.StatusOK, response)
}
