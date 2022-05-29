package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func idsStr(ids []int) string {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.Itoa(id)
	}

	return strings.Join(idsStr, ",")
}

func getJSONObjects(table string, ids []int, db *sql.DB, result interface{}) error {
	if len(ids) == 0 {
		_ = json.Unmarshal([]byte("[]"), &result)
		return nil
	}

	query := fmt.Sprintf("select json_agg(obj.*) from \"%s\" obj where id in (%s)", table, idsStr(ids))

	var data []byte
	err := db.QueryRow(query).Scan(&data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("error unmarshaling json: %s", err)
	}

	return nil
}

func insertJSONObject() {

}
