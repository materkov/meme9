package api

import (
	"context"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/pkg/tracer"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
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
		counters, isLiked, err = store2.GlobalStore.Likes.Get(ctx, postIds, viewerID)
		pkg.LogErr(err)
		wg.Done()
	}()

	usersWrappedMap := map[int]*User{}
	go func() {
		usersMap, err := store2.GlobalStore.Users.Get(utils.UniqueIds(userIds))
		pkg.LogErr(err)

		for userID, user := range usersMap {
			usersWrappedMap[userID], err = transformUser(userID, user, viewerID)
			pkg.LogErr(err)
		}
		wg.Done()
	}()

	pollWrappedMap := map[int]*Poll{}
	go func() {
		pollBytes, err := store2.GlobalStore.Polls.Get(pollIds)
		pkg.LogErr(err)

		polls := make([]*store.Poll, len(pollIds))
		for i, pollID := range pollIds {
			polls[i] = pollBytes[pollID]
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

	err := store2.GlobalStore.Posts.Add(&post)
	if err != nil {
		return nil, fmt.Errorf("error saving post: %w", err)
	}

	err = store2.GlobalStore.Wall.Add(post.UserID, post.ID)
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
	Type     FeedType `json:"type"`
	ByUserID string   `json:"byUserId"`
	ByID     string   `json:"byId"`

	Count     int    `json:"count"`
	PageToken string `json:"pageToken"`
}

func (a *API) PostsList(ctx context.Context, v *Viewer, r *PostsListReq) (*PostsList, error) {
	defer tracer.FromCtx(ctx).StartChild("API.PostsList").Stop()

	if r.ByUserID != "" {
		return a.postsListByUser(ctx, v, r)
	}
	if r.ByID != "" {
		return a.postsListByID(ctx, v, r)
	}

	var err error
	var postIds []int

	if r.Type == "FEED" {
		postIds, err = pkg.GetFeedPostIds(v.UserID)
	} else if r.Type == "DISCOVER" || r.Type == "" {
		postIds, err = store2.GlobalStore.Wall.GetLatest()
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

	postsMap, err := store2.GlobalStore.Posts.Get(postIds)
	if err != nil {
		return nil, err
	}

	var posts []*store.Post
	for _, postID := range postIds {
		posts = append(posts, postsMap[postID])
	}

	result.Items = transformPostBatch(ctx, posts, v.UserID)

	return result, nil
}

func (a *API) postsListByID(ctx context.Context, v *Viewer, r *PostsListReq) (*PostsList, error) {
	postID, _ := strconv.Atoi(r.ByID)

	posts, err := store2.GlobalStore.Posts.Get([]int{postID})
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	} else if posts[postID] == nil || posts[postID].IsDeleted {
		return nil, Error("PostNotFound")
	}

	postsWrapped := transformPostBatch(ctx, []*store.Post{posts[postID]}, v.UserID)

	return &PostsList{Items: postsWrapped}, nil
}

func (a *API) postsListByUser(ctx context.Context, v *Viewer, r *PostsListReq) (*PostsList, error) {
	userID, _ := strconv.Atoi(r.ByUserID)
	if userID <= 0 {
		return nil, Error("IncorrectUserId")
	}

	if r.Count < 0 {
		return nil, Error("IncorrectCount")
	} else if r.Count > 50 {
		return nil, Error("IncorrectCount")
	}

	lessThan, _ := strconv.Atoi(r.PageToken)
	count := 10
	if r.Count != 0 {
		count = r.Count
	}

	postIds, err := store2.GlobalStore.Wall.Get([]int{userID}, lessThan, count)
	if err != nil {
		return nil, fmt.Errorf("error getting posted edges: %w", err)
	}

	postsMap, err := store2.GlobalStore.Posts.Get(postIds)
	if err != nil {
		return nil, fmt.Errorf("error getting posts: %w", err)
	}

	var posts []*store.Post

	for _, edge := range postIds {
		post := postsMap[edge]
		if post != nil {
			posts = append(posts, post)
		}
	}

	result := transformPostBatch(ctx, posts, v.UserID)

	nextAfter := ""
	if len(postIds) == count {
		nextAfter = strconv.Itoa(postIds[len(postIds)-1])
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

	posts, err := store2.GlobalStore.Posts.Get([]int{postID})
	if err != nil {
		return nil, err
	} else if posts[postID] == nil {
		return nil, Error("PostNotFound")
	}

	post := posts[postID]
	if viewer.UserID == 0 {
		return nil, Error("NotAuthorized")
	}
	if post.UserID != viewer.UserID {
		return nil, Error("AccessDenied")
	}

	post.IsDeleted = true

	err = store2.GlobalStore.Posts.Update(post)
	pkg.LogErr(err)

	err = store2.GlobalStore.Wall.Delete(post.UserID, post.ID)
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

	posts, err := store2.GlobalStore.Posts.Get([]int{postID})
	if err != nil {
		return nil, err
	} else if posts[postID] == nil {
		return nil, Error("PostNotFound")
	}

	if r.Action == Unlike {
		err = store2.GlobalStore.Likes.Remove(postID, v.UserID)
		if err != nil {
			return nil, err
		}
	} else {
		err = store2.GlobalStore.Likes.Add(postID, v.UserID)
		if err != nil {
			return nil, err
		}
	}

	return &Void{}, nil
}
