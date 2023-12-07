package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"strings"
	"time"
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

func (s *SqlStore) GetObjectsMany(ids []int) (map[int][]byte, error) {
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

func (s *SqlStore) AddObject(objType int, object interface{}) (int, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return 0, fmt.Errorf("error marshaling to json: %w", err)
	}

	res, err := s.DB.Exec("insert into objects(obj_type, data) values (?, ?)", objType, data)
	if err != nil {
		return 0, fmt.Errorf("error inserting row: %w", err)
	}

	objId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last id from mysql: %w", err)
	}

	return int(objId), nil
}

func (s *SqlStore) AddEdge(fromID, toID, edgeType int) error {
	_, err := s.DB.Exec(
		"insert into edges(from_id, to_id, edge_type, date) values (?, ?, ?, ?)",
		fromID, toID, edgeType, time.Now().Unix(),
	)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
			return ErrDuplicateEdge
		}
		return fmt.Errorf("error inserting edge: %s", err)
	}

	return nil
}

func (s *SqlStore) GetEdge(fromID, toID, edgeType int) (*Edge, error) {
	date := 0
	err := s.DB.QueryRow("select date from edges where from_id = ? and to_id = ? and edge_type = ?", fromID, toID, edgeType).Scan(&date)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEdge
	} else if err != nil {
		return nil, err
	}

	e := &Edge{
		FromID: fromID,
		ToID:   toID,
		Date:   date,
	}
	return e, nil
}

func (s *SqlStore) CountEdges(fromID, edgeType int) (int, error) {
	cnt := 0
	err := s.DB.QueryRow("select count(*) from edges where from_id = ? and edge_type = ?", fromID, edgeType).Scan(&cnt)
	return cnt, err
}

func (s *SqlStore) LoadLikesMany(postIds []int, viewerID int) (counters map[int]int, isLiked map[int]bool, err error) {
	query := `
select from_id, count(*), sum(to_id = %d)
from edges
where from_id in (%s) and edge_type=%d
group by from_id
`
	rows, err := s.DB.Query(fmt.Sprintf(query, viewerID, strings.Join(utils.IdsToStrings(postIds), ","), EdgeTypeLiked))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	isLiked = map[int]bool{}
	counters = map[int]int{}

	for rows.Next() {
		postID, count, isLikedInt := 0, 0, 0
		err = rows.Scan(&postID, &count, &isLikedInt)
		if err != nil {
			return nil, nil, err
		}

		counters[postID] = count
		if isLikedInt > 0 {
			isLiked[postID] = true
		}
	}

	return counters, isLiked, nil
}

func (s *SqlStore) GetEdges(fromID int, edgeType int, limit int, startFrom int) ([]Edge, error) {
	rows, err := s.DB.Query("select to_id, date from edges where from_id = ? and edge_type = ? and to_id < ? order by id desc limit ?", fromID, edgeType, startFrom, limit)
	if err != nil {
		return nil, fmt.Errorf("error selecting rows: %w", err)
	}
	defer rows.Close()

	var results []Edge
	for rows.Next() {
		e := Edge{FromID: fromID}
		err = rows.Scan(&e.ToID, &e.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning edge row: %w", err)
		}

		results = append(results, e)
	}

	return results, err
}

func (s *SqlStore) DelEdge(fromID, toID, edgeType int) error {
	_, err := s.DB.Exec("delete from edges where from_id = ? and edge_type = ? and to_id = ?", fromID, edgeType, toID)
	if err != nil {
		return fmt.Errorf("error deleteing edge: %w", err)
	}

	return nil
}

func (s *SqlStore) UpdateObject(object interface{}, id int) error {
	data, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("error marshaling to json: %w", err)
	}

	_, err = s.DB.Exec("update objects set data = ? where id = ?", data, id)
	if err != nil {
		return fmt.Errorf("error updating row: %w", err)
	}

	return nil
}

func (s *SqlStore) GetUnique(keyType int, key string) (int, error) {
	if key == "" {
		return 0, ErrUniqueNotFound
	}

	objectID := 0
	err := s.DB.QueryRow("select object_id from uniques where type = ? and `key` = ?", keyType, key).Scan(&objectID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrUniqueNotFound
	} else if err != nil {
		return 0, fmt.Errorf("error selecing unique row: %w", err)
	}

	return objectID, nil
}

func (s *SqlStore) AddUnique(keyType int, key string, objectID int) error {
	_, err := s.DB.Exec("insert into uniques(type, `key`, object_id) values (?, ?, ?)", keyType, key, objectID)
	if err != nil {
		return fmt.Errorf("error inserting unique row: %w", err)
	}
	return nil
}
