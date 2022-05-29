package types

import (
	"github.com/materkov/web3/store"
	"sort"
)

type Query struct {
	Feed []*Post `json:"feed,omitempty"`
}

type QueryParams struct {
	Feed     QueryFeed     `json:"feed"`
	Mutation QueryMutation `json:"mutation"`
}

type QueryFeed struct {
	Include bool       `json:"include,omitempty"`
	UserID  int        `json:"userId,omitempty"`
	Inner   PostParams `json:"inner"`
}

type QueryMutation struct {
	Include bool           `json:"include,omitempty"`
	Inner   MutationParams `json:"inner,omitempty"`
}

func ResolveQuery(params QueryParams) (*Query, error) {
	result := &Query{}
	var err error

	if params.Feed.Include {
		userID := params.Feed.UserID
		if userID == 0 {
			userID = 10
		}
		userIds, _ := GlobalStore.ListGet(userID, store.ListSubscribedTo)
		userIds = append(userIds, userID)

		var allPosts []int
		for _, userID := range userIds {
			postIds, _ := GlobalStore.ListGet(userID, store.ListPosted)
			allPosts = append(allPosts, postIds...)
		}

		for _, postID := range allPosts {
			GlobalCachedStore.Need(postID)
		}

		var posts []*store.Post
		for _, postID := range allPosts {
			obj, err := GlobalCachedStore.ObjGet(postID)
			if err == nil {
				if post, ok := obj.(*store.Post); ok {
					posts = append(posts, post)
				}
			}
		}
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Date > posts[j].Date
		})

		for _, post := range posts {
			post, _ := ResolveGraphPost(post.ID, params.Feed.Inner)
			result.Feed = append(result.Feed, post)
		}
	}

	if params.Mutation.Include {
		ResolveMutation(params.Mutation.Inner)
	}

	return result, err
}
