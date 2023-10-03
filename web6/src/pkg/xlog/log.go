package xlog

import (
	"encoding/json"
	"github.com/materkov/meme9/web6/src/store"
)

type Fields map[string]interface{}

func Log(message string, fields Fields) {
	go func() {
		fieldsBytes, _ := json.Marshal(fields)
		_, _ = store.SqlClient.Exec("insert into log(dt, message, params, file) values (now(), ?, ?, ?)", message, fieldsBytes, "")
	}()
}

func ClearOldLogs() {
	_, _ = store.SqlClient.Exec("delete from log where dt <= date_sub(now(), interval 3 day)")
}
