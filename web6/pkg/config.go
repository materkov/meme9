package pkg

type Config struct {
	SaveSecret string
}

var GlobalConfig = &Config{}
var BuildTime = "dev"
