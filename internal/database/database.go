package database

import (
	"fmt"

	"github.com/alukart32/go-fast-key/internal/database/compute"
	"go.uber.org/zap"
)

//go:generate mockery --all --exported

// RequestParser describes the database query parser.
type RequestParser interface {
	Parse(req string) (compute.Query, error)
}

// Engine describes the database engine.
type Engine interface {
	Set(k, v string) error
	Get(k string) (string, error)
	Del(k string) error
}

// Database defines the key-value database.
type Database struct {
	parser RequestParser
	e      Engine

	l *zap.Logger
}

// NewDatabase creates a new Database.
func NewDatabase(parser RequestParser, engine Engine, logger *zap.Logger) (*Database, error) {
	if parser == nil {
		return nil, fmt.Errorf("parser is nil")
	}
	if engine == nil {
		return nil, fmt.Errorf("engine is nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	return &Database{
		parser: parser,
		e:      engine,
		l:      logger,
	}, nil
}

// HandleRequest processes the incoming request and returns the query result.
//
// Errors occur due to an incorrect query or inconsistent data.
func (db *Database) HandleRequest(request string) string {
	db.l.Debug("handle the request", zap.String("request", request))

	query, err := db.parser.Parse(request)
	if err != nil {
		return err.Error()
	}

	var result string
	switch query.CommandID() {
	case compute.SetCommand:
		err = db.doSet(query)
	case compute.GetCommand:
		result, err = db.doGet(query)
	case compute.DelCommand:
		err = db.doDel(query)
	}

	if err != nil {
		result = err.Error()
	}
	if len(result) == 0 {
		result = "ok"
	}

	return result
}

func (db *Database) doSet(q compute.Query) error {
	args := q.Arguments()
	return db.e.Set(args[0], args[1])
}

func (db *Database) doGet(q compute.Query) (string, error) {
	args := q.Arguments()
	val, err := db.e.Get(args[0])
	return val, err
}

func (db *Database) doDel(q compute.Query) error {
	args := q.Arguments()
	return db.e.Del(args[0])
}
