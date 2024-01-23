package pkg

import (
	"fmt"
	"github.com/materkov/meme9/api/src/store2"
)

func GetFeedPostIds(userID int) ([]int, error) {
	userIds, err := store2.GlobalStore.Subs.GetFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting edges: %w", err)
	}

	userIds = append(userIds, userID)

	postIds, err := store2.GlobalStore.Wall.Get(userIds, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("error getting post ids: %w", err)
	}

	if len(postIds) > 100 {
		postIds = postIds[:100]
	}

	return postIds, err
}
