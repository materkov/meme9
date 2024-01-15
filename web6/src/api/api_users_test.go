package api

import (
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestApi_usersList(t *testing.T) {
	api := API{}
	v := Viewer{}

	closer := createTestDB(t)
	defer closer()

	user := store.User{Name: "Test user"}
	_ = store2.GlobalStore.Users.Add(&user)

	resp, err := api.usersList(&v, &UsersListReq{
		UserIds: []string{strconv.Itoa(user.ID)},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Equal(t, resp[0].ID, strconv.Itoa(user.ID))
	require.Equal(t, "Test user", resp[0].Name)
}

func TestAPI_setStatus(t *testing.T) {
	api := API{}

	closer := createTestDB(t)
	defer closer()

	user := store.User{}
	_ = store2.GlobalStore.Users.Add(&user)
	v := Viewer{UserID: user.ID}

	_, err := api.usersSetStatus(&v, &UsersSetStatusReq{
		Status: "Test status",
	})
	require.NoError(t, err)

	resp, err := api.usersList(&v, &UsersListReq{UserIds: []string{strconv.Itoa(user.ID)}})
	require.NoError(t, err)
	require.Equal(t, "Test status", resp[0].Status)
}

func TestAPI_follow(t *testing.T) {
	api := API{}

	closer := createTestDB(t)
	defer closer()

	user1 := store.User{}
	_ = store2.GlobalStore.Users.Add(&user1)
	v := Viewer{UserID: user1.ID}

	user2 := store.User{}
	_ = store2.GlobalStore.Users.Add(&user2)

	// Follow
	_, err := api.usersFollow(&v, &UsersFollow{
		TargetID: strconv.Itoa(user2.ID),
		Action:   Follow,
	})
	require.NoError(t, err)

	resp, err := api.usersList(&v, &UsersListReq{UserIds: []string{strconv.Itoa(user2.ID)}})
	require.NoError(t, err)
	require.True(t, resp[0].IsFollowing)

	// Unfollow
	_, err = api.usersFollow(&v, &UsersFollow{
		TargetID: strconv.Itoa(user2.ID),
		Action:   Unfollow,
	})
	require.NoError(t, err)

	resp, err = api.usersList(&v, &UsersListReq{UserIds: []string{strconv.Itoa(user2.ID)}})
	require.NoError(t, err)
	require.False(t, resp[0].IsFollowing)
}
