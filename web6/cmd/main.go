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

	nodeStore := store2.NewSqlNodes(store.SqlClient)
	store2.GlobalStore = &store2.Store{
		Unique:      store2.NewSqlUniqueStore(store.SqlClient),
		Nodes:       nodeStore,
		TypedNodes:  &store2.TypedNodes{Store: nodeStore},
		Likes:       &store2.SqlLikes{DB: store.SqlClient},
		Subs:        &store2.SqlSubscriptions{DB: store.SqlClient},
		Wall:        &store2.SqlWall{DB: store.SqlClient},
		Votes:       &store2.SqlVotes{DB: store.SqlClient},
		Users:       &store2.SqlUserStore{DB: store.SqlClient},
		Posts:       &store2.SqlPostStore{DB: store.SqlClient},
		Polls:       &store2.SqlPollStore{DB: store.SqlClient},
		PollAnswers: &store2.SqlPollAnswerStore{DB: store.SqlClient},
		Tokens:      &store2.SqlTokenStore{DB: store.SqlClient},
		Configs:     &store2.SqlConfigStore{DB: store.SqlClient},
	}

	results, err := store2.GlobalStore.Configs.Get([]int{store.FakeObjConfig})
	if err != nil {
		log.Fatalf("Error reading config: %s", err)
	} else if results[store.FakeObjConfig] == nil {
		log.Fatalf("Config not found")
	}

	store.GlobalConfig = results[store.FakeObjConfig]

	go func() {
		xlog.ClearOldLogs()
	}()

	s := &api.HttpServer{
		Api: &api.API{},
	}
	s.Serve()
}
