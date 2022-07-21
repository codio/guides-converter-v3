package assessments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codio/guides-converter-v3/internal/utils"
)

const (
	GUIDES_FOLDER = ".guides"
	ASSESSMENTS_DESCRIPTION_FILE = GUIDES_FOLDER + "/assessments.json"
	ASSESSMENTS_FOLDER_NAME = "assessments"
	ASSESSMENTS_FOLDER = GUIDES_FOLDER + "/" + ASSESSMENTS_FOLDER_NAME

	DEFAULT_CONTENT = "[]"
)

func Convert() error {
	workDir, err := os.Getwd()
  if err != nil {
    return err
  }
	pathToAssessmentDescription := filepath.Join(workDir, ASSESSMENTS_DESCRIPTION_FILE)
	var assessments []interface{}
	if err := getAssessments(pathToAssessmentDescription, &assessments); err != nil {
		return err
	}
	newAssessmentsFolder := filepath.Join(workDir, ASSESSMENTS_FOLDER)
	if err := makeDir(newAssessmentsFolder); err != nil {
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

func makeDir(directory string) error {
  if stat, err := os.Stat(directory); err == nil && stat.IsDir() {
    return nil
  }
  if err := os.Mkdir(directory, 0777); err != nil {
		return err
	}
	return nil
}

func createAssessmentJson(fileName string, content map[string]interface{}) error {
	fPath :=  filepath.Join("./", ASSESSMENTS_FOLDER, fileName)
	jsonFile, err := os.OpenFile(fPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	if err := jsonFile.Truncate(0); err != nil {
		return err
	}
	if _, err := jsonFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err := jsonFile.Write(data); err != nil {
		return err
	}
	return nil
}

func getAssessments(path string, assessments *[]interface{}) error {
	jsonFile, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, assessments); err != nil {
    	return err
	}
	return nil
}
