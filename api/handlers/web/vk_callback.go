package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/materkov/meme9/api/handlers"
	"github.com/materkov/meme9/api/pkg"
	"github.com/materkov/meme9/api/store"
)

type VKCallback struct {
	Store  *store.Store
	Config *pkg.Config
}

func writeInternalError(w http.ResponseWriter, err error) {
	log.Printf("[ERROR] VK auth internal error: %s", err)
	_, _ = fmt.Fprint(w, "Internal error")
}

func (v *VKCallback) Handle(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		_, _ = fmt.Fprint(w, "Empty VK code")
		return
	}

	proxyScheme := r.Header.Get("x-forwarded-proto")
	if proxyScheme == "" {
		proxyScheme = "http"
	}

	redirectURI := fmt.Sprintf("%s://%s%s", proxyScheme, r.Host, r.URL.EscapedPath())

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(handlers.VKAppID)},
		"client_secret": []string{v.Config.VKAppSecret},
		"redirect_uri":  []string{redirectURI},
		"code":          []string{code},
	})
	if err != nil {
		writeInternalError(w, fmt.Errorf("http vk error: %v", err))
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		writeInternalError(w, fmt.Errorf("failed reading http body: %v", err))
		return
	}

	body := struct {
		AccessToken string `json:"access_token"`
		UserID      int    `json:"user_id"`
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		writeInternalError(w, fmt.Errorf("incorrect json: %s", bodyBytes))
		return
	} else if body.AccessToken == "" {
		writeInternalError(w, fmt.Errorf("incorrect response: %s", bodyBytes))
		return
	}

	nodeID, err := v.Store.GetUserByVKID(body.UserID)
	if err != nil {
		writeInternalError(w, fmt.Errorf("failed getting by vk id: %v", err))
		return
	}

	var user *store.User

	if nodeID == 0 {
		userID, err := v.Store.GenerateNodeID()
		if err != nil {
			writeInternalError(w, fmt.Errorf("failed generating node id: %v", err))
			return
		}

		err = v.Store.SaveUserVKID(body.UserID, userID)
		if err != nil {
			writeInternalError(w, fmt.Errorf("failed saving vk id: %v", err))
			return
		}

		user = &store.User{
			ID:   userID,
			Name: fmt.Sprintf("User #%d", body.UserID),
			VkID: body.UserID,
		}
		err = v.Store.AddUser(user)
		if err != nil {
			writeInternalError(w, fmt.Errorf("failed saving user: %v", err))
			return
		}
	} else {
		user, err = v.Store.GetUser(nodeID)
		if err != nil {
			writeInternalError(w, fmt.Errorf("failed getting user: %v", err))
			return
		}
	}

	tokenID, err := v.Store.GenerateNodeID()
	if err != nil {
		writeInternalError(w, fmt.Errorf("failed generating node id for token: %v", err))
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
		writeInternalError(w, fmt.Errorf("failed saving token: %v", err))
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
