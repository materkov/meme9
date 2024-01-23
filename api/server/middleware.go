package server

import (
	"context"
	"github.com/materkov/meme9/api/src/pkg"
	"github.com/materkov/meme9/api/src/pkg/tracer"
	"github.com/materkov/meme9/api/src/pkg/xlog"
	"net/http"
	"strings"
)

type Viewer struct {
	UserID       int
	AuthToken    string
	IsCookieAuth bool
	ClientIP     string
}

func getClientIP(r *http.Request) string {
	fwdAddress := r.Header.Get("X-Forwarded-For")
	if fwdAddress != "" {
		return fwdAddress
	}

	return r.RemoteAddr
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := tracer.NewTracer("api")
		defer t.Stop()

		ctx1 := tracer.WithCtx(r.Context(), t)

		userID := 0
		authHeader := r.Header.Get("authorization")
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		if authHeader != "" {
			authToken := pkg.ParseAuthToken(ctx1, authHeader)
			if authToken != nil {
				userID = authToken.UserID
			}
		}

		viewer := &Viewer{
			UserID:   userID,
			ClientIP: getClientIP(r),
		}

		xlog.Log("Processing API request", xlog.Fields{
			"url":       r.URL.String(),
			"userId":    viewer.UserID,
			"ip":        viewer.ClientIP,
			"userAgent": r.UserAgent(),
		})

		method := strings.TrimPrefix(r.URL.Path, "/api/")
		t.Tags["method"] = method

		ctx2 := context.WithValue(ctx1, CtxViewerKey, viewer)
		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}
