package assessments

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codio/guides-converter-v3/internal/guidespaths"
	"github.com/codio/guides-converter-v3/internal/utils"
)

func Convert() error {
	if _, err := os.Stat(guidespaths.GetGuidesPaths().AssessmentsDescriptionFile); os.IsNotExist(err) {
		return nil
	}
	var assessments []interface{}
	if err := utils.GetParsedJson(guidespaths.GetGuidesPaths().AssessmentsDescriptionFile, &assessments); err != nil {
		return err
	}
	if err := utils.MakeDir(guidespaths.GetGuidesPaths().AssessmentsFolder); err != nil {
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
	return nil
}

func createAssessmentJson(fileName string, content map[string]interface{}) error {
	fPath := filepath.Join(guidespaths.GetGuidesPaths().AssessmentsFolder, fileName)
	if err := utils.WriteJson(fPath, content); err != nil {
		return err
	}
	return nil
}
