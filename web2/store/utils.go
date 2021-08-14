package store

import (
	"strconv"
	"strings"
)

func idsStr(ids []int) string {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.Itoa(id)
	}

	return strings.Join(idsStr, ",")
}
