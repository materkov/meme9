package users

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/store"
	"net/http"
	"strconv"
	"time"
)

func GetOrCreateByVKID(vkID int) (int, error) {
	mapKey := fmt.Sprintf("vk2user_map:%d", vkID)
	userIDStr, err := store.RedisClient.Get(context.Background(), mapKey).Result()
	if err == redis.Nil {
		userID, err := Add(vkID)
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

func Add(vkID int) (int, error) {
	user := store.User{
		VkID: vkID,
		Name: fmt.Sprintf("User #%d", vkID),
	}

	id, err := store.NodeInsert(store.ObjectTypeUser, user)
	if err != nil {
		return 0, fmt.Errorf("error saving user: %w", err)
	}

	return id, nil
}

func RefreshFromVk(ctx context.Context, id int) error {
	user := store.CachedStoreFromCtx(ctx).User.Get(id)
	if user == nil {
		return fmt.Errorf("user not found")
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

	err = store.NodeUpdate(user.ID, user)
	if err != nil {
		return err
	}

	return nil
}

func Follow(userID, targetID int) error {
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

func Unfollow(userID, targetID int) error {
	pipe := store.RedisClient.Pipeline()
	pipe.ZRem(context.Background(), fmt.Sprintf("following:%d", userID), targetID)
	pipe.ZRem(context.Background(), fmt.Sprintf("followed_by:%d", targetID), userID)

	_, err := pipe.Exec(context.Background())
	if err != nil {
		return fmt.Errorf("error removing followers from zset: %w", err)
	}

	return nil
}

func IsFollowing(userID int, targetIds []int) (map[int]bool, error) {
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
