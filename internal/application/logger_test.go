package application_test

import (
	"fmt"
	"testing"

	"github.com/alukart32/go-fast-key/internal/application"
	"github.com/alukart32/go-fast-key/internal/configuration"
	"github.com/stretchr/testify/assert"
)

func TestCreateLogger(t *testing.T) {
	tests := map[string]struct {
		cfg           *configuration.Logging
		wantErr       error
		wantNilObject bool
	}{
		"create logger without config": {
			wantErr: nil,
		},
		"create logger with empty config fields": {
			cfg:     &configuration.Logging{},
			wantErr: nil,
		},
		"create logger with config fields": {
			cfg: &configuration.Logging{
				Level:  application.DebugLogLevel,
				Output: "test.log",
			},
			wantErr: nil,
		},
		"create logger with invalid level": {
			cfg:           &configuration.Logging{Level: "invalid"},
			wantErr:       fmt.Errorf("unsupported level: invalid"),
			wantNilObject: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := application.CreateLogger(test.cfg)
			assert.Equal(t, test.wantErr, err)
			if test.wantNilObject {
				assert.Nil(t, logger)
			} else {
				assert.NotNil(t, logger)
			}
		})
	}
}
