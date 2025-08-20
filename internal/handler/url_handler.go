package handler

import (
    "io"
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/Novodanilsk26/shorturl-ya/internal/service"
)

type Handler struct {
    service *service.URLService
}

func NewHandler(service *service.URLService) *Handler {
    return &Handler{service: service}
}

func (h *Handler) RootHandle(c echo.Context) error {
    body, err := io.ReadAll(c.Request().Body)
    if err != nil {
        return c.String(http.StatusBadRequest, "Bad request")
    }

    longURL := string(body)
    shortURL, err := h.service.CreateShortURL(longURL)
    if err != nil {
        return c.String(http.StatusBadRequest, "Bad request")
    }

    return c.String(http.StatusCreated, shortURL)
}

func (h *Handler) RedirectHandle(c echo.Context) error {
    id := c.Param("id")
    longURL, exists := h.service.GetLongURL(id)

    if !exists {
        return c.String(http.StatusNotFound, "Not found")
    }

    return c.Redirect(http.StatusTemporaryRedirect, longURL)
}