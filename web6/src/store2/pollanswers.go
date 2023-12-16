package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
)

type SqlPollAnswerStore struct {
	DB *sql.DB
}

func (u *SqlPollAnswerStore) Get(ids []int) (map[int]*store.PollAnswer, error) {
	rows, err := u.DB.Query(fmt.Sprintf("select id, data from objects where id in (%s)", utils.IdsToString(ids)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]*store.PollAnswer{}
	for rows.Next() {
		objectID := 0
		var data []byte
		err = rows.Scan(&objectID, &data)
		if err != nil {
			return nil, err
		}

		object := store.PollAnswer{}
		err = json.Unmarshal(data, &object)
		if err != nil {
			return nil, err
		}

		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

type PollAnswerStore interface {
	Get(ids []int) (map[int]*store.PollAnswer, error)
}

type MockPollAnswerStore struct {
	objects map[int]*store.PollAnswer
}

func (m *MockPollAnswerStore) Get(ids []int) (map[int]*store.PollAnswer, error) {
	result := map[int]*store.PollAnswer{}

	for _, objectID := range ids {
		result[objectID] = m.objects[objectID]
	}

	return result, nil
}
