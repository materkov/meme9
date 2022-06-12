package types

import (
	"github.com/materkov/web3/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func newMockStore() *store.CachedStore {
	st := &store.MockStore{}
	cachedStore := &store.CachedStore{
		Store:    st,
		Needed:   map[int]bool{},
		ObjCache: map[int]store.CachedItem{},
	}

	return cachedStore
}

func TestResolveUser(t *testing.T) {
	st := newMockStore()
	_ = st.Store.ObjAdd(15, store.ObjectUser, store.User{
		ID:   15,
		Name: "user 15 name",
	})

	user, err := ResolveUser(st, 15, UserParams{
		Name:  &simpleField{},
		Posts: nil,
	})
	require.NoError(t, err)
	require.Equal(t, "User:15", user.ID)
	require.Equal(t, "user 15 name", user.Name)
}
