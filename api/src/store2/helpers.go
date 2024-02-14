package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/api/src/pkg/utils"
)

func LoadObjects(db *sql.DB, ids []int) (map[int][]byte, error) {
	if len(ids) == 0 {
		return make(map[int][]byte), nil
	}

	idsStr := utils.IdsToCommaSeparated(ids)

	query := "select id, data from objects where id in (%s)"
	query = fmt.Sprintf(query, idsStr)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int][]byte{}
	for rows.Next() {
		objectID := 0
		var data []byte
		err = rows.Scan(&objectID, &data)
		if err != nil {
			return nil, err
		}

		result[objectID] = data
	}

	return result, nil
}

func AddObject(db *sql.DB, object interface{}, objectType int) (int, error) {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return 0, fmt.Errorf("error marshaling object to json: %w", err)
	}

	result, err := db.Exec("insert into objects(obj_type, data) values (?, ?)", objectType, objectBytes)
	if err != nil {
		return 0, fmt.Errorf("error inserting object: %w", err)
	}

	objectID, _ := result.LastInsertId()
	return int(objectID), nil
}
