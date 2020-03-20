package testutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Stubin(t *testing.T, value []byte) func() {
	t.Helper()
	inr, inw, err := os.Pipe()
	require.NoError(t, err)
	orgStdin := os.Stdin
	_, err = inw.Write(value)
	require.NoError(t, err)
	require.NoError(t, inw.Close())
	os.Stdin = inr
	return func() { os.Stdin = orgStdin }
}
