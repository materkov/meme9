package pkg

import (
	"database/sql"
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
