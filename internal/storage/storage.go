package storage

import (
	"context"
	"errors"
	"io"

	"github.com/vliubezny/petsera/internal/model"
)

var (
	ErrNotFound = errors.New("not found")
)

type FileStorage interface {
	Put(ctx context.Context, key, contentType string, data io.Reader) error
	Get(ctx context.Context, key string) (io.ReadCloser, string, error)
	Delete(ctx context.Context, key string) error
}

type AnnouncementStorage interface {
	Create(ctx context.Context, announcement *model.Announcement) error
	GetAll(ctx context.Context, filter model.Filter) ([]*model.Announcement, error)
	Delete(ctx context.Context, id string) error
	InTx(ctx context.Context, txFn func(s AnnouncementStorage) error) error
}
