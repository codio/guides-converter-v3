package content

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	copyF "github.com/otiai10/copy"

	"github.com/codio/guides-converter-v3/internal/constants"
	"github.com/codio/guides-converter-v3/internal/utils"
)

var metadataSections = make(map[string]map[string]interface{})
var escapeNameRe = regexp.MustCompile("[^a-zA-Z0-9]")

func Convert() error {
	utils.MakeDir(constants.TmpContentFolder)
	var metadata map[string]interface{}
	var structure map[string]interface{}
	if err := utils.GetParsedJson(constants.GuidesDescriptionFile, &metadata); err != nil {
		return err
	}
	if err := utils.GetParsedJson(constants.GuidesBookFile, &structure); err != nil {
		return err
	}

	if sections, exists := metadata["sections"]; exists {
		for _, item := range sections.([]interface{}) {
			section := item.(map[string]interface{})
			if id, exists := section["id"]; exists {
				metadataSections[id.(string)] = section
			}
		}
	}

	if err := createRootMetadata(metadata, structure); err != nil {
		return err
	}
	if err := copyHtmlHeaderFooter(); err != nil {
		return err
	}
	if children, exists := structure["children"]; exists {
		for _, item := range children.([]interface{}) {
			node := item.((map[string]interface{}))
			convertNodeToNewFormat(constants.TmpContentFolder, &metadata, &node)
		}
	}
	return nil
}

func createRootMetadata(metadata, structure map[string]interface{}) error {
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
	pathToRootJson := filepath.Join(constants.TmpContentFolder, constants.IndexJsonFile)
	return utils.WriteJson(pathToRootJson, newMetadata)
}

func getOrder(children interface{}) ([]string, error) {
	nodes, ok := children.([]interface{})
	if !ok {
		return nil, fmt.Errorf("get order error")
	}
	order := make([]string, 0)
	for _, node := range nodes {
		node := node.(map[string]interface{})
		nodeIsNotExists, err := nodeNotExists(node)
		if err != nil {
			return nil, err
		}
		if (nodeIsNotExists) {
			continue
		}
		nodeFileName, err := getNodeFileName(node)
		if err != nil {
			return nil, err
		}
		order = append(order, nodeFileName)
	}
	return order, nil
}

func nodeNotExists(node map[string]interface{}) (bool, error) {
	if (node["type"].(string) != "page") {
		return false, nil
	}
	pageId := node["pageId"].(string)
	section, sectionExists := metadataSections[pageId]
	if (!sectionExists) {
		return true, nil
	}
	if contentFilePath, exists := section[constants.ContentFile]; exists {
		fullFilePath := filepath.Join(constants.WorkSpace, contentFilePath.(string))
		fileExists, err := utils.FileIsExists(fullFilePath)
		if err != nil {
			return true, err
		}
		return !fileExists, nil
	}
	return true, nil
}

func getNodeFileName(node map[string]interface{}) (string, error) {
	title := node["title"].(string)
	id := node["id"].(string)
	suffix := id[0:4]
	name := fmt.Sprintf("%s-%s", escapeName(title), suffix)
	return name, nil
}

func escapeName(name string) string {
	return escapeNameRe.ReplaceAllString(name, "-")
}

