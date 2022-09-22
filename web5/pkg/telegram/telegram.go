package telegram

import (
	"fmt"
	"github.com/materkov/meme9/web5/store"
	"io"
	"log"
	"net/http"
	"net/url"
)

func SendNotify(text string) error {
	resp, err := http.PostForm(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", store.DefaultConfig.TelegramToken), url.Values{
		"chat_id": []string{"7952464"},
		"text":    []string{text},
	})
	if err != nil {
		return fmt.Errorf("HTTP error: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("%d", resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Body: %s", body)

	return nil
}
