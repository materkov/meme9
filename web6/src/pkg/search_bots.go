package pkg

import "strings"

func IsSearchBot(userAgent string) bool {
	return strings.Contains(userAgent, "facebookexternalhit/")
}
