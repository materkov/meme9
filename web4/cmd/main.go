package main

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web4/fields"
	"github.com/materkov/meme9/web4/store"
	"github.com/materkov/meme9/web4/types"
	"log"
	"net/http"
	"strconv"
)

func logApiCall(method string, args interface{}) {
	argsBytes, _ := json.Marshal(args)
	log.Printf("API %s %s", method, argsBytes)
}

type Post struct {
	ID       string `json:"id,omitempty"`
	TypeName string `json:"typeName,omitempty"`
	Text     string `json:"text,omitempty"`
	Date     int    `json:"date,omitempty"`
	User     *User  `json:"user,omitempty"`
}

type User struct {
	ID       string  `json:"id,omitempty"`
	Typename string  `json:"typeName,omitempty"`
	Name     string  `json:"name,omitempty"`
	Posts    []*Post `json:"posts,omitempty"`
}

type feedListRequest struct {
	Fields string
}

type feedListResponse struct {
	Items []*Post `json:"items,omitempty"`
}

func feedList(req *feedListRequest) *feedListResponse {
	feed := []string{"101", "100"}
	logApiCall("feedList", req)

	f := fields.ParseFields(req.Fields)

	var items []*Post
	if _, itemsFields := f.Has("items"); itemsFields != nil {
		items = postsList(&postsListRequest{
			Id:     feed,
			Fields: itemsFields.ToString(),
		}).Items
		return &feedListResponse{
			Items: items,
		}
	} else {
		items = make([]*Post, len(feed))
		for i, postId := range feed {
			items[i] = &Post{ID: postId}
		}
	}

	return &feedListResponse{
		Items: items,
	}
}

type postsListRequest struct {
	Id     []string
	Fields string
}

type postsListResponse struct {
	Items []*Post `json:"items"`
}

func postsList(req *postsListRequest) *postsListResponse {
	logApiCall("postsList", req)
	f := fields.ParseFields(req.Fields)

	var neededUsers []string
	_, userFields := f.Has("user")

	result := make([]*Post, len(req.Id))
	for i, postID := range req.Id {
		post := &Post{
			ID: postID,
		}

		if _, userFields := f.Has("text"); userFields != nil {
			post.Text = fmt.Sprintf("Post %s", postID)
		}

		if userFields != nil {
			neededUsers = append(neededUsers, "10")
		}
		result[i] = post
	}

	if len(neededUsers) > 0 {
		allUsers := usersList(&usersListRequest{
			Id:     neededUsers,
			Fields: userFields.ToString(),
		}).Items

		for _, post := range result {
			for _, user := range allUsers {
				if user.ID == "10" {
					post.User = user
				}
			}
		}
	}

	return &postsListResponse{Items: result}
}

type usersListRequest struct {
	Id     []string
	Fields string
}

type usersListResponse struct {
	Items []*User `json:"items,omitempty"`
}

func usersList(req *usersListRequest) *usersListResponse {
	logApiCall("usersList", req)
	f := fields.ParseFields(req.Fields)

	ids := make([]int, 0)
	for _, idStr := range req.Id {
		id, _ := strconv.Atoi(idStr)
		ids = append(ids, id)
	}

	users := store.DefaultStore.GetUsers(ids)

	usersMap := map[string]store.User{}
	for _, user := range users {
		usersMap[strconv.Itoa(user.ID)] = user
	}

	neededPosts := make([]string, 0)
	_, postFields := f.Has("posts")

	result := make([]*User, len(req.Id))
	for i, ID := range req.Id {
		user := &User{ID: ID}

		if _, nameFields := f.Has("name"); nameFields != nil {
			user.Name = usersMap[ID].Name
		}

		if postFields != nil {
			neededPosts = append(neededPosts, []string{"100", "101"}...)
		}

		result[i] = user
	}

	if len(neededPosts) > 0 {
		allPosts := postsList(&postsListRequest{
			Id:     neededPosts,
			Fields: postFields.ToString(),
		})

		for _, user := range result {
			for _, post := range allPosts.Items {
				if post.ID == "101" || post.ID == "100" {
					user.Posts = append(user.Posts, post)
				}
			}
		}
	}

	return &usersListResponse{Items: result}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var operations []Operation
	_ = json.NewDecoder(r.Body).Decode(&operations)

	results := make([]interface{}, 0)

	for _, operation := range operations {
		method := operation.Method
		params := operation.Params

		var resp interface{}

		switch method {
		case "feed.list":
			req := feedListRequest{}
			_ = json.Unmarshal(params, &req)
			resp = feedList(&req)

		case "posts.list":
			req := postsListRequest{}
			_ = json.Unmarshal(params, &req)
			resp = postsList(&req)

		case "users.list":
			req := usersListRequest{}
			_ = json.Unmarshal(params, &req)
			resp = usersList(&req)

		case "posts.add":
			resp = Post{
				ID:   "post100",
				Text: "text post 1",
				Date: 1000000,
				//UserID: "user10",
				User: &User{
					ID:   "user10",
					Name: "Name for user 10",
				},
			}
		}

		results = append(results, resp)
	}

	_ = json.NewEncoder(w).Encode(results)
}

type Operation struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

func main() {
	types.DoHandle()
	//http.HandleFunc("/api", apiHandler)
	//_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
