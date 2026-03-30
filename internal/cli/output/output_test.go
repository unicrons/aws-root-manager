package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testHeaders = []string{"Name", "Status"}
	testData    = [][]any{
		{"TrustedAccess", "true"},
		{"RootSessions", "false"},
	}
)

func TestHandleOutput_JSON(t *testing.T) {
	var buf bytes.Buffer
	HandleOutput(&buf, "json", testHeaders, testData)

	out := buf.String()
	assert.Contains(t, out, `"Name"`)
	assert.Contains(t, out, `"TrustedAccess"`)
	assert.Contains(t, out, `"Status"`)
}

func TestHandleOutput_CSV(t *testing.T) {
	var buf bytes.Buffer
	HandleOutput(&buf, "csv", testHeaders, testData)

	out := buf.String()
	assert.Contains(t, out, "Name,Status")
	assert.Contains(t, out, "TrustedAccess,true")
}

func TestHandleOutput_Table(t *testing.T) {
	var buf bytes.Buffer
	HandleOutput(&buf, "table", testHeaders, testData)

	assert.NotEmpty(t, buf.String())
}

func TestHandleOutput_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	HandleOutput(&buf, "xml", testHeaders, testData)

	assert.Empty(t, buf.String())
}

func TestHandleOutput_BoolFormatting(t *testing.T) {
	var buf bytes.Buffer
	data := [][]any{
		{"TrustedAccess", true},
	}
	HandleOutput(&buf, "csv", testHeaders, data)

	assert.Contains(t, buf.String(), "true")
}
