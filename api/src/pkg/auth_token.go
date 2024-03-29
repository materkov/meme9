package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/materkov/meme9/api/src/pkg/tracer"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"time"
)

func ParseAuthToken(ctx context.Context, tokenStr string) *store.Token {
	defer tracer.FromCtx(ctx).StartChild("ParseAuthToken").Stop()

	if tokenStr == "" {
		return nil
	}

	tokenID, err := store2.GlobalStore.Unique.Get(store2.UniqueTypeAuthToken, tokenStr)
	if errors.Is(err, store2.ErrNotFound) {
		return nil
	} else if err != nil {
		LogErr(err)
		return nil
	}

	tokens, err := store2.GlobalStore.Tokens.Get([]int{tokenID})
	if err != nil {
		LogErr(err)
		return nil
	} else if tokens[tokenID] == nil {
		return nil
	}

	return tokens[tokenID]
}

func GenerateAuthToken(userID int) (string, error) {
	token := store.Token{
		ID:     0,
		UserID: userID,
		Date:   int(time.Now().Unix()),
		Token:  store.GenerateToken(userID),
	}

	err := store2.GlobalStore.Tokens.Add(&token)
	if err != nil {
		return "", fmt.Errorf("error storing token: %w", err)
	}

	err = store2.GlobalStore.Unique.Add(store2.UniqueTypeAuthToken, token.Token, token.ID)
	if err != nil {
		return "", fmt.Errorf("error storing edge: %w", err)
	}

	return token.Token, nil
}
