package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/pkg/tracer"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
	"math"
	"net/url"
	"slices"
	"strconv"
	"sync"
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
	Poll *Poll     `json:"poll,omitempty"`
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

	PollID string `json:"pollId"`
}

func transformPostBatch(ctx context.Context, posts []*store.Post, viewerID int) []*Post {
	defer tracer.FromCtx(ctx).StartChild("transformPostBatch").Stop()

	var userIds []int
	var postIds []int
	var pollIds []int
	for _, post := range posts {
		userIds = append(userIds, post.UserID)
		postIds = append(postIds, post.ID)

		if post.PollID != 0 {
			pollIds = append(pollIds, post.PollID)
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	var counters map[int]int
	var isLiked map[int]bool
	go func() {
		var err error
		counters, isLiked, err = store.GlobalStore.LoadLikesMany(ctx, postIds, viewerID)
		pkg.LogErr(err)
		wg.Done()
	}()

	usersWrappedMap := map[int]*User{}
	go func() {
		usersMap, err := store.GlobalStore.GetObjectsMany(ctx, utils.UniqueIds(userIds))
		pkg.LogErr(err)

		for userID, userBytes := range usersMap {
			user := &store.User{}
			err = json.Unmarshal(userBytes, user)
			pkg.LogErr(err)
			if err != nil {
				user = nil
			}

			usersWrappedMap[userID], err = transformUser(userID, user, viewerID)
			pkg.LogErr(err)
		}
		wg.Done()
	}()

	pollWrappedMap := map[int]*Poll{}
	go func() {
		pollBytes, err := store.GlobalStore.GetObjectsMany(ctx, pollIds)
		pkg.LogErr(err)

		polls := make([]*store.Poll, len(pollIds))
		for i, pollID := range pollIds {
			pollData := pollBytes[pollID]

			poll := store.Poll{}
			err = json.Unmarshal(pollData, &poll)
			pkg.LogErr(err)
			poll.ID = pollID

			polls[i] = &poll
		}

		pollsWrapped := transformPollsMany(ctx, polls, viewerID)

		for i, poll := range pollsWrapped {
			pollWrappedMap[pollIds[i]] = poll
		}

		wg.Done()
	}()
	wg.Wait()

	result := make([]*Post, len(posts))

	for i, post := range posts {
		var wrappedLink *PostLink
		if post.Link != nil {
			host := ""

			parsedURL, err := url.Parse(post.Link.URL)
			if err == nil {
				host = parsedURL.Host
			}

			proxiedUrl := fmt.Sprintf("https://3c6ef5be-e5f9-4e47-9a68-bd635323a374.selcdn.net/image-proxy?url=%s", url.QueryEscape(post.Link.ImageURL))

			wrappedLink = &PostLink{
				URL:         post.Link.URL,
				Title:       post.Link.Title,
				Description: post.Link.Description,
				ImageURL:    proxiedUrl,
				Domain:      host,
			}
		}

		pollTransformed := pollWrappedMap[post.PollID]

		result[i] = &Post{
			ID:     strconv.Itoa(post.ID),
			UserID: strconv.Itoa(post.UserID),
			Date:   time.Unix(int64(post.Date), 0).Format(time.RFC3339),
			Text:   post.Text,
			User:   usersWrappedMap[post.UserID],

			LikesCount: counters[post.ID],
			IsLiked:    isLiked[post.ID],

			Link: wrappedLink,
			Poll: pollTransformed,
		}
	}

	return result
}

func (*API) PostsAdd(ctx context.Context, viewer *Viewer, r *PostsAddReq) (*Post, error) {
	if r.Text == "" {
		return nil, Error("TextEmpty")
	}
	if len(r.Text) > 5000 {
		return nil, Error("TextTooLong")
	}
	if viewer.UserID == 0 {
		return nil, Error("NotAuthorized")
	}

	pollID := 0
	if r.PollID != "" {
		pollID, _ = strconv.Atoi(r.PollID)
	}

	post := store.Post{
		UserID: viewer.UserID,
		Date:   int(time.Now().Unix()),
		Text:   r.Text,
		PollID: pollID,
	}

	postID, err := store2.GlobalStore.Nodes.Add(store.ObjTypePost, &post)
	if err != nil {
		return nil, fmt.Errorf("error saving post: %w", err)
	}
	post.ID = postID

	err = store.GlobalStore.AddEdge(store.FakeObjPostedPost, postID, store.EdgeTypePostedPost)
	if err != nil {
		return nil, err
	}

	err = store.GlobalStore.AddEdge(post.UserID, postID, store.EdgeTypePosted)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pkg.TryParseLink(&post)
		pkg.LogErr(err)
	}()

	return transformPostBatch(ctx, []*store.Post{&post}, viewer.UserID)[0], nil
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

func (a *API) PostsList(ctx context.Context, v *Viewer, r *PostsListReq) (*PostsList, error) {
	defer tracer.FromCtx(ctx).StartChild("API.PostsList").Stop()

	var err error
	var postIds []int

	if r.Type == "FEED" {
		postIds, err = pkg.GetFeedPostIds(v.UserID)
	} else if r.Type == "DISCOVER" || r.Type == "" {
		var edges []store.Edge
		edges, err = store.GlobalStore.GetEdges(store.FakeObjPostedPost, store.EdgeTypePostedPost, 1000, math.MaxInt)
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

	postsMap, err := store.GlobalStore.GetObjectsMany(ctx, postIds)
	if err != nil {
		return nil, err
	}

	var posts []*store.Post
	for _, postID := range postIds {
		postBytes := postsMap[postID]
		if postBytes == nil {
			continue
		}

		post := store.Post{}
		err = json.Unmarshal(postBytes, &post)
		if err != nil {
			pkg.LogErr(err)
		}
		post.ID = postID // TODO think about this

		posts = append(posts, &post)
	}

	result.Items = transformPostBatch(ctx, posts, v.UserID)

	return result, nil
}

type PostsListByIdReq struct {
	ID string `json:"id"`
}

func (a *API) PostsListByID(ctx context.Context, v *Viewer, r *PostsListByIdReq) (*Post, error) {
	postID, _ := strconv.Atoi(r.ID)

	post, err := store.GetPost(postID)
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	} else if post == nil || post.IsDeleted {
		return nil, Error("PostNotFound")
	}

	return transformPostBatch(ctx, []*store.Post{post}, v.UserID)[0], nil
}

type PostsListByUserReq struct {
	UserID string `json:"userId"`

	Count int    `json:"count"`
	After string `json:"after"`
}

func (a *API) PostsListByUser(ctx context.Context, v *Viewer, r *PostsListByUserReq) (*PostsList, error) {
	userID, _ := strconv.Atoi(r.UserID)
	if userID <= 0 {
		return nil, Error("IncorrectUserId")
	}

	if r.Count < 0 {
		return nil, Error("IncorrectCount")
	} else if r.Count > 50 {
		return nil, Error("IncorrectCount")
	}

	lessThan := math.MaxInt
	if r.After != "" {
		lessThan, _ = strconv.Atoi(r.After)
	}

	count := 10
	if r.Count != 0 {
		count = r.Count
	}

	edges, err := store.GlobalStore.GetEdges(userID, store.EdgeTypePosted, count, lessThan)
	if err != nil {
		return nil, fmt.Errorf("error getting posted edges: %w", err)
	}

	var posts []*store.Post

	for _, edge := range edges {
		post, err := store.GetPost(edge.ToID)
		if err == nil {
			posts = append(posts, post)
			continue
		}
	}

	result := transformPostBatch(ctx, posts, v.UserID)

	nextAfter := ""
	if len(edges) == count && len(edges) > 0 {
		nextAfter = strconv.Itoa(edges[len(edges)-1].ToID)
	}

	return &PostsList{
		Items:     result,
		PageToken: nextAfter,
	}, nil
}

type PostsDeleteReq struct {
	PostID string `json:"postId"`
}

func (a *API) PostsDelete(viewer *Viewer, r *PostsDeleteReq) (*Void, error) {
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

	post.IsDeleted = true

	err = store2.GlobalStore.Nodes.Update(post.ID, post)
	pkg.LogErr(err)

	err = store.GlobalStore.DelEdge(store.FakeObjPostedPost, post.ID, store.EdgeTypePostedPost)
	pkg.LogErr(err)

	err = store.GlobalStore.DelEdge(post.UserID, post.ID, store.EdgeTypePosted)
	pkg.LogErr(err)

	return &Void{}, nil
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
		err = store.GlobalStore.DelEdge(postID, v.UserID, store.EdgeTypeLiked)
		if err != nil {
			return nil, err
		}
	} else {
		err = store.GlobalStore.AddEdge(postID, v.UserID, store.EdgeTypeLiked)
		if errors.Is(err, store.ErrDuplicateEdge) {
			// Do nothing
		} else if err != nil {
			return nil, err
		}
	}

	return &Void{}, nil
}
