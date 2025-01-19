package initialization

import (
	"fmt"

	"github.com/buurzx/in-mem-kvdb/internal/database"
	"github.com/buurzx/in-mem-kvdb/internal/database/compute"
	"go.uber.org/zap"
)

func CreateDatabase(logger *zap.Logger) (*database.Database, error) {
	compute, err := compute.New(logger)
	if err != nil {
		return nil, fmt.Errorf("initialize compute: %w", err)
	}

	db, err := database.New(compute, logger)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	return db, nil
}
