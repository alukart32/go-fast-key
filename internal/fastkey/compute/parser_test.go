package compute_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alukart32/go-fast-key/internal/fastkey/compute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		req     string
		parser  func() *compute.Parser
		want    compute.Query
		wantErr error
	}{
		{
			name:    "Empty request",
			req:     "",
			wantErr: compute.ErrEmptyRequest,
		},
		{
			name:    "UTF command symbols",
			req:     "字文下",
			wantErr: compute.ErrUnknownCommand,
		},
		{
			name:    "Unknown command",
			req:     "TRUNCATE",
			wantErr: compute.ErrUnknownCommand,
		},
		{
			name:    "SET command invalid args number",
			req:     "SET key",
			wantErr: compute.ErrInvalidArgsNumber,
		},
		{
			name:    "GET command invalid args number",
			req:     "GET key1 key2 key3",
			wantErr: compute.ErrInvalidArgsNumber,
		},
		{
			name:    "DEL command invalid args number",
			req:     "DEL key1 key2 key3",
			wantErr: compute.ErrInvalidArgsNumber,
		},
		{
			name: "Valid SET request",
			req:  "SET key val",
			want: compute.NewQuery(compute.SetCommand, []string{"key", "val"}),
		},
		{
			name: "Valid GET request",
			req:  "GET key",
			want: compute.NewQuery(compute.GetCommand, []string{"key"}),
		},
		{
			name: "Valid DEL request",
			req:  "DEL key",
			want: compute.NewQuery(compute.DelCommand, []string{"key"}),
		},
	}

	parser, err := compute.NewParser(zap.NewNop())
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Parse(tt.req)
			assert.Equal(t, err, tt.wantErr, "Parse() error = %v, wantErr %v", err, tt.wantErr)
			assert.True(t, reflect.DeepEqual(got, tt.want), "Parse() = %v, want %v", got, tt.want)
		})
	}
}

func TestNewParser(t *testing.T) {
	tests := []struct {
		name          string
		logger        *zap.Logger
		wantErr       error
		wantNilObject bool
	}{
		{
			name:          "Logger is nil",
			logger:        nil,
			wantErr:       fmt.Errorf("logger is nil"),
			wantNilObject: true,
		},
		{
			name:          "Valid parser",
			logger:        zap.NewNop(),
			wantErr:       nil,
			wantNilObject: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compute.NewParser(tt.logger)
			assert.Equal(t, err, tt.wantErr, "NewParser() error = %v, wantErr %v", err, tt.wantErr)
			if tt.wantNilObject {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
			}
		})
	}
}
