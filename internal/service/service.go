package service

import (
	"context"
	"errors"

	"github.com/paraumir/shortener/internal/model"
	"github.com/paraumir/shortener/internal/storage"
)

type Service struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateShortURL(
	ctx context.Context,
	originalURL string,
) (model.Link, error) {
	existing, err := s.storage.GetByOriginalURL(ctx, originalURL)

	if err == nil {
		return existing, nil
	}

	if !errors.Is(err, storage.ErrNotFound) {
		return model.Link{}, err
	}

	for {
		shortURL, err := GenerateShortURL()
		if err != nil {
			return model.Link{}, err
		}

		link := model.Link{
			OriginalURL: originalURL,
			ShortURL:    shortURL,
		}

		err = s.storage.Save(ctx, link)

		if err == nil {
			return link, nil
		}

		if errors.Is(err, storage.ErrShortURLExists) {
			continue
		}

		if errors.Is(err, storage.ErrOriginalURLExists) {
			existing, err := s.storage.GetByOriginalURL(ctx, originalURL)
			if err != nil {
				return model.Link{}, err
			}

			return existing, nil
		}

		return model.Link{}, err
	}
}

func (s *Service) GetOriginalURL(
	ctx context.Context,
	shortURL string,
) (string, error) {
	link, err := s.storage.GetByShortURL(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return link.OriginalURL, nil
}
