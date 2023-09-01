package store

type Config struct {
	SaveSecret      string
	AuthTokenSecret string

	VKAppID     int
	VKAppSecret string
}

func GetConfig() (*Config, error) {
	obj := &Config{}
	err := getObject(5, ObjTypeConfig, obj)
	return obj, err
}
