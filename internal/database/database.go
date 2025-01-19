package database

import (
	"errors"
	"fmt"

	"github.com/buurzx/in-mem-kvdb/internal/database/compute"
	"go.uber.org/zap"
)

type Compute interface {
	Parse(request string) (compute.Query, error)
}

type Database struct {
	compute Compute
	logger  *zap.Logger
}

func New(compute Compute, logger *zap.Logger) (*Database, error) {
	if compute == nil {
		return nil, errors.New("invalid compute")
	}

	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	return &Database{
		compute: compute,
		logger:  logger,
	}, nil
}

func (d *Database) HandleRequest(request string) error {
	d.logger.Debug("received request", zap.String("request", request))

	query, err := d.compute.Parse(request)
	if err != nil {
		return fmt.Errorf("parse request: %w", err)
	}

	switch query.CommandID() {
	case compute.SetCommandID:
		d.logger.Debug("set command")
	case compute.GetCommandID:
		d.logger.Debug("get command")
	case compute.DelCommandID:
		d.logger.Debug("del command")
	default:
		d.logger.Debug("unknown command")
	}

	return nil
}
