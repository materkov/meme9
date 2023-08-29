package api

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/pkg"
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

func wrapPage(token *pkg.AuthToken, opts renderOpts) string {
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

	if token != nil {
		opts.Prefetch["authToken"] = token.ToString()
		opts.Prefetch["viewerId"] = token.UserID
		opts.Prefetch["viewerName"] = ""

		user, _ := pkg.GetUser(token.UserID)
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
		prefetch = fmt.Sprintf("<script>window.__prefetchApi = {};</script>")
	}

	page := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link href="/bundle/index.css" rel="stylesheet">
	%s %s
</head>
<body>
	<div id="server-prefetch">%s</div>
	<div id="server-render">%s</div>
	<div id="root"/>
	<script src="/bundle/index.js"></script>

</body>
</html>`

	return fmt.Sprintf(page,
		title,
		openGraph,
		prefetch,
		opts.Content,
	)
}

type HttpServer struct {
}

func (h *HttpServer) userPage(w http.ResponseWriter, r *http.Request, token *pkg.AuthToken) {
	_, _ = fmt.Fprint(w, wrapPage(token, renderOpts{}))
}

func (h *HttpServer) discoverPage(w http.ResponseWriter, r *http.Request, token *pkg.AuthToken) {
	_, _ = fmt.Fprint(w, wrapPage(token, renderOpts{}))
}

func (h *HttpServer) vkCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		_, _ = fmt.Fprint(w, wrapPage(nil, renderOpts{Content: "VK auth fail"}))
		return
	}

	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		proto = "http"
	}

	requestURI := fmt.Sprintf("%s://%s%s", proto, r.Host, r.URL.Path)
	vkUserID, accessToken, err := pkg.ExchangeCode(code, requestURI)
	if err != nil {
		_, _ = fmt.Fprint(w, wrapPage(nil, renderOpts{Content: "VK auth fail"}))
		return
	}

	userName, err := pkg.RefreshFromVk(accessToken, vkUserID)
	if err != nil {
		_, _ = fmt.Fprint(w, wrapPage(nil, renderOpts{Content: "VK auth fail"}))
		return
	}

	userID, err := pkg.GetEdgeByUniqueKey(pkg.FakeObjVkAuth, pkg.EdgeTypeVkAuth, strconv.Itoa(vkUserID))
	if err != nil {
		_, _ = fmt.Fprint(w, wrapPage(nil, renderOpts{Content: "VK auth fail"}))
		return
	}

	if userID == 0 {
		userID, err = pkg.AddObject(pkg.ObjTypeUser, &User{
			Name: "VK Auth user",
		})
		if err != nil {
			_, _ = fmt.Fprint(w, wrapPage(nil, renderOpts{Content: "VK auth fail"}))
			return
		}

		err = pkg.AddEdge(pkg.FakeObjVkAuth, userID, pkg.EdgeTypeVkAuth, strconv.Itoa(vkUserID))
		if err != nil {
			_, _ = fmt.Fprint(w, wrapPage(nil, renderOpts{Content: "VK auth fail"}))
			return
		}
	} else {
		user, err := pkg.GetUser(userID)
		if err != nil {
			_, _ = fmt.Fprint(w, wrapPage(nil, renderOpts{Content: "VK auth fail"}))
			return
		}

		user.Name = userName

		// Already authorized
		pkg.UpdateObject(user, user.ID)
	}

	token := pkg.AuthToken{UserID: userID}

	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    token.ToString(),
		Path:     "/",
		HttpOnly: true,
	})

	_, _ = fmt.Fprint(w, wrapPage(nil, renderOpts{
		Prefetch: map[string]interface{}{
			"authToken":  token.ToString(),
			"viewerId":   strconv.Itoa(userID),
			"viewerName": fmt.Sprintf("User %d", userID),
		},
	}))
}

func (h *HttpServer) postPage(w http.ResponseWriter, r *http.Request, token *pkg.AuthToken) {
	postID, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/posts/"))

	post, err := pkg.GetPost(postID)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	user, err := pkg.GetUser(post.UserID)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	paragraphsHtml := ""
	paragraphsHtml += "<p>" + html.EscapeString(post.Text) + "</p>"

	//if !pkg.IsSearchBot(r.Header.Get("User-Agent")) {
	//	paragraphsHtml = ""
	//}

	page := wrapPage(token, renderOpts{
		Title:         "Post by " + user.Name,
		OGDescription: post.Text,
		OGImage:       "",
		Content:       paragraphsHtml,
		Prefetch: map[string]interface{}{
			"__postPagePost": transformPost(post, user),
		},
	})
	_, _ = fmt.Fprint(w, page)
}

func (h *HttpServer) Serve() {
	// API
	http.HandleFunc("/api/users.list", wrapAPI(h.usersList))

	http.HandleFunc("/api/posts.add", wrapAPI(h.PostsAdd))
	http.HandleFunc("/api/posts.list", wrapAPI(h.PostsList))
	http.HandleFunc("/api/posts.listPostedByUser", wrapAPI(h.PostsListByUser))
	http.HandleFunc("/api/posts.listById", wrapAPI(h.PostsListByID))
	http.HandleFunc("/api/posts.delete", wrapAPI(h.PostsDelete))

	// Web
	http.HandleFunc("/posts/", wrapWeb(h.postPage))
	http.HandleFunc("/users/", wrapWeb(h.userPage))
	http.HandleFunc("/", wrapWeb(h.discoverPage))
	http.HandleFunc("/vk-callback", h.vkCallback)

	// Static (for dev only)
	http.Handle("/bundle/", http.FileServer(http.Dir("../front6/dist")))

	http.ListenAndServe("127.0.0.1:8000", nil)
}
