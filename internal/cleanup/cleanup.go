package cleanup

import (
	"github.com/codio/guides-converter-v3/internal/constants"
	"github.com/codio/guides-converter-v3/internal/utils"
)

func AfterError() error {
	// clean up assessments
	if err := utils.RemoveDirectory(constants.AssessmentsFolder); err != nil {
		return err
	}
	// clean up content
	if err := utils.RemoveDirectory(constants.TmpContentFolder); err != nil {
		return err
	}
	return nil
}

func AfterSuccessfull() error {
	// clean up assessments
	if err := utils.RemoveDirectoryIfEmpty(constants.AssessmentsFolder); err != nil {
		return err
	}
	if err := utils.RemoveFile(constants.AssessmentsDescriptionFile); err != nil {
		return err
	}
	// clean up content
	if err := utils.RemoveFile(constants.GuidesDescriptionFile); err != nil {
		return err
	}
	if err := utils.RemoveFile(constants.GuidesBookFile); err != nil {
		return err
	}
	if err := utils.RemoveDirectory(constants.GuidesContentFolder); err != nil {
		return err
	}
	if err := utils.Rename(constants.TmpContentFolder, constants.GuidesContentFolder); err != nil {
		return err
	}
	return nil
}
