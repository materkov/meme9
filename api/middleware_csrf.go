package api

import (
	"net/http"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/pkg/config"
)

type CSRFMiddleware struct {
	Config *config.Config
}

func (c *CSRFMiddleware) Do(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-csrf-token")
		viewer := r.Context().Value(viewerCtxKey).(*api.Viewer)

		if viewer.User != nil && token != "" {
			viewer.CSRFValidated = api.ValidateCSRFToken(c.Config.CSRFKey, viewer.User.ID, token)
		}

		next(w, r)
	}
}
