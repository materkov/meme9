package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/store"
	"log"
	"strconv"
	"time"
)

func usersGetOrCreateByVKID(vkID int) (int, error) {
	mapKey := fmt.Sprintf("vk2user_map:%d", vkID)
	userIDStr, err := store.RedisClient.Get(context.Background(), mapKey).Result()
	if err == redis.Nil {
		userID, err := usersAdd(vkID)
		if err != nil {
			return 0, err
		}

		_, err = store.RedisClient.Set(context.Background(), mapKey, userID, 0).Result()
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

func usersAdd(vkID int) (int, error) {
	id := nextID()

	user := store.User{
		ID:   id,
		VkID: vkID,
		Name: fmt.Sprintf("User #%d", vkID),
	}

	userBytes, _ := json.Marshal(user)
	_, err := store.RedisClient.Set(context.Background(), fmt.Sprintf("node:%d", user.ID), userBytes, 0).Result()
	if err != nil {
		return 0, fmt.Errorf("error saving user: %w", err)
	}

	return id, nil
}

func nextID() int {
	return int(time.Now().UnixMilli())
}

func usersList(ids []string) []*User {
	keys := make([]string, len(ids))
	for i, userID := range ids {
		keys[i] = fmt.Sprintf("node:%s", userID)
	}

	userBytesList, err := store.RedisClient.MGet(context.Background(), keys...).Result()
	if err != nil {
		log.Printf("Error getting users: %s", err)
	}

	var users []*store.User
	for _, userBytes := range userBytesList {
		if userBytes == nil {
			continue
		}

		user := &store.User{}
		err = json.Unmarshal([]byte(userBytes.(string)), user)
		if err != nil {
			log.Printf("Error unmarshalling user: %s", err)
			continue
		}

		users = append(users, user)
	}

	apiUsers := make([]*User, len(users))
	for i, user := range users {
		apiUsers[i] = &User{
			ID:   strconv.Itoa(user.ID),
			Name: user.Name,
		}
	}

	return apiUsers
}
