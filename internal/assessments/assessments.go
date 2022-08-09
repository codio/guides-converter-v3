package assessments

import (
	"fmt"
	"path/filepath"

	"github.com/codio/guides-converter-v3/internal/utils"
)

const (
	GuidesFolder               = ".guides"
	AssessmentsDescriptionFile = GuidesFolder + "/assessments.json"
	AssessmentsFolderName      = "assessments"
	AssessmentsFolder          = GuidesFolder + "/" + AssessmentsFolderName
)

func Convert() error {
	pathToAssessmentDescription := filepath.Join("./", AssessmentsDescriptionFile)
	var assessments []interface{}
	if err := utils.GetParsedJson(pathToAssessmentDescription, &assessments); err != nil {
		return err
	}
	newAssessmentsFolder := filepath.Join("./", AssessmentsFolder)
	if err := utils.MakeDir(newAssessmentsFolder); err != nil {
		return err
	}
	for _, val := range assessments {
		node, ok := val.(map[string]interface{})
		if !ok {
			return fmt.Errorf("error convert assessments")
		}
		id, ok := node["taskId"].(string)
		if ok {
			createAssessmentJson(id+".json", node)
		}
	}
	if err := utils.RemoveDirectoryIfEmpty(newAssessmentsFolder); err != nil {
		return err
	}
	if err := utils.RemoveFile(pathToAssessmentDescription); err != nil {
		return err
	}
	return nil
}

func createAssessmentJson(fileName string, content map[string]interface{}) error {
	fPath := filepath.Join("./", AssessmentsFolder, fileName)
	if err := utils.WriteJson(fPath, content); err != nil {
		return err
	}
	return nil
}
