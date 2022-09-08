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

func usersList(users []*User) {
	usersMap := map[string]store.User{}
	for _, user := range users {
		userID, _ := strconv.Atoi(user.ID)
		obj := store.User{}
		err := getObject(userID, &obj)
		if err == nil {
			usersMap[user.ID] = obj
		}
	}

	for _, user := range users {
		user.Href = fmt.Sprintf("/users/%s", user.ID)

		stUser, ok := usersMap[user.ID]
		if !ok {
			user.Name = fmt.Sprintf("User #%s", user.ID)
			continue
		}

		user.Name = stUser.Name
	}
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
