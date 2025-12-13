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

	// Set correct MIME type based on file extension
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".html", ".htm":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	default:
		// Let http.ServeFile detect MIME type for other files
	}

	http.ServeFile(w, r, filePath)
}

func (a *API) faviconHandler(w http.ResponseWriter, r *http.Request) {
	// Return 404 for favicon requests
	http.NotFound(w, r)
}
