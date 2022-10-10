package utils

import (
	"log"
	"math/rand"
)

func RandString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func LogIfErr(err error) {
	if err != nil {
		log.Printf("[ERROR] Error: %s", err)
	}
}
