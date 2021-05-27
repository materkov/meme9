package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/materkov/meme9/web/pb"
)

type Feed struct {
}

func (f *Feed) Get(ctx context.Context, request *pb.FeedGetRequest) (*pb.FeedGetResponse, error) {
	panic("implement me")
}

func (f *Feed) GetHeader(ctx context.Context, _ *pb.FeedGetHeaderRequest) (*pb.FeedGetHeaderResponse, error) {
	viewer := GetViewerFromContext(ctx)

	headerRenderer := pb.HeaderRenderer{
		MainUrl:   "/",
		LogoutUrl: "/logout",
	}

	if viewer.UserID != 0 {
		users, err := store.GetUsers([]int{viewer.UserID})
		if err != nil {
			log.Printf("Error getting user: %s", err)
		} else if len(users) == 0 {
			log.Printf("User %d not found", viewer.UserID)
		} else {
			user := users[0]

			headerRenderer.IsAuthorized = true
			headerRenderer.UserAvatar = user.VkAvatar
			headerRenderer.UserName = user.Name
		}
	}

	return &pb.FeedGetHeaderResponse{Renderer: &headerRenderer}, nil
}

type Profile struct {
}

func (p *Profile) Get(ctx context.Context, request *pb.ProfileGetRequest) (*pb.ProfileGetResponse, error) {
	panic("implement me")
}

type Relations struct {
}

func (r *Relations) Follow(ctx context.Context, req *pb.RelationsFollowRequest) (*pb.RelationsFollowResponse, error) {
	viewer := GetViewerFromContext(ctx)

	requestedID, _ := strconv.Atoi(req.UserId)
	if requestedID <= 0 {
		return nil, fmt.Errorf("incorrect follow user_id")
	}

	if viewer.UserID == 0 {
		return nil, fmt.Errorf("need auth")
	}

	followingIds, err := store.GetFollowing(viewer.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed getting following ids: %w", err)
	}

	for _, userID := range followingIds {
		if requestedID == userID {
			return &pb.RelationsFollowResponse{}, nil
		}
	}

	err = store.Follow(viewer.UserID, requestedID)
	if err != nil {
		return nil, fmt.Errorf("failed saving follower: %w", err)
	}

	return &pb.RelationsFollowResponse{}, nil
}

func (r *Relations) Unfollow(ctx context.Context, req *pb.RelationsUnfollowRequest) (*pb.RelationsUnfollowResponse, error) {
	viewer := GetViewerFromContext(ctx)

	requestedID, _ := strconv.Atoi(req.UserId)
	if requestedID <= 0 {
		return nil, fmt.Errorf("incorrect follow user_id")
	}

	if viewer.UserID == 0 {
		return nil, fmt.Errorf("need auth")
	}

	err := store.Unfollow(viewer.UserID, requestedID)
	if err != nil {
		return nil, fmt.Errorf("failed saving to store: %w", err)
	}

	return &pb.RelationsUnfollowResponse{}, nil
}

type Posts struct {
}

func (p *Posts) Add(ctx context.Context, request *pb.PostsAddRequest) (*pb.PostsAddResponse, error) {
	viewer := GetViewerFromContext(ctx)

	postID, err := store.GenerateNextID(ObjectTypePost)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %w", err)
	}

	post := Post{
		ID:     postID,
		UserID: viewer.UserID,
		Date:   int(time.Now().Unix()),
		Text:   request.Text,
	}

	err = store.AddPost(&post)
	if err != nil {
		return nil, fmt.Errorf("error saving post: %w", err)
	}

	return &pb.PostsAddResponse{
		PostUrl: fmt.Sprintf("/profile/%d", viewer.UserID),
	}, nil
}

func (p *Posts) ToggleLike(ctx context.Context, req *pb.ToggleLikeRequest) (*pb.ToggleLikeResponse, error) {
	viewer := GetViewerFromContext(ctx)

	if viewer.UserID == 0 {
		return nil, fmt.Errorf("unathorized user cannot like")
	}

	postID, _ := strconv.Atoi(req.PostId)

	isLiked, err := store.GetIsLiked([]int{postID}, viewer.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting is liked: %w", err)
	}

	if req.Action == pb.ToggleLikeRequest_LIKE && !isLiked[postID] {
		err = store.AddLike(postID, viewer.UserID)
		if err != nil {
			return nil, fmt.Errorf("error saving like: %w", err)
		}
	}

	if req.Action == pb.ToggleLikeRequest_UNLIKE && isLiked[postID] {
		err = store.DeleteLike(postID, viewer.UserID)
		if err != nil {
			return nil, fmt.Errorf("error deleting like: %w", err)
		}
	}

	postLikesCount := 0
	likesCount, err := store.GetLikesCount([]int{postID})
	if err != nil {
		log.Printf("Error getting likes count")
	} else {
		postLikesCount = likesCount[postID]
	}

	return &pb.ToggleLikeResponse{LikesCount: int32(postLikesCount)}, nil
}

func (p *Posts) AddComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.AddCommentResponse, error) {
	viewer := GetViewerFromContext(ctx)

	if viewer.UserID == 0 {
		return nil, fmt.Errorf("need auth")
	}

	if req.Text == "" {
		return nil, fmt.Errorf("text is empty")
	} else if len(req.Text) > 300 {
		return nil, fmt.Errorf("text is too long")
	}

	postID, _ := strconv.Atoi(req.PostId)
	posts, err := store.GetPosts([]int{postID})
	if err != nil {
		return nil, fmt.Errorf("error loading posts: %w", err)
	} else if len(posts) == 0 {
		return nil, fmt.Errorf("post not found")
	}

	comment := Comment{
		PostID: posts[0].ID,
		UserID: viewer.UserID,
		Text:   req.Text,
		Date:   int(time.Now().Unix()),
	}

	err = store.AddComment(&comment)
	if err != nil {
		return nil, fmt.Errorf("error saving comment: %w", err)
	}

	return &pb.AddCommentResponse{}, nil
}

type Utils struct {
}

type handler func(_ string, viewer *Viewer) (*pb.UniversalRenderer, error)

func (u *Utils) ResolveRoute(ctx context.Context, request *pb.ResolveRouteRequest) (*pb.UniversalRenderer, error) {
	viewer := GetViewerFromContext(ctx)

	routes := map[string]handler{
		`^/$`:      handleIndex,
		`^/login$`: handleLogin,

		`^/users/(\d+)$`: handleProfile,

		`^/posts/(\d+)$`: handlePostPage,
	}

	for route, handler := range routes {
		matched, _ := regexp.MatchString(route, request.Url)
		if matched {
			return handler(request.Url, viewer)
		}
	}

	return &pb.UniversalRenderer{}, nil

}
