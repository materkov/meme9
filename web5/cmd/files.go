package main

import (
	"bytes"
	"fmt"
	"github.com/materkov/meme9/web5/store"
	"net/http"
)

func filesSelectelUpload(file []byte, fileName string) error {
	req, _ := http.NewRequest("GET", "https://api.selcdn.ru/auth/v1.0", nil)
	req.Header.Set("X-Auth-User", fmt.Sprintf("%d_%s", store.DefaultConfig.SelectelAccountID, store.DefaultConfig.SelectelUserName))
	req.Header.Set("X-Auth-Key", store.DefaultConfig.SelectelUserPassword)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	_ = resp.Body.Close()

	authToken := resp.Header.Get("X-Auth-Token")
	if authToken == "" {
		return fmt.Errorf("empty auth token")
	}

	url := fmt.Sprintf("https://api.selcdn.ru/v1/SEL_%d/meme-files/avatars/%s", store.DefaultConfig.SelectelAccountID, fileName)
	req, _ = http.NewRequest("PUT", url, bytes.NewReader(file))
	req.Header.Set("X-Auth-Token", authToken)
	req.Header.Set("Content-Type", "image/jpeg")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("incorrect http status: %d", resp.StatusCode)
	}

	return nil
}

func filesGetURL(fileSha string) string {
	return fmt.Sprintf("https://689809.selcdn.ru/meme-files/avatars/%s", fileSha)
}
