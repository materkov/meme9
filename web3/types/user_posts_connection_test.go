package types

import (
	"github.com/materkov/web3/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResolveUserPostsConnection(t *testing.T) {
	st := newMockStore()

	_ = st.Store.ObjAdd(20, store.ObjectPost, &store.Post{})
	_ = st.Store.ObjAdd(21, store.ObjectPost, &store.Post{})
	_ = st.Store.ListAdd(15, store.ListPosted, 20)
	_ = st.Store.ListAdd(15, store.ListPosted, 21)

	posts, err := ResolveUserPostsConnection(st, 15, &UserPostsConnectionFields{
		TotalCount: &simpleField{},
		Edges:      &UserPostsConnectionFieldsEdges{},
	})
	require.NoError(t, err)
	require.Equal(t, 2, posts.TotalCount)
	require.Len(t, posts.Edges, 2)
	require.Equal(t, "Post:21", posts.Edges[0].ID)
	require.Equal(t, "Post:20", posts.Edges[1].ID)
}
