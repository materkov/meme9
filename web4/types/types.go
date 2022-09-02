package types

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web4/store"
	"net/http"
	"strconv"
	"strings"
)

type Composer struct {
}

type Feed struct {
	Posts []string `json:"posts,omitempty"`
	Nodes *Nodes   `json:"nodes,omitempty"`
}

type PostPage struct {
	PagePost string `json:"pagePost,omitempty"`
	Nodes    *Nodes `json:"nodes,omitempty"`
}

type UserPage struct {
	PageUser string   `json:"pageUser,omitempty"`
	NotFound bool     `json:"notFound,omitempty"`
	Nodes    *Nodes   `json:"nodes,omitempty"`
	Posts    []string `json:"posts,omitempty"`
}

type User struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Href string `json:"href,omitempty"`
}

type Nodes struct {
	Users []*User `json:"users,omitempty"`
	Posts []*Post `json:"posts,omitempty"`
}

func wrapUsers(userIds []int) []*User {
	users := store.DefaultStore.GetUsers(userIds)

	result := make([]*User, len(userIds))
	for i, user := range users {
		result[i] = &User{
			ID:   strconv.Itoa(user.ID),
			Name: user.Name,
			Href: fmt.Sprintf("/users/%d", user.ID),
		}
	}

	return result
}

func wrapPostsList(postIds []int) []*Post {
	posts := store.DefaultStore.GetPosts(postIds)

	postsMap := map[int]store.Post{}
	for _, post := range posts {
		postsMap[post.ID] = post
	}

	userIds := map[int]bool{}
	for _, post := range posts {
		userIds[post.UserID] = true
	}

	userIdsList := make([]int, 0, len(userIds))
	for userID := range userIds {
		userIdsList = append(userIdsList, userID)
	}

	users := store.DefaultStore.GetUsers(userIdsList)

	usersMap := map[int]store.User{}
	for _, user := range users {
		usersMap[user.ID] = user
	}

	var results []*Post
	for _, id := range postIds {
		post := postsMap[id]
		//if !ok {
		//	continue
		//}

		user, isUserOk := usersMap[post.UserID]

		wrappedPost := &Post{
			ID:         strconv.Itoa(post.ID),
			FromHref:   fmt.Sprintf("/users/%d", post.UserID),
			Text:       post.Text,
			DetailsURL: fmt.Sprintf("/posts/%d", post.ID),
		}

		if isUserOk {
			wrappedPost.FromName = user.Name
		}

		results = append(results, wrappedPost)
	}

	return results
}

func feedPage() *Feed {
	postIds := []int{101, 100}
	return &Feed{
		Posts: []string{
			"101", "100",
		},
		Nodes: &Nodes{
			Posts: wrapPostsList(postIds),
		},
	}
}

func postPage(id int) *PostPage {
	return &PostPage{
		PagePost: strconv.Itoa(id),
		Nodes: &Nodes{
			Posts: wrapPostsList([]int{id}),
		},
	}
}

type Post struct {
	ID       string `json:"id,omitempty"`
	FromHref string `json:"fromHref"`
	FromName string `json:"fromName"`

	Text       string `json:"text"`
	DetailsURL string `json:"detailsURL"`
}

func userPage(id int) *UserPage {
	result := &UserPage{}

	users := store.DefaultStore.GetUsers([]int{id})
	if len(users) == 0 {
		result.NotFound = true
		return result
	}

	var postIds []string
	for _, post := range users[0].LastPosts {
		postIds = append(postIds, strconv.Itoa(post))
	}

	return &UserPage{
		PageUser: strconv.Itoa(users[0].ID),
		Posts:    postIds,
		Nodes: &Nodes{
			Users: wrapUsers([]int{id}),
			Posts: wrapPostsList(users[0].LastPosts),
		},
	}
}

type BrowseResult struct {
	Feed     *Feed     `json:"feed,omitempty"`
	UserPage *UserPage `json:"userPage,omitempty"`
	PostPage *PostPage `json:"postPage,omitempty"`
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

func DoHandle() {
	// CRUD words: insert, delete, update, list

	http.HandleFunc("/browse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		resp := Browse(r.URL.Query().Get("url"))
		_ = json.NewEncoder(w).Encode(resp)
	})
	http.HandleFunc("/posts.insert", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		resp := Browse(r.URL.Query().Get("url"))
		_ = json.NewEncoder(w).Encode(resp)
	})

	http.ListenAndServe("127.0.0.1:8000", nil)
}
