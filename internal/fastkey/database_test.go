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
			name:    "Database is nil",
			request: "",
			parser:  nil,
			storage: nil,
			want:    "",
			wantErr: fastkey.ErrStandBy,
		},
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
			storage: func() fastkey.Storage { return nil },
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
			var db *fastkey.Database
			if tt.parser != nil || tt.storage != nil {
				db = fastkey.NewDatabase(tt.parser(), tt.storage())
			}

			got, err := db.HandleRequest(tt.request)
			if err != nil && tt.wantErr != nil {
				assert.EqualErrorf(t, err, tt.wantErr.Error(),
					"Engine.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err, "Database.HandleRequest() error = %v", err)

			if got != tt.want {
				t.Errorf("Database.HandleRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
