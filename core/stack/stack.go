package stack

import "sync"

type Stack struct {
	mux   sync.RWMutex
	chats map[int64]*Query
}

func New() *Stack {
	return &Stack{
		chats: make(map[int64]*Query),
	}
}
func (s *Stack) Replace(id int64, city, place string) {
	s.mux.Lock()
	s.chats[id] = &Query{
		City:  city,
		Place: place,
	}
	s.mux.Unlock()
}

func (s *Stack) Get(id int64) (*Query, bool) {
	s.mux.RLock()
	query, ok := s.chats[id]
	s.mux.RUnlock()

	return query, ok
}

func (s *Stack) Delete(id int64) {
	s.mux.Lock()
	delete(s.chats, id)
	s.mux.Unlock()
}

type Query struct {
	City  string
	Place string
}
