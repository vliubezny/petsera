package gcs

import (
	"context"
	"errors"
	"fmt"
	"io"

	gstorage "cloud.google.com/go/storage"

	"github.com/vliubezny/petsera/internal/storage"
)

type GCS struct {
	cli    *gstorage.Client
	bucket string
}

func NewGCS(ctx context.Context, bucket string) (*GCS, error) {
	cli, err := gstorage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &GCS{
		cli:    cli,
		bucket: bucket,
	}, nil
}

func (s *GCS) Put(ctx context.Context, key, contentType string, data io.Reader) (err error) {
	w := s.cli.Bucket(s.bucket).Object(key).NewWriter(ctx)
	defer func() {
		cerr := w.Close()
		if err == nil && cerr != nil {
			err = fmt.Errorf("failed to finish to write object to GCS: %w", cerr)
		}
	}()

	w.ContentType = contentType
	if _, err = io.Copy(w, data); err != nil {
		err = fmt.Errorf("failed to write object to GCS: %w", err)
	}

	return
}

func (s *GCS) Get(ctx context.Context, key string) (io.ReadCloser, string, error) {
	r, err := s.cli.Bucket(s.bucket).Object(key).NewReader(ctx)
	if err != nil {
		if errors.Is(err, gstorage.ErrObjectNotExist) {
			return nil, "", storage.ErrNotFound
		}

		return nil, "", fmt.Errorf("fail to read object from GCS: %w", err)
	}

	return r, r.Attrs.ContentType, nil
}

func (s *GCS) Delete(ctx context.Context, key string) error {
	if err := s.cli.Bucket(s.bucket).Object(key).Delete(ctx); err != nil {
		if errors.Is(err, gstorage.ErrObjectNotExist) {
			return storage.ErrNotFound
		}

		return fmt.Errorf("fail to delete object from GCS: %w", err)
	}

	return nil
}
