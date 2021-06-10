package main

import (
	"database/sql"
	"math/rand"
	"strconv"
	"strings"
)

func RandString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func idsStr(ids []int) string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.Itoa(id)
	}

	return strings.Join(result, ",")
}

func scanIdsList(db *sql.DB, query string) ([]int, error) {
	rows, err := db.Query(query)
	if err != nil {
	    return nil, err
	}
	defer rows.Close()

	result := make([]int, 0)
	for rows.Next() {
		item := 0
		if err = rows.Scan(&item); err != nil {
		    return nil, err
		}

		result = append(result, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, err
}
