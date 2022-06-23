package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func StorageAuthorize() (string, error) {
	body := struct {
		Auth struct {
			Identity struct {
				Methods  []string `json:"methods"`
				Password struct {
					User struct {
						ID       string `json:"id"`
						Password string `json:"password"`
					} `json:"user"`
				} `json:"password"`
			} `json:"identity"`
		} `json:"auth"`
	}{}

	body.Auth.Identity.Methods = []string{"password"}
	body.Auth.Identity.Password.User.ID = GlobalConfig.SelectelStorageUser
	body.Auth.Identity.Password.User.Password = GlobalConfig.SelectelStoragePassword

	bodyBytes, _ := json.Marshal(body)
	resp, err := http.Post("https://api.selcdn.ru/v3/auth/tokens", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("error doping http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("incorrect http code: %d", resp.StatusCode)
	}

	token := resp.Header.Get("x-subject-token")
	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	return token, nil
}

func StorageUpload(filePath string, file []byte, authToken string) error {
	req, err := http.NewRequest("PUT", "https://api.selcdn.ru/v1/SEL_200213/meme-files/"+filePath, bytes.NewReader(file))
	if err != nil {
		return fmt.Errorf("error creating selectel http request: %w", err)
	}

	req.Header.Set("X-Auth-Token", authToken)
	req.Header.Set("Content-Type", "image/jpeg")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error doing http: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		return fmt.Errorf("incorrect http response code: %d, %s", resp.StatusCode, bodyBytes)
	}

	return nil
}
