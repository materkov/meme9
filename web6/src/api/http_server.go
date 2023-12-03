package api

import (
	"context"
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
	"strings"
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

func (h *HttpServer) wrapAPI(w http.ResponseWriter, r *http.Request) {
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

	viewer := &Viewer{
		UserID:   userID,
		ClientIP: getClientIP(r),
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxViewer, viewer)

	xlog.Log("Processing API request", xlog.Fields{
		"url":       r.URL.String(),
		"userId":    viewer.UserID,
		"ip":        viewer.ClientIP,
		"userAgent": r.UserAgent(),
	})

	method := strings.TrimPrefix(r.URL.Path, "/api/")

	switch method {
	case "posts.add":
		req := &PostsAddReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsAdd(viewer, req)
		writeResp(w, resp, err)

	case "posts.list":
		req := &PostsListReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsList(viewer, req)
		writeResp(w, resp, err)

	case "posts.listById":
		req := &PostsListByIdReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsListByID(viewer, req)
		writeResp(w, resp, err)

	case "posts.listByUser":
		req := &PostsListByUserReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsListByUser(viewer, req)
		writeResp(w, resp, err)

	case "posts.delete":
		req := &PostsDeleteReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsDelete(viewer, req)
		writeResp(w, resp, err)

	case "posts.like":
		req := &PostsLikeReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsLike(viewer, req)
		writeResp(w, resp, err)

	case "users.list":
		req := &UsersListReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.usersList(viewer, req)
		writeResp(w, resp, err)

	case "users.setStatus":
		req := &UsersSetStatus{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.usersSetStatus(viewer, req)
		writeResp(w, resp, err)

	case "users.follow":
		req := &UsersFollow{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.usersFollow(viewer, req)
		writeResp(w, resp, err)

	case "auth.login":
		req := &AuthEmailReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.authLogin(viewer, req)
		writeResp(w, resp, err)

	case "auth.register":
		req := &AuthEmailReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.authRegister(viewer, req)
		writeResp(w, resp, err)

	case "auth.vk":
		req := &AuthVkReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.authVk(viewer, req)
		writeResp(w, resp, err)

	case "polls.add":
		req := &PollsAddReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PollsAdd(viewer, req)
		writeResp(w, resp, err)

	case "polls.list":
		req := &PollsListReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PollsList(viewer, req)
		writeResp(w, resp, err)

	case "polls.vote":
		req := &PollsVoteReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PollsVote(viewer, req)
		writeResp(w, resp, err)

	case "polls.deleteVote":
		req := &PollsDeleteVoteReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PollsDeleteVote(viewer, req)
		writeResp(w, resp, err)

	default:
		writeResp(w, nil, Error("UnknownMethod"))
	}
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
