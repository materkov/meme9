package lib

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web2/store"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func ProcessVKkCallback(code string) (int, error) {
	redirectURI := fmt.Sprintf("%s://%s/vk-callback", DefaultConfig.RequestScheme, DefaultConfig.RequestHost)

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(DefaultConfig.VkAppID)},
		"client_secret": []string{DefaultConfig.VkAppSecret},
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

func GetOrCreateUserByVkID(s *store.Store, vkID int) (*store.User, error) {
	user, err := s.User.GetByVkID(vkID)
	if err != nil {
		return nil, fmt.Errorf("error selecting user by vk id: %w", err)
	}

	if user != nil {
		return user, nil
	}

	user = &store.User{
		Name: fmt.Sprintf("VK User #%d", vkID),
		VkID: vkID,
	}
	err = s.User.Add(user)
	if err != nil {
		return nil, fmt.Errorf("error adding user by id")
	}

	return user, nil
}
