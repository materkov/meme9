package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/meme9/web6/api"
	"github.com/materkov/meme9/web6/pkg"
	"log"
)

func main() {
	var err error
	pkg.SqlClient, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/meme9")
	if err != nil {
		log.Fatalf("Error opening mysql: %s", err)
	}

	pkg.GlobalConfig, err = pkg.GetConfig()
	if err != nil {
		log.Fatalf("Error reading config: %s", err)
	}

	s := &api.HttpServer{}
	s.Serve()
}
