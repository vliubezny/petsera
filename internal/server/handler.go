package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"

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
		return echo.NewHTTPError(http.StatusBadRequest, "fail to parse'data'").SetInternal(err)
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

	imageContentType:= file.Header.Get("Content-Type")
	

	doc := &model.Announcement{
		ID:          ksuid.New().String(),
		Text: req.Text,
		Position: model.Location(req.Position),
		CreatedAt:   time.Now().UTC(),
	}

	doc.ImageURL = fmt.Sprintf("/api/images/%s", doc.ID)

	log.Printf("uploading image: %+v", doc)

	if err := srv.announcements.InTx(func(tx storage.AnnouncementStorage) error {
		if err := tx.Create(c.Request().Context(), doc); err != nil {
			return err
		}

		if err := srv.files.Put(c.Request().Context(), doc.ID, imageContentType, src); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, doc)
}

func (srv *Server) getAnnouncementsHandler(c echo.Context) error {
	docs, err := srv.announcements.GetAll(c.Request().Context())
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
