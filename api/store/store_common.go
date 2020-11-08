package store

import (
	"fmt"
	"strconv"

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

func (s *Store) GetUserByVKID(vkID int) (int, error) {
	nodeID, err := s.redis.Get(fmt.Sprintf("vk_id:%d", vkID)).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("error getting redis key: %w", err)
	}

	userID, _ := strconv.Atoi(nodeID)
	return userID, nil
}

func (s *Store) SaveUserVKID(vkID int, userID int) error {
	wasSet, err := s.redis.SetNX(fmt.Sprintf("vk_id:%d", vkID), userID, 0).Result()
	if err != nil {
		return fmt.Errorf("error setting vk_id key: %w", err)
	} else if !wasSet {
		return fmt.Errorf("key vk_id:%d already set", vkID)
	}

	return nil
}
