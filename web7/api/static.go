package api

import (
	"net/http"
	"path/filepath"
	"strings"
)

func (a *API) staticHandler(w http.ResponseWriter, r *http.Request) {
	// Strip /static prefix
	path := strings.TrimPrefix(r.URL.Path, "/static/")
	if path == "" {
		http.NotFound(w, r)
		return
	}

	// Build file path relative to web7 directory
	staticDir := filepath.Join("..", "..", "front7", "dist")
	filePath := filepath.Join(staticDir, path)

	// Prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(staticDir)) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}
