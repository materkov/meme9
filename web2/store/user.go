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

func (s *SqlUserStore) GetByVkID(vkID int) (int, error) {
	userID := 0
	err := s.db.QueryRow(`select user_id from vk_id where vk_id.vk_id = $1`, vkID).Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *SqlUserStore) Add(user *User) error {
	return s.db.QueryRow("insert into \"user\"(name, vkId) values ($1, $2) returning id", user.Name, sql.NullInt32{
		Int32: int32(user.VkID),
		Valid: user.VkID != 0,
	}).Scan(&user.ID)
}
