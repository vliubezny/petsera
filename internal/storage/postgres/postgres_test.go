//go:build integration
// +build integration

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/vliubezny/petsera/internal/model"
	"github.com/vliubezny/petsera/internal/storage"
)

var (
	db  *sql.DB
	ctx = context.Background()
	s   storage.AnnouncementStorage
)

func TestMain(m *testing.M) {
	shutdown := setup()
	s = New(db)

	code := m.Run()

	shutdown()
	os.Exit(code)
}

func setup() func() {
	req := testcontainers.ContainerRequest{
		Image:        "postgis/postgis:13-master",
		Env:          map[string]string{"POSTGRES_PASSWORD": "root"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})
	if err != nil {
		logrus.WithError(err).Fatalf("failed to create container")
	}

	if err := c.Start(ctx); err != nil {
		logrus.WithError(err).Fatal("failed to start container")
	}

	host, err := c.Host(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("failed to get host")
	}

	port, err := c.MappedPort(ctx, "5432")
	if err != nil {
		logrus.WithError(err).Fatal("failed to map port")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=postgres password=root dbname=postgres sslmode=disable", host, port.Int())
	logrus.Info(dsn)
	db = MustSetupDB(dsn, 0, 1, true, "../../../scripts/migrations/postgres/")

	shutdownFn := func() {
		if c != nil {
			c.Terminate(ctx)
		}
	}

	return shutdownFn
}

func TestPg_Create(t *testing.T) {
	a := newAnnouncement("LAX", "image", "2018-09-16T12:00:00", 33.9434, -118.4079)
	err := s.Create(ctx, a)
	assert.NoError(t, err)

	aa, err := s.GetAll(ctx, newFilter(-120, 30, -60, 45, "2018-09-12T12:00:00"))
	require.NoError(t, err)
	assert.Equal(t, []*model.Announcement{a}, aa)
}

func TestPg_Delete(t *testing.T) {
	a := newAnnouncement("LAX2", "image", "2018-09-16T12:00:00", 33.9434, -118.4079)
	err := s.Create(ctx, a)
	assert.NoError(t, err)

	err = s.Delete(ctx, a.ID)
	require.NoError(t, err)

	aa, err := s.GetAll(ctx, newFilter(-120, 30, -60, 45, "2018-09-12T12:00:00"))
	require.NoError(t, err)
	assert.NotContains(t, aa, a)
}

func newAnnouncement(txt, url, ts string, lat, lng float64) *model.Announcement {
	return &model.Announcement{
		ID:       ksuid.New().String(),
		Text:     txt,
		ImageURL: url,
		Position: model.Location{
			Lat: lat,
			Lng: lng,
		},
		CreatedAt: timestamp(ts),
	}
}

func newFilter(west, south, east, north float64, ts string) model.Filter {
	return model.Filter{
		West:  west,
		South: south,
		East:  east,
		North: north,
		After: timestamp(ts),
	}
}

func timestamp(val string) time.Time {
	t, err := time.Parse(time.RFC3339, fmt.Sprintf("%s-00:00", val))
	if err != nil {
		panic(err)
	}

	return t
}
