package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/screenshot-sorter/pkg/fileutils"
	"golang.org/x/time/rate"
)

// ImageProcessor handles the core image processing functionality
type ImageProcessor struct {
	limiter *rate.Limiter
	config  *Config
}

// Config holds the program configuration
type Config struct {
	DryRun    bool
	Verbose   bool
	Recursive bool
	TargetDir string
	SourceDir string
	Version   bool
}

// SupportedFormats defines the image file extensions that the program will process
var SupportedFormats = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".bmp":  true,
}

// NewImageProcessor creates a new image processor instance
func NewImageProcessor(config *Config) *ImageProcessor {
	return &ImageProcessor{
		limiter: rate.NewLimiter(100, 1),
		config:  config,
	}
}

// ProcessDirectory handles the processing of a directory
func (p *ImageProcessor) ProcessDirectory(sourceDir, targetDir string) error {
	// Check if directory exists
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", sourceDir, err)
	}

	// If targetDir is empty, use sourceDir
	if targetDir == "" {
		targetDir = sourceDir
	}

	for _, entry := range entries {
		// Rate limit operations
		p.limiter.Wait(context.Background())

		fullPath := filepath.Join(sourceDir, entry.Name())

		if entry.IsDir() {
			if p.config.Recursive {
				targetSubDir := filepath.Join(targetDir, entry.Name())
				if err := p.ProcessDirectory(fullPath, targetSubDir); err != nil {
					if p.config.Verbose {
						fmt.Printf("Error processing directory %s: %v\n", fullPath, err)
					}
				}
			}
			continue
		}

		if _, err := p.ProcessFile(sourceDir, targetDir, entry); err != nil {
			if p.config.Verbose {
				fmt.Printf("Error processing file %s: %v\n", fullPath, err)
			}
		}
	}

	return nil
}

// ProcessFile handles the processing of a single file
func (p *ImageProcessor) ProcessFile(sourceDir, targetDir string, entry os.DirEntry) (bool, error) {
	// Check if it's a supported image format
	ext := strings.ToLower(filepath.Ext(entry.Name()))
	if !SupportedFormats[ext] {
		return false, nil
	}

	// Get file info for timestamp
	fileInfo, err := entry.Info()
	if err != nil {
		return false, fmt.Errorf("failed to get file info for %s: %w", entry.Name(), err)
	}

	// Get source and target paths
	sourcePath := filepath.Join(sourceDir, entry.Name())

	// Get file's actual timestamp
	fileTime := fileutils.GetFileTime(fileInfo)
	year := fileTime.Format("2006")
	yearDir := filepath.Join(targetDir, year)

	if !p.config.DryRun {
		if err := os.MkdirAll(yearDir, 0755); err != nil {
			return false, fmt.Errorf("failed to create directory %s: %w", yearDir, err)
		}
	}

	// Generate target path and handle conflicts
	targetPath := filepath.Join(yearDir, entry.Name())
	if _, err := os.Stat(targetPath); err == nil {
		// File exists, append timestamp from the original file
		ext := filepath.Ext(entry.Name())
		base := strings.TrimSuffix(entry.Name(), ext)
		timestamp := fileTime.Format("20060102_150405")
		targetPath = filepath.Join(yearDir, fmt.Sprintf("%s_%s%s", base, timestamp, ext))
	}

	if p.config.Verbose {
		fmt.Printf("Moving %s to %s\n", sourcePath, targetPath)
	}

	if !p.config.DryRun {
		if err := os.Rename(sourcePath, targetPath); err != nil {
			return false, fmt.Errorf("failed to move file %s to %s: %w", sourcePath, targetPath, err)
		}
	}

	return true, nil
}
