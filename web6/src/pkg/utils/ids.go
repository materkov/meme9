package utils

import "strconv"

func IdsToStrings(ids []int) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.Itoa(id)
	}

	return result
}

func IdsToInts(ids []string) []int {
	result := make([]int, len(ids))
	for i, id := range ids {
		result[i], _ = strconv.Atoi(id)
	}

	return result
}

func UniqueIds(ids []int) []int {
	uniqueMap := map[int]bool{}
	for _, id := range ids {
		uniqueMap[id] = true
	}

	result := make([]int, len(uniqueMap))
	idx := 0
	for userID := range uniqueMap {
		result[idx] = userID
		idx++
	}

	return result
}
