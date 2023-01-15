package api

import (
	"context"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/auth"
	"github.com/materkov/meme9/web5/pkg/telegram"
	"github.com/materkov/meme9/web5/pkg/users"
	"github.com/materkov/meme9/web5/store"
	"log"
)

type AuthVkCallback struct {
	Code        string
	RedirectURI string
}

type Authorization struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

func handleAuthVkCallback(ctx context.Context, viewerID int, req *AuthVkCallback) (*Authorization, error) {
	code := req.Code
	redirectURI := req.RedirectURI

	vkID, vkAccessToken, err := auth.ExchangeCode(code, redirectURI)
	if err != nil {
		return nil, err
	}

	userID, err := users.GetOrCreateByVKID(vkID)
	if err != nil {
		return nil, err
	}

	user := &store.User{}
	err = store.NodeGet(userID, user)
	if err != nil {
		return nil, err
	}

	user.VkAccessToken = vkAccessToken
	err = store.NodeSave(user.ID, user)
	if err != nil {
		log.Printf("error saving user")
	}

	_, _ = store.RedisClient.RPush(context.Background(), "queue", user.ID).Result()

	authToken, err := auth.CreateToken(userID)
	if err != nil {
		return nil, err
	}

	err = telegram.SendNotify(fmt.Sprintf("meme new login: https://vk.com/id%d", user.VkID))
	if err != nil {
		log.Printf("Error sending telegram notify: %s", err)
	}

	auth := Authorization{
		Token: authToken,
		User:  nil,
	}

	result := handleUserById(ctx, viewerID, fmt.Sprintf("/users/%d", userID))
	wrappedUser := result[0].(User)

	auth.User = &wrappedUser

	return &auth, nil
}

type AuthEmailLogin struct {
	Email    string
	Password string
}

func handleAuthEmailLogin(ctx context.Context, viewerID int, req *AuthEmailLogin) (*Authorization, error) {
	userID, err := auth.EmailAuth(req.Email, req.Password)
	if err == auth.ErrInvalidCredentials {
		return nil, fmt.Errorf("invalid credentials")
	} else if err != nil {
		return nil, err
	}

	token, err := auth.CreateToken(userID)
	if err != nil {
		return nil, err
	}

	auth := Authorization{
		Token: token,
	}

	result := handleUserById(ctx, viewerID, fmt.Sprintf("/users/%d", userID))
	wrappedUser := result[0].(User)

	auth.User = &wrappedUser

	return &auth, nil
}

type AuthEmailRegister struct {
	Email    string
	Password string
}

func handleAuthEmailRegister(ctx context.Context, viewerID int, req *AuthEmailRegister) (*Authorization, error) {
	email, password := req.Email, req.Password
	validateErr := auth.ValidateCredentials(email, password)
	if validateErr != "" {
		return nil, fmt.Errorf("invalid credentials")
	}

	userID, err := auth.Register(email, password)
	if err != nil {
		return nil, err
	}

	token, err := auth.CreateToken(userID)
	if err != nil {
		return nil, err
	}

	auth := Authorization{
		Token: token,
	}

	result := handleUserById(ctx, viewerID, fmt.Sprintf("/users/%d", userID))
	wrappedUser := result[0].(User)

	auth.User = &wrappedUser

	return &auth, nil
}
