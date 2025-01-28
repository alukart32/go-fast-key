package compute

import (
	"strings"
)

// Parser defines the query parser for execution.
type Parser struct{}

// NewParser creates a new Parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse converts the request into a query.
func (p *Parser) Parse(req string) (Query, error) {
	if p == nil {
		return Query{}, ErrStandByParser
	}

	tokens := strings.Fields(strings.TrimSpace(req))
	if len(tokens) == 0 {
		return Query{}, ErrEmptyRequest
	}

	command := tokens[0]
	commandID, err := commandNameToCommandID(command)
	if err != nil {
		return Query{}, err
	}

	query := NewQuery(commandID, tokens[1:])
	argsNumber := commandIDToArgsNumber(query.commandID)
	if len(query.arguments) != argsNumber {
		return Query{}, ErrInvalidArgsNumber
	}
	return query, nil
}
