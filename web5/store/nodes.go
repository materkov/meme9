package store

import (
	"database/sql"
	"encoding/json"
	"github.com/go-redis/redis/v9"
)

var RedisClient *redis.Client
var SqlClient *sql.DB
var DefaultConfig = Config{}

func NodeUpdate(id int, obj interface{}) error {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = SqlClient.Exec("update objects set data = ? where id = ?", objBytes, id)
	if err != nil {
		return err
	}

	return nil
}

func NodeInsert(objType int, obj interface{}) (int, error) {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return 0, err
	}

	result, err := SqlClient.Exec("insert into objects(obj_type, data) values (?, ?)", objType, objBytes)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
