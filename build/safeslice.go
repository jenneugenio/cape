package build

import (
	"sync"
)

type SafeSlice struct {
	values []interface{}
	lock   *sync.Mutex
}

func NewSafeSlice() *SafeSlice {
	return &SafeSlice{
		values: []interface{}{},
		lock:   &sync.Mutex{},
	}
}

func (s *SafeSlice) Add(item interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.values = append(s.values, item)
}

func (s *SafeSlice) Get() []interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.values
}
