# Screenshot Sorter

A simple and efficient tool to organize your screenshots and images into year-based folders.

## Features

- Automatically sorts images into year-based folders
- Supports multiple image formats (PNG, JPG, JPEG, GIF, BMP)
- Dry-run mode to preview changes
- Recursive directory processing
- Rate limiting to prevent system overload
- Handles file naming conflicts automatically

## Usage

```bash
screenshot-sorter [options]

Options:
  -source string    Source directory to process (default: current directory)
  -target string    Target directory for sorted files (default: source directory)
  -dry-run         Show what would be done without making changes
  -recursive       Process subdirectories recursively
  -verbose         Show detailed processing information
  -version         Show version information
```

## Examples

Sort images in the current directory:
```bash
screenshot-sorter
```

Sort images with detailed output:
```bash
screenshot-sorter -verbose
```

Preview changes without moving files:
```bash
screenshot-sorter -dry-run
```

## Installation

Download the latest release from the [releases page](https://github.com/screenshot-sorter/releases) or build from source:

```bash
go install github.com/screenshot-sorter@latest
```

## License

This project is licensed under the terms of the MIT license. See the [LICENSE](../LICENSE) file for details.