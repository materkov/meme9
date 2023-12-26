package api

import (
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
	"golang.org/x/image/draw"
	"html"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"
	"strings"
)

// TODO empty function
func (h *HttpServer) authPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	wrapPage(w, viewer, renderOpts{})
}

func (h *HttpServer) userPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/users/")

	resp1, err := h.Api.PostsList(r.Context(), viewer, &PostsListReq{ByUserID: path})
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
	resp, err := h.Api.PostsList(r.Context(), viewer, &PostsListReq{})
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
	posts, err := store2.GlobalStore.Posts.Get([]int{postID})
	if err != nil || posts[postID] == nil {
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

	post := posts[postID]

	users, _ := store2.GlobalStore.Users.Get([]int{post.UserID})
	user := users[post.UserID]
	if user != nil {
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
			"__postPagePost": transformPostBatch(r.Context(), []*store.Post{post}, viewer.UserID)[0],
		},
	})
}

func (h *HttpServer) imageProxy(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	w.Header().Set("Content-Type", "image/jpeg")
	url := r.URL.Query().Get("url")

	resp, err := http.Get(url)
	if err != nil {
		pkg.LogErr(err)
		w.WriteHeader(400)
		return
	}
	defer resp.Body.Close()

	src, err := jpeg.Decode(resp.Body)
	if err != nil {
		pkg.LogErr(err)
		w.WriteHeader(400)
		return
	}

	if src.Bounds().Size().X <= 200 {
		_ = jpeg.Encode(w, src, nil)
		return
	}

	dst := image.NewRGBA(image.Rect(0, 0, 200, 200))
	draw.BiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	_ = jpeg.Encode(w, dst, nil)
}
