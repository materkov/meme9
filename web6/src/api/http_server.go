package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/web6/src/pkg"
	"hash/crc32"
	"html"
	"log"
	"net/http"
	"strconv"
)

type renderOpts struct {
	Title         string
	OGDescription string
	OGUrl         string
	OGImage       string
	HTTPStatus    int

	Content  string
	Prefetch map[string]interface{}
}

func wrapPage(w http.ResponseWriter, viewer *Viewer, opts renderOpts) {
	buildTime := pkg.BuildTime
	buildCrc := strconv.Itoa(int(crc32.Checksum([]byte(buildTime), crc32.MakeTable(crc32.IEEE))))

	cssPath := fmt.Sprintf("/dist/bundle/index.css?%s", buildCrc)
	jsPath := fmt.Sprintf("/dist/bundle/index.js?%s", buildCrc)
	faviconPath := "/dist/favicon.ico?3"

	w.Header().Set("Link", fmt.Sprintf("<%s>; as=style; rel=preload, <%s>; as=image; rel=preload, <%s>; as=script; rel=preload", cssPath, faviconPath, jsPath))

	openGraph := ""
	if opts.Title != "" {
		openGraph += fmt.Sprintf(`<meta property="og:title" content="%s" />`, html.EscapeString(opts.Title))
	}
	if opts.OGDescription != "" {
		openGraph += fmt.Sprintf(`<meta property="og:description" content="%s" />`, html.EscapeString(opts.OGDescription))
	}
	if opts.OGImage != "" {
		openGraph += fmt.Sprintf(`<meta property="og:image" content="%s" />`, html.EscapeString(opts.OGImage))
	}

	title := ""
	if opts.Title != "" {
		title += "<title>" + html.EscapeString(opts.Title) + "</title>"
	}

	if opts.Prefetch == nil {
		opts.Prefetch = map[string]interface{}{}
	}

	if viewer.UserID != 0 {
		opts.Prefetch["viewerId"] = strconv.Itoa(viewer.UserID)
		opts.Prefetch["viewerName"] = viewer.UserName
	}

	if opts.HTTPStatus != 0 {
		w.WriteHeader(opts.HTTPStatus)
	}

	var prefetchBytes []byte
	if opts.Prefetch != nil {
		var err error
		prefetchBytes, err = json.Marshal(opts.Prefetch)
		if err != nil {
			log.Printf("Error marshaling to json: %s", err)
		}
	} else {
		prefetchBytes = []byte("{}")
	}

	page := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="icon" type="image/x-icon" href="%s">
	<link rel="stylesheet" href="%s">
	%s %s
</head>
<body>
	<div id="server-prefetch">
		<script>
			window.__prefetchApi = %s;
		</script>
	</div>
	<div id="server-render">%s</div>
	<div id="root"/>
	<script src="%s"></script>

</body>
</html>`

	_, _ = fmt.Fprintf(w, page,
		faviconPath,
		cssPath,
		title,
		openGraph,
		prefetchBytes,
		opts.Content,
		jsPath,
	)
}

type HttpServer struct{}

func (h *HttpServer) Serve() {
	// API
	http.HandleFunc("/api/", h.ApiHandler)
	http.HandleFunc("/upload", h.UploadHandler)

	// Web
	http.HandleFunc("/posts/", wrapWeb(h.postPage))
	http.HandleFunc("/users/", wrapWeb(h.userPage))
	http.HandleFunc("/", wrapWeb(h.discoverPage))
	http.HandleFunc("/vk-callback", wrapWeb(h.vkCallback))
	http.HandleFunc("/auth", wrapWeb(h.authPage))

	// Prometheus
	//http.Handle("/metrics", promhttp.Handler())

	// Image
	http.HandleFunc("/image-proxy", h.imageProxy2)

	// Static (for dev only)
	http.Handle("/dist/", http.FileServer(http.Dir("../front6/dist/..")))

	log.Printf("Starting http server: http://localhost:8000")

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}

func getClientIP(r *http.Request) string {
	fwdAddress := r.Header.Get("X-Forwarded-For")
	if fwdAddress != "" {
		return fwdAddress
	}

	return r.RemoteAddr
}

type ContextKey string

var CtxViewer ContextKey = "Viewer"

type webHandler func(w http.ResponseWriter, r *http.Request, viewer *Viewer)

func wrapWeb(handler webHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Version", pkg.BuildTime)

		ctx := r.Context()

		viewer := &Viewer{
			ClientIP: getClientIP(r),
		}

		authCookie, _ := r.Cookie("authToken")
		if authCookie != nil {
			authResp, err := ApiAuthClient.CheckAuth(ctx, &api.CheckAuthReq{Token: authCookie.Value})
			if err == nil {
				viewer.UserID, _ = strconv.Atoi(authResp.UserId)
				viewer.UserName = authResp.UserName
				viewer.AuthToken = authCookie.Value
				viewer.IsCookieAuth = true
			}
		}

		ctx = context.WithValue(ctx, CtxViewer, viewer)

		handler(w, r.WithContext(ctx), viewer)
	}
}

func logAPIPrefetchError(err error) {
	if err == nil {
		return
	}

	log.Printf("API Prefetch error: %s", err)
}
