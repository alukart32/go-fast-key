package application_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alukart32/go-fast-key/internal/application"
	"github.com/alukart32/go-fast-key/internal/configuration"
	"github.com/alukart32/go-fast-key/internal/pkg/datasize"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateNetwork(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg    *configuration.Network
		logger *zap.Logger

		wantErr    error
		wantNilObj bool
	}{
		"create network without logger": {
			wantErr:    errors.New("logger is nil"),
			wantNilObj: true,
		},
		"create network without config": {
			logger:  zap.NewNop(),
			wantErr: nil,
		},
		"create network with empty config fields": {
			logger: zap.NewNop(),
			cfg: &configuration.Network{
				Address: "localhost:20002",
			},
			wantErr: nil,
		},
		"create network with config fields": {
			logger: zap.NewNop(),
			cfg: &configuration.Network{
				Address:        "localhost:10001",
				MaxConnections: 100,
				MaxMessageSize: "2KB",
				IdleTimeout:    time.Second,
			},
			wantErr: nil,
		},
		"create network with incorrect size": {
			logger: zap.NewNop(),
			cfg: &configuration.Network{
				MaxMessageSize: "2incorrect",
			},
			wantErr:    errors.New("parse message size: " + datasize.ErrInvalidSize.Error()),
			wantNilObj: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			network, err := application.CreateNetwork(test.cfg, test.logger)
			assert.Equal(t, test.wantErr, err)
			if test.wantNilObj {
				assert.Nil(t, network)
			} else {
				assert.NotNil(t, network)
			}
		})
	}
}
