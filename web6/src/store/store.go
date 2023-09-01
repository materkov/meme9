package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

const (
	FakeObjPostedPost = -1
	FakeObjVkAuth     = -2

	ObjTypeArticle = 1 // not used
	ObjTypeConfig  = 2
	ObjTypeUser    = 3
	ObjTypePost    = 4

	EdgeTypePosted     = 1
	EdgeTypeLastPosted = 2 // not used
	EdgeTypePostedPost = 3
	EdgeTypeVkAuth     = 4
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

var GlobalConfig = &Config{}
