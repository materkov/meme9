package store

type Config struct {
	SaveSecret      string
	AuthTokenSecret string

	VKAppID     int
	VKAppSecret string
}

var GlobalConfig = &Config{}
