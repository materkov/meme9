package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
)

type SqlConfigStore struct {
	DB *sql.DB
}

func (u *SqlConfigStore) Get(ids []int) (map[int]*store.Config, error) {
	rows, err := u.DB.Query(fmt.Sprintf("select id, data from objects where id in (%s)", utils.IdsToCommaSeparated(ids)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]*store.Config{}
	for rows.Next() {
		objectID := 0
		var data []byte
		err = rows.Scan(&objectID, &data)
		if err != nil {
			return nil, err
		}

		object := store.Config{}
		err = json.Unmarshal(data, &object)
		if err != nil {
			return nil, err
		}

		//object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

type ConfigStore interface {
	Get(ids []int) (map[int]*store.Config, error)
}

type MockConfigStore struct {
	Objects map[int]*store.Config
}

func (m *MockConfigStore) Get(ids []int) (map[int]*store.Config, error) {
	result := map[int]*store.Config{}

	for _, objectID := range ids {
		result[objectID] = m.Objects[objectID]
	}

	return result, nil
}
