package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	VkAppID       int
	VkAppSecret   string
	RequestScheme string
	RequestHost   string
	JwtSecret     string
}

var DefaultConfig *Config

func MustParseConfig() {
	configStr, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config: %s", err)
	}

	err = json.Unmarshal(configStr, &DefaultConfig)
	if err != nil {
		log.Fatalf("Error parsing config: %s", err)
	}
}
