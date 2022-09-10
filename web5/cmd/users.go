package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/store"
	"log"
	"net/http"
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

	err := store.NodeSave(id, user)
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
			ID:     strconv.Itoa(user.ID),
			Name:   user.Name,
			Avatar: user.VkPhoto200,
		}
	}

	return apiUsers
}

func usersRefreshFromVk(id int) error {
	user := store.User{}
	err := store.NodeGet(id, &user)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	args := fmt.Sprintf("v=5.180&access_token=%s&user_ids=%d&fields=photo_200", user.VkAccessToken, user.VkID)
	resp, err := http.Post("https://api.vk.com/method/users.get?"+args, "", nil)
	if err != nil {
		return fmt.Errorf("http error: %s", err)
	}
	defer resp.Body.Close()

	body := struct {
		Response []struct {
			ID        int    `json:"id"`
			Photo200  string `json:"photo_200"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"response"`
		Error json.RawMessage `json:"error"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return fmt.Errorf("incorrect json: %w", err)
	}

	if body.Error != nil {
		return fmt.Errorf("error response from vk: %w", err)
	} else if len(body.Response) == 0 || body.Response[0].ID != user.VkID {
		return fmt.Errorf("user not found")
	}

	user.VkPhoto200 = body.Response[0].Photo200
	user.Name = fmt.Sprintf("%s %s", body.Response[0].FirstName, body.Response[0].LastName)

	err = store.NodeSave(user.ID, user)
	if err != nil {
		return err
	}

	return nil
}
