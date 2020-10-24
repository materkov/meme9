package store

func (s *Store) GetToken(nodeID int) (*Token, error) {
	node := &Token{}
	err := s.doGet(nodeID, node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (s *Store) AddToken(token *Token) error {
	return s.doAdd(token.ID, token)
}

func (s *Store) GetUser(nodeID int) (*User, error) {
	node := &User{}
	err := s.doGet(nodeID, node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (s *Store) AddUser(node *User) error {
	return s.doAdd(node.ID, node)
}

func (s *Store) GetPost(nodeID int) (*Post, error) {
	node := &Post{}
	err := s.doGet(nodeID, node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (s *Store) AddPost(node *Post) error {
	return s.doAdd(node.ID, node)
}
