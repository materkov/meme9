package pkg

import (
	"errors"
	"fmt"
	"github.com/materkov/meme9/web6/src/store"
	"time"
)

func ParseAuthToken(tokenStr string) *store.Token {
	if tokenStr == "" {
		return nil
	}

	tokenID, err := store.GetUnique(store.UniqueTypeAuthToken, tokenStr)
	if errors.Is(err, store.ErrUniqueNotFound) {
		return nil
	} else if err != nil {
		LogErr(err)
		return nil
	}

	token, err := store.GetToken(tokenID)
	if err != nil {
		LogErr(err)
		return nil
	}

	return token
}

func GenerateAuthToken(userID int) (string, error) {
	token := store.Token{
		ID:     0,
		UserID: userID,
		Date:   int(time.Now().Unix()),
		Token:  store.GenerateToken(userID),
	}

	err := store.AddToken(&token)
	if err != nil {
		return "", fmt.Errorf("error storing token: %w", err)
	}

	err = store.AddUnique(store.UniqueTypeAuthToken, token.Token, token.ID)
	if err != nil {
		return "", fmt.Errorf("error storing edge: %w", err)
	}

	return token.Token, nil
}
