package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
)

type SqlUserStore struct {
	DB *sql.DB
}

func (u *SqlUserStore) Get(ids []int) (map[int]*store.User, error) {
	rows, err := u.DB.Query(fmt.Sprintf("select id, data from objects where id in (%s)", utils.IdsToString(ids)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]*store.User{}
	for rows.Next() {
		objectID := 0
		var data []byte
		err = rows.Scan(&objectID, &data)
		if err != nil {
			return nil, err
		}

		object := store.User{}
		err = json.Unmarshal(data, &object)
		if err != nil {
			return nil, err
		}

		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

type UserStore interface {
	Get(ids []int) (map[int]*store.User, error)
}

type MockUserStore struct {
	objects map[int]*store.User
}

func (m *MockUserStore) Get(ids []int) (map[int]*store.User, error) {
	result := map[int]*store.User{}

	for _, objectID := range ids {
		result[objectID] = m.objects[objectID]
	}

	return result, nil
}
