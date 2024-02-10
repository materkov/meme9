package store

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	VKAppID     int
	VKAppSecret string

	TelegramBotToken string
	TelegramChatID   int

	SelectelUploaderPassword string
	UploadTokenSecret        string
}

var GlobalConfig = &Config{}

func ParseConfig() error {
	file, err := os.ReadFile("/Users/m.materkov/projects/meme9/configs/api.json")
	if err != nil {
		file, err = os.ReadFile("/apps/meme9-config/api.json")
	}
	if err != nil {
		return fmt.Errorf("config not found")
	}

	err = json.Unmarshal(file, &GlobalConfig)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	return nil
}
