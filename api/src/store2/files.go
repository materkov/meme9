package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/api/src/pkg/utils"
	"github.com/materkov/meme9/api/src/store"
)

type SqlFileStore struct {
	DB *sql.DB
}

func (s *SqlFileStore) Get(ids []int) (map[int]*store.File, error) {
	if len(ids) == 0 {
		return make(map[int]*store.File), nil
	}

	rows, err := s.DB.Query(fmt.Sprintf("select id, data from objects where id in (%s)", utils.IdsToCommaSeparated(ids)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]*store.File{}
	for rows.Next() {
		objectID := 0
		var data []byte
		err = rows.Scan(&objectID, &data)
		if err != nil {
			return nil, err
		}

		object := store.File{}
		err = json.Unmarshal(data, &object)
		if err != nil {
			return nil, err
		}

		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

func (s *SqlFileStore) Add(object *store.File) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	result, err := s.DB.Exec("insert into objects(obj_type, data) values (?, ?)", store.ObjTypeFile, objectBytes)
	if err != nil {
		return err
	}

	objectID, _ := result.LastInsertId()
	object.ID = int(objectID)

	return nil
}
