package server

import (
	"strings"
	"time"

	"github.com/vliubezny/petsera/internal/model"
	"golang.org/x/exp/constraints"
)

type ValidationError struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

func NewValidationError(msg string) *ValidationError {
	if msg == "" {
		msg = "Validation error(s)"
	}

	return &ValidationError{
		Message: msg,
		Fields:  make(map[string]string),
	}
}

func (e *ValidationError) Error() string {
	var b strings.Builder
	b.WriteString(e.Message)
	b.WriteString(":\n")
	for k, v := range e.Fields {
		b.WriteString(k)
		b.WriteString(" - ")
		b.WriteString(v)
		b.WriteString("\n")
	}

	return b.String()
}

func (e *ValidationError) HasIssues() bool {
	return len(e.Fields) > 0
}

func (e *ValidationError) AppendIssue(field, message string) {
	e.Fields[field] = message
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type AnnouncementRequest struct {
	Text     string   `json:"text"`
	Position Location `json:"position"`
}

func (r AnnouncementRequest) Validate() error {
	e := NewValidationError("")

	if !inRange(r.Position.Lat, -90, 90) {
		e.AppendIssue("position.lat", "Latitude must be a number between -90 and 90")
	}

	if !inRange(r.Position.Lng, -180, 180) {
		e.AppendIssue("position.lng", "Longitude must be a number between -180 and 180")
	}

	if !inRange(len(r.Text), 19, 2001) {
		e.AppendIssue("text", "Text must be between 20 and 2000 characters long")
	}

	if e.HasIssues() {
		return e
	}

	return nil
}

func inRange[T constraints.Ordered](val, min, max T) bool {
	return val > min && val < max
}

type Filter struct {
	North float64   `json:"north" query:"north"`
	East  float64   `json:"east" query:"east"`
	South float64   `json:"south" query:"south"`
	West  float64   `json:"west" query:"west"`
	After time.Time `json:"after" query:"after"`
}

func (f Filter) Validate() error {
	e := NewValidationError("")

	if !inRange(f.North, -90, 90) {
		e.AppendIssue("north", "North latitude must be a number between -90 and 90")
	}

	if !inRange(f.South, -90, 90) {
		e.AppendIssue("south", "South latitude must be a number between -90 and 90")
	}

	if f.North <= f.South {
		e.AppendIssue("north", "North latitude must be greater than south")
	}

	if !inRange(f.East, -180, 180) {
		e.AppendIssue("east", "East longitude must be a number between -180 and 180")
	}

	if !inRange(f.West, -180, 180) {
		e.AppendIssue("west", "West longitude must be a number between -180 and 180")
	}

	if f.East <= f.West {
		e.AppendIssue("east", "East longitude must be greater than west")
	}

	if f.After.IsZero() {
		e.AppendIssue("after", "After parameter is required")
	}

	if !f.After.Before(time.Now().UTC()) {
		e.AppendIssue("after", "After must be before current date")
	}

	if e.HasIssues() {
		return e
	}

	return nil
}

type AnnouncementResponse struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Position  Location  `json:"position"`
	ImageURL  string    `json:"imageURL"`
	CreatedAt time.Time `json:"createdAt"`
}

func ToAnnouncementResponses(announcements []*model.Announcement) []AnnouncementResponse {
	response := make([]AnnouncementResponse, len(announcements))
	for i, a := range announcements {
		response[i] = AnnouncementResponse{
			ID:        a.ID,
			Text:      a.Text,
			Position:  Location(a.Position),
			ImageURL:  a.ImageURL,
			CreatedAt: a.CreatedAt,
		}
	}

	return response
}
