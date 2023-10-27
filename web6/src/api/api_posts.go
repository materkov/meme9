package api

import (
	"errors"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/store"
	"net/url"
	"slices"
	"strconv"
	"time"
)

type Post struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Date   string `json:"date"`
	Text   string `json:"text"`
	User   *User  `json:"user"`

	IsLiked    bool `json:"isLiked,omitempty"`
	LikesCount int  `json:"likesCount,omitempty"`

	Link *PostLink `json:"link,omitempty"`
}

type PostLink struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	Domain      string `json:"domain"`
}

type PostsList struct {
	Items     []*Post `json:"items,omitempty"`
	PageToken string  `json:"pageToken,omitempty"`
}

type PostsAddReq struct {
	Text string `json:"text"`
}

func transformPost(post *store.Post, user *store.User, viewerID int) *Post {
	userWrapped, err := transformUser(post.UserID, user, viewerID)
	pkg.LogErr(err) // TODO think about it

	likesCount, err := store.CountEdges(post.ID, store.EdgeTypeLiked)
	pkg.LogErr(err)

	edge, err := store.GetEdge(post.ID, viewerID, store.EdgeTypeLiked)
	if err != nil && !errors.Is(err, store.ErrNoEdge) {
		pkg.LogErr(err)
	}

	var wrappedLink *PostLink
	if post.Link != nil {
		host := ""

		parsedURL, err := url.Parse(post.Link.URL)
		if err == nil {
			host = parsedURL.Host
		}

		wrappedLink = &PostLink{
			URL:         post.Link.URL,
			Title:       post.Link.Title,
			Description: post.Link.Description,
			ImageURL:    post.Link.ImageURL,
			Domain:      host,
		}
	}

	return &Post{
		ID:     strconv.Itoa(post.ID),
		UserID: strconv.Itoa(post.UserID),
		Date:   time.Unix(int64(post.Date), 0).Format(time.RFC3339),
		Text:   post.Text,
		User:   userWrapped,

		LikesCount: likesCount,
		IsLiked:    edge != nil,

		Link: wrappedLink,
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

	err = store.AddEdge(store.FakeObjPostedPost, postID, store.EdgeTypePostedPost)
	if err != nil {
		return nil, err
	}

	err = store.AddEdge(post.UserID, postID, store.EdgeTypePosted)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(post.UserID)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pkg.TryParseLink(&post)
		pkg.LogErr(err)
	}()

	return transformPost(&post, user, viewer.UserID), nil
}

type FeedType string

const (
	Feed     FeedType = "FEED"
	Discover FeedType = "DISCOVER"
)

type PostsListReq struct {
	Type FeedType `json:"type"`

	Count     int    `json:"count"`
	PageToken string `json:"pageToken"`
}

func (h *API) PostsList(v *Viewer, r *PostsListReq) (*PostsList, error) {
	var err error
	var postIds []int

	if r.Type == "FEED" {
		postIds, err = pkg.GetFeedPostIds(v.UserID)
	} else if r.Type == "DISCOVER" || r.Type == "" {
		var edges []store.Edge
		edges, err = store.GetEdges(store.FakeObjPostedPost, store.EdgeTypePostedPost)
		postIds = store.GetToId(edges)
	} else {
		err = Error("InvalidFeedType")
	}

	if err != nil {
		return nil, err
	}

	if r.PageToken != "" {
		cursor, _ := strconv.Atoi(r.PageToken)
		cursorIdx := slices.Index(postIds, cursor)
		if cursorIdx != -1 && cursorIdx != len(postIds)-1 {
			postIds = postIds[cursorIdx+1:]
		}
	}

	result := &PostsList{}

	count := r.Count
	if count == 0 {
		count = 10
	}

	if len(postIds) > count {
		result.PageToken = strconv.Itoa(postIds[count-1])
		postIds = postIds[:count]
	}

	for _, postID := range postIds {
		post, err := store.GetPost(postID)
		if err != nil {
			continue
		}

		user, _ := store.GetUser(post.UserID)

		result.Items = append(result.Items, transformPost(post, user, v.UserID))
	}

	return result, nil
}

type PostsListByIdReq struct {
	ID string `json:"id"`
}

func (h *API) PostsListByID(v *Viewer, r *PostsListByIdReq) (*Post, error) {
	postID, _ := strconv.Atoi(r.ID)

	post, err := store.GetPost(postID)
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	} else if post == nil {
		return nil, Error("PostNotFound")
	}

	user, _ := store.GetUser(post.UserID)

	return transformPost(post, user, v.UserID), nil
}

type PostsListByUserReq struct {
	UserID string `json:"userId"`
}

func (h *API) PostsListByUser(v *Viewer, r *PostsListByUserReq) ([]*Post, error) {
	userID, _ := strconv.Atoi(r.UserID)
	if userID <= 0 {
		return nil, Error("IncorrectUserId")
	}

	edges, err := store.GetEdges(userID, store.EdgeTypePosted)
	if err != nil {
		return nil, fmt.Errorf("error getting posted edges: %w", err)
	}

	result := make([]*Post, 0)
	for _, edge := range edges {
		post, err := store.GetPost(edge.ToID)
		if err != nil {
			continue
		}

		user, _ := store.GetUser(post.UserID)

		result = append(result, transformPost(post, user, v.UserID))
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

	err = store.DelEdge(store.FakeObjPostedPost, post.ID, store.EdgeTypePostedPost)
	pkg.LogErr(err)

	err = store.DelEdge(post.UserID, post.ID, store.EdgeTypePosted)
	pkg.LogErr(err)

	return Void{}, nil
}

type PostLikeAction string

const (
	Like   PostLikeAction = "LIKE"
	Unlike PostLikeAction = "UNLIKE"
)

type PostsLikeReq struct {
	PostID string         `json:"postId"`
	Action PostLikeAction `json:"action"`
}

func (*API) PostsLike(v *Viewer, r *PostsLikeReq) (*Void, error) {
	if v.UserID == 0 {
		return nil, Error("NotAuthorized")
	}

	postID, _ := strconv.Atoi(r.PostID)
	_, err := store.GetPost(postID)
	if errors.Is(err, store.ErrObjectNotFound) {
		return nil, Error("PostNotFound")
	} else if err != nil {
		return nil, err
	}

	if r.Action == Unlike {
		err = store.DelEdge(postID, v.UserID, store.EdgeTypeLiked)
		if err != nil {
			return nil, err
		}
	} else {
		err = store.AddEdge(postID, v.UserID, store.EdgeTypeLiked)
		if errors.Is(err, store.ErrDuplicateEdge) {
			// Do nothing
		} else if err != nil {
			return nil, err
		}
	}

	return &Void{}, nil
}
