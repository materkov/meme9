package api

import (
	"net/http"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/pb"
)

type CSRFMiddleware struct {
}

func (c *CSRFMiddleware) Do(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		viewer := r.Context().Value(viewerCtxKey).(*api.Viewer)
		if viewer.User != nil {
			token := r.Header.Get("x-csrf-token")

			tokenValid := api.ValidateCSRFToken(viewer.User.ID, token)
			if !tokenValid {
				err := &pb.ErrorRenderer{
					ErrorCode:   "CSRF_VALIDATION_FAILED",
					DisplayText: "Error validating CSRF token",
				}
				writeResponse(w, err)
				return
			}
		}

		next(w, r)
	}
}
