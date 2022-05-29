package store

func (s *SqlObjectStore) GetUser(id int) (*User, error) {
	objects, err := s.GetByIdMany([]int{id})
	if err != nil {
		return nil, err
	}

	if len(objects) == 0 {
		return nil, nil
	}

	user, ok := objects[0].(*User)
	if !ok {
		return nil, nil
	}

	return user, err
}

func (s *SqlObjectStore) GetUsers(ids []int) ([]*User, error) {
	objects, err := s.GetByIdMany(ids)
	if err != nil {
		return nil, err
	}

	var result []*User
	for _, object := range objects {
		user, ok := object.(*User)
		if ok {
			result = append(result, user)
		}
	}

	return result, nil
}

func (s *SqlObjectStore) GetPost(id int) (*Post, error) {
	objects, err := s.GetByIdMany([]int{id})
	if err != nil {
		return nil, err
	}

	if len(objects) == 0 {
		return nil, nil
	}

	user, ok := objects[0].(*Post)
	if !ok {
		return nil, nil
	}

	return user, err
}
