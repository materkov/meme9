package api

import (
	"github.com/materkov/meme9/web6/pb/github.com/materkov/meme9/api"
	"html"
	"log"
	"net/http"
	"strings"
)

// TODO empty function
func (h *HttpServer) authPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	wrapPage(w, viewer, renderOpts{})
}

func (h *HttpServer) userPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/users/")

	resp1, err := ApiPostsClient.List(r.Context(), &api.ListReq{ByUserId: path})
	logAPIPrefetchError(err)

	resp2, err := ApiUsersClient.List(r.Context(), &api.UsersListReq{UserIds: []string{path}})
	logAPIPrefetchError(err)

	if resp2 == nil || resp2.Users[0].Name == "" {
		wrapPage(w, viewer, renderOpts{
			HTTPStatus: 404,
			Content:    "User not found",
		})
	}

	wrapPage(w, viewer, renderOpts{
		Prefetch: map[string]interface{}{
			"__userPage": map[string]interface{}{
				"user_id": path,
				"user":    resp2.Users[0],
				"posts":   resp1,
			},
		},
	})
}

func (h *HttpServer) discoverPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	resp, err := ApiPostsClient.List(r.Context(), &api.ListReq{
		Type: api.FeedType_DISCOVER,
	})
	logAPIPrefetchError(err)

	wrapPage(w, viewer, renderOpts{
		Prefetch: map[string]interface{}{
			"__postsList": resp,
		},
	})
}

func (h *HttpServer) vkCallback(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	wrapPage(w, viewer, renderOpts{})
}

func (h *HttpServer) postPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	postID := strings.TrimPrefix(r.URL.Path, "/posts/")

	apiResp, err := ApiPostsClient.List(r.Context(), &api.ListReq{ById: postID})
	logAPIPrefetchError(err)
	if err != nil || len(apiResp.Items) == 0 {
		wrapPage(w, viewer, renderOpts{
			Title:      "Post not found",
			Content:    "Post not found",
			HTTPStatus: 404,
			Prefetch: map[string]interface{}{
				"__postPagePost": nil,
			},
		})
		return
	}

	post := apiResp.Items[0]

	paragraphsHtml := ""
	paragraphsHtml += "<p>" + html.EscapeString(post.Text) + "</p>"

	//if !pkg.IsSearchBot(r.Header.Get("User-Agent")) {
	//	paragraphsHtml = ""
	//}

	wrapPage(w, viewer, renderOpts{
		Title:         "Post by " + post.User.Name,
		OGDescription: post.Text,
		OGImage:       "",
		Content:       paragraphsHtml,
		Prefetch: map[string]interface{}{
			"__postPagePost": post,
		},
	})
}

func (h *HttpServer) imageProxy2(w http.ResponseWriter, r *http.Request) {
	resp, err := ImageProxyClient.Resize(r.Context(), &api.ResizeReq{
		ImageUrl: r.URL.Query().Get("url"),
	})
	if err != nil {
		log.Printf("[ERROR] Error: %s", err)
		w.WriteHeader(400)
	} else {
		w.Header().Set("Content-Type", "image/jpeg")
		_, _ = w.Write(resp.Image)
	}
}
