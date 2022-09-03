package types

type Config struct {
	VKAppID     int
	VKAppSecret string
}

var DefaultConfig = Config{}

type Viewer struct {
	UserID int

	Origin string
}
