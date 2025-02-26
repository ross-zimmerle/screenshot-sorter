package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/screenshot-sorter/pkg/core"
)

const version = "1.0.0"

func main() {
	config := parseFlags()

	if config.Version {
		fmt.Printf("Screenshot Sorter v%s\n", version)
		return
	}

	processor := core.NewImageProcessor(config)
	if err := processor.ProcessDirectory(config.SourceDir, config.TargetDir); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nScreenshot sorting complete!")
	fmt.Println("Press Enter to exit...")
	if _, err := fmt.Scanln(); err != nil && err.Error() != "unexpected newline" {
		log.Printf("Error reading input: %v", err)
	}
}

func parseFlags() *core.Config {
	config := &core.Config{}

	// Get executable directory as default
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path:", err)
	}
	defaultDir := filepath.Dir(exePath)

	flag.BoolVar(&config.DryRun, "dry-run", false, "Show what would be done without making changes")
	flag.BoolVar(&config.Verbose, "verbose", false, "Show detailed processing information")
	flag.BoolVar(&config.Recursive, "recursive", false, "Process subdirectories recursively")
	flag.StringVar(&config.TargetDir, "target", "", "Target directory for sorted files (default: source directory)")
	flag.StringVar(&config.SourceDir, "source", defaultDir, "Source directory to process (default: executable directory)")
	flag.BoolVar(&config.Version, "version", false, "Show version information")
	flag.Parse()

	// If target is not specified, use source directory
	if config.TargetDir == "" {
		config.TargetDir = config.SourceDir
	}

	return config
}
