package api

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/auth"
	"github.com/materkov/meme9/web5/pkg/metrics"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Edges struct {
	URL        string `json:"url,omitempty"`
	TotalCount int    `json:"totalCount,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`

	Items []string `json:"items,omitempty"`
}

func handleQuery(requestID int, viewerID int, urls []string) []interface{} {
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

	results := map[string]interface{}{}
	for len(urls) > 0 {
		url := ""
		for _, _url := range urls {
			url = _url
			urls = urls[1:]
			break
		}

		if _, alreadyResolved := results[url]; alreadyResolved {
			continue
		}

		path := url
		idx := strings.Index(path, "?")
		if idx != -1 {
			path = path[:idx]
		}

		started := time.Now()

		for _, r := range routes {
			if m, _ := regexp.MatchString("^"+r.Pattern+"$", path); m {
				localResults := r.Handler(viewerID, url)
				for _, result := range localResults {
					if related, ok := result.(string); ok {
						urls = append(urls, related)
					} else {
						results[url] = result
					}
				}
				break
			}
		}

		metrics.WriteSpan(requestID, url, time.Since(started))
	}

	resultsList := make([]interface{}, len(results))
	idx := 0
	for _, resource := range results {
		resultsList[idx] = resource
		idx++
	}

	return resultsList
}

func HandleAPI(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	requestID := rand.Int()
	defer func() {
		metrics.WriteSpan(requestID, "API Request", time.Since(started))
	}()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", fmt.Sprintf("%x", requestID))

	if r.Method == "OPTIONS" {
		return
	}

	authToken := r.Header.Get("authorization")
	authToken = strings.TrimPrefix(authToken, "Bearer ")
	userID, _ := auth.CheckToken(authToken)

	urls := strings.Split(r.URL.Query().Get("urls"), ",")
	results := handleQuery(requestID, userID, urls)

	_ = json.NewEncoder(w).Encode(results)
}
