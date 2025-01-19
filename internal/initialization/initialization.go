package initialization

import (
	"fmt"

	"github.com/buurzx/in-mem-kvdb/internal/config"

	"go.uber.org/zap"
)

type Initializer struct {
	Logger *zap.Logger
}

func NewInitializer(config *config.Config) (*Initializer, error) {
	logger, err := CreateLogger(config.Logger.Level, config.Logger.OutputFilePath)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	return &Initializer{
		Logger: logger,
	}, nil
}

func (i *Initializer) Start() {
	i.Logger.Info("Application started")
}
