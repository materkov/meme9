package pkg

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/api/src/store"
	"image"
	_ "image/jpeg"
	"io"
	"net/http"
	"strings"
)

func SelectelUpload(file []byte, fileHash string) error {
	body := `{"auth":{"identity":{"methods":["password"],"password":{"user":{"name":"Meme-files-uploader","domain":{"name":"200213"},"password":"%s"}}},"scope":{"project":{"name":"My First Project","domain":{"name":"200213"}}}}}`
	body = fmt.Sprintf(body, store.GlobalConfig.SelectelUploaderPassword)

	authReq, err := http.NewRequest("POST", "https://cloud.api.selcloud.ru/identity/v3/auth/tokens", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating http request: %w", err)
	}
	authReq.Header.Set("Content-Type", "application/json")

	authResp, err := http.DefaultClient.Do(authReq)
	if err != nil {
		return fmt.Errorf("error doing http auth request: %w", err)
	}

	authBody, _ := io.ReadAll(authResp.Body)
	_ = authResp.Body.Close()

	authToken := authResp.Header.Get("x-subject-token")
	if authToken == "" {
		return fmt.Errorf("no x-subject-token header in response: %d, %s", authResp.StatusCode, authBody)
	}

	uploadReq, err := http.NewRequest("PUT", "https://swift.ru-1.storage.selcloud.ru/v1/7061fb856af24df0be194a8f8bebb303/meme-files/b/"+fileHash, bytes.NewReader(file))
	if err != nil {
		return fmt.Errorf("error creating upload http request: %w", err)
	}

	uploadReq.Header.Set("content-type", "image/jpeg")
	uploadReq.Header.Set("x-auth-token", authToken)

	uploadResp, err := http.DefaultClient.Do(uploadReq)
	if err != nil {
		return fmt.Errorf("error doing upload http request: %w", err)
	}

	uploadBody, _ := io.ReadAll(uploadResp.Body)
	if uploadResp.StatusCode != 201 {
		return fmt.Errorf("bad upload http response code: %d, %s", uploadResp.StatusCode, uploadBody)
	}

	return nil
}

func GetFileHash(file []byte) string {
	h := sha256.New()
	h.Write(file)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetFilePath(fileHash string) string {
	return fmt.Sprintf("https://03e51a90-4d56-4906-80cc-edda75797b32.selstorage.ru/b/%s", fileHash)
}

func ValidatePhoto(file []byte) (int, int, error) {
	im, format, err := image.DecodeConfig(bytes.NewReader(file))
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse photo meta: %w", err)
	} else if format != "jpeg" {
		return 0, 0, fmt.Errorf("only jpeg pictures are allowed")
	}

	return im.Width, im.Height, nil
}

type UploadToken struct {
	Hash   string
	Width  int
	Height int
	Size   int
}

func (u *UploadToken) ToString() string {
	tokenStr, _ := json.Marshal(u)

	h := hmac.New(sha256.New, []byte(store.GlobalConfig.UploadTokenSecret))
	h.Write(tokenStr)

	return fmt.Sprintf("%x.%s", h.Sum(nil), tokenStr)
}
