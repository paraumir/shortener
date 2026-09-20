package storage

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/paraumir/shortener/internal/model"
)

func newTestPostgres(t *testing.T) *pgxpool.Pool {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		databaseURL = "postgres://shortener:shortener@localhost:5433/shortener"
	}

	ctx := context.Background()

	db, err := NewPostgresPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func cleanLinks(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	_, err := db.Exec(
		context.Background(),
		"TRUNCATE TABLE links RESTART IDENTITY",
	)
	if err != nil {
		t.Fatalf("failed to clean links table: %v", err)
	}
}

func TestPostgresStorage_SaveAndGet(t *testing.T) {
	db := newTestPostgres(t)
	cleanLinks(t, db)

	store := NewPostgresStorage(db)
	ctx := context.Background()

	link := model.Link{
		OriginalURL: "https://example.com",
		ShortURL:    "aB123_XyZ9",
	}

	err := store.Save(ctx, link)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	got, err := store.GetByShortURL(ctx, link.ShortURL)
	if err != nil {
		t.Fatalf("GetByShortURL() returned error: %v", err)
	}

	if got != link {
		t.Fatalf(
			"GetByShortURL() = %+v, want %+v",
			got,
			link,
		)
	}
}

func TestPostgresStorage_GetByOriginalURL(t *testing.T) {
	db := newTestPostgres(t)
	cleanLinks(t, db)

	store := NewPostgresStorage(db)
	ctx := context.Background()

	link := model.Link{
		OriginalURL: "https://example.com",
		ShortURL:    "aB123_XyZ9",
	}

	if err := store.Save(ctx, link); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	got, err := store.GetByOriginalURL(ctx, link.OriginalURL)
	if err != nil {
		t.Fatalf("GetByOriginalURL() returned error: %v", err)
	}

	if got != link {
		t.Fatalf(
			"GetByOriginalURL() = %+v, want %+v",
			got,
			link,
		)
	}
}

func TestPostgresStorage_NotFound(t *testing.T) {
	db := newTestPostgres(t)
	cleanLinks(t, db)

	store := NewPostgresStorage(db)
	ctx := context.Background()

	_, err := store.GetByShortURL(ctx, "aB123_XyZ9")

	if err != ErrNotFound {
		t.Fatalf(
			"GetByShortURL() error = %v, want %v",
			err,
			ErrNotFound,
		)
	}
}

func TestPostgresStorage_DuplicateOriginalURL(t *testing.T) {
	db := newTestPostgres(t)
	cleanLinks(t, db)

	store := NewPostgresStorage(db)
	ctx := context.Background()

	first := model.Link{
		OriginalURL: "https://example.com",
		ShortURL:    "aB123_XyZ9",
	}

	second := model.Link{
		OriginalURL: "https://example.com",
		ShortURL:    "Qw456_AbC7",
	}

	if err := store.Save(ctx, first); err != nil {
		t.Fatalf(
			"first Save() returned error: %v",
			err,
		)
	}

	err := store.Save(ctx, second)

	if err != ErrOriginalURLExists {
		t.Fatalf(
			"second Save() error = %v, want %v",
			err,
			ErrOriginalURLExists,
		)
	}
}

func TestPostgresStorage_DuplicateShortURL(t *testing.T) {
	db := newTestPostgres(t)
	cleanLinks(t, db)

	store := NewPostgresStorage(db)
	ctx := context.Background()

	first := model.Link{
		OriginalURL: "https://example.com",
		ShortURL:    "aB123_XyZ9",
	}

	second := model.Link{
		OriginalURL: "https://google.com",
		ShortURL:    "aB123_XyZ9",
	}

	if err := store.Save(ctx, first); err != nil {
		t.Fatalf(
			"first Save() returned error: %v",
			err,
		)
	}

	err := store.Save(ctx, second)

	if err != ErrShortURLExists {
		t.Fatalf(
			"second Save() error = %v, want %v",
			err,
			ErrShortURLExists,
		)
	}
}
