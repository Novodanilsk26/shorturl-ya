package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandle_POST_Success(t *testing.T) {
	originalURLs := make(map[string]string)
	for k, v := range urls {
		originalURLs[k] = v
	}
	defer func() { urls = originalURLs }()

	reqBody := "https://practicum.yandex.ru/"
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	rootHandle(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	expectedContentType := "text/plain"
	if contentType := rec.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("expected content-type %s, got %s", expectedContentType, contentType)
	}

	respBody := rec.Body.String()
	if !strings.HasPrefix(respBody, "http://localhost:8080/") {
		t.Errorf("response should start with 'http://localhost:8080/', got: %s", respBody)
	}

	shortID := strings.TrimPrefix(respBody, "http://localhost:8080/")
	if urls[shortID] != reqBody {
		t.Errorf("expected URL %s to be saved with ID %s", reqBody, shortID)
	}
}

func TestRootHandle_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	rootHandle(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestRootHandle_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()

	rootHandle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestRedirectHandle_GET_Success(t *testing.T) {
	
	testID := "test123"
	testURL := "https://example.com"
	urls[testID] = testURL

	req := httptest.NewRequest(http.MethodGet, "/"+testID, nil)
	rec := httptest.NewRecorder()

	redirectHandle(rec, req)

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != testURL {
		t.Errorf("expected location %s, got %s", testURL, location)
	}
}

func TestRedirectHandle_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test123", nil)
	rec := httptest.NewRecorder()

	redirectHandle(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestRedirectHandle_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()

	redirectHandle(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGenerateShortID(t *testing.T) {
	iterations := 100
	generated := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		id := generateShortID()
		
		if len(id) != 8 {
			t.Errorf("expected ID length 8, got %d", len(id))
		}
		
		if generated[id] {
			t.Errorf("duplicate ID generated: %s", id)
		}
		generated[id] = true
		
		for _, c := range id {
			if !(c >= '0' && c <= '9') && !(c >= 'a' && c <= 'f') {
				t.Errorf("invalid character in ID: %c", c)
				break
			}
		}
	}
}