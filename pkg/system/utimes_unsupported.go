//go:build !linux && !freebsd
// +build !linux,!freebsd

package system

import (
	"syscall"

	errsys "github.com/bhojpur/errors/pkg/system"
)

// LUtimesNano is only supported on linux and freebsd.
func LUtimesNano(path string, ts []syscall.Timespec) error {
	return errsys.ErrNotSupportedPlatform
}
