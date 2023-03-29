package http

import (
	"context"
	"fmt"
	"net/http"
	"server/internal/app"
	"time"

	internallogger "server/internal/logger"

	"github.com/labstack/echo/v4"
)

type Server struct {
	host   string
	port   string
	e      *echo.Echo
	app    *app.App
	router *Router
	// logg   internallogger.Logger
}

func NewServer(host, port string, app *app.App, router *Router, logg internallogger.Logger) *Server {
	e := echo.New()
	e.HideBanner = true

	e.Use(MiddlwareLogger(logg))

	return &Server{
		host:   host,
		port:   port,
		e:      e,
		app:    app,
		router: router,
	}
}

func (f *Server) BuildRouters() {
	f.e.Static("/", "uploads")
	//fs := http.FileServer(http.Dir("/Users/user/work/dev/backups/uploads"))
	fs := http.FileServer(http.Dir("/home/user/work/backup/uploads"))
	f.e.GET("/uploads/*", echo.WrapHandler(http.StripPrefix("/uploads/", fs)))

	api := f.e.Group("/api")
	api.GET("/command", f.router.CommandHandler)
	api.POST("/upload", f.router.UploadFile)
}

func (f *Server) Start() error {
	if err := f.e.Start(fmt.Sprintf(":%s", f.port)); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("client stopped: %w", err)
	}

	return nil
}

func (f *Server) Stop() error {
	optCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := f.e.Shutdown(optCtx); err != nil {
		return fmt.Errorf("could not shutdown client gracefuly: %w", err)
	}

	return nil
}
