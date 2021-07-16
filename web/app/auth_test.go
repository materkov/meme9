package app

import (
	"context"
	"net/http"
	"testing"

	"github.com/materkov/meme9/web/store"
	"github.com/stretchr/testify/require"
)

func TestTryVkAuth(t *testing.T) {
	setupDB(t)
	auth := Auth{Store: ObjectStore}
	ctx := context.Background()

	// Not registered
	userID1, err := auth.tryVkAuth(ctx, "https://meme.mmaks.me/?vk_user_id=123&vk_data=15&not_vk_param=568&sign=Ee4LRaYNMC41nqNamNr1dziR5xYemmd2tAo-eAjRsNg")
	require.NoError(t, err)
	require.NotZero(t, userID1)

	// Already registered
	userID2, err := auth.tryVkAuth(ctx, "https://meme.mmaks.me/?vk_user_id=123&vk_data=15&not_vk_param=568&sign=Ee4LRaYNMC41nqNamNr1dziR5xYemmd2tAo-eAjRsNg")
	require.NoError(t, err)
	require.Equal(t, userID1, userID2)

	// Another VK id
	userID3, err := auth.tryVkAuth(ctx, "https://meme.mmaks.me/?vk_user_id=124&vk_data=15&not_vk_param=568&sign=JsUrOlIzwBrNUe4zViaTFxUqRkn0zW5h3CqyngaxeNE")
	require.NoError(t, err)
	require.NotZero(t, userID3)
	require.NotEqual(t, userID2, userID3)

	// Empty urls
	_, err = auth.tryVkAuth(ctx, "")
	require.Equal(t, ErrNotApplicable, err)

	_, err = auth.tryVkAuth(ctx, "https://meme.mmaks.me/?some_param=123")
	require.Equal(t, ErrNotApplicable, err)

	_, err = auth.tryVkAuth(ctx, "https://meme.mmaks.me/")
	require.Equal(t, ErrNotApplicable, err)

	_, err = auth.tryVkAuth(ctx, "/")
	require.Equal(t, ErrNotApplicable, err)

	// Invalid URL
	_, err = auth.tryVkAuth(ctx, ":")
	require.Equal(t, ErrNotApplicable, err)

	// Incorrect sign
	_, err = auth.tryVkAuth(ctx, "https://meme.mmaks.me/?vk_user_id=124&vk_data=15&not_vk_param=568&sign=BAD-SIGN")
	require.Equal(t, ErrAuthFailed, err)
}

func TestTryCookieAuth(t *testing.T) {
	setupDB(t)
	auth := Auth{Store: ObjectStore}

	ObjectStore.ObjAdd(&store.StoredObject{ID: 1, Token: &store.Token{
		ID:     1,
		Token:  "1-test-token",
		UserID: 167,
	}})

	r, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	r.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: "1-test-token",
	})
	r.Header.Set("x-csrf-token", "yERDvFm5LK747IUA/0Q1hHk3VnSVI4EeszOvK7x0/Ig=")

	token, err := auth.tryCookieAuth(r)
	require.NoError(t, err)
	require.Equal(t, "1-test-token", token.Token)
	require.Equal(t, 167, token.UserID)
}

func TestTryCookieAuth_Failed(t *testing.T) {
	setupDB(t)
	auth := Auth{}

	r, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	// No cookie
	_, err = auth.tryCookieAuth(r)
	require.Equal(t, ErrNotApplicable, err)

	r.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: "test-token",
	})

	// Incorrect CSRF
	_, err = auth.tryCookieAuth(r)
	require.Equal(t, ErrAuthFailed, err)

	// Token not found
	r.Header.Set("x-csrf-token", "yERDvFm5LK747IUA/0Q1hHk3VnSVI4EeszOvK7x0/Ig=")
	_, err = auth.tryCookieAuth(r)
	require.Equal(t, ErrAuthFailed, err)
}

func TestTryHeaderAuth(t *testing.T) {
	setupDB(t)
	auth := Auth{}

	err := ObjectStore.ObjAdd(&store.StoredObject{ID: 1, Token: &store.Token{
		ID:     1,
		Token:  "1-test-token",
		UserID: 12,
	}})
	require.NoError(t, err)

	ctx := context.Background()

	token, err := auth.tryHeaderAuth(ctx, "1-test-token")
	require.NoError(t, err)
	require.Equal(t, token.ID, 1)

	token, err = auth.tryHeaderAuth(ctx, "Bearer 1-test-token")
	require.NoError(t, err)
	require.Equal(t, token.ID, 1)

	// Empty header
	_, err = auth.tryHeaderAuth(ctx, "")
	require.ErrorIs(t, err, ErrNotApplicable)

	// Not found
	_, err = auth.tryHeaderAuth(ctx, "incorrect-token")
	require.Equal(t, ErrAuthFailed, err)

	_, err = auth.tryHeaderAuth(ctx, "1-incorrect-token")
	require.Equal(t, ErrAuthFailed, err)
}
