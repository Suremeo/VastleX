package blocks

import (
	"encoding/json"
	"hash/fnv"
	"sync"
)

// Store is a block store, used for translating the block state ids.
type Store struct {
	blocks        sync.Map
	blocksReverse sync.Map
	I             int64
}

// Import imports an array of blocks into the block store.
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

// RuntimeFromHash converts a block hash to a block runtime id.
func (store *Store) RuntimeFromHash(hash int64) int64 {
	r, ok := store.blocks.Load(hash)
	if ok {
		return r.(int64)
	} else {
		return 0
	}
}

// HashFromRuntime converts a runtime id to a block hash.
func (store *Store) HashFromRuntime(runtime int64) int64 {
	r, ok := store.blocksReverse.Load(runtime)
	if ok {
		return r.(int64)
	} else {
		return 0
	}
}

// hash returns a fnv64a of the supplied bytes.
func hash(bytes []byte) int64 {
	h := fnv.New64a()
	_, _ = h.Write(bytes)
	return int64(h.Sum64())
}
