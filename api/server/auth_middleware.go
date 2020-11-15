package server

import (
	"net/http"

	"github.com/materkov/meme9/api/pkg"
	"github.com/materkov/meme9/api/pkg/api"
	"github.com/materkov/meme9/api/pkg/csrf"
)

type csrfMiddleware struct {
	Config *pkg.Config
}

func (c *csrfMiddleware) Do(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-csrf-token")
		viewer := r.Context().Value(viewerCtxKey).(*api.Viewer)

		if viewer.User != nil && token != "" {
			viewer.CSRFValidated = csrf.ValidateToken(c.Config.CSRFKey, viewer.User.ID, token)
		}

		next(w, r)
	}
}
