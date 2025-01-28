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
		name     string
		eng      *storage.Engine
		args     pair
		testData pair
		wantErr  error
	}{
		{
			name: "Nil engine",
			eng:  nil,
			args: pair{
				k: "key",
				v: "val",
			},
			wantErr: storage.ErrStandByEngine,
		},
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
			testData: pair{
				k: "key",
				v: "val_1",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zeroTestData := pair{}
			if tt.testData != zeroTestData {
				require.NoError(t, tt.eng.Set(tt.testData.k, tt.testData.v))
			}

			err := tt.eng.Set(tt.args.k, tt.args.v)
			if err != nil && tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "Engine.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err, "Engine.Set() error = %v, want nil error", err)

			gotVal, err := tt.eng.Get(tt.args.k)
			require.NoError(t, err, "Engine.Get() error = %v for %v key", err, tt.args.k)
			assert.EqualValues(t, tt.args.v, gotVal, "Engine.Get() got %v by %v key, want %v", gotVal, tt.args.k, tt.args.v)
		})
	}
}

func TestEngine_Get(t *testing.T) {
	type pair struct {
		k string
		v string
	}
	tests := []struct {
		name     string
		key      string
		eng      *storage.Engine
		testData pair
		want     string
		wantErr  error
	}{
		{
			name:    "Nil engine",
			key:     "key",
			eng:     nil,
			wantErr: storage.ErrStandByEngine,
		},
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
			name: "Found by key",
			key:  "key",
			eng:  storage.NewEngine(0),
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
			zeroTestData := pair{}
			if tt.testData != zeroTestData {
				require.NoError(t, tt.eng.Set(tt.testData.k, tt.testData.v))
			}

			got, err := tt.eng.Get(tt.key)
			if err != nil && tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "Engine.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err, "Engine.Get() error = %v for %v key", err, tt.key)
			assert.EqualValues(t, tt.want, got, "Engine.Get() got %v by %v key, want %v", got, tt.key, tt.want)
		})
	}
}

func TestEngine_Del(t *testing.T) {
	type pair struct {
		k string
		v string
	}
	tests := []struct {
		name     string
		key      string
		eng      *storage.Engine
		testData pair
		wantErr  error
	}{
		{
			name:    "Nil engine",
			key:     "key",
			eng:     nil,
			wantErr: storage.ErrStandByEngine,
		},
		{
			name:    "Empty key",
			key:     "",
			eng:     storage.NewEngine(0),
			wantErr: storage.ErrInvalidEntityID,
		},
		{
			name: "Deleted by key",
			key:  "key",
			eng:  storage.NewEngine(0),
			testData: pair{
				k: "key",
				v: "val",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zeroTestData := pair{}
			if tt.testData != zeroTestData {
				require.NoError(t, tt.eng.Set(tt.testData.k, tt.testData.v))
			}

			err := tt.eng.Del(tt.key)
			if err != nil && tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "Engine.Del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err, "Engine.Del() error = %v for %v key", err, tt.key)

			_, err = tt.eng.Get(tt.key)
			require.ErrorIs(t, err, storage.ErrNotFound,
				"Engine.Get() error = %v for %v key, want %v", err, tt.key, storage.ErrNotFound)
		})
	}
}
