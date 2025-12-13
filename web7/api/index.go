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
	// Check if user is authenticated - if yes, redirect to feed, otherwise show auth page
	cookie := r.Header.Get("Cookie")
	isAuthenticated := false
	if cookie != "" {
		cookies := parseCookies(cookie)
		if token, ok := cookies["auth_token"]; ok && token != "" {
			_, err := a.tokensService.VerifyToken(r.Context(), "Bearer "+token)
			if err == nil {
				isAuthenticated = true
			}
		}
	}

	if isAuthenticated {
		// Redirect authenticated users to feed, preserving query parameters
		redirectURL := "/feed"
		if r.URL.RawQuery != "" {
			redirectURL += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	// Serve auth page for unauthenticated users
	a.authPageHandler(w, r)
}
