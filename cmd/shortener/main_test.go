package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
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

	e := echo.New()
	c := e.NewContext(req, rec)

	err := rootHandle(c)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	respBody := rec.Body.String()
	if !strings.HasPrefix(respBody, "http://") {
		t.Errorf("Response should start with 'http://', got: %s", respBody)
	}

	shortID := respBody[strings.LastIndex(respBody, "/")+1:]
	if urls[shortID] != reqBody {
		t.Errorf("Expected URL %s to be saved with ID %s", reqBody, shortID)
	}
}

func TestRootHandle_InvalidMethod(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()

    e := echo.New()
    c := e.NewContext(req, rec)

    err := rootHandle(c)
    if err != nil {
        t.Fatal(err)
    }

    if rec.Code != http.StatusBadRequest {
        t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
    }
}

func TestRedirectHandle_GET_Success(t *testing.T) {
	testID := "test123"
	testURL := "https://example.com"
	urls[testID] = testURL

	req := httptest.NewRequest(http.MethodGet, "/"+testID, nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(testID)

	err := redirectHandle(c)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != testURL {
		t.Errorf("Expected location %s, got %s", testURL, location)
	}
}

func TestRedirectHandle_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	err := redirectHandle(c)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGenerateShortID(t *testing.T) {
	iterations := 100
	generated := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		id := generateShortID()
		
		if len(id) != 8 {
			t.Errorf("Expected ID length 8, got %d", len(id))
		}
		
		if generated[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		generated[id] = true
		
		for _, c := range id {
			if !(c >= '0' && c <= '9') && !(c >= 'a' && c <= 'f') {
				t.Errorf("Invalid character in ID: %c", c)
				break
			}
		}
	}
}