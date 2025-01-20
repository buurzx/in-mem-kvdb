package storage

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	errNotFound = errors.New("not found")
)

type Engine interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key, value string)
	Del(ctx context.Context, key string)
}

type Storage struct {
	engine Engine
	logger *zap.Logger
}

func New(logger *zap.Logger, engine Engine) (*Storage, error) {
	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	return &Storage{
		engine: engine,
		logger: logger,
	}, nil
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	value, found := s.engine.Get(ctx, key)
	if found {
		return value, nil
	}

	return "", errNotFound
}

func (s *Storage) Set(ctx context.Context, key, value string) {
	s.engine.Set(ctx, key, value)
}

func (s *Storage) Del(ctx context.Context, key string) {
	s.engine.Del(ctx, key)
}
