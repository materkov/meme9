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

	configEnv := os.Getenv("CONFIG")
	var err error
	if configEnv != "" {
		err = json.Unmarshal([]byte(configEnv), &conf)
	} else {
		homeDir, _ := os.UserHomeDir()
		dat, err := ioutil.ReadFile(homeDir + "/.meme")
		if err == nil {
			err = json.Unmarshal(dat, &conf)
		}
	}

	if err != nil {
		panic("Error parsing Config: " + err.Error())
	}

	m := api.Main{Config: &conf}
	m.Main()
}
