package compute_test

import (
	"testing"

	"github.com/alukart32/go-fast-key/internal/database/compute"
	"github.com/stretchr/testify/assert"
)

func TestNewQuery(t *testing.T) {
	gotQuery := compute.NewQuery(compute.SetCommand, []string{"key", "value"})
	expectedCommandID := compute.SetCommand
	expectedArguments := []string{"key", "value"}

	assert.EqualValues(t, expectedCommandID, gotQuery.CommandID(),
		"gotQuery.CommandID() = %v, want = %v", gotQuery.CommandID(), expectedCommandID)
	assert.EqualValues(t, expectedArguments, gotQuery.Arguments(),
		"gotQuery.Arguments() = %v, want = %v", gotQuery.Arguments(), expectedArguments)
}
