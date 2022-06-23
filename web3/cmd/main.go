package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/web3/pkg"
	"github.com/materkov/web3/store"
	"github.com/materkov/web3/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func gqlFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

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
		Store:    pkg.GlobalStore,
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

func handleUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	authTokenStr := r.Header.Get("Authorization")
	authTokenStr = strings.TrimPrefix(authTokenStr, "Bearer ")

	if authTokenStr == "" {
		return
	}

	authToken, err := pkg.ParseAuthToken(authTokenStr)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = pkg.UpdateAvatar(body, authToken.UserID)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "OK")
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/meme9")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pkg.GlobalStore = &store.SqlStore{DB: db}

	objects, err := pkg.GlobalStore.ObjGet([]int{store.ObjectIDConfig})
	if err != nil {
		log.Fatalf("Failed reading config: %s", err)
	}

	config, ok := objects[store.ObjectIDConfig].(*store.Config)
	if !ok {
		log.Fatalf("Error parsing config")
	}
	pkg.GlobalConfig = config

	http.HandleFunc("/gql", gqlFunc)
	http.HandleFunc("/upload", handleUpload)

	log.Printf("Starting http server 127.0.0.1:8000")
	http.ListenAndServe("127.0.0.1:8000", nil)
}
