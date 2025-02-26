package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/screenshot-sorter/pkg/fileutils"
)

type FileInfoDirEntry struct {
	info os.FileInfo
}

func (f FileInfoDirEntry) Name() string               { return f.info.Name() }
func (f FileInfoDirEntry) IsDir() bool                { return f.info.IsDir() }
func (f FileInfoDirEntry) Type() os.FileMode          { return f.info.Mode().Type() }
func (f FileInfoDirEntry) Info() (os.FileInfo, error) { return f.info, nil }

func TestImageProcessor_ProcessFile(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file with a specific creation/modification time
	testFileName := "test.png"
	testFilePath := filepath.Join(tempDir, testFileName)

	// Create file with content
	if err := os.WriteFile(testFilePath, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Get file info for creation time
	fileInfo, err := os.Stat(testFilePath)
	if err != nil {
		t.Fatal(err)
	}

	fileTime := fileutils.GetFileTime(fileInfo)
	expectedYear := fileTime.Format("2006")

	config := &Config{
		DryRun:    false,
		Verbose:   true,
		TargetDir: tempDir,
	}

	processor := NewImageProcessor(config)
	dirEntry := FileInfoDirEntry{info: fileInfo}
	processed, err := processor.ProcessFile(tempDir, tempDir, dirEntry)

	if err != nil {
		t.Errorf("ProcessFile() error = %v", err)
	}
	if !processed {
		t.Error("ProcessFile() file was not processed")
	}

	// Check if file was moved to the correct year folder based on file time
	expectedPath := filepath.Join(tempDir, expectedYear, testFileName)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("File was not moved to the expected location: %s", expectedPath)
	}
}

func TestImageProcessor_ProcessFileWithCustomTime(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files for different years
	testFiles := []struct {
		name string
		year int
	}{
		{"old.png", 2020},
		{"new.png", 2023},
		{"future.png", 2025},
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf.name)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}

		// Set file times (this might not work on all platforms)
		testTime := time.Date(tf.year, time.January, 1, 0, 0, 0, 0, time.UTC)
		if err := os.Chtimes(filePath, testTime, testTime); err != nil {
			t.Logf("Warning: Could not set file time for %s: %v", tf.name, err)
		}
	}

	config := &Config{
		DryRun:    false,
		Verbose:   true,
		TargetDir: tempDir,
	}

	processor := NewImageProcessor(config)

	// Process each file and verify it goes to the correct year folder
	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf.name)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			t.Fatal(err)
		}

		dirEntry := FileInfoDirEntry{info: fileInfo}
		processed, err := processor.ProcessFile(tempDir, tempDir, dirEntry)
		if err != nil {
			t.Errorf("ProcessFile() error = %v", err)
			continue
		}
		if !processed {
			t.Errorf("ProcessFile() file %s was not processed", tf.name)
			continue
		}

		// Check if file was moved to the correct year folder
		expectedPath := filepath.Join(tempDir, fileutils.GetFileTime(fileInfo).Format("2006"), tf.name)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("File %s was not moved to the expected location: %s", tf.name, expectedPath)
		}
	}
}

