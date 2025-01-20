package inmemory

import "sync"

type HashTable struct {
	mutex sync.RWMutex
	data  map[string]string
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]string),
	}
}

func (h *HashTable) Set(key, value string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.data[key] = value
}

func (h *HashTable) Get(key string) (string, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	value, ok := h.data[key]
	return value, ok
}

func (h *HashTable) Del(key string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.data, key)
}
