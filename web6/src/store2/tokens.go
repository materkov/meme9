package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
)

type SqlTokenStore struct {
	DB *sql.DB
}

func (u *SqlTokenStore) Get(ids []int) (map[int]*store.Token, error) {
	rows, err := u.DB.Query(fmt.Sprintf("select id, data from objects where id in (%s)", utils.IdsToString(ids)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]*store.Token{}
	for rows.Next() {
		objectID := 0
		var data []byte
		err = rows.Scan(&objectID, &data)
		if err != nil {
			return nil, err
		}

		object := store.Token{}
		err = json.Unmarshal(data, &object)
		if err != nil {
			return nil, err
		}

		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

type TokenStore interface {
	Get(ids []int) (map[int]*store.Token, error)
}

type MockTokenStore struct {
	objects map[int]*store.Token
}

func (m *MockTokenStore) Get(ids []int) (map[int]*store.Token, error) {
	result := map[int]*store.Token{}

	for _, objectID := range ids {
		result[objectID] = m.objects[objectID]
	}

	return result, nil
}
