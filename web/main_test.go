package main

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	db.SetMaxOpenConns(1)

	_, err = db.Exec(migrations)
	require.NoError(t, err)

	// TODO replace this global
	store = Store{db: db}
}

func TestHandleIndex(t *testing.T) {
	setupDB(t)

	require.NoError(t, store.Follow(10, 11))

	require.NoError(t, store.AddPost(&Post{ID: 1, UserID: 10}))
	require.NoError(t, store.AddPost(&Post{ID: 2, UserID: 11}))
	require.NoError(t, store.AddPost(&Post{ID: 3, UserID: 12}))

	resp, err := handleIndex("", &Viewer{UserID: 10})
	require.NoError(t, err)
	require.NotNil(t, resp)

	posts := resp.GetFeedRenderer().Posts
	require.Len(t, posts, 2)
	require.Equal(t, "2", posts[0].Id)
	require.Equal(t, "1", posts[1].Id)
}

func TestHandleIndex_NotAuthorized(t *testing.T) {
	setupDB(t)

	resp, err := handleIndex("", &Viewer{})
	require.NoError(t, err)
	require.Len(t, resp.GetFeedRenderer().Posts, 0)
}
