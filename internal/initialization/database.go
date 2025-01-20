package initialization

import (
	"fmt"

	"github.com/buurzx/in-mem-kvdb/internal/database"
	"github.com/buurzx/in-mem-kvdb/internal/database/compute"
	"github.com/buurzx/in-mem-kvdb/internal/database/storage"
	inmemory "github.com/buurzx/in-mem-kvdb/internal/database/storage/engine/in_memory"
	"go.uber.org/zap"
)

func CreateDatabase(logger *zap.Logger) (*database.Database, error) {
	compute, err := compute.New(logger)
	if err != nil {
		return nil, fmt.Errorf("initialize compute: %w", err)
	}

	engine := inmemory.NewEngine(logger)

	storage, err := storage.New(logger, engine)
	if err != nil {
		return nil, fmt.Errorf("initialize storage: %w", err)
	}

	db, err := database.New(compute, storage, logger)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	return db, nil
}
