package web

import (
	"net/http"
	"time"
)

type Logout struct{}

func (l *Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "DELETED",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}
