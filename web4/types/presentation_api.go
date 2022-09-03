package types

import (
	"github.com/materkov/meme9/web4/store"
	"log"
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

type Feed struct {
	Route Route    `json:"route,omitempty"`
	Posts []string `json:"posts,omitempty"`
	Nodes *Nodes   `json:"nodes,omitempty"`
}

type PostPage struct {
	Route    Route  `json:"route,omitempty"`
	PagePost string `json:"pagePost,omitempty"`
	Nodes    *Nodes `json:"nodes,omitempty"`
}

type UserPage struct {
	Route    Route    `json:"route,omitempty"`
	PageUser string   `json:"pageUser,omitempty"`
	NotFound bool     `json:"notFound,omitempty"`
	Nodes    *Nodes   `json:"nodes,omitempty"`
	Posts    []string `json:"posts,omitempty"`
}

type Nodes struct {
	Users []*User `json:"users,omitempty"`
	Posts []*Post `json:"posts,omitempty"`
}

func feedPage() *Feed {
	postIds, err := postsGetFeed()
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	posts := postsList(postIds)
	users := usersList(getUsersFromPosts(posts))

	postIdsStr := make([]string, len(postIds))
	for i, postID := range postIds {
		postIdsStr[i] = strconv.Itoa(postID)
	}

	return &Feed{
		//Route: RouteFeed,
		Posts: postIdsStr,
		Nodes: &Nodes{
			Posts: posts,
			Users: users,
		},
	}
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

func postPage(id int) *PostPage {
	posts := postsList([]int{id})
	users := usersList(getUsersFromPosts(posts))

	return &PostPage{
		Route:    RoutePostsId,
		PagePost: strconv.Itoa(id),
		Nodes: &Nodes{
			Posts: posts,
			Users: users,
		},
	}
}

func userPage(id int) *UserPage {
	result := &UserPage{}

	user := &store.User{}
	err := getObject(id, &user)
	if err != nil {
		result.NotFound = true
		return result
	}

	postIds, err := postsGetFeedByUsers([]int{user.ID})
	if err != nil {
		log.Printf("Error getting user feed: %s", err)
	}

	postIdsStr := make([]string, len(postIds))
	for i, postID := range postIds {
		postIdsStr[i] = strconv.Itoa(postID)
	}

	return &UserPage{
		Route:    RouteUserPage,
		PageUser: strconv.Itoa(user.ID),
		Posts:    postIdsStr,
		Nodes: &Nodes{
			Users: usersList([]int{id}),
			Posts: postsList(postIds),
		},
	}
}

type postsAddRequest struct {
	Text string `json:"text"`
}

type AddPostResult struct {
	Post *Post `json:"post,omitempty"`
}

func addPost(req *postsAddRequest, viewer *Viewer) *AddPostResult {
	postID, _ := postsAdd(req, viewer)
	post := postsList([]int{postID})[0]
	return &AddPostResult{Post: post}
}

type BrowseResult struct {
	UserID string `json:"userId,omitempty"`

	Feed       *Feed       `json:"feed,omitempty"`
	UserPage   *UserPage   `json:"userPage,omitempty"`
	PostPage   *PostPage   `json:"postPage,omitempty"`
	VkCallback *VkCallback `json:"vkCallback,omitempty"`
}

type VkCallback struct {
	UserID    string `json:"userId,omitempty"`
	AuthToken string `json:"authToken,omitempty"`
}

func Browse(url string, viewer *Viewer) *BrowseResult {
	if url == "/" {
		return &BrowseResult{
			UserID: strconv.Itoa(viewer.UserID),
			Feed:   feedPage(),
		}
	}

	if strings.HasPrefix(url, "/posts/") {
		postIDStr := strings.TrimPrefix(url, "/posts/")
		postID, _ := strconv.Atoi(postIDStr)
		return &BrowseResult{PostPage: postPage(postID)}
	}

	if strings.HasPrefix(url, "/users/") {
		postIDStr := strings.TrimPrefix(url, "/users/")
		postID, _ := strconv.Atoi(postIDStr)
		return &BrowseResult{UserPage: userPage(postID)}
	}

	if strings.HasPrefix(url, "/vk-callback") {
		code := strings.TrimPrefix(url, "/vk-callback?code=")
		vkID, _ := authExchangeCode("http://localhost:3000", code)
		if vkID != 0 {
			userID, _ := usersGetOrCreateByVKID(vkID)
			token, _ := authCreateToken(userID)

			return &BrowseResult{VkCallback: &VkCallback{
				UserID:    strconv.Itoa(userID),
				AuthToken: token,
			}}
		}
	}

	return &BrowseResult{}
}
