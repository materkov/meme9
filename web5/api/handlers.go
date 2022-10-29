package api

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/auth"
	"github.com/materkov/meme9/web5/pkg/metrics"
	"math/rand"
	"net/http"
	"net/url"
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

func handleResource(requestID int, viewerID int, resource string) (interface{}, []string) {
	type route struct {
		Pattern string
		Handler func(viewerID int, url string) []interface{}
	}

	started := time.Now()
	defer func() {
		metrics.WriteSpan(requestID, resource, time.Since(started))
	}()

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

	parsedURL, err := url.Parse(resource)
	if err != nil {
		return nil, nil
	}

	var result interface{}
	var related []string

	for _, r := range routes {
		if m, _ := regexp.MatchString("^"+r.Pattern+"$", parsedURL.Path); m {
			localResults := r.Handler(viewerID, resource)
			for _, item := range localResults {
				if item, ok := item.(string); ok {
					related = append(related, item)
				} else {
					result = item
				}
			}
			break
		}
	}

	return result, related
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
	results := DoAsync(requestID, userID, urls)

	_ = json.NewEncoder(w).Encode(results)
}
