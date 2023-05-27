package pkg

import "fmt"

func GetEdges(fromID int, edgeType int) ([]int, error) {
	rows, err := SqlClient.Query("select to_id from edges where from_id = ? and edge_type = ?", fromID, edgeType)
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
