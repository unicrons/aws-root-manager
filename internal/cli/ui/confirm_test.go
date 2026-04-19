package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests run without a TTY (stdin is a pipe), so Confirm always hits the non-TTY branch.

func TestConfirm_NonTTY_ReturnsError(t *testing.T) {
	confirmed, err := Confirm("Delete credentials?")

	require.Error(t, err)
	assert.False(t, confirmed)
	assert.True(t, strings.Contains(err.Error(), "--yes"), "error should mention --yes flag")
}
