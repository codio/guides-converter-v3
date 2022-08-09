package content

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	copyF "github.com/otiai10/copy"

	"github.com/codio/guides-converter-v3/internal/utils"
)

const (
	GuidesFolder          = ".guides"
	GuidesContentFolder   = GuidesFolder + "/content"
	GuidesDescriptionFile = GuidesFolder + "/metadata.json"
	GuidesBookFile        = GuidesFolder + "/book.json"

	TmpContentFolderName = "content_v3"
	TmpContentFolder     = GuidesFolder + "/" + TmpContentFolderName
	IndexJsonFile        = "index.json"
	IndexFile            = "index"
	ContentFile          = "content-file"

	ContentHeaderFile = "header.html"
	ContentFooterFile = "footer.html"
)

var metadataSections = make(map[string]map[string]interface{})
var escapeNameRe = regexp.MustCompile("[^a-zA-Z0-9]")

func Convert() error {
	pathToTmpContent := filepath.Join("./", TmpContentFolder)
	utils.MakeDir(pathToTmpContent)

	pathToGuidesDescriptionFile := filepath.Join("./", GuidesDescriptionFile)
	pathToBookFile := filepath.Join("./", GuidesBookFile)
	var metadata map[string]interface{}
	var structure map[string]interface{}
	if err := utils.GetParsedJson(pathToGuidesDescriptionFile, &metadata); err != nil {
		return err
	}
	if err := utils.GetParsedJson(pathToBookFile, &structure); err != nil {
		return nil
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
			convertNodeToNewFormat(pathToTmpContent, &metadata, &node)
		}
	}

	if err := utils.RemoveFile(pathToGuidesDescriptionFile); err != nil {
		return err
	}
	if err := utils.RemoveFile(pathToBookFile); err != nil {
		return err
	}
	pathToContent := filepath.Join("./", GuidesContentFolder)
	if err := utils.RemoveDirectory(pathToContent); err != nil {
		return err
	}
	if err := utils.Rename(pathToTmpContent, pathToContent); err != nil {
		return err
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
	pathToRootJson := filepath.Join("./", TmpContentFolder, IndexJsonFile)
	return utils.WriteJson(pathToRootJson, newMetadata)
}

func getOrder(children interface{}) ([]string, error) {
	nodes, ok := children.([]interface{})
	if !ok {
		return nil, fmt.Errorf("get order error")
	}
	var order []string
	for _, node := range nodes {
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
	name := fmt.Sprintf("%s-%s", escapeName(title), suffix)
	return name, nil
}

func escapeName(name string) string {
	return escapeNameRe.ReplaceAllString(name, "-")
}

func copyHtmlHeaderFooter() error {
	contentFolder := filepath.Join("./", GuidesContentFolder)
	tmpContentFolder := filepath.Join("./", TmpContentFolder)
	if err := copyF.Copy(
		filepath.Join(contentFolder, ContentHeaderFile),
		filepath.Join(tmpContentFolder, ContentHeaderFile)); err != nil {
		log.Printf("error copy html header: %s\n", err)
	}
	if err := copyF.Copy(
		filepath.Join(contentFolder, ContentFooterFile),
		filepath.Join(tmpContentFolder, ContentFooterFile)); err != nil {
		log.Printf("error copy html footer: %s\n", err)
	}
	return nil
}

func convertNodeToNewFormat(parentPath string, metadataPtr, nodePtr *map[string]interface{}) error {
	createNodeMetadata(parentPath, metadataPtr, nodePtr)
	convertNodeContent(parentPath, metadataPtr, nodePtr)
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
			metadataFileName = IndexFile
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
	if title, exists := node["title"]; exists {
		nodeMetadata["title"] = title.(string)
	}
	if pageId, exists := node["pageId"]; exists {
		if section, exists := metadataSections[pageId.(string)]; exists {
			nodeMetadata = section
		} else {
			return nil, fmt.Errorf("metadata for node %s not found", pageId)
		}
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
	if contentFile, exists := nodeMetadata[ContentFile]; exists {
		extension := filepath.Ext(contentFile.(string))
		newContentPath := fmt.Sprintf("%s/%s%s", filePath, metadataFileName, extension)
		newContentPath = strings.Replace(newContentPath, TmpContentFolder, GuidesContentFolder, 1)
		newMetadata["contentFile"] = newContentPath
		delete(newMetadata, ContentFile)
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
	nodeContentPath, exists := nodeMetadata[ContentFile].(string)
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

	fileName := IndexFile + extension
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
