package posts

import (
	"context"
	"fmt"
	"github.com/materkov/meme9/web5/store"
)

func FeedLen(userID int) (int, error) {
	res, err := store.RedisClient.LLen(context.Background(), fmt.Sprintf("feed:%d", userID)).Result()
	if err != nil {
		return 0, fmt.Errorf("error getting feed len: %w", err)
	}

	return int(res), nil
}
