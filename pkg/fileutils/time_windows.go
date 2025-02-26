package fileutils

import (
	"os"
	"syscall"
	"time"
)

func getPlatformSpecificTime(fi os.FileInfo) (time.Time, bool) {
	if sys := fi.Sys(); sys != nil {
		if winInfo, ok := sys.(*syscall.Win32FileAttributeData); ok {
			return time.Unix(0, winInfo.CreationTime.Nanoseconds()), true
		}
	}
	return time.Time{}, false
}
