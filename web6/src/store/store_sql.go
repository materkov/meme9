package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/tracer"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"strings"
)

type SqlStore struct {
	DB *sql.DB
}

func (s *SqlStore) getObject(id int, objType int, obj interface{}) error {
	var data []byte
	err := s.DB.QueryRow("select data from objects where id = ? and obj_type = ?", id, objType).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrObjectNotFound
	} else if err != nil {
		return fmt.Errorf("error selecting database: %w", err)
	}

	err = json.Unmarshal(data, obj)
	if err != nil {
		return fmt.Errorf("error unmarshaling object %d: %w", objType, err)
	}

	return nil

}

func (s *SqlStore) GetObjectsMany(ctx context.Context, ids []int) (map[int][]byte, error) {
	span := tracer.FromCtx(ctx).StartChild("GetObjectsMany")
	defer span.Stop()

	span.Tags["ids"] = strings.Join(utils.IdsToStrings(ids), ",")

	if len(ids) == 0 {
		return map[int][]byte{}, nil
	}

	rows, err := s.DB.Query(fmt.Sprintf("select id, data from objects where id in (%s)", strings.Join(utils.IdsToStrings(ids), ",")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultMap := map[int][]byte{}
	for rows.Next() {
		id := 0
		var data []byte
		err = rows.Scan(&id, &data)
		if err != nil {
			return nil, err
		}

		resultMap[id] = data
	}

	return resultMap, nil
}
