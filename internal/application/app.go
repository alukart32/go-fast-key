package application

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/alukart32/go-fast-key/internal/configuration"
	"github.com/alukart32/go-fast-key/internal/database"
	"github.com/alukart32/go-fast-key/internal/database/compute"
	"github.com/alukart32/go-fast-key/internal/network"
	"go.uber.org/zap"
)

type App struct {
	dbEngine database.Engine
	server   *network.TCPServer
	logger   *zap.Logger
}

func NewApp(cfg *configuration.Config) (*App, error) {
	if cfg == nil {
		return nil, errors.New("new application: config is invalid")
	}

	logger, _ := zap.NewDevelopment()
	// logger, err := CreateLogger(cfg.Logging)
	// if err != nil {
	// 	return nil, fmt.Errorf("create logger: %w", err)
	// }

	engine, err := CreateEngine(cfg.Engine, logger)
	if err != nil {
		return nil, fmt.Errorf("create database engine: %w", err)
	}

	server, err := CreateNetwork(cfg.Network, logger)
	if err != nil {
		return nil, fmt.Errorf("create network: %w", err)
	}

	app := App{
		dbEngine: engine,
		server:   server,
		logger:   logger,
	}

	return &app, nil
}

func (a *App) Run(ctx context.Context) error {
	requestParser, err := compute.NewParser(a.logger)
	if err != nil {
		return fmt.Errorf("create the request parser: %v", err)
	}

	db, err := database.NewDatabase(requestParser, a.dbEngine, a.logger)
	if err != nil {
		return fmt.Errorf("create the database: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		a.server.HandleQueries(ctx, func(_ context.Context, request []byte) []byte {
			response, err := db.HandleRequest(string(request))
			if err != nil {
				return []byte(err.Error())
			}
			return []byte(response)
		})
	}()

	wg.Wait()

	a.logger.Sync()
	return err
}
