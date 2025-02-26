//go:build !windows

package main

import (
	"os"
	"testing"
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
	_, ok := getPlatformSpecificTime(fi)
	if ok {
		t.Error("Should not get platform-specific time on non-Windows platforms")
	}
}