func TestImageProcessor_ProcessFileWithConflictingTimes(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create two files with the same name but different times
	testFile1 := filepath.Join(tempDir, "test.png")
	if err := os.WriteFile(testFile1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create year directory and existing file
	fileInfo1, err := os.Stat(testFile1)
	if err != nil {
		t.Fatal(err)
	}

	year := fileutils.GetFileTime(fileInfo1).Format("2006")
	yearDir := filepath.Join(tempDir, year)
	if err := os.MkdirAll(yearDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a file with the same name in the year directory
	existingFile := filepath.Join(yearDir, "test.png")
	if err := os.WriteFile(existingFile, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}

	config := &Config{
		DryRun:    false,
		Verbose:   true,
		TargetDir: tempDir,
	}

	processor := NewImageProcessor(config)
	dirEntry := FileInfoDirEntry{info: fileInfo1}
	processed, err := processor.ProcessFile(tempDir, tempDir, dirEntry)

	if err != nil {
		t.Errorf("ProcessFile() error = %v", err)
	}
	if !processed {
		t.Error("ProcessFile() file was not processed")
	}

	// Check that both files exist in the year directory
	files, err := os.ReadDir(yearDir)
	if err != nil {
		t.Fatal(err)
	}

	// Should have the original file and the renamed one
	if len(files) != 2 {
		t.Errorf("Expected 2 files in year directory, got %d", len(files))
	}

	// Verify the original wasn't overwritten
	content, err := os.ReadFile(existingFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "content2" {
		t.Error("Original file was overwritten")
	}
}

func TestImageProcessor_ProcessDirectoryRecursive(t *testing.T) {
	// Set test environment to use ModTime
	os.Setenv("SCREENSHOT_SORTER_TEST_USE_MODTIME", "1")
	defer os.Unsetenv("SCREENSHOT_SORTER_TEST_USE_MODTIME")

	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create nested test directories
	dirs := []string{
		filepath.Join(tempDir, "folder1"),
		filepath.Join(tempDir, "folder1", "subfolder"),
		filepath.Join(tempDir, "folder2"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create test files with different years and types
	testFiles := []struct {
		path    string
		content string
		year    int
		wanted  bool // whether it should be processed
	}{
		{filepath.Join(tempDir, "root2020.png"), "content1", 2020, true},
		{filepath.Join(tempDir, "root2023.jpg"), "content2", 2023, true},
		{filepath.Join(tempDir, "ignore.txt"), "content3", 2023, false},
		{filepath.Join(dirs[0], "folder2021.png"), "content4", 2021, true},
		{filepath.Join(dirs[1], "sub2022.png"), "content5", 2022, true},
		{filepath.Join(dirs[2], "other2024.jpg"), "content6", 2024, true},
	}

	for _, tf := range testFiles {
		if err := os.WriteFile(tf.path, []byte(tf.content), 0644); err != nil {
			t.Fatal(err)
		}

		// Set file times
		testTime := time.Date(tf.year, time.January, 1, 0, 0, 0, 0, time.UTC)
		if err := os.Chtimes(tf.path, testTime, testTime); err != nil {
			t.Logf("Warning: Could not set file time for %s: %v", tf.path, err)
		}
	}

	config := &Config{
		DryRun:    false,
		Verbose:   true,
		Recursive: true,
		TargetDir: tempDir,
		SourceDir: tempDir,
	}

	processor := NewImageProcessor(config)
	if err := processor.ProcessDirectory(tempDir, tempDir); err != nil {
		t.Errorf("ProcessDirectory() error = %v", err)
	}

	// Verify each file is in the correct year directory
	for _, tf := range testFiles {
		baseName := filepath.Base(tf.path)
		year := fmt.Sprintf("%d", tf.year)
		sourceDir := filepath.Dir(tf.path)
		expectedPath := filepath.Join(sourceDir, year, baseName)

		exists := true
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			exists = false
		}

		// For image files, they should be moved
		if tf.wanted {
			if !exists {
				t.Errorf("File %s should have been moved to %s", tf.path, expectedPath)
			}

			// Original file should not exist
			if _, err := os.Stat(tf.path); !os.IsNotExist(err) {
				t.Errorf("Original file %s should have been moved", tf.path)
			}
		} else {
			// Non-image files should stay in place
			if _, err := os.Stat(tf.path); os.IsNotExist(err) {
				t.Errorf("Non-image file %s should not have been moved", tf.path)
			}
		}
	}
}

func TestImageProcessor_FileConflictWithTimestamp(t *testing.T) {
	// Set test environment to use ModTime
	os.Setenv("SCREENSHOT_SORTER_TEST_USE_MODTIME", "1")
	defer os.Unsetenv("SCREENSHOT_SORTER_TEST_USE_MODTIME")

	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files with specific timestamps
	time1 := time.Date(2023, 5, 15, 14, 30, 0, 0, time.UTC)
	time2 := time.Date(2023, 6, 20, 16, 45, 0, 0, time.UTC)

	// Create first file and set its time
	file1 := filepath.Join(tempDir, "test.png")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(file1, time1, time1); err != nil {
		t.Logf("Warning: Could not set file time for %s: %v", file1, err)
	}

	// Process first file
	config := &Config{
		DryRun:    false,
		Verbose:   true,
		TargetDir: tempDir,
	}

	processor := NewImageProcessor(config)
	fileInfo1, _ := os.Stat(file1)
	dirEntry1 := FileInfoDirEntry{info: fileInfo1}
	if _, err := processor.ProcessFile(tempDir, tempDir, dirEntry1); err != nil {
		t.Fatal(err)
	}

	// Create second file with same name but different time
	file2 := filepath.Join(tempDir, "test.png")
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(file2, time2, time2); err != nil {
		t.Logf("Warning: Could not set file time for %s: %v", file2, err)
	}

	// Process second file
	fileInfo2, _ := os.Stat(file2)
	dirEntry2 := FileInfoDirEntry{info: fileInfo2}
	if _, err := processor.ProcessFile(tempDir, tempDir, dirEntry2); err != nil {
		t.Fatal(err)
	}

	// Check that both files exist with their proper timestamps
	yearDir := filepath.Join(tempDir, "2023")
	files, err := os.ReadDir(yearDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	// Convert expected timestamp to UTC for comparison
	timestampTime := time2.UTC()
	expectedTimestamp := timestampTime.Format("20060102_150405")
	expectedNames := map[string]bool{
		"test.png": true,
		fmt.Sprintf("test_%s.png", expectedTimestamp): true,
	}

	for _, f := range files {
		if !expectedNames[f.Name()] {
			t.Errorf("Unexpected filename: %s", f.Name())
		}
	}

	// Verify file contents are preserved
	for _, f := range files {
		content, err := os.ReadFile(filepath.Join(yearDir, f.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if len(content) == 0 {
			t.Errorf("File %s is empty", f.Name())
		}
	}
}
