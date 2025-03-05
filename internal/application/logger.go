package application

import (
	"fmt"

	"github.com/alukart32/go-fast-key/internal/configuration"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLogLevel = "debug"
	InfoLogLevel  = "info"
	WarnLogLevel  = "warn"
	ErrorLogLevel = "error"
)

const (
	defaultLogEncoding   = "json"
	defaultLogLevel      = zapcore.DebugLevel
	defaultLogOutputPath = "fastkey.log"
)

func CreateLogger(cfg *configuration.Logging) (*zap.Logger, error) {
	level := defaultLogLevel
	output := defaultLogOutputPath

	if cfg != nil {
		if cfg.Level != "" {
			supportedLoggingLevels := map[string]zapcore.Level{
				DebugLogLevel: zapcore.DebugLevel,
				InfoLogLevel:  zapcore.InfoLevel,
				WarnLogLevel:  zapcore.WarnLevel,
				ErrorLogLevel: zapcore.ErrorLevel,
			}

			var found bool
			if level, found = supportedLoggingLevels[cfg.Level]; !found {
				return nil, fmt.Errorf("unsupported level: %v", cfg.Level)
			}
		}

		if cfg.Output != "" {
			output = cfg.Output
		}
	}

	loggerCfg := zap.Config{
		Encoding:    defaultLogEncoding,
		Development: true,
		Level:       zap.NewAtomicLevelAt(level),
		OutputPaths: []string{output},
	}

	return loggerCfg.Build()
}
