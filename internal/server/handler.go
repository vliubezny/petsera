package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/vliubezny/petsera/internal/model"
	"github.com/vliubezny/petsera/internal/storage"
)

func (srv *Server) uploadHandler(c echo.Context) error {
	data := c.FormValue("data")
	if data == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "'data' parameter is missing")
	}

	var req AnnouncementRequest
	if err := json.Unmarshal([]byte(data), &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "fail to parse 'data'").SetInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err).SetInternal(err)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	imageContentType := file.Header.Get("Content-Type")

	announcement := &model.Announcement{
		ID:        ksuid.New().String(),
		Text:      req.Text,
		Position:  model.Location(req.Position),
		CreatedAt: time.Now().UTC(),
	}

	announcement.ImageURL = fmt.Sprintf("/api/images/%s", announcement.ID)

	logrus.Infof("save announcement: %+v", announcement)

	if err := srv.announcements.InTx(c.Request().Context(), func(tx storage.AnnouncementStorage) error {
		if err := tx.Create(c.Request().Context(), announcement); err != nil {
			return err
		}

		if err := srv.files.Put(c.Request().Context(), announcement.ID, imageContentType, src); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, announcement)
}

func (srv *Server) getAnnouncementsHandler(c echo.Context) error {
	var filter Filter
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "fail to parse filter").SetInternal(err)
	}

	if err := filter.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err).SetInternal(err)
	}

	docs, err := srv.announcements.GetAll(c.Request().Context(), model.Filter(filter))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, docs)
}

func (srv *Server) getImageHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.ErrBadRequest
	}

	data, contentType, err := srv.files.Get(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return echo.ErrNotFound
		}

		return err
	}

	defer data.Close()

	return c.Stream(http.StatusOK, contentType, data)
}

func (srv *Server) getHealthHandler(c echo.Context) error {
	if err := srv.checker.Check(c.Request().Context()); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "UNHEALTHY",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": "HEALTHY",
	})
}

func (srv *Server) indexHandler(frontendCfg map[string]any, statics http.FileSystem) (echo.HandlerFunc, error) {
	cfg, err := json.Marshal(frontendCfg)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal frontend config: %w", err)
	}

	f, err := statics.Open("index.html")
	if err != nil {
		return nil, fmt.Errorf("fail to open index.html: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("fail to read index.html: %w", err)
	}

	index := strings.Replace(string(data), "__APP_CONFIG__", string(cfg), 1)

	return func(c echo.Context) error {
		return c.HTML(http.StatusOK, index)
	}, nil
}
