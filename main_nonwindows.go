//go:build !windows

package main

import (
	"os"
	"time"
)

func getPlatformSpecificTime(fi os.FileInfo) (time.Time, bool) {
	return time.Time{}, false
}
