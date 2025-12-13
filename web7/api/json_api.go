package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/twitchtv/twirp"
	"golang.org/x/crypto/bcrypt"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	json_api "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/json_api"
	postsservice "github.com/materkov/meme9/web7/services/posts"
)

type contextKey string

const UserIDKey contextKey = "userID"

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetFeed implements the JsonAPI GetFeed method
func (a *API) GetFeed(ctx context.Context, req *json_api.FeedRequest) (*json_api.FeedResponse, error) {
	userID := GetUserIDFromContext(ctx)

	feedType := req.Type
	if feedType == "" {
		feedType = "all"
	}

	var postsList []posts.Post
	var err error

	if feedType == "subscriptions" {
		if userID == "" {
			return nil, twirp.NewError(twirp.Unauthenticated, "authentication required for subscriptions feed")
		}

		followingIDs, err := a.subscriptions.GetFollowing(ctx, userID)
		if err != nil {
			followingIDs = []string{}
		}

		subscribedUserIDs := append(followingIDs, userID)
		postsList, err = a.posts.GetByUserIDs(ctx, subscribedUserIDs)
		if err != nil {
			return nil, twirp.NewError(twirp.Internal, err.Error())
		}
	} else {
		postsList, err = a.posts.GetAll(ctx)
		if err != nil {
			return nil, twirp.NewError(twirp.Internal, err.Error())
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
	feedPosts := make([]*json_api.FeedPostResponse, len(postsList))
	for i, post := range postsList {
		username := ""
		if user := usersMap[post.UserID]; user != nil {
			username = user.Username
		}
		feedPosts[i] = &json_api.FeedPostResponse{
			Id:        post.ID,
			Text:      post.Text,
			UserId:    post.UserID,
			Username:  username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		}
	}

	return &json_api.FeedResponse{
		Posts: feedPosts,
	}, nil
}

// Publish implements the JsonAPI Publish method
func (a *API) Publish(ctx context.Context, req *json_api.PublishRequest) (*json_api.PublishResponse, error) {
	userID := GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	post, err := a.postsService.CreatePost(ctx, req.Text, userID)
	if err != nil {
		if errors.Is(err, postsservice.ErrTextEmpty) {
			return nil, twirp.NewError(twirp.InvalidArgument, "text_empty")
		}
		if errors.Is(err, postsservice.ErrTextTooLong) {
			return nil, twirp.NewError(twirp.InvalidArgument, "text_too_long")
		}
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &json_api.PublishResponse{
		Id: post.ID,
	}, nil
}

// Login implements the JsonAPI Login method
func (a *API) Login(ctx context.Context, req *json_api.LoginRequest) (*json_api.LoginResponse, error) {
	if req.Username == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "username_required")
	}
	if req.Password == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "password_required")
	}

	user, err := a.users.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, twirp.NewError(twirp.Unauthenticated, "invalid_credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, twirp.NewError(twirp.Unauthenticated, "invalid_credentials")
	}

	tokenValue, err := generateToken()
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	_, err = a.tokens.Create(ctx, tokens.Token{
		Token:     tokenValue,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	return &json_api.LoginResponse{
		Token:    tokenValue,
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}

// Register implements the JsonAPI Register method
func (a *API) Register(ctx context.Context, req *json_api.RegisterRequest) (*json_api.LoginResponse, error) {
	if req.Username == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "username_required")
	}
	if req.Password == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "password_required")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	user, err := a.users.Create(ctx, users.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	})
	if err != nil {
		if errors.Is(err, users.ErrUsernameExists) {
			return nil, twirp.NewError(twirp.AlreadyExists, "username_exists")
		}
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	tokenValue, err := generateToken()
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	_, err = a.tokens.Create(ctx, tokens.Token{
		Token:     tokenValue,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	return &json_api.LoginResponse{
		Token:    tokenValue,
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}

// GetUserPosts implements the JsonAPI GetUserPosts method
func (a *API) GetUserPosts(ctx context.Context, req *json_api.UserPostsRequest) (*json_api.UserPostsResponse, error) {
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id_required")
	}

	postsList, err := a.posts.GetByUserID(ctx, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	user, err := a.users.GetByID(ctx, req.UserId)
	username := ""
	if err == nil && user != nil {
		username = user.Username
	}

	userPosts := make([]*json_api.UserPostResponse, len(postsList))
	for i, post := range postsList {
		userPosts[i] = &json_api.UserPostResponse{
			Id:        post.ID,
			Text:      post.Text,
			UserId:    post.UserID,
			Username:  username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		}
	}

	return &json_api.UserPostsResponse{
		Posts: userPosts,
	}, nil
}

// Subscribe implements the JsonAPI Subscribe method
func (a *API) Subscribe(ctx context.Context, req *json_api.SubscribeRequest) (*json_api.SubscribeResponse, error) {
	followerID := GetUserIDFromContext(ctx)
	if followerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
	}

	err := a.subscriptions.Subscribe(ctx, followerID, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &json_api.SubscribeResponse{
		Subscribed: true,
	}, nil
}

// Unsubscribe implements the JsonAPI Unsubscribe method
func (a *API) Unsubscribe(ctx context.Context, req *json_api.SubscribeRequest) (*json_api.SubscribeResponse, error) {
	followerID := GetUserIDFromContext(ctx)
	if followerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
	}

	err := a.subscriptions.Unsubscribe(ctx, followerID, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &json_api.SubscribeResponse{
		Subscribed: false,
	}, nil
}

// GetSubscriptionStatus implements the JsonAPI GetSubscriptionStatus method
func (a *API) GetSubscriptionStatus(ctx context.Context, req *json_api.SubscribeRequest) (*json_api.SubscribeResponse, error) {
	followerID := GetUserIDFromContext(ctx)
	if followerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
	}

	isSubscribed, err := a.subscriptions.IsSubscribed(ctx, followerID, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &json_api.SubscribeResponse{
		Subscribed: isSubscribed,
	}, nil
}

// GetFollowing implements the JsonAPI GetFollowing method
func (a *API) GetFollowing(ctx context.Context, req *json_api.GetFollowingRequest) (*json_api.GetFollowingResponse, error) {
	userID := GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	// If no user_id provided, use the authenticated user's ID
	targetUserID := req.UserId
	if targetUserID == "" {
		targetUserID = userID
	}

	// Only allow users to see their own following list, or we could allow public access
	// For now, require authentication and allow seeing own following list
	if targetUserID != userID {
		return nil, twirp.NewError(twirp.PermissionDenied, "can only view own following list")
	}

	followingIDs, err := a.subscriptions.GetFollowing(ctx, targetUserID)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &json_api.GetFollowingResponse{
		UserIds: followingIDs,
	}, nil
}

// GetPost implements the JsonAPI GetPost method
func (a *API) GetPost(ctx context.Context, req *json_api.GetPostRequest) (*json_api.GetPostResponse, error) {
	if req.PostId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "post_id is required")
	}

	post, err := a.posts.GetByID(ctx, req.PostId)
	if err != nil {
		return nil, twirp.NewError(twirp.NotFound, "post not found")
	}

	return &json_api.GetPostResponse{
		Id:        post.ID,
		Text:      post.Text,
		UserId:    post.UserID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetUser implements the JsonAPI GetUser method
func (a *API) GetUser(ctx context.Context, req *json_api.GetUserRequest) (*json_api.GetUserResponse, error) {
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
	}

	user, err := a.users.GetByID(ctx, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.NotFound, "user not found")
	}

	return &json_api.GetUserResponse{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}

// VerifyToken implements the JsonAPI VerifyToken method
func (a *API) VerifyToken(ctx context.Context, req *json_api.VerifyTokenRequest) (*json_api.VerifyTokenResponse, error) {
	if req.Token == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "token is required")
	}

	userID, err := a.tokensService.VerifyToken(ctx, req.Token)
	if err != nil {
		return nil, twirp.NewError(twirp.Unauthenticated, "invalid token")
	}

	return &json_api.VerifyTokenResponse{
		UserId: userID,
	}, nil
}

