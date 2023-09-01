package api

import (
	"encoding/json"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/store"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	ID   string
	Name string
}

func transformUser(userID int, user *store.User) *User {
	result := &User{
		ID: strconv.Itoa(userID),
	}
	if user == nil {
		return result
	}

	result.Name = user.Name

	return result
}

func (*HttpServer) usersList(w http.ResponseWriter, r *http.Request, t *pkg.AuthToken) (interface{}, error) {
	req := struct {
		UserIds []string
	}{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, &Error{400, "error parsing request"}
	}

	result := make([]*User, len(req.UserIds))
	for i, userIdStr := range req.UserIds {
		userId, _ := strconv.Atoi(userIdStr)
		user, err := store.GetUser(userId)
		if err != nil {
			log.Printf("[ERROR] Error loading user: %s", err)
		}

		result[i] = transformUser(userId, user)
	}

	return result, nil
}
