package subscriptions

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/subscriptions-service/api"
	subscriptionsapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/subscriptions"
)

type SubscriptionsAdapter interface {
	Subscribe(ctx context.Context, followerID, followingID string) error
	Unsubscribe(ctx context.Context, followerID, followingID string) error
	IsSubscribed(ctx context.Context, followerID, followingID string) (bool, error)
	GetFollowing(ctx context.Context, userID string) ([]string, error)
}

type Service struct {
	subscriptions SubscriptionsAdapter
}

func NewService(subscriptionsAdapter SubscriptionsAdapter) *Service {
	return &Service{subscriptions: subscriptionsAdapter}
}

func (s *Service) Subscribe(ctx context.Context, req *subscriptionsapi.SubscribeRequest) (*subscriptionsapi.SubscribeResponse, error) {
	viewerID := api.GetUserIDFromContext(ctx)
	if viewerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "target_id_empty")
	}
	if viewerID == req.UserId {
		return nil, twirp.NewError(twirp.InvalidArgument, "cannot_subscribe_to_yourself")
	}

	err := s.subscriptions.Subscribe(ctx, viewerID, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &subscriptionsapi.SubscribeResponse{
		Subscribed: true,
	}, nil
}

func (s *Service) Unsubscribe(ctx context.Context, req *subscriptionsapi.SubscribeRequest) (*subscriptionsapi.SubscribeResponse, error) {
	viewerID := api.GetUserIDFromContext(ctx)
	if viewerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "target_id_empty")
	}

	err := s.subscriptions.Unsubscribe(ctx, viewerID, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &subscriptionsapi.SubscribeResponse{
		Subscribed: false,
	}, nil
}

func (s *Service) GetStatus(ctx context.Context, req *subscriptionsapi.SubscribeRequest) (*subscriptionsapi.SubscribeResponse, error) {
	viewerID := api.GetUserIDFromContext(ctx)
	if viewerID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "target_id_empty")
	}

	isSubscribed, err := s.subscriptions.IsSubscribed(ctx, viewerID, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &subscriptionsapi.SubscribeResponse{
		Subscribed: isSubscribed,
	}, nil
}
