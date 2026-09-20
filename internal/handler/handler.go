package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/paraumir/shortener/internal/service"
	"github.com/paraumir/shortener/internal/storage"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

type originalURLResponse struct {
	URL string `json:"url"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request shortenRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	request.URL = strings.TrimSpace(request.URL)

	if request.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	link, err := h.service.CreateShortURL(
		r.Context(),
		request.URL,
	)
	if err != nil {
		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	response := shortenResponse{
		ShortURL: link.ShortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (h *Handler) GetOriginalURL(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortURL := strings.TrimPrefix(r.URL.Path, "/")

	if !service.IsValidShortURL(shortURL) {
		http.Error(w, "invalid short URL", http.StatusBadRequest)
		return
	}

	originalURL, err := h.service.GetOriginalURL(
		r.Context(),
		shortURL,
	)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(
				w,
				"link not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	response := originalURLResponse{
		URL: originalURL,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
