package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/meme9/web6/src/api"
	"github.com/materkov/meme9/web6/src/pkg/xlog"
	"github.com/materkov/meme9/web6/src/store"
	"log"
)

func main() {
	var err error
	store.SqlClient, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/meme9")
	if err != nil {
		log.Fatalf("Error opening mysql: %s", err)
	}

	store.GlobalConfig, err = store.GetConfig()
	if err != nil {
		log.Fatalf("Error reading config: %s", err)
	}

	go func() {
		xlog.ClearOldLogs()
	}()

	s := &api.HttpServer{
		Api: &api.API{},
	}
	s.Serve()
}
