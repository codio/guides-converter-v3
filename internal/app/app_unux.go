//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris

package app

import (
	"os"
	"syscall"

	"github.com/codio/guides-converter-v3/internal/guidespaths"
)

func alreadyInProgress() (bool, error) {
	f, err := os.OpenFile(guidespaths.GetGuidesPaths().AlreadyInProgressFlag, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return true, err
	}
	fileDescriptor := int(f.Fd())
	if err := syscall.Flock(fileDescriptor, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return true, nil
	}
	return false, nil
}
