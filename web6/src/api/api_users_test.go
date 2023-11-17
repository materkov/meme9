package api

import (
	"github.com/materkov/meme9/web6/src/store"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestApi_usersList(t *testing.T) {
	api := API{}
	v := Viewer{}

	closer := createTestDB(t)
	defer closer()

	userID, _ := store.GlobalStore.AddObject(store.ObjTypeUser, &store.User{Name: "Test user"})

	resp, err := api.usersList(&v, &UsersListReq{
		UserIds: []string{strconv.Itoa(userID)},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Equal(t, resp[0].ID, strconv.Itoa(userID))
	require.Equal(t, "Test user", resp[0].Name)
}

func TestAPI_setStatus(t *testing.T) {
	api := API{}

	closer := createTestDB(t)
	defer closer()

	userID, _ := store.GlobalStore.AddObject(store.ObjTypeUser, &store.User{})
	v := Viewer{UserID: userID}

	_, err := api.usersSetStatus(&v, &UsersSetStatus{
		Status: "Test status",
	})
	require.NoError(t, err)

	resp, err := api.usersList(&v, &UsersListReq{UserIds: []string{strconv.Itoa(userID)}})
	require.NoError(t, err)
	require.Equal(t, "Test status", resp[0].Status)
}

func TestAPI_follow(t *testing.T) {
	api := API{}

	closer := createTestDB(t)
	defer closer()

	user1ID, _ := store.GlobalStore.AddObject(store.ObjTypeUser, &store.User{})
	v := Viewer{UserID: user1ID}

	user2ID, _ := store.GlobalStore.AddObject(store.ObjTypeUser, &store.User{})

	// Follow
	_, err := api.usersFollow(&v, &UsersFollow{
		TargetID: strconv.Itoa(user2ID),
		Action:   Follow,
	})
	require.NoError(t, err)

	resp, err := api.usersList(&v, &UsersListReq{UserIds: []string{strconv.Itoa(user2ID)}})
	require.NoError(t, err)
	require.True(t, resp[0].IsFollowing)

	// Unfollow
	_, err = api.usersFollow(&v, &UsersFollow{
		TargetID: strconv.Itoa(user2ID),
		Action:   Unfollow,
	})
	require.NoError(t, err)

	resp, err = api.usersList(&v, &UsersListReq{UserIds: []string{strconv.Itoa(user2ID)}})
	require.NoError(t, err)
	require.False(t, resp[0].IsFollowing)
}
