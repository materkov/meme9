package app

import (
	"context"
	"testing"

	"github.com/materkov/meme9/web/pb"
	"github.com/materkov/meme9/web/store"
	"github.com/stretchr/testify/require"
)

func TestRelations_Follow(t *testing.T) {
	setupDB(t)
	r := Relations{App: &App{Store: ObjectStore}}

	err := r.App.Store.ObjAdd(&store.StoredObject{ID: 15, User: &store.User{}})
	require.NoError(t, err)

	ctx := WithViewerContext(context.Background(), &Viewer{UserID: 44})

	// Follow
	_, err = r.Follow(ctx, &pb.RelationsFollowRequest{UserId: "15"})
	require.NoError(t, err)

	// Check status
	respCheck, err := r.Check(ctx, &pb.RelationsCheckRequest{UserId: "15"})
	require.NoError(t, err)
	require.True(t, respCheck.IsFollowing)

	// Check another user
	respCheck, err = r.Check(ctx, &pb.RelationsCheckRequest{UserId: "16"})
	require.NoError(t, err)
	require.False(t, respCheck.IsFollowing)

	// Follow again (no-op
	_, err = r.Follow(ctx, &pb.RelationsFollowRequest{UserId: "15"})
	require.NoError(t, err)

	// Unfollow
	_, err = r.Unfollow(ctx, &pb.RelationsUnfollowRequest{UserId: "15"})
	require.NoError(t, err)

	// Check status again
	respCheck, err = r.Check(ctx, &pb.RelationsCheckRequest{UserId: "15"})
	require.NoError(t, err)
	require.False(t, respCheck.IsFollowing)
}

func TestRelations_Follow_Errors(t *testing.T) {
	setupDB(t)
	r := Relations{App: &App{Store: ObjectStore}}

	ctx := WithViewerContext(context.Background(), &Viewer{UserID: 0})

	// Not auth
	_, err := r.Follow(ctx, &pb.RelationsFollowRequest{UserId: "15"})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "auth")

	// User not exists
	ctx = WithViewerContext(context.Background(), &Viewer{UserID: 22})
	_, err = r.Follow(ctx, &pb.RelationsFollowRequest{UserId: "15"})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "exists")

	// Cannot follow yourself
	_, err = r.Follow(ctx, &pb.RelationsFollowRequest{UserId: "22"})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "yourself")
}
