package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/vliubezny/petsera/internal/model"
	"github.com/vliubezny/petsera/internal/storage"
)

type extContext interface {
	sqlx.ExtContext
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	GetContext(ctx context.Context, dest any, query string, args ...any) error
}

type PG struct {
	dbx *sqlx.DB
	ext extContext
}

// New creates postgres storage.
func New(db *sql.DB) PG {
	dbx := sqlx.NewDb(db, "postgres")
	return PG{
		dbx: dbx,
		ext: dbx,
	}
}

func (p PG) Close() error {
	if p.dbx != nil {
		return p.dbx.Close()
	}

	return nil
}

func (p PG) Create(ctx context.Context, a *model.Announcement) error {
	query := `
		INSERT INTO announcement (id, text, image_url, position, created_at)
		VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326), $6);
	`
	_, err := p.ext.ExecContext(ctx, query, a.ID, a.Text, a.ImageURL, a.Position.Lng, a.Position.Lat, a.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create announcement: %w", err)
	}

	return nil
}

func (p PG) GetAll(ctx context.Context, filter model.Filter) ([]*model.Announcement, error) {
	query := `
		SELECT id, text, image_url, ST_X(position::geometry) AS lng, ST_Y(position::geometry) AS lat, created_at
		FROM announcement
		WHERE ST_Intersects
			( position
			, ST_MakeEnvelope ( $1 -- xmin (min lng)
							, $2 -- ymin (min lat)
							, $3 -- xmax (max lng)
							, $4 -- ymax (max lat)
							, 4326 -- projection epsg-code
							)::geography('POLYGON')
			)
			AND created_at > $5;
	`
	var announcements []announcement
	if err := p.ext.SelectContext(
		ctx,
		&announcements,
		query,
		filter.West,
		filter.South,
		filter.East,
		filter.North,
		filter.After,
	); err != nil {
		return nil, fmt.Errorf("failed to get announcements: %w", err)
	}

	data := make([]*model.Announcement, len(announcements))
	for i, a := range announcements {
		data[i] = a.toModel()
	}

	return data, nil
}

func (p PG) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM announcement WHERE id = $1;
	`
	if _, err := p.ext.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}

	return nil
}

func (p PG) InTx(ctx context.Context, action func(s storage.AnnouncementStorage) error) error {
	tx, err := p.dbx.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = action(PG{dbx: p.dbx, ext: tx})

	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v root: %w", rbErr, err)
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (p PG) Check(ctx context.Context) error {
	if err := p.dbx.PingContext(ctx); err != nil {
		return fmt.Errorf("postgres check failed: %w", err)
	}

	return nil
}
