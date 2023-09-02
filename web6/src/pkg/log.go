package pkg

import "log"

func LogErr(e error) {
	if e != nil {
		log.Printf("[ERROR] %s", e.Error())
	}
}
