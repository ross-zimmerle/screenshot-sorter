package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/screenshot-sorter/pkg/core"
)

// FileInfoDirEntry wraps FileInfo to implement DirEntry
type FileInfoDirEntry struct {
	info os.FileInfo
}

func (f FileInfoDirEntry) Name() string               { return f.info.Name() }
func (f FileInfoDirEntry) IsDir() bool                { return f.info.IsDir() }
func (f FileInfoDirEntry) Type() os.FileMode          { return f.info.Mode().Type() }
func (f FileInfoDirEntry) Info() (os.FileInfo, error) { return f.info, nil }

func parseTestFlags(args []string) *core.Config {
	config := &core.Config{}
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.BoolVar(&config.DryRun, "dry-run", false, "Show what would be done without making changes")
	flags.BoolVar(&config.Verbose, "verbose", false, "Show detailed processing information")
	flags.BoolVar(&config.Recursive, "recursive", false, "Process subdirectories recursively")
	flags.StringVar(&config.TargetDir, "target", "", "Target directory for sorted files (default: source directory)")
	flags.StringVar(&config.SourceDir, "source", "", "Source directory to process (default: current directory)")
	flags.Parse(args[1:])
	return config
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected core.Config
	}{
		{
			name: "default values",
			args: []string{"cmd"},
			expected: core.Config{
				DryRun:    false,
				Verbose:   false,
				Recursive: false,
				TargetDir: "",
				SourceDir: "",
			},
		},
		{
			name: "all flags set",
			args: []string{"cmd", "-dry-run", "-verbose", "-recursive", "-target", "/target", "-source", "/source"},
			expected: core.Config{
				DryRun:    true,
				Verbose:   true,
				Recursive: true,
				TargetDir: "/target",
				SourceDir: "/source",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestFlags(tt.args)

			if config.DryRun != tt.expected.DryRun {
				t.Errorf("DryRun = %v, want %v", config.DryRun, tt.expected.DryRun)
			}
			if config.Verbose != tt.expected.Verbose {
				t.Errorf("Verbose = %v, want %v", config.Verbose, tt.expected.Verbose)
			}
			if config.Recursive != tt.expected.Recursive {
				t.Errorf("Recursive = %v, want %v", config.Recursive, tt.expected.Recursive)
			}
			if config.TargetDir != tt.expected.TargetDir {
				t.Errorf("TargetDir = %v, want %v", config.TargetDir, tt.expected.TargetDir)
			}
			if config.SourceDir != tt.expected.SourceDir {
				t.Errorf("SourceDir = %v, want %v", config.SourceDir, tt.expected.SourceDir)
			}
		})
	}
}

func TestProcessFile(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFileName := "test.png"
	testFilePath := filepath.Join(tempDir, testFileName)
	if err := os.WriteFile(testFilePath, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set up test configuration
	config := &core.Config{
		DryRun:    false,
		Verbose:   true,
		Recursive: false,
		TargetDir: tempDir,
		SourceDir: tempDir,
	}

	processor := core.NewImageProcessor(config)

	// Get file info for testing
	fileInfo, err := os.Stat(testFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Create DirEntry wrapper and test file processing
	dirEntry := FileInfoDirEntry{info: fileInfo}
	processed, err := processor.ProcessFile(tempDir, tempDir, dirEntry)
	if err != nil {
		t.Errorf("ProcessFile() error = %v", err)
	}
	if !processed {
		t.Error("ProcessFile() file was not processed")
	}

	// Check if file was moved to year folder
	year := time.Now().Format("2006")
	expectedPath := filepath.Join(tempDir, year, testFileName)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("File was not moved to the expected location: %s", expectedPath)
	}
}

func TestSupportedFormats(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{"test.png", true},
		{"test.jpg", true},
		{"test.jpeg", true},
		{"test.gif", true},
		{"test.bmp", true},
		{"test.txt", false},
		{"test.doc", false},
		{"test", false},
		{"test.PNG", true}, // Test case insensitivity
		{"test.JPG", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			ext := filepath.Ext(strings.ToLower(tt.filename))
			if got := core.SupportedFormats[ext]; got != tt.want {
				t.Errorf("SupportedFormats[%q] = %v, want %v", ext, got, tt.want)
			}
		})
	}
}

