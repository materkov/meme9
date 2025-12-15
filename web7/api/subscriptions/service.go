package subscriptions

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/adapters/subscriptions"
	"github.com/materkov/meme9/web7/api"
	subscriptionsapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/subscriptions"
)

type Service struct {
	subscriptions *subscriptions.Adapter
}

func NewService(subscriptionsAdapter *subscriptions.Adapter) *Service {
	return &Service{
		subscriptions: subscriptionsAdapter,
	}
}

// Subscribe implements the Subscriptions Subscribe method
func (s *Service) Subscribe(ctx context.Context, req *subscriptionsapi.SubscribeRequest) (*subscriptionsapi.SubscribeResponse, error) {
	followerID := api.GetUserIDFromContext(ctx)
	if followerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
	}

	err := s.subscriptions.Subscribe(ctx, followerID, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &subscriptionsapi.SubscribeResponse{
		Subscribed: true,
	}, nil
}

// Unsubscribe implements the Subscriptions Unsubscribe method
func (s *Service) Unsubscribe(ctx context.Context, req *subscriptionsapi.SubscribeRequest) (*subscriptionsapi.SubscribeResponse, error) {
	followerID := api.GetUserIDFromContext(ctx)
	if followerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
	}

	err := s.subscriptions.Unsubscribe(ctx, followerID, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &subscriptionsapi.SubscribeResponse{
		Subscribed: false,
	}, nil
}

// GetStatus implements the Subscriptions GetStatus method
func (s *Service) GetStatus(ctx context.Context, req *subscriptionsapi.SubscribeRequest) (*subscriptionsapi.SubscribeResponse, error) {
	followerID := api.GetUserIDFromContext(ctx)
	if followerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
	}

	isSubscribed, err := s.subscriptions.IsSubscribed(ctx, followerID, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &subscriptionsapi.SubscribeResponse{
		Subscribed: isSubscribed,
	}, nil
}

// GetFollowing implements the Subscriptions GetFollowing method
func (s *Service) GetFollowing(ctx context.Context, req *subscriptionsapi.GetFollowingRequest) (*subscriptionsapi.GetFollowingResponse, error) {
	userID := api.GetUserIDFromContext(ctx)
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

	followingIDs, err := s.subscriptions.GetFollowing(ctx, targetUserID)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &subscriptionsapi.GetFollowingResponse{
		UserIds: followingIDs,
	}, nil
}

// IsSubscribed implements the Subscriptions IsSubscribed method
func (s *Service) IsSubscribed(ctx context.Context, req *subscriptionsapi.IsSubscribedRequest) (*subscriptionsapi.IsSubscribedResponse, error) {
	if req.SubscriberId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "subscriber_id is required")
	}
	if req.TargetUserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "target_user_id is required")
	}

	isSubscribed, err := s.subscriptions.IsSubscribed(ctx, req.SubscriberId, req.TargetUserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &subscriptionsapi.IsSubscribedResponse{
		Subscribed: isSubscribed,
	}, nil
}
