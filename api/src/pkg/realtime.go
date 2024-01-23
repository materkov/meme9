package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func PushRealtimeEvent(userID int, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling to json: %w", err)
	}

	dataEscaped := url.QueryEscape(string(dataBytes))
	resp, err := http.Get(fmt.Sprintf("http://localhost:8001/push?userId=%d&data=%s", userID, dataEscaped))
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
