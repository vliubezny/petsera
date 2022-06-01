package model

import "time"

type Location struct {
	Lat float64
	Lng float64
}

type Announcement struct {
	ID string
	Text string
	CreatedAt time.Time
	Position Location
	ImageURL string
}