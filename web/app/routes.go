package app

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/materkov/meme9/web/pb"
	"github.com/materkov/meme9/web/store"
)

// /
func handleIndex(ctx context.Context, _ string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	if viewer.UserID == 0 {
		return &pb.UniversalRenderer{
			Renderer: &pb.UniversalRenderer_FeedRenderer{FeedRenderer: &pb.FeedRenderer{
				PlaceholderText: "Залогиньтесь, чтобы увидеть ленту.",
			}},
		}, nil
	}

	assocs, err := ObjectStore.AssocRange(ctx, viewer.UserID, store.Assoc_Following, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed getting following ids: %w", err)
	}

	followingIds := make([]int, len(assocs))
	for i, assoc := range assocs {
		followingIds[i] = assoc.Following.ID2
	}

	followingIds = append(followingIds, viewer.UserID)

	postIds := make([]int, 0)
	for _, userID := range followingIds {
		assocs, err := ObjectStore.AssocRange(ctx, userID, store.AssocPosted, 30)
		if err != nil {
			return nil, fmt.Errorf("error getting assocs: %w", err)
		}

		for _, assoc := range assocs {
			postIds = append(postIds, assoc.Posted.ID2)
		}
	}

	sort.Slice(postIds, func(i, j int) bool {
		return postIds[i] > postIds[j]
	})

	posts := make([]*store.Post, 0)
	for _, postId := range postIds {
		obj, err := ObjectStore.ObjGet(ctx, postId)
		if err != nil {
			log.Printf("error selcting post: %s", err)
			continue
		} else if obj == nil || obj.Post == nil {
			continue
		}

		posts = append(posts, obj.Post)
	}

	wrappedPosts := convertPosts(ctx, posts, viewer.UserID, true)

	return &pb.UniversalRenderer{
		Renderer: &pb.UniversalRenderer_FeedRenderer{FeedRenderer: &pb.FeedRenderer{
			Posts: wrappedPosts,
		}},
	}, nil
}

// /login
func handleLogin(ctx context.Context, _ string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	requestScheme := viewer.RequestScheme
	requestHost := viewer.RequestHost
	redirectURL := url.QueryEscape(fmt.Sprintf("%s://%s/vk-callback", requestScheme, requestHost))
	vkURL := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%d&response_type=code&redirect_uri=%s", DefaultConfig.VKAppID, redirectURL)

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_LoginPageRenderer{
		LoginPageRenderer: &pb.LoginPageRenderer{
			AuthUrl: vkURL,
			Text:    "Войти через ВК",
		},
	}}, nil
}

// /users/{id}
func handleProfile(ctx context.Context, url string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	req := &pb.ProfileGetRequest{
		Id: strings.TrimPrefix(url, "/users/"),
	}

	userID, _ := strconv.Atoi(req.Id)
	obj, err := ObjectStore.ObjGet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error selecting user: %w", err)
	} else if obj == nil || obj.User == nil {
		return nil, fmt.Errorf("user not found")
	}

	user := obj.User

	assocs, err := ObjectStore.AssocRange(ctx, userID, store.AssocPosted, 50)
	if err != nil {
		log.Printf("Error selecting user posts: %s", err)
	}

	posts := make([]*store.Post, 0)
	for _, assoc := range assocs {
		obj, err := ObjectStore.ObjGet(ctx, assoc.Posted.ID2)
		if err != nil {
			log.Printf("Error selecting post: %s", err)
			continue
		} else if obj == nil || obj.Post == nil {
			continue
		}

		posts = append(posts, obj.Post)
	}

	wrappedPosts := convertPosts(ctx, posts, viewer.UserID, false)

	assocs, err = ObjectStore.AssocRange(ctx, viewer.UserID, store.Assoc_Following, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed getting following ids: %w", err)
	}

	followingIds := make([]int, len(assocs))
	for i, assoc := range assocs {
		followingIds[i] = assoc.Following.ID2
	}

	isFollowing := false
	for _, userID := range followingIds {
		if userID == user.ID {
			isFollowing = true
		}
	}

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_ProfileRenderer{ProfileRenderer: &pb.ProfileRenderer{
		Id:          strconv.Itoa(user.ID),
		Name:        user.Name,
		Avatar:      user.VkAvatar,
		Posts:       wrappedPosts,
		IsFollowing: isFollowing,
	}}}, nil
}

// /posts/{id}
func handlePostPage(ctx context.Context, url string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	postIDStr := strings.TrimPrefix(url, "/posts/")
	postID, _ := strconv.Atoi(postIDStr)

	obj, err := ObjectStore.ObjGet(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("error selecting post: %s", err)
	} else if obj == nil || obj.Post == nil {
		return nil, fmt.Errorf("post not found")
	}

	assocs, err := ObjectStore.AssocRange(ctx, postID, store.Assoc_Commended, 100)
	if err != nil {
		log.Printf("Error selecting comment ids: %s", err)
	}

	commentIds := make([]int, len(assocs))
	for i, assoc := range assocs {
		commentIds[i] = assoc.Commented.ID2
	}

	var comments []*store.Comment
	for _, commentID := range commentIds {
		obj, err := ObjectStore.ObjGet(ctx, commentID)
		if err != nil || obj == nil || obj.Comment == nil {
			log.Printf("Error selecting comments objects: %s", err)
			continue
		}

		comments = append(comments, obj.Comment)
	}

	// TODO
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].ID > comments[j].ID
	})

	wrappedPosts := convertPosts(ctx, []*store.Post{obj.Post}, viewer.UserID, false)
	wrappedComments := convertComments(comments)

	composerPlaceholder := ""
	var composer *pb.CommentComposerRenderer

	if viewer.UserID != 0 {
		composer = &pb.CommentComposerRenderer{
			PostId:      postIDStr,
			Placeholder: "Напишите здесь свой комментарий...",
		}
	} else {
		composerPlaceholder = "Авторизуйтесь, чтрбы оставить комментарий."
	}

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_PostRenderer{PostRenderer: &pb.PostRenderer{
		Post:                wrappedPosts[0],
		Comments:            wrappedComments,
		Composer:            composer,
		ComposerPlaceholder: composerPlaceholder,
	}}}, nil
}

// /sandbox
func handleSandbox(ctx context.Context, url string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_SandboxRenderer{SandboxRenderer: &pb.SandboxRenderer{}}}, nil
}
