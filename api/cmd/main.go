package main

import (
	"math/rand"
	"time"

	"github.com/materkov/meme9/api"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	m := api.Main{}
	m.Main()
}
