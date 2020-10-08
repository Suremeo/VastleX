package entity

import (
	"sync"
)

type Store struct {
	entities sync.Map
}

func (store *Store) Get(eid int64) int64 {
	r, ok := store.entities.Load(eid)
	if !ok {
		return 0
	} else {
		return r.(int64)
	}
}
func (store *Store) Set(one, two int64) *Store {
	store.entities.Store(one, two)
	return store
}

func (store *Store) Delete(eid int64) *Store {
	store.entities.Delete(eid)
	return store
}

func (store *Store) Range(fun func(key, value interface{}) bool) {
	store.entities.Range(fun)
}

func (store *Store) Clear() *Store {
	store.entities = sync.Map{}
	return store
}
