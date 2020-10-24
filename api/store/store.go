package store

import (
	"fmt"
	"strconv"
)

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

func (s *Store) AddToFeed(postID int) error {
	_, err := s.redis.LPush("feed", postID).Result()
	if err != nil {
		return fmt.Errorf("error adding to redis list: %w", err)
	}
	return nil
}

func (s *Store) GetFeed() ([]int, error) {
	nodeIdsStr, err := s.redis.LRange("feed", 0, 50).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting feed list: %w", err)
	}

	nodeIds := make([]int, len(nodeIdsStr))
	for i, nodeID := range nodeIdsStr {
		nodeIds[i], _ = strconv.Atoi(nodeID)
	}

	return nodeIds, nil
}
