package types

import (
	"fmt"
	"strconv"
	"strings"
)

type Node interface {
	IsNode()
}

func (*User) IsNode() {}
func (*Post) IsNode() {}

type NodeParams struct {
	OnPost PostParams `json:"onPost,omitempty"`
	OnUser UserParams `json:"onUser,omitempty"`
}

func ResolveNode(id string, params NodeParams) (Node, error) {
	if strings.HasPrefix(id, "Post:") {
		id = strings.TrimPrefix(id, "Post:")
		idInt, _ := strconv.Atoi(id)
		return ResolveGraphPost(idInt, params.OnPost)
	} else if strings.HasPrefix(id, "User:") {
		id = strings.TrimPrefix(id, "User:")
		idInt, _ := strconv.Atoi(id)
		return ResolveUser(idInt, params.OnUser)
	} else {
		return nil, fmt.Errorf("incorrect id")
	}
}
