package types

import (
	"github.com/materkov/web3/pkg"
	"github.com/materkov/web3/store"
	"sort"
)

type Query struct {
	Type      string    `json:"type"`
	ID        string    `json:"id"`
	Viewer    *User     `json:"viewer,omitempty"`
	Feed      []*Post   `json:"feed,omitempty"`
	VkAuthURL string    `json:"vkAuthUrl,omitempty"`
	Mutation  *Mutation `json:"mutation,omitempty"`
	Node      Node      `json:"node,omitempty"`
}

type QueryParams struct {
	Viewer    *QueryViewer   `json:"viewer"`
	Feed      *QueryFeed     `json:"feed"`
	Mutation  *QueryMutation `json:"mutation"`
	VkAuthURL *simpleField   `json:"vkAuthUrl"`
	Node      *QueryNode     `json:"node"`
}

type QueryNode struct {
	ID    string     `json:"id"`
	Inner NodeParams `json:"inner"`
}

type QueryFeed struct {
	UserID int        `json:"userId,omitempty"`
	Inner  PostParams `json:"inner"`
}

type QueryMutation struct {
	Inner MutationParams `json:"inner,omitempty"`
}

type QueryViewer struct {
	Inner UserParams `json:"inner"`
}

func ResolveQuery(viewer pkg.Viewer, params QueryParams) (*Query, error) {
	result := &Query{
		Type: "Query",
		ID:   "query",
	}
	var err error

	if params.Feed != nil {
		userID := params.Feed.UserID
		if userID == 0 {
			userID = viewer.UserID
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

	if params.Mutation != nil {
		result.Mutation = ResolveMutation(viewer, params.Mutation.Inner)
	}

	if params.VkAuthURL != nil {
		result.VkAuthURL = pkg.GetRedirectURL(viewer.Origin)
	}

	if params.Viewer != nil {
		result.Viewer, _ = ResolveUser(viewer.UserID, params.Viewer.Inner)
	}

	if params.Node != nil {
		result.Node, _ = ResolveNode(params.Node.ID, params.Node.Inner)
	}

	return result, err
}
