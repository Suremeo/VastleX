package blocks

import (
	"encoding/json"
	"hash/fnv"
	"sync"
)

type Store struct {
	blocks        sync.Map
	blocksReverse sync.Map
	I             int64
}

func (store *Store) Import(blocks []interface{}) {
	for _, block := range blocks {
		dat, _ := json.Marshal(block)
		h := hash(dat)
		store.blocks.Store(h, store.I)
		store.blocksReverse.Store(store.I, h)
		//println("Registered: ", store.I)
		store.I++
	}
}

func (store *Store) RuntimeFromHash(hash int64) int64 {
	r, ok := store.blocks.Load(hash)
	if ok {
		return r.(int64)
	} else {
		return 0
	}
}

func (store *Store) HashFromRuntime(hash int64) int64 {
	r, ok := store.blocksReverse.Load(hash)
	if ok {
		return r.(int64)
	} else {
		return 0
	}
}

func hash(s []byte) int64 {
	h := fnv.New64a()
	_, _ = h.Write(s)
	return int64(h.Sum64())
}
