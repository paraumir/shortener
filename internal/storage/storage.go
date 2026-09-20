package storage

import (
	"context"
	"errors"

	"github.com/paraumir/shortener/internal/model"
)

var (
	ErrNotFound          = errors.New("link not found")
	ErrOriginalURLExists = errors.New("original URL already exists")
	ErrShortURLExists    = errors.New("short URL already exists")
)

type Storage interface {
	Save(ctx context.Context, link model.Link) error
	GetByShortURL(ctx context.Context, shortURL string) (model.Link, error)
	GetByOriginalURL(ctx context.Context, originalURL string) (model.Link, error)
}
