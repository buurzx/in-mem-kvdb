package inmemory

import (
	"context"

	"go.uber.org/zap"
)

type Engine struct {
	hashTable *HashTable
	logger    *zap.Logger
}

func NewEngine(logger *zap.Logger) *Engine {
	return &Engine{
		hashTable: NewHashTable(),
		logger:    logger,
	}
}

func (e *Engine) Get(ctx context.Context, key string) (string, bool) {
	return e.hashTable.Get(key)
}

func (e *Engine) Set(ctx context.Context, key, value string) {
	e.hashTable.Set(key, value)
}

func (e *Engine) Del(ctx context.Context, key string) {
	e.hashTable.Del(key)
}
