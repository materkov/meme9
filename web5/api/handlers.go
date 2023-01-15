package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/auth"
	"github.com/materkov/meme9/web5/pkg/metrics"
	"github.com/materkov/meme9/web5/store"
	"log"
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

func handleQuery(ctx context.Context, requestID int, viewerID int, urls []string) []interface{} {
	type route struct {
		Pattern string
		Handler func(ctx context.Context, viewerID int, url string) []interface{}
	}

	routes := []route{
		{"/feed", handleFeed},

		{"/users/(\\w+)", handleUserById},
		{"/users/(\\w+)/followers", handleUserFollowers},
		{"/users/(\\w+)/following", handleUserFollowing},
		{"/users/(\\w+)/posts", handleUserPosts},
		{"/users/(\\w+)/online", handleUserOnline},

		{"/posts/(\\w+)", handlePostsId},
		{"/posts/(\\w+)/liked", handlePostsLiked},

		{"/photos/(\\w+)", handlePhotosId},

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
				localResults := r.Handler(ctx, viewerID, url)
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
	requestID := rand.Int()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", fmt.Sprintf("%x", requestID))

	if r.Method == "OPTIONS" {
		return
	}

	var resources []string
	_ = json.NewDecoder(r.Body).Decode(&resources)

	started := time.Now()
	defer func() {
		metrics.WriteSpan(requestID, "API Request", time.Since(started), "resources", strings.Join(resources, ","))
	}()

	authToken := r.Header.Get("authorization")
	authToken = strings.TrimPrefix(authToken, "Bearer ")
	userID, _ := auth.CheckToken(authToken)

	resultsCh := make(chan []interface{})
	for _, resource := range resources {
		resourceCopy := resource

		go func() {
			ctx := r.Context()
			ctx = store.WithCachedStore(ctx)

			resultsCh <- handleQuery(ctx, requestID, userID, []string{resourceCopy})
		}()
	}

	var results []interface{}
	for range resources {
		results = append(results, <-resultsCh...)
	}

	_ = json.NewEncoder(w).Encode(results)
}

func HandleAPI2(w http.ResponseWriter, r *http.Request) {
	requestID := rand.Int()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", fmt.Sprintf("%x", requestID))

	if r.Method == "OPTIONS" {
		return
	}

	method := strings.TrimPrefix(r.URL.Path, "/api2/")

	started := time.Now()
	defer func() {
		metrics.WriteSpan(requestID, "API Request", time.Since(started), "method", method)
	}()

	authToken := r.Header.Get("authorization")
	authToken = strings.TrimPrefix(authToken, "Bearer ")
	userID, _ := auth.CheckToken(authToken)

	ctx := r.Context()
	ctx = store.WithCachedStore(ctx)

	var resp interface{}
	var err error

	switch method {
	case "posts.add":
		req := PostsAdd{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handlePostsAdd(ctx, userID, &req)
	case "posts.delete":
		req := PostsDelete{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		err = handlePostsDelete(ctx, userID, &req)
	case "posts.like":
		req := PostsLike{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handlePostsLike(ctx, userID, &req)
	case "posts.unlike":
		req := PostsUnlike{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handlePostsUnlike(ctx, userID, &req)
	default:
		err = fmt.Errorf("unknown method")
	}

	response := struct {
		Data  interface{} `json:"data"`
		Error string      `json:"error,omitempty"`
	}{
		Data: resp,
	}
	if err != nil {
		response.Error = err.Error()
	}

	err = json.NewEncoder(w).Encode(response)
	log.Printf("%s", err)
}
