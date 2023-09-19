package pkg

import (
	"fmt"
	"github.com/materkov/meme9/web6/src/store"
	"time"
)

func ParseAuthToken(tokenStr string) *store.Token {
	if tokenStr == "" {
		return nil
	}

	tokenID, err := store.GetEdgeByUniqueKey(store.FakeObjToken, 0, tokenStr)
	if err != nil {
		return nil
	} else if tokenID == 0 {
		return nil
	}

	token, err := store.GetToken(tokenID)
	if err != nil {
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

	err = store.AddEdge(store.FakeObjToken, token.ID, 0, token.Token)
	if err != nil {
		return "", fmt.Errorf("error storing edge: %w", err)
	}

	return token.Token, nil
}
