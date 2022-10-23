package auth

import (
	"github.com/materkov/meme9/web5/pkg/testutils"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestAuthEmailAuth(t *testing.T) {
	testutils.SetupRedis(t)

	// Register
	registeredUserID, err := Register("my@email.com", "pass")
	require.NoError(t, err)
	require.NotEmpty(t, registeredUserID)

	// Login, good password
	userID, err := EmailAuth("my@email.com", "pass")
	require.NoError(t, err)
	require.Equal(t, registeredUserID, userID)

	// Bad password
	_, err = EmailAuth("my@email.com", "bad-password")
	require.Equal(t, ErrInvalidCredentials, err)

	// Bad email
	_, err = EmailAuth("bad@email.com", "pass")
	require.Equal(t, ErrInvalidCredentials, err)
}

func TestAuthCheckCredentials(t *testing.T) {
	testutils.SetupRedis(t)

	_, err := Register("good@mail.com", "")
	require.NoError(t, err)

	table := []struct {
		Email, Password, Err string
	}{
		{"", "", "empty email"},
		{"bad!email", "", "incorrect email"},
		{"e@mail.com", "", "empty password"},
		{"e@mail.com" + strings.Repeat("m", 400), "", "email too long"},
		{"good@mail.com", "", "email already registered"},
	}

	for _, item := range table {
		t.Run(item.Err, func(t *testing.T) {
			err := ValidateCredentials(item.Email, item.Password)
			require.Equal(t, item.Err, err)
		})
	}
}
