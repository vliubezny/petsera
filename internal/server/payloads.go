package server

import (
	"strings"
	"time"

	"github.com/vliubezny/petsera/internal/model"
	"golang.org/x/exp/constraints"
)

type ValidationError struct {
	Message string `json:"message"`
	Fields map[string]string `json:"fields"`
}

func NewValidationError(msg string) *ValidationError {
	if msg == "" {
		msg = "Validation error(s)"
	}

	return &ValidationError{
		Message: msg,
		Fields: make(map[string]string),
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
	Text string `json:"text"`
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

type AnnouncementResponse struct {
	ID string `json:"id"`
	Text string `json:"text"`
	Position Location `json:"position"`
	ImageURL string `json:"imageURL"`
	CreatedAt time.Time `json:"createdAt"`
}

func ToAnnouncementResponses(announcements []*model.Announcement) []AnnouncementResponse {
	response := make([]AnnouncementResponse, len(announcements))
	for i, a := range announcements {
		response[i] = AnnouncementResponse {
			ID: a.ID,
			Text: a.Text,
			Position: Location(a.Position),
			ImageURL: a.ImageURL,
			CreatedAt: a.CreatedAt,
		}
	}

	return response
}