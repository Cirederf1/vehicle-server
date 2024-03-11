package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func EncodeJSON(t *testing.T, v any) io.Reader {
	t.Helper()

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(v)
	require.NoError(t, err)

	return &buf
}
