package cleanup

import (
	"github.com/codio/guides-converter-v3/internal/guidespaths"
	"github.com/codio/guides-converter-v3/internal/utils"
)

func AfterError() error {
	// clean up assessments
	if err := utils.RemoveDirectory(guidespaths.AssessmentsFolder); err != nil {
		return err
	}
	// clean up content
	if err := utils.RemoveDirectory(guidespaths.TmpContentFolder); err != nil {
		return err
	}
	// clean up AlreadyInProgressFlag
	if err := utils.RemoveFile(guidespaths.AlreadyInProgressFlag); err != nil {
		return err
	}
	return nil
}

func AfterSuccessfull() error {
	// clean up assessments
	if err := utils.RemoveDirectoryIfEmpty(guidespaths.AssessmentsFolder); err != nil {
		return err
	}
	if err := utils.RemoveFile(guidespaths.AssessmentsDescriptionFile); err != nil {
		return err
	}
	// clean up content
	if err := utils.RemoveFile(guidespaths.GuidesDescriptionFile); err != nil {
		return err
	}
	if err := utils.RemoveFile(guidespaths.GuidesBookFile); err != nil {
		return err
	}
	if err := utils.RemoveDirectory(guidespaths.GuidesContentFolder); err != nil {
		return err
	}
	if err := utils.Rename(guidespaths.TmpContentFolder, guidespaths.GuidesContentFolder); err != nil {
		return err
	}
	// clean up AlreadyInProgressFlag
	if err := utils.RemoveFile(guidespaths.AlreadyInProgressFlag); err != nil {
		return err
	}
	return nil
}
