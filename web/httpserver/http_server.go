package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/materkov/meme9/web/app"
	"github.com/materkov/meme9/web/pb"
	"github.com/materkov/meme9/web/store"
	"github.com/materkov/meme9/web/tracer"
	"github.com/materkov/meme9/web/utils"
)

type HttpServer struct {
	Store    *store.ObjectStore
	FeedSrv  *app.Feed
	UtilsSrv *app.Utils

	App *app.App
}

func (h *HttpServer) Serve() {
	http.Handle("/vk-callback", h.middleware(h.handleVKCallback))
	http.Handle("/logout", h.middleware(h.handleLogout))
	http.Handle("/upload", h.middleware(h.handleUpload))
	http.Handle("/api", h.middleware(h.handleAPI))
	http.Handle("/", h.middleware(h.handleDefault))

	log.Printf("[INFO] Starting http server at 8000")

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}

func (h *HttpServer) middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		viewer := app.Viewer{
			RequestHost: r.Host,
		}

		viewer.RequestScheme = "http"
		if r.Header.Get("x-forwarded-proto") == "https" {
			viewer.RequestScheme = "https"
		}

		viewer.Token, viewer.UserID, _ = h.App.TryAuth(r)

		ctx := app.WithViewerContext(r.Context(), &viewer)

		trc := tracer.NewTracer("HTTP request")
		defer trc.Stop()
		ctx = context.WithValue(ctx, utils.RequestIdKey{}, trc.TraceID)

		app.Logf(ctx, "HTTP Request: %s, user %d",
			r.URL.RequestURI(),
			viewer.UserID,
		)

		w.Header().Set("x-request-id", fmt.Sprintf("%x", trc.TraceID))

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (h *HttpServer) handleVKCallback(w http.ResponseWriter, r *http.Request) {
	viewer := app.GetViewerFromContext(r.Context())
	accessToken, err := h.App.DoVKCallback(r.Context(), r.URL.Query().Get("code"), viewer)
	if err != nil {
		log.Printf("Error: %s", err)
		_, _ = fmt.Fprint(w, "Failed to authorize")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *HttpServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *HttpServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	viewer := app.GetViewerFromContext(r.Context())
	if viewer.UserID == 0 {
		fmt.Fprint(w, "no auth")
		return
	}

	file, err := ioutil.ReadAll(r.Body)
	if err != nil || len(file) == 0 {
		fmt.Fprint(w, "no file")
		return
	}

	photo, err := h.App.UploadPhoto(file, viewer.UserID)
	if err != nil {
		log.Printf("Error uploading file: %s", err)
		fmt.Fprint(w, "error uploading file")
		return
	}

	fmt.Fprintf(w, "%d", photo.ID)
}

func (h *HttpServer) handleAPI(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := request.Context()
	method := request.URL.Query().Get("method")
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writeAPIError(w, fmt.Errorf("failed reading body"))
		return
	}

	resp, err := h.App.HandleJSONRequest(ctx, method, body)
	if err != nil {
		writeAPIError(w, fmt.Errorf("failed reading body"))
		return
	}

	_, _ = w.Write(resp)
}

func (h *HttpServer) handleDefault(w http.ResponseWriter, r *http.Request) {
	respRoute, _ := h.UtilsSrv.ResolveRoute(r.Context(), &pb.ResolveRouteRequest{Url: r.URL.Path})
	resp, _ := h.FeedSrv.GetHeader(r.Context(), nil)

	//initialDataHeader, _ := protojson.Marshal(resp)
	//initialData, _ := protojson.Marshal(respRoute)

	initialDataHeader, _ := json.Marshal(resp)
	initialData, _ := json.Marshal(respRoute)

	const page = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <title>meme</title>
    <meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
<script>
    window.initialDataHeader = %s;
    window.initialData = %s;
</script>
<div id="root"></div>
<script src="/static/App.js"></script>
</body>
</html>`

	_, _ = fmt.Fprintf(w, page, initialDataHeader, initialData)
}

func writeAPIError(w http.ResponseWriter, err error) {
	response := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}
	w.WriteHeader(400)
	_ = json.NewEncoder(w).Encode(response)
}
