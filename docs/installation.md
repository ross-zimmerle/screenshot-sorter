# Installation Guide

## Prerequisites

- Go 1.24 or higher (for building from source)

## Installation Methods

### Pre-built Binaries (v1.0.0)

1. Visit the [releases page](https://github.com/screenshot-sorter/releases)
2. Download the appropriate version for your operating system
3. Extract the archive
4. Add the executable to your system's PATH (optional)

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/screenshot-sorter.git
   cd screenshot-sorter
   ```

2. Build the project:
   ```bash
   go build
   ```

3. Install globally (optional):
   ```bash
   go install
   ```

## Verifying Installation

To verify the installation and check your version:
```bash
screenshot-sorter -version
```

The command should output: `Screenshot Sorter v1.0.0`