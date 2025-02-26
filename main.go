package main

import (
	"flag"
	"fmt"
	"log"

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
	fmt.Scanln()
}

func parseFlags() *core.Config {
	config := &core.Config{}
	flag.BoolVar(&config.DryRun, "dry-run", false, "Show what would be done without making changes")
	flag.BoolVar(&config.Verbose, "verbose", false, "Show detailed processing information")
	flag.BoolVar(&config.Recursive, "recursive", false, "Process subdirectories recursively")
	flag.StringVar(&config.TargetDir, "target", "", "Target directory for sorted files (default: source directory)")
	flag.StringVar(&config.SourceDir, "source", "", "Source directory to process (default: current directory)")
	flag.BoolVar(&config.Version, "version", false, "Show version information")
	flag.Parse()
	return config
}
