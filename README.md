# Screenshot Sorter

A simple command-line tool written in Go that automatically organizes image files into year-based folders based on their creation date.

## Installation

### Prerequisites

- Go 1.20 or later

### Using Go Install
```bash
go install github.com/screenshot-sorter@latest
```

### Building from Source
1. Clone the repository
```bash
git clone https://github.com/screenshot-sorter.git
cd screenshot-sorter
```

2. Build the executable
```bash
go build
```

## Features

- 📁 Automatically sorts image files into folders by year
- 🖼️ Supports common image formats (PNG, JPG, JPEG, GIF, BMP)
- 🏷️ Preserves original filenames
- ⚡ Handles duplicate filenames automatically
- 📂 Recursive directory processing
- 🔍 Dry-run mode to preview changes
- 📌 Custom source and target directory support
- 📝 Verbose logging option
- 🚦 Rate limiting to prevent system overload (100 operations/second)
- 🎯 Platform-specific timestamp handling

## Usage

### Basic Usage
Simply run the executable and it will process images in its directory:
```bash
screenshot-sorter
```

### Advanced Options
```bash
screenshot-sorter [options]

Options:
  -source string    Source directory to process (default: executable directory)
  -target string    Target directory for sorted files (default: source directory)
  -dry-run         Show what would be done without making changes
  -recursive       Process subdirectories recursively
  -verbose         Show detailed processing information
  -version         Show version information
```

### Examples

Sort files in current directory:
```bash
screenshot-sorter
```

Preview changes without moving files:
```bash
screenshot-sorter -dry-run
```

Process files recursively with detailed output:
```bash
screenshot-sorter -recursive -verbose
```

Sort files from one directory to another:
```bash
screenshot-sorter -source ~/Downloads -target ~/Pictures
```

## Supported Image Formats

The following image formats are supported (case-insensitive):
- PNG (.png)
- JPEG (.jpg, .jpeg)
- GIF (.gif)
- BMP (.bmp)

## Notes

- File timestamps are based on:
  - Creation time on Windows
  - Modification time on other platforms
- Files are organized into year-based folders
- Operations are rate-limited to 100 per second to prevent system overload
- Duplicate filenames are handled automatically
- Non-image files are ignored

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)