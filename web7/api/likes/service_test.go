package likes

import (
	"context"
	"errors"
	"testing"

	"github.com/materkov/meme9/web7/adapters/users"
	"github.com/materkov/meme9/web7/api"
	"github.com/materkov/meme9/web7/api/likes/mocks"
	likesapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/likes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func initService(t *testing.T) (*Service, *mocks.MockLikesAdapter, *mocks.MockUsersAdapter, func()) {
	ctrl := gomock.NewController(t)
	closer := func() {
		ctrl.Finish()
	}

	mockLikes := mocks.NewMockLikesAdapter(ctrl)
	mockUsers := mocks.NewMockUsersAdapter(ctrl)
	return NewService(mockLikes, mockUsers), mockLikes, mockUsers, closer
}

func TestService_GetLikers(t *testing.T) {
	service, mockLikes, mockUsers, closer := initService(t)
	defer closer()

	t.Run("success with usernames", func(t *testing.T) {
		ctx := context.Background()
		postID := "post-123"
		userID1 := "user-1"
		userID2 := "user-2"
		username1 := "alice"
		username2 := "bob"

		mockLikes.EXPECT().
			GetLikers(ctx, postID, "", 20).
			Return([]string{userID1, userID2}, "", nil).
			Times(1)

		mockUsers.EXPECT().
			GetByIDs(ctx, []string{userID1, userID2}).
			Return(map[string]*users.User{
				userID1: {ID: userID1, Username: username1},
				userID2: {ID: userID2, Username: username2},
			}, nil).
			Times(1)

		resp, err := service.GetLikers(ctx, &likesapi.GetLikersRequest{
			PostId: postID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Likers, 2)
		require.Equal(t, userID1, resp.Likers[0].UserId)
		require.Equal(t, username1, resp.Likers[0].Username)
		require.Equal(t, userID2, resp.Likers[1].UserId)
		require.Equal(t, username2, resp.Likers[1].Username)
		require.Empty(t, resp.PageToken)
	})

	t.Run("success with pagination", func(t *testing.T) {
		ctx := context.Background()
		postID := "post-123"
		userID1 := "user-1"
		userID2 := "user-2"
		pageToken := "token-123"

		mockLikes.EXPECT().
			GetLikers(ctx, postID, pageToken, 10).
			Return([]string{userID1, userID2}, "next-token", nil).
			Times(1)

		mockUsers.EXPECT().
			GetByIDs(ctx, []string{userID1, userID2}).
			Return(map[string]*users.User{
				userID1: {ID: userID1, Username: "alice"},
				userID2: {ID: userID2, Username: "bob"},
			}, nil).
			Times(1)

		resp, err := service.GetLikers(ctx, &likesapi.GetLikersRequest{
			PostId:    postID,
			PageToken: pageToken,
			Count:     10,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Likers, 2)
		require.Equal(t, "next-token", resp.PageToken)
	})

	t.Run("empty post id", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.GetLikers(ctx, &likesapi.GetLikersRequest{
			PostId: "",
		})
		api.RequireError(t, err, "post_id_required")
	})

	t.Run("empty result", func(t *testing.T) {
		ctx := context.Background()
		postID := "post-123"

		mockLikes.EXPECT().
			GetLikers(ctx, postID, "", 20).
			Return([]string{}, "", nil).
			Times(1)

		mockUsers.EXPECT().
			GetByIDs(ctx, []string{}).
			Return(map[string]*users.User{}, nil).
			Times(1)

		resp, err := service.GetLikers(ctx, &likesapi.GetLikersRequest{
			PostId: postID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Empty(t, resp.Likers)
		require.Empty(t, resp.PageToken)
	})

	t.Run("default count when zero", func(t *testing.T) {
		ctx := context.Background()
		postID := "post-123"

		mockLikes.EXPECT().
			GetLikers(ctx, postID, "", 20).
			Return([]string{}, "", nil).
			Times(1)

		mockUsers.EXPECT().
			GetByIDs(ctx, []string{}).
			Return(map[string]*users.User{}, nil).
			Times(1)

		resp, err := service.GetLikers(ctx, &likesapi.GetLikersRequest{
			PostId: postID,
			Count:  0,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("default count when negative", func(t *testing.T) {
		ctx := context.Background()
		postID := "post-123"

		mockLikes.EXPECT().
			GetLikers(ctx, postID, "", 20).
			Return([]string{}, "", nil).
			Times(1)

		mockUsers.EXPECT().
			GetByIDs(ctx, []string{}).
			Return(map[string]*users.User{}, nil).
			Times(1)

		resp, err := service.GetLikers(ctx, &likesapi.GetLikersRequest{
			PostId: postID,
			Count:  -5,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("max count capped at 100", func(t *testing.T) {
		ctx := context.Background()
		postID := "post-123"

		mockLikes.EXPECT().
			GetLikers(ctx, postID, "", 100).
			Return([]string{}, "", nil).
			Times(1)

		mockUsers.EXPECT().
			GetByIDs(ctx, []string{}).
			Return(map[string]*users.User{}, nil).
			Times(1)

		resp, err := service.GetLikers(ctx, &likesapi.GetLikersRequest{
			PostId: postID,
			Count:  200,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("missing usernames handled gracefully", func(t *testing.T) {
		ctx := context.Background()
		postID := "post-123"
		userID1 := "user-1"
		userID2 := "user-2"

		mockLikes.EXPECT().
			GetLikers(ctx, postID, "", 20).
			Return([]string{userID1, userID2}, "", nil).
			Times(1)

		// Return only one user, missing the other
		mockUsers.EXPECT().
			GetByIDs(ctx, []string{userID1, userID2}).
			Return(map[string]*users.User{
				userID1: {ID: userID1, Username: "alice"},
			}, nil).
			Times(1)

		resp, err := service.GetLikers(ctx, &likesapi.GetLikersRequest{
			PostId: postID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Likers, 2)
		require.Equal(t, userID1, resp.Likers[0].UserId)
		require.Equal(t, "alice", resp.Likers[0].Username)
		require.Equal(t, userID2, resp.Likers[1].UserId)
		require.Empty(t, resp.Likers[1].Username) // Missing username should be empty
	})

	t.Run("user lookup error handled gracefully", func(t *testing.T) {
		ctx := context.Background()
		postID := "post-123"
		userID1 := "user-1"

		mockLikes.EXPECT().
			GetLikers(ctx, postID, "", 20).
			Return([]string{userID1}, "", nil).
			Times(1)

		// User lookup fails but service continues
		mockUsers.EXPECT().
			GetByIDs(ctx, []string{userID1}).
			Return(nil, errors.New("database error")).
			Times(1)

		resp, err := service.GetLikers(ctx, &likesapi.GetLikersRequest{
			PostId: postID,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Likers, 1)
		require.Equal(t, userID1, resp.Likers[0].UserId)
		require.Empty(t, resp.Likers[0].Username) // Username empty due to error
	})

	t.Run("likes adapter error", func(t *testing.T) {
		ctx := context.Background()
		postID := "post-123"

		mockLikes.EXPECT().
			GetLikers(ctx, postID, "", 20).
			Return(nil, "", errors.New("database error")).
			Times(1)

		_, err := service.GetLikers(ctx, &likesapi.GetLikersRequest{
			PostId: postID,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get likers")
	})
}
