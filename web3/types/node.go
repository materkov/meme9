package types

import (
	"fmt"
	"github.com/materkov/web3/pkg"
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

func ResolveNode(cachedStore *store.CachedStore, id string, params NodeParams) (Node, error) {
	objectType, objectID, _ := pkg.ParseGlobalID(id)

	switch objectType {
	case pkg.GlobalIDPost:
		return ResolveGraphPost(cachedStore, objectID, params.OnPost)
	case pkg.GlobalIDUser:
		return ResolveUser(cachedStore, objectID, params.OnUser)
	default:
		return nil, fmt.Errorf("incorrect id")
	}
}
