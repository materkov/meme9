package types

import (
	"fmt"
	"github.com/materkov/web3/pkg"
	"github.com/materkov/web3/pkg/globalid"
	"github.com/materkov/web3/store"
)

type Node interface {
	IsNode()
}

func (*User) IsNode() {}
func (*Post) IsNode() {}

type NodeParams struct {
	OnPost *PostParams `json:"onPost,omitempty"`
	OnUser *UserParams `json:"onUser,omitempty"`
}

func ResolveNode(cachedStore *store.CachedStore, id string, params NodeParams, viewer *pkg.Viewer) (Node, error) {
	objectID, _ := globalid.Parse(id)

	switch objectID := objectID.(type) {
	case *globalid.PostID:
		return ResolveGraphPost(cachedStore, objectID.PostID, params.OnPost, viewer)
	case *globalid.UserID:
		return ResolveUser(cachedStore, objectID.UserID, params.OnUser, viewer)
	default:
		return nil, fmt.Errorf("incorrect id")
	}
}
