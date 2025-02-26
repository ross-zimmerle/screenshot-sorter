package fileutils

import (
	"os"
	"time"
)

// GetFileTime attempts to get the most appropriate timestamp for the file
func GetFileTime(fi os.FileInfo) time.Time {
	if platformTime, ok := getPlatformSpecificTime(fi); ok {
		return platformTime
	}
	return fi.ModTime()
}
