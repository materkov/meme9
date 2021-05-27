package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func handleVKCallback(w http.ResponseWriter, r *http.Request) {
	viewer := GetViewerFromContext(r.Context())
	accessToken, err := doVKCallback(r.URL.Query().Get("code"), viewer)
	if err != nil {
		log.Printf("Error: %s", err)
		_, _ = fmt.Fprint(w, "Failed to authorize")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func doVKCallback(code string, viewer *Viewer) (string, error) {
	if code == "" {
		return "", fmt.Errorf("empty VK code")
	}

	vkAppID := 7260220

	redirectURI := fmt.Sprintf("%s://%s/vk-callback", viewer.RequestScheme, viewer.RequestHost)

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(vkAppID)},
		"client_secret": []string{config.VKAppSecret},
		"redirect_uri":  []string{redirectURI},
		"code":          []string{code},
	})
	if err != nil {
		return "", fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}

	body := struct {
		AccessToken string `json:"access_token"`
		UserID      int    `json:"user_id"`
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return "", fmt.Errorf("error parsing json %s: %s", bodyBytes, err)
	} else if body.AccessToken == "" {
		return "", fmt.Errorf("empty access token: %s", bodyBytes)
	}

	userID, err := store.GetByVkID(body.UserID)
	if err != nil {
		return "", fmt.Errorf("error selecting by vk id: %w", err)
	}

	var user *User
	users, err := store.GetUsers([]int{userID})
	if err != nil {
		return "", fmt.Errorf("error getting users: %w", err)
	} else if len(users) == 1 {
		user = users[0]
	} else {
		userID, err = store.GenerateNextID(ObjectTypeUser)
		if err != nil {
			return "", fmt.Errorf("error generating user id: %w", err)
		}

		user = &User{
			ID:   userID,
			VkID: body.UserID,
		}

		err = store.AddUserByVK(user)
		if err != nil {
			return "", fmt.Errorf("error saving user: %w", err)
		}
	}

	vkName, vkAvatar, err := fetchVkData(body.UserID, body.AccessToken)
	if err != nil {
		log.Printf("Error getting vk data: %s", err)
	} else {
		user.Name = vkName
		user.VkAvatar = vkAvatar
		err = store.UpdateNameAvatar(user)
		if err != nil {
			return "", fmt.Errorf("failed updating name and avatar: %w", err)
		}
	}

	token := Token{
		Token:  RandString(50),
		UserID: userID,
	}
	err = store.AddToken(&token)
	if err != nil {
		return "", fmt.Errorf("failed saving token: %w", err)
	}

	return token.Token, nil
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
