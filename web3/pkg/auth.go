package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func GetRedirectURL(origin string) string {
	vkAppID := GlobalConfig.VKAppID

	redirectURL := fmt.Sprintf("%s/vk-callback", origin)
	redirectURL = url.QueryEscape(redirectURL)
	vkURL := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%d&response_type=code&redirect_uri=%s", vkAppID, redirectURL)

	return vkURL
}

func ExchangeCode(origin string, code string) (int, error) {
	vkAppID := GlobalConfig.VKAppID
	vkAppSecret := GlobalConfig.VKAppSecret

	redirectURI := fmt.Sprintf("%s/vk-callback", origin)

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
