package nodes

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	cache  map[int]int
	nextID int

	m       sync.Mutex
	waiting map[int]func() int
}

var ErrNodeNotFound = fmt.Errorf("node not found")

func (s *Store) Get(id int, node any) int {
	s.m.Lock()
	defer s.m.Unlock()

	result, ok := s.waiting[id]
	if ok {
		return result()
	}

	result = sync.OnceValue(func() int {
		time.Sleep(time.Second * 3)
		return 900
	})

	s.waiting[id] = result
	return s.waiting[id]()
}
