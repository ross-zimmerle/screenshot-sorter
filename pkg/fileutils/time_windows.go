package fileutils

import (
	"os"
	"syscall"
	"time"
)

func getPlatformSpecificTime(fi os.FileInfo) (time.Time, bool) {
	if os.Getenv("SCREENSHOT_SORTER_TEST_USE_MODTIME") == "1" {
		return fi.ModTime(), true
	}

	if sys := fi.Sys(); sys != nil {
		if winInfo, ok := sys.(*syscall.Win32FileAttributeData); ok {
			return time.Unix(0, winInfo.CreationTime.Nanoseconds()), true
		}
	}
	return time.Time{}, false
}
