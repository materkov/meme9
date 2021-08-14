package main

import (
	"database/sql"
	"github.com/materkov/meme9/web2/lib"
	"log"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/meme9/web2/controller"
	"github.com/materkov/meme9/web2/store"
)

func main() {
	lib.MustParseConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

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
