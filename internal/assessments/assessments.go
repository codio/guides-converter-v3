package assessments

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codio/guides-converter-v3/internal/utils"
)

const (
	GUIDES_FOLDER = ".guides"
	ASSESSMENTS_DESCRIPTION_FILE = GUIDES_FOLDER + "/assessments.json"
	ASSESSMENTS_FOLDER_NAME = "assessments"
	ASSESSMENTS_FOLDER = GUIDES_FOLDER + "/" + ASSESSMENTS_FOLDER_NAME
)

func Convert() error {
	workDir, err := os.Getwd()
  if err != nil {
    return err
  }
	pathToAssessmentDescription := filepath.Join(workDir, ASSESSMENTS_DESCRIPTION_FILE)
	var assessments []interface{}
	if err := utils.GetParsedJson(pathToAssessmentDescription, &assessments); err != nil {
		return err
	}
	newAssessmentsFolder := filepath.Join(workDir, ASSESSMENTS_FOLDER)
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
			createAssessmentJson(id + ".json", node)
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
	fPath :=  filepath.Join("./", ASSESSMENTS_FOLDER, fileName)
	if err := utils.WriteJson(fPath, content); err != nil {
		return err
	}
	return nil
}
