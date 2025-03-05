package application

import (
	"errors"
	"fmt"

	"github.com/alukart32/go-fast-key/internal/configuration"
	"github.com/alukart32/go-fast-key/internal/database"
	"github.com/alukart32/go-fast-key/internal/database/engine"
	"go.uber.org/zap"
)

func CreateEngine(cfg *configuration.Engine, logger *zap.Logger) (database.Engine, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	if cfg == nil {
		return engine.NewMemEngine(256), nil
	}

	if cfg.Type != "" {
		supportedTypes := map[string]struct{}{
			"in_memory": {},
		}

		if _, found := supportedTypes[cfg.Type]; !found {
			return nil, fmt.Errorf("unsupported engine type: %v", cfg.Type)
		}
	}

	return engine.NewMemEngine(256), nil
}
