package api

import (
	"context"
	"github.com/materkov/meme9/web5/store"
	"log"
	"net/url"
	"strconv"
)

// /feed
func handleFeed(viewerID int, reqUrl string) []interface{} {
	parsedURL, _ := url.Parse(reqUrl)
	cursor, _ := strconv.Atoi(parsedURL.Query().Get("cursor"))
	count := 10

	postIds, err := store.RedisClient.LRange(context.Background(), "feed", int64(cursor), int64(cursor+count-1)).Result()
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	nextCursor := ""
	if len(postIds) == count {
		nextCursor = strconv.Itoa(cursor + count)
	}

	feed := Edges{
		URL:        reqUrl,
		TotalCount: 20,
		NextCursor: nextCursor,
		Items:      postIds,
	}

	var results []interface{}
	results = append(results, feed)

	for _, postID := range postIds {
		results = append(results, handlePostsId(viewerID, "/posts/"+postID)...)
	}

	return results
}
