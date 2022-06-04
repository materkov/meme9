package pkg

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAuthToken(t *testing.T) {
	token := AuthToken{
		UserID:   48142,
		IssuedAt: int(time.Now().Unix()),
	}
	tokenStr := token.ToString()
	require.NotEmpty(t, tokenStr)

	parsedToken, err := ParseAuthToken(tokenStr)
	require.NoError(t, err)
	require.Equal(t, token.UserID, parsedToken.UserID)
	require.Equal(t, token.IssuedAt, parsedToken.IssuedAt)
}

func TestParseAuthToken_Errors(t *testing.T) {
	_, err := ParseAuthToken("incorrect token")
	require.Equal(t, err, ErrIncorrectToken)

	token := AuthToken{
		UserID:   48142,
		IssuedAt: int(time.Now().Unix()),
	}
	tokenStr := token.ToString()

	// incorrect signature
	_, err = ParseAuthToken(tokenStr + "x")
	require.Equal(t, ErrIncorrectToken, err)

	_, err = ParseAuthToken("x" + tokenStr)
	require.Equal(t, ErrIncorrectToken, err)
}

func TestParseAuthToken_VeryOld(t *testing.T) {
	token := AuthToken{
		UserID:   1,
		IssuedAt: 12334,
	}
	tokenStr := token.ToString()

	_, err := ParseAuthToken(tokenStr)
	require.Equal(t, ErrIncorrectToken, err)
}
