package server

import (
	"context"
	"fmt"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/pkg"
	"github.com/materkov/meme9/api/src/pkg/tracer"
	"github.com/materkov/meme9/api/src/pkg/utils"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/twitchtv/twirp"
	"net/url"
	"slices"
	"strconv"
	"sync"
	"time"
)

type ctxKey string

var (
	// TODO func ViewerFromContext
	CtxViewerKey ctxKey = "viewer"
)

type PostsServer struct{}

func transformDate(date int) string {
	return time.Unix(int64(date), 0).Format(time.RFC3339)
}

func transformPostBatch(ctx context.Context, posts []*store.Post, viewerID int) []*api.Post {
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
	wg.Add(4)

	var counters map[int]int
	var isLiked map[int]bool
	go func() {
		var err error
		counters, isLiked, err = store2.GlobalStore.Likes.Get(ctx, postIds, viewerID)
		pkg.LogErr(err)
		wg.Done()
	}()

	usersWrappedMap := map[int]*api.User{}
	go func() {
		usersMap, err := store2.GlobalStore.Users.Get(utils.UniqueIds(userIds))
		pkg.LogErr(err)

		for userID, user := range usersMap {
			usersWrappedMap[userID], err = transformUser(userID, user, viewerID)
			pkg.LogErr(err)
		}
		wg.Done()
	}()

	pollWrappedMap := map[int]*api.Poll{}
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

	isBookmarked := map[int]bool{}
	go func() {
		var err error
		isBookmarked, err = store2.GlobalStore.Bookmarks.IsBookmarked(postIds, viewerID)
		pkg.LogErr(err)

		wg.Done()
	}()
	wg.Wait()

	result := make([]*api.Post, len(posts))

	for i, post := range posts {
		if post.IsDeleted {
			result[i] = &api.Post{
				Id:        strconv.Itoa(post.ID),
				IsDeleted: true,
			}
			continue
		}

		var wrappedLink *api.PostLink
		if post.Link != nil {
			host := ""

			parsedURL, err := url.Parse(post.Link.URL)
			if err == nil {
				host = parsedURL.Host
			}

			proxiedUrl := fmt.Sprintf("https://3c6ef5be-e5f9-4e47-9a68-bd635323a374.selcdn.net/image-proxy?url=%s", url.QueryEscape(post.Link.ImageURL))

			wrappedLink = &api.PostLink{
				Url:         post.Link.URL,
				Title:       post.Link.Title,
				Description: post.Link.Description,
				ImageUrl:    proxiedUrl,
				Domain:      host,
			}
		}

		pollTransformed := pollWrappedMap[post.PollID]

		result[i] = &api.Post{
			Id:     strconv.Itoa(post.ID),
			UserId: strconv.Itoa(post.UserID),
			Date:   transformDate(post.Date),
			Text:   post.Text,
			User:   usersWrappedMap[post.UserID],

			LikesCount: int32(counters[post.ID]),
			IsLiked:    isLiked[post.ID],

			Link: wrappedLink,
			Poll: pollTransformed,

			IsBookmarked: isBookmarked[post.ID],
		}
	}

	return result
}

func (p *PostsServer) Add(ctx context.Context, req *api.AddReq) (*api.Post, error) {
	if req.Text == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "TextEmpty")
	}
	if len(req.Text) > 5000 {
		return nil, twirp.NewError(twirp.InvalidArgument, "TextTooLong")
	}

	viewer := ctx.Value(CtxViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "NotAuthorized")
	}

	pollID := 0
	if req.PollId != "" {
		pollID, _ = strconv.Atoi(req.PollId)
	}

	post := store.Post{
		UserID: viewer.UserID,
		Date:   int(time.Now().Unix()),
		Text:   req.Text,
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
	go func() {
		err := pkg.PushRealtimeEvent(post.UserID, map[string]interface{}{
			"type":   "NEW_POST",
			"postId": post.ID,
		})
		pkg.LogErr(err)
	}()

	return transformPostBatch(ctx, []*store.Post{&post}, viewer.UserID)[0], nil
}

