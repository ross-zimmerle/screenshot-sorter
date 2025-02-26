# Advanced Usage

## Batch Processing

For batch processing multiple directories, you can use the provided scripts in the `examples` folder:

- Windows: `batch_process.ps1`
- Linux/MacOS: `batch_process.sh`

## Rate Limiting

The tool includes built-in rate limiting (100 operations per second) to prevent system overload. This is particularly useful when processing large directories.

## Directory Structure

When organizing files, the tool creates the following structure:

```
target_directory/
├── 2021/
│   ├── screenshot1.png
│   └── screenshot2.jpg
├── 2022/
│   ├── screenshot3.png
│   └── screenshot4.jpg
└── 2023/
    ├── screenshot5.png
    └── screenshot6.jpg
```

## Handling Duplicates

When a file with the same name exists in the destination folder, the tool automatically creates a unique filename by appending a timestamp:

- Original: `screenshot.png`
- Duplicate: `screenshot_20240315_143022.png`

## Command Examples

### Process Multiple Source Directories
```bash
screenshot-sorter -source "C:\Pictures" -target "D:\Sorted Pictures" -recursive
```

### Verbose Output with Dry Run
```bash
screenshot-sorter -verbose -dry-run -recursive
```

### Sorting Specific File Types
The tool automatically processes these image formats:
- PNG (.png)
- JPEG (.jpg, .jpeg)
- GIF (.gif)
- BMP (.bmp)