package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
)

type SqlPollStore struct {
	DB *sql.DB
}

func (u *SqlPollStore) Get(ids []int) (map[int]*store.Poll, error) {
	rows, err := u.DB.Query(fmt.Sprintf("select id, data from objects where id in (%s)", utils.IdsToString(ids)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]*store.Poll{}
	for rows.Next() {
		objectID := 0
		var data []byte
		err = rows.Scan(&objectID, &data)
		if err != nil {
			return nil, err
		}

		object := store.Poll{}
		err = json.Unmarshal(data, &object)
		if err != nil {
			return nil, err
		}

		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

type PollStore interface {
	Get(ids []int) (map[int]*store.Poll, error)
}

type MockPollStore struct {
	objects map[int]*store.Poll
}

func (m *MockPollStore) Get(ids []int) (map[int]*store.Poll, error) {
	result := map[int]*store.Poll{}

	for _, objectID := range ids {
		result[objectID] = m.objects[objectID]
	}

	return result, nil
}
