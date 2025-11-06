package api

import (
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
	CreatedAt string `json:"createdAd"`
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

func (a *API) feedHandler(w http.ResponseWriter, r *http.Request) {
	postsList, err := a.posts.GetAll(r.Context())
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
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
