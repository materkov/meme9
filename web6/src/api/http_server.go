package api

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/pkg/xlog"
	"github.com/materkov/meme9/web6/src/store"
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
		opts.Prefetch["viewerName"] = ""

		user, _ := store.GetUser(viewer.UserID)
		if user != nil {
			opts.Prefetch["viewerName"] = user.Name
		}
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

	buildTime := pkg.BuildTime
	buildCrc := strconv.Itoa(int(crc32.Checksum([]byte(buildTime), crc32.MakeTable(crc32.IEEE))))

	page := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="icon" type="image/x-icon" href="/dist/favicon.ico?3">
	<link rel="stylesheet" href="/dist/bundle/index.css?%s">
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
	<script src="/dist/bundle/index.js?%s"></script>

</body>
</html>`

	_, _ = fmt.Fprintf(w, page,
		buildCrc,
		title,
		openGraph,
		prefetchBytes,
		opts.Content,
		buildCrc,
	)
}

type HttpServer struct {
	Api *API
}

func (h *HttpServer) Serve() {
	// API
	http.HandleFunc("/api/", h.ApiHandler)

	// Web
	http.HandleFunc("/posts/", wrapWeb(h.postPage))
	http.HandleFunc("/users/", wrapWeb(h.userPage))
	http.HandleFunc("/", wrapWeb(h.discoverPage))
	http.HandleFunc("/vk-callback", wrapWeb(h.vkCallback))
	http.HandleFunc("/auth", wrapWeb(h.authPage))

	// Image
	http.HandleFunc("/image-proxy", wrapWeb(h.imageProxy))

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

type webHandler func(w http.ResponseWriter, r *http.Request, viewer *Viewer)

func wrapWeb(handler webHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Version", pkg.BuildTime)

		viewer := &Viewer{
			ClientIP: getClientIP(r),
		}

		authCookie, _ := r.Cookie("authToken")
		if authCookie != nil {
			authToken := pkg.ParseAuthToken(authCookie.Value)
			if authToken != nil {
				viewer.UserID = authToken.UserID
				viewer.AuthToken = authCookie.Value
				viewer.IsCookieAuth = true
			}
		}

		xlog.Log("Processing web request", xlog.Fields{
			"url":       r.URL.String(),
			"userId":    viewer.UserID,
			"ip":        viewer.ClientIP,
			"userAgent": r.UserAgent(),
		})

		handler(w, r, viewer)
	}
}

func logAPIPrefetchError(err error) {
	if err == nil {
		return
	}

	log.Printf("API Prefetch error: %s", err)
}
