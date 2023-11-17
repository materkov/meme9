package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/xlog"
	"github.com/materkov/meme9/web6/src/store"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var HTTPClient = &http.Client{}

var VKEndpoint = "https://oauth.vk.com"
var VKAPIEndpoint = "https://api.vk.com"

func ExchangeCode(code string, redirectURI string) (int, string, error) {
	vkAppID := store.GlobalConfig.VKAppID
	vkAppSecret := store.GlobalConfig.VKAppSecret

	xlog.Log("Exchanging VK auth code", xlog.Fields{
		"code":        code,
		"redirectURI": redirectURI,
	})

	resp, err := HTTPClient.PostForm(VKEndpoint+"/access_token", url.Values{
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

	xlog.Log("VK auth response", xlog.Fields{
		"response": string(bodyBytes),
	})

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

func RefreshFromVk(accessToken string, vkUserID int) (string, error) {
	args := fmt.Sprintf("v=5.180&access_token=%s&user_ids=%d&fields=photo_200", accessToken, vkUserID)
	resp, err := HTTPClient.Post(VKAPIEndpoint+"/method/users.get?"+args, "", nil)
	if err != nil {
		return "", fmt.Errorf("http error: %s", err)
	}
	defer resp.Body.Close()

	body := struct {
		Response []struct {
			ID        int    `json:"id"`
			Photo200  string `json:"photo_200"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"response"`
		Error json.RawMessage `json:"error"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", fmt.Errorf("incorrect json: %w", err)
	}

	if body.Error != nil {
		return "", fmt.Errorf("error response from vk: %w", err)
	} else if len(body.Response) == 0 || body.Response[0].ID != vkUserID {
		return "", fmt.Errorf("user not found")
	}

	//user.VkPhoto200 = body.Response[0].Photo200
	userName := fmt.Sprintf("%s %s", body.Response[0].FirstName, body.Response[0].LastName)

	return userName, nil
}
