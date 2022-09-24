package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/store"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func authExchangeCode(code string, redirectURI string) (int, string, error) {
	vkAppID := store.DefaultConfig.VKAppID
	vkAppSecret := store.DefaultConfig.VKAppSecret

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(vkAppID)},
		"client_secret": []string{vkAppSecret},
		"redirect_uri":  []string{redirectURI},
		"code":          []string{code},
	})
	if err != nil {
		return 0, "", fmt.Errorf("http error: %s", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("error reading body: %s", err)
	}

	body := struct {
		AccessToken string `json:"access_token"`
		UserID      int    `json:"user_id"`
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return 0, "", fmt.Errorf("error parsing json: %w", err)
	} else if body.AccessToken == "" {
		return 0, "", fmt.Errorf("no access_token in response")
	}

	return body.UserID, body.AccessToken, nil
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

var ErrInvalidCredentials = fmt.Errorf("invalid credentials")

func authEmailAuth(email, password string) (int, error) {
	userIDStr, err := store.RedisClient.Get(context.Background(), fmt.Sprintf("map_email2id:%s", email)).Result()
	if err == redis.Nil {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	userID, _ := strconv.Atoi(userIDStr)

	user := store.User{}
	err = store.NodeGet(userID, &user)
	if err != nil {
		return 0, fmt.Errorf("error getting user: %s", err)
	}

	if password != user.PasswordHash {
		return 0, ErrInvalidCredentials
	}

	return user.ID, nil
}

func authRegister(email, password string) (int, error) {
	id := int(time.Now().UnixMilli())
	user := store.User{
		ID:           id,
		Name:         fmt.Sprintf("User #%d", id),
		Email:        email,
		PasswordHash: password,
	}
	err := store.NodeSave(user.ID, user)
	if err != nil {
		return 0, fmt.Errorf("error saving user: %w", err)
	}

	key := fmt.Sprintf("map_email2id:%s", email)
	wasSet, err := store.RedisClient.SetNX(context.Background(), key, user.ID, 0).Result()
	if err != nil {
		return 0, fmt.Errorf("error saving map key: %s", err)
	} else if !wasSet {
		return 0, fmt.Errorf("key was not set")
	}

	return user.ID, nil
}

func authValidateCredentials(email, password string) string {
	if email == "" {
		return "empty email"
	}
	if !strings.Contains(email, "@") {
		return "incorrect email"
	}
	if len(email) > 200 {
		return "email too long"
	}

	key := fmt.Sprintf("map_email2id:%s", email)
	_, err := store.RedisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// ok
	} else if err != nil {
		log.Printf("Error: %s", err)
		return "cannot check email"
	} else {
		return "email already registered"
	}

	if password == "" {
		return "empty password"
	}

	return ""
}
