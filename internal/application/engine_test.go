package application_test

import (
	"errors"
	"testing"

	"github.com/alukart32/go-fast-key/internal/application"
	"github.com/alukart32/go-fast-key/internal/configuration"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateEngine(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg        *configuration.Engine
		logger     *zap.Logger
		wantErr    error
		wantNilObj bool
	}{
		"create engine without logger": {
			wantErr:    errors.New("logger is nil"),
			wantNilObj: true,
		},
		"create engine without config": {
			logger:  zap.NewNop(),
			wantErr: nil,
		},
		"create engine with empty config fields": {
			cfg:     &configuration.Engine{},
			logger:  zap.NewNop(),
			wantErr: nil,
		},
		"create engine with config fields": {
			cfg: &configuration.Engine{
				Type: "in_memory",
			},
			logger:  zap.NewNop(),
			wantErr: nil,
		},
		"create engine with incorrect type": {
			cfg:        &configuration.Engine{Type: "invalid"},
			logger:     zap.NewNop(),
			wantErr:    errors.New("unsupported engine type: invalid"),
			wantNilObj: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			engine, err := application.CreateEngine(test.cfg, test.logger)
			assert.Equal(t, test.wantErr, err)
			if test.wantNilObj {
				assert.Nil(t, engine)
			} else {
				assert.NotNil(t, engine)
			}
		})
	}
}
