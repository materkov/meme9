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

func usersList(ids []int, viewerID int, includeIsFollowing bool, includeFollowersCount bool) []*User {
	chanUsersMap := make(chan map[int]*store.User)

	go func() {
		keys := make([]string, len(ids))
		for i, userID := range ids {
			keys[i] = fmt.Sprintf("node:%d", userID)
		}

		userBytesList, err := store.RedisClient.MGet(context.Background(), keys...).Result()
		if err != nil {
			log.Printf("Error getting users: %s", err)
		}

		usersMap := map[int]*store.User{}
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

			usersMap[user.ID] = user
		}

		chanUsersMap <- usersMap
	}()

	chanIsFollowing := make(chan map[int]bool)
	go func() {
		if !includeIsFollowing {
			chanIsFollowing <- nil
			return
		}

		isFollowing, err := usersIsFollowing(viewerID, ids)
		if err != nil {
			log.Printf("Error getting is followed: %s", err)
		}
		chanIsFollowing <- isFollowing
	}()

	chanFollowingCount := make(chan map[int]int)
	go func() {
		if !includeFollowersCount {
			chanFollowingCount <- nil
			return
		}

		result := map[int]int{}
		for _, id := range ids {
			count, err := usersFollowingCount(id)
			if err != nil {
				log.Printf("Error getting following count: %s", err)
			}
			result[id] = count
		}

		chanFollowingCount <- result
	}()

	chanFollowedByCount := make(chan map[int]int)
	go func() {
		if !includeFollowersCount {
			chanFollowedByCount <- nil
			return
		}

		result := map[int]int{}
		for _, id := range ids {
			count, err := usersFollowedByCount(id)
			if err != nil {
				log.Printf("Error getting followedBy count: %s", err)
			}
			result[id] = count
		}

		chanFollowedByCount <- result
	}()

	usersMap := <-chanUsersMap
	isFollowing := <-chanIsFollowing
	followedByCount := <-chanFollowedByCount
	followingCount := <-chanFollowingCount

	apiUsers := make([]*User, len(ids))
	for i, userID := range ids {
		apiUser := &User{
			ID: strconv.Itoa(userID),
		}

		apiUsers[i] = apiUser

		user, ok := usersMap[userID]
		if !ok {
			continue
		}

		apiUser.Name = user.Name
		apiUser.Avatar = user.VkPhoto200
		apiUser.Bio = user.Bio
		apiUser.IsFollowing = isFollowing[userID]
		apiUser.FollowingCount = followingCount[userID]
		apiUser.FollowedByCount = followedByCount[userID]
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

func usersFollow(userID, targetID int) error {
	score := float64(time.Now().UnixMilli())

	pipe := store.RedisClient.Pipeline()
	pipe.ZAdd(context.Background(), fmt.Sprintf("following:%d", userID), redis.Z{
		Score:  score,
		Member: targetID,
	})
	pipe.ZAdd(context.Background(), fmt.Sprintf("followed_by:%d", targetID), redis.Z{
		Score:  score,
		Member: userID,
	})

	_, err := pipe.Exec(context.Background())
	if err != nil {
		return fmt.Errorf("error storing followers key: %w", err)
	}

	return nil
}

func usersUnfollow(userID, targetID int) error {
	pipe := store.RedisClient.Pipeline()
	pipe.ZRem(context.Background(), fmt.Sprintf("following:%d", userID), targetID)
	pipe.ZRem(context.Background(), fmt.Sprintf("followed_by:%d", targetID), userID)

	_, err := pipe.Exec(context.Background())
	if err != nil {
		return fmt.Errorf("error removing followers from zset: %w", err)
	}

	return nil
}

func usersIsFollowing(userID int, targetIds []int) (map[int]bool, error) {
	targetsStr := make([]string, len(targetIds))
	for i, targetID := range targetIds {
		targetsStr[i] = strconv.Itoa(targetID)
	}

	scores, err := store.RedisClient.ZMScore(context.Background(), fmt.Sprintf("following:%d", userID), targetsStr...).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting scores: %w", err)
	}

	result := map[int]bool{}
	for i, targetID := range targetIds {
		if scores[i] != 0 {
			result[targetID] = true
		}
	}

	return result, nil
}

func usersFollowingCount(userID int) (int, error) {
	result, err := store.RedisClient.ZCard(context.Background(), fmt.Sprintf("following:%d", userID)).Result()
	return int(result), err
}

func usersFollowedByCount(userID int) (int, error) {
	result, err := store.RedisClient.ZCard(context.Background(), fmt.Sprintf("followed_by:%d", userID)).Result()
	return int(result), err
}
