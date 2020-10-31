package handlers

import (
	"net/http"
	"time"

	login "github.com/materkov/meme9/api/pb"
)

type LogoutApi struct{}

func (l *LogoutApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "DELETED",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})
	writeResponse(w, &login.AnyRenderer{Renderer: &login.AnyRenderer_LogoutRenderer{
		LogoutRenderer: &login.LogoutRenderer{},
	}})
}
