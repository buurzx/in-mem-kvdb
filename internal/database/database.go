package database

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

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

type CommandHandler func(context.Context, []string) string

type Database struct {
	compute  Compute
	storage  Storage
	logger   *zap.Logger
	commands map[string]CommandHandler
}

func New(compute Compute, storage Storage, logger *zap.Logger) (*Database, error) {
	if compute == nil {
		return nil, errors.New("invalid compute")
	}

	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	db := &Database{
		compute:  compute,
		storage:  storage,
		logger:   logger,
		commands: make(map[string]CommandHandler),
	}

	// Register commands
	db.commands["GET"] = db.handleGetRequest
	db.commands["SET"] = db.handleSetRequest
	db.commands["DEL"] = db.handleDelRequest

	return db, nil
}

func (d *Database) HandleRequest(ctx context.Context, request string) string {
	parts := strings.Fields(request)
	if len(parts) == 0 {
		return "Empty command"
	}

	command := strings.ToUpper(parts[0])
	handler, exists := d.commands[command]
	if !exists {
		validCommands := make([]string, 0, len(d.commands))
		for cmd := range d.commands {
			validCommands = append(validCommands, cmd)
		}
		sort.Strings(validCommands) // Sort for consistent output
		return fmt.Sprintf("Invalid command. Available commands: %s", strings.Join(validCommands, ", "))
	}

	return handler(ctx, parts[1:])
}

func (d *Database) handleSetRequest(ctx context.Context, query []string) string {
	if len(query) < 2 {
		return "Invalid SET command. Usage: SET <key> <value>"
	}

	key := query[0]
	value := query[1]

	d.storage.Set(ctx, key, value)

	return "[OK]"
}

func (d *Database) handleGetRequest(ctx context.Context, query []string) string {
	if len(query) < 1 {
		return "Invalid GET command. Usage: GET <key>"
	}

	key := query[0]

	value, err := d.storage.Get(ctx, key)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return value
}

func (d *Database) handleDelRequest(ctx context.Context, query []string) string {
	if len(query) < 1 {
		return "Invalid DEL command. Usage: DEL <key>"
	}

	key := query[0]

	d.storage.Del(ctx, key)

	return "[OK]"
}
