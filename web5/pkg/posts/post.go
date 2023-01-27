package posts

import (
	"context"
	"fmt"
	"github.com/materkov/meme9/web5/store"
	"log"
	"time"
)

func CanSee(post *store.Post, viewerID int) bool {
	if post == nil {
		return false
	}
	if post.IsDeleted {
		return false
	}

	return true
}

func CanEdit(post *store.Post, viewerID int) bool {
	if !CanSee(post, viewerID) {
		return false
	}

	return post.UserID == viewerID
}

func Add(text string, userID int, photoID int) (int, error) {
	nextId := int(time.Now().UnixMilli())

	post := store.Post{
		ID:      nextId,
		Text:    text,
		UserID:  userID,
		Date:    int(time.Now().Unix()),
		PhotoID: photoID,
	}
	err := store.NodeSave(post.ID, store.ObjectTypePost, post)
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

func Delete(post *store.Post) error {
	pipe := store.RedisClient.Pipeline()
	pipe.LRem(context.Background(), "feed", 0, post.ID)
	pipe.LRem(context.Background(), fmt.Sprintf("feed:%d", post.UserID), 0, post.ID)

	_, err := pipe.Exec(context.Background())
	if err != nil {
		return fmt.Errorf("error removing from feed: %w", err)
	}

	post.IsDeleted = true
	err = store.NodeSave(post.ID, store.ObjectTypePost, post)
	if err != nil {
		return fmt.Errorf("error updating post node: %w", err)
	}

	return nil
}