// IsSubscribed implements the JsonAPI IsSubscribed method
func (a *API) IsSubscribed(ctx context.Context, req *json_api.IsSubscribedRequest) (*json_api.IsSubscribedResponse, error) {
	if req.SubscriberId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "subscriber_id is required")
	}
	if req.TargetUserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "target_user_id is required")
	}

	isSubscribed, err := a.subscriptions.IsSubscribed(ctx, req.SubscriberId, req.TargetUserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &json_api.IsSubscribedResponse{
		Subscribed: isSubscribed,
	}, nil
}

// GetPostsByUserIDs implements the JsonAPI GetPostsByUserIDs method
func (a *API) GetPostsByUserIDs(ctx context.Context, req *json_api.GetPostsByUserIDsRequest) (*json_api.GetPostsByUserIDsResponse, error) {
	if len(req.UserIds) == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_ids is required")
	}

	postsList, err := a.posts.GetByUserIDs(ctx, req.UserIds)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	posts := make([]*json_api.GetPostResponse, len(postsList))
	for i, p := range postsList {
		posts[i] = &json_api.GetPostResponse{
			Id:        p.ID,
			Text:      p.Text,
			UserId:    p.UserID,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		}
	}

	return &json_api.GetPostsByUserIDsResponse{
		Posts: posts,
	}, nil
}

// GetUsersByIDs implements the JsonAPI GetUsersByIDs method
func (a *API) GetUsersByIDs(ctx context.Context, req *json_api.GetUsersByIDsRequest) (*json_api.GetUsersByIDsResponse, error) {
	if len(req.UserIds) == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_ids is required")
	}

	usersMap, err := a.users.GetByIDs(ctx, req.UserIds)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	result := make(map[string]*json_api.GetUserResponse)
	for userID, u := range usersMap {
		result[userID] = &json_api.GetUserResponse{
			Id:       u.ID,
			Username: u.Username,
		}
	}

	return &json_api.GetUsersByIDsResponse{
		Users: result,
	}, nil
}
