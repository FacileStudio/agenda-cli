package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintJSON(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	data := map[string]string{"key": "value"}
	err := PrintJSON(data)

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, `"key": "value"`) {
		t.Errorf("expected json containing key: value, got %s", output)
	}
}
