package cleanup

import (
	"github.com/codio/guides-converter-v3/internal/guidespaths"
	"github.com/codio/guides-converter-v3/internal/utils"
)

func AfterError() error {
	// clean up assessments
	if err := utils.RemoveDirectory(guidespaths.GetGuidesPaths().AssessmentsFolder); err != nil {
		return err
	}
	// clean up content
	if err := utils.RemoveDirectory(guidespaths.GetGuidesPaths().TmpContentFolder); err != nil {
		return err
	}
	// clean up AlreadyInProgressFlag
	if err := utils.RemoveFile(guidespaths.GetGuidesPaths().AlreadyInProgressFlag); err != nil {
		return err
	}
	return nil
}

func AfterSuccessfull() error {
	// clean up assessments
	if err := utils.RemoveDirectoryIfEmpty(guidespaths.GetGuidesPaths().AssessmentsFolder); err != nil {
		return err
	}
	if err := utils.RemoveFile(guidespaths.GetGuidesPaths().AssessmentsDescriptionFile); err != nil {
		return err
	}
	// clean up content
	if err := utils.RemoveFile(guidespaths.GetGuidesPaths().GuidesDescriptionFile); err != nil {
		return err
	}
	if err := utils.RemoveFile(guidespaths.GetGuidesPaths().GuidesBookFile); err != nil {
		return err
	}
	if err := utils.RemoveDirectory(guidespaths.GetGuidesPaths().GuidesContentFolder); err != nil {
		return err
	}
	if err := utils.Rename(guidespaths.GetGuidesPaths().TmpContentFolder, guidespaths.GetGuidesPaths().GuidesContentFolder); err != nil {
		return err
	}
	// clean up AlreadyInProgressFlag
	if err := utils.RemoveFile(guidespaths.GetGuidesPaths().AlreadyInProgressFlag); err != nil {
		return err
	}
	return nil
}
