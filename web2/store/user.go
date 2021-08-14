package store

import (
	"database/sql"
)

type User struct {
	ID   int
	Name string
	VkID int
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

func (s *SqlUserStore) GetByVkID(vkID int) (*User, error) {
	userID := 0
	err := s.db.QueryRow("select id from user where vk_id = ?", vkID).Scan(&userID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return s.GetById(userID)
}

func (s *SqlUserStore) Add(user *User) error {
	result, err := s.db.Exec("insert into user(name, vk_id) values (?, ?)", user.Name, sql.NullInt32{
		Int32: int32(user.VkID),
		Valid: user.VkID != 0,
	})
	if err != nil {
		return err
	}

	userID, _ := result.LastInsertId()
	user.ID = int(userID)

	return nil
}
