package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/users"
)

type Post struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
}

func mapPostToAPIPost(post posts.Post, username string) Post {
	return Post{
		ID:        post.ID,
		Text:      post.Text,
		UserID:    post.UserID,
		Username:  username,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}
}

type FeedRequest struct {
	Type string `json:"type"`
}

func (a *API) feedHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GET requests with feed page (no auth required)
	if r.Method == http.MethodGet {
		a.feedPageHandler(w, r)
		return
	}

	// Handle POST requests with API endpoint (auth required)
	if r.Method != http.MethodPost {
		writeErrorCode(w, "method_not_allowed", "")
		return
	}

	// Check authentication for POST requests
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	userID, authErr := a.tokensService.VerifyToken(r.Context(), authHeader)
	if authErr != nil {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	// Set userID in context for use in handler
	ctx := context.WithValue(r.Context(), userIDKey, userID)
	r = r.WithContext(ctx)

	var req FeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	// Default to "global" if not specified
	feedType := req.Type
	if feedType == "" {
		feedType = "global"
	}

	var postsList []posts.Post
	var err error

	// For subscriptions feed, require authentication
	if feedType == "subscriptions" {
		userID := getUserID(r)
		if userID == "" {
			writeErrorCode(w, "unauthorized", "Authentication required for subscriptions feed")
			return
		}

		// Get subscriptions for the current user
		followingIDs, err := a.subscriptions.GetFollowing(r.Context(), userID)
		if err != nil {
			log.Printf("Error fetching subscriptions: %v", err)
			// Continue with empty subscriptions
			followingIDs = []string{}
		}

		// Include own posts and posts from subscribed users
		subscribedUserIDs := append(followingIDs, userID)
		postsList, err = a.posts.GetByUserIDs(r.Context(), subscribedUserIDs)
		if err != nil {
			writeInternalServerError(w, "internal_server_error", "")
			return
		}
	} else {
		// Global feed - show all posts
		postsList, err = a.posts.GetAll(r.Context())
		if err != nil {
			writeInternalServerError(w, "internal_server_error", "")
			return
		}
	}

	// Collect unique user IDs
	userIDSet := make(map[string]bool)
	for _, post := range postsList {
		if post.UserID != "" {
			userIDSet[post.UserID] = true
		}
	}

	// Convert set to slice
	userIDs := make([]string, 0, len(userIDSet))
	for userID := range userIDSet {
		userIDs = append(userIDs, userID)
	}

	// Fetch all users in a single batch query
	usersMap, err := a.users.GetByIDs(r.Context(), userIDs)
	if err != nil {
		// Log error but continue with empty usernames
		log.Printf("Error fetching users: %v", err)
		usersMap = make(map[string]*users.User)
	}

	// Build username map
	usernameMap := make(map[string]string)
	for userID, user := range usersMap {
		usernameMap[userID] = user.Username
	}

	// Map posts to API posts with usernames
	apiPosts := make([]Post, len(postsList))
	for i, post := range postsList {
		username := usernameMap[post.UserID]
		apiPosts[i] = mapPostToAPIPost(post, username)
	}

	json.NewEncoder(w).Encode(apiPosts)
}
