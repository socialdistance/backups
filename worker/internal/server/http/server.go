package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"worker/internal/app"

	internallogger "worker/internal/logger"

	"github.com/labstack/echo/v4"
)

type Server struct {
	host   string
	port   string
	e      *echo.Echo
	app    *app.App
	router *Router
}

func NewServer(host, port string, app *app.App, router *Router, logg internallogger.Logger) *Server {
	e := echo.New()
	e.HideBanner = true

	// e.Use(middleware.CORS())
	e.Use(MiddlwareLogger(logg))

	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"http://localhost:3000"},
	//	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	//	AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	//}))

	//e.Use(middleware.Static("/home/user/work/mountaineering/uploads"))

	return &Server{
		host:   host,
		port:   port,
		e:      e,
		app:    app,
		router: router,
	}
}

func (f *Server) BuildRouters() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				f.router.RequestToControlServer()
			}
		}
	}()
}

func (f *Server) Start() error {
	if err := f.e.Start(fmt.Sprintf(":%s", f.port)); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server stopped: %w", err)
	}

	return nil
}

func (f *Server) Stop() error {
	optCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := f.e.Shutdown(optCtx); err != nil {
		return fmt.Errorf("could not shutdown server gracefuly: %w", err)
	}

	return nil
}