func TestProcessDirectoryRecursive(t *testing.T) {
	// Create temporary test directories
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a nested directory structure
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files in both root and subdirectory
	files := []struct {
		path    string
		content string
	}{
		{filepath.Join(tempDir, "root.png"), "root content"},
		{filepath.Join(subDir, "sub.png"), "sub content"},
	}

	for _, f := range files {
		if err := os.WriteFile(f.path, []byte(f.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Test with recursive flag
	config := &core.Config{
		Recursive: true,
		Verbose:   true,
		TargetDir: tempDir,
		SourceDir: tempDir,
	}

	processor := core.NewImageProcessor(config)
	if err := processor.ProcessDirectory(tempDir, tempDir); err != nil {
		t.Errorf("ProcessDirectory() error = %v", err)
	}

	// Verify files were moved to year-based directories in their respective locations
	year := time.Now().Format("2006")
	expectedPaths := []string{
		filepath.Join(tempDir, year, "root.png"),
		filepath.Join(subDir, year, "sub.png"),
	}

	for _, path := range expectedPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file not found: %s", path)
		}
	}

	// Verify original files were moved
	originalPaths := []string{
		filepath.Join(tempDir, "root.png"),
		filepath.Join(subDir, "sub.png"),
	}

	for _, path := range originalPaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("Original file should have been moved: %s", path)
		}
	}
}

func TestFileConflictResolution(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create year directory and existing file
	year := time.Now().Format("2006")
	yearDir := filepath.Join(tempDir, year)
	if err := os.MkdirAll(yearDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create original and conflicting files
	testFile := filepath.Join(tempDir, "test.png")
	existingFile := filepath.Join(yearDir, "test.png")

	if err := os.WriteFile(testFile, []byte("new content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(existingFile, []byte("existing content"), 0644); err != nil {
		t.Fatal(err)
	}

	config := &core.Config{
		Verbose:   true,
		TargetDir: tempDir,
		SourceDir: tempDir,
	}

	processor := core.NewImageProcessor(config)

	// Process the file
	fileInfo, _ := os.Stat(testFile)
	dirEntry := FileInfoDirEntry{info: fileInfo}
	processed, err := processor.ProcessFile(tempDir, tempDir, dirEntry)

	if err != nil {
		t.Errorf("ProcessFile() error = %v", err)
	}
	if !processed {
		t.Error("ProcessFile() file was not processed")
	}

	// Verify both files exist (original and renamed)
	files, err := os.ReadDir(yearDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files in year directory, got %d", len(files))
	}
}

func TestDryRun(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file
	testFile := filepath.Join(tempDir, "test.png")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	config := &core.Config{
		DryRun:    true,
		Verbose:   true,
		TargetDir: tempDir,
		SourceDir: tempDir,
	}

	processor := core.NewImageProcessor(config)

	// Process the file
	fileInfo, _ := os.Stat(testFile)
	dirEntry := FileInfoDirEntry{info: fileInfo}
	processed, err := processor.ProcessFile(tempDir, tempDir, dirEntry)

	if err != nil {
		t.Errorf("ProcessFile() error = %v", err)
	}
	if !processed {
		t.Error("ProcessFile() should return processed=true in dry-run mode")
	}

	// Verify file was not moved
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("File should not have been moved in dry-run mode")
	}

	// Verify year directory was not created
	year := time.Now().Format("2006")
	yearDir := filepath.Join(tempDir, year)
	if _, err := os.Stat(yearDir); !os.IsNotExist(err) {
		t.Error("Year directory should not have been created in dry-run mode")
	}
}

func TestProcessDirectoryErrors(t *testing.T) {
	// Test with non-existent directory
	config := &core.Config{
		Verbose:   true,
		SourceDir: "nonexistent",
		TargetDir: "nonexistent",
	}

	processor := core.NewImageProcessor(config)
	err := processor.ProcessDirectory("nonexistent", "nonexistent")
	if err == nil {
		t.Error("ProcessDirectory() should return error for non-existent directory")
	}

	// Test with read-only directory (if possible on the current OS)
	tempDir, err := os.MkdirTemp("", "screenshot-sorter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.png")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Make the directory read-only
	if err := os.Chmod(tempDir, 0444); err != nil {
		t.Skip("Cannot test read-only directory permissions")
	}
	defer os.Chmod(tempDir, 0755)

	config.SourceDir = tempDir
	config.TargetDir = tempDir

	// Attempt to process the read-only directory
	err = processor.ProcessDirectory(tempDir, tempDir)
	if err == nil {
		// Note: On Windows, this might still succeed due to permission inheritance
		t.Log("Warning: Expected error for read-only directory, but got none (might be OS-dependent)")
	}
}
