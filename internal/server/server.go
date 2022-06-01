package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/vliubezny/petsera/internal/storage"
)

type Config struct {
	Port        string
	AnnouncementStorage  storage.AnnouncementStorage
	FileStorage storage.FileStorage
	Statics     http.FileSystem
}

type Server struct {
	e          *echo.Echo
	address    string
	announcements       storage.AnnouncementStorage
	files      storage.FileStorage
}

func New(cfg Config) (*Server, error) {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	srv := &Server{
		e:          e,
		address:    fmt.Sprintf(":%s", cfg.Port),
		announcements:       cfg.AnnouncementStorage,
		files:      cfg.FileStorage,
	}

	assetHandler := http.FileServer(cfg.Statics)
	e.GET("/*", echo.WrapHandler(assetHandler))

	e.GET("/api/announcements", srv.getAnnouncementsHandler)
	e.POST("/api/announcements", srv.uploadHandler)
	e.GET("/api/images/:id", srv.getImageHandler)

	return srv, nil
}

func (srv *Server) Start(ctx context.Context) error {
	err := srv.e.Start(srv.address)
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.e.Shutdown(ctx)
}
