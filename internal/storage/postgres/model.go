package postgres

import (
	"time"

	"github.com/vliubezny/petsera/internal/model"
)

type announcement struct {
	ID        string    `db:"id"`
	Text      string    `db:"text"`
	ImageURL  string    `db:"image_url"`
	Lat       float64   `db:"lat"`
	Lng       float64   `db:"lng"`
	CreatedAt time.Time `db:"created_at"`
}

func (a announcement) toModel() *model.Announcement {
	return &model.Announcement{
		ID:       a.ID,
		Text:     a.Text,
		ImageURL: a.ImageURL,
		Position: model.Location{
			Lat: a.Lat,
			Lng: a.Lng,
		},
		CreatedAt: a.CreatedAt,
	}
}
