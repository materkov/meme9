package types

import (
	"github.com/hashicorp/go-multierror"
	"github.com/materkov/web3/pkg"
	"github.com/materkov/web3/pkg/globalid"
	"github.com/materkov/web3/store"
)

type UserPostsConnection struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	TotalCount *int    `json:"totalCount,omitempty"`
	Edges      []*Post `json:"edges,omitempty"`
}

type UserPostsConnectionFields struct {
	TotalCount *simpleField                    `json:"count"`
	Edges      *UserPostsConnectionFieldsEdges `json:"edges"`
}

type UserPostsConnectionFieldsEdges struct {
	Inner *PostParams `json:"inner"`
}

func ResolveUserPostsConnection(st *store.CachedStore, userID int, fields *UserPostsConnectionFields, viewer *pkg.Viewer) (*UserPostsConnection, error) {
	result := &UserPostsConnection{
		ID:   globalid.Create(globalid.Stub{}),
		Type: "UserPostsConnection",
	}

	var errors error

	if fields.TotalCount != nil {
		count, err := st.Store.ListCount(userID, store.ListPosted)
		if err != nil {
			errors = multierror.Append(errors, err)
		}

		result.TotalCount = &count
	}

	if fields.Edges != nil {
		postIds, err := st.Store.ListGet(userID, store.ListPosted)
		if err != nil {
			errors = multierror.Append(errors, err)
		}

		for _, postID := range postIds {
			st.Need(postID)
		}

		for _, postID := range postIds {
			post, err := ResolveGraphPost(st, postID, fields.Edges.Inner, viewer)
			if err != nil {
				errors = multierror.Append(errors, err)
			}

			result.Edges = append(result.Edges, post)
		}
	}

	return result, errors
}
