package api

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/pkg"
	"net/http"
	"strconv"
	"time"
)

type Post struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Date   string `json:"date"`
	Text   string `json:"text"`
	User   *User  `json:"user"`
}

type PostsAddReq struct {
	Text string `json:"text"`
}

func transformPost(post *pkg.Post, user *pkg.User) *Post {
	return &Post{
		ID:     strconv.Itoa(post.ID),
		UserID: strconv.Itoa(post.UserID),
		Date:   time.Unix(int64(post.Date), 0).Format(time.RFC3339),
		Text:   post.Text,
		User:   transformUser(post.UserID, user),
	}
}

func (h *HttpServer) PostsAdd(w http.ResponseWriter, r *http.Request, t *pkg.AuthToken) (interface{}, error) {
	req := PostsAddReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, &Error{Code: 400, Message: "cannot parse request"}
	}
	if req.Text == "" {
		return nil, &Error{Code: 400, Message: "empty text"}
	}
	if t == nil {
		return nil, &Error{Code: 400, Message: "not authorized"}
	}

	post := pkg.Post{
		UserID: t.UserID,
		Date:   int(time.Now().Unix()),
		Text:   req.Text,
	}

	postID, err := pkg.AddObject(pkg.ObjTypePost, &post)
	if err != nil {
		return nil, fmt.Errorf("error saving post: %w", err)
	}
	post.ID = postID

	_ = pkg.AddEdge(pkg.FakeObjPostedPost, postID, pkg.EdgeTypePostedPost, "")

	user, _ := pkg.GetUser(post.UserID)

	return transformPost(&post, user), nil
}

type PostsListReq struct {
	Text string `json:"text"`
}

func (h *HttpServer) PostsList(w http.ResponseWriter, r *http.Request, t *pkg.AuthToken) (interface{}, error) {
	req := PostsAddReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, &Error{Code: 400, Message: "cannot parse request"}
	}

	postIds, err := pkg.GetEdges(pkg.FakeObjPostedPost, pkg.EdgeTypePostedPost)
	if err != nil {
		return nil, fmt.Errorf("error getting posted edges: %w", err)
	}

	result := make([]*Post, 0)
	for _, postID := range postIds {
		post, err := pkg.GetPost(postID)
		if err != nil {
			continue
		}

		user, _ := pkg.GetUser(post.UserID)

		result = append(result, transformPost(post, user))
	}

	return result, nil
}
