package main

import (
    "fmt"

    "github.com/labstack/echo/v4"
    "github.com/Novodanilsk26/shorturl-ya/internal/handler"
    "github.com/Novodanilsk26/shorturl-ya/internal/repository"
    "github.com/Novodanilsk26/shorturl-ya/internal/service"
)

func main() {
    cfg := parseFlags()

    if err := run(cfg); err != nil {
        panic(err)
    }
}

func run(cfg *Config) error {
    e := echo.New()

    repo := repository.NewURLRepository()
    urlService := service.NewURLService(repo, cfg.BaseURL)
    urlHandler := handler.NewHandler(urlService)

    e.POST("/", urlHandler.RootHandle)
    e.GET("/:id", urlHandler.RedirectHandle)

    fmt.Printf("Running server on %s\n", cfg.RunAddr)
    fmt.Printf("Base URL: %s\n", cfg.BaseURL)

    return e.Start(cfg.RunAddr)
}