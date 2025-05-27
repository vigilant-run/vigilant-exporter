package tailer

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	READ_TIMEOUT     = 100 * time.Millisecond
	FILE_PERMISSIONS = 0644
)

func TestNewTail_ExistingFileFromStart(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.log")
	content := "line 1\nline 2\nline 3\n"
	err := os.WriteFile(testFile, []byte(content), FILE_PERMISSIONS)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := TailConfig{
		Path:        testFile,
		StartOffset: 0,
	}

	tail, err := NewTail(config)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
		return
	}
	defer tail.Stop()

	select {
	case line := <-tail.Lines:
		if line.Text != "line 1" {
			t.Errorf("expected line 'line 1', got '%s'", line.Text)
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for line")
	}
}

func TestNewTail_ExistingFileWithOffset(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "offset_test.log")
	content := "line 1\nline 2\nline 3\n"
	err := os.WriteFile(testFile, []byte(content), FILE_PERMISSIONS)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := TailConfig{
		Path:        testFile,
		StartOffset: 7,
	}

	tail, err := NewTail(config)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
		return
	}
	defer tail.Stop()

	select {
	case line := <-tail.Lines:
		if line.Text != "line 2" {
			t.Errorf("expected line 'line 2', got '%s'", line.Text)
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for line")
	}
}

func TestNewTail_NonexistentFile(t *testing.T) {
	config := TailConfig{
		Path:        "/path/that/does/not/exist/test.log",
		StartOffset: 0,
	}

	tail, err := NewTail(config)
	if tail != nil {
		t.Error("expected tail to be nil for nonexistent file")
	}

	var tailFileDoesNotExistError TailFileDoesNotExistError
	if !errors.As(err, &tailFileDoesNotExistError) {
		t.Errorf("expected TailFileDoesNotExistError but got: %T", err)
		return
	}
}

func TestNewTail_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty.log")
	err := os.WriteFile(testFile, []byte(""), FILE_PERMISSIONS)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	config := TailConfig{
		Path:        testFile,
		StartOffset: 0,
	}

	tail, err := NewTail(config)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
		return
	}
	defer tail.Stop()

	select {
	case <-tail.Lines:
		t.Error("expected no line")
	case <-time.After(READ_TIMEOUT):
		return
	}
}

func TestNewTail_InvalidPermission(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "invalid_permission.log")
	err := os.WriteFile(testFile, []byte(""), 0000)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := TailConfig{
		Path:        testFile,
		StartOffset: 0,
	}

	tail, err := NewTail(config)
	if tail != nil {
		t.Error("expected tail to be nil for invalid permission")
	}

	var tailFileInvalidPermissionError TailFileInvalidPermissionError
	if !errors.As(err, &tailFileInvalidPermissionError) {
		t.Errorf("expected TailFileInvalidPermissionError but got: %T", err)
		return
	}
}

func TestNewTail_RealTimeUpdates_AppendSingleLine(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "real_time_updates.log")
	err := os.WriteFile(testFile, []byte(""), FILE_PERMISSIONS)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tail, err := NewTail(TailConfig{
		Path:        testFile,
		StartOffset: 0,
	})
	if err != nil {
		t.Fatalf("Failed to create tail: %v", err)
	}
	defer tail.Stop()

	time.Sleep(1 * time.Second)

	err = os.WriteFile(testFile, []byte("line 1\n"), FILE_PERMISSIONS)
	if err != nil {
		t.Fatalf("Failed to append line to test file: %v", err)
	}

	select {
	case line := <-tail.Lines:
		if line.Text != "line 1" {
			t.Errorf("expected line 'line 1', got '%s'", line.Text)
		}
	case <-time.After(READ_TIMEOUT):
		t.Error("timeout waiting for line")
	}
}

func TestNewTail_RealTimeUpdates_AppendMultipleLines(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "real_time_updates.log")
	err := os.WriteFile(testFile, []byte(""), FILE_PERMISSIONS)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tail, err := NewTail(TailConfig{
		Path:        testFile,
		StartOffset: 0,
	})
	if err != nil {
		t.Fatalf("Failed to create tail: %v", err)
	}
	defer tail.Stop()

	time.Sleep(1 * time.Second)

	err = os.WriteFile(testFile, []byte("line 1\nline 2\nline 3\n"), FILE_PERMISSIONS)
	if err != nil {
		t.Fatalf("Failed to append lines to test file: %v", err)
	}

	for _, line := range []string{"line 1", "line 2", "line 3"} {
		select {
		case l := <-tail.Lines:
			if l.Text != line {
				t.Errorf("expected line '%s', got '%s'", line, l.Text)
			}
		case <-time.After(READ_TIMEOUT):
			t.Errorf("timeout waiting for line: %s", line)
		}
	}
}
