package types

import (
	"github.com/materkov/web3/pkg"
	"github.com/materkov/web3/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResolveMutation_AddPost(t *testing.T) {
	st := store.CachedStore{
		Store:    &store.MockStore{},
		Needed:   map[int]bool{},
		ObjCache: map[int]store.CachedItem{},
	}
	viewer := pkg.Viewer{}

	result := ResolveMutation(&st, viewer, MutationParams{
		AddPost: &MutationAddPost{
			Text: "hello world",
		},
	})

	require.NotEmpty(t, result.AddPost.ID)

	node, err := ResolveNode(&st, result.AddPost.ID, NodeParams{OnPost: &PostParams{Text: &PostText{}}})
	require.NoError(t, err)
	require.Equal(t, result.AddPost.ID, node.(*Post).ID)
	require.Equal(t, "hello world", node.(*Post).Text)
}
