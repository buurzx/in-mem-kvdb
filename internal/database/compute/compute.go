package compute

import (
	"errors"
	"strings"

	"go.uber.org/zap"
)

var (
	errInvalidLogger    = errors.New("invalid logger")
	errEmptyRequest     = errors.New("empty request")
	errInvalidCommand   = errors.New("invalid command")
	errInvalidArguments = errors.New("invalid arguments")
)

type Compute struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) (*Compute, error) {
	if logger == nil {
		return nil, errInvalidLogger
	}

	return &Compute{
		logger: logger,
	}, nil
}

func (c *Compute) Parse(request string) (Query, error) {
	tokens := strings.Fields(request)
	if len(tokens) == 0 {
		c.logger.Error("empty request")
		return Query{}, errEmptyRequest
	}

	command := tokens[0]

	commandID := commandNameToCommandID(command)
	if commandID == UnknownCommandID {
		c.logger.Debug("invalid command", zap.String("query", request))
		return Query{}, errInvalidCommand
	}

	query := NewQuery(commandID, tokens[1:])

	argumentsNumber := commandArgumentsNumber(commandID)
	if len(query.Arguments()) != argumentsNumber {
		c.logger.Debug("invalid arguments for query", zap.String("query", request))
		return Query{}, errInvalidArguments
	}

	return query, nil
}
