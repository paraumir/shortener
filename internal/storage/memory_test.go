package storage

import (
	"context"
	"testing"

	"github.com/paraumir/shortener/internal/model"
)

func TestMemoryStorage_SaveAndGetByShortURL(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	link := model.Link{
		OriginalURL: "https://example.com",
		ShortURL:    "abc123_XYZ",
	}

	err := storage.Save(ctx, link)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	got, err := storage.GetByShortURL(ctx, link.ShortURL)
	if err != nil {
		t.Fatalf("GetByShortURL() returned error: %v", err)
	}

	if got != link {
		t.Fatalf("GetByShortURL() = %+v, want %+v", got, link)
	}
}

func TestMemoryStorage_GetByOriginalURL(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	link := model.Link{
		OriginalURL: "https://example.com",
		ShortURL:    "abc123_XYZ",
	}

	err := storage.Save(ctx, link)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	got, err := storage.GetByOriginalURL(ctx, link.OriginalURL)
	if err != nil {
		t.Fatalf("GetByOriginalURL() returned error: %v", err)
	}

	if got != link {
		t.Fatalf("GetByOriginalURL() = %+v, want %+v", got, link)
	}
}

func TestMemoryStorage_NotFound(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	_, err := storage.GetByShortURL(ctx, "does_not_exist")

	if err != ErrNotFound {
		t.Fatalf("GetByShortURL() error = %v, want %v", err, ErrNotFound)
	}
}

func TestMemoryStorage_DuplicateOriginalURL(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	first := model.Link{
		OriginalURL: "https://example.com",
		ShortURL:    "abc123_XYZ",
	}

	second := model.Link{
		OriginalURL: "https://example.com",
		ShortURL:    "different_1",
	}

	if err := storage.Save(ctx, first); err != nil {
		t.Fatalf("first Save() returned error: %v", err)
	}

	err := storage.Save(ctx, second)
	if err != ErrOriginalURLExists {
		t.Fatalf(
			"second Save() error = %v, want %v",
			err,
			ErrOriginalURLExists,
		)
	}
}

func TestMemoryStorage_DuplicateShortURL(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	first := model.Link{
		OriginalURL: "https://example.com",
		ShortURL:    "abc123_XYZ",
	}

	second := model.Link{
		OriginalURL: "https://google.com",
		ShortURL:    "abc123_XYZ",
	}

	if err := storage.Save(ctx, first); err != nil {
		t.Fatalf("first Save() returned error: %v", err)
	}

	err := storage.Save(ctx, second)
	if err != ErrShortURLExists {
		t.Fatalf(
			"second Save() error = %v, want %v",
			err,
			ErrShortURLExists,
		)
	}
}
