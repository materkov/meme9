package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/tracer"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
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

	id, err := store2.GlobalStore.Nodes.Add(store.ObjTypeToken, token)
	if err != nil {
		return "", fmt.Errorf("error storing token: %w", err)
	}
	token.ID = id

	err = store2.GlobalStore.Unique.Add(store2.UniqueTypeAuthToken, token.Token, token.ID)
	if err != nil {
		return "", fmt.Errorf("error storing edge: %w", err)
	}

	return token.Token, nil
}
