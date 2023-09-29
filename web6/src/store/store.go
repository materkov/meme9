package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
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

	EdgeTypeFollowing  = 7
	EdgeTypeFollowedBy = 8
)

var SqlClient *sql.DB

var (
	ErrObjectNotFound = fmt.Errorf("object not found")
	ErrDuplicateEdge  = fmt.Errorf("edge duplicate")
)

func getObject(id int, objType int, obj interface{}) error {
	var data []byte
	err := SqlClient.QueryRow("select data from objects where id = ? and obj_type = ?", id, objType).Scan(&data)
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

func UpdateObject(object interface{}, id int) error {
	data, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("error marshaling to json: %w", err)
	}

	_, err = SqlClient.Exec("update objects set data = ? where id = ?", data, id)
	if err != nil {
		return fmt.Errorf("error updating row: %w", err)
	}

	return nil
}

func AddObject(objType int, object interface{}) (int, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return 0, fmt.Errorf("error marshaling to json: %w", err)
	}

	res, err := SqlClient.Exec("insert into objects(obj_type, data) values (?, ?)", objType, data)
	if err != nil {
		return 0, fmt.Errorf("error inserting row: %w", err)
	}

	objId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last id from mysql: %w", err)
	}

	return int(objId), nil
}

type Edge struct {
	FromID int
	ToID   int
	Date   int
}

func GetEdges(fromID int, edgeType int) ([]Edge, error) {
	rows, err := SqlClient.Query("select to_id, date from edges where from_id = ? and edge_type = ? order by id desc", fromID, edgeType)
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

func GetToId(edges []Edge) []int {
	result := make([]int, len(edges))
	for i, edge := range edges {
		result[i] = edge.ToID
	}
	return result
}

var ErrNoEdge = fmt.Errorf("no edge")

func GetEdge(fromID, toID, edgeType int) (Edge, error) {
	e := Edge{
		FromID: fromID,
		ToID:   toID,
	}

	err := SqlClient.QueryRow("select date from edges where from_id = ? and to_id = ? and edge_type = ?", fromID, toID, edgeType).Scan(&e.Date)
	if errors.Is(err, sql.ErrNoRows) {
		return Edge{}, ErrNoEdge
	} else if err != nil {
		return Edge{}, err
	}

	return e, nil
}

func AddEdge(fromID, toID, edgeType int, uniqueKey string) error {
	_, err := SqlClient.Exec(
		"insert into edges(from_id, to_id, edge_type, unique_key, date) values (?, ?, ?, ?, ?)",
		fromID, toID, edgeType, sql.NullString{String: uniqueKey, Valid: uniqueKey != ""}, time.Now().Unix(),
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

func DelEdge(fromID, toID, edgeType int) error {
	_, err := SqlClient.Exec("delete from edges where from_id = ? and edge_type = ? and to_id = ?", fromID, edgeType, toID)
	if err != nil {
		return fmt.Errorf("error deleteing edge: %w", err)
	}

	return nil
}
