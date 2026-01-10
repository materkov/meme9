package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/twitchtv/twirp"
	"go.uber.org/mock/gomock"

	authapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/auth"
	"github.com/materkov/meme9/auth-service/adapters/mongo"
	"github.com/materkov/meme9/auth-service/api/mocks"
	"golang.org/x/crypto/bcrypt"
)

func reauireError(t *testing.T, err error, msg string) {
	require.Error(t, err)
	twirpErr := err.(twirp.Error)
	require.Equal(t, msg, twirpErr.Msg())

}

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mocks.NewMockMongoAdapter(ctrl)
	service := NewService(mockAdapter)

	require.NotNil(t, service)
	require.Equal(t, mockAdapter, service.mongo)
}

func TestService_Login(t *testing.T) {
	ctx := context.Background()

	t.Run("successful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		password := "testpassword"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		mockAdapter.EXPECT().
			GetUserByUsername(gomock.Any(), "testuser").
			Return(&mongo.User{
				ID:           "user123",
				Username:     "testuser",
				PasswordHash: string(hashedPassword),
				CreatedAt:    time.Now(),
			}, nil)

		mockAdapter.EXPECT().
			CreateToken(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, token mongo.Token) (string, error) {
				require.NotEmpty(t, token.Token)
				require.Equal(t, "user123", token.UserID)
				return "token-id", nil
			})

		req := &authapi.LoginRequest{
			Username: "testuser",
			Password: password,
		}

		resp, err := service.Login(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "user123", resp.UserId)
		require.Equal(t, "testuser", resp.Username)
		require.NotEmpty(t, resp.Token)
	})

	t.Run("empty username", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &authapi.LoginRequest{
			Username: "",
			Password: "password",
		}

		_, err := service.Login(ctx, req)
		reauireError(t, err, "invalid_credentials")
	})

	t.Run("empty password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &authapi.LoginRequest{
			Username: "testuser",
			Password: "",
		}

		_, err := service.Login(ctx, req)
		reauireError(t, err, "invalid_credentials")
	})

	t.Run("user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetUserByUsername(gomock.Any(), "nonexistent").
			Return(nil, mongo.ErrUserNotFound)

		req := &authapi.LoginRequest{
			Username: "nonexistent",
			Password: "password",
		}

		_, err := service.Login(ctx, req)
		reauireError(t, err, "invalid_credentials")
	})

	t.Run("wrong password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		password := "correctpassword"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		mockAdapter.EXPECT().
			GetUserByUsername(gomock.Any(), "testuser").
			Return(&mongo.User{
				ID:           "user123",
				Username:     "testuser",
				PasswordHash: string(hashedPassword),
				CreatedAt:    time.Now(),
			}, nil)

		req := &authapi.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}

		_, err = service.Login(ctx, req)
		reauireError(t, err, "invalid_credentials")
	})

	t.Run("database error getting user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetUserByUsername(gomock.Any(), "testuser").
			Return(nil, errors.New("database error"))

		req := &authapi.LoginRequest{
			Username: "testuser",
			Password: "password",
		}

		resp, err := service.Login(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "error getting user")
	})

	t.Run("error creating token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		password := "testpassword"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		user := &mongo.User{
			ID:           "user123",
			Username:     "testuser",
			PasswordHash: string(hashedPassword),
			CreatedAt:    time.Now(),
		}

		mockAdapter.EXPECT().
			GetUserByUsername(gomock.Any(), "testuser").
			Return(user, nil)

		mockAdapter.EXPECT().
			CreateToken(gomock.Any(), gomock.Any()).
			Return("", errors.New("token creation error"))

		req := &authapi.LoginRequest{
			Username: "testuser",
			Password: password,
		}

		resp, err := service.Login(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "error creating auth token")
	})
}

