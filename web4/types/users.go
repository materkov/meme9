package types

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web4/store"
	"strconv"
)

type User struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Href string `json:"href,omitempty"`
}

func usersList(ids []int) []*User {
	usersMap := map[int]store.User{}
	for _, postID := range ids {
		obj := store.User{}
		err := getObject(postID, &obj)
		if err == nil {
			usersMap[obj.ID] = obj
		}
	}

	results := make([]*User, len(ids))
	for i, userID := range ids {
		result := &User{
			ID:   strconv.Itoa(userID),
			Href: fmt.Sprintf("/users/%d", userID),
		}
		results[i] = result

		user, ok := usersMap[userID]
		if !ok {
			continue
		}

		result.Name = user.Name
	}

	return results
}

func usersAdd(vkID int) (int, error) {
	id := nextID()

	user := store.User{
		ID:   id,
		VkID: vkID,
		Name: fmt.Sprintf("User #%d", vkID),
	}

	err := saveObject(id, user)
	if err != nil {
		return 0, fmt.Errorf("error saving user: %w", err)
	}

	return id, nil
}

func usersGetOrCreateByVKID(vkID int) (int, error) {
	mapKey := fmt.Sprintf("vk2user_map:%d", vkID)
	userIDStr, err := redisClient.Get(context.Background(), mapKey).Result()
	if err == redis.Nil {
		userID, err := usersAdd(vkID)
		if err != nil {
			return 0, err
		}

		_, err = redisClient.Set(context.Background(), mapKey, userID, 0).Result()
		if err != nil {
			return 0, err
		}

		return userID, nil
	} else if err != nil {
		return 0, err
	}

	userID, _ := strconv.Atoi(userIDStr)
	return userID, nil
}
