package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Error("NewLogger returned nil")
	}
}

func TestLogger_Info(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalOutput)

	logger := NewLogger()
	logger.Info("test info message")

	output := buf.String()
	if !strings.Contains(output, "test info message") {
		t.Errorf("Expected log to contain 'test info message', got: %s", output)
	}
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected log to contain 'INFO', got: %s", output)
	}
}

func TestLogger_Error(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalOutput)

	logger := NewLogger()
	logger.Error("test error message")

	output := buf.String()
	if !strings.Contains(output, "test error message") {
		t.Errorf("Expected log to contain 'test error message', got: %s", output)
	}
	if !strings.Contains(output, "ERROR") {
		t.Errorf("Expected log to contain 'ERROR', got: %s", output)
	}
}

func TestLogger_Fatal(t *testing.T) {
	// We can't easily test Fatal since it calls os.Exit(1)
	// But we can test that the method exists and is callable
	logger := NewLogger()
	if logger == nil {
		t.Error("NewLogger returned nil")
	}

	// Just verify the method exists - we can't actually call it in tests
	// since it would exit the test process
	_ = logger.Fatal
}

func TestLogger_WithMultipleMessages(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalOutput)

	logger := NewLogger()
	logger.Info("first message")
	logger.Error("second message")

	output := buf.String()

	// Check that both messages are present
	if !strings.Contains(output, "first message") {
		t.Errorf("Expected log to contain 'first message', got: %s", output)
	}
	if !strings.Contains(output, "second message") {
		t.Errorf("Expected log to contain 'second message', got: %s", output)
	}

	// Check that both log levels are present
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected log to contain 'INFO', got: %s", output)
	}
	if !strings.Contains(output, "ERROR") {
		t.Errorf("Expected log to contain 'ERROR', got: %s", output)
	}
}

func TestLogger_EmptyMessage(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalOutput)

	logger := NewLogger()
	logger.Info("")

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected log to contain 'INFO' even with empty message, got: %s", output)
	}
}

// Test that logger respects original log settings
func TestLogger_PreservesOriginalSettings(t *testing.T) {
	// Save original settings
	originalFlags := log.Flags()
	originalPrefix := log.Prefix()
	originalOutput := log.Writer()

	// Create logger
	logger := NewLogger()
	logger.Info("test")

	// Check that original settings are preserved
	if log.Flags() != originalFlags {
		t.Errorf("Logger changed log flags. Expected %d, got %d", originalFlags, log.Flags())
	}
	if log.Prefix() != originalPrefix {
		t.Errorf("Logger changed log prefix. Expected '%s', got '%s'", originalPrefix, log.Prefix())
	}
	if log.Writer() != originalOutput {
		t.Errorf("Logger changed log output")
	}
}