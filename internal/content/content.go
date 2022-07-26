package content

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	copyF "github.com/otiai10/copy"

	"github.com/codio/guides-converter-v3/internal/utils"
)

const (
	GUIDES_FOLDER = ".guides"
	GUIDES_CONTENT_FOLDER = GUIDES_FOLDER + "/content"
	GUIDES_DESCRIPTION_FILE = GUIDES_FOLDER + "/metadata.json"
	GUIDES_BOOK_FILE = GUIDES_FOLDER + "/book.json"

	TMP_CONTENT_FOLDER_NAME = "content_v3"
	TMP_CONTENT_FOLDER = GUIDES_FOLDER + "/" + TMP_CONTENT_FOLDER_NAME
	INDEX_JSON_FILE = "index.json"
	INDEX_FILE = "index"

	CONTENT_HEADER_FILE = "header.html"
	CONTENT_FOOTER_FILE = "footer.html"
)

var escapeNameRe = regexp.MustCompile("[^a-zA-Z0-9]")


func Convert() error {
	workDir, err := os.Getwd()
  if err != nil {
    return err
  }
	pathToTmpContent := filepath.Join(workDir, TMP_CONTENT_FOLDER)
	utils.MakeDir(pathToTmpContent)

	pathToGuidesDescriptionFile := filepath.Join(workDir, GUIDES_DESCRIPTION_FILE)
	pathToBookFile := filepath.Join(workDir, GUIDES_BOOK_FILE)
	var metadata map[string]interface{}
	var structure map[string]interface{}
	if err := utils.GetParsedJson(pathToGuidesDescriptionFile, &metadata); err != nil {
		return err
	}
	if err := utils.GetParsedJson(pathToBookFile, &structure); err != nil {
		return nil
	}
	
	if err := createRootMetadata(workDir, metadata, structure); err != nil {
		return err
	}
	if err := copyHtmlHeaderFooter(workDir); err != nil {
		return err
	}
	if children, exists := structure["children"]; exists {
		for _, node := range(children.([]map[string]interface{})) {
			convertNodeToNewFormat(pathToTmpContent, metadata, node)
		}
	}

	if err := utils.RemoveFile(pathToGuidesDescriptionFile); err != nil {
		return err
	}
	if err := utils.RemoveFile(pathToBookFile); err != nil {
		return err
	}
	pathToContent := filepath.Join(workDir, GUIDES_CONTENT_FOLDER)
	if err := utils.RemoveDirectory(pathToContent); err != nil {
		return err
	}
	if err := utils.Rename(pathToTmpContent, pathToContent); err != nil {
		return err
	}
	return nil
}

func createRootMetadata(workDir string, metadata, structure map[string]interface{}) error {
	newMetadata := make(map[string]interface{})
	for k, v := range metadata {
		if k != "sections" {
			newMetadata[k] = v
		}
	}
	newMetadata["title"] = structure["name"]
	order, err := getOrder(structure["children"])
	newMetadata["order"] = order
	if err != nil {
		return err
	}
	pathToRootJson := filepath.Join(workDir, TMP_CONTENT_FOLDER, INDEX_JSON_FILE)
	return utils.WriteJson(pathToRootJson, newMetadata)
}

func getOrder(children interface{}) ([]string, error) {
	nodes, ok := children.([]interface{})
	if !ok {
		return nil, fmt.Errorf("get order error")
	}
	var order []string
	for _, node := range(nodes) {
		node := node.(map[string]interface{})
		nodeFileName, err := getNodeFileName(node)
		if err != nil {
			return nil, err
		}
		order = append(order, nodeFileName)
	}
	return order, nil
}

func getNodeFileName(node map[string]interface{}) (string, error) {
	title := node["title"].(string)
	id := node["id"].(string)
	suffix := id[0:4]
	name := escapeName(title) + "-" + suffix
	return name, nil
}

func escapeName(name string) string {
	return escapeNameRe.ReplaceAllString(name, "-")
}

func copyHtmlHeaderFooter(workDir string) error {
	contentFolder := filepath.Join(workDir, GUIDES_CONTENT_FOLDER)
	tmpContentFolder := filepath.Join(workDir, TMP_CONTENT_FOLDER)
	if err := copyF.Copy(
		filepath.Join(contentFolder, CONTENT_HEADER_FILE),
		filepath.Join(tmpContentFolder, CONTENT_HEADER_FILE)); err != nil {
		return err
	}
	if err := copyF.Copy(
		filepath.Join(contentFolder, CONTENT_FOOTER_FILE),
		filepath.Join(tmpContentFolder, CONTENT_FOOTER_FILE)); err != nil {
		return err
	}
	return nil
}

func convertNodeToNewFormat(parentPath string, metadata, node map[string]interface{}) error {
	createNodeMetadata(parentPath, metadata, node)
	convertNodeContent()
	nodeFileName, err := getNodeFileName(node)
	if err != nil {
		return err
	}
	nodePath := filepath.Join(parentPath, nodeFileName)
	if children, exists := node["children"]; exists {
		for _, node := range(children.([]map[string]interface{})) {
			convertNodeToNewFormat(nodePath, metadata, node)
		}
	}
	return nil
}

func createNodeMetadata(parentPath string, metadata, node map[string]interface{}) error {
	nodeFileName, err := getNodeFileName(node)
	if err != nil {
		return err
	}
	filePath := parentPath
  metadataFileName := nodeFileName
	if nType, exists := node["type"]; exists {
		if nType.(string) != "page" {
			metadataFileName = INDEX_FILE
			filePath = filepath.Join(parentPath, nodeFileName)
		}
	}
	content, err := getNodeMetadataContent(filePath, metadataFileName, metadata, node)
	if err != nil {
		return err
	}
	pathToFile := filepath.Join(filePath, metadataFileName) + ".json"
	if err := utils.WriteJson(pathToFile, content); err != nil {
		return err
	}
	return nil
}

func getNodeMetadataContent(filePath, metadataFileName string, metadata, node map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func convertNodeContent()  {
	
}
