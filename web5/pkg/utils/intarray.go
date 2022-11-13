package utils

import "strconv"

func StrToIntArray(ids []string) []int {
	result := make([]int, 0)
	for _, id := range ids {
		parsedID, _ := strconv.Atoi(id)
		if parsedID > 0 {
			result = append(result, parsedID)
		}
	}

	return result
}

func IntToStrArray(ids []int) []string {
	result := make([]string, 0)
	for _, id := range ids {
		if id > 0 {
			result = append(result, strconv.Itoa(id))
		}
	}

	return result
}
