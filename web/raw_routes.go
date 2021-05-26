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
	code := r.URL.Query().Get("code")
	if code == "" {
		_, _ = fmt.Fprint(w, "Empty VK code")
		return
	}

	// TODO
	proxyScheme := "http"
	if r.Header.Get("x-forwarded-proto") != "" {
		proxyScheme = r.Header.Get("x-forwarded-proto")
	}

	vkAppID := 7260220

	redirectURI := fmt.Sprintf("%s://%s%s", proxyScheme, r.Host, r.URL.EscapedPath())

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(vkAppID)},
		"client_secret": []string{config.VKAppSecret},
		"redirect_uri":  []string{redirectURI},
		"code":          []string{code},
	})
	if err != nil {
		fmt.Fprintf(w, "http vk error: %v", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "failed reading http body: %v", err)
		return
	}

	body := struct {
		AccessToken string `json:"access_token"`
		UserID      int    `json:"user_id"`
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		fmt.Fprintf(w, "incorrect json: %s", bodyBytes)
		return
	} else if body.AccessToken == "" {
		fmt.Fprintf(w, "incorrect response: %s", bodyBytes)
		return
	}

	userID, err := store.GetByVkID(body.UserID)
	if err != nil {
		log.Printf("Error selecting by vk id: %s", err)
		fmt.Fprintf(w, "internal error")
		return
	}

	users, err := store.GetUsers([]int{userID})
	if err != nil {
		log.Printf("Error selecting user: %s", err)
		fmt.Fprintf(w, "internal error")
		return
	} else if len(users) == 0 {
		log.Printf("User %d not found", userID)
		fmt.Fprintf(w, "internal error")
		return
	}

	user := users[0]

	vkName, vkAvatar, err := fetchVkData(body.UserID, body.AccessToken)
	if err == nil {
		user.Name = vkName
		user.VkAvatar = vkAvatar
		err = store.UpdateNameAvatar(user)
		if err != nil {
			log.Printf("Failed saving new name&avatar: %s", err)
		}
	}

	token := Token{
		Token:  RandString(50),
		UserID: userID,
	}
	err = store.AddToken(&token)
	if err != nil {
		log.Printf("error saving token: %s", err)
		fmt.Fprintf(w, "internal error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.Token,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
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