func (p *PostsServer) List(ctx context.Context, req *api.ListReq) (*api.PostsList, error) {
	defer tracer.FromCtx(ctx).StartChild("API.PostsList").Stop()

	if req.ByUserId != "" {
		return p.postsListByUser(ctx, req)
	}
	if req.ById != "" {
		return p.postsListByID(ctx, req)
	}

	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	var err error
	var postIds []int

	if req.Type == api.FeedType_FEED {
		postIds, err = pkg.GetFeedPostIds(viewer.UserID)
	} else if req.Type == api.FeedType_DISCOVER {
		postIds, err = store2.GlobalStore.Wall.GetLatest()
	} else {
		err = twirp.NewError(twirp.InvalidArgument, "InvalidFeedType")
	}

	if err != nil {
		return nil, err
	}

	if req.PageToken != "" {
		cursor, _ := strconv.Atoi(req.PageToken)
		cursorIdx := slices.Index(postIds, cursor)
		if cursorIdx != -1 && cursorIdx != len(postIds)-1 {
			postIds = postIds[cursorIdx+1:]
		}
	}

	result := &api.PostsList{}

	count := int(req.Count)
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

	result.Items = transformPostBatch(ctx, posts, viewer.UserID)

	return result, nil
}

func (p *PostsServer) postsListByID(ctx context.Context, r *api.ListReq) (*api.PostsList, error) {
	postID, _ := strconv.Atoi(r.ById)

	posts, err := store2.GlobalStore.Posts.Get([]int{postID})
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	}

	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	postsWrapped := transformPostBatch(ctx, []*store.Post{posts[postID]}, viewer.UserID)

	return &api.PostsList{Items: postsWrapped}, nil
}

func (p *PostsServer) postsListByUser(ctx context.Context, r *api.ListReq) (*api.PostsList, error) {
	userID, _ := strconv.Atoi(r.ByUserId)
	if userID <= 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "IncorrectUserId")
	}

	if r.Count < 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "IncorrectCount")
	} else if r.Count > 50 {
		return nil, twirp.NewError(twirp.InvalidArgument, "IncorrectCount")
	}

	lessThan, _ := strconv.Atoi(r.PageToken)
	count := 10
	if r.Count != 0 {
		count = int(r.Count)
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

	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	result := transformPostBatch(ctx, posts, viewer.UserID)

	nextAfter := ""
	if len(postIds) == count {
		nextAfter = strconv.Itoa(postIds[len(postIds)-1])
	}

	return &api.PostsList{
		Items:     result,
		PageToken: nextAfter,
	}, nil
}

func (p *PostsServer) Delete(ctx context.Context, r *api.PostsDeleteReq) (*api.Void, error) {
	postID, _ := strconv.Atoi(r.PostId)

	posts, err := store2.GlobalStore.Posts.Get([]int{postID})
	if err != nil {
		return nil, err
	} else if posts[postID] == nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "PostNotFound")
	}

	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	post := posts[postID]
	if viewer.UserID == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "NotAuthorized")
	}
	if post.UserID != viewer.UserID {
		return nil, twirp.NewError(twirp.InvalidArgument, "AccessDenied")
	}

	post.IsDeleted = true

	err = store2.GlobalStore.Posts.Update(post)
	pkg.LogErr(err)

	err = store2.GlobalStore.Wall.Delete(post.UserID, post.ID)
	pkg.LogErr(err)

	return &api.Void{}, nil
}

func (*PostsServer) Like(ctx context.Context, r *api.PostsLikeReq) (*api.Void, error) {
	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	if viewer.UserID == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "NotAuthorized")
	}

	postID, _ := strconv.Atoi(r.PostId)

	if r.Action == api.PostLikeAction_UNLIKE {
		err := store2.GlobalStore.Likes.Remove(postID, viewer.UserID)
		if err != nil {
			return nil, err
		}
	} else {
		posts, err := store2.GlobalStore.Posts.Get([]int{postID})
		if err != nil {
			return nil, err
		} else if posts[postID] == nil || posts[postID].IsDeleted {
			return nil, twirp.NewError(twirp.InvalidArgument, "PostNotFound")
		}

		err = store2.GlobalStore.Likes.Add(postID, viewer.UserID)
		if err != nil {
			return nil, err
		}
	}

	return &api.Void{}, nil
}
