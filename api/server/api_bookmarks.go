package server

import (
	"context"
	"fmt"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/twitchtv/twirp"
	"strconv"
)

type Bookmarks struct{}

func (b *Bookmarks) Add(ctx context.Context, req *api.BookmarksAddReq) (*api.Void, error) {
	viewer := ctx.Value(CtxViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		return nil, ErrNotAuthorized
	}

	postID, _ := strconv.Atoi(req.PostId)
	posts, err := store2.GlobalStore.Posts.Get([]int{postID})
	if err != nil {
		return nil, fmt.Errorf("error getting post: %s", err)
	} else if posts[postID] == nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "PostNotFound")
	}

	err = store2.GlobalStore.Bookmarks.Add(postID, viewer.UserID)
	if err != nil {
		return nil, fmt.Errorf("error saving bookmark: %w", err)
	}

	return &api.Void{}, err
}

func (b *Bookmarks) Remove(ctx context.Context, req *api.BookmarksAddReq) (*api.Void, error) {
	viewer := ctx.Value(CtxViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		return nil, ErrNotAuthorized
	}

	postID, _ := strconv.Atoi(req.PostId)

	err := store2.GlobalStore.Bookmarks.Remove(postID, viewer.UserID)
	if err != nil {
		return nil, fmt.Errorf("error deleting bookmark: %w", err)
	}

	return &api.Void{}, err
}

func (b *Bookmarks) List(ctx context.Context, req *api.BookmarkListReq) (*api.BookmarkList, error) {
	viewer := ctx.Value(CtxViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		return nil, ErrNotAuthorized
	}

	count := 10

	after, _ := strconv.Atoi(req.PageToken)
	bookmarks, err := store2.GlobalStore.Bookmarks.List(viewer.UserID, after, count)
	if err != nil {
		return nil, fmt.Errorf("errog etting bookmarks: %w", err)
	}

	nextPage := ""
	if len(bookmarks) == 10 {
		nextPage = strconv.Itoa(bookmarks[9].Date)
	}

	postIds := make([]int, len(bookmarks))
	for i, item := range bookmarks {
		postIds[i] = item.PostID
	}

	posts, err := store2.GlobalStore.Posts.Get(postIds)
	if err != nil {
		return nil, fmt.Errorf("error loading posts: %w", err)
	}

	var postsList []*store.Post
	for _, postId := range postIds {
		if posts[postId] != nil {
			postsList = append(postsList, posts[postId])
		}
	}

	wrappedPosts := transformPostBatch(ctx, postsList, viewer.UserID)
	wrappedBookmarks := make([]*api.Bookmark, len(postsList))
	for i := range postsList {
		wrappedBookmarks[i] = &api.Bookmark{
			Date: transformDate(bookmarks[i].Date),
			Post: wrappedPosts[i],
		}
	}

	return &api.BookmarkList{
		Items:     wrappedBookmarks,
		PageToken: nextPage,
	}, nil
}
