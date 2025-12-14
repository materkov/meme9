package api

import (
	"context"
	"time"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/subscriptions"
	"github.com/materkov/meme9/web7/adapters/users"
	feed "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/feed"
)

type FeedService struct {
	posts         *posts.Adapter
	users         *users.Adapter
	subscriptions *subscriptions.Adapter
}

func NewFeedService(postsAdapter *posts.Adapter, usersAdapter *users.Adapter, subscriptionsAdapter *subscriptions.Adapter) *FeedService {
	return &FeedService{
		posts:         postsAdapter,
		users:         usersAdapter,
		subscriptions: subscriptionsAdapter,
	}
}

// GetFeed implements the Feed GetFeed method
func (s *FeedService) GetFeed(ctx context.Context, req *feed.FeedRequest) (*feed.FeedResponse, error) {
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

		followingIDs, err := s.subscriptions.GetFollowing(ctx, userID)
		if err != nil {
			followingIDs = []string{}
		}

		subscribedUserIDs := append(followingIDs, userID)
		postsList, err = s.posts.GetByUserIDs(ctx, subscribedUserIDs)
		if err != nil {
			return nil, twirp.NewError(twirp.Internal, err.Error())
		}
	} else {
		postsList, err = s.posts.GetAll(ctx)
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

	usersMap, err := s.users.GetByIDs(ctx, userIDs)
	if err != nil {
		usersMap = make(map[string]*users.User)
	}

	// Build feed posts
	feedPosts := make([]*feed.FeedPostResponse, len(postsList))
	for i, post := range postsList {
		username := ""
		if user := usersMap[post.UserID]; user != nil {
			username = user.Username
		}
		feedPosts[i] = &feed.FeedPostResponse{
			Id:        post.ID,
			Text:      post.Text,
			UserId:    post.UserID,
			Username:  username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		}
	}

	return &feed.FeedResponse{
		Posts: feedPosts,
	}, nil
}
