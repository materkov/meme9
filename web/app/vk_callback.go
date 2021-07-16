package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/materkov/meme9/web/store"
)

func (a *App) DoVKCallback(ctx context.Context, code string, viewer *Viewer) (string, error) {
	if code == "" {
		return "", fmt.Errorf("empty VK code")
	}

	redirectURI := fmt.Sprintf("%s://%s/vk-callback", viewer.RequestScheme, viewer.RequestHost)

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(DefaultConfig.VKAppID)},
		"client_secret": []string{DefaultConfig.VKAppSecret},
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

	assocType := store.Assoc_VK_ID + strconv.Itoa(body.UserID)
	assocs, err := a.Store.AssocRange(ctx, 0, assocType, 1)
	if err != nil {
		return "", fmt.Errorf("error selecting by vk id: %w", err)
	}

	userID := 0
	var user *store.User
	if len(assocs) > 0 {
		userID = assocs[0].VkID.ID2
		obj, err := a.Store.ObjGet(ctx, userID)
		if err != nil {
			return "", fmt.Errorf("error getting users: %w", err)
		} else if obj == nil || obj.User == nil {
			return "", fmt.Errorf("nil or not user object")
		}

		user = obj.User
	} else {
		userID, err = a.Store.GenerateNextID()
		if err != nil {
			return "", fmt.Errorf("error generating user id: %w", err)
		}

		user = &store.User{
			ID: userID,
		}

		err = a.Store.ObjAdd(&store.StoredObject{ID: userID, User: &store.User{
			ID: userID,
		}})
		if err != nil {
			return "", fmt.Errorf("error saving obj: %w", err)
		}

		assocType := store.Assoc_VK_ID + strconv.Itoa(body.UserID)
		err = a.Store.AssocAdd(0, userID, assocType, &store.StoredAssoc{VkID: &store.VkID{
			ID1:  0,
			ID2:  userID,
			Type: assocType,
		}})
		if err != nil {
			return "", fmt.Errorf("error saving assoc: %w", err)
		}
	}

	vkName, vkAvatar, err := fetchVkData(body.UserID, body.AccessToken)
	if err != nil {
		log.Printf("Error getting vk data: %s", err)
	} else {
		user.Name = vkName
		user.VkAvatar = vkAvatar
		err = a.Store.ObjUpdate(&store.StoredObject{ID: user.ID, User: user})
		if err != nil {
			return "", fmt.Errorf("failed updating name and avatar: %w", err)
		}
	}

	objectID, err := a.Store.GenerateNextID()
	if err != nil {
		return "", fmt.Errorf("failed generating object id: %w", err)
	}

	token := fmt.Sprintf("%d-%s", objectID, RandString(40))
	err = a.Store.ObjAdd(&store.StoredObject{ID: objectID, Token: &store.Token{
		ID:     objectID,
		Token:  token,
		UserID: userID,
	}})
	if err != nil {
		return "", fmt.Errorf("failed saving token: %w", err)
	}

	return token, nil
}

