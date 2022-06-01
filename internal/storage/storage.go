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
	Create(ctx context.Context, doc *model.Announcement) error
	GetAll(ctx context.Context) ([]*model.Announcement, error)
	Delete(ctx context.Context, id string) error
	InTx(txFn func(tx AnnouncementStorage) error) error
}
