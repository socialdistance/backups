package http

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"server/internal/app"
	"time"
)

type Server struct {
	host   string
	port   string
	e      *echo.Echo
	app    *app.App
	router *Router
}

func NewServer(host, port string, app *app.App, router *Router) *Server {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORS())

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

// BuildRouters TODO: serve static files
func (f *Server) BuildRouters() {
	//f.e.Static("/", "tmp")
	//fs := http.FileServer(http.Dir("/home/user/work/class_attachments/server/tmp"))
	//f.e.GET("/uploads/*", echo.WrapHandler(http.StripPrefix("/uploads/", fs)))
	//
	//fsAPI := f.e.Group("/api")
	//
	//fsAPI.POST("/upload", f.router.Upload)
	//fsAPI.DELETE("/delete", f.router.Delete)
	//fsAPI.GET("/data", f.router.ListData) // get data from table attachments
	//fsAPI.POST("/name", f.router.GetUser) // get user from table info
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
