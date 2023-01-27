package store

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type CachedObject interface {
	Post | User | Photo
}

type GenericCachedStore[T CachedObject] struct {
	cache   map[int]*T
	objType int
}

func (p *GenericCachedStore[T]) Preload(ids []int) {
	var neededIds []string
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := p.cache[id]; ok {
			continue
		}

		neededIds = append(neededIds, strconv.Itoa(id))
	}
	if len(neededIds) == 0 {
		return
	}

	query := "select id, data from objects where id in (%s) and obj_type = %d"
	query = fmt.Sprintf(query, strings.Join(neededIds, ","), p.objType)
	rows, err := SqlClient.Query(query)

	for rows.Next() {
		id := 0
		var data string
		_ = rows.Scan(&id, &data)

		// TODO: fixme
		data = strings.Replace(data, "\"ID\": 0", fmt.Sprintf("\"ID\":%d", id), 1)

		obj := new(T)
		err = json.Unmarshal([]byte(data), obj)
		if err != nil {
			log.Printf("Error unmarshaling obj: %s", err)
		}

		p.cache[id] = obj
	}

	for _, id := range ids {
		if _, ok := p.cache[id]; !ok {
			p.cache[id] = nil
		}
	}
}

func (p *GenericCachedStore[T]) Get(id int) *T {
	p.Preload([]int{id})
	return p.cache[id]
}
