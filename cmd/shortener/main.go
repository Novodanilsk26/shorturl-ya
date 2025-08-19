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
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	e := echo.New()

	e.POST("/", rootHandle)
	e.GET("/:id", redirectHandle)

	fmt.Printf("Running server on %s\n", flagRunAddr)
	fmt.Printf("Base URL: %s\n", flagBaseURL)

	return e.Start(flagRunAddr)
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

	shortURL := fmt.Sprintf("%s/%s", flagBaseURL, shortID)

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