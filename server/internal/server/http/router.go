package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	internalapp "server/internal/app"

	"github.com/google/uuid"
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
	task := new(WorkerTaskDTO)

	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, "Error bind")
	}

	taskIdStr, err := uuid.Parse(task.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Can't parse uuid")
	}

	taskResponse, err := r.app.CommandHandlerApp(c.Request().Context(), taskIdStr, task.Address, task.Command, task.Hostname)
	if err != nil {
		// TODO:
		return c.JSON(http.StatusBadRequest, "something wrong")
	}

	return c.JSON(http.StatusOK, taskResponse)
}

func (r *Router) UploadFile(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(fmt.Sprintf("/Users/user/work/dev/backups/uploads/%s", file.Filename))
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

	}

	return c.JSON(http.StatusOK, "test")
}
