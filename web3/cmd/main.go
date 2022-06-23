package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/web3/pkg"
	"github.com/materkov/web3/store"
	"github.com/materkov/web3/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var GlobalStore store.Store

func gqlFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Content-Type", "application/json")

	authToken := r.Header.Get("Authorization")
	authToken = strings.TrimPrefix(authToken, "Bearer ")

	viewer := pkg.Viewer{}
	if authToken != "" {
		authToken, err := pkg.ParseAuthToken(authToken)
		if err == nil {
			viewer.UserID = authToken.UserID
		}
	}

	viewer.Origin = r.Header.Get("Origin")

	body, _ := ioutil.ReadAll(r.Body)
	json.NewEncoder(w).Encode(runGQL(viewer, body))
}

func runGQL(viewer pkg.Viewer, req []byte) interface{} {
	cachedStore := &store.CachedStore{
		Store:    GlobalStore,
		ObjCache: map[int]store.CachedItem{},
		Needed:   map[int]bool{},
	}

	var errors []error
	var fields = types.QueryParams{}

	json.Unmarshal(req, &fields)

	result, err := types.ResolveQuery(cachedStore, viewer, fields)
	if err != nil {
		errors = append(errors, err)
	}

	return result
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/meme9")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	GlobalStore = &store.SqlStore{DB: db}

	objects, err := GlobalStore.ObjGet([]int{store.ObjectIDConfig})
	if err != nil {
		log.Fatalf("Failed reading config: %s", err)
	}

	config, ok := objects[store.ObjectIDConfig].(*store.Config)
	if !ok {
		log.Fatalf("Error parsing config")
	}
	pkg.GlobalConfig = config

	http.HandleFunc("/gql", gqlFunc)

	log.Printf("Starting http server 127.0.0.1:8000")
	http.ListenAndServe("127.0.0.1:8000", nil)
}
