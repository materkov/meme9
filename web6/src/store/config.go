package store

type Config struct {
	SaveSecret      string
	AuthTokenSecret string

	VKAppID     int
	VKAppSecret string

	TelegramBotToken string
	TelegramChatID   int
}

var GlobalConfig = &Config{}
