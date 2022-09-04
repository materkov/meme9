package types

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Composer struct {
}

type Route string

const (
	RoutePostsId  Route = "PostPage"
	RouteFeed     Route = "Feed"
	RouteUserPage Route = "UserPage"
)

type Nodes struct {
	Users []*User `json:"users,omitempty"`
	Posts []*Post `json:"posts,omitempty"`
}

func getUsersFromPosts(posts []*Post) []int {
	userIds := map[string]bool{}
	for _, post := range posts {
		userIds[post.FromID] = true
	}

	userIdsList := make([]int, 0)
	for userIdStr := range userIds {
		userID, _ := strconv.Atoi(userIdStr)
		if userID > 0 {
			userIdsList = append(userIdsList, userID)
		}
	}

	return userIdsList
}

type BrowseResult struct {
	UserID string `json:"userId,omitempty"`

	VkCallback *VkCallbackResponse `json:"vkCallback,omitempty"`
	AddPost    *AddPostResponse    `json:"addPost,omitempty"`

	ComponentName string      `json:"componentName"`
	ComponentData interface{} `json:"componentData"`
}

func Browse(url string, q string, viewer *Viewer) *BrowseResult {
	result := &BrowseResult{
		UserID: strconv.Itoa(viewer.UserID),
	}

	if url == "/" {
		result.ComponentName = "Feed"
		result.ComponentData, _ = Feed(&FeedRequest{})
	}

	if strings.HasPrefix(url, "/posts/") {
		postIDStr := strings.TrimPrefix(url, "/posts/")
		result.ComponentName = "PostPage"
		result.ComponentData, _ = PostPage(&PostPageRequest{PostID: postIDStr})
	}

	if strings.HasPrefix(url, "/users/") {
		userIDStr := strings.TrimPrefix(url, "/users/")
		result.ComponentName = "UserPage"
		result.ComponentData, _ = UserPage(&UserPageRequest{UserID: userIDStr})
	}

	if strings.HasPrefix(url, "/vk-callback") {
		code := strings.TrimPrefix(url, "/vk-callback?code=")
		result.VkCallback, _ = VkCallback(&VkCallbackRequest{Code: code})
	}

	if strings.HasPrefix(url, "/posts/add") {
		req := &AddPostRequest{}
		_ = json.Unmarshal([]byte(q), req)
		result.AddPost, _ = AddPost(req, viewer)
	}

	return result
}

// Queries: PostPage, UserPage, Feed, VkCallback, AddPost
// Routing: с бека надо ответить, какой компонент показать и какие пропсы на нем
