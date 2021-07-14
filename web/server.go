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
	store2 "github.com/materkov/meme9/web/store"
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

	assocs, err := objectStore.AssocRange(viewer.UserID, store2.Assoc_Following, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed getting following ids: %w", err)
	}

	followingIds := make([]int, len(assocs))
	for i, assoc := range assocs {
		followingIds[i] = assoc.Following.ID2
	}

	for _, userID := range followingIds {
		if requestedID == userID {
			return &pb.RelationsFollowResponse{}, nil
		}
	}

	err = objectStore.AssocAdd(viewer.UserID, requestedID, store2.Assoc_Following, &store2.StoredAssoc{Following: &store2.Following{
		ID1:  viewer.UserID,
		ID2:  requestedID,
		Type: store2.Assoc_Liked,
	}})
	if err != nil {
		return nil, fmt.Errorf("failed saving assoc: %w", err)
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

	err := objectStore.AssocDelete(viewer.UserID, requestedID, store2.Assoc_Following)
	if err != nil {
		return nil, fmt.Errorf("failed deleting assoc: %w", err)
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
		obj, err := objectStore.ObjGet(photoID)
		if obj == nil || obj.Photo == nil {
			return nil, fmt.Errorf("photo not found")
		} else if err != nil {
			return nil, fmt.Errorf("error getting photos: %w", err)
		}

		if obj.Photo.UserID != viewer.UserID {
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

	data, err := objectStore.AssocGet(postID, store2.Assoc_Liked, viewer.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting is liked: %w", err)
	}

	isLiked := data != nil && data.Liked != nil

	if req.Action == pb.ToggleLikeRequest_LIKE && !isLiked {
		err = objectStore.AssocAdd(postID, viewer.UserID, store2.Assoc_Liked, &store2.StoredAssoc{Liked: &store2.Liked{
			ID1:  postID,
			ID2:  viewer.UserID,
			Type: store2.Assoc_Liked,
		}})
		if err != nil {
			return nil, fmt.Errorf("error saving like: %w", err)
		}
	}

	if req.Action == pb.ToggleLikeRequest_UNLIKE && isLiked {
		err = objectStore.AssocDelete(postID, viewer.UserID, store2.Assoc_Liked)
		if err != nil {
			return nil, fmt.Errorf("error deleting like: %w", err)
		}
	}

	likesCount, err := objectStore.AssocCount(postID, store2.Assoc_Liked)
	if err != nil {
		log.Printf("Error getting likes count")
	}

	return &pb.ToggleLikeResponse{LikesCount: int32(likesCount)}, nil
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

	comment := store2.Comment{
		ID:     objectID,
		PostID: posts[0].ID,
		UserID: viewer.UserID,
		Text:   req.Text,
		Date:   int(time.Now().Unix()),
	}

	err = objectStore.ObjAdd(&store2.StoredObject{ID: objectID, Comment: &comment})
	if err != nil {
		return nil, fmt.Errorf("error saving comment: %w", err)
	}

	err = objectStore.AssocAdd(comment.PostID, comment.ID, store2.Assoc_Commended, &store2.StoredAssoc{Commented: &store2.Commented{
		ID1:  comment.PostID,
		ID2:  comment.ID,
		Type: store2.Assoc_Commended,
	}})
	if err != nil {
		return nil, fmt.Errorf("error saving assoc: %w", err)
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
