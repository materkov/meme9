package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/server"
	"github.com/materkov/meme9/api/src/pkg/xlog"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/twitchtv/twirp"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/meme9")
	if err != nil {
		log.Fatalf("Error opening mysql: %s", err)
	}

	store2.GlobalStore = &store2.Store{
		Unique:      &store2.SqlUniqueStore{DB: db},
		Likes:       &store2.SqlLikes{DB: db},
		Subs:        &store2.SqlSubscriptions{DB: db},
		Wall:        &store2.SqlWall{DB: db},
		Votes:       &store2.SqlVotes{DB: db},
		Users:       &store2.SqlUserStore{DB: db},
		Posts:       &store2.SqlPostStore{DB: db},
		Polls:       &store2.SqlPollStore{DB: db},
		PollAnswers: &store2.SqlPollAnswerStore{DB: db},
		Tokens:      &store2.SqlTokenStore{DB: db},
		Configs:     &store2.SqlConfigStore{DB: db},
	}

	store.SqlClient = db

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

	postsSrv := api.NewPostsServer(&server.PostsServer{}, twirp.WithServerJSONSkipDefaults(true))
	authSrv := api.NewAuthServer(&server.AuthServer{}, twirp.WithServerJSONSkipDefaults(true))
	pollSrv := api.NewPollsServer(&server.PollServer{}, twirp.WithServerJSONSkipDefaults(true))
	userSrv := api.NewUsersServer(&server.UserServer{}, twirp.WithServerJSONSkipDefaults(true))

	http.Handle(postsSrv.PathPrefix(), server.AuthMiddleware(postsSrv))
	http.Handle(authSrv.PathPrefix(), server.AuthMiddleware(authSrv))
	http.Handle(pollSrv.PathPrefix(), server.AuthMiddleware(pollSrv))
	http.Handle(userSrv.PathPrefix(), server.AuthMiddleware(userSrv))

	http.ListenAndServe("127.0.0.1:8002", nil)
}
