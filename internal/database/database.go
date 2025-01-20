package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/buurzx/in-mem-kvdb/internal/database/compute"
	"go.uber.org/zap"
)

type Compute interface {
	Parse(request string) (compute.Query, error)
}

type Storage interface {
	Set(context.Context, string, string)
	Del(context.Context, string)
	Get(context.Context, string) (string, error)
}

type Database struct {
	compute Compute
	storage Storage
	logger  *zap.Logger
}

func New(compute Compute, storage Storage, logger *zap.Logger) (*Database, error) {
	if compute == nil {
		return nil, errors.New("invalid compute")
	}

	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	return &Database{
		compute: compute,
		storage: storage,
		logger:  logger,
	}, nil
}

func (d *Database) HandleRequest(ctx context.Context, request string) string {
	d.logger.Debug("received request", zap.String("request", request))

	query, err := d.compute.Parse(request)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	switch query.CommandID() {
	case compute.SetCommandID:
		return d.handleSetRequest(ctx, query)
	case compute.GetCommandID:
		return d.handleGetRequest(ctx, query)
	case compute.DelCommandID:
		return d.handleDelRequest(ctx, query)
	}

	d.logger.Error("unknown command", zap.String("request", request))

	return "[error] internal error"
}

func (d *Database) handleSetRequest(ctx context.Context, query compute.Query) string {
	key := query.Arguments()[0]
	value := query.Arguments()[1]

	d.storage.Set(ctx, key, value)

	return "[OK]"
}

func (d *Database) handleGetRequest(ctx context.Context, query compute.Query) string {
	key := query.Arguments()[0]

	value, err := d.storage.Get(ctx, key)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return value
}

func (d *Database) handleDelRequest(ctx context.Context, query compute.Query) string {
	key := query.Arguments()[0]

	d.storage.Del(ctx, key)

	return "[OK]"
}
