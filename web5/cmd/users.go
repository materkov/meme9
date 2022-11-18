package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/files"
	"github.com/materkov/meme9/web5/pkg/users"
	"github.com/materkov/meme9/web5/store"
	"log"
	"strconv"
)

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

		isFollowing, err := users.IsFollowing(viewerID, ids)
		if err != nil {
			log.Printf("Error getting is followed: %s", err)
		}
		chanIsFollowing <- isFollowing
	}()

	usersMap := <-chanUsersMap
	isFollowing := <-chanIsFollowing

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
		apiUser.Bio = user.Bio
		apiUser.IsFollowing = isFollowing[userID]

		if user.AvatarSha != "" {
			apiUser.Avatar = files.GetURL(user.AvatarSha)
		} else if user.VkPhoto200 != "" {
			apiUser.Avatar = user.VkPhoto200
		}
	}

	return apiUsers
}
