package api

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/pkg/utils"
	"github.com/materkov/meme9/web5/store"
	"log"
	"net/url"
	"sort"
	"strconv"
	"sync"
)

const (
	feedTypeDiscover = "DISCOVER"
	feedTypeFeed     = "FEED"
)

// /feed
func handleFeed(ctx context.Context, viewerID int, reqUrl string) []interface{} {
	parsedURL, _ := url.Parse(reqUrl)
	cursor, _ := strconv.Atoi(parsedURL.Query().Get("cursor"))

	feedType := parsedURL.Query().Get("feedType")
	if feedType != feedTypeDiscover && feedType != feedTypeFeed {
		feedType = feedTypeDiscover
	}
	if feedType == feedTypeFeed && viewerID == 0 {
		feedType = feedTypeDiscover
	}

	count := 10

	var postIds []int
	if feedType == feedTypeDiscover {
		postIdsStr, err := store.RedisClient.LRange(context.Background(), "feed", 0, 10000).Result()
		if err != nil {
			log.Printf("Error getting feed: %s", err)
		}

		postIds = utils.StrToIntArray(postIdsStr)
	} else if feedType == feedTypeFeed {
		key := fmt.Sprintf("following:%d", viewerID)
		followingIdsStr, err := store.RedisClient.ZRange(ctx, key, 0, -1).Result()
		if err != nil {
			log.Printf("Error getting following ids: %s", err)
		}

		pipe := store.RedisClient.Pipeline()
		var cmds []*redis.StringSliceCmd

		for _, userIDStr := range followingIdsStr {
			userID, _ := strconv.Atoi(userIDStr)
			if userID > 0 {
				cmd := pipe.LRange(context.Background(), fmt.Sprintf("feed:%d", userID), 0, int64(count))
				cmds = append(cmds, cmd)
			}
		}

		cmd := pipe.LRange(context.Background(), fmt.Sprintf("feed:%d", viewerID), 0, int64(count))
		cmds = append(cmds, cmd)

		_, err = pipe.Exec(ctx)
		if err != nil {
			log.Printf("Errog etting feeds: %s", err)
		}

		for _, cmd := range cmds {
			currentIds := utils.StrToIntArray(cmd.Val())
			postIds = append(postIds, currentIds...)
		}

		sort.Sort(sort.Reverse(sort.IntSlice(postIds)))
	}

	if cursor > 0 && cursor < len(postIds) {
		postIds = postIds[cursor:]
	}
	if count < len(postIds) {
		postIds = postIds[:count]
	}

	nextCursor := ""
	if len(postIds) == count {
		nextCursor = strconv.Itoa(cursor + count)
	}

	feed := Edges{
		URL:        reqUrl,
		TotalCount: 20,
		NextCursor: nextCursor,
		Items:      utils.IntToStrArray(postIds),
	}

	var results []interface{}
	results = append(results, feed)

	for _, postID := range postIds {
		results = append(results, fmt.Sprintf("/posts/%d", postID))
	}

	var userIds []int
	store.CachedStoreFromCtx(ctx).Post.Preload(postIds)
	for _, postID := range postIds {
		post := store.CachedStoreFromCtx(ctx).Post.Get(postID)
		if post != nil {
			userIds = append(userIds, post.UserID)
			store.CachedStoreFromCtx(ctx).Online.Preload(post.UserID)
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		store.CachedStoreFromCtx(ctx).User.Preload(userIds)
		wg.Done()
	}()
	go func() {
		store.CachedStoreFromCtx(ctx).Liked.Preload(viewerID, postIds)
		wg.Done()
	}()
	wg.Wait()

	return results
}
