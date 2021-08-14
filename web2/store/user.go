package store

import (
	"database/sql"
)

type User struct {
	ID   int
	Name string
}

type SqlUserStore struct {
	db *sql.DB
}

func (s *SqlUserStore) GetById(id int) (*User, error) {
	users, err := s.GetByIdMany([]int{id})
	return users[id], err
}

func (s *SqlUserStore) GetByIdMany(ids []int) (map[int]*User, error) {
	rows, err := s.db.Query("select id, name from user where id in (" + idsStr(ids) + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]*User{}
	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.ID, &u.Name)
		if err != nil {
			return nil, err
		}

		result[u.ID] = &u
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, err
}
