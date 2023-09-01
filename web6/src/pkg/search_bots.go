package pkg

import "strings"

func IsSearchBot(userAgent string) bool {
	if strings.Contains(userAgent, "facebookexternalhit/") {
		return true
	}

	return false
}
