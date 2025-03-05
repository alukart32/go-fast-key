package configuration

import (
	"errors"
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Engine  *Engine  `yaml:"engine"`
	Network *Network `yaml:"network"`
	Logging *Logging `yaml:"logging"`
}

type Engine struct {
	Type string `yaml:"type"`
}

type Network struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type Logging struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

func Load(reader io.Reader) (*Config, error) {
	if reader == nil {
		return nil, errors.New("incorrect reader")
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("fail to read buffer")
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
