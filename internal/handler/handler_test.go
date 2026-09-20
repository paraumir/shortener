package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/paraumir/shortener/internal/service"
	"github.com/paraumir/shortener/internal/storage"
)

func setupHandler() *Handler {
	store := storage.NewMemoryStorage()
	svc := service.New(store)

	return New(svc)
}

func TestHandler_Shorten(t *testing.T) {
	h := setupHandler()

	body := `{"url":"https://example.com"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	h.Shorten(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"status code = %d, want %d",
			recorder.Code,
			http.StatusCreated,
		)
	}

	var response shortenResponse

	err := json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.ShortURL) != 10 {
		t.Fatalf(
			"short URL length = %d, want 10",
			len(response.ShortURL),
		)
	}
}

func TestHandler_GetOriginalURL(t *testing.T) {
	h := setupHandler()

	createBody := `{"url":"https://example.com"}`

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(createBody),
	)

	createReq.Header.Set("Content-Type", "application/json")

	createRecorder := httptest.NewRecorder()

	h.Shorten(createRecorder, createReq)

	var createResponse shortenResponse

	err := json.NewDecoder(createRecorder.Body).Decode(&createResponse)
	if err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/"+createResponse.ShortURL,
		nil,
	)

	recorder := httptest.NewRecorder()

	h.GetOriginalURL(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"status code = %d, want %d",
			recorder.Code,
			http.StatusOK,
		)
	}

	var response originalURLResponse

	err = json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.URL != "https://example.com" {
		t.Fatalf(
			"URL = %q, want %q",
			response.URL,
			"https://example.com",
		)
	}
}

func TestHandler_GetOriginalURL_NotFound(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/aB123_XyZ9",
		nil,
	)

	recorder := httptest.NewRecorder()

	h.GetOriginalURL(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"status code = %d, want %d",
			recorder.Code,
			http.StatusNotFound,
		)
	}
}

func TestHandler_Shorten_InvalidJSON(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(`{"url":`),
	)

	recorder := httptest.NewRecorder()

	h.Shorten(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"status code = %d, want %d",
			recorder.Code,
			http.StatusBadRequest,
		)
	}
}

func TestHandler_Shorten_EmptyURL(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(`{"url":""}`),
	)

	recorder := httptest.NewRecorder()

	h.Shorten(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"status code = %d, want %d",
			recorder.Code,
			http.StatusBadRequest,
		)
	}
}

func TestHandler_Shorten_MethodNotAllowed(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/shorten",
		nil,
	)

	recorder := httptest.NewRecorder()

	h.Shorten(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"status code = %d, want %d",
			recorder.Code,
			http.StatusMethodNotAllowed,
		)
	}
}

func TestHandler_Shorten_SameURL(t *testing.T) {
	h := setupHandler()

	body := `{"url":"https://example.com"}`

	firstReq := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(body),
	)

	firstRecorder := httptest.NewRecorder()

	h.Shorten(firstRecorder, firstReq)

	var firstResponse shortenResponse

	err := json.NewDecoder(firstRecorder.Body).Decode(&firstResponse)
	if err != nil {
		t.Fatalf("failed to decode first response: %v", err)
	}

	secondReq := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(body),
	)

	secondRecorder := httptest.NewRecorder()

	h.Shorten(secondRecorder, secondReq)

	var secondResponse shortenResponse

	err = json.NewDecoder(secondRecorder.Body).Decode(&secondResponse)
	if err != nil {
		t.Fatalf("failed to decode second response: %v", err)
	}

	if firstResponse.ShortURL != secondResponse.ShortURL {
		t.Fatalf(
			"short URLs are different: first=%q second=%q",
			firstResponse.ShortURL,
			secondResponse.ShortURL,
		)
	}
}

func TestHandler_GetOriginalURL_InvalidShortURL(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	h.GetOriginalURL(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"status code = %d, want %d",
			recorder.Code,
			http.StatusBadRequest,
		)
	}
}
