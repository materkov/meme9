package api

import (
	"fmt"
	"net/http"
)

const apiHost = ""
const staticHost = "/static"

func indexHTML() string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>meme9</title>
  <link rel="stylesheet" href="%s/index.css">
</head>
<body>
  <script>
    window.API_BASE_URL = "%s";
  </script>
  <div id="root"></div>
  <script src="%s/index.js"></script>
</body>
</html>`, staticHost, apiHost, staticHost)
}

func (a *API) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serve index.html for all routes (client-side routing)
	// API routes and static files are handled by other handlers registered before this one
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(indexHTML()))
}
