package compute_test

import (
	"reflect"
	"testing"

	"github.com/alukart32/go-fast-key/internal/fastkey/compute"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		req     string
		parser  *compute.Parser
		want    compute.Query
		wantErr error
	}{
		{
			name:    "Parser is nil",
			req:     "",
			parser:  nil,
			wantErr: compute.ErrStandByParser,
		},
		{
			name:    "Empty request",
			req:     "",
			parser:  compute.NewParser(),
			wantErr: compute.ErrEmptyRequest,
		},
		{
			name:    "UTF command symbols",
			req:     "字文下",
			parser:  compute.NewParser(),
			wantErr: compute.ErrUnknownCommand,
		},
		{
			name:    "Unknown command",
			req:     "TRUNCATE",
			parser:  compute.NewParser(),
			wantErr: compute.ErrUnknownCommand,
		},
		{
			name:    "SET command invalid args number",
			req:     "SET key",
			parser:  compute.NewParser(),
			wantErr: compute.ErrInvalidArgsNumber,
		},
		{
			name:    "GET command invalid args number",
			req:     "GET key1 key2 key3",
			parser:  compute.NewParser(),
			wantErr: compute.ErrInvalidArgsNumber,
		},
		{
			name:    "DEL command invalid args number",
			req:     "DEL key1 key2 key3",
			parser:  compute.NewParser(),
			wantErr: compute.ErrInvalidArgsNumber,
		},
		{
			name:   "Valid SET request",
			req:    "SET key val",
			parser: compute.NewParser(),
			want:   compute.NewQuery(compute.SetCommand, []string{"key", "val"}),
		},
		{
			name:   "Valid GET request",
			req:    "GET key",
			parser: compute.NewParser(),
			want:   compute.NewQuery(compute.GetCommand, []string{"key"}),
		},
		{
			name:   "Valid DEL request",
			req:    "DEL key",
			parser: compute.NewParser(),
			want:   compute.NewQuery(compute.DelCommand, []string{"key"}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.Parse(tt.req)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
