package server

import (
	"context"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func createViewerContext(userID int) context.Context {
	return context.WithValue(context.Background(), CtxViewerKey, &Viewer{UserID: userID})
}

func TestApi_usersList(t *testing.T) {
	srv := UserServer{}
	ctx := createViewerContext(0)

	closer := createTestDB(t)
	defer closer()

	user := store.User{Name: "Test user"}
	_ = store2.GlobalStore.Users.Add(&user)

	resp, err := srv.List(ctx, &api.UsersListReq{
		UserIds: []string{strconv.Itoa(user.ID)},
	})
	require.NoError(t, err)
	require.Len(t, resp.Users, 1)
	require.Equal(t, resp.Users[0].Id, strconv.Itoa(user.ID))
	require.Equal(t, "Test user", resp.Users[0].Name)
}

func TestAPI_setStatus(t *testing.T) {
	srv := UserServer{}

	closer := createTestDB(t)
	defer closer()

	user := store.User{}
	_ = store2.GlobalStore.Users.Add(&user)
	ctx := createViewerContext(user.ID)

	_, err := srv.SetStatus(ctx, &api.UsersSetStatus{
		Status: "Test status",
	})
	require.NoError(t, err)

	resp, err := srv.List(ctx, &api.UsersListReq{
		UserIds: []string{strconv.Itoa(user.ID)},
	})
	require.NoError(t, err)
	require.Equal(t, "Test status", resp.Users[0].Status)
}

func TestAPI_follow(t *testing.T) {
	srv := UserServer{}

	closer := createTestDB(t)
	defer closer()

	user1 := store.User{}
	_ = store2.GlobalStore.Users.Add(&user1)
	ctx := createViewerContext(user1.ID)

	user2 := store.User{}
	_ = store2.GlobalStore.Users.Add(&user2)

	// Follow
	_, err := srv.Follow(ctx, &api.UsersFollowReq{
		TargetId: strconv.Itoa(user2.ID),
		Action:   api.SubscribeAction_FOLLOW,
	})
	require.NoError(t, err)

	resp, err := srv.List(ctx, &api.UsersListReq{
		UserIds: []string{strconv.Itoa(user2.ID)},
	})
	require.NoError(t, err)
	require.True(t, resp.Users[0].IsFollowing)

	// Unfollow
	_, err = srv.Follow(ctx, &api.UsersFollowReq{
		TargetId: strconv.Itoa(user2.ID),
		Action:   api.SubscribeAction_UNFOLLOW,
	})
	require.NoError(t, err)

	resp, err = srv.List(ctx, &api.UsersListReq{
		UserIds: []string{strconv.Itoa(user2.ID)},
	})
	require.NoError(t, err)
	require.False(t, resp.Users[0].IsFollowing)
}
