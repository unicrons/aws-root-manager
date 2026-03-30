package output

import (
	"bytes"
	"strings"
	"testing"
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
	if !strings.Contains(out, `"Name"`) {
		t.Errorf("expected Name key in JSON output, got: %s", out)
	}
	if !strings.Contains(out, `"TrustedAccess"`) {
		t.Errorf("expected TrustedAccess value in JSON output, got: %s", out)
	}
	if !strings.Contains(out, `"Status"`) {
		t.Errorf("expected Status key in JSON output, got: %s", out)
	}
}

func TestHandleOutput_CSV(t *testing.T) {
	var buf bytes.Buffer
	HandleOutput(&buf, "csv", testHeaders, testData)

	out := buf.String()
	if !strings.Contains(out, "Name,Status") {
		t.Errorf("expected CSV header row, got: %s", out)
	}
	if !strings.Contains(out, "TrustedAccess,true") {
		t.Errorf("expected CSV data row, got: %s", out)
	}
}

func TestHandleOutput_Table(t *testing.T) {
	var buf bytes.Buffer
	HandleOutput(&buf, "table", testHeaders, testData)

	if buf.Len() == 0 {
		t.Error("expected table output, got empty buffer")
	}
}

func TestHandleOutput_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	HandleOutput(&buf, "xml", testHeaders, testData)

	if buf.Len() != 0 {
		t.Errorf("expected no output for unknown format, got: %s", buf.String())
	}
}

func TestHandleOutput_BoolFormatting(t *testing.T) {
	var buf bytes.Buffer
	data := [][]any{
		{"TrustedAccess", true},
	}
	HandleOutput(&buf, "csv", testHeaders, data)

	out := buf.String()
	if !strings.Contains(out, "true") {
		t.Errorf("expected bool to be formatted as string in CSV, got: %s", out)
	}
}
