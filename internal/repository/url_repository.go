package repository

import (
    "sync"
)

type URLRepository struct {
    mu   sync.RWMutex
    urls map[string]string
}

func NewURLRepository() *URLRepository {
    return &URLRepository{
        urls: make(map[string]string),
    }
}

func (r *URLRepository) Store(id, url string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.urls[id] = url
}

func (r *URLRepository) Get(id string) (string, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    url, exists := r.urls[id]
    return url, exists
}