package store

type Config struct {
	VKAppID     int
	VKAppSecret string

	TelegramBotToken string
	TelegramChatID   int
}

var GlobalConfig = &Config{}
