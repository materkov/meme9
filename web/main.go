package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/meme9/web/app"
	"github.com/materkov/meme9/web/httpserver"
	"github.com/materkov/meme9/web/store"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := app.DefaultConfig.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	db, err := sql.Open("mysql", "root:root@/meme9")
	if err != nil {
		log.Fatalf("failed opening mysql connection: %s", err)
	}

	app.ObjectStore = store.NewObjectStore(db)

	app.Main()

	httpSrv := &httpserver.HttpServer{
		Store:    app.ObjectStore,
		FeedSrv:  app.FeedSrv,
		UtilsSrv: app.UtilsSrv,
	}
	httpSrv.Serve()

}
