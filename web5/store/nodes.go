package store

import (
	"database/sql"
	"encoding/json"
	"github.com/go-redis/redis/v9"
)

var RedisClient *redis.Client
var SqlClient *sql.DB
var DefaultConfig = Config{}

func NodeSave(id int, objType int, obj interface{}) error {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = SqlClient.Exec("insert into objects(id, obj_type, data) values (?, ?, ?) on duplicate key update data = ?", id, objType, objBytes, objBytes)
	if err != nil {
		return err
	}

	return nil
}
