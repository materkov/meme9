package apiwrapper

import (
	"net/http"
	"path/filepath"
	"strings"
)

func (r *Router) StaticHandler(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/static/")
	if path == "" {
		http.NotFound(w, req)
		return
	}

	staticDir := filepath.Join("..", "..", "front7", "dist")
	filePath := filepath.Join(staticDir, path)

	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(staticDir)) {
		http.NotFound(w, req)
		return
	}

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
	}

	http.ServeFile(w, req, filePath)
}

func (r *Router) FaviconHandler(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
}
