package pkg

import (
	"database/sql"
	"errors"
	"fmt"
)

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
