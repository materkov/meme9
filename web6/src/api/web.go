package api

import (
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/store"
	"html"
	"net/http"
	"strconv"
	"strings"
)

func (h *HttpServer) authPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{}))
}

func (h *HttpServer) userPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/users/")

	resp1, _ := h.Api.PostsListByUser(viewer, &PostsListByUserReq{UserID: path})
	resp2, _ := h.Api.usersList(viewer, &UsersListReq{UserIds: []string{path}})

	_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{
		Prefetch: map[string]interface{}{
			"__userPage": map[string]interface{}{
				"user_id": path,
				"user":    resp2[0],
				"posts":   resp1,
			},
		},
	}))

}

func (h *HttpServer) discoverPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	resp, _ := h.Api.PostsList(viewer, &PostsListReq{})

	_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{
		Prefetch: map[string]interface{}{
			"__postsList": resp,
		},
	}))
}

func (h *HttpServer) vkCallback(w http.ResponseWriter, r *http.Request) {
	viewer := &Viewer{}
	code := r.URL.Query().Get("code")
	if code == "" {
		_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{Content: "VK auth fail"}))
		return
	}

	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		proto = "http"
	}

	requestURI := fmt.Sprintf("%s://%s%s", proto, r.Host, r.URL.Path)
	vkUserID, accessToken, err := pkg.ExchangeCode(code, requestURI)
	if err != nil {
		_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{Content: "VK auth fail"}))
		return
	}

	userName, err := pkg.RefreshFromVk(accessToken, vkUserID)
	if err != nil {
		_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{Content: "VK auth fail"}))
		return
	}

	userID, err := store.GetEdgeByUniqueKey(store.FakeObjVkAuth, store.EdgeTypeVkAuth, strconv.Itoa(vkUserID))
	if err != nil {
		_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{Content: "VK auth fail"}))
		return
	}

	if userID == 0 {
		userID, err = store.AddObject(store.ObjTypeUser, &User{
			Name: "VK Auth user",
		})
		if err != nil {
			_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{Content: "VK auth fail"}))
			return
		}

		err = store.AddEdge(store.FakeObjVkAuth, userID, store.EdgeTypeVkAuth, strconv.Itoa(vkUserID))
		if err != nil {
			_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{Content: "VK auth fail"}))
			return
		}
	} else {
		user, err := store.GetUser(userID)
		if err != nil {
			_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{Content: "VK auth fail"}))
			return
		}

		user.Name = userName

		// Already authorized
		err = store.UpdateObject(user, user.ID)
		pkg.LogErr(err)
	}

	token := pkg.AuthToken{UserID: userID}

	http.SetCookie(w, &http.Cookie{
		Name:  "authToken",
		Value: token.ToString(),
		Path:  "/",
	})
	http.Redirect(w, r, "/", 302)
}

func (h *HttpServer) postPage(w http.ResponseWriter, r *http.Request, viewer *Viewer) {
	postID, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/posts/"))

	post, err := store.GetPost(postID)
	if err != nil {
		_, _ = fmt.Fprint(w, wrapPage(viewer, renderOpts{
			Title:   "Post not found",
			Content: "Post not found",
		}))
		return
	}

	user, err := store.GetUser(post.UserID)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	paragraphsHtml := ""
	paragraphsHtml += "<p>" + html.EscapeString(post.Text) + "</p>"

	//if !pkg.IsSearchBot(r.Header.Get("User-Agent")) {
	//	paragraphsHtml = ""
	//}

	page := wrapPage(viewer, renderOpts{
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
