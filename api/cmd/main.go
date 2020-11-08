package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/pkg/config"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	conf := config.Config{}

	configJson := []byte(os.Getenv("CONFIG"))
	if len(configJson) == 0 {
		homeDir, _ := os.UserHomeDir()
		configJson, _ = ioutil.ReadFile(homeDir + "/.meme")
	}

	err := json.Unmarshal(configJson, &conf)
	if err != nil {
		panic("Error parsing config: " + err.Error())
	}

	m := api.Main{Config: &conf}
	m.Main()
}
