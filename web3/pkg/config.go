package pkg

type Config struct {
	VKAppID     int
	VKAppSecret string

	AuthTokenSecret string
}

var GlobalConfig = Config{}
