package main

import (
	"fmt"
	"strings"
	"io"
	"crypto/rand"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"net/http"
)

var urls = make(map[string]string) 

func main() {
	e := echo.New()

	e.POST("/", rootHandle)
	e.GET("/:id", redirectHandle)

	e.Logger.Fatal(e.Start(":8080"))
}

func rootHandle(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)

	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	longURL := strings.TrimSpace(string(body))
	if longURL == "" {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	shortID := generateShortID()
	urls[shortID] = longURL

	shortURL := fmt.Sprintf("http://%s/%s", c.Request().Host, shortID)

	return c.String(http.StatusCreated, shortURL)
}

func redirectHandle(c echo.Context) error {
	id := c.Param("id")
	longURL, exists := urls[id]

	if !exists {
		return c.String(http.StatusNotFound, "Not found")
	}

	return c.Redirect(http.StatusTemporaryRedirect, longURL)
}


func generateShortID() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%x", len(urls)+1)
	}
	return hex.EncodeToString(b)
}