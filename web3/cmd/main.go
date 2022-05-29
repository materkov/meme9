package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/web3/store"
	"github.com/materkov/web3/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func gqlFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	body, _ := ioutil.ReadAll(r.Body)
	json.NewEncoder(w).Encode(runGQL(body))
}

func runGQL(req []byte) interface{} {
	var errors []error
	var fields = types.QueryParams{}

	json.Unmarshal(req, &fields)
	log.Printf("%+v", fields)

	result, err := types.ResolveQuery(fields)
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

	types.GlobalStore = &store.Store{DB: db}
	types.GlobalCachedStore = &store.CachedStore{
		Store:    types.GlobalStore,
		ObjCache: map[int]store.CachedItem{},
		Needed:   map[int]bool{},
	}

	r := runGQL([]byte(`
{
  "feed": {
    "include": true,
    "userId": 10,
    "inner": {
      "text": {
        "include": true,
        "maxLength": 19
      },
      "user": {
        "include": true,
        "inner": {
          "name": {
            "include": true
          }
        }
      }
    }
  }
}
`))
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "  ")
	w.Encode(r)

	//http.HandleFunc("/", httpFunc)
	http.HandleFunc("/gql", gqlFunc)

	http.ListenAndServe("127.0.0.1:8000", nil)
}
