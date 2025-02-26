package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/screenshot-sorter/pkg/core"
)

// FileInfoDirEntry wrapper is already defined in main_test.go

func BenchmarkProcessFile(b *testing.B) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file
	testFile := filepath.Join(tempDir, "test.png")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		b.Fatal(err)
	}

	// Set up configuration
	config := &core.Config{
		DryRun:    false,
		Verbose:   false,
		Recursive: false,
		TargetDir: tempDir,
		SourceDir: tempDir,
	}

	processor := core.NewImageProcessor(config)

	// Get file info for testing
	fileInfo, err := os.Stat(testFile)
	if err != nil {
		b.Fatal(err)
	}

	dirEntry := FileInfoDirEntry{info: fileInfo}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.ProcessFile(tempDir, tempDir, dirEntry)
		// Recreate file for next iteration
		os.WriteFile(testFile, []byte("test content"), 0644)
	}
}

func BenchmarkProcessDirectory(b *testing.B) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	for i := 0; i < 100; i++ {
		filename := filepath.Join(tempDir, time.Now().Format("20060102_150405.999999999.png"))
		if err := os.WriteFile(filename, []byte("test content"), 0644); err != nil {
			b.Fatal(err)
		}
		time.Sleep(time.Nanosecond) // Ensure unique filenames
	}

	config := &core.Config{
		DryRun:    false,
		Verbose:   false,
		Recursive: false,
		TargetDir: tempDir,
		SourceDir: tempDir,
	}

	processor := core.NewImageProcessor(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := processor.ProcessDirectory(tempDir, tempDir); err != nil {
			b.Fatal(err)
		}

		// Recreate files for next iteration
		if i < b.N-1 { // Don't recreate files after last iteration
			for j := 0; j < 100; j++ {
				filename := filepath.Join(tempDir, time.Now().Format("20060102_150405.999999999.png"))
				if err := os.WriteFile(filename, []byte("test content"), 0644); err != nil {
					b.Fatal(err)
				}
				time.Sleep(time.Nanosecond) // Ensure unique filenames
			}
		}
	}
}
