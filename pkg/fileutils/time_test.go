package fileutils

import (
	"os"
	"testing"
	"time"
)

func TestGetFileTime(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "time-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Get file info
	fi, err := os.Stat(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Test file time retrieval
	fileTime := GetFileTime(fi)
	if fileTime.IsZero() {
		t.Error("Expected non-zero file time")
	}

	// Time should be recent
	if time.Since(fileTime) > time.Minute {
		t.Error("File time seems incorrect")
	}
}
