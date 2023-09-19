package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	FakeObjPostedPost = -1
	FakeObjVkAuth     = -2
	FakeObjEmailAuth  = -3
	FakeObjToken      = -4

	ObjTypeArticle = 1 // not used
	ObjTypeConfig  = 2
	ObjTypeUser    = 3
	ObjTypePost    = 4
	ObjTypeToken   = 5

	EdgeTypePosted     = 1
	EdgeTypeLastPosted = 2 // not used
	EdgeTypePostedPost = 3
	EdgeTypeVkAuth     = 4 // not used
	EdgeTypeEmailAuth  = 5 // not used
	EdgeTypeToken      = 6 // not used
)

var SqlClient *sql.DB

var ErrObjectNotFound = fmt.Errorf("object not found")

func getObject(id int, objType int, obj interface{}) error {
	var data []byte
	err := SqlClient.QueryRow("select data from objects where id = ? and obj_type = ?", id, objType).Scan(&data)
	if err == sql.ErrNoRows {
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

func UpdateObject(object interface{}, id int) error {
	data, _ := json.Marshal(object)
	_, err := SqlClient.Exec("update objects set data = ? where id = ?", data, id)
	if err != nil {
		return fmt.Errorf("error updating row: %w", err)
	}

	return nil
}

func AddObject(objType int, object interface{}) (int, error) {
	data, _ := json.Marshal(object)
	res, err := SqlClient.Exec("insert into objects(obj_type, data) values (?, ?)", objType, data)
	if err != nil {
		return 0, fmt.Errorf("error inserting row: %w", err)
	}

	objId, _ := res.LastInsertId()

	return int(objId), nil
}

func GetEdges(fromID int, edgeType int) ([]int, error) {
	rows, err := SqlClient.Query("select to_id from edges where from_id = ? and edge_type = ? order by id desc", fromID, edgeType)
	if err != nil {
		return nil, fmt.Errorf("error selecting rows: %w", err)
	}
	defer rows.Close()

	var results []int
	for rows.Next() {
		objID := 0
		err = rows.Scan(&objID)
		if err != nil {
			return nil, fmt.Errorf("error scanning edge row: %w", err)
		}

		results = append(results, objID)
	}

	return results, err
}

func AddEdge(fromID, toID, edgeType int, uniqueKey string) error {
	_, err := SqlClient.Exec(
		"insert into edges(from_id, to_id, edge_type, unique_key) values (?, ?, ?, ?)",
		fromID, toID, edgeType, sql.NullString{String: uniqueKey, Valid: uniqueKey != ""},
	)
	if err != nil {
		return fmt.Errorf("error inserting edge: %s", err)
	}

	return nil
}

func GetEdgeByUniqueKey(fromID int, edgeType int, uniqueKey string) (int, error) {
	toID := 0
	err := SqlClient.QueryRow("select to_id from edges where from_id = ? and edge_type = ? and unique_key = ? limit 1", fromID, edgeType, uniqueKey).Scan(&toID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("error selecting row: %w", err)
	}

	return toID, nil
}

func DelEdge(fromID, edgeType, toID int) error {
	_, err := SqlClient.Exec("delete from edges where from_id = ? and edge_type = ? and to_id = ?", fromID, edgeType, toID)
	if err != nil {
		return fmt.Errorf("error deleteing edge: %w", err)
	}

	return nil
}
