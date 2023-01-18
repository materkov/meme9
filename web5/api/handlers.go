package api

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/auth"
	"github.com/materkov/meme9/web5/pkg/metrics"
	"github.com/materkov/meme9/web5/store"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Edges struct {
	URL        string `json:"url,omitempty"`
	TotalCount int    `json:"totalCount,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`

	Items []string `json:"items,omitempty"`
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

	method := r.URL.Query().Get("method")

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
	case "posts.getLikesConnection":
		req := PostsGetLikesConnection{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handlePostsLikesConnection(ctx, userID, &req)
	case "users.follow":
		req := UsersFollow{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		err = handleUsersFollow(ctx, userID, &req)
	case "users.unfollow":
		req := UsersUnfollow{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		err = handleUsersUnfollow(ctx, userID, &req)
	case "users.edit":
		req := UsersEdit{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		err = handleUsersEdit(ctx, userID, &req)
	case "auth.vkCallback":
		req := AuthVkCallback{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handleAuthVkCallback(ctx, userID, &req)
	case "auth.emailLogin":
		req := AuthEmailLogin{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handleAuthEmailLogin(ctx, userID, &req)
	case "auth.emailRegister":
		req := AuthEmailRegister{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handleAuthEmailRegister(ctx, userID, &req)
	case "auth.viewer":
		req := AuthViewer{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handleAuthViewer(ctx, userID, &req)
	case "users.setOnline":
		req := UsersSetOnline{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		err = handleUsersSetOnline(ctx, userID, &req)
	case "users.setAvatar":
		req := UsersSetAvatar{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handleUsersSetAvatar(ctx, userID, &req)
	case "feed.list":
		req := FeedList{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handleFeedList(ctx, userID, &req)
	case "users.posts.list":
		req := UsersPostsList{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err = handleUsersPostsList(ctx, userID, &req)
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
