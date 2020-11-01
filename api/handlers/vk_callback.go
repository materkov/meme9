package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type VKCallback struct {
	Store *store.Store
}

func (v *VKCallback) Handle(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		writeError(w, &login.ErrorRenderer{DisplayText: "empty code"})
		return
	}

	proxyScheme := r.Header.Get("x-forwarded-proto")
	if proxyScheme == "" {
		proxyScheme = "http"
	}

	clientSecret, err := v.Store.GetVkSecret()
	if err != nil {
		writeInternalError(w, err)
		return
	}

	redirectURI := fmt.Sprintf("%s://%s%s", proxyScheme, r.Host, r.URL.EscapedPath())

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(VKAppID)},
		"client_secret": []string{clientSecret},
		"redirect_uri":  []string{redirectURI},
		"code":          []string{code},
	})
	if err != nil {
		writeInternalError(w, err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	body := struct {
		AccessToken string `json:"access_token"`
		UserID      int    `json:"user_id"`
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		writeInternalError(w, err)
		return
	} else if body.AccessToken == "" {
		writeInternalError(w, fmt.Errorf("error from vk: %s", bodyBytes))
		return
	}

	nodeID, err := v.Store.GetUserByVKID(body.UserID)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	var user *store.User

	if nodeID == 0 {
		userID, err := v.Store.GenerateNodeID()
		if err != nil {
			writeInternalError(w, err)
			return
		}

		err = v.Store.SaveUserVKID(body.UserID, userID)
		if err != nil {
			writeInternalError(w, err)
			return
		}

		user = &store.User{
			ID:   userID,
			Name: fmt.Sprintf("User #%d", body.UserID),
			VkID: body.UserID,
		}
		err = v.Store.AddUser(user)
		if err != nil {
			writeInternalError(w, err)
			return
		}
	} else {
		user, err = v.Store.GetUser(nodeID)
		if err != nil {
			writeInternalError(w, fmt.Errorf("error getting user node: %w", err))
			return
		}
	}

	tokenID, err := v.Store.GenerateNodeID()
	if err != nil {
		writeInternalError(w, err)
		return
	}

	token := store.Token{
		ID:       tokenID,
		IssuedAt: int(time.Now().Unix()),
		UserID:   user.ID,
	}
	token.GenerateRandomToken()

	err = v.Store.AddToken(&token)
	if err != nil {
		writeInternalError(w, err)
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
