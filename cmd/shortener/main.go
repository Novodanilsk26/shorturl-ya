package main

import (
	"fmt"
	"net/http"
	"strings"
	"io"
)

var urls = make(map[string]string) 

func main() {
	http.HandleFunc("/", rootHandle)
	http.HandleFunc("/{id}", redirectHandle)
	http.ListenAndServe(":8080", nil)
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	longURL := strings.TrimSpace(string(body))
	if longURL == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	shortID := generateShortID()
	urls[shortID] = longURL

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://localhost:8080/%s", shortID)
}

func redirectHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	longURL, exists := urls[id]
	if !exists {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

func getURLFromRequest(r *http.Request) (string, error) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		return "", fmt.Errorf("invalid content type")
	}

	buf := make([]byte, 1024)
	n, err := r.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}

	url := strings.TrimSpace(string(buf[:n]))
	if url == "" {
		return "", fmt.Errorf("empty URL")
	}

	return url, nil
}

func generateShortID() string {
	return fmt.Sprintf("%x", len(urls)+1)
}