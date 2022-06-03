package model

import "time"

type Filter struct {
	North float64
	East  float64
	South float64
	West  float64
	After time.Time
}

type Location struct {
	Lat float64
	Lng float64
}

type Announcement struct {
	ID        string
	Text      string
	CreatedAt time.Time
	Position  Location
	ImageURL  string
}
