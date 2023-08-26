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
	Prefetch interface{}
}

func wrapPage(opts renderOpts) string {
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

func (h *HttpServer) userPage(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, wrapPage(renderOpts{}))
}

func (h *HttpServer) discoverPage(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, wrapPage(renderOpts{}))
}

func (h *HttpServer) vkCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		_, _ = fmt.Fprint(w, wrapPage(renderOpts{Content: "VK auth fail"}))
		return
	}

	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		proto = "http"
	}

	requestURI := fmt.Sprintf("%s://%s%s", proto, r.Host, r.URL.Path)
	vkUserID, accessToken, err := pkg.ExchangeCode(code, requestURI)
	if err != nil {
		_, _ = fmt.Fprint(w, wrapPage(renderOpts{Content: "VK auth fail"}))
		return
	}

	userName, err := pkg.RefreshFromVk(accessToken, vkUserID)
	if err != nil {
		_, _ = fmt.Fprint(w, wrapPage(renderOpts{Content: "VK auth fail"}))
		return
	}

	userID, err := pkg.GetEdgeByUniqueKey(pkg.FakeObjVkAuth, pkg.EdgeTypeVkAuth, strconv.Itoa(vkUserID))
	if err != nil {
		_, _ = fmt.Fprint(w, wrapPage(renderOpts{Content: "VK auth fail"}))
		return
	}

	if userID == 0 {
		userID, err = pkg.AddObject(pkg.ObjTypeUser, &User{
			Name: "VK Auth user",
		})
		if err != nil {
			_, _ = fmt.Fprint(w, wrapPage(renderOpts{Content: "VK auth fail"}))
			return
		}

		err = pkg.AddEdge(pkg.FakeObjVkAuth, userID, pkg.EdgeTypeVkAuth, strconv.Itoa(vkUserID))
		if err != nil {
			_, _ = fmt.Fprint(w, wrapPage(renderOpts{Content: "VK auth fail"}))
			return
		}
	} else {
		user, err := pkg.GetUser(userID)
		if err != nil {
			_, _ = fmt.Fprint(w, wrapPage(renderOpts{Content: "VK auth fail"}))
			return
		}

		user.Name = userName

		// Already authorized
		pkg.UpdateObject(user, user.ID)
	}

	token := pkg.AuthToken{UserID: userID}

	_, _ = fmt.Fprint(w, wrapPage(renderOpts{
		Prefetch: map[string]interface{}{
			"authToken":  token.ToString(),
			"viewerId":   strconv.Itoa(userID),
			"viewerName": fmt.Sprintf("User %d", userID),
		},
	}))
}

func (h *HttpServer) articlePage(w http.ResponseWriter, r *http.Request) {
	articleID, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/article/"))

	article, err := pkg.GetArticle(articleID)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	articleImage := ""
	paragraphsHtml := ""
	paragraphsHtml += "<h1>" + html.EscapeString(article.Title) + "</h1>"

	articleDescription := ""
	for _, paragraph := range article.Paragraphs {
		if paragraph.ParagraphImage != nil {
			paragraphsHtml += "<img src=\"" + html.EscapeString(paragraph.ParagraphImage.URL) + "\" />"

			if articleImage == "" {
				articleImage = paragraph.ParagraphImage.URL
			}
		} else if paragraph.ParagraphText != nil {
			paragraphsHtml += "<p>" + html.EscapeString(paragraph.ParagraphText.Text) + "</p>"

			if articleDescription == "" {
				articleDescription = paragraph.ParagraphText.Text
			}
		}
	}

	wrappedArticle := transformArticle(strconv.Itoa(article.ID), article)

	if !pkg.IsSearchBot(r.Header.Get("User-Agent")) {
		paragraphsHtml = ""
	}

	page := wrapPage(renderOpts{
		Title:         article.Title,
		OGDescription: articleDescription,
		OGImage:       articleImage,
		Content:       paragraphsHtml,
		Prefetch:      wrappedArticle,
	})
	_, _ = fmt.Fprint(w, page)
}

func (h *HttpServer) Serve() {
	// API
	http.HandleFunc("/api/users.list", wrapAPI(h.usersList))
	http.HandleFunc("/api/articles.list", wrapAPI(h.ArticlesList))
	http.HandleFunc("/api/articles.listPostedByUser", wrapAPI(h.listPostedByUser))
	http.HandleFunc("/api/articles.save", wrapAPI(h.ArticlesSave))
	http.HandleFunc("/api/articles.lastPosted", wrapAPI(h.ArticlesLastPosted))

	http.HandleFunc("/api/posts.add", wrapAPI(h.PostsAdd))
	http.HandleFunc("/api/posts.list", wrapAPI(h.PostsList))

	// Web
	http.HandleFunc("/article/", h.articlePage)
	http.HandleFunc("/users/", h.userPage)
	http.HandleFunc("/", h.discoverPage)
	http.HandleFunc("/vk-callback", h.vkCallback)

	// Static (for dev only)
	http.Handle("/bundle/", http.FileServer(http.Dir("../front6/dist")))

	http.ListenAndServe("127.0.0.1:8000", nil)
}
