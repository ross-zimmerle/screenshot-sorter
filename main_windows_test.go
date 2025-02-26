//go:build windows

package main

import (
	"os"
	"testing"
	"time"
)

func TestGetPlatformSpecificTime(t *testing.T) {
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

	// Test platform-specific time
	creationTime, ok := getPlatformSpecificTime(fi)
	if !ok {
		t.Error("Expected to get creation time on Windows")
	}

	// Creation time should be close to now
	if time.Since(creationTime) > time.Minute {
		t.Error("Creation time seems incorrect")
	}
}
