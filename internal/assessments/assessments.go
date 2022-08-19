package assessments

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/codio/guides-converter-v3/internal/constants"
	"github.com/codio/guides-converter-v3/internal/utils"
)

func Convert() error {
	fmt.Println("QQQQQQQQQQQq")
	time.Sleep(10 * time.Second)
	fmt.Println("wwwwwwwwwwwwwww")
	if _, err := os.Stat(constants.AssessmentsDescriptionFile); os.IsNotExist(err) {
		return nil
	}
	var assessments []interface{}
	if err := utils.GetParsedJson(constants.AssessmentsDescriptionFile, &assessments); err != nil {
		return err
	}
	if err := utils.MakeDir(constants.AssessmentsFolder); err != nil {
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
	fPath := filepath.Join("./", constants.AssessmentsFolder, fileName)
	if err := utils.WriteJson(fPath, content); err != nil {
		return err
	}
	return nil
}
