package types

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web4/store"
	"log"
	"sort"
	"strconv"
	"time"
)

type Post struct {
	ID     string `json:"id,omitempty"`
	FromID string `json:"fromId,omitempty"`

	Text       string `json:"text"`
	DetailsURL string `json:"detailsURL"`
}

func nextID() int {
	return int(time.Now().UnixMilli())
}

func postsList(ids []int) []*Post {
	postsMap := map[int]store.Post{}
	for _, postID := range ids {
		post := store.Post{}
		err := getObject(postID, &post)
		if err == nil {
			postsMap[post.ID] = post
		}
	}

	results := make([]*Post, len(ids))
	for i, postID := range ids {
		result := &Post{
			ID:         strconv.Itoa(postID),
			DetailsURL: fmt.Sprintf("/posts/%d", postID),
		}
		results[i] = result

		post, ok := postsMap[postID]
		if !ok {
			continue
		}

		result.FromID = strconv.Itoa(post.UserID)
		result.Text = post.Text
	}

	return results
}

func postsAdd(req *postsAddRequest, viewerID int) (int, error) {
	postID := nextID()

	post := store.Post{
		ID:     postID,
		UserID: viewerID,
		Text:   req.Text,
	}
	postBytes, err := json.Marshal(post)
	if err != nil {
		return 0, fmt.Errorf("error serializing post to json: %w", err)
	}

	_, err = redisClient.Set(context.Background(), fmt.Sprintf("node:%d", postID), postBytes, 0).Result()
	if err != nil {
		return 0, fmt.Errorf("error saving post to redis: %w", err)
	}

	_, err = redisClient.LPush(context.Background(), "feed", post.ID).Result()
	if err != nil {
		log.Printf("Error saving feed key: %s", err)
	}

	_, err = redisClient.LPush(context.Background(), fmt.Sprintf("feed:%d", post.UserID), post.ID).Result()
	if err != nil {
		log.Printf("Error saving user feed key: %s", err)
	}

	return postID, nil
}

func postsGetFeed() ([]int, error) {
	postIdsStr, err := redisClient.LRange(context.Background(), "feed", 0, 10).Result()
	if err != nil {
		return nil, fmt.Errorf("error reading feed: %w", err)
	}

	postIds := make([]int, len(postIdsStr))
	for i, postID := range postIdsStr {
		postIds[i], _ = strconv.Atoi(postID)
	}

	return postIds, nil
}

func postsGetFeedByUsers(users []int) ([]int, error) {
	p := redisClient.Pipeline()

	results := make([]*redis.StringSliceCmd, len(users))
	for i, userID := range users {
		results[i] = p.LRange(context.Background(), fmt.Sprintf("feed:%d", userID), 0, 10)
	}

	_, err := p.Exec(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error reading feed keys: %w", err)
	}

	var postIds []int
	for _, resultItem := range results {
		for _, postIDStr := range resultItem.Val() {
			postID, _ := strconv.Atoi(postIDStr)
			postIds = append(postIds, postID)
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(postIds)))

	return postIds, nil
}
