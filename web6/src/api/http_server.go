package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/store"
	"hash/crc32"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type renderOpts struct {
	Title         string
	OGDescription string
	OGUrl         string
	OGImage       string

	Content  string
	Prefetch map[string]interface{}
}

func wrapPage(viewer *Viewer, opts renderOpts) string {
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
		opts.Prefetch["authToken"] = viewer.AuthToken
		opts.Prefetch["viewerId"] = viewer.UserID
		opts.Prefetch["viewerName"] = ""

		user, _ := store.GetUser(viewer.UserID)
		if user != nil {
			opts.Prefetch["viewerName"] = user.Name
		}
	}

	prefetch := ""
	if opts.Prefetch != nil {
		prefetchBytes, err := json.Marshal(opts.Prefetch)
		if err != nil {
			log.Printf("Error marshaling to json: %s", err)
		}
		prefetch = fmt.Sprintf("<script>window.__prefetchApi = %s</script>", prefetchBytes)
	} else {
		prefetch = fmt.Sprint("<script>window.__prefetchApi = {};</script>")
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
	<div id="server-prefetch">%s</div>
	<div id="server-render">%s</div>
	<div id="root"/>
	<script src="/dist/bundle/index.js?%s"></script>

</body>
</html>`

	return fmt.Sprintf(page,
		buildCrc,
		title,
		openGraph,
		prefetch,
		opts.Content,
		buildCrc,
	)
}

type HttpServer struct {
	Api *API
}

func (h *HttpServer) Serve() {
	// API
	http.HandleFunc("/api/users.list", wrapAPI(h.usersList))

	http.HandleFunc("/api/posts.add", wrapAPI(h.PostsAdd))
	http.HandleFunc("/api/posts.list", wrapAPI(h.PostsList))
	http.HandleFunc("/api/posts.listPostedByUser", wrapAPI(h.PostsListByUser))
	http.HandleFunc("/api/posts.listById", wrapAPI(h.PostsListByID))
	http.HandleFunc("/api/posts.delete", wrapAPI(h.PostsDelete))

	http.HandleFunc("/api/auth.login", wrapAPI(h.authLogin))
	http.HandleFunc("/api/auth.register", wrapAPI(h.authRegister))

	// Web
	http.HandleFunc("/posts/", wrapWeb(h.postPage))
	http.HandleFunc("/users/", wrapWeb(h.userPage))
	http.HandleFunc("/", wrapWeb(h.discoverPage))
	http.HandleFunc("/vk-callback", h.vkCallback)
	http.HandleFunc("/logout", h.logout)

	// Static (for dev only)
	http.Handle("/dist/", http.FileServer(http.Dir("../front6/dist/..")))

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}

func wrapAPI(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Version", pkg.BuildTime)

		userID := 0
		authHeader := r.Header.Get("authorization")
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		if authHeader != "" {
			authToken := pkg.ParseAuthToken(authHeader)
			if authToken != nil {
				userID = authToken.UserID
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxViewer, &Viewer{
			UserID: userID,
		})

		handler(w, r.WithContext(ctx))
	}
}

type webHandler func(w http.ResponseWriter, r *http.Request, viewer *Viewer)

func wrapWeb(handler webHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Version", pkg.BuildTime)

		viewer := &Viewer{}

		authCookie, _ := r.Cookie("authToken")
		if authCookie != nil {
			authToken := pkg.ParseAuthToken(authCookie.Value)
			if authToken != nil {
				viewer.UserID = authToken.UserID
				viewer.AuthToken = authCookie.Value
			}
		}

		handler(w, r, viewer)
	}
}
