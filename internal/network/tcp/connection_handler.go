package network

import (
	"context"
	"net"
	"time"

	"github.com/buurzx/in-mem-kvdb/internal/database"
	"go.uber.org/zap"
)

type ConnectionHandler struct {
	conn        net.Conn
	db          *database.Database
	logger      *zap.Logger
	idleTimeout time.Duration
	bufferSize  int
}

func NewConnectionHandler(
	conn net.Conn,
	db *database.Database,
	logger *zap.Logger,
	idleTimeout time.Duration,
	bufferSize int,
) *ConnectionHandler {
	return &ConnectionHandler{
		conn:        conn,
		db:          db,
		logger:      logger,
		idleTimeout: idleTimeout,
		bufferSize:  bufferSize,
	}
}

func (h *ConnectionHandler) Handle(ctx context.Context, request []byte) string {
	return h.db.HandleRequest(ctx, string(request))
}
