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
	_ = pkg.AddEdge(post.UserID, postID, pkg.EdgeTypePosted, "")

	user, _ := pkg.GetUser(post.UserID)

	return transformPost(&post, user), nil
}

type PostsList struct {
}

func (h *HttpServer) PostsList(w http.ResponseWriter, r *http.Request, t *pkg.AuthToken) (interface{}, error) {
	req := PostsList{}
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

type PostsListById struct {
	ID string `json:"id"`
}

func (h *HttpServer) PostsListByID(w http.ResponseWriter, r *http.Request, t *pkg.AuthToken) (interface{}, error) {
	req := PostsListById{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, &Error{Code: 400, Message: "cannot parse request"}
	}

	postID, _ := strconv.Atoi(req.ID)

	post, err := pkg.GetPost(postID)
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	} else if post == nil {
		return nil, &Error{Code: 400, Message: "post not found"}
	}

	user, _ := pkg.GetUser(post.UserID)

	return transformPost(post, user), nil
}

type PostsListByUserReq struct {
	UserID string `json:"userId"`
}

func (h *HttpServer) PostsListByUser(w http.ResponseWriter, r *http.Request, t *pkg.AuthToken) (interface{}, error) {
	req := PostsListByUserReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, &Error{Code: 400, Message: "cannot parse request"}
	}

	userID, _ := strconv.Atoi(req.UserID)
	if userID <= 0 {
		return nil, &Error{Code: 400, Message: "incorrect user id"}
	}

	postIds, err := pkg.GetEdges(userID, pkg.EdgeTypePosted)
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

type PostsDeleteReq struct {
	ID string `json:"id"`
}

func (h *HttpServer) PostsDelete(w http.ResponseWriter, r *http.Request, t *pkg.AuthToken) (interface{}, error) {
	req := PostsDeleteReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, &Error{Code: 400, Message: "cannot parse request"}
	}

	postID, _ := strconv.Atoi(req.ID)

	post, err := pkg.GetPost(postID)
	if err != nil {
		return nil, err
	} else if post == nil {
		return nil, &Error{Code: 400, Message: "post not found"}
	}
	if t == nil {
		return nil, &Error{Code: 400, Message: "not authorized"}
	}
	if post.UserID != t.UserID {
		return nil, &Error{Code: 400, Message: "no access to this post"}
	}

	_ = pkg.DelEdge(post.UserID, pkg.EdgeTypePosted, post.ID)

	return Void{}, nil
}
