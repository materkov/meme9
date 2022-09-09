package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"strconv"
	"time"
)

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

func usersAdd(vkID int) (int, error) {
	id := nextID()

	user := User{
		ID:   id,
		VkID: vkID,
		Name: fmt.Sprintf("User #%d", vkID),
	}

	userBytes, _ := json.Marshal(user)
	_, err := redisClient.Set(context.Background(), fmt.Sprintf("node:%d", user.ID), userBytes, 0).Result()
	if err != nil {
		return 0, fmt.Errorf("error saving user: %w", err)
	}

	return id, nil
}

func nextID() int {
	return int(time.Now().UnixMilli())
}
