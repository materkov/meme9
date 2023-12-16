package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
)

type SqlPostStore struct {
	DB *sql.DB
}

func (u *SqlPostStore) Get(ids []int) (map[int]*store.Post, error) {
	rows, err := u.DB.Query(fmt.Sprintf("select id, data from objects where id in (%s)", utils.IdsToString(ids)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]*store.Post{}
	for rows.Next() {
		objectID := 0
		var data []byte
		err = rows.Scan(&objectID, &data)
		if err != nil {
			return nil, err
		}

		object := store.Post{}
		err = json.Unmarshal(data, &object)
		if err != nil {
			return nil, err
		}

		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

type PostStore interface {
	Get(ids []int) (map[int]*store.Post, error)
}

type MockPostStore struct {
	objects map[int]*store.Post
}

func (m *MockPostStore) Get(ids []int) (map[int]*store.Post, error) {
	result := map[int]*store.Post{}

	for _, objectID := range ids {
		result[objectID] = m.objects[objectID]
	}

	return result, nil
}
