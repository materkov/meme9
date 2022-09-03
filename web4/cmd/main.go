package main

import (
	"encoding/json"
	"github.com/materkov/meme9/web4/types"
	"os"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	dat, _ := os.ReadFile(homeDir + "/mypage/config.json")
	if len(dat) > 0 {
		_ = json.Unmarshal(dat, &types.DefaultConfig)
	}

	types.DoHandle()
}
