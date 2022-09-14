package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web5/store"
	"log"
	"strconv"
	"time"
)

func postsList(ids []int, viewerID int) []*Post {
	if len(ids) == 0 {
		return nil
	}

	keys := make([]string, len(ids))
	for i, postID := range ids {
		keys[i] = fmt.Sprintf("node:%d", postID)
	}

	postsBytes, err := store.RedisClient.MGet(context.Background(), keys...).Result()
	if err != nil {
		log.Printf("error getting posts: %s", err)
	}

	posts := map[int]*store.Post{}
	for _, postBytes := range postsBytes {
		if postBytes == nil {
			continue
		}

		post := &store.Post{}
		err = json.Unmarshal([]byte(postBytes.(string)), post)
		if err != nil {
			continue
		}

		posts[post.ID] = post
	}

	apiPosts := make([]*Post, len(ids))
	for i, postID := range ids {
		result := &Post{ID: strconv.Itoa(postID)}
		apiPosts[i] = result

		post, ok := posts[postID]
		if !ok {
			continue
		} else if post.IsDeleted {
			result.Text = "DELETED"
			continue
		}

		result.Text = post.Text
		result.Date = time.Unix(int64(post.Date), 0).UTC().Format(time.RFC3339)
		result.UserID = strconv.Itoa(post.UserID)
		result.CanDelete = post.UserID == viewerID
	}

	return apiPosts
}

func postsAdd(text string, userID int) (int, error) {
	nextId := int(time.Now().UnixMilli())

	post := store.Post{
		ID:     nextId,
		Text:   text,
		UserID: userID,
		Date:   int(time.Now().Unix()),
	}
	err := store.NodeSave(post.ID, post)
	if err != nil {
		return 0, fmt.Errorf("error creating post node: %w", err)
	}

	doneFeed := make(chan bool)
	doneUserFeed := make(chan bool)
	go func() {
		_, err = store.RedisClient.LPush(context.Background(), "feed", post.ID).Result()
		if err != nil {
			log.Printf("Error saving post to feed")
		}
		doneFeed <- true
	}()
	go func() {
		_, err = store.RedisClient.LPush(context.Background(), fmt.Sprintf("feed:%d", post.UserID), post.ID).Result()
		if err != nil {
			log.Printf("Error saving user feed key: %s", err)
		}
		doneUserFeed <- true
	}()

	<-doneFeed
	<-doneUserFeed

	return post.ID, nil
}

func postsDelete(post *store.Post) error {
	pipe := store.RedisClient.Pipeline()
	pipe.LRem(context.Background(), "feed", 0, post.ID)
	pipe.LRem(context.Background(), fmt.Sprintf("feed:%d", post.UserID), 0, post.ID)

	_, err := pipe.Exec(context.Background())
	if err != nil {
		return fmt.Errorf("error removing from feed: %w", err)
	}

	post.IsDeleted = true
	err = store.NodeSave(post.ID, post)
	if err != nil {
		return fmt.Errorf("error updating post node: %w", err)
	}

	return nil
}
