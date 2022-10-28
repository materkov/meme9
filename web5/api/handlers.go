package api

import (
	"encoding/json"
	"github.com/materkov/meme9/web5/pkg/auth"
	"net/http"
	"regexp"
	"strings"
)

type Edges struct {
	URL        string `json:"url,omitempty"`
	TotalCount int    `json:"totalCount,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`

	Items []string `json:"items,omitempty"`
}

func handleQuery(viewerID int, url string) []interface{} {
	type route struct {
		Pattern string
		Handler func(viewerID int, url string) []interface{}
	}

	routes := []route{
		{"/feed", handleFeed},

		{"/users/(\\w+)", handleUserById},
		{"/users/(\\w+)/followers", handleUserFollowers},
		{"/users/(\\w+)/following", handleUserFollowing},
		{"/users/(\\w+)/posts", handleUserPosts},

		{"/posts/(\\w+)", handlePostsId},
		{"/posts/(\\w+)/liked", handlePostsLiked},

		{"/viewer", handleViewer},
	}

	path := url
	idx := strings.Index(path, "?")
	if idx != -1 {
		path = path[:idx]
	}

	for _, r := range routes {
		if m, _ := regexp.MatchString("^"+r.Pattern+"$", path); m {
			return r.Handler(viewerID, url)
		}
	}

	return nil
}

func HandleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	authToken := r.Header.Get("authorization")
	authToken = strings.TrimPrefix(authToken, "Bearer ")
	userID, _ := auth.CheckToken(authToken)

	urls := strings.Split(r.URL.Query().Get("urls"), ",")

	results := make([]interface{}, 0)
	for _, query := range urls {
		results = append(results, handleQuery(userID, query)...)
	}

	_ = json.NewEncoder(w).Encode(results)
}
