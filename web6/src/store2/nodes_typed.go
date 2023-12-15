package store2

import (
	"encoding/json"
	"github.com/materkov/meme9/web6/src/store"
)

type TypedNodes struct {
	Store NodeStore
}

func (t *TypedNodes) GetUser(id int) (*store.User, error) {
	result, err := t.Store.Get([]int{id})
	if err != nil {
		return nil, err
	}

	userBytes := result[id]
	if userBytes == nil {
		return nil, ErrNotFound
	}

	user := store.User{}
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return &user, nil
}
