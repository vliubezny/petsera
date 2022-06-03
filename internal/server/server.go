package server

import (
	"context"
	"net"
	"net/http"

	echologrus "github.com/davrux/echo-logrus/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sirupsen/logrus"

	"github.com/vliubezny/petsera/internal/health"
	"github.com/vliubezny/petsera/internal/storage"
)

type Config struct {
	Port                string
	AnnouncementStorage storage.AnnouncementStorage
	FileStorage         storage.FileStorage
	Statics             http.FileSystem
	DevMode             bool
	Checker             health.Checker
}

type Server struct {
	e             *echo.Echo
	address       string
	announcements storage.AnnouncementStorage
	files         storage.FileStorage
	checker       health.Checker
}

func New(cfg Config) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.Debug = cfg.DevMode

	echologrus.Logger = logrus.New()
	e.Logger = echologrus.GetEchoLogger()
	e.Use(echologrus.Middleware())

	e.Use(middleware.Recover())

	var host string
	if cfg.DevMode {
		host = "127.0.0.1"
	}

	srv := &Server{
		e:             e,
		address:       net.JoinHostPort(host, cfg.Port),
		announcements: cfg.AnnouncementStorage,
		files:         cfg.FileStorage,
		checker:       cfg.Checker,
	}

	e.GET("/health", srv.getHealthHandler)

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
