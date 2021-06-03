package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	VKAppID         int
	VKAppSecret     string
	VKMiniAppSecret string
	CSRFKey         string
	AWSKeyID        string
	AWSKeySecret    string
}

func (c *Config) Load() error {
	configStr, _ := os.LookupEnv("CONFIG")
	if configStr == "" {
		configFile, _ := ioutil.ReadFile("config.json")
		configStr = string(configFile)
	}

	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return fmt.Errorf("failed parsing json: %w", err)
	}

	return nil
}

var config Config
