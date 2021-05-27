package main

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/materkov/meme9/web/pb"
)

// /
func handleIndex(_ string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	if viewer.UserID == 0 {
		return &pb.UniversalRenderer{
			Renderer: &pb.UniversalRenderer_FeedRenderer{FeedRenderer: &pb.FeedRenderer{
				PlaceholderText: "Залогиньтесь, чтобы увидеть ленту.",
			}},
		}, nil
	}

	followingIds, err := store.GetFollowing(viewer.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting following user ids: %w", err)
	}

	followingIds = append(followingIds, viewer.UserID)

	postIds, err := store.GetPostsByUsers(followingIds)
	if err != nil {
		return nil, fmt.Errorf("error getting post ids: %w", err)
	}

	posts, err := store.GetPosts(postIds)
	if err != nil {
		return nil, fmt.Errorf("error getting post ids: %w", err)
	}

	postsMap := map[int]*Post{}
	for _, post := range posts {
		postsMap[post.ID] = post
	}

	var postsOrdered []*Post
	for _, postId := range postIds {
		post := postsMap[postId]
		if post != nil {
			postsOrdered = append(postsOrdered, post)
		}
	}

	wrappedPosts := convertPosts(postsOrdered, viewer.UserID, true)

	return &pb.UniversalRenderer{
		Renderer: &pb.UniversalRenderer_FeedRenderer{FeedRenderer: &pb.FeedRenderer{
			Posts: wrappedPosts,
		}},
	}, nil
}

// /login
func handleLogin(_ string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	requestScheme := viewer.RequestScheme
	requestHost := viewer.RequestHost
	redirectURL := url.QueryEscape(fmt.Sprintf("%s://%s/vk-callback", requestScheme, requestHost))
	vkURL := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%d&response_type=code&redirect_uri=%s", config.VKAppID, redirectURL)

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_LoginPageRenderer{
		LoginPageRenderer: &pb.LoginPageRenderer{
			AuthUrl: vkURL,
			Text:    "Войти через ВК",
		},
	}}, nil
}

// /users/{id}
func handleProfile(url string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	req := &pb.ProfileGetRequest{
		Id: strings.TrimPrefix(url, "/users/"),
	}

	userID, _ := strconv.Atoi(req.Id)
	users, err := store.GetUsers([]int{userID})
	if err != nil {
		return nil, fmt.Errorf("error selecting user: %w", err)
	} else if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := users[0]

	postIds, err := store.GetPostsByUsers([]int{userID})
	if err != nil {
		log.Printf("Error selecting user posts: %s", err)
	}

	posts, err := store.GetPosts(postIds)
	if err != nil {
		log.Printf("Error selecting posts: %s", err)
	}

	wrappedPosts := convertPosts(posts, viewer.UserID, false)

	followingIds, err := store.GetFollowing(viewer.UserID)
	if err != nil {
		log.Printf("Error getting following users: %s", err)
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
func handlePostPage(url string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	postIDStr := strings.TrimPrefix(url, "/posts/")
	postID, _ := strconv.Atoi(postIDStr)

	posts, err := store.GetPosts([]int{postID})
	if err != nil {
		return nil, fmt.Errorf("error selecting post: %s", err)
	} else if len(posts) == 0 {
		return nil, fmt.Errorf("post not found")
	}

	commentIds, err := store.GetCommentsByPost(postID)
	if err != nil {
		log.Printf("Error selecting comment ids: %s", err)
	}

	comments, err := store.GetComments(commentIds)
	if err != nil {
		log.Printf("Error selecting comments objects: %s", err)
	}

	// TODO
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].ID > comments[j].ID
	})

	wrappedPosts := convertPosts(posts, viewer.UserID, false)
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
