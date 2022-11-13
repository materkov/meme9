package files

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/materkov/meme9/web5/store"
	"net/http"
)

const mainFolder = "b"

func GetURL(fileSha string) string {
	return fmt.Sprintf("https://689809.selcdn.ru/meme-files/%s/%s", mainFolder, fileSha)
}

func Hash(file []byte) string {
	hash := sha256.Sum256(file)
	hashHex := hex.EncodeToString(hash[:])

	return hashHex
}

func SelectelUpload(file []byte, fileHash string) error {
	req, _ := http.NewRequest("GET", "https://api.selcdn.ru/auth/v1.0", nil)

	userName := fmt.Sprintf(
		"%d_%s",
		store.DefaultConfig.SelectelAccountID, store.DefaultConfig.SelectelUserName,
	)
	req.Header.Set("X-Auth-User", userName)
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

	url := fmt.Sprintf(
		"https://api.selcdn.ru/v1/SEL_%d/meme-files/%s/%s",
		store.DefaultConfig.SelectelAccountID, mainFolder, fileHash,
	)
	req, _ = http.NewRequest("PUT", url, bytes.NewReader(file))
	req.Header.Set("X-Auth-Token", authToken)
	req.Header.Set("Content-Type", "image/jpeg") // TODO mime types

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
