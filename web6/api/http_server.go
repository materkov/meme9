package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/materkov/meme9/web6/pkg"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HttpServer struct {
	Api *API
}

func (h *HttpServer) articlesList(w http.ResponseWriter, r *http.Request) {
	req := &articlesListReq{}
	_ = json.NewDecoder(r.Body).Decode(req)
	resp, err := h.Api.ArticlesList(req)
	h.writeResp(w, resp, err)
}

func (h *HttpServer) articlesSave(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("authorization")
	authToken = strings.TrimPrefix(authToken, "Bearer ")
	if authToken != pkg.GlobalConfig.SaveSecret {
		h.writeResp(w, nil, &Error{
			Code:    403,
			Message: "no access",
		})
		return
	}

	req := &InputArticle{}
	_ = json.NewDecoder(r.Body).Decode(req)
	resp, err := h.Api.ArticlesSave(req)
	h.writeResp(w, resp, err)
}

func (h *HttpServer) articlePage(w http.ResponseWriter, r *http.Request) {
	articleID, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/article/"))

	article, err := pkg.GetArticle(articleID)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	resp, _ := h.Api.ArticlesList(&articlesListReq{ID: strconv.Itoa(articleID)})

	articleImage := ""
	paragraphsHtml := ""
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

	result := "<!DOCTYPE html>"
	result += "<html><head>"
	result += "<meta charset=\"UTF-8\">"
	result += "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">"
	result += "<link href=\"/bundle/index.css\" rel=\"stylesheet\">"
	result += "<title>" + html.EscapeString(article.Title) + "</title>"
	result += fmt.Sprintf(`<meta property="og:title" content="%s" />`, html.EscapeString(article.Title))
	result += fmt.Sprintf(`<meta property="og:description" content="%s" />`, html.EscapeString(articleDescription))
	result += fmt.Sprintf(`<meta property="og:url" content="%s" />`, r.URL.RequestURI())

	if articleImage != "" {
		result += fmt.Sprintf(`<meta property="og:image" content="%s" />`, html.EscapeString(articleImage))
	}

	result += "</head><body>"
	result += "<div id=\"root\"/>"

	isSearchBot := pkg.IsSearchBot(r.UserAgent())
	if isSearchBot {
		result += "<h1>" + html.EscapeString(article.Title) + "</h1>"
		result += paragraphsHtml
	}

	respBytes, _ := json.Marshal(resp)
	result += "<script>window.__prefetchApi =" + string(respBytes) + "</script>"

	result += "<script src=\"/bundle/index.js\"></script>"
	result += "</body></html>"

	w.Write([]byte(result))
}

func (h *HttpServer) writeResp(w http.ResponseWriter, resp interface{}, err error) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Version", pkg.BuildTime)

	if err != nil {
		w.WriteHeader(400)
		log.Printf("Error: %s", err)

		code := 0
		message := ""

		var publicErr *Error
		if ok := errors.As(err, &publicErr); ok {
			code = publicErr.Code
			message = publicErr.Message
		} else {
			code = 400
			message = "Internal server error"
		}

		resp = map[string]interface{}{
			"code":    code,
			"message": message,
		}
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *HttpServer) Serve() {
	http.HandleFunc("/api/articles.list", h.articlesList)
	http.HandleFunc("/api/articles.save", h.articlesSave)
	http.HandleFunc("/article/", h.articlePage)
	http.Handle("/bundle/", http.FileServer(http.Dir("../front6/dist")))

	http.ListenAndServe("127.0.0.1:8000", nil)
}
