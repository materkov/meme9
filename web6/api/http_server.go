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

	// Web
	http.HandleFunc("/article/", h.articlePage)
	http.HandleFunc("/users/", h.userPage)
	http.HandleFunc("/", h.discoverPage)

	// Static (for dev only)
	http.Handle("/bundle/", http.FileServer(http.Dir("../front6/dist")))

	http.ListenAndServe("127.0.0.1:8000", nil)
}
