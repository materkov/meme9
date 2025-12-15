package auth

import (
	"context"
	"testing"

	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	authapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/auth"
	"github.com/stretchr/testify/require"
	"github.com/twitchtv/twirp"
)

// Mock implementations for testing
type mockUsersAdapter struct {
	users map[string]*users.User
}

func newMockUsersAdapter() *mockUsersAdapter {
	return &mockUsersAdapter{
		users: make(map[string]*users.User),
	}
}

func (m *mockUsersAdapter) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, users.ErrNotFound
}

func (m *mockUsersAdapter) Create(ctx context.Context, user users.User) (string, error) {
	userID := "user-" + user.Username
	user.ID = userID
	m.users[userID] = &user
	return userID, nil
}

type mockTokensAdapter struct {
	tokens map[string]*tokens.Token
}

func newMockTokensAdapter() *mockTokensAdapter {
	return &mockTokensAdapter{
		tokens: make(map[string]*tokens.Token),
	}
}

func (m *mockTokensAdapter) Create(ctx context.Context, token tokens.Token) (string, error) {
	tokenID := "token-" + token.Token
	token.ID = tokenID
	m.tokens[token.Token] = &token
	return tokenID, nil
}

func (m *mockTokensAdapter) GetByValue(ctx context.Context, tokenValue string) (*tokens.Token, error) {
	token, ok := m.tokens[tokenValue]
	if !ok {
		return nil, tokens.ErrNotFound
	}
	return token, nil
}

func requireApiError(t *testing.T, err error, msg string) {
	t.Helper()

	twirpErr, ok := err.(twirp.Error)
	require.True(t, ok, "error should be a twirp.Error, got: %T", err)
	require.Equal(t, msg, twirpErr.Msg())
}

func initService() (*Service, *mockUsersAdapter) {
	mockUsers := newMockUsersAdapter()
	mockTokens := newMockTokensAdapter()

	return NewService(mockUsers, mockTokens, nil), mockUsers
}

func TestService_RegisterAndLogin(t *testing.T) {
	service, _ := initService()

	// Register
	respRegister, err := service.Register(context.Background(), &authapi.RegisterRequest{
		Username: "test-user",
		Password: "test-password",
	})
	require.NoError(t, err)
	require.NotEmpty(t, respRegister.UserId)
	require.NotEmpty(t, respRegister.Token)
	require.Equal(t, "test-user", respRegister.Username)

	// Login
	respLogin, err := service.Login(context.Background(), &authapi.LoginRequest{
		Username: "test-user",
		Password: "test-password",
	})
	require.NoError(t, err)
	require.Equal(t, respRegister.UserId, respLogin.UserId)
	require.Equal(t, "test-user", respLogin.Username)
	require.NotEmpty(t, respLogin.Token)

	// Verify Register token
	respVerify, err := service.VerifyToken(context.Background(), &authapi.VerifyTokenRequest{
		Token: respRegister.Token,
	})
	require.NoError(t, err)
	require.Equal(t, respRegister.UserId, respVerify.UserId)

	// Verify Login token
	respVerify, err = service.VerifyToken(context.Background(), &authapi.VerifyTokenRequest{
		Token: respLogin.Token,
	})
	require.NoError(t, err)
	require.Equal(t, respRegister.UserId, respVerify.UserId)
}

func TestService_Login_Invalid(t *testing.T) {
	service, mockUsers := initService()
	mockUsers.users["test-user"] = &users.User{
		Username:     "test-user",
		PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password"
	}

	t.Run("empty username", func(t *testing.T) {
		_, err := service.Login(context.Background(), &authapi.LoginRequest{
			Username: "",
			Password: "test-password",
		})
		requireApiError(t, err, "invalid_credentials")
	})

	t.Run("empty password", func(t *testing.T) {
		_, err := service.Login(context.Background(), &authapi.LoginRequest{
			Username: "test-user",
			Password: "",
		})
		requireApiError(t, err, "invalid_credentials")
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := service.Login(context.Background(), &authapi.LoginRequest{
			Username: "non-existing-user",
			Password: "test-password",
		})
		requireApiError(t, err, "invalid_credentials")
	})

	t.Run("invalid password", func(t *testing.T) {
		_, err := service.Login(context.Background(), &authapi.LoginRequest{
			Username: "test-user",
			Password: "wrong-password",
		})
		requireApiError(t, err, "invalid_credentials")
	})
}

func TestService_Register_Invalid(t *testing.T) {
	service, mockUsers := initService()
	mockUsers.users["test-user"] = &users.User{
		Username: "test-user",
	}

	t.Run("empty username", func(t *testing.T) {
		_, err := service.Register(context.Background(), &authapi.RegisterRequest{
			Username: "",
			Password: "test-password",
		})
		requireApiError(t, err, "username_required")
	})

	t.Run("empty password", func(t *testing.T) {
		_, err := service.Register(context.Background(), &authapi.RegisterRequest{
			Username: "test-user",
			Password: "",
		})
		requireApiError(t, err, "password_required")
	})

	t.Run("username exists", func(t *testing.T) {
		_, err := service.Register(context.Background(), &authapi.RegisterRequest{
			Username: "test-user",
			Password: "test-password",
		})
		requireApiError(t, err, "username_exists")
	})
}

func TestService_VerifyToken_Invalid(t *testing.T) {
	service, _ := initService()

	t.Run("empty token", func(t *testing.T) {
		_, err := service.VerifyToken(context.Background(), &authapi.VerifyTokenRequest{
			Token: "",
		})
		requireApiError(t, err, "invalid_token")
	})

	t.Run("token not found", func(t *testing.T) {
		_, err := service.VerifyToken(context.Background(), &authapi.VerifyTokenRequest{
			Token: "nonexistent-token",
		})
		requireApiError(t, err, "invalid_token")
	})
}

func TestService_VerifyToken_WithBearerPrefix(t *testing.T) {
	service, _ := initService()

	resp, err := service.Register(context.Background(), &authapi.RegisterRequest{
		Username: "test-user",
		Password: "test-password",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)

	validTokens := []string{
		"Bearer " + resp.Token,
		"Bearer " + resp.Token + " ",
		"Bearer    " + resp.Token + " ",
		" " + resp.Token,
		" " + resp.Token + "  ",
	}

	for _, token := range validTokens {
		_, err = service.VerifyToken(context.Background(), &authapi.VerifyTokenRequest{
			Token: token,
		})
		require.NoError(t, err)
	}
}
