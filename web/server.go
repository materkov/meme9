package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/materkov/meme9/web/pb"
)

type Feed struct {
}

func GenerateCSRFToken(token string) string {
	mac := hmac.New(sha256.New, []byte(config.CSRFKey))
	_, _ = mac.Write([]byte(token))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (f *Feed) GetHeader(ctx context.Context, _ *pb.FeedGetHeaderRequest) (*pb.FeedGetHeaderResponse, error) {
	viewer := GetViewerFromContext(ctx)

	headerRenderer := pb.HeaderRenderer{
		MainUrl:   "/",
		LoginUrl:  "/login",
		LogoutUrl: "/logout",
	}

	if viewer.UserID != 0 {
		if viewer.Token != nil {
			headerRenderer.CsrfToken = GenerateCSRFToken(viewer.Token.Token)
		}

		users, err := store.User.Get([]int{viewer.UserID})
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

	followingIds, err := store.Followers.GetFollowing(viewer.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed getting following ids: %w", err)
	}

	for _, userID := range followingIds {
		if requestedID == userID {
			return &pb.RelationsFollowResponse{}, nil
		}
	}

	objectID, err := store.GenerateNextID(ObjectTypeFollower)
	if err != nil {
		return nil, fmt.Errorf("error generating object id: %w", err)
	}

	err = store.Followers.Add(&Followers{
		ID:         objectID,
		User1ID:    viewer.UserID,
		User2ID:    requestedID,
		FollowDate: int(time.Now().Unix()),
	})
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

	err := store.Followers.Unfollow(viewer.UserID, requestedID)
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

	photoID := 0
	if request.PhotoId != "" {
		photoID, _ = strconv.Atoi(request.PhotoId)
		photos, err := store.Photo.Get([]int{photoID})
		if err != nil {
			return nil, fmt.Errorf("error getting photos: %w", err)
		} else if len(photos) == 0 {
			return nil, fmt.Errorf("photo not found")
		}

		if photos[0].UserID != viewer.UserID {
			return nil, fmt.Errorf("photo from another user")
		}
	}

	post := Post{
		ID:      postID,
		UserID:  viewer.UserID,
		Date:    int(time.Now().Unix()),
		Text:    request.Text,
		PhotoID: photoID,
	}

	err = store.Post.Add(&post)
	if err != nil {
		return nil, fmt.Errorf("error saving post: %w", err)
	}

	return &pb.PostsAddResponse{
		PostUrl: fmt.Sprintf("/posts/%d", post.ID),
	}, nil
}

func (p *Posts) ToggleLike(ctx context.Context, req *pb.ToggleLikeRequest) (*pb.ToggleLikeResponse, error) {
	viewer := GetViewerFromContext(ctx)

	if viewer.UserID == 0 {
		return nil, fmt.Errorf("unathorized user cannot like")
	}

	postID, _ := strconv.Atoi(req.PostId)

	isLiked, err := store.Likes.GetIsLiked([]int{postID}, viewer.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting is liked: %w", err)
	}

	if req.Action == pb.ToggleLikeRequest_LIKE && !isLiked[postID] {
		objectID, err := store.GenerateNextID(ObjectTypeLike)
		if err != nil {
			return nil, fmt.Errorf("error generating object id: %w", err)
		}

		err = store.Likes.Add(&Likes{
			ID:     objectID,
			PostID: postID,
			UserID: viewer.UserID,
			Time:   int(time.Now().Unix()),
		})
		if err != nil {
			return nil, fmt.Errorf("error saving like: %w", err)
		}
	}

	if req.Action == pb.ToggleLikeRequest_UNLIKE && isLiked[postID] {
		err = store.Likes.Delete(postID, viewer.UserID)
		if err != nil {
			return nil, fmt.Errorf("error deleting like: %w", err)
		}
	}

	postLikesCount := 0
	likesCount, err := store.Likes.GetCount([]int{postID})
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
	posts, err := store.Post.Get([]int{postID})
	if err != nil {
		return nil, fmt.Errorf("error loading posts: %w", err)
	} else if len(posts) == 0 {
		return nil, fmt.Errorf("post not found")
	}

	objectID, err := store.GenerateNextID(ObjectTypeComment)
	if err != nil {
		return nil, fmt.Errorf("error generating object id: %w", err)
	}

	comment := Comment{
		ID:     objectID,
		PostID: posts[0].ID,
		UserID: viewer.UserID,
		Text:   req.Text,
		Date:   int(time.Now().Unix()),
	}

	err = store.Comment.Add(&comment)
	if err != nil {
		return nil, fmt.Errorf("error saving comment: %w", err)
	}

	return &pb.AddCommentResponse{}, nil
}

type Utils struct {
}

func (u *Utils) ResolveRoute(ctx context.Context, request *pb.ResolveRouteRequest) (*pb.UniversalRenderer, error) {
	type handler func(_ string, viewer *Viewer) (*pb.UniversalRenderer, error)

	routes := map[string]handler{
		`^/$`:      handleIndex,
		`^/login$`: handleLogin,

		`^/users/(\d+)$`: handleProfile,

		`^/posts/(\d+)$`: handlePostPage,

		`^/sandbox$`: handleSandbox,
	}

	viewer := GetViewerFromContext(ctx)

	for route, handler := range routes {
		matched, _ := regexp.MatchString(route, request.Url)
		if matched {
			return handler(request.Url, viewer)
		}
	}

	return &pb.UniversalRenderer{}, nil
}
