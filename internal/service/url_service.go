package service

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "net/url"
    "strings"

    "github.com/Novodanilsk26/shorturl-ya/internal/repository"
)

type URLService struct {
    repo    *repository.URLRepository
    baseURL string
}

func NewURLService(repo *repository.URLRepository, baseURL string) *URLService {
    return &URLService{
        repo:    repo,
        baseURL: strings.TrimSuffix(baseURL, "/"),
    }
}

func (s *URLService) CreateShortURL(longURL string) (string, error) {
    longURL = strings.TrimSpace(longURL)
    if longURL == "" {
        return "", fmt.Errorf("empty URL")
    }

    if !isValidURL(longURL) {
        return "", fmt.Errorf("invalid URL")
    }

    shortID := s.generateShortID()
    s.repo.Store(shortID, longURL)

    return fmt.Sprintf("%s/%s", s.baseURL, shortID), nil
}

func (s *URLService) GetLongURL(shortID string) (string, bool) {
    return s.repo.Get(shortID)
}

func (s *URLService) generateShortID() string {
    b := make([]byte, 4)
    _, err := rand.Read(b)
    if err != nil {
        return fmt.Sprintf("%x", len(s.getAllURLs())+1)
    }
    return hex.EncodeToString(b)
}

func (s *URLService) getAllURLs() map[string]string {
    return make(map[string]string)
}

func isValidURL(u string) bool {
    parsed, err := url.ParseRequestURI(u)
    if err != nil {
        return false
    }
    return parsed.Scheme != "" && parsed.Host != ""
}