package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/paraumir/shortener/internal/model"
)

const (
	originalURLConstraint = "links_original_url_key"
	shortURLConstraint    = "links_short_url_key"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

func NewPostgresStorage(db *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) Save(
	ctx context.Context,
	link model.Link,
) error {
	_, err := s.db.Exec(
		ctx,
		`
		INSERT INTO links (original_url, short_url)
		VALUES ($1, $2)
		`,
		link.OriginalURL,
		link.ShortURL,
	)

	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.ConstraintName == originalURLConstraint {
			return ErrOriginalURLExists
		}

		if pgErr.ConstraintName == shortURLConstraint {
			return ErrShortURLExists
		}
	}

	return err
}

func (s *PostgresStorage) GetByShortURL(
	ctx context.Context,
	shortURL string,
) (model.Link, error) {
	var link model.Link

	err := s.db.QueryRow(
		ctx,
		`
		SELECT original_url, short_url
		FROM links
		WHERE short_url = $1
		`,
		shortURL,
	).Scan(
		&link.OriginalURL,
		&link.ShortURL,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.Link{}, ErrNotFound
	}

	if err != nil {
		return model.Link{}, err
	}

	return link, nil
}

func (s *PostgresStorage) GetByOriginalURL(
	ctx context.Context,
	originalURL string,
) (model.Link, error) {
	var link model.Link

	err := s.db.QueryRow(
		ctx,
		`
		SELECT original_url, short_url
		FROM links
		WHERE original_url = $1
		`,
		originalURL,
	).Scan(
		&link.OriginalURL,
		&link.ShortURL,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.Link{}, ErrNotFound
	}

	if err != nil {
		return model.Link{}, err
	}

	return link, nil
}

func NewPostgresPool(
	ctx context.Context,
	databaseURL string,
) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
