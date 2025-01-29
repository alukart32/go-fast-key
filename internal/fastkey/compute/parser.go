package compute

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// Parser defines the query parser for execution.
type Parser struct {
	l *zap.Logger
}

// NewParser creates a new Parser.
func NewParser(logger *zap.Logger) (*Parser, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}
	return &Parser{l: logger}, nil
}

// Parse converts the request into a query.
func (p *Parser) Parse(req string) (Query, error) {
	tokens := strings.Fields(strings.TrimSpace(req))
	if len(tokens) == 0 {
		p.l.Debug("empty tokens", zap.String("request", req))
		return Query{}, ErrEmptyRequest
	}

	command := tokens[0]
	commandID, err := commandNameToCommandID(command)
	if err != nil {
		p.l.Debug("invalid command", zap.String("request", req))
		return Query{}, err
	}

	query := NewQuery(commandID, tokens[1:])
	argsNumber := commandIDToArgsNumber(query.commandID)
	if len(query.arguments) != argsNumber {
		p.l.Debug("invalid arguments for query", zap.String("request", req))
		return Query{}, ErrInvalidArgsNumber
	}
	return query, nil
}
