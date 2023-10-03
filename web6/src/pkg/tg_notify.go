package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/xlog"
	"github.com/materkov/meme9/web6/src/store"
	"io"
	"net/http"
)

func SendTelegramNotify(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", store.GlobalConfig.TelegramBotToken)
	body := struct {
		ChatID int    `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: store.GlobalConfig.TelegramChatID,
		Text:   text,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling json body: %w", err)
	}

	resp, err := http.DefaultClient.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	xlog.Log("Telegram notify", xlog.Fields{
		"req":    string(bodyBytes),
		"resp":   string(respBytes),
		"status": resp.StatusCode,
	})

	return nil
}

func SendTelegramNotifyAsync(text string) {
	go func() {
		err := SendTelegramNotify(text)
		if err != nil {
			xlog.Log("Telegram error", xlog.Fields{
				"err": err.Error(),
			})
		}
	}()
}
