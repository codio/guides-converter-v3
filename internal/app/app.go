package app

import (
	"github.com/codio/guides-converter-v3/internal/assessments"
	"github.com/codio/guides-converter-v3/internal/cleanup"
	"github.com/codio/guides-converter-v3/internal/content"
)

func Run() error {
	if err := assessments.Convert(); err != nil {
		cleanup.AfterError()
		return err
	}
	if err := content.Convert(); err != nil {
		cleanup.AfterError()
		return err
	}
	if err := cleanup.AfterSuccessfull(); err != nil {
		return err
	}
	return nil
}
