package app

import (
	"fmt"
	"os"
	"syscall"

	"github.com/codio/guides-converter-v3/internal/assessments"
	"github.com/codio/guides-converter-v3/internal/cleanup"
	"github.com/codio/guides-converter-v3/internal/constants"
	"github.com/codio/guides-converter-v3/internal/content"
)

func Run() error {
	inProgress, err := alreadyInProgress()
	if err != nil {
		return fmt.Errorf("alreadyInProgress error")
	}
	if inProgress {
		return nil
	}
	if err := assessments.Convert(); err != nil {
		cleanup.AfterError()
		return fmt.Errorf("assessments convert error")
	}
	if err := content.Convert(); err != nil {
		cleanup.AfterError()
		return fmt.Errorf("content convert error")
	}
	if err := cleanup.AfterSuccessfull(); err != nil {
		return fmt.Errorf("cleanup error")
	}
	return nil
}

func alreadyInProgress() (bool, error) {
	f, err := os.OpenFile(constants.AlreadyInProgressFlag, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return true, err
	}
	fileDescriptor := int(f.Fd())
	if err := syscall.Flock(fileDescriptor, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return true, nil
	}
	return false, nil
}
