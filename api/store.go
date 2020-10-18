package api

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/tinylib/msgp/msgp"
)

var ErrNodeNotFound = fmt.Errorf("node not found")

type Store struct {
	redis *redis.Client
}

func NewStore(redis *redis.Client) *Store {
	return &Store{redis: redis}
}

func (s *Store) doGet(nodeID int, node msgp.Unmarshaler) error {
	nodeSerialized, err := s.redis.Get(fmt.Sprintf("node:%d", nodeID)).Result()
	if err == redis.Nil {
		return ErrNodeNotFound
	} else if err != nil {
		return fmt.Errorf("error selecting key from redis: %w", err)
	}

	_, err = node.UnmarshalMsg([]byte(nodeSerialized))
	if err != nil {
		return fmt.Errorf("error unserializing node: %w", err)
	}

	return nil
}

func (s *Store) doAdd(id int, node msgp.Marshaler) error {
	nodeMarshaled, err := node.MarshalMsg(nil)
	if err != nil {
		return fmt.Errorf("error marshaling token: %w", err)
	}

	_, err = s.redis.Set(fmt.Sprintf("node:%d", id), nodeMarshaled, 0).Result()
	if err != nil {
		return fmt.Errorf("error saving token to redis: %w", err)
	}

	return nil
}

func (s *Store) GenerateNodeID() (int, error) {
	nodeID, err := s.redis.Incr("node_ids").Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementing redis key: %w", err)
	}

	return int(nodeID), nil
}

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
