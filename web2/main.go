package main

import (
	"database/sql"
	"github.com/materkov/meme9/web2/lib"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/meme9/web2/controller"
	"github.com/materkov/meme9/web2/store"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	lib.MustParseConfig()

	db, err := sql.Open("mysql", "root:root@/meme9")
	if err != nil {
		log.Fatal(err)
	}

	st := store.NewStore(db)
	srv := controller.Server{
		Store: st,
	}

	srv.Serve()
}