func copyHtmlHeaderFooter() error {
	if err := copyF.Copy(
		filepath.Join(constants.GuidesContentFolder, constants.ContentHeaderFile),
		filepath.Join(constants.TmpContentFolder, constants.ContentHeaderFile)); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := copyF.Copy(
		filepath.Join(constants.GuidesContentFolder, constants.ContentFooterFile),
		filepath.Join(constants.TmpContentFolder, constants.ContentFooterFile)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func convertNodeToNewFormat(parentPath string, metadataPtr, nodePtr *map[string]interface{}) error {
	err := convertNodeContent(parentPath, metadataPtr, nodePtr)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	err = createNodeMetadata(parentPath, metadataPtr, nodePtr)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	nodeFileName, err := getNodeFileName(*nodePtr)
	if err != nil {
		return err
	}
	nodePath := filepath.Join(parentPath, nodeFileName)
	node := *nodePtr
	if children, exists := node["children"]; exists {
		for _, item := range children.([]interface{}) {
			node := item.(map[string]interface{})
			convertNodeToNewFormat(nodePath, metadataPtr, &node)
		}
	}
	return nil
}

func createNodeMetadata(parentPath string, metadataPtr, nodePtr *map[string]interface{}) error {
	node := *nodePtr
	nodeFileName, err := getNodeFileName(node)
	if err != nil {
		return err
	}
	filePath := parentPath
	metadataFileName := nodeFileName
	if nType, exists := node["type"]; exists {
		if nType.(string) != "page" {
			metadataFileName = constants.IndexFile
			filePath = filepath.Join(parentPath, nodeFileName)
		}
	}
	content, err := getNodeMetadataContent(filePath, metadataFileName, metadataPtr, nodePtr)
	if err != nil {
		return err
	}
	utils.MakeDir(filePath)
	pathToFile := filepath.Join(filePath, metadataFileName) + ".json"
	if err := utils.WriteJson(pathToFile, content); err != nil {
		return err
	}
	return nil
}

func getNodeMetadataContent(filePath, metadataFileName string, metadataPtr, nodePtr *map[string]interface{}) (map[string]interface{}, error) {
	var nodeMetadata = make(map[string]interface{})
	node := *nodePtr
	if pageId, exists := node["pageId"]; exists {
		if section, exists := metadataSections[pageId.(string)]; exists {
			nodeMetadata = section
		} else {
			return nil, fmt.Errorf("metadata for node %s not found", pageId)
		}
	}
	if title, exists := node["title"]; exists {
		nodeMetadata["title"] = title.(string)
	}
	if contentType, exists := nodeMetadata["type"]; exists {
		nodeMetadata["contentType"] = contentType.(string)
	}
	newMetadata := make(map[string]interface{})
	for k, v := range nodeMetadata {
		if k != "sections" {
			newMetadata[k] = v
		}
	}
	if typeN, exists := node["type"]; exists {
		newMetadata["type"] = typeN.(string)
	}
	if id, exists := node["id"]; exists {
		newMetadata["id"] = id.(string)
	}
	if _, exists := nodeMetadata[constants.ContentFile]; exists {
		delete(newMetadata, constants.ContentFile)
	}
	if children, exists := node["children"]; exists {
		order, err := getOrder(children)
		if err != nil {
			return nil, err
		}
		newMetadata["order"] = order
	}
	return newMetadata, nil
}

func convertNodeContent(parentPath string, metadataPtr, nodePtr *map[string]interface{}) error {
	var nodeMetadata = make(map[string]interface{})
	node := *nodePtr
	pageId, exists := node["pageId"].(string)
	if !exists {
		return nil
	}
	if section, exists := metadataSections[pageId]; exists {
		nodeMetadata = section
	} else {
		return fmt.Errorf("metadata for node %s not found", pageId)
	}
	nodeContentPath, exists := nodeMetadata[constants.ContentFile].(string)
	if !exists {
		return fmt.Errorf("content for node %s not found", pageId)
	}
	if _, err := os.Stat(filepath.Join("./", nodeContentPath)); os.IsNotExist(err) {
		return fmt.Errorf("content for node %s not found", pageId)
	}
	extension := filepath.Ext(nodeContentPath)
	nodeFileName, err := getNodeFileName(node)
	if err != nil {
		return err
	}

	fileName := constants.IndexFile + extension
	if typeN, exists := node["type"]; exists && typeN.(string) == "page" {
		fileName = nodeFileName + extension
		if err := copyF.Copy(
			filepath.Join("./", nodeContentPath),
			filepath.Join(parentPath, fileName)); err != nil {
			log.Printf("error copy node content: %s\n", err)
		}
		return nil
	}

	if err := copyF.Copy(
		filepath.Join("./", nodeContentPath),
		filepath.Join(parentPath, nodeFileName, fileName)); err != nil {
		log.Printf("error copy node content: %s\n", err)
	}
	return nil
}
