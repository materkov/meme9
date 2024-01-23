package api

/*
import (
	"github.com/materkov/meme9/api/src/pkg"
	"github.com/materkov/meme9/api/src/pkg/tracer"
	"github.com/materkov/meme9/api/src/pkg/xlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strings"
)

var apiCalls = promauto.NewCounter(prometheus.CounterOpts{
	Name: "api_calls",
	Help: "The total number of processed events",
})


func (h *HttpServer) ApiHandler(w http.ResponseWriter, r *http.Request) {
	t := tracer.NewTracer("api")
	defer t.Stop()
	apiCalls.Inc()

	ctx := tracer.WithCtx(r.Context(), t)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Version", pkg.BuildTime)

	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
		return
	}

	userID := 0
	authHeader := r.Header.Get("authorization")
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader != "" {
		authToken := pkg.ParseAuthToken(ctx, authHeader)
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
}
*/
