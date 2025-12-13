package twirp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/twitchtv/twirp"
	"golang.org/x/crypto/bcrypt"

	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	"github.com/materkov/meme9/web7/api"
	json_api "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/json_api"
)

type contextKey string

const userIDKey contextKey = "userID"

// Server implements the JsonAPI Twirp service
type Server struct {
	api *api.API
}

// NewServer creates a new Twirp server
func NewServer(api *api.API) *Server {
	return &Server{
		api: api,
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetFeed implements the GetFeed RPC
func (s *Server) GetFeed(ctx context.Context, req *json_api.FeedRequest) (*json_api.FeedResponse, error) {
	// Extract user ID from context (set by auth hook)
	userID := getUserIDFromContext(ctx)

	apiReq := api.FeedRequest{
		Type: req.Type,
	}

	feedPosts, err := s.api.GetFeed(ctx, apiReq, userID)
	if err != nil {
		if err.Error() == "authentication required for subscriptions feed" {
			return nil, twirp.NewError(twirp.Unauthenticated, "authentication required for subscriptions feed")
		}
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	// Convert to protobuf response
	posts := make([]*json_api.FeedPostResponse, len(feedPosts))
	for i, post := range feedPosts {
		posts[i] = &json_api.FeedPostResponse{
			Id:        post.ID,
			Text:      post.Text,
			UserId:    post.UserID,
			Username:  post.Username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		}
	}

	return &json_api.FeedResponse{
		Posts: posts,
	}, nil
}

// Publish implements the Publish RPC
func (s *Server) Publish(ctx context.Context, req *json_api.PublishRequest) (*json_api.PublishResponse, error) {
	userID := getUserIDFromContext(ctx)
	if userID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	apiReq := api.PublishRequest{
		Text: req.Text,
	}

	resp, err := s.api.Publish(ctx, apiReq, userID)
	if err != nil {
		if err.Error() == "text_empty" {
			return nil, twirp.NewError(twirp.InvalidArgument, "text_empty")
		}
		if err.Error() == "text_too_long" {
			return nil, twirp.NewError(twirp.InvalidArgument, "text_too_long")
		}
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &json_api.PublishResponse{
		Id: resp.ID,
	}, nil
}

// Login implements the Login RPC
func (s *Server) Login(ctx context.Context, req *json_api.LoginRequest) (*json_api.LoginResponse, error) {
	if req.Username == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "username_required")
	}
	if req.Password == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "password_required")
	}

	user, err := s.api.GetUserByUsername(ctx, req.Username)
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

	err = s.api.CreateToken(ctx, tokens.Token{
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

// Register implements the Register RPC
func (s *Server) Register(ctx context.Context, req *json_api.RegisterRequest) (*json_api.LoginResponse, error) {
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

	user, err := s.api.CreateUser(ctx, users.User{
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

	err = s.api.CreateToken(ctx, tokens.Token{
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

// GetUserPosts implements the GetUserPosts RPC
func (s *Server) GetUserPosts(ctx context.Context, req *json_api.UserPostsRequest) (*json_api.UserPostsResponse, error) {
	apiReq := api.UserPostsRequest{
		UserID: req.UserId,
	}

	userPosts, err := s.api.GetUserPosts(ctx, apiReq)
	if err != nil {
		if err.Error() == "user_id is required" {
			return nil, twirp.NewError(twirp.InvalidArgument, "user_id_required")
		}
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	posts := make([]*json_api.UserPostResponse, len(userPosts))
	for i, post := range userPosts {
		posts[i] = &json_api.UserPostResponse{
			Id:        post.ID,
			Text:      post.Text,
			UserId:    post.UserID,
			Username:  post.Username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		}
	}

	return &json_api.UserPostsResponse{
		Posts: posts,
	}, nil
}

// Subscribe implements the Subscribe RPC
func (s *Server) Subscribe(ctx context.Context, req *json_api.SubscribeRequest) (*json_api.SubscribeResponse, error) {
	followerID := getUserIDFromContext(ctx)
	if followerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	apiReq := api.SubscribeRequest{
		UserID: req.UserId,
	}

	resp, err := s.api.Subscribe(ctx, apiReq, followerID)
	if err != nil {
		if err.Error() == "user_id is required" {
			return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
		}
		if err.Error() == "unauthorized" {
			return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
		}
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &json_api.SubscribeResponse{
		Subscribed: resp.Subscribed,
	}, nil
}

// Unsubscribe implements the Unsubscribe RPC
func (s *Server) Unsubscribe(ctx context.Context, req *json_api.SubscribeRequest) (*json_api.SubscribeResponse, error) {
	followerID := getUserIDFromContext(ctx)
	if followerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	apiReq := api.SubscribeRequest{
		UserID: req.UserId,
	}

	resp, err := s.api.Unsubscribe(ctx, apiReq, followerID)
	if err != nil {
		if err.Error() == "user_id is required" {
			return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
		}
		if err.Error() == "unauthorized" {
			return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
		}
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &json_api.SubscribeResponse{
		Subscribed: resp.Subscribed,
	}, nil
}

// GetSubscriptionStatus implements the GetSubscriptionStatus RPC
func (s *Server) GetSubscriptionStatus(ctx context.Context, req *json_api.SubscribeRequest) (*json_api.SubscribeResponse, error) {
	followerID := getUserIDFromContext(ctx)
	if followerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	apiReq := api.SubscribeRequest{
		UserID: req.UserId,
	}

	resp, err := s.api.GetSubscriptionStatus(ctx, apiReq, followerID)
	if err != nil {
		if err.Error() == "user_id is required" {
			return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
		}
		if err.Error() == "unauthorized" {
			return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
		}
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &json_api.SubscribeResponse{
		Subscribed: resp.Subscribed,
	}, nil
}

// getUserIDFromContext extracts user ID from context
func getUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}
