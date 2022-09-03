package types

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web4/store"
	"log"
	"strconv"
	"strings"
	"time"
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
	postIdsStr, err := redisClient.LRange(context.Background(), "feed", 0, 10).Result()
	if err != nil {
		log.Printf("ERR")
		return &Feed{}
	}

	postIds := make([]int, len(postIdsStr))
	for i, id := range postIdsStr {
		postIds[i], _ = strconv.Atoi(id)
	}

	userIds := map[string]bool{}
	posts := postsList(postIds)
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

	return &Feed{
		//Route: RouteFeed,
		Posts: postIdsStr,
		Nodes: &Nodes{
			Posts: postsList(postIds),
			Users: usersList(userIdsList),
		},
	}
}

func postPage(id int) *PostPage {
	userIds := map[string]bool{}
	posts := postsList([]int{id})
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

	return &PostPage{
		Route:    RoutePostsId,
		PagePost: strconv.Itoa(id),
		Nodes: &Nodes{
			Posts: posts,
			Users: usersList(userIdsList),
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

	var postIds []string
	for _, post := range user.LastPosts {
		postIds = append(postIds, strconv.Itoa(post))
	}

	return &UserPage{
		Route:    RouteUserPage,
		PageUser: strconv.Itoa(user.ID),
		Posts:    postIds,
		Nodes: &Nodes{
			Users: usersList([]int{id}),
			Posts: postsList(user.LastPosts),
		},
	}
}

type AddPostRequest struct {
	Text string `json:"text"`
}

type AddPostResult struct {
}

func addPost(req *AddPostRequest) *AddPostResult {
	postID := int(time.Now().Unix())

	post := store.Post{
		ID:     postID,
		UserID: 324825265,
		Text:   req.Text,
	}
	postBytes, _ := json.Marshal(post)

	err := redisClient.Set(context.Background(), fmt.Sprintf("node:%d", postID), postBytes, 0)
	log.Printf("%s", err)

	r2 := redisClient.LPush(context.Background(), "feed", post.ID)
	log.Printf("%s", r2)

	return &AddPostResult{}
}

type BrowseResult struct {
	Feed     *Feed          `json:"feed,omitempty"`
	UserPage *UserPage      `json:"userPage,omitempty"`
	PostPage *PostPage      `json:"postPage,omitempty"`
	AddPost  *AddPostResult `json:"addPost,omitempty"`
}

func Browse(url string) *BrowseResult {
	if url == "/" {
		return &BrowseResult{Feed: feedPage()}
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

	return &BrowseResult{}
}
