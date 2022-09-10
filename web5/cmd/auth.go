package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/store"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func authExchangeCode(code string, redirectURI string) (int, error) {
	vkAppID := store.DefaultConfig.VKAppID
	vkAppSecret := store.DefaultConfig.VKAppSecret

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(vkAppID)},
		"client_secret": []string{vkAppSecret},
		"redirect_uri":  []string{redirectURI},
		"code":          []string{code},
	})
	if err != nil {
		return 0, fmt.Errorf("http error: %s", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading body: %s", err)
	}

	body := struct {
		AccessToken string `json:"access_token"`
		UserID      int    `json:"user_id"`
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return 0, fmt.Errorf("error parsing json: %w", err)
	} else if body.AccessToken == "" {
		return 0, fmt.Errorf("no access_token in response")
	}

	return body.UserID, nil
}

func randStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func authCreateToken(userID int) (string, error) {
	token := store.AuthToken{
		//ID:     nextID(),
		UserID: userID,
		Token:  randStringRunes(30),
		Date:   int(time.Now().Unix()),
	}

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	_, err = store.RedisClient.Set(context.Background(), fmt.Sprintf("auth_token:%s", token.Token), tokenBytes, time.Minute*10).Result()
	if err != nil {
		return "", err
	}

	return token.Token, err
}

func authCheckToken(tokenStr string) (int, error) {
	if tokenStr == "" {
		return 0, nil
	}

	tokenBytes, err := store.RedisClient.Get(context.Background(), fmt.Sprintf("auth_token:%s", tokenStr)).Bytes()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	token := store.AuthToken{}
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return 0, err
	}

	return token.UserID, err
}
