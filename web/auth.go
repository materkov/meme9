package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrNotApplicable = fmt.Errorf("this auth not applicable")
	ErrAuthFailed    = fmt.Errorf("auth failed")
)

type Auth struct {
	store *Store
}

func (a *Auth) tryVkAuth(authUrl string) (int, error) {
	parsedUrl, err := url.Parse(authUrl)
	if err != nil || authUrl == "" {
		return 0, ErrNotApplicable
	}

	vkUserID, _ := strconv.Atoi(parsedUrl.Query().Get("vk_user_id"))
	if vkUserID == 0 {
		return 0, ErrNotApplicable
	}

	keys := make([]string, 0)
	for key := range parsedUrl.Query() {
		if strings.HasPrefix(key, "vk_") {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	for i, key := range keys {
		keys[i] = fmt.Sprintf("%s=%s", key, parsedUrl.Query().Get(key))
	}

	signString := strings.Join(keys, "&")

	mac := hmac.New(sha256.New, []byte(config.VKMiniAppSecret))
	_, _ = mac.Write([]byte(signString))
	computedSign := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	actualSign := parsedUrl.Query().Get("sign")
	if computedSign != actualSign {
		return 0, ErrAuthFailed
	}

	userID, err := store.GetByVkID(vkUserID)
	if err != nil {
		return 0, fmt.Errorf("failed getting user by vk id: %w", err)
	} else if userID != 0 {
		return userID, nil
	}

	// New user
	userID, err = store.GenerateNextID(ObjectTypeUser)
	if err != nil {
		return 0, fmt.Errorf("error generating object ID: %w", err)
	}

	user := User{
		ID:   userID,
		Name: fmt.Sprintf("VK User %d", vkUserID),
		VkID: vkUserID,
	}
	err = store.AddUserByVK(&user)
	if err != nil {
		return 0, fmt.Errorf("error saving user: %w", err)
	}

	return userID, nil
}

func (a *Auth) tryCookieAuth(r *http.Request) (*Token, error) {
	accessCookie, err := r.Cookie("access_token")
	if err != nil || accessCookie.Value == "" {
		return nil, ErrNotApplicable
	}

	headerToken := r.Header.Get("x-csrf-token")
	validCSRFToken := GenerateCSRFToken(accessCookie.Value)

	if r.URL.Path == "/api" && headerToken != validCSRFToken {
		return nil, ErrAuthFailed
	}

	return a.tryTokenAuth(accessCookie.Value)
}

func (a *Auth) tryHeaderAuth(authHeader string) (*Token, error) {
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, ErrNotApplicable
	}

	return a.tryTokenAuth(authHeader)
}

func (a *Auth) tryTokenAuth(tokenStr string) (*Token, error) {
	token, err := store.GetToken(tokenStr)
	if err == ErrObjectNotFound {
		return nil, ErrAuthFailed
	} else if err != nil {
		return nil, fmt.Errorf("error selecting token: %w", err)
	}

	return token, nil
}

func (a *Auth) tryAuth(r *http.Request) (*Token, int, error) {
	token, err := a.tryHeaderAuth(r.Header.Get("authorization"))
	if err == nil {
		return token, token.UserID, err
	}

	token, err = a.tryCookieAuth(r)
	if err == nil {
		return token, token.UserID, err
	}

	userID, err := a.tryVkAuth(r.URL.String())
	if err == nil {
		return nil, userID, err
	}

	userID, err = a.tryVkAuth(r.Header.Get("x-vk-auth"))
	if err == nil {
		return nil, userID, err
	}

	return nil, 0, fmt.Errorf("not authorized")
}
