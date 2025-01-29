package fastkey_test

import (
	"fmt"
	"testing"

	"github.com/alukart32/go-fast-key/internal/fastkey"
	"github.com/alukart32/go-fast-key/internal/fastkey/compute"
	fastkey_mocks "github.com/alukart32/go-fast-key/internal/fastkey/mocks"
	"github.com/alukart32/go-fast-key/internal/fastkey/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDatabase_HandleRequest(t *testing.T) {
	tests := []struct {
		name    string
		request string
		parser  func() fastkey.RequestParser
		storage func() fastkey.Storage
		want    string
		wantErr error
	}{
		{
			name:    "Handle command with parser error",
			request: "TRUNCATE key",
			parser: func() fastkey.RequestParser {
				m := fastkey_mocks.NewRequestParser(t)
				m.
					On("Parse", "TRUNCATE key").
					Return(compute.Query{}, fmt.Errorf("parser error")).
					Once()
				return m
			},
			storage: func() fastkey.Storage { return fastkey_mocks.NewStorage(t) },
			want:    "",
			wantErr: fmt.Errorf("parser error"),
		},
		{
			name:    "Valid SET query",
			request: "SET key val",
			parser: func() fastkey.RequestParser {
				m := fastkey_mocks.NewRequestParser(t)
				query := compute.NewQuery(compute.SetCommand, []string{"key", "val"})
				m.On("Parse", "SET key val").Return(query, nil).Once()
				return m
			},
			storage: func() fastkey.Storage {
				m := fastkey_mocks.NewStorage(t)
				m.On("Set", "key", "val").Return(nil).Once()
				return m
			},
			want:    "",
			wantErr: nil,
		},
		{
			name:    "SET query with storage error",
			request: "SET key val",
			parser: func() fastkey.RequestParser {
				m := fastkey_mocks.NewRequestParser(t)
				query := compute.NewQuery(compute.SetCommand, []string{"key", "val"})
				m.On("Parse", "SET key val").Return(query, nil).Once()
				return m
			},
			storage: func() fastkey.Storage {
				m := fastkey_mocks.NewStorage(t)
				m.On("Set", "key", "val").Return(fmt.Errorf("storage error")).Once()
				return m
			},
			want:    "",
			wantErr: fmt.Errorf("storage error"),
		},
		{
			name:    "Valid GET query",
			request: "GET key",
			parser: func() fastkey.RequestParser {
				m := fastkey_mocks.NewRequestParser(t)

				query := compute.NewQuery(compute.GetCommand, []string{"key"})
				m.On("Parse", "GET key").Return(query, nil).Once()
				return m
			},
			storage: func() fastkey.Storage {
				m := fastkey_mocks.NewStorage(t)
				m.On("Get", "key").Return("val", nil).Once()
				return m
			},
			want:    "val",
			wantErr: nil,
		},
		{
			name:    "GET query with storage error",
			request: "GET key",
			parser: func() fastkey.RequestParser {
				m := fastkey_mocks.NewRequestParser(t)
				query := compute.NewQuery(compute.GetCommand, []string{"key"})
				m.On("Parse", "GET key").Return(query, nil).Once()
				return m
			},
			storage: func() fastkey.Storage {
				m := fastkey_mocks.NewStorage(t)
				m.On("Get", "key").Return("", fmt.Errorf("storage error")).Once()
				return m
			},
			want:    "",
			wantErr: fmt.Errorf("storage error"),
		},
		{
			name:    "GET query with not found error",
			request: "GET key",
			parser: func() fastkey.RequestParser {
				m := fastkey_mocks.NewRequestParser(t)
				query := compute.NewQuery(compute.GetCommand, []string{"key"})
				m.On("Parse", "GET key").Return(query, nil).Once()
				return m
			},
			storage: func() fastkey.Storage {
				m := fastkey_mocks.NewStorage(t)
				m.On("Get", "key").Return("", storage.ErrNotFound).Once()
				return m
			},
			want:    "",
			wantErr: storage.ErrNotFound,
		},
		{
			name:    "Valid DEL query",
			request: "DEL key",
			parser: func() fastkey.RequestParser {
				m := fastkey_mocks.NewRequestParser(t)

				query := compute.NewQuery(compute.DelCommand, []string{"key"})
				m.On("Parse", "DEL key").Return(query, nil).Once()
				return m
			},
			storage: func() fastkey.Storage {
				m := fastkey_mocks.NewStorage(t)
				m.On("Del", "key").Return(nil).Once()
				return m
			},
			want:    "",
			wantErr: nil,
		},
		{
			name:    "DEL query with storage error",
			request: "DEL key",
			parser: func() fastkey.RequestParser {
				m := fastkey_mocks.NewRequestParser(t)
				query := compute.NewQuery(compute.DelCommand, []string{"key"})
				m.On("Parse", "DEL key").Return(query, nil).Once()
				return m
			},
			storage: func() fastkey.Storage {
				m := fastkey_mocks.NewStorage(t)
				m.On("Del", "key").Return(fmt.Errorf("storage error")).Once()
				return m
			},
			want:    "",
			wantErr: fmt.Errorf("storage error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := fastkey.NewDatabase(tt.parser(), tt.storage(), zap.NewNop())
			require.NoError(t, err)

			got, err := db.HandleRequest(tt.request)
			assert.Equal(t, err, tt.wantErr, "HandleRequest() error = %v, wantErr %v", err, tt.wantErr)
			assert.True(t, got == tt.want, "HandleRequest() = %v, want %v", got, tt.want)
		})
	}
}

func TestNewDatabase(t *testing.T) {
	tests := []struct {
		name          string
		parser        fastkey.RequestParser
		db            fastkey.Storage
		logger        *zap.Logger
		wantErr       error
		wantNilObject bool
	}{
		{
			name:          "Create without parser",
			parser:        nil,
			db:            fastkey_mocks.NewStorage(t),
			logger:        zap.NewNop(),
			wantErr:       fmt.Errorf("parser is nil"),
			wantNilObject: true,
		},
		{
			name:          "Create without storage",
			parser:        fastkey_mocks.NewRequestParser(t),
			db:            nil,
			logger:        zap.NewNop(),
			wantErr:       fmt.Errorf("storage is nil"),
			wantNilObject: true,
		},
		{
			name:          "Create without logger",
			parser:        fastkey_mocks.NewRequestParser(t),
			db:            fastkey_mocks.NewStorage(t),
			logger:        nil,
			wantErr:       fmt.Errorf("logger is nil"),
			wantNilObject: true,
		},
		{
			name:          "Created",
			parser:        fastkey_mocks.NewRequestParser(t),
			db:            fastkey_mocks.NewStorage(t),
			logger:        zap.NewNop(),
			wantErr:       nil,
			wantNilObject: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fastkey.NewDatabase(tt.parser, tt.db, tt.logger)

			assert.Equal(t, err, tt.wantErr, "NewDatabase() error = %v, wantErr %v", err, tt.wantErr)
			if tt.wantNilObject {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
			}
		})
	}
}
