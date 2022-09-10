package app

import (
	"context"
	"testing"

	"github.com/materkov/meme9/web/store"
	"github.com/stretchr/testify/require"
)

func TestFeed_GetHeader(t *testing.T) {
	setupDB(t)
	f := Feed{App: &App{Store: ObjectStore}}

	require.NoError(t, f.App.Store.ObjAdd(&store.StoredObject{
		ID:   14,
		User: &store.User{},
	}))

	// No auth
	ctx := WithViewerContext(context.Background(), &Viewer{})

	resp, err := f.GetHeader(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, "/", resp.GetRenderer().MainUrl)
	require.False(t, resp.GetRenderer().IsAuthorized)

	// With auth
	ctx = WithViewerContext(context.Background(), &Viewer{UserID: 14})

	resp, err = f.GetHeader(ctx, nil)
	require.NoError(t, err)
	require.True(t, resp.GetRenderer().IsAuthorized)
}
