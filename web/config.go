package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	VKAppID         int
	VKAppSecret     string
	VKMiniAppSecret string
	CSRFKey         string
}

func (c *Config) MustLoad() {
	configStr, _ := os.LookupEnv("CONFIG")
	if configStr == "" {
		configFile, _ := ioutil.ReadFile("config.json")
		configStr = string(configFile)
	}

	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		panic(err)
	}
}

var config Config
