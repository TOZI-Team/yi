//go:build !windows

package platformCore

import (
	"os"
)

func isRoot() bool {
	if os.Getuid() == 0 {
		return true
	}
	return false
}
