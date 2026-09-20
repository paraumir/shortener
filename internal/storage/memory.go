package storage

import (
	"context"
	"sync"

	"github.com/paraumir/shortener/internal/model"
)

type MemoryStorage struct {
	mu sync.RWMutex

	byShortURL    map[string]model.Link
	byOriginalURL map[string]model.Link
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		byShortURL:    make(map[string]model.Link),
		byOriginalURL: make(map[string]model.Link),
	}
}

func (s *MemoryStorage) Save(ctx context.Context, link model.Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.byOriginalURL[link.OriginalURL]; exists {
		return ErrOriginalURLExists
	}

	if _, exists := s.byShortURL[link.ShortURL]; exists {
		return ErrShortURLExists
	}

	s.byShortURL[link.ShortURL] = link
	s.byOriginalURL[link.OriginalURL] = link

	return nil
}

func (s *MemoryStorage) GetByShortURL(
	ctx context.Context,
	shortURL string,
) (model.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, ok := s.byShortURL[shortURL]
	if !ok {
		return model.Link{}, ErrNotFound
	}

	return link, nil
}

func (s *MemoryStorage) GetByOriginalURL(
	ctx context.Context,
	originalURL string,
) (model.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, ok := s.byOriginalURL[originalURL]
	if !ok {
		return model.Link{}, ErrNotFound
	}

	return link, nil
}