func TestService_Register(t *testing.T) {
	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetUserByUsername(gomock.Any(), "newuser").
			Return(nil, mongo.ErrUserNotFound)

		mockAdapter.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, user mongo.User) (string, error) {
				require.Equal(t, "newuser", user.Username)
				require.NotEmpty(t, user.PasswordHash)
				require.NotZero(t, user.CreatedAt)
				// Verify password is hashed
				err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123"))
				require.NoError(t, err)
				return "user456", nil
			})

		mockAdapter.EXPECT().
			CreateToken(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, token mongo.Token) (string, error) {
				require.Equal(t, "user456", token.UserID)
				require.NotEmpty(t, token.Token)
				return "token-id", nil
			})

		req := &authapi.RegisterRequest{
			Username: "newuser",
			Password: "password123",
		}

		resp, err := service.Register(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "user456", resp.UserId)
		require.Equal(t, "newuser", resp.Username)
		require.NotEmpty(t, resp.Token)
	})

	t.Run("empty username", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &authapi.RegisterRequest{
			Username: "",
			Password: "password",
		}

		_, err := service.Register(ctx, req)
		reauireError(t, err, "username_required")
	})

	t.Run("empty password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &authapi.RegisterRequest{
			Username: "testuser",
			Password: "",
		}

		_, err := service.Register(ctx, req)
		reauireError(t, err, "password_required")
	})

	t.Run("username already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		existingUser := &mongo.User{
			ID:           "existing123",
			Username:     "existinguser",
			PasswordHash: "hash",
			CreatedAt:    time.Now(),
		}

		mockAdapter.EXPECT().
			GetUserByUsername(gomock.Any(), "existinguser").
			Return(existingUser, nil)

		req := &authapi.RegisterRequest{
			Username: "existinguser",
			Password: "password",
		}

		_, err := service.Register(ctx, req)
		reauireError(t, err, "username_exists")
	})

	t.Run("database error getting user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetUserByUsername(gomock.Any(), "testuser").
			Return(nil, errors.New("database error"))

		req := &authapi.RegisterRequest{
			Username: "testuser",
			Password: "password",
		}

		resp, err := service.Register(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "error getting user")
	})

	t.Run("error creating user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetUserByUsername(gomock.Any(), "newuser").
			Return(nil, mongo.ErrUserNotFound)

		mockAdapter.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return("", errors.New("user creation error"))

		req := &authapi.RegisterRequest{
			Username: "newuser",
			Password: "password",
		}

		resp, err := service.Register(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "error creating user")
	})

	t.Run("error creating token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetUserByUsername(gomock.Any(), "newuser").
			Return(nil, mongo.ErrUserNotFound)

		mockAdapter.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return("user789", nil)

		mockAdapter.EXPECT().
			CreateToken(gomock.Any(), gomock.Any()).
			Return("", errors.New("token creation error"))

		req := &authapi.RegisterRequest{
			Username: "newuser",
			Password: "password",
		}

		resp, err := service.Register(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "error creating auth token")
	})
}

func TestService_VerifyToken(t *testing.T) {
	ctx := context.Background()

	t.Run("successful token verification", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		token := &mongo.Token{
			ID:        "token123",
			Token:     "valid-token",
			UserID:    "user123",
			CreatedAt: time.Now(),
		}

		mockAdapter.EXPECT().
			GetTokenByValue(gomock.Any(), "valid-token").
			Return(token, nil)

		req := &authapi.VerifyTokenRequest{
			Token: "valid-token",
		}

		resp, err := service.VerifyToken(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "user123", resp.UserId)
	})

	t.Run("token with Bearer prefix", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		token := &mongo.Token{
			ID:        "token123",
			Token:     "valid-token",
			UserID:    "user123",
			CreatedAt: time.Now(),
		}

		mockAdapter.EXPECT().
			GetTokenByValue(gomock.Any(), "valid-token").
			Return(token, nil)

		req := &authapi.VerifyTokenRequest{
			Token: "Bearer valid-token",
		}

		resp, err := service.VerifyToken(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "user123", resp.UserId)
	})

	t.Run("token with Bearer prefix and spaces", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		token := &mongo.Token{
			ID:        "token123",
			Token:     "valid-token",
			UserID:    "user123",
			CreatedAt: time.Now(),
		}

		mockAdapter.EXPECT().
			GetTokenByValue(gomock.Any(), "valid-token").
			Return(token, nil)

		req := &authapi.VerifyTokenRequest{
			Token: "Bearer  valid-token  ",
		}

		resp, err := service.VerifyToken(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "user123", resp.UserId)
	})

	t.Run("empty token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &authapi.VerifyTokenRequest{
			Token: "",
		}

		_, err := service.VerifyToken(ctx, req)
		reauireError(t, err, "invalid_token")
	})

	t.Run("token not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetTokenByValue(gomock.Any(), "invalid-token").
			Return(nil, mongo.ErrTokenNotFound)

		req := &authapi.VerifyTokenRequest{
			Token: "invalid-token",
		}

		_, err := service.VerifyToken(ctx, req)
		reauireError(t, err, "invalid_token")
	})

	t.Run("database error getting token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetTokenByValue(gomock.Any(), "token").
			Return(nil, errors.New("database error"))

		req := &authapi.VerifyTokenRequest{
			Token: "token",
		}

		resp, err := service.VerifyToken(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "error verifying token")
	})
}

func TestGenerateToken(t *testing.T) {
	t.Run("generates unique tokens", func(t *testing.T) {
		token1, err := generateToken()
		require.NoError(t, err)
		require.NotEmpty(t, token1)

		token2, err := generateToken()
		require.NoError(t, err)
		require.NotEmpty(t, token2)

		require.NotEqual(t, token1, token2)
	})

	t.Run("generates valid base64 tokens", func(t *testing.T) {
		token, err := generateToken()
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.Greater(t, len(token), 20) // Base64 encoded 32 bytes should be ~43 chars
	})
}
