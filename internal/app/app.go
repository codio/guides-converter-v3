package app

import (
	"fmt"
	"os"

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
	if _, err := os.Stat(constants.AlreadyInProgressFlag); err == nil {
		return true, nil
	}
	f, err := os.Create(constants.AlreadyInProgressFlag)
	if err != nil {
		return false, err
	}
	defer f.Close()
	return false, nil
}
