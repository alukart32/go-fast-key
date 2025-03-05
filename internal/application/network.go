package application

import (
	"errors"
	"fmt"

	"github.com/alukart32/go-fast-key/internal/configuration"
	"github.com/alukart32/go-fast-key/internal/network"
	"github.com/alukart32/go-fast-key/internal/pkg/datasize"
	"go.uber.org/zap"
)

const defaultServerAddress = ":3223"

func CreateNetwork(cfg *configuration.Network, logger *zap.Logger) (*network.TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	address := defaultServerAddress
	var options []network.TCPServerOption

	if cfg != nil {
		if cfg.Address != "" {
			address = cfg.Address
		}

		if cfg.MaxConnections != 0 {
			options = append(options, network.WithServerMaxConnectionsNumber(uint(cfg.MaxConnections)))
		}

		if cfg.MaxMessageSize != "" {
			size, err := datasize.Parse(cfg.MaxMessageSize)
			if err != nil {
				return nil, fmt.Errorf("parse message size: %v", err)
			}

			options = append(options, network.WithServerBufferSize(uint(size)))
		}

		if cfg.IdleTimeout != 0 {
			options = append(options, network.WithServerIdleTimeout(cfg.IdleTimeout))
		}
	}

	return network.NewTCPServer(address, logger, options...)
}
