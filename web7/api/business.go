package api

import (
	"context"
	"errors"
	"time"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	postsservice "github.com/materkov/meme9/web7/services/posts"
)

// Business logic methods (no HTTP dependencies)

// FeedRequest represents a feed request
type FeedRequest struct {
	Type string
}

// FeedPost represents a post in the feed response
type FeedPost struct {
	ID        string
	Text      string
	UserID    string
	Username  string
	CreatedAt time.Time
}

// GetFeed returns posts for the feed
func (a *API) GetFeed(ctx context.Context, req FeedRequest, userID string) ([]FeedPost, error) {
	feedType := req.Type
	if feedType == "" {
		feedType = "global"
	}

	var postsList []posts.Post
	var err error

	if feedType == "subscriptions" {
		if userID == "" {
			return nil, errors.New("authentication required for subscriptions feed")
		}

		followingIDs, err := a.subscriptions.GetFollowing(ctx, userID)
		if err != nil {
			followingIDs = []string{}
		}

		subscribedUserIDs := append(followingIDs, userID)
		postsList, err = a.posts.GetByUserIDs(ctx, subscribedUserIDs)
		if err != nil {
			return nil, err
		}
	} else {
		postsList, err = a.posts.GetAll(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Collect unique user IDs
	userIDSet := make(map[string]bool)
	for _, post := range postsList {
		if post.UserID != "" {
			userIDSet[post.UserID] = true
		}
	}

	userIDs := make([]string, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	usersMap, err := a.users.GetByIDs(ctx, userIDs)
	if err != nil {
		usersMap = make(map[string]*users.User)
	}

	// Build feed posts
	feedPosts := make([]FeedPost, len(postsList))
	for i, post := range postsList {
		username := ""
		if user := usersMap[post.UserID]; user != nil {
			username = user.Username
		}
		feedPosts[i] = FeedPost{
			ID:        post.ID,
			Text:      post.Text,
			UserID:    post.UserID,
			Username:  username,
			CreatedAt: post.CreatedAt,
		}
	}

	return feedPosts, nil
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string
	Password string
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token    string
	UserID   string
	Username string
}

// GetUserByUsername returns a user by username (for login)
func (a *API) GetUserByUsername(ctx context.Context, username string) (*users.User, error) {
	return a.users.GetByUsername(ctx, username)
}

// CreateToken creates a new token
func (a *API) CreateToken(ctx context.Context, token tokens.Token) error {
	_, err := a.tokens.Create(ctx, token)
	return err
}

// CreateUser creates a new user
func (a *API) CreateUser(ctx context.Context, user users.User) (*users.User, error) {
	return a.users.Create(ctx, user)
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string
	Password string
}

// PublishRequest represents a publish request
type PublishRequest struct {
	Text string
}

// PublishResponse represents a publish response
type PublishResponse struct {
	ID string
}

// Publish creates a new post
func (a *API) Publish(ctx context.Context, req PublishRequest, userID string) (*PublishResponse, error) {
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	post, err := a.postsService.CreatePost(ctx, req.Text, userID)
	if err != nil {
		if errors.Is(err, postsservice.ErrTextEmpty) {
			return nil, errors.New("text_empty")
		}
		if errors.Is(err, postsservice.ErrTextTooLong) {
			return nil, errors.New("text_too_long")
		}
		return nil, err
	}

	return &PublishResponse{ID: post.ID}, nil
}

// UserPostsRequest represents a user posts request
type UserPostsRequest struct {
	UserID string
}

// UserPost represents a post in user posts response
type UserPost struct {
	ID        string
	Text      string
	UserID    string
	Username  string
	CreatedAt time.Time
}

// GetUserPosts returns posts for a user
func (a *API) GetUserPosts(ctx context.Context, req UserPostsRequest) ([]UserPost, error) {
	if req.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	postsList, err := a.posts.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	user, err := a.users.GetByID(ctx, req.UserID)
	username := ""
	if err == nil && user != nil {
		username = user.Username
	}

	userPosts := make([]UserPost, len(postsList))
	for i, post := range postsList {
		userPosts[i] = UserPost{
			ID:        post.ID,
			Text:      post.Text,
			UserID:    post.UserID,
			Username:  username,
			CreatedAt: post.CreatedAt,
		}
	}

	return userPosts, nil
}

// SubscribeRequest represents a subscription request
type SubscribeRequest struct {
	UserID string
}

// SubscribeResponse represents a subscription response
type SubscribeResponse struct {
	Subscribed bool
}

// Subscribe subscribes to a user
func (a *API) Subscribe(ctx context.Context, req SubscribeRequest, followerID string) (*SubscribeResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("user_id is required")
	}
	if followerID == "" {
		return nil, errors.New("unauthorized")
	}

	err := a.subscriptions.Subscribe(ctx, followerID, req.UserID)
	if err != nil {
		return nil, err
	}

	return &SubscribeResponse{Subscribed: true}, nil
}

// Unsubscribe unsubscribes from a user
func (a *API) Unsubscribe(ctx context.Context, req SubscribeRequest, followerID string) (*SubscribeResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("user_id is required")
	}
	if followerID == "" {
		return nil, errors.New("unauthorized")
	}

	err := a.subscriptions.Unsubscribe(ctx, followerID, req.UserID)
	if err != nil {
		return nil, err
	}

	return &SubscribeResponse{Subscribed: false}, nil
}

// GetSubscriptionStatus returns subscription status
func (a *API) GetSubscriptionStatus(ctx context.Context, req SubscribeRequest, followerID string) (*SubscribeResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("user_id is required")
	}
	if followerID == "" {
		return nil, errors.New("unauthorized")
	}

	isSubscribed, err := a.subscriptions.IsSubscribed(ctx, followerID, req.UserID)
	if err != nil {
		return nil, err
	}

	return &SubscribeResponse{Subscribed: isSubscribed}, nil
}

