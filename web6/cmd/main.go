package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/meme9/web6/src/api"
	"github.com/materkov/meme9/web6/src/pkg/xlog"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
	"log"
)

func main() {
	var err error
	store.SqlClient, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/meme9")
	if err != nil {
		log.Fatalf("Error opening mysql: %s", err)
	}

	store.GlobalStore = &store.SqlStore{DB: store.SqlClient}

	nodeStore := store2.NewSqlNodes(store.SqlClient)
	store2.GlobalStore = &store2.Store{
		Unique:     store2.NewSqlUniqueStore(store.SqlClient),
		Nodes:      nodeStore,
		TypedNodes: &store2.TypedNodes{Store: nodeStore},
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
