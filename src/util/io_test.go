package util

import (
	"bytes"
	"io"
	"log"
	"strings"
	"testing"
)

func TestOutputWithLogCapture(t *testing.T) {
	input := "test line\n"
	reader := io.NopCloser(strings.NewReader(input))
	var output string

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	Output(&output, reader, 123)

	if output != input {
		t.Errorf("Expected output %q, got %q", input, output)
	}

	if !strings.Contains(logBuf.String(), "123: test line") {
		t.Errorf("Expected log to contain process ID and message, got: %s", logBuf.String())
	}
}
