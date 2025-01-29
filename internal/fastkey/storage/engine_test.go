package storage_test

import (
	"testing"

	"github.com/alukart32/go-fast-key/internal/fastkey/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngine_Set(t *testing.T) {
	type pair struct {
		k string
		v string
	}
	tests := []struct {
		name        string
		eng         *storage.Engine
		args        pair
		testData    pair
		setTestData bool
		wantErr     error
	}{
		{
			name: "Empty key",
			eng:  storage.NewEngine(0),
			args: pair{
				k: "",
				v: "val",
			},
			wantErr: storage.ErrInvalidEntityID,
		},
		{
			name: "Empty val",
			eng:  storage.NewEngine(0),
			args: pair{
				k: "key",
				v: "",
			},
			wantErr: storage.ErrInvalidEntityData,
		},
		{
			name: "Reset val",
			eng:  storage.NewEngine(0),
			args: pair{
				k: "key",
				v: "val_2",
			},
			setTestData: true,
			testData: pair{
				k: "key",
				v: "val_1",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setTestData {
				require.NoError(t, tt.eng.Set(tt.testData.k, tt.testData.v))
			}

			err := tt.eng.Set(tt.args.k, tt.args.v)
			if tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr, "Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err, "Set() error = %v, want nil error", err)

			gotVal, err := tt.eng.Get(tt.args.k)
			require.NoError(t, err, "Get() error = %v for %v key", err, tt.args.k)
			assert.EqualValues(t, tt.args.v, gotVal, "Get() got %v by %v key, want %v", gotVal, tt.args.k, tt.args.v)
		})
	}
}

func TestEngine_Get(t *testing.T) {
	type pair struct {
		k string
		v string
	}
	tests := []struct {
		name        string
		key         string
		eng         *storage.Engine
		testData    pair
		setTestData bool
		want        string
		wantErr     error
	}{
		{
			name:    "Empty key",
			key:     "",
			eng:     storage.NewEngine(0),
			want:    "",
			wantErr: storage.ErrInvalidEntityID,
		},
		{
			name:    "Not found by key",
			key:     "key",
			eng:     storage.NewEngine(0),
			want:    "",
			wantErr: storage.ErrNotFound,
		},
		{
			name:        "Found by key",
			key:         "key",
			eng:         storage.NewEngine(0),
			setTestData: true,
			testData: pair{
				k: "key",
				v: "val",
			},
			want:    "val",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setTestData {
				require.NoError(t, tt.eng.Set(tt.testData.k, tt.testData.v))
			}

			got, err := tt.eng.Get(tt.key)
			if tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr, "Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.NoError(t, err, "Get() error = %v for %v key", err, tt.key)
			assert.EqualValues(t, tt.want, got, "Get() got %v by %v key, want %v", got, tt.key, tt.want)
		})
	}
}

func TestEngine_Del(t *testing.T) {
	type pair struct {
		k string
		v string
	}
	tests := []struct {
		name        string
		key         string
		eng         *storage.Engine
		testData    pair
		setTestData bool
		wantErr     error
	}{
		{
			name:    "Empty key",
			key:     "",
			eng:     storage.NewEngine(0),
			wantErr: storage.ErrInvalidEntityID,
		},
		{
			name:        "Deleted by key",
			key:         "key",
			eng:         storage.NewEngine(0),
			setTestData: true,
			testData: pair{
				k: "key",
				v: "val",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setTestData {
				require.NoError(t, tt.eng.Set(tt.testData.k, tt.testData.v))
			}

			err := tt.eng.Del(tt.key)
			if tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr, "Del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.NoError(t, err, "Del() error = %v for %v key", err, tt.key)
			_, err = tt.eng.Get(tt.key)
			require.ErrorIs(t, err, storage.ErrNotFound,
				"Get() error = %v for %v key, want %v", err, tt.key, storage.ErrNotFound)
		})
	}
}
