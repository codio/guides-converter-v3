// +build windows

package app

import (
	"github.com/juju/fslock"

	"github.com/codio/guides-converter-v3/internal/constants"
)

func alreadyInProgress() (bool, error) {
	lock := fslock.New(constants.AlreadyInProgressFlag)
	if err := lock.Lock(); err != nil {
		return true, nil
	}
	return false, nil
}
