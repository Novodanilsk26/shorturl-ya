package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/Novodanilsk26/shorturl-ya/internal/handler"
    "github.com/Novodanilsk26/shorturl-ya/internal/repository"
    "github.com/Novodanilsk26/shorturl-ya/internal/service"
)

func TestRootHandle_POST_Success(t *testing.T) {
    repo := repository.NewURLRepository()
    service := service.NewURLService(repo, "http://localhost:8080")
    handler := handler.NewHandler(service)

    reqBody := "https://practicum.yandex.ru/"
    req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
    rec := httptest.NewRecorder()

    e := echo.New()
    c := e.NewContext(req, rec)

    err := handler.RootHandle(c)
    if err != nil {
        t.Fatalf("Handler returned error: %v", err)
    }

    if rec.Code != http.StatusCreated {
        t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
    }

    respBody := rec.Body.String()
    if !strings.HasPrefix(respBody, "http://localhost:8080/") {
        t.Errorf("Response should start with 'http://localhost:8080/', got: %s", respBody)
    }
}

func TestRootHandle_InvalidMethod(t *testing.T) {
    repo := repository.NewURLRepository()
    service := service.NewURLService(repo, "http://localhost:8080")
    handler := handler.NewHandler(service)

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()

    e := echo.New()
    c := e.NewContext(req, rec)

    err := handler.RootHandle(c)
    if err != nil {
        t.Fatal(err)
    }

    if rec.Code != http.StatusBadRequest {
        t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
    }
}

func TestRedirectHandle_GET_Success(t *testing.T) {
    repo := repository.NewURLRepository()
    service := service.NewURLService(repo, "http://localhost:8080")
    handler := handler.NewHandler(service)

    testID := "test123"
    testURL := "https://example.com"
    repo.Store(testID, testURL)

    req := httptest.NewRequest(http.MethodGet, "/"+testID, nil)
    rec := httptest.NewRecorder()

    e := echo.New()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(testID)

    err := handler.RedirectHandle(c)
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
    repo := repository.NewURLRepository()
    service := service.NewURLService(repo, "http://localhost:8080")
    handler := handler.NewHandler(service)

    req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
    rec := httptest.NewRecorder()

    e := echo.New()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("nonexistent")

    err := handler.RedirectHandle(c)
    if err != nil {
        t.Fatal(err)
    }

    if rec.Code != http.StatusNotFound {
        t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
    }
}

