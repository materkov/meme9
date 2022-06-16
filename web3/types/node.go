package types

import (
	"fmt"
	"github.com/materkov/web3/store"
	"strconv"
	"strings"
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
	if strings.HasPrefix(id, "Post:") {
		id = strings.TrimPrefix(id, "Post:")
		idInt, _ := strconv.Atoi(id)
		return ResolveGraphPost(cachedStore, idInt, params.OnPost)
	} else if strings.HasPrefix(id, "User:") {
		id = strings.TrimPrefix(id, "User:")
		idInt, _ := strconv.Atoi(id)
		return ResolveUser(cachedStore, idInt, params.OnUser)
	} else {
		return nil, fmt.Errorf("incorrect id")
	}
}