// VerifyToken verifies a token and returns user ID
func (a *API) VerifyToken(ctx context.Context, token string) (string, error) {
	return a.tokensService.VerifyToken(ctx, token)
}

// GetFollowing returns the list of user IDs that a user is following
func (a *API) GetFollowing(ctx context.Context, userID string) ([]string, error) {
	return a.subscriptions.GetFollowing(ctx, userID)
}

// IsSubscribed checks if subscriberID is subscribed to targetUserID
func (a *API) IsSubscribed(ctx context.Context, subscriberID, targetUserID string) (bool, error) {
	return a.subscriptions.IsSubscribed(ctx, subscriberID, targetUserID)
}

// PostData represents post data for HTML rendering (no adapter dependency)
type PostData struct {
	ID        string
	Text      string
	UserID    string
	CreatedAt time.Time
}

// UserData represents user data for HTML rendering (no adapter dependency)
type UserData struct {
	ID       string
	Username string
}

// GetAllPostsHTML returns all posts as PostData
func (a *API) GetAllPostsHTML(ctx context.Context) ([]PostData, error) {
	adapterPosts, err := a.posts.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]PostData, len(adapterPosts))
	for i, p := range adapterPosts {
		result[i] = PostData{
			ID:        p.ID,
			Text:      p.Text,
			UserID:    p.UserID,
			CreatedAt: p.CreatedAt,
		}
	}
	return result, nil
}

// GetPostsByUserIDHTML returns posts for a user as PostData
func (a *API) GetPostsByUserIDHTML(ctx context.Context, userID string) ([]PostData, error) {
	adapterPosts, err := a.posts.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]PostData, len(adapterPosts))
	for i, p := range adapterPosts {
		result[i] = PostData{
			ID:        p.ID,
			Text:      p.Text,
			UserID:    p.UserID,
			CreatedAt: p.CreatedAt,
		}
	}
	return result, nil
}

// GetPostsByUserIDsHTML returns posts for multiple users as PostData
func (a *API) GetPostsByUserIDsHTML(ctx context.Context, userIDs []string) ([]PostData, error) {
	adapterPosts, err := a.posts.GetByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	result := make([]PostData, len(adapterPosts))
	for i, p := range adapterPosts {
		result[i] = PostData{
			ID:        p.ID,
			Text:      p.Text,
			UserID:    p.UserID,
			CreatedAt: p.CreatedAt,
		}
	}
	return result, nil
}

// GetPostByIDHTML returns a single post as PostData
func (a *API) GetPostByIDHTML(ctx context.Context, postID string) (*PostData, error) {
	adapterPost, err := a.posts.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	return &PostData{
		ID:        adapterPost.ID,
		Text:      adapterPost.Text,
		UserID:    adapterPost.UserID,
		CreatedAt: adapterPost.CreatedAt,
	}, nil
}

// GetUserByIDHTML returns a user as UserData
func (a *API) GetUserByIDHTML(ctx context.Context, userID string) (*UserData, error) {
	adapterUser, err := a.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &UserData{
		ID:       adapterUser.ID,
		Username: adapterUser.Username,
	}, nil
}

// GetUsersByIDsHTML returns multiple users as UserData
func (a *API) GetUsersByIDsHTML(ctx context.Context, userIDs []string) (map[string]*UserData, error) {
	adapterUsers, err := a.users.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*UserData)
	for userID, u := range adapterUsers {
		result[userID] = &UserData{
			ID:       u.ID,
			Username: u.Username,
		}
	}
	return result, nil
}
