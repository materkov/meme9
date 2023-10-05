package pkg

import (
	"fmt"
	"github.com/materkov/meme9/web6/src/store"
	"sort"
)

func GetFeedPostIds(userID int) ([]int, error) {
	edges, err := store.GetEdges(userID, store.EdgeTypeFollowing)
	if err != nil {
		return nil, fmt.Errorf("error getting edges: %w", err)
	}

	userIds := store.GetToId(edges)
	userIds = append(userIds, userID)

	allPostsCh := make(chan []store.Edge)
	for _, userId := range userIds {
		userIdCopy := userId
		go func() {
			posts, err := store.GetEdges(userIdCopy, store.EdgeTypePosted)

			LogErr(err)
			allPostsCh <- posts
		}()
	}

	var allPosts []store.Edge
	for range userIds {
		allPosts = append(allPosts, <-allPostsCh...)
	}

	sort.Slice(allPosts, func(i, j int) bool {
		return allPosts[i].Date > allPosts[j].Date
	})

	if len(allPosts) > 100 {
		allPosts = allPosts[:100]
	}

	return store.GetToId(allPosts), err
}
