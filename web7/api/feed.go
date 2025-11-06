package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/materkov/meme9/web7/adapters/mongo"
)

type Post struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAd"`
}

func mapMongoPostToAPIPost(mongoPost mongo.Post, username string) Post {
	return Post{
		ID:        mongoPost.ID,
		Text:      mongoPost.Text,
		UserID:    mongoPost.UserID,
		Username:  username,
		CreatedAt: mongoPost.CreatedAt.Format(time.RFC3339),
	}
}

func (a *API) feedHandler(w http.ResponseWriter, r *http.Request) {
	mongoPosts, err := a.mongo.GetAllPosts(r.Context())
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Collect unique user IDs
	userIDSet := make(map[string]bool)
	for _, post := range mongoPosts {
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
	users, err := a.mongo.GetUsersByIDs(r.Context(), userIDs)
	if err != nil {
		// Log error but continue with empty usernames
		log.Printf("Error fetching users: %v", err)
		users = make(map[string]*mongo.User)
	}

	// Build username map
	usernameMap := make(map[string]string)
	for userID, user := range users {
		usernameMap[userID] = user.Username
	}

	// Map posts to API posts with usernames
	apiPosts := make([]Post, len(mongoPosts))
	for i, mongoPost := range mongoPosts {
		username := usernameMap[mongoPost.UserID]
		apiPosts[i] = mapMongoPostToAPIPost(mongoPost, username)
	}

	json.NewEncoder(w).Encode(apiPosts)
}
