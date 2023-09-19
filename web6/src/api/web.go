package api

import (
	"github.com/materkov/meme9/web6/src/store"
	"html"
	"net/http"
	"strconv"
	"strings"
)

func (h *HttpServer) authPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	wrapPage(w, viewer, renderOpts{})
}

func (h *HttpServer) userPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/users/")

	resp1, err := h.Api.PostsListByUser(viewer, &PostsListByUserReq{UserID: path})
	logAPIPrefetchError(err)

	resp2, err := h.Api.usersList(viewer, &UsersListReq{UserIds: []string{path}})
	logAPIPrefetchError(err)

	// TODO think about this
	if resp2[0].Name == "" {
		wrapPage(w, viewer, renderOpts{
			HTTPStatus: 404,
			Content:    "User not found",
		})
	}

	wrapPage(w, viewer, renderOpts{
		Prefetch: map[string]interface{}{
			"__userPage": map[string]interface{}{
				"user_id": path,
				"user":    resp2[0],
				"posts":   resp1,
			},
		},
	})
}

func (h *HttpServer) discoverPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	resp, err := h.Api.PostsList(viewer, &PostsListReq{})
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
	postID, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/posts/"))

	// TODO think about this
	post, err := store.GetPost(postID)
	if err != nil {
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

	user, err := store.GetUser(post.UserID)
	if err != nil {
		wrapPage(w, viewer, renderOpts{
			Title:   "Internal error",
			Content: "Internal error",
		})
		return
	}

	paragraphsHtml := ""
	paragraphsHtml += "<p>" + html.EscapeString(post.Text) + "</p>"

	//if !pkg.IsSearchBot(r.Header.Get("User-Agent")) {
	//	paragraphsHtml = ""
	//}

	wrapPage(w, viewer, renderOpts{
		Title:         "Post by " + user.Name,
		OGDescription: post.Text,
		OGImage:       "",
		Content:       paragraphsHtml,
		Prefetch: map[string]interface{}{
			"__postPagePost": transformPost(post, user),
		},
	})
}
