package datasize_test

import (
	"testing"

	"github.com/alukart32/go-fast-key/internal/pkg/datasize"
	"github.com/stretchr/testify/require"
)

func TestParseSizeWithBytes(t *testing.T) {
	t.Parallel()

	size, err := datasize.Parse("20B")
	require.NoError(t, err)
	require.Equal(t, 20, size)

	size, err = datasize.Parse("20b")
	require.NoError(t, err)
	require.Equal(t, 20, size)

	size, err = datasize.Parse("20")
	require.NoError(t, err)
	require.Equal(t, 20, size)
}

func TestParseSizeWithKiloBytes(t *testing.T) {
	t.Parallel()

	size, err := datasize.Parse("20KB")
	require.NoError(t, err)
	require.Equal(t, 20*1024, size)

	size, err = datasize.Parse("20Kb")
	require.NoError(t, err)
	require.Equal(t, 20*1024, size)

	size, err = datasize.Parse("20kb")
	require.NoError(t, err)
	require.Equal(t, 20*1024, size)
}

func TestParseSizeWithMegaBytes(t *testing.T) {
	t.Parallel()

	size, err := datasize.Parse("20MB")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024, size)

	size, err = datasize.Parse("20Mb")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024, size)

	size, err = datasize.Parse("20mb")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024, size)
}

func TestParseSizeWithGigaBytes(t *testing.T) {
	t.Parallel()

	size, err := datasize.Parse("20GB")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024*1024, size)

	size, err = datasize.Parse("20Gb")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024*1024, size)

	size, err = datasize.Parse("20gb")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024*1024, size)
}

func TestParseIncorrectSize(t *testing.T) {
	t.Parallel()

	_, err := datasize.Parse("-20")
	require.Error(t, err)

	_, err = datasize.Parse("-20TB")
	require.Error(t, err)
}
