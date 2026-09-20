package service

import (
	"context"
	"testing"

	"github.com/paraumir/shortener/internal/storage"
)

func TestService_CreateShortURL(t *testing.T) {
	store := storage.NewMemoryStorage()
	service := New(store)
	ctx := context.Background()

	link, err := service.CreateShortURL(
		ctx,
		"https://example.com",
	)
	if err != nil {
		t.Fatalf("CreateShortURL() returned error: %v", err)
	}

	if link.OriginalURL != "https://example.com" {
		t.Fatalf(
			"OriginalURL = %q, want %q",
			link.OriginalURL,
			"https://example.com",
		)
	}

	if len(link.ShortURL) != 10 {
		t.Fatalf(
			"ShortURL length = %d, want 10",
			len(link.ShortURL),
		)
	}
}

func TestService_SameURLReturnsSameShortURL(t *testing.T) {
	store := storage.NewMemoryStorage()
	service := New(store)
	ctx := context.Background()

	first, err := service.CreateShortURL(
		ctx,
		"https://example.com",
	)
	if err != nil {
		t.Fatalf(
			"first CreateShortURL() returned error: %v",
			err,
		)
	}

	second, err := service.CreateShortURL(
		ctx,
		"https://example.com",
	)
	if err != nil {
		t.Fatalf(
			"second CreateShortURL() returned error: %v",
			err,
		)
	}

	if first.ShortURL != second.ShortURL {
		t.Fatalf(
			"short URLs are different: first=%q second=%q",
			first.ShortURL,
			second.ShortURL,
		)
	}
}

func TestService_GetOriginalURL(t *testing.T) {
	store := storage.NewMemoryStorage()
	service := New(store)
	ctx := context.Background()

	link, err := service.CreateShortURL(
		ctx,
		"https://example.com",
	)
	if err != nil {
		t.Fatalf(
			"CreateShortURL() returned error: %v",
			err,
		)
	}

	originalURL, err := service.GetOriginalURL(
		ctx,
		link.ShortURL,
	)
	if err != nil {
		t.Fatalf(
			"GetOriginalURL() returned error: %v",
			err,
		)
	}

	if originalURL != "https://example.com" {
		t.Fatalf(
			"original URL = %q, want %q",
			originalURL,
			"https://example.com",
		)
	}
}

func TestService_ConcurrentCreate(t *testing.T) {
	store := storage.NewMemoryStorage()
	service := New(store)
	ctx := context.Background()

	const goroutines = 100

	results := make(chan string, goroutines)
	errors := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			link, err := service.CreateShortURL(
				ctx,
				"https://example.com",
			)
			if err != nil {
				errors <- err
				return
			}

			results <- link.ShortURL
		}()
	}

	var shortURL string

	for i := 0; i < goroutines; i++ {
		select {
		case err := <-errors:
			t.Fatalf("CreateShortURL() returned error: %v", err)

		case result := <-results:
			if shortURL == "" {
				shortURL = result
				continue
			}

			if result != shortURL {
				t.Fatalf(
					"got different short URLs: %q and %q",
					shortURL,
					result,
				)
			}
		}
	}
}
