package api

import (
	"errors"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/store"
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

func transformPost(post *store.Post, user *store.User) *Post {
	return &Post{
		ID:     strconv.Itoa(post.ID),
		UserID: strconv.Itoa(post.UserID),
		Date:   time.Unix(int64(post.Date), 0).Format(time.RFC3339),
		Text:   post.Text,
		User:   transformUser(post.UserID, user),
	}
}

func (*API) PostsAdd(viewer *Viewer, r *PostsAddReq) (*Post, error) {
	if r.Text == "" {
		return nil, Error("TextEmpty")
	}
	if len(r.Text) > 5000 {
		return nil, Error("TextTooLong")
	}
	if viewer.UserID == 0 {
		return nil, Error("NotAuthorized")
	}

	post := store.Post{
		UserID: viewer.UserID,
		Date:   int(time.Now().Unix()),
		Text:   r.Text,
	}

	postID, err := store.AddObject(store.ObjTypePost, &post)
	if err != nil {
		return nil, fmt.Errorf("error saving post: %w", err)
	}
	post.ID = postID

	err = store.AddEdge(store.FakeObjPostedPost, postID, store.EdgeTypePostedPost, "")
	if err != nil {
		return nil, err
	}

	err = store.AddEdge(post.UserID, postID, store.EdgeTypePosted, "")
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(post.UserID)
	if err != nil {
		return nil, err
	}

	return transformPost(&post, user), nil
}

type PostsListReq struct {
}

func (h *API) PostsList(_ *Viewer, r *PostsListReq) ([]*Post, error) {
	postIds, err := store.GetEdges(store.FakeObjPostedPost, store.EdgeTypePostedPost)
	if err != nil {
		return nil, fmt.Errorf("error getting posted edges: %w", err)
	}

	result := make([]*Post, 0)
	for _, postID := range postIds {
		post, err := store.GetPost(postID)
		if err != nil {
			continue
		}

		user, _ := store.GetUser(post.UserID)

		result = append(result, transformPost(post, user))
	}

	return result, nil
}

type PostsListByIdReq struct {
	ID string `json:"id"`
}

func (h *API) PostsListByID(_ *Viewer, r *PostsListByIdReq) (*Post, error) {
	postID, _ := strconv.Atoi(r.ID)

	post, err := store.GetPost(postID)
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	} else if post == nil {
		return nil, Error("PostNotFound")
	}

	user, _ := store.GetUser(post.UserID)

	return transformPost(post, user), nil
}

type PostsListByUserReq struct {
	UserID string `json:"userId"`
}

func (h *API) PostsListByUser(_ *Viewer, r *PostsListByUserReq) ([]*Post, error) {
	userID, _ := strconv.Atoi(r.UserID)
	if userID <= 0 {
		return nil, Error("IncorrectUserId")
	}

	postIds, err := store.GetEdges(userID, store.EdgeTypePosted)
	if err != nil {
		return nil, fmt.Errorf("error getting posted edges: %w", err)
	}

	result := make([]*Post, 0)
	for _, postID := range postIds {
		post, err := store.GetPost(postID)
		if err != nil {
			continue
		}

		user, _ := store.GetUser(post.UserID)

		result = append(result, transformPost(post, user))
	}

	return result, nil
}

type PostsDeleteReq struct {
	PostID string `json:"postId"`
}

func (h *API) PostsDelete(viewer *Viewer, r *PostsDeleteReq) (interface{}, error) {
	postID, _ := strconv.Atoi(r.PostID)

	post, err := store.GetPost(postID)
	if errors.Is(err, store.ErrObjectNotFound) {
		return nil, Error("PostNotFound")
	} else if err != nil {
		return nil, err
	}

	if viewer.UserID == 0 {
		return nil, Error("NotAuthorized")
	}
	if post.UserID != viewer.UserID {
		return nil, Error("AccessDenied")
	}

	err = store.DelEdge(store.FakeObjPostedPost, store.EdgeTypePostedPost, post.ID)
	pkg.LogErr(err)

	err = store.DelEdge(post.UserID, store.EdgeTypePosted, post.ID)
	pkg.LogErr(err)

	return Void{}, nil
}
