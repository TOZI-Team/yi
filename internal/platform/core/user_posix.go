//go:build !windows

package platformCore

import (
	"os"
	"runtime"
)

func isRoot() bool {
	if runtime.GOOS != "windows" && os.Getuid() == 0 {
		return true
	}
	return false
}
